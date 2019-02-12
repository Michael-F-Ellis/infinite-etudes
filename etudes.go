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
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
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

// Global map of output filenames and file objects.
var _gOutputs = make(map[string]*os.File)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// within returns True if val between lo and hi, inclusive.
func within(lo int, val int, hi int) bool {
	return lo <= val && val <= hi
}

func main() {
	// Ensure we exit with an error code and log message
	// when needed after deferred cleanups have run.
	// Credit: https://tinyurl.com/ycv9zpbn
	var err error
	defer func() {
		if err != nil {
			removeOutputFiles()
			log.Fatalln(err)
		}
	}()

	// Close any opened output files on exit.
	defer closeOutputFiles()

	// Parse command line
	flag.Usage = usage

	var instrument int
	flag.IntVar(&instrument, "i", 0, "General Midi instrument number: 0 ... 127")

	var keyname string
	keynames := []string{"c", "dflat", "d", "eflat", "e", "f", "gflat", "g", "aflat", "a", "bflat", "b"}
	h := fmt.Sprintf("Key name: one of %v", keynames)
	flag.StringVar(&keyname, "k", "c", h)

	var midilo int
	flag.IntVar(&midilo, "l", 36, "Lowest desired Midi pitch")

	var midihi int
	flag.IntVar(&midihi, "u", 84, "Highest desired Midi pitch")

	var tempo int
	flag.IntVar(&tempo, "t", 120, "tempo in beats per minute")

	flag.Parse()

	// validate flags
	if !within(0, instrument, 127) {
		log.Fatalln("instrument must be in range 0 to 127")
	}

	keynum := -1
	for i, name := range keynames {
		if keyname == name {
			keynum = i
			break
		}
	}
	if keynum == -1 {
		log.Fatalf("invalid key name %s", keyname)
	}

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

	// Create and write all the output files
	for i := 0; i < 12; i++ {
		mkKeyEtudes(i, midilo, midihi, tempo, instrument)
	}

}

// usage extends the flag package's default help message.
func usage() {
	fmt.Println(copyright)
	fmt.Printf("Usage: etudes [OPTIONS]\n  -h    print this help message.\n")
	flag.PrintDefaults()
	fmt.Println(description)

}

// mkEtudes generates the six files associated with keynum where
// 0->c, 1->dflat, 2->d, ... 11->b
func mkKeyEtudes(keynum int, midilo int, midihi int, tempo int, instrument int) {
	for _, sequence := range generateSequences(keynum) {
		mkMidi(&sequence, midilo, midihi, tempo, instrument)
	}
}

func getScale(keynum int, isminor bool) []int {
	scale := []int{0, 2, 4, 5, 7, 9, 11}
	if isminor {
		scale[4] = 8 // raised 5th is minor leading tone
	}
	for i, p := range scale {
		scale[i] = (p + keynum) % 12
	}
	sort.Ints(scale)
	return scale
}

type midiTriple [3]int

type midiTripleSequence struct {
	seq  []midiTriple
	name string
}

func permute3(scale []int) []midiTriple {
	var permutations []midiTriple
	for i, p := range scale {
		for j, q := range scale {
			if j == i {
				continue
			}
			for k, r := range scale {
				if k == i || k == j {
					continue
				}
				t := midiTriple{p, q, r}
				permutations = append(permutations, t)
			}
		}
	}
	return permutations
}

