package main

import (
	"fmt"
	"strings"
)

// Note: The instrument numbers in this file are 0-indexed.
var gmSoundNameToNum0 = map[string]int{"Acoustic Grand Piano": 0,
	"Bright Acoustic Piano":   1,
	"Electric Grand Piano":    2,
	"Honky-tonk Piano":        3,
	"Electric Piano 1":        4,
	"Electric Piano 2":        5,
	"Harpsichord":             6,
	"Clavinet":                7,
	"Celesta":                 8,
	"Glockenspiel":            9,
	"Music Box":               10,
	"Vibraphone":              11,
	"Marimba":                 12,
	"Xylophone":               13,
	"Tubular Bells":           14,
	"Dulcimer":                15,
	"Drawbar Organ":           16,
	"Percussive Organ":        17,
	"Rock Organ":              18,
	"Church Organ":            19,
	"Reed Organ":              20,
	"Accordion":               21,
	"Harmonica":               22,
	"Tango Accordion":         23,
	"Acoustic Guitar (nylon)": 24,
	"Acoustic Guitar (steel)": 25,
	"Electric Guitar (jazz)":  26,
	"Electric Guitar (clean)": 27,
	"Electric Guitar (muted)": 28,
	"Overdriven Guitar":       29,
	"Distortion Guitar":       30,
	"Guitar Harmonics":        31,
	"Acoustic Bass":           32,
	"Electric Bass (finger)":  33,
	"Electric Bass (pick)":    34,
	"Fretless Bass":           35,
	"Slap Bass 1":             36,
	"Slap Bass 2":             37,
	"Synth Bass 1":            38,
	"Synth Bass 2":            39,
	"Violin":                  40,
	"Viola":                   41,
	"Cello":                   42,
	"Contrabass":              43,
	"Tremolo Strings":         44,
	"Pizzicato Strings":       45,
	"Orchestral Harp":         46,
	"Timpani":                 47,
	"String Ensemble 1":       48,
	"String Ensemble 2":       49,
	"SynthStrings 1":          50,
	"SynthStrings 2":          51,
	"Choir Aahs":              52,
	"Voice Oohs":              53,
	"Synth Voice":             54,
	"Orchestra Hit":           55,
	"Trumpet":                 56,
	"Trombone":                57,
	"Tuba":                    58,
	"Muted Trumpet":           59,
	"French Horn":             60,
	"Brass Section":           61,
	"Synth Brass 1":           62,
	"Synth Brass 2":           63,
	"Soprano Sax":             64,
	"Alto Sax":                65,
	"Tenor Sax":               66,
	"Baritone Sax":            67,
	"Oboe":                    68,
	"English Horn":            69,
	"Bassoon":                 70,
	"Clarinet":                71,
	"Piccolo":                 72,
	"Flute":                   73,
	"Recorder":                74,
	"Pan Flute":               75,
	"Blown Bottle":            76,
	"Shakuhachi":              77,
	"Whistle":                 78,
	"Ocarina":                 79,
	"Lead 1 (square)":         80,
	"Lead 2 (sawtooth)":       81,
	"Lead 3 (calliope)":       82,
	"Lead 4 (chiff)":          83,
	"Lead 5 (charang)":        84,
	"Lead 6 (voice)":          85,
	"Lead 7 (fifths)":         86,
	"Lead 8 (bass+lead":       87,
	"Pad 1 (new age)":         88,
	"Pad 2 (warm)":            89,
	"Pad 3 (polysynth)":       90,
	"Pad 4 (choir)":           91,
	"Pad 5 (bowed)":           92,
	"Pad 6 (metallic)":        93,
	"Pad 7 (halo)":            94,
	"Pad 8 (sweep)":           95,
	"FX 1 (train)":            96,
	"FX 2 (soundtrack)":       97,
	"FX 3 (crystal)":          98,
	"FX 4 (atmosphere)":       99,
	"FX 5 (brightness)":       100,
	"FX 6 (goblins)":          101,
	"FX 7 (echoes)":           102,
	"FX 8 (sci-fi)":           103,
	"Sitar":                   104,
	"Banjo":                   105,
	"Shamisen":                106,
	"Koto":                    107,
	"Kalimba":                 108,
	"Bagpipe":                 109,
	"Fiddle":                  110,
	"Shanai":                  111,
	"Tinkle Bell":             112,
	"Agogo":                   113,
	"Steel Drums":             114,
	"Woodblock":               115,
	"Tailo Drum":              116,
	"Melodic Drum":            117,
	"Synth Drum":              118,
	"Reverse Cymbal":          119,
	"Guitar Fret Noise":       120,
	"Breath Noise":            121,
	"Seashore":                122,
	"Bird Tweet":              123,
	"Telephone Ring":          124,
	"Helicopter":              125,
	"Applause":                126,
	"Gunshot":                 127}

var gmFileNamePrefixToNum = make(map[string]int)

// Fill in the map that lets us look up midi instrument
// numbers from the alternate instrument names we use
// in etude file names.
func init() {
	for name, num := range gmSoundNameToNum0 {
		pfx := gmSoundFileNamePrefix(name)
		gmFileNamePrefixToNum[pfx] = num
	}
}

// gmSoundName looks up the sound name from the number.
// We do it with a loop since this is an infrequent operation.
func gmSoundName(num int) (string, error) {
	var err error
	for name, number := range gmSoundNameToNum0 {
		if num == number {
			return name, err
		}
	}
	// failed if we get to here
	err = fmt.Errorf("%d is not a valid GM sound number", num)
	return "", err
}

// gmSoundFileNamePrefix takes a sound name returned from
// gmSoundName and returns a clean version without spaces,
// capitals or parentheses that's suitable for use as a file
// name prefix e.g. "FX 4 (atmosphere)" -> "fx_4_atmosphere"
func gmSoundFileNamePrefix(name string) string {
	clean := strings.ToLower(name)
	clean = strings.Replace(clean, "(", "", -1)
	clean = strings.Replace(clean, ")", "", -1)
	clean = strings.Replace(clean, " ", "_", -1)
	return clean
}
