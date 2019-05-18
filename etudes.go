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
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
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

// for midi key signatures sharps are positive, flats are negative.
var keySharps = map[string]int{
	"c": 0, "dflat": -5, "d": 2, "eflat": -3,
	"e": 4, "f": -1, "gflat": -6, "g": 1,
	"aflat": -4, "a": 3, "bflat": -2, "b": 5,
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// within returns True if val between lo and hi, inclusive.
func within(lo int, val int, hi int) bool {
	return lo <= val && val <= hi
}

func main() {
	// Parse command line
	flag.Usage = usage

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
	}
}

// mkFinalEtudes generates the 12 files associated with pitch numbers where
// 0->c, 1->dflat, 2->d, ... 11->b
func mkFinalEtudes(midilo int, midihi int, tempo int, instrument int) {
	for _, sequence := range generateFinalSequences(midilo, midihi, tempo, instrument) {
		mkMidi(&sequence)
	}
}

// getScale returns the major or harmonic minor
// scale in the specified key signature.
func getScale(keynum int, isminor bool) []int {
	scale := []int{0, 2, 4, 5, 7, 9, 11}
	if isminor {
		scale[4] = 8 // raised 5th is minor leading tone
	}
	for i, p := range scale {
		scale[i] = (p + keynum) % 12
	}
	return scale
}

// getChromaticScale returns the chromatic scale
func getChromaticScale() (scale []int) {
	scale = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	return
}

// permute3 returns a slice of midiTriple containing
// all possible permutations of 3 distinct notes in the
// scale.
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

// generateFinalSequences returns a slice of 12 etudeSequences as described in the usage instructions.
// Each sequence consists of all possible triples with a final pitch corresponding to pitchnum.
func generateFinalSequences(midilo int, midihi int, tempo int, instrument int) (sequences []etudeSequence) {
	// Get the chromatic scale as midi numbers in the range 0 - 11
	midiChromaticScaleNums := getChromaticScale()
	// Generate all 3 note permutations
	triples := permute3(midiChromaticScaleNums)

	sname, err := gmSoundName(instrument)
	if err != nil {
		panic("instrument number should have already been validated")
	}
	iname := gmSoundFileNamePrefix(sname)

	// construct the sequences
	for pitch := 0; pitch < 12; pitch++ {
		pitchname := keyNames[pitch]
		sequences = append(sequences, etudeSequence{
			filename:   pitchname + "_final" + "_" + iname + ".mid",
			midilo:     midilo,
			midihi:     midihi,
			tempo:      tempo,
			instrument: instrument,
			keyname:    pitchname,
		})
	}

	// filter the triples into the corresponding etude sequences
	// starting with majors
	for _, t := range triples {
		final := t[2]
		sequences[final].seq = append(sequences[final].seq, t)
	}
	return
}

// generateKeySequences returns a slice of six etudeSequences as described in the usage instructions.
func generateKeySequences(keynum int, midilo int, midihi int, tempo int, instrument int) []etudeSequence {
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
			if t[0] == four || t[0] == seven || t[2] == four || t[2] == seven {
				raisedFiveWithFourOrSeven.seq = append(raisedFiveWithFourOrSeven.seq, t)
				continue
			} else {
				raisedFive.seq = append(raisedFive.seq, t)
				continue
			}
		}

		// if we get to here, neither t[0] or t[1] are sharpfive
		if t[2] == sharpfive {
			if t[0] == four || t[0] == seven || t[1] == four || t[1] == seven {
				raisedFiveWithFourOrSeven.seq = append(raisedFiveWithFourOrSeven.seq, t)
			} else {
				raisedFive.seq = append(raisedFive.seq, t)
			}
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

// mkMidi shuffles a sequence and then offsets each triple as needed to keep
// the pitches within the limits specified in the sequence. Finally, it calls
// writeMidi file to convert the data to Standard Midi form and write it to
// disk.
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

// shuffle puts a slice of midiTriple in random order using the Fisher-Yates
// algorithm,
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

// writeMidiFile creates a midi file from an etudeSequence.
// Each midiTriple in the sequence is placed on beats 1, 2, 3 of
// a 4/4 measure with rest on beat 4. Each measure is played
// 4 times accompanied by a metronome track.  The etude begins
// with a one-bar count-in.
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
		trackInstrument(sequence),
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
	// end of track (note last delta is already in place)
	eot := []byte{0xff, 0x2f, 0x00}
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

// fourBarMusic returns a byte buffer containing four bars of  one midiTriple
// as described in function comment for writeMidiFile.
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

// metronomeBars returns a byte buffer containing n bars of metronome click.
// Downbeats use a High Wood Block sound. Other beats use a Low Wood Block,
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

// keySignature returns a MIDI KeySignature event preceeded by zero delta time.
func keySignature(s *etudeSequence) []byte {
	sharps := keySharps[s.keyname]
	sf := byte(sharps & 0xFF) // because flats are negative ints
	mi := byte(0)             // always major in this code
	return []byte{0x0, 0xFF, 0x59, 0x02, sf, mi}
}

// trackInstrument returns a Program Change event with the instrument specified
// in s preceeded by 0 delta time,
func trackInstrument(s *etudeSequence) []byte {
	return []byte{0x00, 0xC0, byte(s.instrument)}
}

// composeFileName returns a string containing a filename of the form
// key_scaledescription_intrument.mid, e.g
// "gflat_pentatonic_acoustic_grand_piano.mid"
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

// constrain applies tighten to a midiTriple, then adjusts its octave
// so that the first pitch is as close as possible to the last pitch
// of a previous triple. Then it checks to see if any of the adjusted
// pitches are above midihi or below midilow and re-adjusts the octave
// as needed to keep the pitches within the limits.
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
