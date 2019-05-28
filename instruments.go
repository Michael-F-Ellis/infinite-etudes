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
		displayName: "Bass",
		gmnumber:    33,
		name:        "acoustic_bass",
		midilo:      28,
		midihi:      55,
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
		displayName: "Electric Bass",
		gmnumber:    34,
		name:        "electric_bass_finger",
		midilo:      28,
		midihi:      67,
	},
	{
		displayName: "Electric Guitar",
		gmnumber:    27,
		name:        "electric_guitar_jazz",
		midilo:      40,
		midihi:      88,
	},
	{
		displayName: "Flute",
		gmnumber:    74,
		name:        "flute",
		midilo:      60,
		midihi:      98,
	},
	{
		displayName: "Guitar",
		gmnumber:    26,
		name:        "acoustic_guitar_steel",
		midilo:      44,
		midihi:      76,
	},
	{
		displayName: "Piano",
		gmnumber:    1,
		name:        "acoustic_grand_piano",
		midilo:      36,
		midihi:      96,
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
}
