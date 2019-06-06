package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	. "github.com/Michael-F-Ellis/infinite-etudes/internal/ht" // dot import makes sense here
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
		Link(`rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css"`),
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
	header := Div(`style="text-align:center; margin-bottom:2vh;"`,
		H2("class=title", SC("Infinite Etudes Web Demo")),
		Em("", SC("Ear training for your fingers")),
	)
	// Etude menus:
	// Key
	keys := []Content{}
	for _, k := range keyInfo {
		value := fmt.Sprintf(`value="%s"`, k.fileName)
		keys = append(keys, Option(value, SC(k.uiName)))
	}
	keySelect := Select("id=key-select", keys...)

	// Scale pattern
	scales := []Content{}
	for _, ptn := range scaleInfo { // scaleInfo is defined in server.go
		value := fmt.Sprintf(`value="%s"`, ptn.fileName)
		scales = append(scales, Option(value, SC(ptn.uiName)))
	}
	scaleSelect := Select("id=scale-select", scales...)

	// Instrument sound
	sounds := []Content{}
	for _, iinfo := range supportedInstruments {
		name := iinfo.displayName
		value := fmt.Sprintf(`value="%s"`, iinfo.name)
		sounds = append(sounds, Option(value, SC(name)))
	}
	soundSelect := Select("id=sound-select", sounds...)

	// Controls
	playBtn := Button(`onclick="playStart()"`, SC("Play"))
	stopBtn := Button(`onclick="playStop()"`, SC("Stop"))
	downloadBtn := Button(`onclick="downloadEtude()"`, SC("Download"))

	// Assemble everything into the body element.
	body = Body("", header,
		Div("", keySelect, scaleSelect, soundSelect),
		Div("", playBtn, stopBtn, downloadBtn),
		quickStart(),
		forTheCurious(),
		intervalsOctavesRanges(),
		tempo(),
		biblio(),
		coda(),
	)
	return
}

func quickStart() (div *ElementTree) {
	div = Div("",
		H3("", SC("For the impatient")),
		Ol("",
			Li("", SC(`Choose a key,`)),
			Li("", SC(`Choose a scale pattern,`)),
			Li("", SC(`Choose an instrument sound,`)),
			Li("", SC(`Click 'Play' and play along.`)),
		),
	)
	return
}

