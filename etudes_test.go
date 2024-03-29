package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func TestGetScale(t *testing.T) {
	keynum := 0
	expected := []int{0, 2, 4, 5, 7, 9, 11}
	isminor := false
	got := getScale(keynum, isminor)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("for keynum=%d, isminor=%t,  expected %v, got %v", keynum, isminor, expected, got)
	}

	expected = []int{0, 2, 4, 5, 8, 9, 11}
	isminor = true
	got = getScale(keynum, isminor)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("for keynum=%d, isminor=%t,  expected %v, got %v", keynum, isminor, expected, got)
	}

	keynum = 11
	expected = []int{11, 1, 3, 4, 6, 8, 10}
	isminor = false
	got = getScale(keynum, isminor)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("for keynum=%d, isminor=%t,  expected %v, got %v", keynum, isminor, expected, got)
	}

}

func TestTripleFromIntervals(t *testing.T) {
	type testcase struct {
		in  []int
		up  bool
		exp midiPattern
	}
	cases := []testcase{
		{[]int{0, 1}, true, midiPattern{0, 1, 0}},
		{[]int{0, 1}, false, midiPattern{0, -11, 0}},
		{[]int{1, 0}, false, midiPattern{1, 0, 1}},
		{[]int{7, 0}, true, midiPattern{7, 12, 7}},
		{[]int{7, 7}, true, midiPattern{7, 19, 7}},
		{[]int{7, 7}, false, midiPattern{7, -5, 7}},
	}
	for _, c := range cases {
		got := tripleFromPitchPair(c.in[0], c.in[1], c.up)
		if diff := deep.Equal(got, c.exp); diff != nil {
			t.Errorf("%v input, %v", c.in, diff)
		}
	}
}

func TestFlip(t *testing.T) {
	var (
		heads int
		tails int
	)
	for i := 0; i < 100; i++ {
		if flip() {
			heads++
		} else {
			tails++
		}
	}
	// Crude test that flip is approximately fair
	if heads < 25 || heads > 75 {
		t.Errorf("Unexpected outcome: %d heads, %d tails", heads, tails)
	}
}
func TestPermute2(t *testing.T) {
	p := permute2(getChromaticScale())
	if len(p) != 144 {
		t.Errorf("expected 132 permutations, got %d", len(p))
	}
	exp := midiPattern{0, 12, 0}
	exp2 := midiPattern{0, -12, 0}
	got := p[0]
	diff := deep.Equal(got, exp)
	diff1 := deep.Equal(got, exp2)
	if diff != nil && diff1 != nil {
		t.Errorf("expected first triple to be one of %v, %v. Got %v", exp, exp2, got)
	}
}

func TestPermute3(t *testing.T) {
	p := permute3([]int{0, 2, 4, 5, 7, 9, 11})
	if len(p) != 210 {
		t.Errorf("expected 210 permutations, got %d", len(p))
	}
	if diff := deep.Equal(p[0], midiPattern{0, 2, 4}); diff != nil {
		t.Errorf("expected first triple to be %v, got %v", [3]int{0, 2, 4}, p[0])
	}
}

func TestPermute4(t *testing.T) {
	p := permute4([]int{0, 2, 4, 5, 7, 9, 11})
	if len(p) != 840 {
		t.Errorf("expected 840 permutations, got %d", len(p))
	}
	if diff := deep.Equal(p[0], midiPattern{0, 2, 4, 5}); diff != nil {
		t.Errorf("expected first triple to be %v, got %v", [4]int{0, 2, 4, 5}, p[0])
	}
}

