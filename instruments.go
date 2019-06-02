package main

import "fmt"

type instrumentInfo struct {
	displayName string // what we show in the UI
	gmnumber    int    // General Midi Sound number (1-indexed)
	name        string // used in file names
	midilo      int    // lowest midi pitch to be used
	midihi      int    // highest midi pitch to be used
}

// getSupportedInstrumentByName returns the instrumentInfo
// struct that matches the name argument. It returns a non=nil
// error if no match is found.
func getSupportedInstrumentByName(name string) (iInfo instrumentInfo, err error) {
	for _, i := range supportedInstruments {
		if i.name == name {
			iInfo = i
			return
		}
	}
	// not found if we get to here
	err = fmt.Errorf("%s is not a supported instrument", name)
	return
}

// Here are the instruments we currently support.
var supportedInstruments = []instrumentInfo{
	{
		displayName: "Bass, Acoustic",
		gmnumber:    33,
		name:        "acoustic_bass",
		midilo:      28,
		midihi:      55,
	},
	{
		displayName: "Bass, Electric",
		gmnumber:    34,
		name:        "electric_bass_finger",
		midilo:      28,
		midihi:      67,
	},
	{
		displayName: "Bassoon",
		gmnumber:    71,
		name:        "bassoon",
		midilo:      34,
		midihi:      72,
	},
	{
		displayName: "Cello",
		gmnumber:    43,
		name:        "cello",
		midilo:      36,
		midihi:      72,
	},
	{
		displayName: "Clarinet",
		gmnumber:    72,
		name:        "clarinet",
		midilo:      50,
		midihi:      79,
	},
	{
		displayName: "Flute",
		gmnumber:    74,
		name:        "flute",
		midilo:      60,
		midihi:      98,
	},
	{
		displayName: "Guitar, Acoustic",
		gmnumber:    26,
		name:        "acoustic_guitar_steel",
		midilo:      40,
		midihi:      76,
	},
	{
		displayName: "Guitar, Electric",
		gmnumber:    27,
		name:        "electric_guitar_jazz",
		midilo:      40,
		midihi:      88,
	},
	{
		displayName: "Oboe",
		gmnumber:    69,
		name:        "oboe",
		midilo:      58,
		midihi:      92,
	},
	{
		displayName: "Piano",
		gmnumber:    1,
		name:        "acoustic_grand_piano",
		midilo:      36,
		midihi:      96,
	},
	{
		displayName: "Sax, Soprano",
		gmnumber:    65,
		name:        "soprano_sax",
		midilo:      56,
		midihi:      87,
	},
	{
		displayName: "Sax, Alto",
		gmnumber:    66,
		name:        "alto_sax",
		midilo:      49,
		midihi:      80,
	},
	{
		displayName: "Sax, Tenor",
		gmnumber:    67,
		name:        "tenor_sax",
		midilo:      44,
		midihi:      75,
	},
	{
		displayName: "Sax, Baritone",
		gmnumber:    68,
		name:        "baritone_sax",
		midilo:      36,
		midihi:      68,
	},
	{
		displayName: "Trombone",
		gmnumber:    58,
		name:        "trombone",
		midilo:      40,
		midihi:      77,
	},
	{
		displayName: "Trumpet",
		gmnumber:    57,
		name:        "trumpet",
		midilo:      54,
		midihi:      86,
	},
	{
		displayName: "Violin",
		gmnumber:    41,
		name:        "violin",
		midilo:      55,
		midihi:      91,
	},
	{
		displayName: "Viola",
		gmnumber:    42,
		name:        "viola",
		midilo:      48,
		midihi:      84,
	},
	{
		displayName: "Vocal, Soprano",
		gmnumber:    53,
		name:        "choir_aahs_soprano",
		midilo:      60,
		midihi:      84,
	},
	{
		displayName: "Vocal, Alto",
		gmnumber:    53,
		name:        "choir_aahs_alto",
		midilo:      52,
		midihi:      76,
	},
	{
		displayName: "Vocal, Tenor",
		gmnumber:    53,
		name:        "choir_aahs_tenor",
		midilo:      46,
		midihi:      72,
	},
	{
		displayName: "Vocal, Bass",
		gmnumber:    53,
		name:        "choir_aahs_bass",
		midilo:      40,
		midihi:      64,
	},
}
