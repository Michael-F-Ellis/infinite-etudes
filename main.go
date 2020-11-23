// Copyright 2019 Ellis & Grant, Inc. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
/*
etudes generates a set of 6 midi files for each of 12 key signature. Each set
covers all possible combinations of 3 pitches within the key.

Command line usage is

   etudes [-h] [-t tempo] [-l midilow] [-u midihi ] [-i instrument]

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const copyright = `
Copyright 2019 Ellis & Grant, Inc. All rights reserved.  Use of the source
code is governed by an MIT-style license that can be found in the LICENSE
file.
`
const description = `
infinite-etudes generates ear training exercises for instrumentalists.
Infinite-etudes is a high-performance self-contained web server
that provides a simple user interface that allows the user to choose a key, a
scale pattern and an instrument sound and play a freshly-generated etude in
the web browser. A publically available instance is running at 

https://etudes.ellisandgrant.com

See the file server.go for details including environment variables needed
for https service.
`

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// userHomeDir returns the user's home directory name on Windows, Linux or Mac.
// Credit: https://stackoverflow.com/a/7922977/426853
func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
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
	flag.StringVar(&imgPath, "g", filepath.Join(userHomeDir(), "go", "src", "github.com", "Michael-F-Ellis", "infinite-etudes", "img"), "Path to img files on your host (server-mode only)")

	var midijsPath string
	flag.StringVar(&midijsPath, "m", filepath.Join(userHomeDir(), "go", "src", "github.com", "Michael-F-Ellis", "infinite-etudes", "midijs"), "Path to midijs files on your host (server-mode only)")

	var hostport string
	flag.StringVar(&hostport, "p", "localhost:8080", "hostname (or IP) and port to serve on. (server-mode only)")

	flag.IntVar(&expireSeconds, "x", 10, "Maximum age in seconds for generated files (server-mode only)")

	// make sure all flags are defined before calling this
	flag.Parse()

	serveEtudes(hostport, midijsPath, imgPath)

}

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

// usage extends the flag package's default help message.
func usage() {
	fmt.Println(copyright)
	fmt.Printf("Usage: etudes [OPTIONS]\n  -h    print this help message.\n")
	flag.PrintDefaults()
	fmt.Println(description)

}

// mkRequestedEtude creates the requested etude in the current directory. The
// arguments are assumed to be previously vetted and are not checked.
func mkRequestedEtude(midilo, midihi, tempo, instrument int, r etudeRequest) {
	iname := r.instrument
	switch r.pattern {
	case "allintervals":
		s := generateIntervalSequence(midilo, midihi, tempo, instrument, r)
		mkMidi(&s, true)
	case "interval":
		s := generateEqualIntervalSequence(midilo, midihi, tempo, instrument, r)
		mkMidi(&s, true)
	case "intervalpair":
		i1 := intervalSizeByName(r.interval1)
		i2 := intervalSizeByName(r.interval2)
		s := generateTwoIntervalSequence(midilo, midihi, tempo, instrument, iname, i1, i2)
		s.req = r
		mkMidi(&s, true) // no tighten
	case "intervaltriple":
		i1 := intervalSizeByName(r.interval1)
		i2 := intervalSizeByName(r.interval2)
		i3 := intervalSizeByName(r.interval3)
		s := generateThreeIntervalSequence(midilo, midihi, tempo, instrument, iname, i1, i2, i3)
		s.req = r
		mkMidi(&s, true) // no tighten
	default:
		panic(fmt.Sprintf("%s is not a supported etude pattern", r.pattern))
	}
}

// iToBools converts the first length bits of v to
// a slice of bool, e.g. iToBools(4,3) -> [true, false, false]
func iToBools(v, length int) (b []bool) {
	for i := length - 1; i >= 0; i-- {
		b = append(b, (v&(1<<i) > 0))
	}
	return
}