func TestGenerateIntervalSequences(t *testing.T) {
	req := etudeRequest{
		instrument:  "acoustic_grand_piano",
		tempo:       "120",
		repeats:     3,
		pattern:     "allintervals",
		tonalCenter: "c",
	}
	s := generateIntervalSequence(36, 84, 120, 0, req)
	if s.req.midiFilename() != "c_allintervals_acoustic_grand_piano_on_120_3_0.mid" {
		t.Errorf("expected name of first sequence to be c_intervals, got %s", s.filename)
	}
}
func TestGenerateEqualIntervalSequences(t *testing.T) {
	req := etudeRequest{
		instrument: "acoustic_grand_piano",
		tempo:      "120",
		interval1:  "unison",
		pattern:    "interval",
		repeats:    3,
	}
	s := generateEqualIntervalSequence(36, 84, 120, 0, req)
	fname_got := s.req.midiFilename()
	fname_exp := "interval_unison_acoustic_grand_piano_on_120_3_0.mid"
	if fname_got != fname_exp {
		t.Errorf("expected file name of first sequence to be %s, got %s", fname_exp, fname_got)
	}
}

// checkPatternOrders verifies the expected number of patterns that are ordered
// ascending, descending and unordered.
func checkPatternOrders(t *testing.T, aexp, dexp, uexp int, seq etudeSequence) {
	var dump bool
	ptns := seq.ptns
	nexp := aexp + dexp + uexp
	if len(ptns) != nexp {
		t.Errorf("expected %d patterns, got %d", nexp, len(ptns))
		dump = true
	}
	// Count the number of each pattern type
	var asc, desc, unord int
	for _, p := range ptns {
		switch sliceOrder(p) {
		case SliceUnordered:
			unord++
		case SliceAscending:
			asc++
		case SliceDescending:
			desc++
		}
	}
	if unord != uexp {
		t.Errorf("expected %d unordered, got %d", uexp, unord)
		dump = true
	}
	if asc != aexp {
		t.Errorf("expected %d ascending, got %d", aexp, asc)
		dump = true
	}
	if desc != dexp {
		t.Errorf("expected %d descending, got %d", dexp, desc)
		dump = true
	}
	if dump {
		log.Println(seq)
	}

}
func TestGenerateTwoIntervalSequenceAllOrders(t *testing.T) {
	s := generateTwoIntervalSequenceAllOrders(36, 84, 120, 0, "", 2, 2)
	// There should be exactly 2 ascending and 2 descending patterns and 8 unordered patterns.
	checkPatternOrders(t, 2, 2, 8, s)
}

func TestGenerateTwoIntervalSequenceUpDown(t *testing.T) {
	s := generateTwoIntervalSequenceUpDown(36, 84, 120, 0, "", 2, 2)
	// There should be exactly 6 ascending and 6 descending patterns and 0 unordered patterns.
	checkPatternOrders(t, 6, 6, 0, s)
}

func TestGenerateThreeIntervalSequenceAllOrders(t *testing.T) {
	s := generateThreeIntervalSequenceAllOrders(36, 84, 120, 0, "", 2, 2, 2)
	// There should be exactly 1 ascending and 1 descending patterns and 22 unordered patterns.
	checkPatternOrders(t, 1, 1, 22, s)
}

func TestGenerateThreeIntervalSequenceUpDown(t *testing.T) {
	s := generateThreeIntervalSequenceUpDown(36, 84, 120, 0, "", 2, 2, 2)
	// There should be exactly 6 ascending and 6 descending patterns and 22 unordered patterns.
	checkPatternOrders(t, 6, 6, 0, s)
}

func TestTighten(t *testing.T) {
	x := midiPattern{1, 2, 3}
	exp := midiPattern{1, 2, 3}
	tighten(&x)
	if diff := deep.Equal(x, exp); diff != nil {
		t.Errorf("expected %v, got %v", exp, x)
	}
	x = midiPattern{0, 11, 10}
	exp = midiPattern{0, -1, -2}
	tighten(&x)
	if diff := deep.Equal(x, exp); diff != nil {
		t.Errorf("expected %v, got %v", exp, x)
	}
	x = midiPattern{11, 0, 10}
	exp = midiPattern{11, 12, 10}
	tighten(&x)
	if diff := deep.Equal(x, exp); diff != nil {
		t.Errorf("expected %v, got %v", exp, x)
	}
}

