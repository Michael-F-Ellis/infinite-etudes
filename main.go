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

You can run it from the command line (cli mode) or as a web server (server mode).

In cli mode, infinite-etudes generates a set of 7 midi files for each of 12 key
signatures. Each set covers all possible combinations of 3 pitches within the
key. The files are generated in the current working directory.

In server mode, infinite-etudes is a high-performance self-contained web server
that provides a simple user interface that allows the user to choose a key, a
scale pattern and an instrument sound and play a freshly-generated etude in
the web browser. A publically available instance is running at 

https://etudes.ellisandgrant.com

See the file server.go for details including environment variables needed
for https service.

The midi file names structure is '<key>_<scalepattern>_<instrument>.mid'. For example,
	
	eflat_pentatonic_trumpet.mid
	eflat_final_trumpet.mid
	eflat_plus_four_trumpet.mid
	eflat_plus_seven_trumpet.mid
	eflat_four_and_seven_trumpet.mid
	eflat_raised_five_trumpet.mid
	eflat_raised_five_with_four_or_seven_trumpet.mid

The 12 keynames used are:

	a, b_flat, b, c, dflat, d, eflat, e, f, gflat, g, aflat

The scale pattern describes the scale degrees used.

	pentatonic
		all 3 note permutations of [1, 2, 3, 5, 6]

	final
		all 3 note permutations from the chromatic scale that end
		on the tonic of the key.

    plus_four
		all 3 note permutations of [1, 2, 3, 4, 5, 6] that contain 4

    plus_seven
		all 3 note permutations of [1, 2, 3, 5, 6, 7] that contain 7

    four_and_seven
		all 3 note permutations of [1, 2, 3, 4, 5, 6, 7] that contain 4 and 7

	raised_five
		all 3 note permutations of [1, 2, 3, #5, 6] that contain #5

    raised_five_with_four_or_seven
		all 3 note permutations of [1, 2, 3, 4, #5, 6, 7] that contain #5 and
		at least one of [4, 7]
	
The instrument names correspond to the names in the General Midi Sound Set.
`

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// within returns True if val between lo and hi, inclusive.
func within(lo int, val int, hi int) bool {
	return lo <= val && val <= hi
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

var debug bool // enables some diagnostic output when true

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
	flag.BoolVar(&debug, "d", false, "Enable diagnostic output")

	var advancing bool
	flag.BoolVar(&advancing, "a", false, "Use advancing rhythm pattern")

	var instrument int
	flag.IntVar(&instrument, "i", 1, "General Midi instrument number: 1 ... 128")

	var midilo int
	flag.IntVar(&midilo, "l", 36, "Lowest desired Midi pitch")

	var midihi int
	flag.IntVar(&midihi, "u", 84, "Highest desired Midi pitch")

	var tempo int
	flag.IntVar(&tempo, "t", 120, "tempo in beats per minute")

	// server mode flags
	var serve bool
	flag.BoolVar(&serve, "s", false, "Run application as a server.")

	var midijsPath string
	flag.StringVar(&midijsPath, "m", filepath.Join(userHomeDir(), "go", "src", "github.com", "Michael-F-Ellis", "infinite-etudes", "midijs"), "Path to midijs files on your host (server-mode only)")

	var hostport string
	flag.StringVar(&hostport, "p", "localhost:8080", "hostname (or IP) and port to serve on. (server-mode only)")

	var expireSeconds int
	flag.IntVar(&expireSeconds, "x", 3600, "Maximum age in seconds for generated files (server-mode only)")

	// make sure all flags are defined before calling this
	flag.Parse()

	// validate flags
	if !within(1, instrument, 128) {
		log.Fatalln("instrument must be in range 1 to 128")
	}
	instrument-- // convert to 0 indexed

	if !within(0, midilo, 93) {
		log.Fatalln("midilo must be between 0 and 93")
	}

	if !within(24, midihi, 127) {
		log.Fatalln("midihi must be between 24 and 127")
	}

	if midihi-midilo < 24 {
		log.Fatalln("midihi must be at least 24 semitones above midilo")
	}

	if !within(20, tempo, 300) {
		log.Fatalln("tempo must be between 20 and 300 bpm")
	}

	if serve {
		serveEtudes(hostport, expireSeconds, midijsPath)
	} else {
		// create the midi files
		mkAllEtudes(midilo, midihi, tempo, instrument, "", advancing)
	}

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

// mkAllEtudes creates in the current directory all the etude files we support
// for the specified instrument. The arguments are assumed to be previously
// vetted and are not checked.
func mkAllEtudes(midilo, midihi, tempo, instrument int, iname string, advancing bool) {
	// Create and write all the tonal output files for all 12 key signatures

	for i := 0; i < 12; i++ {
		mkKeyEtudes(i, midilo, midihi, tempo, instrument, iname, advancing)
	}
	mkFinalEtudes(midilo, midihi, tempo, instrument, iname, advancing)
	mkIntervalEtudes(midilo, midihi, tempo, instrument, iname, advancing)
}

// mkKeyEtudes generates the six files associated with keynum where
// 0->c, 1->dflat, 2->d, ... 11->b
func mkKeyEtudes(keynum int, midilo int, midihi int, tempo int,
	instrument int, iname string, advancing bool) {
	for _, sequence := range generateKeySequences(keynum, midilo, midihi, tempo, instrument, iname) {
		mkMidi(&sequence, advancing, false)
		if debug {
			fmt.Println(pitchHistogram(sequence))
		}
	}
}

// mkIntervalEtudes generates the 12 interval files associated with pitch numbers where
// 0->c, 1->dflat, 2->d, ... 11->b and the 12 interval files associated with interval
// sizes from m2 to P8
func mkIntervalEtudes(midilo int, midihi int, tempo int,
	instrument int, iname string, advancing bool) {
	for _, sequence := range generateIntervalSequences(midilo, midihi, tempo, instrument, iname) {
		mkMidi(&sequence, advancing, true)
		if debug {
			fmt.Println(pitchHistogram(sequence))
		}
	}
	for _, sequence := range generateEqualIntervalSequences(midilo, midihi, tempo, instrument, iname) {
		mkMidi(&sequence, advancing, true)
		if debug {
			fmt.Println(pitchHistogram(sequence))
		}
	}
}

// mkFinalEtudes generates the 12 files associated with pitch numbers where
// 0->c, 1->dflat, 2->d, ... 11->b
func mkFinalEtudes(midilo int, midihi int, tempo int,
	instrument int, iname string, advancing bool) {
	for _, sequence := range generateFinalSequences(midilo, midihi, tempo, instrument, iname) {
		mkMidi(&sequence, advancing, false)
		if debug {
			fmt.Println(pitchHistogram(sequence))
		}
	}
}

// pitchHistorgram counts the pitches in each octave and returns a string with
// the filename followed by the counts in octaves 0-11. It panics if any pitch
// is outside the valid midi range 0-127. This func is primarily a debug tool
// to verify that the pitch distribution is roughly uniform over the desired
// octave range.
func pitchHistogram(e etudeSequence) (histo string) {
	var counts [11]int
	for _, triple := range e.seq {
		for _, p := range triple {
			bin := p / 12
			if bin < 0 || bin > 11 {
				panic(fmt.Sprintf("impossible midi pitch %d in sequence for file %s", p, e.filename))
			}
			counts[bin]++
		}
		histo = fmt.Sprintf("%s %v", e.filename, counts)
	}
	return
}
