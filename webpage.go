package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	. "github.com/Michael-F-Ellis/goht" // dot import makes sense here
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
	Render(page, &buf, 0)
	err = ioutil.WriteFile("index.html", buf.Bytes(), 0644)
	return
}

func indexBody() (body *HtmlTree) {
	header := Div(`style="text-align:center; margin-bottom:2vh;"`,
		H2("class=title", "Infinite Etudes Web App"),
		Em("", "Ear training for your fingers"),
	)
	// Etude menus:
	// Key
	var keys []interface{}
	for _, k := range keyInfo {
		value := fmt.Sprintf(`value="%s" aria-label="%s"`, k.fileName, k.uiAria)
		keys = append(keys, Option(value, k.uiName))
	}
	keySelect := Select("id=key-select", keys...)

	// Scale pattern
	var scales []interface{}
	for _, ptn := range scaleInfo { // scaleInfo is defined in server.go
		value := fmt.Sprintf(`value="%s"`, ptn.fileName)
		scales = append(scales, Option(value, ptn.uiName))
	}
	scaleSelect := Select("id=scale-select", scales...)

	// Instrument sound
	var sounds []interface{}
	for _, iinfo := range supportedInstruments {
		name := iinfo.displayName
		value := fmt.Sprintf(`value="%s"`, iinfo.name)
		sounds = append(sounds, Option(value, name))
	}
	soundSelect := Select("id=sound-select", sounds...)

	// Rythhm pattern
	var rhythms []interface{}
	for _, rhy := range []string{"steady", "advancing"} {
		name := rhy + " rhythm"
		value := fmt.Sprintf(`value="%s"`, rhy)
		rhythms = append(rhythms, Option(value, name))
	}
	rhythmSelect := Select("id=rhythm-select", rhythms...)

	// Controls
	playBtn := Button(`onclick="playStart()"`, "Play")
	stopBtn := Button(`onclick="playStop()"`, "Stop")
	downloadBtn := Button(`onclick="downloadEtude()"`, "Download")

	// Assemble everything into the body element.
	body = Body("", header,
		Div("", keySelect, scaleSelect, soundSelect, rhythmSelect),
		Div("", playBtn, stopBtn, downloadBtn),
		quickStart(),
		forTheCurious(),
		forVocalists(),
		rhythmPatterns(),
		intervalsOctavesRanges(),
		tempo(),
		variations(),
		faq(),
		biblio(),
		coda(),
	)
	return
}

func quickStart() (div *HtmlTree) {
	div = Div("",
		H3("", "For the impatient"),
		Ol("",
			Li("", `Choose a key,`),
			Li("", `Choose a scale pattern,`),
			Li("", `Choose an instrument sound,`),
			Li("", `Click 'Play' and play along.`),
		),
	)
	return
}

