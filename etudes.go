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
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
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

type midiTriple [3]int

type etudeSequence struct {
	seq        []midiTriple
	midilo     int
	midihi     int
	tempo      int
	instrument int
	keyname    string
	filename   string
}

var keyNames = []string{"c", "dflat", "d", "eflat", "e", "f", "gflat", "g", "aflat", "a", "bflat", "b"}

// for midi key sig sharps are positive, flats negative,
var keySharps = map[string]int{"c": 0, "dflat": -5, "d": 2, "eflat": -3, "e": 4, "f": -1, "gflat": -6, "g": 1, "aflat": -4, "a": 3, "bflat": -2, "b": 5}

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
	for _, sequence := range generateSequences(keynum, midilo, midihi, tempo, instrument) {
		mkMidi(&sequence)
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
	// sort.Ints(scale)
	return scale
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

func generateSequences(keynum int, midilo int, midihi int, tempo int, instrument int) []etudeSequence {
	// Look up the keyname string
	keyname := keyNames[keynum]
	// Get the major and harmonic minor scales as midi numbers in the range 0 - 11
	midiMajorScaleNums := getScale(keynum, false)
	midiMinorScaleNums := getScale(keynum, true)
	// Generate all 3 note permutations
	majors := permute3(midiMajorScaleNums)
	minors := permute3(midiMinorScaleNums)

	sname, err := gmSoundName(instrument)
	if err != nil {
		panic("instrument number should have already been validated")
	}
	iname := gmSoundFileNamePrefix(sname)

	// declare the sequences
	pentatonic := etudeSequence{
		filename:   keyname + "_pentatonic" + "_" + iname + ".mid",
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
	}
	plusFour := etudeSequence{
		filename:   keyname + "_plus_four" + "_" + iname + ".mid",
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
	}
	plusSeven := etudeSequence{
		filename:   keyname + "_plus_seven" + "_" + iname + ".mid",
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
	}
	fourAndSeven := etudeSequence{
		filename:   keyname + "_four_and_seven" + "_" + iname + ".mid",
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
	}
	raisedFive := etudeSequence{
		filename:   keyname + "_raised_five" + "_" + iname + ".mid",
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
	}
	raisedFiveWithFourOrSeven := etudeSequence{
		filename:   keyname + "_raised_five_with_four_or_seven" + "_" + iname + ".mid",
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
	}

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
	sequences := []etudeSequence{
		pentatonic, plusFour, plusSeven, fourAndSeven,
		raisedFive, raisedFiveWithFourOrSeven}

	return sequences
}

// mkMidi shuffles the sequence and then offsets themreplicates each triple 3 times so that
// the expanded sequence will play each item 4 times.
//
func mkMidi(sequence *etudeSequence) {
	// Shuffle the sequence
	shuffle(sequence.seq)

	// Constrain the sequence assuming a prior pitch of middle C (60)
	prior := 60
	seqlen := len(sequence.seq)
	for i := 0; i < seqlen; i++ {
		t := &(sequence.seq[i])
		constrain(t, prior, sequence.midilo, sequence.midihi)
		prior = t[2]
	}
	// Write the etude
	writeMidiFile(sequence)

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

// low3 returns a 3 byte array representing the lower
// 3 bytes of n, e.g. as a 24 bit number
func low3(n uint32) (u24 [3]byte) {
	u24[0] = byte((n & 0x00FFFFFF) >> 16)
	u24[1] = byte((n & 0x0000FFFF) >> 8)
	u24[2] = byte((n & 0x000000FF))
	return u24
}

func writeMidiFile(sequence *etudeSequence) {
	// open the file
	fd, err := os.Create(sequence.filename)
	if err != nil {
		msg := fmt.Sprintf("Couldn't open output file %s", sequence.filename)
		panic(msg)
	}
	defer fd.Close()
	// write the header "MThd len=6, format=1, tracks=6, ticks=960"
	header := []byte{0x4d, 0x54, 0x68, 0x64, 0, 0, 0, 6, 0, 1, 0, 3, 3, 192}
	n, err := fd.Write(header)
	if err != nil {
		panic(err)
	}
	if n != len(header) {
		panic("failed to write header")
	}
	// write the tempo track
	microseconds := low3(uint32(60000000 / sequence.tempo)) //microseconds per beat
	var record = []interface{}{
		// Time signature event
		byte(0),                // delta time
		low3(uint32(0xFF5804)), // tempo event
		byte(4),                // beats per measure
		byte(2),                // quarter note beat (because 2^2 = 4)
		byte(24),               // clocks per tick
		byte(8),                // 32nd's per quarter note
		// Tempo event
		byte(0),                // delta time
		low3(uint32(0xFF5103)), // tempo event
		microseconds,
		// EOT event
		byte(0),                // delta time
		low3(uint32(0xFF2F00)), // End of track
	}
	// write the track data to a temporary buffer
	// so we can compute its length
	buf := new(bytes.Buffer)
	for _, v := range record {
		err = binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			panic(err)
		}
	}
	// prepend the track header.
	var track = []interface{}{
		[]byte{'M', 'T', 'r', 'k'},
		uint32(buf.Len()), // length of track data
		buf.Bytes(),
	}
	// write tempo track to file
	for _, v := range track {
		err = binary.Write(fd, binary.BigEndian, v)
		if err != nil {
			panic(err)
		}
	}

	// write the instrument track
	buf = new(bytes.Buffer)
	record = []interface{}{
		keySignature(sequence),
		byte(0x9e), // four beats hi byte
		byte(0x00), // four beats lo byte
	}
	for _, v := range record {
		err = binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			panic(err)
		}
	}
	for _, t := range sequence.seq {
		music := fourBarsMusic(t).Bytes()
		err = binary.Write(buf, binary.BigEndian, music)
		if err != nil {
			panic(err)
		}
	}
	// end of track
	eot := []byte{0x00, 0xff, 0x2f, 0x00}
	err = binary.Write(buf, binary.BigEndian, eot)
	if err != nil {
		panic(err)
	}
	// prepend the track header.
	track = []interface{}{
		[]byte{'M', 'T', 'r', 'k'},
		uint32(buf.Len()), // length of track data
		buf.Bytes(),
	}
	// write instrument track to file
	for _, v := range track {
		err = binary.Write(fd, binary.BigEndian, v)
		if err != nil {
			panic(err)
		}
	}

	// compose the metronome track
	buf = new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, byte(0x00))
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(sequence.seq); i++ {
		var music []byte
		switch i {
		case 0:
			music = metronomeBars(5).Bytes()
		default:
			music = metronomeBars(4).Bytes()
		}
		err = binary.Write(buf, binary.BigEndian, music)
		if err != nil {
			panic(err)
		}
	}
	// end of track
	err = binary.Write(buf, binary.BigEndian, eot)
	if err != nil {
		panic(err)
	}

	// write the metronome track
	// prepend the track header.
	track = []interface{}{
		[]byte{'M', 'T', 'r', 'k'},
		uint32(buf.Len()), // length of track data
		buf.Bytes(),
	}
	// write metronome track to file
	for _, v := range track {
		err = binary.Write(fd, binary.BigEndian, v)
		if err != nil {
			panic(err)
		}
	}

}
func fourBarsMusic(t midiTriple) *bytes.Buffer {
	// These are the only variable length delta times we need.
	noBeats := byte(0x00)
	oneBeatHi := byte(0x87)
	oneBeatLo := byte(0x40)
	// fourBeats := []byte{0x9e, 0x00}

	velocity1 := byte(0x65) // downbeat
	velocity2 := byte(0x51) // other beats

	on := byte(0x90)  // Note On, channel 1
	off := byte(0x80) // Note off, channel 1.

	buf := new(bytes.Buffer)
	check := func(e error) {
		if e != nil {
			panic(e)
		}
	}
	// mkBeat writes MIDI for one beat with note on and off events.
	mkBeat := func(buf *bytes.Buffer, pitch byte, velocity byte, after int) {
		var b []byte
		switch after {
		case 0:
			b = []byte{on, pitch, velocity, oneBeatHi, oneBeatLo, off, pitch, velocity, noBeats}
		case 1:
			b = []byte{on, pitch, velocity, oneBeatHi, oneBeatLo, off, pitch, velocity, oneBeatHi, oneBeatLo}
		default:
			panic(errors.New("programming error: arg after must be 0 or 1"))
		}

		check(binary.Write(buf, binary.BigEndian, b))
	}

	// write all 4 bars for this triple
	for i := 0; i < 4; i++ {
		// first beat
		pitch := byte(t[0])
		mkBeat(buf, pitch, velocity1, 0)
		// 2nd beat
		pitch = byte(t[1])
		mkBeat(buf, pitch, velocity2, 0)
		// 3rd beat (4th beat is a rest, so we append a one beat delay after the Note Off event.
		pitch = byte(t[2])
		mkBeat(buf, pitch, velocity2, 1)
	}
	return buf
}

