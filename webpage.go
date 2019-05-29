package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	. "github.com/Michael-F-Ellis/infinite-etudes/internal/ht"
)

// mkWebPages constructs the application web pages in the current working
// directory.
func mkWebPages() (err error) {
	err = mkIndex()
	return
}

func mkIndex() (err error) {
	var buf bytes.Buffer
	// <head>
	head := Head("",
		Meta(`name="viewport" content="width=device-width, initial-scale=1"`),
		Meta(`name="description", content="Infinite Etudes demo"`),
		Meta(`name="keywords", content="music,notation,midi,tbon"`),
		indexCSS(),
		indexJS(), // js for this page
		// js midi libraries
		Script("src=/midijs/libtimidity.js charset=UTF-8"),
		Script("src=/midijs/midi.js charset=UTF-8"),
	)

	// <html>
	page := Html("", head, indexBody())
	page.Render(&buf, 0)
	err = ioutil.WriteFile("index.html", buf.Bytes(), 0644)
	return
}

func indexBody() (body *ElementTree) {
	header := H2("", SC("Infinite Etudes"))
	keys := []Content{}
	for _, k := range keyNames {
		keys = append(keys, Option("", SC(k)))
	}
	keySelect := Select("id=key-select", keys...)

	scales := []Content{}
	patterns := []string{"pentatonic", "final"}
	for _, ptn := range patterns {
		scales = append(scales, Option("", SC(ptn)))
	}
	scaleSelect := Select("id=scale-select", scales...)

	sounds := []Content{}
	for _, iinfo := range supportedInstruments {
		name := iinfo.displayName
		value := fmt.Sprintf(`value="%s"`, iinfo.name)
		sounds = append(sounds, Option(value, SC(name)))
	}
	soundSelect := Select("id=sound-select", sounds...)

	playBtn := Button(`onclick="playToggle()"`, SC("Play"))
	body = Body("", header, keySelect, scaleSelect, soundSelect, playBtn)
	return
}

func indexCSS() *ElementTree {
	return Style("", SC(`
    body {margin: 0; height: 100%; overflow: hidden}
    h1 {font-size: 300%; margin-bottom: 1vh}
    h2 {font-size: 200%}
    h3 {font-size: 150%; margin-left: 2vw}
    h4 {
        font-size: 120%;
        margin-left: 2vw;
        margin-top: 1vw;
        margin-bottom: 1vw;
    }
    p {
        font-size: 100%;
        margin-left: 5%;
        margin-right: 10%;
        margin-top: 1%;
        margin-bottom: 1%;
    }
    img.example {
        margin-left: 5%;
        margin-right: 10%;
        width: 85vw;
    }
    select {margin-left: 5%; margin-bottom: 1%}
    a {font-size: 100%}
    button.nav {
        font-size: 120%;
        margin-right: 1%;
        background-color: #DFD;
    }
    input {font-size: 100%}
    li {
        font-size: 100%;
        margin-left: 5%;
        margin-right: 10%;
        margin-bottom: 0.5%;
    }
    pre {font-size: 75%; margin-left: 5%}
	/* hover color for buttons */
    input[type=submit]:hover {background-color: #0a0}
    input[type=button]:hover {background-color: #0a0}
	`))
}

func indexJS() (script *ElementTree) {
	script = Script("",
		SC(`
		var playing = false
		// chores at start-up
		function start() {
		  // Chrome and other browsers now disallow AudioContext until
		  // after a user action.
		  document.body.addEventListener("click", MIDIjs.resumeAudioContext);
		}

		// Read the selects and return the URL for the etude to be played or downloaded.
		function etudeURL() {
		  key = document.getElementById("key-select").value
		  scale = document.getElementById("scale-select").value
		  sound = document.getElementById("sound-select").value
		  return "/etude/" + key + "/" + scale + "/" + sound
		}

		function playToggle() {
		  if (playing==true) {
		    MIDIjs.stop()
			playing = false
		  } else {
		    MIDIjs.play(etudeURL())
			playing = true
		  }
		}

		// Run start when the doc is fully loaded.
		document.addEventListener("DOMContentLoaded", start);
	`))
	return
}