func forTheCurious() (div *HtmlTree) {
	heading := "For the curious"
	p1 := `Infinite Etudes generates ear/finger training
	etudes for instrumentalists and vocalists. By default, the etudes follow a simple four
	bar form: a sequence of 3 different notes is played on beats 1, 2, and 3
	and a rest on beat 4. Each bar is played four times before moving on -- so
	you have 3 chances to play the sequence after the first hearing.
	`

	p2 := `Each etude contains all possible 3-note sequences in the key for
	the chosen scale pattern. The sequences are presented in random order. New
	etudes are generated every hour. The program is called 'Infinite Etudes'
	because the number of possible orderings of the sequences easily exceeds
	the number of stars in the universe. Luckily, the goal is to learn to
	recognize and play the individual 3-note sequences. That turns out to be a
	much more reasonable task (and the infinite sequence orderings are actually
	helpful because they prevent you from relying on muscle memory.)
    `

	p3 := `So how many sequences are there? Well, there are 12 pitches in
	the Western equal-tempered octave, so there are 12 * 11 * 10 = 1320
	possible sequences of 3 different pitches.`

	p4 := `Playing them all in one sitting with 4 bars devoted to each
	sequence would take just under 3 hours at 120 bpm. I think you'd have to be
	a little crazy to do that, but who am I to rein in your passion? For the
	rest of us, breaking it down into keys and scale patterns allows practicing
	in manageable chunks.`

	p5 := `<strong>Pentatonic:</strong> If any scale can be said to be
	universal across history and cultures, this is it. This pattern is also the
	easiest because you're only dealing with 5 pitches at a time. There are 60
	possible 3-note sequences in each key.  Each etude takes 8 minutes to
	play.`

	p6 := `<strong>Chromatic Final:</strong> This one's special. It's
	composed of all the sequences that end on the note you choose with the
	'key' selector without regard for any particular scale or key. It's the longest
	and most challenging of the patterns. For any given final note, there are
	110 possible sequences. Each etude takes just under 15 minutes to play.`

	p7 := `The good news is that this pattern is most efficient way to play
	every possible sequence because there's no overlap between the etudes for
	different final notes. <strong><em>In fact, you could stop reading right
	here and just start playing this pattern with a different final note every
	day.</em></strong> In 12 days that will take you through every possible 3-note chord
	in every inversion and every possible 3 note fragment of every possible
	scale. Not bad for only 15 minutes a day.`

	p7a := `The other scale patterns don't introduce any new sequences.
	They're included because you may find them useful for developing a sense of
	how the sequences function in a tonal context.`

	p8 := `<strong>Plus Four, Plus Seven, Four and Seven:</strong> These
	patterns connect the pentatonic scale to the major scale (and its relative
	minor. In music-speak, the pentatonic scale is degrees 1,2,3,5,6 of the
	major scale that starts on the same note. So C pentatonic is C D E G A and
	C major is C D E <strong>F</strong> G A <strong>B</strong>. F and B are 4 and 7
	in C major.`

	p9 := `<strong>Plus Four</strong> contains all the sequences that
	consist of 4 and any two of 1,2,3,5,6. As there are 60 such sequences, the
	etudes are 8 minutes long. <strong>Plus Seven</strong> is analogous. It
	adds 7 instead of 4, creating another 60 sequences.`

	p10 := `<strong>Four and Seven</strong> contains all the sequences that
	contain both 4 and 7 plus one other note from the pentatonic scale.There
	are only 30 such sequences. The etudes take 4 minutes to play.
	Interestingly, these are the only sequences that can be said to exist in
	exactly one key.  I'll leave it to you to work out why that's so :-)`

	p11 := `<strong>Harmonic Minor 1</strong> and <strong>Harmonic Minor 2</strong>
    explore the relative harmonic minor scale that's common in Middle
    Eastern music. They're included because they complete the coverage of all
    possible sequences (except pairs of adjacent half-steps) in a tonal context.`

	p12 := `<strong>Harmonic Minor 1</strong> contains all the sequences (36
    total) from 1,2,3,♯5,6 that contain ♯5. It takes just under 5 minutes to play.`

	p13 := `<strong>Harmonic Minor 2</strong> contains all the sequences (55
	total) from 1,2,3,4,♯5,6,7 that contain ♯5 and one or both of 4 and 7. It takes 7:20 to play.`

	p14 := `<strong>Intervals</strong> is the simplest and shortest pattern. Each etude presents all twelve chromatic pitches
	relative to the chosen key note and takes just over 90 seconds to play.`

	div = Div("",
		H3("", heading),
		P("", p1),
		P("", p2),
		P("", p3),
		P("", p4),
		P("", "Here are the patterns."),
		P("", p14),
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
func forVocalists() (div *HtmlTree) {
	p1 := `I conceived Infinite Etudes as an aid for instrumentalists. I've since
	found it's also quite useful as a daily vocal workout for intonation. The
	instrument selection menu has choir ahh sounds for soprano, alto, tenor and bass ranges.`
	div = Div("",
		H3("", "For Vocalists"),
		P("", p1),
	)
	return
}

func rhythmPatterns() (div *HtmlTree) {
	p1 := `Infinite Etudes supports two ryhthm patterns via the rhythm select
	menu. The "steady rhythm" pattern is the default. Each sequence of three
	notes is played 4 times with the first note falling on the downbeat of a
	measure and a rest on the fourth beat.`

	p2 := `The "advancing rhythm" pattern adds some rhythmic variety. It begins
	in the same way as the steady rhythm pattern but omits the rest at the end
	of the fourth repetition. The creates a cycle so that the second sequence has the
	downbeat on the second note, the third sequence has the downbeat on the third note,
	then fourth sequence has the downbeat on a rest. The point is to allow you to hear
	and play the sequences at different starting points within a measure.`

	p3 := `You'll probably want to use the steady pattern until you're getting the
    most of the notes right on the first or second repetition.`

	div = Div("",
		H3("", "Rhythm Patterns"),
		P("", p1),
		P("", p2),
		P("", p3),
	)

	return
}

func intervalsOctavesRanges() (div *HtmlTree) {
	p1 := `What I said earlier about covering all possible sequences of
	notes needs some clarification.  First, the program puts the notes of every
	sequence in close voicing so that each note is no more than 6 semitones
	from the note that came before it.`

	p2 := `For example, the sequences 'E G C' and 'C G E' will always be
	voiced so that the E is the lowest note and the C is the highest. The same
	rule applies between the last note of a sequence and the first note of the
	next. If the program generates 'E G C' followed by 'E F D' the E in the
	second sequence will be an octave above the E in the first sequence.`

	p3 := `If you find that explanation confusing, don't worry. What's going
	on will be obvious after a few minutes of playing along. Just be aware that
	none of the sequences presented will contain leaps of a 5th (7 semitones) or larger.
	A simple way to incorporate larger leaps is to play one of the notes an octave higher
	or lower.`

	p4 := `This voicing rule has a couple of good consequences: The
	notes of each sequence will always fit within one octave and the sequences,
	being randomly chosen, will wander over the entire pitch range of your
	instrument. The normal limits of each instrument are known to the program
	and it will keep everything within the bounds of what's playable.`

	div = Div("",
		H3("", "Intervals, Octaves, Ranges"),
		P("", p1),
		P("", p4),
		P("", p2),
		P("", p3),
	)
	return
}

func tempo() (div *HtmlTree) {
	p1 := `This web demo generates MIDI files in 4/4 time with the tempo fixed
	at 120 beats per minute. If you need it slower or faster, the easiest solution
	is to use the download button to save a local copy of a file and play it with
	a program that allows you to adjust the tempo.  I recommend QMidi for Mac. I don't
	know what's good on PC but a little Googling should turn up something appropriate. Downloading
	also allows you to play the files through better equipment for more realistic sound.`

	p2 := `You might also consider installing MuseScore, the excellent open
	source notation editor. Version 3.1 does a very good job importing Infinite
	Etudes midi files. Besides controlling tempo, you can print the etude as sheet
	music or play it back with real-time highlighting of each note as it's
	played.`

	p3 := `A third option, if you have software skills, is to install Infinite
	Etudes on your computer from the source code on <a
	href="https://github.com/Michael-F-Ellis/infinite-etudes">GitHub.</a>. You
	can control both tempo and the number of octaves from the command line
	version.`

	p4 := `Having said all that, let me offer a reason not to work at a much
	slower tempo. In a jam session you're going to encounter 120 bpm and faster
	and you'll be playing in eighths or sixteenths. I know there's a saying
	among music teachers that "if you practice slowly and never make a mistake,
	you'll never make a mistake" -- and slow practice certainly has its place -- but
	that saying doesn't hold up in the research on learning. Immediate feedback
	about mistakes is the important consideration [Brown 2014, p90]. Assuming
	you can hear when two notes are in unison, you'll always know when you've
	played a sequence right or wrong with Infinite Etudes. So have at it boldly
	and smile at the notes you miss. You'll have more fun and make better progress.`
	div = Div("",
		H3("", "Tempo"),
		P("", p1),
		P("", p2),
		P("", p3),
		P("", p4),
	)
	return
}

func variations() (div *HtmlTree) {
	p1 := `As you progress, some sequences will become easy to recognize and
	play before others. When you nail a particular sequence correctly and
	confidently on first hearing (hooray!), you can put the remaining two bars
	to good use in a variety of ways. Here are a few suggestions, some simple
	and some difficult:`
	var variants = []string{
		`Finger it differently.`,
		`Change the bowing or picking.`,
		`Play it with the other hand (keyboard).`,
		`Play it in the same octave on different strings (string instruments).`,
		`Play one note up or down an octave.`,
		`Play the whole sequence up or down an octave.`,
		`Play it in both hands one or more octaves apart.`,
		`Play it up or down a fifth (or fourth, third, ...).`,
		`Play it as a chord.`,
		`Find a bass note or chord that works with the sequence.`,
		`Mess with the rhythm, accents, dynamics, timbre, ...`,
		`Shred it in sixteenth note cross-rhythm, e.g. 1231 2312 3123`,
		`Fill in between the notes.`,
		`Invent a counter-melody,`,
		`or simply take a deep breath and relax your fingers.`,
	}
	p2 := `Above all, make some music whenever possible!`
	// need []interface{} to pass strings as Li elements to Ul()
	var ivariants []interface{}
	for _, s := range variants {
		ivariants = append(ivariants, Li("", s))
	}
	div = Div("",
		H3("", "Variations"),
		P("", p1),
		Ul("", ivariants...),
		P("", p2),
	)
	return
}

func faq() (div *HtmlTree) {
	qa := func(q string, a ...string) (div *HtmlTree) {
		var item []interface{}
		item = append(item, P("", Strong("", Em("", q))))
		for _, s := range a {
			item = append(item, P("", s))
		}
		div = Div("", item...)
		return
	}
	q1 := qa(`Why 3 notes rather than 4 or 5 ...?`,
		`The math is less friendly for longer sequences. There are 11880 possible sequences of 4 notes and
	95040 sequences of 5 notes. You can get through all 1320 3-note sequences every 12 days in 15 minutes/day playing the
	Chromatic Final etudes. To do that with 4-notes sequences would take over 3 months (108 days).`,
		`With 5-note sequences it would take more than 2 years.`)

	q2 := qa(`Are 3 notes enough to be of benefit?`,
		`My own experience says 'yes'. At the piano, I've experienced
	very noticeable improvement in my ability to play by ear as well as in my sight-reading. I attribute both
	to having to devote less mental effort to fingering.`,
		`As a singer, I use the etudes as a daily exercise to work on intonation through my full vocal range.`)

	div = Div("",
		H3("", `FAQ`),
		q1,
		q2,
	)
	return
}

func coda() (div *HtmlTree) {
	p1 := `I wrote Infinite Etudes for two reasons: First, as a tool for my
	own practice on piano and viola; second as a small project to develop a
	complete application in the Go programming language. I'm happy with it on
	both counts and I hope you find it useful also. The source code is available
	on <a href="https://github.com/Michael-F-Ellis/infinite-etudes">GitHub.</a>
	<br><br>Mike Ellis<br>Weaverville NC<br>May 2019`
	div = Div("",
		H3("", "Coda"),
		P("", p1),
	)
	return
}

func cite(citation, comment string) (div *HtmlTree) {
	div = Div("",
		P("", "<em>"+citation+"</em>"),
		P("", "<small>"+comment+"</small>"),
	)
	return
}

func biblio() (div *HtmlTree) {
	div = Div("",
		H3("", "References"),
		P("", "A few good books that influenced the development of Infinite Etudes:"),

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

func indexCSS() *HtmlTree {
	return Style("", `
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
	`)
}

func indexJS() (script *HtmlTree) {
	script = Script("",
		`
		// chores at start-up
		function start() {
		  // Chrome and other browsers now disallow AudioContext until
		  // after a user action.
		  document.body.addEventListener("click", MIDIjs.resumeAudioContext);
		}

		// Read the selects and return the URL for the etude to be played or downloaded.
		function etudeURL() {
		  key = document.getElementById("key-select").value
		  if (key=="random") {
			  key=randomKey()
			  };
		  scale = document.getElementById("scale-select").value
		  sound = document.getElementById("sound-select").value
		  rhythm = document.getElementById("rhythm-select").value
		  return "/etude/" + key + "/" + scale + "/" + sound + "/" + rhythm
		}

		// Read the selects and returns a proposed filename for the etude to be downloaded.
		function etudeFileName() {
		  key = document.getElementById("key-select").value
		  scale = document.getElementById("scale-select").value
		  sound = document.getElementById("sound-select").value
		  rhythm = document.getElementById("rhythm-select").value
		  return key + "_" + scale + "_" + sound + "_" + rhythm + ".midi"
		}
		// randomKey returns a keyname chosen randomly from a list of supported
		// keys.
		function randomKey() {
			keys = ['c', 'dflat', 'd', 'eflat', 'e', 'f',
			'gflat', 'g', 'aflat', 'a', 'bflat', 'b']
			return keys[Math.floor(Math.random() * keys.length)]
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
	`)
	return
}
