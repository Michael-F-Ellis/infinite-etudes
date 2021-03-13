// Package valid provides variables and functions used for validation and web page construction

package valid

type NameInfo struct {
	FileName string
	UiName   string
	UiAria   string // alternate text for screen readers
	Size     int    // interval size in half steps. Not meaningful for other parameters.
}

var PatternInfo = []NameInfo{
	{"interval", "One Interval", "One Interval", 0},
	{"allintervals", "Tonic Intervals", "Tonic Intervals", 0},
	{"intervalpair", "Two Intervals", "Two Intervals", 0},
	{"intervaltriple", "Three Intervals", "Three Intervals", 0},
}

// Pattern returns true if the scale name is in the ones we support.
func Pattern(name string) (ok bool) {
	for _, s := range PatternInfo {
		if s.FileName == name {
			ok = true
			break
		}
	}
	return
}

var IntervalInfo = []NameInfo{
	{"unison", "Unison", "Unison", 0},
	{"minor2", "Minor 2", "Minor Second", 1},
	{"major2", "Major 2", "Major Second", 2},
	{"minor3", "Minor 3", "Minor Third", 3},
	{"major3", "Major 3", "Major Third", 4},
	{"perfect4", "Perfect 4", "Perfect Fourth", 5},
	{"tritone", "Tritone", "Tritone", 6},
	{"perfect5", "Perfect 5", "Perfect Fifth", 7},
	{"minor6", "Minor 6", "Minor Sixth", 8},
	{"major6", "Major 6", "Major Sixth", 9},
	{"minor7", "Minor 7", "Minor Seventh", 10},
	{"major7", "Major 7", "Major Seventh", 11},
	{"octave", "Octave", "Octave", 12},
}

// IntervalName returns true if the interval name is in the ones we support.
func IntervalName(name string) (ok bool) {
	for _, k := range IntervalInfo {
		if k.FileName == name {
			ok = true
			break
		}
	}
	return
}

var KeyInfo = []NameInfo{
	{"c", "C", "C", 0},
	{"dflat", "D♭", "D-flat", 0},
	{"d", "D", "D", 0},
	{"eflat", "E♭", "E-flat", 0},
	{"e", "E", "E", 0},
	{"f", "F", "F", 0},
	{"gflat", "G♭", "G-flat", 0},
	{"g", "G", "G", 0},
	{"aflat", "A♭", "A-flat", 0},
	{"a", "A", "A", 0},
	{"bflat", "B♭", "B-flat", 0},
	{"b", "B", "B", 0},
	{"random", "Random", "Random", 0},
}

// KeyName returns true if the interval name is in the ones we support.
func KeyName(name string) (ok bool) {
	for _, k := range KeyInfo {
		if k.FileName == name {
			ok = true
			break
		}
	}
	return
}

func MetronomePattern(name string) (ok bool) {
	switch name {
	case "on", "downbeat", "off":
		ok = true
	}
	return
}
func Tempo(tBPM int) (ok bool) {
	return tBPM >= 60 && tBPM <= 480 // our aribtrary limits
}
