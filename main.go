// Copyright 2019 Ellis & Grant, Inc. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const copyright = `
Copyright 2019-2020 Ellis & Grant, Inc. All rights reserved.  Use of the source
code is governed by an MIT-style license that can be found in the LICENSE
file.
`
const description = `
Infinite-etudes generates ear training exercises for instrumentalists. The
program contains a high-performance self-contained web server that provides a
simple user interface that allows the user to choose a pattern of intervals,
an instrument sound, and tempo to generate and play a freshly-generated etude
in the web browser. A public instance is running at

https://etudes.ellisandgrant.com

See the file server.go for details including environment variables needed
for https service.
`

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var expireSeconds int // max age for generated etude files

func main() {
	// initialize standard logger to write to "etudes.log"
	logf, err := os.OpenFile("etudes.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logf.Close()
	// Start logging to new file.
	log.SetOutput(logf)

	// Parse command line
	flag.Usage = usage

	// Command mode flags

	var imgPath string
	// flag.StringVar(&imgPath, "g", filepath.Join(userHomeDir(), "go", "src", "github.com", "Michael-F-Ellis", "infinite-etudes", "img"), "Path to img files on your host (server-mode only)")
	flag.StringVar(&imgPath, "g", "assets/img", "Path to img files on your host (server-mode only)")

	var midijsPath string
	// flag.StringVar(&midijsPath, "m", filepath.Join(userHomeDir(), "go", "src", "github.com", "Michael-F-Ellis", "infinite-etudes", "midijs"), "Path to midijs files on your host (server-mode only)")
	flag.StringVar(&midijsPath, "m", "assets/midijs", "Path to midijs files on your host (server-mode only)")

	var hostport string
	flag.StringVar(&hostport, "p", ":8080", "hostname (or IP) and port to serve on. (server-mode only)")

	flag.IntVar(&expireSeconds, "x", 10, "Maximum age in seconds for generated files (server-mode only)")

	// make sure all flags are defined before calling this
	flag.Parse()

	serveEtudes(hostport)

}

/* OUT: No longer needed
// validDirPath returns a non-nil error if path is not a directory on the host.
func validDirPath(path string) (err error) {
	finfo, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("invalid path %s: %v", path, err)
		return
	}
	if !finfo.Mode().IsDir() {
		err = fmt.Errorf("%s is not a directory", path)
	}
	return
}
*/
// usage extends the flag package's default help message.
func usage() {
	fmt.Println(copyright)
	fmt.Printf("Usage: etudes [OPTIONS]\n  -h    print this help message.\n")
	flag.PrintDefaults()
	fmt.Println(description)

}