func TestConstrain(t *testing.T) {
	x := midiPattern{1, 2, 3}
	prior := 60
	exp := midiPattern{61, 62, 63}
	constrain(&x, prior, 36, 84, false)
	if diff := deep.Equal(x, exp); diff != nil {
		t.Errorf("expected %v, got %v", exp, x)
	}
	x = midiPattern{1, 2, 3}
	prior = 83
	exp = midiPattern{73, 74, 75}
	constrain(&x, prior, 36, 84, false)
	if diff := deep.Equal(x, exp); diff != nil {
		t.Errorf("expected %v, got %v", exp, x)
	}
	x = midiPattern{0, 7, 2}
	prior = 37
	exp = midiPattern{48, 43, 38}
	constrain(&x, prior, 36, 84, false)
	if diff := deep.Equal(x, exp); diff != nil {
		t.Errorf("expected %v, got %v", exp, x)
	}
}
func TestMkMidi(t *testing.T) {
	var x etudeSequence
	var exp etudeSequence
	var exp2 etudeSequence
	x.tempo = 120
	x.midilo = 36
	x.midihi = 84
	x.req = etudeRequest{
		tonalCenter: "c",
		pattern:     "pentatonic",
		instrument:  "trumpet",
		tempo:       "120",
	}
	x.ptns = []midiPattern{{1, 2, 3}, {4, 5, 6}}
	mkMidi(&x, false)
	// verify that the pitches in both sequences have been shifted
	// modulo 12 and that they are between midihi and midilo.
	modulus := x.ptns[0][0] / 12
	if modulus*12 < 36 {
		t.Errorf("midi pitches too low: %v", x)
	}
	if modulus*12 > 72 {
		t.Errorf("midi pitches too high: %v", x)
	}
	// now translate the values back to the original
	// using the modulus we calculated
	var y []midiPattern
	for _, ptn := range x.ptns {
		var yptn midiPattern
		for _, v := range ptn {
			yptn = append(yptn, v-modulus*12)
		}
		y = append(y, yptn)
	}
	exp.ptns = []midiPattern{{1, 2, 3}, {4, 5, 6}}
	exp2.ptns = []midiPattern{{4, 5, 6}, {1, 2, 3}}
	if !(reflect.DeepEqual(y, exp.ptns) || reflect.DeepEqual(y, exp2.ptns)) {
		t.Errorf("expected %v or %v, got %v", exp.ptns, exp2.ptns, y)
	}
	x.ptns = []midiPattern{{1, 2, 3, 4}, {4, 5, 6, 7}}
	mkMidi(&x, false)
	// verify that the pitches in both sequences have been shifted
	// modulo 12 and that they are between midihi and midilo.
	modulus = x.ptns[0][0] / 12
	if modulus*12 < 36 {
		t.Errorf("midi pitches too low: %v", x)
	}
	if modulus*12 > 72 {
		t.Errorf("midi pitches too high: %v", x)
	}
	// now translate the values back to the original
	// using the modulus we calculated
	y = []midiPattern{}
	for _, ptn := range x.ptns {
		var yptn midiPattern
		for _, v := range ptn {
			yptn = append(yptn, v-modulus*12)
		}
		y = append(y, yptn)
	}
	exp.ptns = []midiPattern{{1, 2, 3, 4}, {4, 5, 6, 7}}
	exp2.ptns = []midiPattern{{4, 5, 6, 7}, {1, 2, 3, 4}}
	if !(reflect.DeepEqual(y, exp.ptns) || reflect.DeepEqual(y, exp2.ptns)) {
		t.Errorf("expected %v or %v, got %v", exp.ptns, exp2.ptns, y)
	}

}

func TestShuffle(t *testing.T) {
	var x etudeSequence
	var y etudeSequence
	x.ptns = []midiPattern{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}}
	y.ptns = []midiPattern{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}}
	shufflePatterns(x.ptns)
	if reflect.DeepEqual(x.ptns, y.ptns) {
		t.Errorf("shuffle did not change sequence, could be chance, so try again")
	}
}

func TestGMSoundName(t *testing.T) {
	x, err := gmSoundName(99)
	exp := "FX 4 (atmosphere)"
	if err != nil {
		t.Errorf("lookup failed: %v", err)
	} else if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}
}

