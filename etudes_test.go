package main

import (
	"bytes"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"
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

func TestTripleFromInterval(t *testing.T) {
	type testcase struct {
		in  []int
		up  bool
		exp midiTriple
	}
	cases := []testcase{
		{[]int{0, 1}, true, midiTriple{0, 1, 0}},
		{[]int{0, 1}, false, midiTriple{0, -11, 0}},
		{[]int{1, 0}, false, midiTriple{1, 0, 1}},
		{[]int{7, 0}, true, midiTriple{7, 12, 7}},
		{[]int{7, 7}, true, midiTriple{7, 19, 7}},
		{[]int{7, 7}, false, midiTriple{7, -5, 7}},
	}
	for _, c := range cases {
		got := tripleFromPitchPair(c.in[0], c.in[1], c.up)
		if got != c.exp {
			t.Errorf("%v input, exp %v, got %v", c.in, c.exp, got)
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
	exp := midiTriple{0, 12, 0}
	exp2 := midiTriple{0, -12, 0}
	got := p[0]
	if exp != got && exp2 != got {
		t.Errorf("expected first triple to be one of %v, %v. Got %v", exp, exp2, got)
	}
}

func TestPermute3(t *testing.T) {
	p := permute3([]int{0, 2, 4, 5, 7, 9, 11})
	if len(p) != 210 {
		t.Errorf("expected 210 permutations, got %d", len(p))
	}
	if !(p[0] == [3]int{0, 2, 4}) {
		t.Errorf("expected first triple to be %v, got %v", [3]int{0, 2, 4}, p[0])
	}
}

func TestGenerateKeySequences(t *testing.T) {
	req := etudeRequest{
		instrument: "acoustic_grand_piano",
		tempo:      "120",
		rhythm:     "steady",
	}
	s := generateKeySequences(0, 36, 84, 120, 0, req)
	if len(s) != 6 {
		t.Errorf("expected 6 sequences, got %d", len(s))
	}
	if s[0].req.midiFilename() != "c_pentatonic_acoustic_grand_piano_steady_120.mid" {
		t.Errorf("expected name of first sequence to be c_pentatonic, got %s", s[0].req.midiFilename())
	}
	// verify that all 300 permutations are accounted for
	n := 0
	for _, seq := range s {
		n += len(seq.seq)
	}
	if n != 300 {
		t.Errorf("expected 300 midiTriples total, got %d", n)
	}
	// verify that no triple in Raised5 contains 4 or 7
	midiMajorScaleNums := getScale(0, false)
	four := midiMajorScaleNums[3]
	seven := midiMajorScaleNums[6]
	r5 := s[4]
	if r5.req.midiFilename() != "c_raised_five_acoustic_grand_piano_steady_120.mid" {
		t.Errorf("expected fifth sequence filename to start with 'c_raised_5', got %s", r5.req.midiFilename())
	}
	for _, x := range r5.seq {
		for i, v := range x {
			if v == four || v == seven {
				t.Errorf("raised5 triples should not contain fourth or seventh scale degrees, found %v at index %d", x, i)
			}
		}
	}
}

func TestGenerateIntervalSequences(t *testing.T) {
	req := etudeRequest{
		instrument: "acoustic_grand_piano",
		tempo:      "120",
		rhythm:     "steady",
	}
	s := generateIntervalSequences(36, 84, 120, 0, req)
	if len(s) != 12 {
		t.Errorf("expected 12 sequences, got %d", len(s))
	}
	if s[0].req.midiFilename() != "c_allintervals_acoustic_grand_piano_steady_120.mid" {
		t.Errorf("expected name of first sequence to be c_intervals, got %s", s[0].filename)
	}
	// verify that all 144 permutations are accounted for
	n := 0
	for _, seq := range s {
		n += len(seq.seq)
	}
	if n != 144 {
		t.Errorf("expected 144 midiTriples total, got %d", n)
	}
}
func TestGenerateEqualIntervalSequences(t *testing.T) {
	req := etudeRequest{
		instrument: "acoustic_grand_piano",
		tempo:      "120",
		rhythm:     "steady",
		interval1:  "minor2",
		pattern:    "interval",
	}
	s := generateEqualIntervalSequences(36, 84, 120, 0, req)
	if len(s) != 12 {
		t.Errorf("expected 12 sequences, got %d", len(s))
	}
	fname_got := s[0].req.midiFilename()
	fname_exp := "interval_minor2_acoustic_grand_piano_steady_120.mid"
	if fname_got != fname_exp {
		t.Errorf("expected file name of first sequence to be %s, got %s", fname_exp, fname_got)
	}
	// verify that all 144 permutations are accounted for
	n := 0
	for _, seq := range s {
		n += len(seq.seq)
	}
	if n != 144 {
		t.Errorf("expected 144 midiTriples total, got %d", n)
	}
}

func TestGenerateFinalSequences(t *testing.T) {
	req := etudeRequest{
		instrument: "acoustic_grand_piano",
		tempo:      "120",
		rhythm:     "steady",
	}
	s := generateFinalSequences(36, 84, 120, 0, req)
	if len(s) != 12 {
		t.Errorf("expected 12 sequences, got %d", len(s))
	}
	if s[0].req.midiFilename() != "c_final_acoustic_grand_piano_steady_120.mid" {
		t.Errorf("expected name of first sequence to be c_final, got %s", s[0].filename)
	}
	// verify that all 1320 permutations are accounted for
	n := 0
	for _, seq := range s {
		n += len(seq.seq)
	}
	if n != 1320 {
		t.Errorf("expected 1320 midiTriples total, got %d", n)
	}
}
func TestGenerateIntervalPairSequence(t *testing.T) {
	s := generateIntervalPairlSequence(36, 84, 120, 0, "", 2, 2)
	if len(s.seq) != 12 {
		t.Errorf("expected 12 triples, got %d", len(s.seq))
	}
	log.Println(s)
}

func TestTighten(t *testing.T) {
	x := midiTriple{1, 2, 3}
	exp := midiTriple{1, 2, 3}
	tighten(&x)
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}
	x = midiTriple{0, 11, 10}
	exp = midiTriple{0, -1, -2}
	tighten(&x)
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}
	x = midiTriple{11, 0, 10}
	exp = midiTriple{11, 12, 10}
	tighten(&x)
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}
}

