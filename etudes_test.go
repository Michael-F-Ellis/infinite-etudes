package main

import (
	"fmt"
	"reflect"
	"testing"
)

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
	expected = []int{1, 3, 4, 6, 8, 10, 11}
	isminor = false
	got = getScale(keynum, isminor)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("for keynum=%d, isminor=%t,  expected %v, got %v", keynum, isminor, expected, got)
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

func TestGenerateSequences(t *testing.T) {
	s := generateSequences(0, 36, 84, 120, 0)
	if len(s) != 6 {
		t.Errorf("expected 6 sequences, got %d", len(s))
	}
	if s[0].filename != "c_pentatonic_acoustic_grand_piano.mid" {
		t.Errorf("expected name of first sequence to be c_pentatonic, got %s", s[0].filename)
	}
	// verify that all 300 permutations are accounted for
	n := 0
	for _, seq := range s {
		n += len(seq.seq)
	}
	if n != 300 {
		t.Errorf("expected 210 midiTriples total, got %d", n)
	}
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
	constrain(&x, prior, 36, 84)
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}
	x = midiTriple{1, 2, 3}
	prior = 83
	exp = midiTriple{73, 74, 75}
	constrain(&x, prior, 36, 84)
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}
	x = midiTriple{0, 7, 2}
	prior = 37
	exp = midiTriple{48, 43, 38}
	constrain(&x, prior, 36, 84)
	if x != exp {
		t.Errorf("expected %v, got %v", exp, x)
	}

}

func TestMkMidi(t *testing.T) {
	var x etudeSequence
	var exp etudeSequence
	var exp2 etudeSequence
	x.seq = []midiTriple{{1, 2, 3}, midiTriple{4, 5, 6}}
	x.tempo = 120
	x.midilo = 36
	x.midihi = 84
	x.filename = "/tmp/testmkmidi.mid"
	exp.seq = []midiTriple{{61, 62, 63}, midiTriple{64, 65, 66}}
	exp2.seq = []midiTriple{{64, 65, 66}, midiTriple{61, 62, 63}}
	mkMidi(&x)
	if !(reflect.DeepEqual(x.seq, exp.seq) || reflect.DeepEqual(x.seq, exp2.seq)) {
		t.Errorf("expected %v or %v, got %v", exp.seq, exp2.seq, x.seq)
	}

}

func TestShuffle(t *testing.T) {
	var x etudeSequence
	var y etudeSequence
	x.seq = []midiTriple{{1, 2, 3}, midiTriple{4, 5, 6}, midiTriple{7, 8, 9}, midiTriple{10, 11, 12}, midiTriple{13, 14, 15}}
	y.seq = []midiTriple{{1, 2, 3}, midiTriple{4, 5, 6}, midiTriple{7, 8, 9}, midiTriple{10, 11, 12}, midiTriple{13, 14, 15}}
	shuffle(x.seq)
	if reflect.DeepEqual(x.seq, y.seq) {
		t.Errorf("shuffle did not change sequence, could be chance, so try again")
	}
	fmt.Printf("%v", x.seq)

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

func TestFourBarsMusic(t *testing.T) {
	pitches := midiTriple{1, 2, 3}
	// just the first bar
	exp := []byte{0x90, 0x01, 0x65, 0x87, 0x40, 0x80, 0x01, 0x65, 0x00, 0x90, 0x02, 0x51, 0x87, 0x40, 0x80, 0x02, 0x51, 0x00, 0x90, 0x03, 0x51, 0x87, 0x40, 0x80, 0x03, 0x51, 0x87, 0x40}
	x := fourBarsMusic(pitches)
	if !reflect.DeepEqual(x.Bytes()[:len(exp)], exp) {
		t.Errorf("expected % x, got % x", exp, x)
	}
	n := len(x.Bytes())
	if n != 4*len(exp) {
		t.Errorf("expected %d bytes, got %d", 4*len(exp), n)
	}

}
func TestMetronomeBars(t *testing.T) {
	// just the first bar
	exp := []byte{0x99, 0x4c, 0x65, 0x87, 0x40, 0x89, 0x4c, 0x65, 0x00, 0x99, 0x4d, 0x51, 0x87, 0x40, 0x89, 0x4d, 0x51, 0x00, 0x99, 0x4d, 0x51, 0x87, 0x40, 0x89, 0x4d, 0x51, 0x00, 0x99, 0x4d, 0x51, 0x87, 0x40, 0x89, 0x4d, 0x51, 0x00}
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

//func TestExtractFileNames(t *testing.T) {
//	// Good lines
//	line := "// FTEE foo.bar baz.txt"
//	expected := []string{"foo.bar", "baz.txt"}
//	names, err := extractFileNames("FTEE", line)
//	if err != nil {
//		t.Errorf("unexpected error \"%q\" parsing \"%s\"", err, line)
//	}
//	if !reflect.DeepEqual(names, expected) {
//		t.Errorf("expected %s got %s", expected, names)
//	}
//	// Output lines (no FTEE)
//	line = "lorem ipsum sit amet ..."
//	names, err = extractFileNames("FTEE", line)
//	if err != nil {
//		t.Errorf("unexpected error \"%q\" parsing \"%s\"", err, line)
//	}
//
//	// Bad lines
//	line = "//FTEE foo.bar baz.txt"
//	errexp := fmt.Errorf("delimiter FTEE must be surrounded by whitespace")
//	names, err = extractFileNames("FTEE", line)
//	if !reflect.DeepEqual(err, errexp) {
//		t.Errorf("expected %q got %q", errexp, err)
//	}
//	line = "// FTEE foo.bar FTEE baz.txt"
//	errexp = fmt.Errorf("found more than one delimiter FTEE in line")
//	names, err = extractFileNames("FTEE", line)
//	if !reflect.DeepEqual(err, errexp) {
//		t.Errorf("Expected %q got %q", errexp, err)
//	}
//	line = "// FTEE"
//	errexp = fmt.Errorf("no file names found after delimiter FTEE")
//	names, err = extractFileNames("FTEE", line)
//	if !reflect.DeepEqual(err, errexp) {
//		t.Errorf("expected %q got %q", errexp, err)
//	}
//
//}
//
//func TestOpenOutputFiles(t *testing.T) {
//	//outputs := make(map[string]*os.File)
//	names := []string{"/tmp/foo.txt", "/tmp/bar.txt"}
//	err := openOutputFiles(names)
//	if err != nil {
//		t.Errorf("unexpected error: %q", err)
//	}
//	if len(_gOutputs) != 2 {
//		t.Errorf("expected 2 opened files, got %d", len(_gOutputs))
//	}
//	err = openOutputFiles(names)
//	if err != nil {
//		t.Errorf("unexpected error: %q", err)
//	}
//	if len(_gOutputs) != 2 {
//		t.Errorf("expected 2 opened files, got %d", len(_gOutputs))
//	}
//	closeOutputFiles()
//	removeOutputFiles()
//}
//
//func BenchmarkProcessInputFile(b *testing.B) {
//	// includes file I/O
//	for n := 0; n < b.N; n++ {
//		infd, err := os.Open("bigfile.txt")
//		if err != nil {
//			err = fmt.Errorf("couldn't open input file: %q", err)
//			return
//		}
//		processInputFile(infd, "FTEE")
//		infd.Close()
//		removeOutputFiles()
//	}
//}