func generateSequences(keynum int) []midiTripleSequence {
	// Look up the keyname string
	keynames := []string{"c", "dflat", "d", "eflat", "e", "f", "gflat", "g", "aflat", "a", "bflat", "b"}
	keyname := keynames[keynum]
	// Get the major and harmonic minor scales as midi numbers in the range 0 - 11
	midiMajorScaleNums := getScale(keynum, false)
	midiMinorScaleNums := getScale(keynum, true)
	// Generate all 3 note permutations
	majors := permute3(midiMajorScaleNums)
	minors := permute3(midiMinorScaleNums)

	// declare the sequences
	pentatonic := midiTripleSequence{name: keyname + "_pentatonic"}
	plusFour := midiTripleSequence{name: keyname + "_plus_four"}
	plusSeven := midiTripleSequence{name: keyname + "_plus_seven"}
	fourAndSeven := midiTripleSequence{name: keyname + "_four_and_seven"}
	raisedFive := midiTripleSequence{name: keyname + "_raised_five"}
	raisedFiveWithFourOrSeven := midiTripleSequence{name: keyname + "_raised_five_with_four_or_seven"}

	// scale degree midi values for this key
	four := midiMajorScaleNums[3]
	seven := midiMajorScaleNums[6]
	sharpfive := midiMinorScaleNums[4]

	// filter the triples in majors and minors to the appropriate sequences
	// starting with majors
	for _, t := range majors {
		switch t[0] {
		case four:
			// it's either plus4 or 4and7
			if t[1] == seven || t[2] == seven {
				fourAndSeven.seq = append(fourAndSeven.seq, t)
				continue
			} else {
				plusFour.seq = append(plusFour.seq, t)
				continue
			}

		case seven:
			// it's either plus7 or 4and7
			if t[1] == four || t[2] == four {
				fourAndSeven.seq = append(fourAndSeven.seq, t)
				continue
			} else {
				plusSeven.seq = append(plusSeven.seq, t)
				continue
			}

		}
		// if we get to here, t[0] is not 4 or seven
		switch t[1] {
		case four:
			// it's either plus4 or 4and7
			if t[2] == seven {
				fourAndSeven.seq = append(fourAndSeven.seq, t)
				continue
			} else {
				plusFour.seq = append(plusFour.seq, t)
				continue
			}

		case seven:
			// it's either plus7 or 4and7
			if t[2] == four {
				fourAndSeven.seq = append(fourAndSeven.seq, t)
				continue
			} else {
				plusSeven.seq = append(plusSeven.seq, t)
				continue
			}
		}
		// if we get to here, neither t[0] or t[1] are four or seven
		switch t[2] {
		case four:
			plusFour.seq = append(plusFour.seq, t)
		case seven:
			plusSeven.seq = append(plusSeven.seq, t)
		default:
			pentatonic.seq = append(pentatonic.seq, t)
		}

	}

	// now deal with minors
	for _, t := range minors {
		if t[0] == sharpfive {
			if t[1] == four || t[1] == seven || t[2] == four || t[2] == seven {
				raisedFiveWithFourOrSeven.seq = append(raisedFiveWithFourOrSeven.seq, t)
				continue
			} else {
				raisedFive.seq = append(raisedFive.seq, t)
				continue
			}
		}
		// if we get to here, t[0] is not sharpfive
		if t[1] == sharpfive {
			if t[2] == four || t[2] == seven {
				raisedFiveWithFourOrSeven.seq = append(raisedFiveWithFourOrSeven.seq, t)
				continue
			} else {
				raisedFive.seq = append(raisedFive.seq, t)
				continue
			}
		}

		// if we get to here, neither t[0] or t[1] are sharpfive
		if t[2] == sharpfive {
			raisedFive.seq = append(raisedFive.seq, t)
		}
		// we don't care about any other triples because they've all been
		// accounted for by processing the majors.

	}

	// create the slice of sequences
	sequences := []midiTripleSequence{
		pentatonic, plusFour, plusSeven, fourAndSeven,
		raisedFive, raisedFiveWithFourOrSeven}

	return sequences
}

// mkMidi shuffles the sequence and then offsets themreplicates each triple 3 times so that
// the expanded sequence will play each item 4 times.
//
func mkMidi(sequence *midiTripleSequence, midilo int, midihi int, tempo int, instrument int) {
	// check for programming errors
	if tempo < 20 || tempo > 300 {
		msg := fmt.Sprintf("refusing ridiculous tempo value of %d beats per minute", tempo)
		panic(msg)
	}
	if instrument > 127 {
		msg := fmt.Sprintf("expected midi instrument number <= 127, got %d", instrument)
		panic(msg)
	}
	// Note: first call to constrain will check pitch limits

	// Shuffle the sequence
	shuffle(sequence.seq)

	// Constrain the sequence assuming a prior pitch of middle C (60)
	prior := 60
	seqlen := len(sequence.seq)
	for i := 0; i < seqlen; i++ {
		t := &(sequence.seq[i])
		constrain(t, prior, midilo, midihi)
		prior = t[2]
	}
	// Write the etude
	writeMidiFile(sequence, tempo, instrument)

}

// Fisher-Yates shuffle
func shuffle(slc []midiTriple) {
	N := len(slc)
	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		slc[r], slc[i] = slc[i], slc[r]
	}
}

func writeMidiFile(sequence *midiTripleSequence, tempo int, instrument int) {
	// compose the file name

	// open the file

	// write the header

	// write the tempo track

	// write the sequence

	// compose the metronome track

	// write the metronome track
}

func composeFileName(sequence *midiTripleSequence, instrument int) string {
	front := sequence.name
	sname, err := gmSoundName(instrument)
	if err != nil {
		msg := fmt.Sprintf("couldn't get instrument name: %v", err)
		panic(msg)
	}
	iname := gmSoundFileNamePrefix(sname)
	extension := ".mid"
	return front + "_" + iname + extension
}