func forTheCurious() (div *ElementTree) {
	heading := SC("For the curious")
	p1 := SC(`Infinite Etudes generates ear/finger training
	etudes for instrumentalists. All the etudes follow a simple four bar
	form: a sequence of 3 different notes is played on beats 1, 2, and 3 and a
	rest on beat 4. Each bar is played four times before moving on
	 -- so you have 3 chances to play the sequence after the first
	hearing.
	`)

	p2 := SC(`Each etude contains all possible 3-note sequences in the key for
	the chosen scale pattern. The sequences are presented in random order. New
	etudes are generated every hour. The program is called 'Infinite Etudes'
	because the number of possible orderings of the sequences easily exceeds
	the number of stars in the universe. Luckily, the goal is to learn to
	recognize and play the individual 3-note sequences. That turns out to be a
	much more reasonable task (and the infinite sequence orderings are actually
	helpful because they prevent you from relying on muscle memory.)
    `)

	p3 := SC(`So how many sequences are there? Well, there are 12 pitches in
	the Western equal-tempered octave, so there are 12 * 11 * 10 = 1320
	possible sequences of 3 different pitches.`)

	p4 := SC(`Playing them all in one sitting with 4 bars devoted to each
	sequence would take just under 3 hours at 120 bpm. I think you'd have to be
	a little crazy to do that, but who am I to rein in your passion? For the
	rest of us, breaking it down into keys and scale patterns allows practicing
	in manageable chunks.`)

	p5 := SC(`<strong>Pentatonic:</strong> If any scale can be said to be
	universal across history and cultures, this is it. This pattern is also the
	easiest because you're only dealing with 5 pitches at a time. There are 60
	possible 3-note sequences in each key.  Each etude takes 8 minutes to
	play.`)

	p6 := SC(`<strong>Chromatic Final:</strong> This one's special. It's
	composed of all the sequences that end on the note you choose with the
	'key' selector without regard for any particular scale or key. It's the longest
	and most challenging of the patterns. For any given final note, there are
	110 possible sequences. Each etude takes just under 15 minutes to play.`)

	p7 := SC(`The good news is that this pattern is most efficient way to play
	every possible sequence because there's no overlap between the etudes for
	different final notes. <strong><em>In fact, you could stop reading right
	here and just start playing this pattern with a different final note every
	day.</em></strong> In 12 days that will take you through every possible 3-note chord
	in every inversion and every possible 3 note fragment of every possible
	scale. Not bad for only 15 minutes a day.`)

	p7a := SC(`The other scale patterns don't introduce any new sequences.
	They're included because you may find them useful for developing a sense of
	how the sequences function in a tonal context.`)

	p8 := SC(`<strong>Plus Four, Plus Seven, Four and Seven:</strong> These
	patterns connect the pentatonic scale to the major scale (and its relative
	minor). In music-speak, the pentatonic scale is degrees 1,2,3,5,6 of the
	major scale that starts on the same note. So C pentatonic is C D E G A and
	C major is C D E <strong>F</strong> G A <strong>B</strong>. F and B are 4 and 7
	in C major.`)

	p9 := SC(`<strong>Plus Four</strong> contains all the sequences that
	consist of 4 and any two of 1,2,3,5,6. As there are 60 such sequences, the
	etudes are 8 minutes long. <strong>Plus Seven</strong> is analogous. It
	adds 7 instead of 4, creating another 60 sequences.`)

	p10 := SC(`<strong>Four and Seven</strong> contains all the sequences that
	contain both 4 and 7 plus one other note from the pentatonic scale.There
	are only 30 such sequences. The etudes take 4 minutes to play.
	Interestingly, these are the only sequences that can be said to exist in
	exactly one key.  I'll leave it to you to work out why that's so :-)`)

	p11 := SC(`<strong>Harmonic Minor 1</strong> and <strong>Harmonic Minor 2</strong>
    explore the relative harmonic minor scale that's common in Middle
    Eastern music. They're included because they complete the coverage of all
    possible sequences (except pairs of adjacent half-steps) in a tonal context.`)

	p12 := SC(`<strong>Harmonic Minor 1</strong> contains all the sequences (36
    total) from 1,2,3,#5,6 that contain #5. It takes just under 5 minutes to play.`)

	p13 := SC(`<strong>Harmonic Minor 2</strong> contains all the sequences (55
    total) from 1,2,3,4,#5,6,7 that contain #5 and one or both of 4 and 7., It takes 7:20 to play.`)

	div = Div("",
		H3("", heading),
		P("", p1),
		P("", p2),
		P("", p3),
		P("", p4),
		P("", SC("Here are the patterns.")),
		P("", p5),
		P("", p6),
		P("", p7),
		P("", p7a),
		P("", p8),
		P("", p9),
		P("", p10),
		P("", p11),
		P("", p12),
		P("", p13),
	)
	return
}

func intervalsOctavesRanges() (div *ElementTree) {
	p1 := SC(`What I said earlier about covering all possible sequences of
	notes needs some clarification.  First, the program puts the notes of every
	sequence in close voicing so that each note is no more than 6 semitones
	from the note that came before it.`)

	p2 := SC(`For example, the sequences 'E G C' and 'C G E' will always be
	voiced so that the E is the lowest note and the C is the highest. The same
	rule applies between the last note of a sequence and the first note of the
	next. If the program generates 'E G C' followed by 'E F D' the E in the
	second sequence will be an octave above the E in the first sequence.`)

	p3 := SC(`If you find that explanation confusing, don't worry. What's going
	on will be obvious after a few minutes of playing along. Just be aware that
	none of the sequences presented will contain leaps of a 5th (7 semitones) or larger.
	A simple way to incorporate larger leaps is to play one of the notes an octave higher
	or lower.`)

	p4 := SC(`This voicing rule has a couple of good consequences: The
	notes of each sequence will always fit within one octave and the sequences,
	being randomly chosen, will wander over the entire pitch range of your
	instrument. The normal limits of each instrument are known to the program
	and it will keep everything within the bounds of what's playable.`)

	div = Div("",
		H3("", SC("Intervals, Octaves, Ranges")),
		P("", p1),
		P("", p4),
		P("", p2),
		P("", p3),
	)
	return
}

func tempo() (div *ElementTree) {
	p1 := SC(`This web demo generates MIDI files in 4/4 time with the tempo fixed
	at 120 beats per minute. If you need it slower or faster, the easiest solution
	is to use the download button to save a local copy of a file and play it with
	a program that allows you to adjust the tempo.  I recommend QMidi for Mac. I don't
	know what's good on PC but a little Googling should turn up something appropriate. Downloading
	also allows you to play the files through better equipment for more realistic sound.`)

	p2 := SC(`You might also consider installing MuseScore, the excellent open
	source notation editor. Version 3.1 does a very good job importing Infinite
	Etudes midi files. Besides controlling tempo, you print the etude as sheet
	music or play it back with real-time highlighting of each note as it's
	played.`)
	div = Div("",
		H3("", SC("Tempo")),
		P("", p1),
		P("", p2),
	)
	return
}

