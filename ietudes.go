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
	"fmt"
	"math/rand"
	"os"
)

type midiPattern []int

type etudeSequence struct {
	seq        []midiPattern
	midilo     int
	midihi     int
	tempo      int
	instrument int
	keyname    string
	filename   string
	req        etudeRequest
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

// tripleFrom2Intervals builds a midiTripe from an initial pitch, p, and a pair
// of intervals, i1 and i2, (measured in half steps). The construction always
// ascends from p.
func tripleFrom2Intervals(p, i1, i2 int) (t midiPattern) {
	q := p + i1
	r := q + i2
	t = midiPattern{p, q, r}
	return
}

// quadFrom3Intervals builds a 4 element midiPattern from an initial pitch, p, and three
// intervals, i1, i2, i3, (measured in half steps). The construction always
// ascends from p.
func quadFrom3Intervals(p, i1, i2, i3 int) (t midiPattern) {
	q := p + i1
	r := q + i2
	s := r + i3
	t = midiPattern{p, q, r, s}
	return
}

// tripleFromPitchPair builds a midiTriple from a pair of pitch numbers.
// The up argument determines if second pitch is to be adjusted an octave up or down, depending on the
// interval.
func tripleFromPitchPair(p, q int, up bool) (t midiPattern) {
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
	t = midiPattern{p, q, p}
	return
}

// flip simulates a fair coin flip
func flip() (up bool) {
	if rand.Intn(2) == 1 {
		up = true
	}
	return
}

// permute2 returns a slice of midiPattern. Each element is built from an interval
// with the first note repeated as the third note, e.g. {3, 5, 3} or {3, -7, 3}
func permute2(scale []int) (permutations []midiPattern) {
	for _, p := range scale {
		for i, q := range scale {
			// append midiPatterns alternately inverting them.
			permutations = append(permutations, tripleFromPitchPair(p, q, i%2 == 0))
		}
	}
	return
}

// permute3 returns a slice of midiPattern containing
// all possible permutations of 3 distinct notes in the
// scale.
func permute3(scale []int) []midiPattern {
	var permutations []midiPattern
	for i, p := range scale {
		for j, q := range scale {
			if j == i {
				continue
			}
			for k, r := range scale {
				if k == i || k == j {
					continue
				}
				t := midiPattern{p, q, r}
				permutations = append(permutations, t)
			}
		}
	}
	return permutations
}

// permute4 returns a slice of midiPattern containing all
// possible permutations of 4 distinct notes in the scale
func permute4(scale []int) []midiPattern {
	var permutations []midiPattern
	for i, p := range scale {
		for j, q := range scale {
			if j == i {
				continue
			}
			for k, r := range scale {
				if k == i || k == j {
					continue
				}
				for l, s := range scale {
					if l == i || l == j || l == k {
						continue
					}
					t := midiPattern{p, q, r, s}
					permutations = append(permutations, t)
				}
			}
		}
	}
	return permutations

}

// generateEqualIntervalSequences returns a slice of etudeSequences as described in the usage instructions.
// Each sequence consists of triples of equal interval sizes
func generateEqualIntervalSequences(midilo int, midihi int, tempo int, instrument int, req etudeRequest) (sequences []etudeSequence) {
	// Get the chromatic scale as midi numbers in the range 0 - 11
	midiChromaticScaleNums := getChromaticScale()
	// Generate all intervals
	triples := permute2(midiChromaticScaleNums)
	// include unisons (permute2 doesn't generate them)
	for n := range midiChromaticScaleNums {
		triples = append(triples, midiPattern{n, n, n})
	}

	var intervalNames []string
	for _, iinfo := range intervalInfo {
		intervalNames = append(intervalNames, iinfo.fileName)
	}

	// construct the sequences
	for i := 0; i < len(intervalNames); i++ {
		sequences = append(sequences, etudeSequence{
			midilo:     midilo,
			midihi:     midihi,
			tempo:      tempo,
			instrument: instrument,
			req:        req,
		})
		sequences[i].req.interval1 = intervalNames[i]
		sequences[i].req.pattern = "interval"

	}

	// filter the triples into the corresponding etude sequences
	for _, t := range triples {
		diff := t[0] - t[1]
		if diff < 0 {
			diff = -diff
		}
		// diff -= 1
		sequences[diff].seq = append(sequences[diff].seq, t)
	}
	return
}

// generateIntervalSequences returns a slice of 12 etudeSequences as described in the usage instructions.
// Each sequence consists of 12 triples with the middle pitch corresponding to pitchnum.
func generateIntervalSequences(midilo int, midihi int, tempo int, instrument int, req etudeRequest) (sequences []etudeSequence) {
	// Get the chromatic scale as midi numbers in the range 0 - 11
	midiChromaticScaleNums := getChromaticScale()
	// Generate all intervals
	triples := permute2(midiChromaticScaleNums)
	// include unisons (permute2 doesn't generate them)
	for n := range midiChromaticScaleNums {
		triples = append(triples, midiPattern{n, n, n})
	}
	// construct the sequences
	for pitch := 0; pitch < 12; pitch++ {
		pitchname := keyNames[pitch]
		sequences = append(sequences, etudeSequence{
			midilo:     midilo,
			midihi:     midihi,
			tempo:      tempo,
			instrument: instrument,
			keyname:    pitchname,
			req:        req,
		})
		sequences[pitch].req.tonalCenter = pitchname
		sequences[pitch].req.pattern = "allintervals"
	}

	// filter the triples into the corresponding etude sequences
	for _, t := range triples {
		pitch := t[0] % 12
		if pitch < 0 {
			pitch += 12
		}
		sequences[pitch].seq = append(sequences[pitch].seq, t)
	}
	return
}

// generateTwoIntervalSequence returns an etudeSequence with 12 triples of
// equal interval sizes, one beginning on each pitch in the Chromatic scale.
func generateTwoIntervalSequence(midilo int, midihi int, tempo int, instrument int, iname string, i1, i2 int) (sequence etudeSequence) {
	// Get the chromatic scale as midi numbers in the range 0 - 11
	midiChromaticScaleNums := getChromaticScale()
	// Generate all triples
	patterns := []midiPattern{}
	for p := range midiChromaticScaleNums {
		t := tripleFrom2Intervals(p, i1, i2)
		shufflePatternPitches(&t) // randomize the order
		patterns = append(patterns, t)
	}
	indices := permute3([]int{0, 1, 2})   // 6 possible note orders
	indices = append(indices, indices...) // double the list
	shufflePatterns(indices)              // shuffle the note orders
	// now rearrange the pattern pitches using the list of shuffled note orders to
	// guarantee each possible note order will appear exactly twice.
	for i, p := range patterns {
		ptn := make(midiPattern, 3)
		copy(ptn, patterns[i])
		idx := indices[i]
		for j := range p {
			ptn[j] = p[idx[j]]
		}
		patterns[i] = ptn
	}

	// construct the sequence
	sequence = etudeSequence{
		seq:        patterns,
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
	}
	return
}

// generateThreeIntervalSequence returns an etudeSequence with 12 quads of
// equal interval sizes, one beginning on each pitch in the Chromatic scale.
func generateThreeIntervalSequence(midilo int, midihi int, tempo int, instrument int, iname string, i1, i2, i3 int) (sequence etudeSequence) {
	// Get the chromatic scale as midi numbers in the range 0 - 11
	midiChromaticScaleNums := getChromaticScale()
	// Generate all triples
	patterns := []midiPattern{}
	for p := range midiChromaticScaleNums {
		q := quadFrom3Intervals(p, i1, i2, i3)
		patterns = append(patterns, q, q)
	}
	indices := permute4([]int{0, 1, 2, 3}) // 24 possible note orders
	shufflePatterns(indices)               // shuffle the note orders
	// now rearrange the pattern pitches using the list of shuffled note orders to
	// guarantee each possible note order will appear exactly once.
	for i, p := range patterns {
		ptn := make(midiPattern, 4)
		copy(ptn, patterns[i])
		idx := indices[i]
		for j := range p {
			ptn[j] = p[idx[j]]
		}
		patterns[i] = ptn
	}

	// construct the sequence
	sequence = etudeSequence{
		seq:        patterns,
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
	}
	return
}

// generateFinalSequences returns a slice of 12 etudeSequences as described in the usage instructions.
// Each sequence consists of all possible triples with a final pitch corresponding to pitchnum.
func generateFinalSequences(midilo int, midihi int, tempo int, instrument int, req etudeRequest) (sequences []etudeSequence) {
	// Get the chromatic scale as midi numbers in the range 0 - 11
	midiChromaticScaleNums := getChromaticScale()
	// Generate all 3 note permutations
	triples := permute3(midiChromaticScaleNums)

	// construct the sequences
	for pitch := 0; pitch < 12; pitch++ {
		pitchname := keyNames[pitch]
		sequences = append(sequences, etudeSequence{
			midilo:     midilo,
			midihi:     midihi,
			tempo:      tempo,
			instrument: instrument,
			keyname:    pitchname,
			req:        req,
		})
		sequences[pitch].req.tonalCenter = pitchname
		sequences[pitch].req.pattern = "final"
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
func generateKeySequences(keynum int, midilo int, midihi int, tempo int, instrument int, req etudeRequest) []etudeSequence {
	// Look up the keyname string
	keyname := keyNames[keynum]
	// Get the major and harmonic minor scales as midi numbers in the range 0 - 11
	midiMajorScaleNums := getScale(keynum, false)
	midiMinorScaleNums := getScale(keynum, true)
	// Generate all 3 note permutations
	majors := permute3(midiMajorScaleNums)
	minors := permute3(midiMinorScaleNums)

	// declare the sequences
	pentatonic := etudeSequence{
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
		req:        req,
	}
	pentatonic.req.tonalCenter = keyname
	pentatonic.req.pattern = "pentatonic"

	plusFour := etudeSequence{
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
		req:        req,
	}
	plusFour.req.tonalCenter = keyname
	plusFour.req.pattern = "plus_four"

	plusSeven := etudeSequence{
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
		req:        req,
	}
	plusSeven.req.tonalCenter = keyname
	plusSeven.req.pattern = "plus_seven"

	fourAndSeven := etudeSequence{
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
		req:        req,
	}
	fourAndSeven.req.tonalCenter = keyname
	fourAndSeven.req.pattern = "four_and_seven"

	raisedFive := etudeSequence{
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
		req:        req,
	}
	raisedFive.req.tonalCenter = keyname
	raisedFive.req.pattern = "raised_five"

	raisedFiveWithFourOrSeven := etudeSequence{
		midilo:     midilo,
		midihi:     midihi,
		tempo:      tempo,
		instrument: instrument,
		keyname:    keyname,
		req:        req,
	}
	raisedFiveWithFourOrSeven.req.tonalCenter = keyname
	raisedFiveWithFourOrSeven.req.pattern = "raised_five_with_four_or_seven"

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
func mkMidi(sequence *etudeSequence, noTighten bool) {
	// Shuffle the sequence
	shufflePatterns(sequence.seq)

	// Constrain the sequence assuming random prior pitch within the
	// instrumen's midi range.
	prior := rand.Intn(1+sequence.midihi-sequence.midilo) + sequence.midilo
	seqlen := len(sequence.seq)
	for i := 0; i < seqlen; i++ {
		t := &(sequence.seq[i])
		constrain(t, prior, sequence.midilo, sequence.midihi, noTighten)
		prior = (*t)[2]
		/*
			// for the special case of an "allintervals" request swap
			// the middle pitch (the tonic) with the first and last pitches.
			if sequence.req.pattern == "allintervals" {
				(*t)[0] = (*t)[1]
				(*t)[1] = (*t)[2]
				(*t)[2] = (*t)[0]
			}
		*/
	}
	// Write the etude
	writeMidiFile(sequence)

}

// shufflePatternPitches puts the pitches of a midiPattern in random order using
// the Fisher-Yates algorithm.
func shufflePatternPitches(t *midiPattern) {
	N := len(*t)
	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		(*t)[r], (*t)[i] = (*t)[i], (*t)[r]
	}
}

// shufflePatterns puts a slice of midiPatterns in random order using the
// Fisher-Yates algorithm.
func shufflePatterns(slc []midiPattern) {
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
	// update the filename with the rhythm pattern
	sequence.filename = sequence.req.midiFilename()
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
	for _, t := range sequence.seq {
		music := nBarsMusic(1+sequence.req.repeats, t).Bytes()
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
	bufferMusic := func(bytes []byte) {
		err = binary.Write(buf, binary.BigEndian, bytes)
		if err != nil {
			panic(err)
		}
	}
	bufferMusic([]byte{0x00})

	// one bar count-in
	countin := metronomeBars(1, &etudeRequest{metronome: metronomeOn}).Bytes()
	bufferMusic(countin)
	//
	nbars := 1 + sequence.req.repeats
	for i := 0; i < len(sequence.seq); i++ {
		music := metronomeBars(nbars, &sequence.req).Bytes()
		bufferMusic(music)
	}
	// end of track
	bufferMusic(eot)

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

// nBarsMusic returns a byte buffer containing four bars of  one midiPattern
// Two rhythm formats are supported with selection determined by the boolean
// argument named advancing. When advancing is false, the rhythm is constant with
// the pitches of each triple falling on beats 1, 2, 3 and a rest on beat 4. This
// pattern is repeated in 4 identical bars.
//
// When advancing is true, the rhythm advances by 1 beat in the 4th bar so that
// a sequence of patterns may be cycled through 4 different patterns by successively
// incrementing the offset argument.
//
// The advancing format expects the arguments t and u to be successive patterns from
// the sequence. The non-advancing format ignores u.
//
// For the final pattern in the sequence, it is recommended to pass the same value
// for t and u.
func nBarsMusic(nbars int, ptn midiPattern) *bytes.Buffer {
	// There is no valid reason to call this function with nbars < 1, so panic if that happens.
	if nbars < 1 {
		panic(fmt.Sprintf("attempted to create etude with %d bars per pattern.", nbars))
	}
	// These are the only variable length delta times we need.
	noBeats := byte(0x00)
	oneBeatHiByte := byte(0x87)
	oneBeatLoByte := byte(0x40)
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
	// mkBeat writes MIDI for one beat with note on and off events with
	// the specified pitch and velocity. If addRest is true, it appends
	// a second beat of silence.
	mkBeat := func(buf *bytes.Buffer, pitch byte, velocity byte, addRest bool) {
		var b []byte
		switch addRest {
		case false:
			b = []byte{on, pitch, velocity, oneBeatHiByte, oneBeatLoByte, off, pitch, velocity, noBeats}
		case true:
			b = []byte{on, pitch, velocity, oneBeatHiByte, oneBeatLoByte, off, pitch, velocity, oneBeatHiByte, oneBeatLoByte}
		}
		check(binary.Write(buf, binary.BigEndian, b))
	}

	// write all n bars for this triple
	for i := 0; i < nbars; i++ {
		var pitch byte
		// first beat
		pitch = byte(ptn[0])
		mkBeat(buf, pitch, velocity1, false)
		// 2nd beat
		pitch = byte(ptn[1])
		mkBeat(buf, pitch, velocity2, false)
		switch len(ptn) {
		case 3: // triple pattern
			// 3rd beat (4th beat is a rest, so we append a one beat of silence.
			pitch = byte(ptn[2])
			mkBeat(buf, pitch, velocity2, true)
		case 4: // quad pattern
			// 3rd and 4th beats
			pitch = byte(ptn[2])
			mkBeat(buf, pitch, velocity2, false)
			pitch = byte(ptn[3])
			mkBeat(buf, pitch, velocity2, false)

		}
	}
	return buf
}

// metronomeBars returns a byte buffer containing n bars of metronome click.
// Downbeats use a High Wood Block sound. Other beats use a Low Wood Block,
func metronomeBars(n int, req *etudeRequest) *bytes.Buffer {
	// These are the only variable length delta times we need.
	noBeats := byte(0x00)
	oneBeatHi := byte(0x87)
	oneBeatLo := byte(0x40)
	// adjust velocities according to request
	var velocity1, velocity2 byte
	switch req.metronome {
	case metronomeOn:
		velocity1 = byte(0x30) // downbeat
		velocity2 = byte(0x10) // other beats
		// no adjusment
	case metronomeDownbeatOnly:
		velocity1 = byte(0x30) // downbeat
		velocity2 = byte(0x00) // other beats
	case metronomeOff:
		velocity1, velocity2 = 0, 0
	default:
		panic("programming error: %d is not a supported value for etudeRequest.metronome.")
	}

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
	mkBeat := func(buf *bytes.Buffer, pitch byte, velocity byte) {
		b := []byte{on, pitch, velocity, oneBeatHi, oneBeatLo, off, pitch, velocity, noBeats}
		check(binary.Write(buf, binary.BigEndian, b))
	}

	// write as many bars as requested
	for i := 0; i < n; i++ {
		// first beat
		mkBeat(buf, wbh, velocity1)
		// 2nd beat
		mkBeat(buf, wbl, velocity2)
		// 3rd beat
		mkBeat(buf, wbl, velocity2)
		// 4th beat
		mkBeat(buf, wbl, velocity2)
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
func tighten(t *midiPattern) {
	// adjust second pitch relative to first
	(*t)[1] = adjustSuccessor((*t)[0], (*t)[1])
	// adjust third pitch relative to (adjusted) second.
	(*t)[2] = adjustSuccessor((*t)[1], (*t)[2])
}

// constrain applies tighten to a midiTriple, then adjusts its octave
// so that the first pitch is as close as possible to the last pitch
// of a previous triple. Then it checks to see if any of the adjusted
// pitches are above midihi or below midilow and re-adjusts the octave
// as needed to keep the pitches within the limits.
func constrain(t *midiPattern, prior int, midilo int, midihi int, noTighten bool) {
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
	offset := adjustSuccessor(prior, (*t)[0]) - (*t)[0]
	for i := 0; i < len(*t); i++ {
		(*t)[i] += offset
	}
	// If needed, shift pitches by octaves until all are between midilo and midihi inclusive.
	lo := int(midilo)
	// anylow tests if any pitches are too low
	anylow := func() bool {
		for _, p := range *t {
			if p < lo {
				return true
			}
		}
		return false
	}
	// adjust until none are too low
	for anylow() {
		for i := range *t {
			(*t)[i] += 12
		}
	}
	hi := int(midihi)
	// anyhigh tests if any pitches are too high
	anyhigh := func() bool {
		for _, p := range *t {
			if p > hi {
				return true
			}
		}
		return false
	}
	// adjust until none are too high
	for anyhigh() {
		for i := range *t {
			(*t)[i] -= 12
		}
	}
}