func TestGMSoundFileNamePrefix(t *testing.T) {
	exp := "fx_4_atmosphere"
	x := gmSoundFileNamePrefix("FX 4 (atmosphere)")
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}
}

func TestQuadNormalRhythm(t *testing.T) {
	pitches := []midiPattern{{1, 2, 3, 4}, {4, 5, 6, 7}}
	exp := []byte{
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00,
		0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00,
		0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x00,
		0x90, 0x04, 0x51, 0x87, 0x40, 0x80, 0x04, 0x51, 0x00,
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00,
		0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00,
		0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x00,
		0x90, 0x04, 0x51, 0x87, 0x40, 0x80, 0x04, 0x51, 0x00,
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00,
		0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00,
		0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x00,
		0x90, 0x04, 0x51, 0x87, 0x40, 0x80, 0x04, 0x51, 0x00,
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00,
		0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00,
		0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x00,
		0x90, 0x04, 0x51, 0x87, 0x40, 0x80, 0x04, 0x51, 0x00,
	}
	x := nBarsMusic(pitches[0], &etudeRequest{repeats: 3})
	if diff := deep.Equal(x.Bytes()[:], exp); diff != nil {
		t.Errorf("%v", diff)
	}
	n := len(x.Bytes())
	if n != len(exp) {
		t.Errorf("expected %d bytes, got %d", 4*len(exp), n)
	}

}
func TestNBarsNormalRhythm(t *testing.T) {
	pitches := []midiPattern{{1, 2, 3}, {4, 5, 6}}
	oneBarMidi := []byte{
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00,
		0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00,
		0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40,
	}
	exp := []byte{}
	exp = append(exp, oneBarMidi...)
	for n := 2; n < 5; n++ {
		exp = append(exp, oneBarMidi...)

		got := nBarsMusic(pitches[0], &etudeRequest{repeats: n - 1})
		if diff := deep.Equal(got.Bytes()[:], exp); diff != nil {
			t.Errorf("%d: %v", n, diff)
		}
	}
}

func TestMetronomeBars(t *testing.T) {
	//  Metronome On
	exp := []byte{
		0x99, 0x4c, 0x30, 0x87, 0x40, 0x89, 0x4c, 0x30, 0x00,
		0x99, 0x4d, 0x10, 0x87, 0x40, 0x89, 0x4d, 0x10, 0x00,
		0x99, 0x4d, 0x10, 0x87, 0x40, 0x89, 0x4d, 0x10, 0x00,
		0x99, 0x4d, 0x10, 0x87, 0x40, 0x89, 0x4d, 0x10, 0x00}
	x := metronomeBars(4, &(etudeRequest{metronome: metronomeOn}))
	if diff := deep.Equal(x.Bytes()[:len(exp)], exp); diff != nil {
		t.Errorf("%v", diff)
	}
	n := len(x.Bytes())
	if n != 4*len(exp) {
		t.Errorf("expected %d bytes, got %d", 4*len(exp), n)
	}
	exp2 := []byte{
		0x99, 0x4c, 0x30, 0x87, 0x40, 0x89, 0x4c, 0x30, 0x00,
		0x99, 0x4d, 0x00, 0x87, 0x40, 0x89, 0x4d, 0x00, 0x00, // silent
		0x99, 0x4d, 0x00, 0x87, 0x40, 0x89, 0x4d, 0x00, 0x00, // silent
		0x99, 0x4d, 0x00, 0x87, 0x40, 0x89, 0x4d, 0x00, 0x00, // silent
	}
	x = metronomeBars(4, &(etudeRequest{metronome: metronomeDownbeatOnly}))
	if diff := deep.Equal(x.Bytes()[:len(exp2)], exp2); diff != nil {
		t.Errorf("%v", diff)
	}
	n = len(x.Bytes())
	if n != 4*len(exp) {
		t.Errorf("expected %d bytes, got %d", 4*len(exp), n)
	}
	exp3 := []byte{
		0x99, 0x4c, 0x00, 0x87, 0x40, 0x89, 0x4c, 0x00, 0x00,
		0x99, 0x4d, 0x00, 0x87, 0x40, 0x89, 0x4d, 0x00, 0x00,
		0x99, 0x4d, 0x00, 0x87, 0x40, 0x89, 0x4d, 0x00, 0x00,
		0x99, 0x4d, 0x00, 0x87, 0x40, 0x89, 0x4d, 0x00, 0x00,
	}
	x = metronomeBars(4, &(etudeRequest{metronome: metronomeOff}))
	if diff := deep.Equal(x.Bytes()[:len(exp3)], exp3); diff != nil {
		t.Errorf("%v", diff)
	}
	n = len(x.Bytes())
	if n != 4*len(exp) {
		t.Errorf("expected %d bytes, got %d", 4*len(exp), n)
	}

}