func coda() (div *ElementTree) {
	p1 := SC(`I wrote Infinite Etudes for two reasons: First, as a tool for my
	own practice on piano and viola; second as a small project to develop a
	complete application in the Go programming language. I'm happy with it on
	both counts and I hope you find it useful also. The source code is available
	on <a href="https://github.com/Michael-F-Ellis/infinite-etudes">GitHub.</a>
	<br><br>Mike Ellis<br>Weaverville NC<br>May 2019`)
	div = Div("",
		H3("", SC("Coda")),
		P("", p1),
	)
	return
}

func cite(citation, comment string) (div *ElementTree) {
	div = Div("",
		P("", SC("<em>"+citation+"</em>")),
		P("", SC("<small>"+comment+"</small>")),
	)
	return
}

func biblio() (div *ElementTree) {
	div = Div("",
		H3("", SC("References")),
		P("", SC("A few good books that influenced the development of Infinite Etudes:")),

		cite(`Brown, Peter C. Make It Stick : the Science of Successful
        Learning. Cambridge, Massachusetts :The Belknap Press of Harvard University
        Press, 2014.`,
			`An exceedingly readable and practical summary of what works and doesn't
	   work for efficient learning.  I've attempted to incorporate the core
	   principles (frequent low stakes testing, interleaving and spaced
	   repetition) into the design of Infinite Etudes.`),

		cite(`Huron, David. Sweet Anticipation: Music and the Psychology of
		Expectation.  Cambridge, Massachusetts : The MIT Press, 2006`,
			`The book's central theme is a theory that a large part of what makes
		music enjoyable is a combination of satisfaction from predicting what
		comes next and delight when our predictions are occasionally confounded
		in interesting ways. Regarding the development of Infinite Etudes,
		Chapter 4 "Auditory Learning" and Chapter 10 "Tonality" were quite
		useful.`),

		cite(`Werner, Kenny. Effortless Mastery : Liberating the Master Musician Within. Innermusic Publishing, 2011`,
			`Jazz pianist Kenny Werner's autobiographical take on his own road to
		mastery through mindfulness.  The title refers the sensation of
		effortlessness that accompanies mastery rather than to some magical
		method of learning without practicing. I took from it a sense of the
		value of patience in allowing musicianship to develop.`),

		cite(`Wooten, Victor L. The Music Lesson : a Spiritual Search for
		Growth Through Music. New York, New York : The Berkley Publishing
		Group, 2008.`,
			`Pearls of musical wisdom are threaded throughout Grammy Award winning
		bassist Victor Wooten's fanciful tale of adventures with music teachers
		who show up unannounced at his home (think Carlos Castaneda without the
		hallucinogens.) Read it as a counterpoint to the emphasis on notes
		embodied in Infinite Etudes. Good music making is also about
		articulation, technique, emotion, dynamics, tempo, tone, phrasing and
		space, or as one of Wooten's teachers puts it: "Never lose the groove
		in order to find a note!"`),
	)
	return
}

func indexCSS() *ElementTree {
	return Style("", SC(`
    body {
	  margin: 0;
	  height: 100%;
	  overflow: auto;
	  background-color: #DDA;
	  }
    h1 {font-size: 300%; margin-bottom: 1vh}
    h2 {font-size: 200%}
	h2.title {text-align: center;}
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
    select {
	  margin-left: 5%;
	  margin-bottom: 1%;
	  background-color: white;
	}
    button {
	  margin-left: 5%;
	  margin-bottom: 1%;
	  background-color: #ADA;
	}
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

		// Read the selects and returns a proposed filename for the etude to be downloaded.
		function etudeFileName() {
		  key = document.getElementById("key-select").value
		  scale = document.getElementById("scale-select").value
		  sound = document.getElementById("sound-select").value
		  return key + "_" + scale + "_" + sound + ".midi"
		}

		function playStart() {
		    MIDIjs.stop()
		    MIDIjs.play(etudeURL())
		}

		function playStop() {
		    MIDIjs.stop()
		}
        
		function downloadEtude() {
		  // adapted from https://stackoverflow.com/a/49917066/426853
		  let a = document.createElement('a')
		  a.href = etudeURL()
		  a.download = etudeFileName()
		  document.body.appendChild(a)
		  a.click()
		  document.body.removeChild(a)
		}

		// Run start when the doc is fully loaded.
		document.addEventListener("DOMContentLoaded", start);
	`))
	return
}
