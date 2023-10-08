// Copyright (c) 2014-2020 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"net/http"
	_ "net/http/pprof"
	"runtime/trace"

	"github.com/canonical/pebble/internals/cli"
)

func main() {
	traceFile := fmt.Sprintf("/var/snap/maas/common/pebble/%d.trace", time.Now().UnixNano())
	f, err := os.Create(traceFile)
	if err != nil {
		log.Fatal(err)
	}
	trace.Start(f)
	timer := time.NewTimer(5 * time.Minute)
	go func() {
		<-timer.C
		trace.Stop()
		log.Println("trace file", traceFile)
	}()

	defer trace.Stop()
	go func() {
		fmt.Println(http.ListenAndServe(":6060", nil))
	}()

	if err := cli.Run(); err != nil {
		fmt.Fprintf(cli.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