func TestConstrain(t *testing.T) {
	x := midiTriple{1, 2, 3}
	prior := 60
	exp := midiTriple{61, 62, 63}
	constrain(&x, prior, 36, 84, false)
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}
	x = midiTriple{1, 2, 3}
	prior = 83
	exp = midiTriple{73, 74, 75}
	constrain(&x, prior, 36, 84, false)
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}
	x = midiTriple{0, 7, 2}
	prior = 37
	exp = midiTriple{48, 43, 38}
	constrain(&x, prior, 36, 84, false)
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}

}

func TestMkMidi(t *testing.T) {
	var x etudeSequence
	var exp etudeSequence
	var exp2 etudeSequence
	x.seq = []midiTriple{{1, 2, 3}, {4, 5, 6}}
	x.tempo = 120
	x.midilo = 36
	x.midihi = 84
	x.req = etudeRequest{
		tonalCenter: "c",
		pattern:     "pentatonic",
		instrument:  "trumpet",
		rhythm:      "steady",
		tempo:       "120",
	}
	exp.seq = []midiTriple{{61, 62, 63}, {64, 65, 66}}
	exp2.seq = []midiTriple{{64, 65, 66}, {61, 62, 63}}
	mkMidi(&x, false)
	if !(reflect.DeepEqual(x.seq, exp.seq) || reflect.DeepEqual(x.seq, exp2.seq)) {
		t.Errorf("expected %v or %v, got %v", exp.seq, exp2.seq, x.seq)
	}

}

func TestShuffle(t *testing.T) {
	var x etudeSequence
	var y etudeSequence
	x.seq = []midiTriple{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}}
	y.seq = []midiTriple{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}}
	shuffleTriples(x.seq)
	if reflect.DeepEqual(x.seq, y.seq) {
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

func TestComposeFileName(t *testing.T) {
	s := etudeSequence{filename: "eflat_pentatonic"}
	exp := "eflat_pentatonic_electric_grand_piano.mid"
	x := composeFileName(&s, 2)
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}

}

func TestFourBarsNormalRhythm(t *testing.T) {
	pitches := []midiTriple{{1, 2, 3}, {4, 5, 6}}
	exp := []byte{
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40,
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40,
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40,
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40,
	}
	x := fourBarsMusic(pitches[0], pitches[1], false, 0)
	if !reflect.DeepEqual(x.Bytes()[:], exp) {
		t.Errorf("expected % x, got % x", exp, x)
	}
	n := len(x.Bytes())
	if n != len(exp) {
		t.Errorf("expected %d bytes, got %d", 4*len(exp), n)
	}

}

