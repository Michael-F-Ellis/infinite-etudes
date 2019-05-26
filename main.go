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
	"time"
)

const copyright = `
Copyright 2019 Ellis & Grant, Inc. All rights reserved.  Use of the source
code is governed by an MIT-style license that can be found in the LICENSE
file.
`
const description = `
	etudes generates a set of 6 midi files for each of 12 key signature. Each set
	covers all possible combinations of 3 pitches within the key. 
  
    The file names for each are structured as follows:
		
		eflat_pentatonic.mid
		eflat_plus_four.mid
		eflat_plus_seven.mid
		eflat_four_and_seven.mid
		eflat_raised_five.mid
		eflat_raised_five_with_four_or_seven.mid

	Each file name begins with keyname. The 12 keynames used are:

		a, b_flat, b, c, dflat, d, eflat, e, f, gflat, g, aflat

	The remainder of the file name describes the scale degrees used.

		pentatonic
			all 3 note permutations of [1, 2, 3, 5, 6]

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
	
`

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// within returns True if val between lo and hi, inclusive.
func within(lo int, val int, hi int) bool {
	return lo <= val && val <= hi
}

var debug bool // enables some diagnostic output when true

func main() {
	// Parse command line
	flag.Usage = usage

	flag.BoolVar(&debug, "d", true, "Enable diagnostic output")

	var instrument int
	flag.IntVar(&instrument, "i", 1, "General Midi instrument number: 1 ... 128")

	var midilo int
	flag.IntVar(&midilo, "l", 36, "Lowest desired Midi pitch")

	var midihi int
	flag.IntVar(&midihi, "u", 84, "Highest desired Midi pitch")

	var tempo int
	flag.IntVar(&tempo, "t", 120, "tempo in beats per minute")

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

	// Create and write all the output files for all 12 key signatures
	for i := 0; i < 12; i++ {
		mkKeyEtudes(i, midilo, midihi, tempo, instrument)
	}
	mkFinalEtudes(midilo, midihi, tempo, instrument)

}

// usage extends the flag package's default help message.
func usage() {
	fmt.Println(copyright)
	fmt.Printf("Usage: etudes [OPTIONS]\n  -h    print this help message.\n")
	flag.PrintDefaults()
	fmt.Println(description)

}

// mkKeyEtudes generates the six files associated with keynum where
// 0->c, 1->dflat, 2->d, ... 11->b
func mkKeyEtudes(keynum int, midilo int, midihi int, tempo int, instrument int) {
	for _, sequence := range generateKeySequences(keynum, midilo, midihi, tempo, instrument) {
		mkMidi(&sequence)
		if debug {
			fmt.Println(pitchHistogram(sequence))
		}
	}
}

// mkFinalEtudes generates the 12 files associated with pitch numbers where
// 0->c, 1->dflat, 2->d, ... 11->b
func mkFinalEtudes(midilo int, midihi int, tempo int, instrument int) {
	for _, sequence := range generateFinalSequences(midilo, midihi, tempo, instrument) {
		mkMidi(&sequence)
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
