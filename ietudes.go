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
	"fmt"
	"math/rand"
	"os"
	"strings"
)

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

// tripleFromInterval builds a midiTriple from a pair of pitch numbers.
// The up argument determines if second pitch is to be adjusted an octave up or down, depending on the
// interval.
func tripleFromInterval(p, q int, up bool) (t midiTriple) {
	switch {
	case p == q: // convert unison to an octave
		if up {
			q += 12
		} else {
			q -= 12
		}
	case p < q:
		if !up {
			q -= 12
		}
	case p > q:
		if up {
			q += 12
		}
	default:
		panic("This should be impossible")
	}
	t = midiTriple{p, q, p}
	return
}

// flip simulates a fair coin flip
func flip() (up bool) {
	if rand.Intn(2) == 1 {
		up = true
	}
	return
}

// permute2 returns a slice of midiTriple. Each triple is built from an interval
// with the first note repeated as the third note, e.g. {3, 5, 3} or {3, -7, 3}
func permute2(scale []int) (permutations []midiTriple) {
	for _, p := range scale {
		for _, q := range scale {
			permutations = append(permutations, tripleFromInterval(p, q, flip()))
		}
	}
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

// generateEqualIntervalSequences returns a slice of etudeSequences as described in the usage instructions.
// Each sequence consists of triples of equal interval sizes
func generateEqualIntervalSequences(midilo int, midihi int, tempo int, instrument int, iname string) (sequences []etudeSequence) {
	// Get the chromatic scale as midi numbers in the range 0 - 11
	midiChromaticScaleNums := getChromaticScale()
	// Generate all intervals
	triples := permute2(midiChromaticScaleNums)

	if iname == "" {
		sname, err := gmSoundName(instrument)
		if err != nil {
			panic("instrument number should have already been validated")
		}
		iname = gmSoundFileNamePrefix(sname)
	}
	var intervalNames []string
	for _, iinfo := range intervalInfo {
		intervalNames = append(intervalNames, iinfo.fileName)
	}

	// construct the sequences
	for i := 0; i < 12; i++ {
		intervalName := intervalNames[i]
		sequences = append(sequences, etudeSequence{
			filename:   intervalName + "_intervals" + "_" + iname + ".mid",
			midilo:     midilo,
			midihi:     midihi,
			tempo:      tempo,
			instrument: instrument,
			keyname:    intervalName,
		})
	}

	// filter the triples into the corresponding etude sequences
	for _, t := range triples {
		diff := t[0] - t[1]
		if diff < 0 {
			diff = -diff
		}
		diff -= 1
		sequences[diff].seq = append(sequences[diff].seq, t)
	}
	return
}

// generateIntervalSequences returns a slice of 12 etudeSequences as described in the usage instructions.
// Each sequence consists of 12 triples with the middle pitch corresponding to pitchnum.
func generateIntervalSequences(midilo int, midihi int, tempo int, instrument int, iname string) (sequences []etudeSequence) {
	// Get the chromatic scale as midi numbers in the range 0 - 11
	midiChromaticScaleNums := getChromaticScale()
	// Generate all intervals
	triples := permute2(midiChromaticScaleNums)

	if iname == "" {
		sname, err := gmSoundName(instrument)
		if err != nil {
			panic("instrument number should have already been validated")
		}
		iname = gmSoundFileNamePrefix(sname)
	}
	// construct the sequences
	for pitch := 0; pitch < 12; pitch++ {
		pitchname := keyNames[pitch]
		sequences = append(sequences, etudeSequence{
			filename:   pitchname + "_intervals" + "_" + iname + ".mid",
			midilo:     midilo,
			midihi:     midihi,
			tempo:      tempo,
			instrument: instrument,
			keyname:    pitchname,
		})
	}

	// filter the triples into the corresponding etude sequences
	for _, t := range triples {
		pitch := t[1] % 12
		if pitch < 0 {
			pitch += 12
		}
		sequences[pitch].seq = append(sequences[pitch].seq, t)
	}
	return
}

// generateFinalSequences returns a slice of 12 etudeSequences as described in the usage instructions.
// Each sequence consists of all possible triples with a final pitch corresponding to pitchnum.
func generateFinalSequences(midilo int, midihi int, tempo int, instrument int, iname string) (sequences []etudeSequence) {
	// Get the chromatic scale as midi numbers in the range 0 - 11
	midiChromaticScaleNums := getChromaticScale()
	// Generate all 3 note permutations
	triples := permute3(midiChromaticScaleNums)

	if iname == "" {
		sname, err := gmSoundName(instrument)
		if err != nil {
			panic("instrument number should have already been validated")
		}
		iname = gmSoundFileNamePrefix(sname)
	}
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
func generateKeySequences(keynum int, midilo int, midihi int, tempo int, instrument int, iname string) []etudeSequence {
	// Look up the keyname string
	keyname := keyNames[keynum]
	// Get the major and harmonic minor scales as midi numbers in the range 0 - 11
	midiMajorScaleNums := getScale(keynum, false)
	midiMinorScaleNums := getScale(keynum, true)
	// Generate all 3 note permutations
	majors := permute3(midiMajorScaleNums)
	minors := permute3(midiMinorScaleNums)

	if iname == "" {
		sname, err := gmSoundName(instrument)
		if err != nil {
			panic("instrument number should have already been validated")
		}
		iname = gmSoundFileNamePrefix(sname)
	}

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
func mkMidi(sequence *etudeSequence, advancing bool, noTighten bool) {
	// Shuffle the sequence
	shuffle(sequence.seq)

	// Constrain the sequence assuming a prior pitch halfway between the limits.
	prior := (sequence.midilo + sequence.midihi) / 2
	seqlen := len(sequence.seq)
	for i := 0; i < seqlen; i++ {
		t := &(sequence.seq[i])
		constrain(t, prior, sequence.midilo, sequence.midihi, noTighten)
		prior = t[2]
	}
	// Write the etude
	writeMidiFile(sequence, advancing)

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
func writeMidiFile(sequence *etudeSequence, advancing bool) {
	// update the filename with the rhythm pattern
	var newend = "_steady.mid"
	if advancing {
		newend = "_advancing.mid"
	}
	sequence.filename = strings.Replace(sequence.filename, ".mid", newend, 1)

	// open the file
	fd, err := os.Create(sequence.filename)
	if err != nil {
		msg := fmt.Sprintf("Couldn't open output file %s", sequence.filename)
		panic(msg)
	}
	defer fd.Close()
	// write the header "MThd len=6, format=1, tracks=3, ticks=960"
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
	var offset int
	for i, t := range sequence.seq {
		var u midiTriple
		if i < len(sequence.seq)-1 {
			u = sequence.seq[i+1]
		} else {
			// last triple in sequence, so pass a copy as the successor
			u = t
		}
		music := fourBarsMusic(t, u, advancing, offset).Bytes()
		err = binary.Write(buf, binary.BigEndian, music)
		if err != nil {
			panic(err)
		}
		if advancing {
			offset = (offset + 1) % 4
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
		switch {
		case i == 0:
			music = metronomeBars(5).Bytes()
		case advancing && (i%4 == 3):
			music = metronomeBars(3).Bytes()
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

// fourBarsMusic returns a byte buffer containing four bars of  one midiTriple
// Two rhythm formats are supported with selection determined by the boolean
// argument named advancing. When advancing is false, the rhythm is constant with
// the pitches of each triple falling on beats 1, 2, 3 and a rest on beat 4. This
// pattern is repeated in 4 identical bars.
//
// When advancing is true, the rhythm advances by 1 beat in the 4th bar so that
// a sequence of triples may be cycled through 4 different patterns by successively
// incrementing the offset argument.
//
// The advancing format expects the arguments t and u to be successive triples from
// the sequence. The non-advancing format ignores u.
//
// For the final triple in the sequence, it is recommended to pass the same value
// for t and u.
func fourBarsMusic(t, u midiTriple, advancing bool, offset int) *bytes.Buffer {
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
		var pitch byte
		if !advancing {
			// first beat
			pitch = byte(t[0])
			mkBeat(buf, pitch, velocity1, 0)
			// 2nd beat
			pitch = byte(t[1])
			mkBeat(buf, pitch, velocity2, 0)

			// 3rd beat (4th beat is a rest, so we append a one beat delay after the Note Off event.
			pitch = byte(t[2])
			mkBeat(buf, pitch, velocity2, 1)
			continue
		}
		// advancing
		// 3rd beat
		switch offset {
		case 0:
			// first beat
			pitch = byte(t[0])
			mkBeat(buf, pitch, velocity1, 0)
			// 2nd beat
			pitch = byte(t[1])
			mkBeat(buf, pitch, velocity2, 0)

			if i < 3 {
				// 3rd beat
				pitch = byte(t[2])
				mkBeat(buf, pitch, velocity2, 1)
				continue
			}
			pitch = byte(t[2])
			mkBeat(buf, pitch, velocity2, 0)

			// 4th beat
			pitch = byte(u[0])
			mkBeat(buf, pitch, velocity2, 0)
		case 1:
			// first beat
			pitch = byte(t[1])
			mkBeat(buf, pitch, velocity1, 0)
			// 2nd beat
			if i < 3 {
				pitch = byte(t[2])
				mkBeat(buf, pitch, velocity2, 1)

				// 4th beat
				pitch = byte(t[0])
				mkBeat(buf, pitch, velocity2, 0)
				continue
			}
			// 2nd beat
			pitch = byte(t[2])
			mkBeat(buf, pitch, velocity2, 0)

			// 3rd beat
			pitch = byte(u[0])
			mkBeat(buf, pitch, velocity2, 0)

			// 4th beat
			pitch = byte(u[1])
			mkBeat(buf, pitch, velocity2, 0)

		case 2:
			if i < 3 {
				// first beat
				pitch = byte(t[2])
				mkBeat(buf, pitch, velocity1, 1)

				// 2nd beat is a rest

				// 3rd beat
				pitch = byte(t[0])
				mkBeat(buf, pitch, velocity2, 0)

				// 4th beat
				pitch = byte(t[1])
				mkBeat(buf, pitch, velocity2, 0)
				continue
			}
			// 1st beat
			pitch = byte(t[2])
			mkBeat(buf, pitch, velocity1, 0)

			// 2nd beat
			pitch = byte(u[0])
			mkBeat(buf, pitch, velocity2, 0)

			// 3rd beat
			pitch = byte(u[1])
			mkBeat(buf, pitch, velocity2, 0)

			// 4th beat
			pitch = byte(u[2])
			mkBeat(buf, pitch, velocity2, 1)

		case 3:
			// only 3 bars for this offset
			if i == 3 {
				continue
			}
			// 2nd beat
			pitch = byte(t[0])
			mkBeat(buf, pitch, velocity2, 0)

			// 3rd beat
			pitch = byte(t[1])
			mkBeat(buf, pitch, velocity2, 0)

			if i < 2 {
				// 4th beat
				pitch = byte(t[2])
				mkBeat(buf, pitch, velocity2, 1)
				continue
			}

			// 4th beat
			pitch = byte(t[2])
			mkBeat(buf, pitch, velocity2, 0)
		}
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
func constrain(t *midiTriple, prior int, midilo int, midihi int, noTighten bool) {
	if midilo > 127 || midihi > 127 || midihi-midilo < 24 {
		msg := fmt.Sprintf("Invalid midi limits %v, %v", midilo, midihi)
		panic(msg) // Programming error. Bad limits should be rejected at startup
	}
	// Tighten the triple to close position
	if !noTighten {
		tighten(t)
	}
	// Shift tightened triple so that first pitch is as
	// close as possible to prior.
	offset := adjustSuccessor(prior, t[0]) - t[0]
	for i := 0; i < len(t); i++ {
		t[i] += offset
	}
	// If needed, shift pitches by octaves until all are between midilo and midihi inclusive.
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
