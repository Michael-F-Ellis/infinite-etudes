package miditempo

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var fName = "c_pentatonic_acoustic_bass_steady.mid"

func TestGetTempo(t *testing.T) {
	addr, µs, err := GetTempo(fName)
	if err != nil {
		t.Errorf("%v", err)
	}
	if µs != 500000 {
		t.Errorf("exp %d, got %d", 500000, µs)
	}
	fmt.Printf("0x%x, %d\n", addr, µs)
}

func TestSetTempo(t *testing.T) {
	µs := uint(60000000 / 100)
	bytes, err := SetTempo(fName, µs)
	if err != nil {
		t.Errorf("%v", err)
	}
	outfile := "/tmp/test.mid"
	err = ioutil.WriteFile(outfile, bytes, 0644)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer os.Remove(outfile)
	_, gotµs, err := GetTempo(outfile)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if µs != gotµs {
		t.Errorf("exp %d, got %d", µs, gotµs)
	}
}