func TestKeySignature(t *testing.T) {
	exp := []byte{0x0, 0xFF, 0x59, 0x02, 0xfe, 0x00}
	s := etudeSequence{keyname: "bflat"}
	x := keySignature(&s)
	if !reflect.DeepEqual(x, exp) {
		t.Errorf("expected %v, got %v", exp, x)
	}
}

func TestTrackInstrument(t *testing.T) {
	exp := []byte{0x00, 0xC0, 0x00}
	s := etudeSequence{instrument: 0}
	x := trackInstrument(&s)
	if !reflect.DeepEqual(x, exp) {
		t.Errorf("expected %v, got %v", exp, x)
	}
	s.instrument = 41 // viola
	exp = []byte{0x00, 0xC0, 0x29}
	x = trackInstrument(&s)
	if !reflect.DeepEqual(x, exp) {
		t.Errorf("expected %v, got %v", exp, x)
	}

}
func TestQuadFrom3Intervals(t *testing.T) {
	got := quadFrom3Intervals(60, 4, 3, 3) // C E G from middle C
	exp := midiPattern{60, 64, 67, 70}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
func TestTripleFrom2Intervals(t *testing.T) {
	got := tripleFrom2Intervals(60, 4, 3) // C E G from middle C
	exp := midiPattern{60, 64, 67}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
func TestShufflePatternPitches(t *testing.T) {
	// loop until we see all 6 possible orders
	m := map[[3]int]int{
		{1, 2, 3}: 0,
		{1, 3, 2}: 0,
		{2, 1, 3}: 0,
		{2, 3, 1}: 0,
		{3, 1, 2}: 0,
		{3, 2, 1}: 0,
	}
	var done bool
	for i := 0; i < 1000; i++ {
		done = true // assumption
		trpl := midiPattern{1, 2, 3}
		shufflePatternPitches(&trpl)
		var key [3]int
		copy(key[:], trpl)
		m[key] += 1
		for _, n := range m {
			if n == 0 {
				done = false
				break
			}
		}
		if done {
			goto SUCCESS
		}
	}
	t.Errorf("%v", m)
SUCCESS:
	// log.Println(m)
}
func TestIntervalPairEtude(t *testing.T) {
	// generate a midi file with root position major triads
	s := generateTwoIntervalSequenceAllOrders(36, 84, 120, 0, "", 4, 3)
	mkMidi(&s, false) // steady rhythm, no tighten
}
func TestExtractIntervalPair(t *testing.T) {
	type testcase struct {
		s  string
		ok bool // got two intervals, both valid
	}
	tcs := []testcase{
		{"3-4", true},    // ok
		{"3", false},     // too few intervals
		{"3-4-5", false}, // too many intervals
		{"x-4", false},   // bad i1
		{"3-x", false},   // bad i2
		{"0-4", false},   // i1 too low
		{"3-13", false},  // i2 too high
	}
	for _, tc := range tcs {
		_, _, err := extractIntervalPair(tc.s)
		switch tc.ok {
		case true:
			if err != nil {
				t.Errorf("for input %s: %v", tc.s, err)
			}
		case false:
			if err == nil {
				t.Errorf("input %s should have yielded an error", tc.s)
			}
		}

	}
}

func TestReverse(t *testing.T) {
	type testcase struct {
		s   []int
		exp []int
	}
	tcases := []testcase{
		{[]int{1, 2, 3}, []int{3, 2, 1}},
	}
	for _, tc := range tcases {
		Reverse(tc.s)
		if diff := deep.Equal(tc.s, tc.exp); diff != nil {
			t.Errorf("%v", diff)
			continue
		}
	}
}

type SliceOrder int
type Numeric interface {
	int | uint |
		uint8 | uint16 | uint32 | uint64 |
		int8 | int16 | int32 | int64 |
		float32 | float64
}

const (
	SliceDescending = iota - 1
	SliceUnordered
	SliceAscending
)

// sliceOrder inspects a numeric slice and returns SliceOrder constant to
// indicate whether it is ascending or descending or unordered.
func sliceOrder[T ~[]E, E Numeric](s T) SliceOrder {
	a := s[0]
	// See if it is strictly ascending
	for i, b := range s[1:] {
		switch {
		case a == b:
			return SliceUnordered
		case a > b:
			if i == 0 {
				goto checkDescending
			} else {
				return SliceUnordered
			}
		}
		a = b
	}
	// if we get to here, it's ascending
	return SliceAscending

checkDescending:
	for _, b := range s[1:] {
		if a <= b {
			return SliceUnordered
		}
		a = b
	}
	return SliceDescending

}
func TestSliceOrder(t *testing.T) {
	type testcase struct {
		s   []int
		exp SliceOrder
	}
	tcases := []testcase{
		{[]int{1, 2, 3}, SliceAscending},
		{[]int{3, 2, 1}, SliceDescending},
		{[]int{1, 2, 0}, SliceUnordered},
		{[]int{1, 2, 2}, SliceUnordered},
	}
	for _, tc := range tcases {
		got := sliceOrder(tc.s)
		if diff := deep.Equal(got, tc.exp); diff != nil {
			t.Errorf("%v", diff)
			continue
		}
	}
}
func TestRangeCheck(t *testing.T) {
	type testcase struct {
		req etudeRequest
		exp bool
		err error
	}
	tcases := []testcase{
		// valid instrument, small pattern
		{
			req: etudeRequest{
				instrument: "acoustic_grand_piano",
				pattern:    "interval",
				interval1:  "major3",
			},
			exp: true,
			err: nil,
		},
		// invalid instrument should return an error
		{
			req: etudeRequest{
				instrument: "bobonker",
				pattern:    "interval",
				interval1:  "major3",
			},
			exp: true,
			err: fmt.Errorf(""),
		},
		// valid instrument, pattern too large
		{
			req: etudeRequest{
				instrument: "melodica",
				pattern:    "intervaltriple",
				interval1:  "octave",
				interval2:  "octave",
				interval3:  "octave",
			},
			exp: false,
			err: nil,
		},
	}
	for _, tc := range tcases {
		got, err := rangeCheck(tc.req)
		if err != nil {
			if tc.err == nil {
				t.Errorf("unexpected error: %v", err)
			}
			continue
		}
		if diff := deep.Equal(got, tc.exp); diff != nil {
			t.Errorf("%v", diff)
		}
	}
}
func TestPitchName(t *testing.T) {
	type testcase struct {
		p      int
		sharps bool
		exp    string
	}
	tcases := []testcase{
		{61, true, "C♯4"},
		{61, false, "D♭4"},
	}
	for _, tc := range tcases {
		got := pitchName(tc.p, tc.sharps)
		if diff := deep.Equal(got, tc.exp); diff != nil {
			t.Errorf("%v", diff)
			continue
		}
	}
}
func TestStandardTempii(t *testing.T) {
	type testcase struct {
		exp   int // value we expect to find at index
		index int // position in tempii slice
	}
	tcases := []testcase{
		{40, 0},
		{63, 11},
		{76, 15},
		{126, 27},
		{152, 31},
		{220, 39},
		{336, 49},
		{480, 58},
	}
	tempii := standardTempii()
	for _, tc := range tcases {
		got := tempii[tc.index]
		if diff := deep.Equal(got, tc.exp); diff != nil {
			t.Errorf("index %d: %v", tc.index, diff)
			continue
		}
	}
}