func TestFourBarsAdvancingRhythm(t *testing.T) {
	var offset int
	var music *bytes.Buffer
	pitches := []midiTriple{{1, 2, 3}, {4, 5, 6}}
	music = fourBarsMusic(pitches[0], pitches[1], true, offset)
	exp := []byte{
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40,
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40,
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40,
		0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x00, 0x90, 0x04, 0x51, 0x87, 0x40, 0x80, 0x04, 0x51, 0x00,
	}
	if !reflect.DeepEqual(music.Bytes()[:], exp) {
		t.Errorf("offset 0:\nexp % x\ngot % x\n", exp, music)
	}
	n := len(music.Bytes())
	if n != len(exp) {
		t.Errorf("expected %d bytes, got %d", len(exp), n)
	}

	offset = 1
	music = fourBarsMusic(pitches[0], pitches[1], true, offset)
	exp = []byte{
		0x90, 0x02, 0x65, 0x87, 0x40, 0x80, 0x02, 0x65, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40, 0x90, 0x01, 0x51, 0x87, 0x40, 0x80, 0x01, 0x51, 0x00,
		0x90, 0x02, 0x65, 0x87, 0x40, 0x80, 0x02, 0x65, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40, 0x90, 0x01, 0x51, 0x87, 0x40, 0x80, 0x01, 0x51, 0x00,
		0x90, 0x02, 0x65, 0x87, 0x40, 0x80, 0x02, 0x65, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40, 0x90, 0x01, 0x51, 0x87, 0x40, 0x80, 0x01, 0x51, 0x00,
		0x90, 0x02, 0x65, 0x87, 0x40, 0x80, 0x02, 0x65, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x00, 0x90, 0x04, 0x51, 0x87, 0x40, 0x80, 0x04, 0x51, 0x00, 0x90, 0x05, 0x51, 0x87, 0x40, 0x80, 0x05, 0x51, 0x00,
	}
	if !reflect.DeepEqual(music.Bytes()[:], exp) {
		t.Errorf("offset 1:\nexp % x\ngot % x\n", exp, music)
	}
	n = len(music.Bytes())
	if n != len(exp) {
		t.Errorf("expected %d bytes, got %d", len(exp), n)
	}

	offset = 2
	music = fourBarsMusic(pitches[0], pitches[1], true, offset)
	exp = []byte{
		0x90, 0x03, 0x65, 0x87, 0x40, 0x80, 0x03, 0x65, 0x87, 0x40, 0x90, 0x01, 0x51, 0x87, 0x40, 0x80, 0x01, 0x51, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00,
		0x90, 0x03, 0x65, 0x87, 0x40, 0x80, 0x03, 0x65, 0x87, 0x40, 0x90, 0x01, 0x51, 0x87, 0x40, 0x80, 0x01, 0x51, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00,
		0x90, 0x03, 0x65, 0x87, 0x40, 0x80, 0x03, 0x65, 0x87, 0x40, 0x90, 0x01, 0x51, 0x87, 0x40, 0x80, 0x01, 0x51, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00,
		0x90, 0x03, 0x65, 0x87, 0x40, 0x80, 0x03, 0x65, 0x00, 0x90, 0x04, 0x51, 0x87, 0x40, 0x80, 0x04, 0x51, 0x00, 0x90, 0x05, 0x51, 0x87, 0x40, 0x80, 0x05, 0x51, 0x00, 0x90, 0x06, 0x51, 0x87, 0x40, 0x80, 0x06, 0x51, 0x87, 0x40,
	}
	if !reflect.DeepEqual(music.Bytes()[:], exp) {
		t.Errorf("offset 2:\nexp % x\ngot % x\n", exp, music)
	}
	n = len(music.Bytes())
	if n != len(exp) {
		t.Errorf("expected %d bytes, got %d", len(exp), n)
	}

	offset = 3
	music = fourBarsMusic(pitches[0], pitches[1], true, offset)
	exp = []byte{
		0x90, 0x01, 0x51, 0x87, 0x40, 0x80, 0x01, 0x51, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40,
		0x90, 0x01, 0x51, 0x87, 0x40, 0x80, 0x01, 0x51, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40,
		0x90, 0x01, 0x51, 0x87, 0x40, 0x80, 0x01, 0x51, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x00,
	}
	if !reflect.DeepEqual(music.Bytes()[:], exp) {
		t.Errorf("offset 3:\nexp % x\ngot % x\n", exp, music)
	}
	n = len(music.Bytes())
	if n != len(exp) {
		t.Errorf("expected %d bytes, got %d", len(exp), n)
	}

}
func TestMetronomeBars(t *testing.T) {
	// just the first bar
	exp := []byte{0x99, 0x4c, 0x30, 0x87, 0x40, 0x89, 0x4c, 0x30, 0x00, 0x99, 0x4d, 0x10, 0x87, 0x40, 0x89, 0x4d, 0x10, 0x00, 0x99, 0x4d, 0x10, 0x87, 0x40, 0x89, 0x4d, 0x10, 0x00, 0x99, 0x4d, 0x10, 0x87, 0x40, 0x89, 0x4d, 0x10, 0x00}
	x := metronomeBars(4)
	if !reflect.DeepEqual(x.Bytes()[:len(exp)], exp) {
		t.Errorf("expected % x, got % x", exp, x)
	}
	n := len(x.Bytes())
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
func TestTripleFrom2Intervals(t *testing.T) {
	got := tripleFrom2Intervals(60, 4, 3) // C E G from middle C
	exp := midiTriple{60, 64, 67}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
func TestShuffleTriplePitches(t *testing.T) {
	// loop until we see all 6 possible orders
	m := map[midiTriple]int{
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
		trpl := midiTriple{1, 2, 3}
		shuffleTriplePitches(&trpl)
		m[trpl] += 1
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
	s := generateIntervalPairlSequence(36, 84, 120, 0, "", 4, 3)
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