func metronomeBars(n int) *bytes.Buffer {
	// These are the only variable length delta times we need.
	noBeats := byte(0x00)
	oneBeatHi := byte(0x87)
	oneBeatLo := byte(0x40)
	// fourBeats := []byte{0x9e, 0x00}

	velocity1 := byte(0x65) // downbeat
	velocity2 := byte(0x51) // other beats

	on := byte(0x99)  // Note On, channel 10
	off := byte(0x89) // Note off, channel 10

	wbh := byte(0x4c) // wood block hi for downbeats
	wbl := byte(0x4d) // wood block lo for other beats

	buf := new(bytes.Buffer)
	check := func(e error) {
		if e != nil {
			panic(e)
		}
	}
	// mkBeat writes MIDI for one beat with note on and off events.
	mkBeat := func(buf *bytes.Buffer, pitch byte, velocity byte, after int) {
		var b []byte
		switch after {
		case 0:
			b = []byte{on, pitch, velocity, oneBeatHi, oneBeatLo, off, pitch, velocity, noBeats}
		case 1:
			b = []byte{on, pitch, velocity, oneBeatHi, oneBeatLo, off, pitch, velocity, oneBeatHi, oneBeatLo}
		default:
			panic(errors.New("programming error: arg after must be 0 or 1"))
		}

		check(binary.Write(buf, binary.BigEndian, b))
	}

	// write as many bars as requested
	for i := 0; i < n; i++ {
		// first beat
		mkBeat(buf, wbh, velocity1, 0)
		// 2nd beat
		mkBeat(buf, wbl, velocity2, 0)
		// 3rd beat
		mkBeat(buf, wbl, velocity2, 0)
		// 4th beat
		mkBeat(buf, wbl, velocity2, 0)
	}
	return buf
}

func keySignature(s *etudeSequence) []byte {
	sharps := keySharps[s.keyname]
	sf := byte(sharps & 0xFF) // because flats are negative ints
	mi := byte(0)             // always major in this code
	return []byte{0x0, 0xFF, 0x59, 0x02, sf, mi}
}

func composeFileName(sequence *etudeSequence, instrument int) string {
	front := sequence.filename
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
