package valid

import "fmt"

type InstrumentInfo struct {
	DisplayName string // what we show in the UI
	GMNumber    int    // General Midi Sound number (1-indexed)
	Name        string // used in file names
	MidiLo      int    // lowest midi pitch to be used
	MidiHi      int    // highest midi pitch to be used
}

// InstrumentName returns true if the instrument name is in the ones we
// support.
func InstrumentName(name string) (ok bool) {
	_, err := InstrumentByName(name)
	if err == nil {
		ok = true
	}
	return
}

// InstrumentByName returns the instrumentInfo
// struct that matches the name argument. It returns a non=nil
// error if no match is found.
func InstrumentByName(name string) (iInfo InstrumentInfo, err error) {
	for _, i := range Instruments {
		if i.Name == name {
			iInfo = i
			return
		}
	}
	// not found if we get to here
	err = fmt.Errorf("%s is not a supported instrument", name)
	return
}

// Here are the instruments we currently support.
var Instruments = []InstrumentInfo{
	{
		DisplayName: "Bass, Acoustic",
		GMNumber:    33,
		Name:        "acoustic_bass",
		MidiLo:      28,
		MidiHi:      55,
	},
	{
		DisplayName: "Bass, Electric",
		GMNumber:    34,
		Name:        "electric_bass_finger",
		MidiLo:      28,
		MidiHi:      67,
	},
	{
		DisplayName: "Bassoon",
		GMNumber:    71,
		Name:        "bassoon",
		MidiLo:      34,
		MidiHi:      72,
	},
	{
		DisplayName: "Cello",
		GMNumber:    43,
		Name:        "cello",
		MidiLo:      36,
		MidiHi:      72,
	},
	{
		DisplayName: "Clarinet",
		GMNumber:    72,
		Name:        "clarinet",
		MidiLo:      50,
		MidiHi:      79,
	},
	{
		DisplayName: "Flute",
		GMNumber:    74,
		Name:        "flute",
		MidiLo:      60,
		MidiHi:      98,
	},
	{
		DisplayName: "Guitar, Acoustic",
		GMNumber:    26,
		Name:        "acoustic_guitar_steel",
		MidiLo:      40,
		MidiHi:      76,
	},
	{
		DisplayName: "Guitar, Electric",
		GMNumber:    27,
		Name:        "electric_guitar_jazz",
		MidiLo:      40,
		MidiHi:      88,
	},
	{
		DisplayName: "Oboe",
		GMNumber:    69,
		Name:        "oboe",
		MidiLo:      58,
		MidiHi:      92,
	},
	{
		DisplayName: "Piano",
		GMNumber:    1,
		Name:        "acoustic_grand_piano",
		MidiLo:      36,
		MidiHi:      96,
	},
	{
		DisplayName: "Sax, Soprano",
		GMNumber:    65,
		Name:        "soprano_sax",
		MidiLo:      56,
		MidiHi:      87,
	},
	{
		DisplayName: "Sax, Alto",
		GMNumber:    66,
		Name:        "alto_sax",
		MidiLo:      49,
		MidiHi:      80,
	},
	{
		DisplayName: "Sax, Tenor",
		GMNumber:    67,
		Name:        "tenor_sax",
		MidiLo:      44,
		MidiHi:      75,
	},
	{
		DisplayName: "Sax, Baritone",
		GMNumber:    68,
		Name:        "baritone_sax",
		MidiLo:      36,
		MidiHi:      68,
	},
	{
		DisplayName: "Trombone",
		GMNumber:    58,
		Name:        "trombone",
		MidiLo:      40,
		MidiHi:      77,
	},
	{
		DisplayName: "Trumpet",
		GMNumber:    57,
		Name:        "trumpet",
		MidiLo:      54,
		MidiHi:      86,
	},
	{
		DisplayName: "Violin",
		GMNumber:    41,
		Name:        "violin",
		MidiLo:      55,
		MidiHi:      91,
	},
	{
		DisplayName: "Viola",
		GMNumber:    42,
		Name:        "viola",
		MidiLo:      48,
		MidiHi:      84,
	},
	{
		DisplayName: "Vocal, Soprano",
		GMNumber:    53,
		Name:        "choir_aahs_soprano",
		MidiLo:      60,
		MidiHi:      84,
	},
	{
		DisplayName: "Vocal, Alto",
		GMNumber:    53,
		Name:        "choir_aahs_alto",
		MidiLo:      52,
		MidiHi:      76,
	},
	{
		DisplayName: "Vocal, Tenor",
		GMNumber:    53,
		Name:        "choir_aahs_tenor",
		MidiLo:      46,
		MidiHi:      72,
	},
	{
		DisplayName: "Vocal, Bass",
		GMNumber:    53,
		Name:        "choir_aahs_bass",
		MidiLo:      40,
		MidiHi:      64,
	},
}
