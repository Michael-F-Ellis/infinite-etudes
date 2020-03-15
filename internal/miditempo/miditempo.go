package miditempo

import (
	"fmt"
	"io/ioutil"
	"os"
)

func getFileBytes(filepath string) (bytes []byte, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		err = fmt.Errorf("error reading %v: %v", file, err)
		return
	}
	defer file.Close()

	// midifiles are small, so read the whole thing into memory
	bytes, err = ioutil.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("error reading %v: %v", file, err)
		return
	}
	return
}

// GetTempo finda and returns address and value of the first midi microseconds
// per beat event in bytes.
func GetTempo(filepath string) (addr int, tempoMs uint, err error) {
	addr, tempoMs, err = getFileTempo(filepath)
	return
}
func getFileTempo(filepath string) (addr int, tempoMs uint, err error) {
	bytes, err := getFileBytes(filepath)
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}
	// tempo events start with 0xFF5103 followed by 3 bytes whose
	// value is the tempo in µsec.
	var state int // will be 5 when we have the entire sequence
	for i, b := range bytes {
		switch state {
		case 0:
			if b == 0xFF {
				state = 1
			}
		case 1:
			switch b {
			case 0x51:
				state = 2
			default:
				state = 0
			}
		case 2:
			switch b {
			case 0x03:
				state = 3
			default:
				state = 0
			}
		case 3: // found it. i is now the offset of the most significant byte
			addr = i
			tempoMs = uint(b) << 16
			state = 4
		case 4:
			tempoMs += uint(b) << 8
			state = 5
		case 5:
			tempoMs += uint(b)
			return // Success!
		}
	}

	err = fmt.Errorf("tempo event not found")
	return
}

// low3 returns a 3 byte array representing the lower
// 3 bytes of n, e.g. as a 24 bit number
func low3(n uint) (u24 [3]byte) {
	u24[0] = byte((n & 0xFFFFFF) >> 16)
	u24[1] = byte((n & 0xFFFF) >> 8)
	u24[2] = byte((n & 0xFF))
	return u24
}

// SetTempo returns a new copy of the file's content with the tempo
// event altered so that its value is the requested number of microseconds
func SetTempo(filepath string, µs uint) (bytes []byte, err error) {
	if µs == 0 {
		err = fmt.Errorf("%d is too small for a midi SetTempo event value", µs)
		return
	}
	if µs > 0xFFFFFFF {
		err = fmt.Errorf("%d is too large for a midi SetTempo event value", µs)
		return
	}
	bytes, err = getFileBytes(filepath)
	if err != nil {
		return
	}
	addr, _, err := getFileTempo(filepath)
	if err != nil {
		return
	}
	for i, b := range low3(µs) {
		bytes[i + addr] = b
	}

    return
}