// extractFileNames parses a line of text. If the line doesn't contain the
// delimiter, it returns an empty slice and a nil error to indicate that this
// line is to be output to whatever file targets are currently in effect.
// Otherwise it splits the line on whitespace. Each field after the delimiter
// is presumed to to be a file name and is appended to the names slice. Non-nil
// errors are returned unless the delimiter is found in exactly one field and
// there is at least on field following it.
func extractFileNames(delimiter string, line string) (names []string, err error) {
	// Short circuit if line doesn't contain delimiter
	if !strings.Contains(line, delimiter) {
		return
	}
	fields := strings.Fields(line)
	dfound := false
	for _, field := range fields {
		if !dfound {
			if field == delimiter {
				dfound = true
			}
			continue
		}
		if field == delimiter {
			err = fmt.Errorf("found more than one delimiter %s in line", delimiter)
			return names, err
		}
		names = append(names, field)
	}
	switch dfound {
	case false:
		err = fmt.Errorf("Delimiter %s must be surrounded by whitespace", delimiter)
	case true:
		if len(names) == 0 {
			err = fmt.Errorf("No file names found after delimiter %s", delimiter)
		}
	}
	return names, err
}

// adjustSuccessor returns a pitch adjusted to be within
// 6 semitones of its predecessor.
func adjustSuccessor(p0 int, p1 int) (adjustedP1 int) {
	for p1-p0 > 6 {
		p1 -= 12
	}
	for p1-p0 < -6 {
		p1 += 12
	}
	return p1
}

// tighten puts the pitches of a midiTriple in close sequential position,
// e.g. 13 19 16 -> 13 7 4, so that each note is within 6 semitones of the
// its predecessor.
func tighten(t *midiTriple) {
	// adjust second pitch relative to first
	t[1] = adjustSuccessor(t[0], t[1])
	// adjust third pitch relative to (adjusted) second.
	t[2] = adjustSuccessor(t[1], t[2])
}

func constrain(t *midiTriple, prior int, midilo int, midihi int) {
	if midilo > 127 || midihi > 127 || midihi-midilo < 24 {
		msg := fmt.Sprintf("Invalid midi limits %v, %v", midilo, midihi)
		panic(msg) // Programming error. Bad limits should be rejected at startup
	}
	// Adjust first pitch relative to prior.
	t[0] = adjustSuccessor(prior, t[0])
	// Tighten remaining pitches
	tighten(t)
	// Shift pitches by octaves until all are between midilo and midihi inclusive.
	lo := int(midilo)
	for t[0] < lo || t[1] < lo || t[2] < lo {
		t[0] += 12
		t[1] += 12
		t[2] += 12
	}
	hi := int(midihi)
	for t[0] > hi || t[1] > hi || t[2] > hi {
		t[0] -= 12
		t[1] -= 12
		t[2] -= 12
	}
}

// openOutputFiles is called with results from extractFileNames. For each name
// in the list, It checks the outputs map to see if the file is already opened.
// If so, it ignores the name and moves on to the next one.  Otherwise it
// attempts to open the file for writing, truncating it if it exists. If
// successful it adds it to outputs map. On failure, it returns the error from
// os.Create immediately without attempting to open any further files from the
// names list.
func openOutputFiles(names []string) error {
	var err error
	for _, name := range names {
		isnew := true
		for oname := range _gOutputs {
			if oname == name {
				isnew = false
				break
			}
		}
		if isnew {
			fd, err := os.Create(name)
			if err != nil {
				return err
			}
			_gOutputs[name] = fd
		}
	}
	return err
}

// closeOutputFiles is used as a deferred call in main to ensure that all
// output files are closed on exit.
func closeOutputFiles() {
	for _, fd := range _gOutputs {
		fd.Close()
	}
}

// removeOutputFiles is used in main to ensure that all output files are
// removed if an error has occurred.
func removeOutputFiles() {
	for name := range _gOutputs {
		os.Remove(name)
	}
}

// processInputFile reads every line from fd and scans to see
// if it contains delimiter. If not, the line is output to
// all currently active target files. If the line contains delimiter
// it parses the remainder of the line as a list of space-delimited
// filenames for output. These are passed to openOutputFiles() to be
// opened if they haven't been opened already. If openOutputFiles()
// succeeds, the files are set as the current output targets for
// following lines until the next delimiter line is encountered.
//
// If parsing fails, the error from extractFileNames() is returned.
// Similarly, processing ends if openOutputFiles() fails.
// Processing ends normally when all lines in the file have been
// read and processed.
func processInputFile(fd *os.File, delimiter string) error {
	defer fd.Close()
	var err error
	reader := bufio.NewReader(fd)
	var targets = make([]*os.File, 0)
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return err
		}
		names, err := extractFileNames(delimiter, line)
		if err != nil {
			return err
		}
		if len(names) == 0 {
			// lineout := line + "\n"
			for _, f := range targets {
				f.WriteString(line)
			}
		} else {
			err = openOutputFiles(names)
			if err != nil {
				return err
			}
			targets = make([]*os.File, 0)
			for _, name := range names {
				targets = append(targets, _gOutputs[name])
			}
		}
	}
	return err
}
