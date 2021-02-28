// +build mage

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	. "github.com/Michael-F-Ellis/goht" // dot import makes sense here
	"github.com/Michael-F-Ellis/infinite-etudes/internal/valid"
)

const (
	crossMark string = "&#x2717;"
	checkMark string = "&#x2713;"
)

type SilenceOption struct {
	value int    // binary mask. 1-bits are silent
	html  string // three circles (white or green) indicate which repeats are silent.
}

var silencePatterns = []SilenceOption{
	{0, checkMark + checkMark + checkMark},
	{1, checkMark + checkMark + crossMark},
	{2, checkMark + crossMark + checkMark},
	{4, crossMark + checkMark + checkMark},
	{3, checkMark + crossMark + crossMark},
	{5, crossMark + checkMark + crossMark},
	{6, crossMark + crossMark + checkMark},
	{7, crossMark + crossMark + crossMark},
}

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
	err = Render(page, &buf, 0)
	if err != nil {
		return
	}
	err = ioutil.WriteFile("assets/index.html", buf.Bytes(), 0644)
	return
}

func indexBody() (body *HtmlTree) {
	header := Div(`style="text-align:center; margin-bottom:2vh;"`,
		A(`name="top"`),
		H2("class=title", "Infinite Etudes"),
		Em("", "Ear training for your fingers"),
	)
	// Etude menus:
	// Scale pattern
	var scales []interface{}
	for _, ptn := range valid.PatternInfo { // scaleInfo is defined in server.go
		value := fmt.Sprintf(`value="%s"`, ptn.FileName)
		scales = append(scales, Option(value, ptn.UiName))
	}
	scaleSelect := Div(`class="Column" id=scale-div`, Label(``, "Pattern", Select("id=scale-select", scales...)))
	// scaleSelectLabel := Label(`class=sel-label`, "Pattern")

	// Key
	var keys []interface{}
	for _, k := range valid.KeyInfo {
		value := fmt.Sprintf(`value="%s" aria-label="%s"`, k.FileName, k.UiAria)
		keys = append(keys, Option(value, k.UiName))
	}
	keySelect := Div(`class="Column" id="key-div"`, Label(``, "Tonal Center", Select("id=key-select", keys...)))
	// Interval1 and Interval2
	var intervals []interface{}
	for _, v := range valid.IntervalInfo {
		value := fmt.Sprintf(`value="%s" aria-label="%s"`, v.FileName, v.UiAria)
		uival := fmt.Sprintf("%d (%s)", v.Size, v.UiName)
		intervals = append(intervals, Option(value, uival))
	}
	interval1Select := Div(`class="Column" id="interval1-div"`, Label(``, "Interval 1", Select("id=interval1-select", intervals...)))
	interval2Select := Div(`class="Column" id="interval2-div"`, Label(``, "Interval 2", Select("id=interval2-select", intervals...)))
	interval3Select := Div(`class="Column" id="interval3-div"`, Label(``, "Interval 3", Select("id=interval3-select", intervals...)))
	// Instrument sound
	var sounds []interface{}
	for _, iinfo := range valid.Instruments {
		name := iinfo.DisplayName
		value := fmt.Sprintf(`value="%s"`, iinfo.Name)
		sounds = append(sounds, Option(value, name))
	}
	soundSelect := Div(`class="Column" id="sound-div"`, Label(``, "Instrument", Select("id=sound-select", sounds...)))

	// Metronome
	var metros []interface{}
	for _, ptn := range []string{"on", "downbeat", "off"} {
		attrs := fmt.Sprintf(`value="%s"`, ptn)
		metros = append(metros, Option(attrs, ptn))
	}
	metroSelect := Div(`class="Column" id="metro-div"`, Label(``, "Metronome", Select("id=metro-select", metros...))) // Metronome control

	var tempos []interface{}
	var tempoValues []int
	for i := 60; i < 484; i += 4 {
		tempoValues = append(tempoValues, i)
	}
	for _, bpm := range tempoValues {
		name := fmt.Sprintf("%d", bpm)
		value := fmt.Sprintf(`value="%d"`, bpm)
		if bpm == 120 {
			value += " selected" // use 120 as the default value
		}
		tempos = append(tempos, Option(value, name))
	}
	tempoSelect := Div(`class="Column" id="tempo-div"`, Label(``, "Tempo", Select("id=tempo-select", tempos...)))

	// Repeats
	var repeats []interface{}
	for _, reps := range []string{"3", "2", "1", "0"} {
		attrs := fmt.Sprintf(`value="%s"`, reps)
		repeats = append(repeats, Option(attrs, reps))
	}
	repeatSelect := Div(`class="Column" id="repeat-div"`, Label(``, "Repeats", Select("id=repeat-select", repeats...)))

	// Silences
	var silences []interface{}
	for _, ptn := range silencePatterns {
		attrs := fmt.Sprintf(`value="%d"`, ptn.value)
		silences = append(silences, Option(attrs, ptn.html))
	}
	silenceSelect := Div(`class="Column" id="silence-div"`, Label(``, "Muting", Select("id=silence-select", silences...)))

	// Controls
	playBtn := Button(`onclick="playStart()"`, "Play")
	stopBtn := Button(`onclick="playStop()"`, "Stop")
	downloadBtn := Button(`onclick="downloadEtude()"`, "Download")

	// Assemble everything into the body element.
	body = Body("", header,
		Div(`class="Row" id="scale-row"`, scaleSelect, keySelect, interval1Select, interval2Select, interval3Select),
		Div(`class="Row"`, soundSelect, metroSelect),
		Div(`class="Row"`, tempoSelect, repeatSelect, silenceSelect),
		Div(`style="padding-top:1vh;"`, playBtn, stopBtn, downloadBtn),
		quickStart(),
		forTheCurious(),
		toTop(),
		forTheSerious(),
		toTop(),
		forBeginners(),
		toTop(),
		forVocalists(),
		toTop(),
		userInterface(),
		toTop(),
		custom(),
		toTop(),
		variations(),
		toTop(),
		// faq(),
		biblio(),
		toTop(),
		coda(),
		toTop(),
	)
	return
}
func toTop() (div *HtmlTree) {
	div = Div(``, A(`href="#top" style="color:#88F;"`, Em(``, `top`)))
	return
}
func quickStart() (div *HtmlTree) {
	div = Div("",
		H3("", "For the impatient"),
		Ol("",
			Li("", `Choose a pattern,`),
			Li("", `Choose the interval(s) (or tonal center),`),
			Li("", `Choose an instrument sound,`),
			Li("", `Click 'Play' and play along.`),
		),
		P(``, `See `, A(`href="#ui"`, "User Interface"), `for more about the
		other selectors above. To learn more about the etude patterns, see `,
			A(`href="#patterns"`, "Patterns."), `&nbsp;If you're new to your
		instrument or to ear training, see `, A(`href="#beginners"`, `For
		Beginners.`)),
	)
	return
}

func forTheCurious() (div *HtmlTree) {
	heading := "For the curious"
	p0 := `Infinite Etudes generates ear/finger training etudes for
	instrumentalists. The emphasis is on improving your ability to play what
	you hear by thoroughly exploring combinations of 2, 3 and 4 pitches over
	the full range of your instrument.`

	p0a := `Infinite Etudes doesn't try to teach music theory, melodic
	structure or how to read and write notation and it's not necessary to
	know these things to use the program, though it certainly doesn't hurt.`

	p1 := `The etudes follow a simple four bar form: a sequence of different
	notes is played on beats 1, 2, and 3 and a rest on beat 4. By default,
	each bar is played four times before moving on -- so you have 3 chances
	to play the sequence after the first hearing. You can control the number
	of repeats. You can also choose to silence one or more of the repeated
	measures.`

	p2 := `The program is called 'Infinite Etudes' because the number of
	possible orderings of the sequences is so large that you'll never play
	the same etude twice. That's an important part of the design. It prevents
	you from relying on the muscle memory that develops if you play the same
	etude repeatedly.
    `
	p3 := `Here are brief descriptions of the patterns currently supported by
	Infinite Etudes accompanied by notated examples. For brevity, the score
	examples shown here are captured with a repeat count of zero. `

	p14 := `<strong>One Interval</strong> presents 12 instances of the same
	interval pair, i.e. 3 notes, in random order. Each instance begins on a
	different pitch so that all 12 pitches are covered. <em>Note: For brevity, the
	score examples shown here are captured with a repeat count of zero.</em>`

	p15 := `<strong>Tonic Intervals</strong> presents 13 different intervals,
	i.e., all possible pitches relative to the chosen tonic pitch. Use this
	pattern as a self-test to gauge your progress at distinguishing the
	between the intervals.`

	p16 := `<strong>Two Intervals</strong> is, as you might expect, a series
	of three pitches specified by the interval1 and interval2 selectors. The
	score example shows a typical etude produced by choosing 4 half steps (a
	major third) for the lower interval and 3 half steps (a minor third) for
	the upper interval, i.e. a a major triad in root position. It's important
	to be able to recognize and play the notes of any interval pattern in any
	order. For 3 notes, there are 6 possible orderings. The program
	arranges for each ordering to occur twice among the 12 sequences
	presented.`

	p17 := `<strong>Three Intervals</strong> is similar to the Two Interval
	pattern but uses 3 intervals to produce 4-note sequence. There are 24
	examples in etudes produced with this pattern because you can play 4
	notes in 24 different orders. The example shows a typical etude
	constructed with a 2-2-1 pattern of half steps, corresponding to the
	first 4 notes of a major scale.`

	div = Div("",
		H3("", heading),
		P("", p0),
		P("", p0a),
		P("", p1),
		P("", p2),
		A(`name="patterns"`, H4("", "Patterns")),
		P(``, p3),
		H4("", "One Interval"),
		P("", p14),
		Img(`src="img/one_interval_excerpt.png" class="example"`),
		H4("", "Tonic Intervals"),
		P("", p15),
		Img(`src="img/allintervals_excerpt.png" class="example"`),
		H4("", "Two Intervals"),
		P("", p16),
		Img(`src="img/two_interval_excerpt.png" class="example"`),
		H4("", "Three Intervals"),
		P("", p17),
		Img(`src="img/three_interval_excerpt.png" class="example"`),
	)
	return
}
func forTheSerious() (div *HtmlTree) {
	p0 := `First the bad news: Playing by ear has a lot in common with
	learning to a speak a language fluently: you need to learn a lot of
	vocabulary (among other things). For a language the vocabulary elements
	are words, for music the elements are short sequences of pitches &mdash;
	and there are a lot of them.`

	p1 := `If you include unisons and octaves in the set of intervals to be
	used, there are 13*13 = 169 possible ways to put two intervals together
	to form a sequence of 3 pitches. For sequences of 4 pitches, there are
	13*13*13 = 2197 combinations of 3 intervals.`

	p2 := `Now consider that
	each sequence can be started on any pitch. If your instrument has a 3
	octave range, you're looking at (approximately) 6000 ways to play 3 notes
	and 79,000 ways to play 4. And, remember, each 3 note sequence can be
	played in 6 different orders and each 4 note sequence in 24 orders.`

	p3 := `So here's the good news: while it's probably true that if you
	searched all the music ever written or recorded you'd find at least one
	instance of each of those combinations, it's also certainly true that
	most music uses a much smaller subset most of the time. In fact, getting
	fluent on the 12 intervals plus 18 interval pairs (3 notes) and 35 sets
	of three intervals (4 notes) will take you a very long way. It's still a
	lot of work but you can get there if you're willing to devote five to ten
	minutes of your practice time every day.`

	p4 := `Your first goal should be to master recognising and playing single
	intervals at a brisk tempo over the full range of your instrument. The
	One Interval and Tonic Intervals patterns are your friends. Use the
	former to get solid on each of the 12 intervals and the latter to relate
	all the intervals to specific pitches.`

	p5 := `The next step is patterns of three notes created with the Two
	Intervals pattern. Here's a list of 18 common patterns that cover scale
	fragments and simple chords. <strong>Important:</strong> The patterns in
	the list below are in half-steps, e.g. "4-3" rather than "Major Third,
	Minor Third" to indicate a root position major triad.`

	ul1 := Ul(``,
		Li(``, "Scalar Diatonic and Pentatonic: 1-2, 2-1, 2-2, 2-3, 3-2"),
		Li(``, "Scalar Chromatic: 1-1, 1-3, 3-1"),
		Li(``, `Root Position Triads: 4-3, 3-4, 3-3, 4-4`),
		Li(``, `First Inversion Triads: 3-5, 4-5, 3-6`),
		Li(``, `Second Inversion Triads: 5-4, 5-3, 6-3`),
	)
	p6 := `Once the three note patterns start to feel easy, start working on
	the following four-note combinations with the Three Intervals pattern.
	These 35 patterns cover all the common scale fragments and chords,
	including dominant, minor, diminished and augmented 7ths.`

	ul2 := Ul(``,
		Li(``, `Scalar Diatonic and Pentatonic: 1-2-2, 2-1-2, 2-2-1, 2-2-2, 2-3-2, 3-2-2, 2-2-3`),
		Li(``, `Scalar Chromatic: 1-1-1, 1-2-1, 1-3-1, 2-1-3, 3-1-2`),
		Li(``, `Root Position Triads &amp; 7ths: 4-3-5, 4-3-4, 4-3-3, 3-4-5, 3-4-4, 3-4-3, 3-3-3, 3-3-4, 4-4-2`),
		Li(``, `First Inversion 7ths : 3-4-1, 3-3-2, 4-4-1, 3-4-2, 4-2-2`),
		Li(``, `Second Inversion 7ths: 4-1-4, 3-2-4, 4-1-3, 4-2-3, 2-2-4`),
		Li(``, `Third Inversion 7ths: 1-4-3, 2-4-3, 1-3-4, 2-3-3, 2-4-4`),
	)

	div = Div("",
		H3("", "For the serious"),
		P("", p0),
		P("", p1),
		P("", p2),
		P("", p3),
		P("", p4),
		H4(``, `Two Intervals (3 notes)`),
		P("", p5),
		ul1,
		H4(``, `Three Intervals (4 notes)`),
		P("", p6),
		ul2,
	)
	return
}
func forBeginners() (div *HtmlTree) {
	p1a := `The only real prerequisite is being able hear when a note you play
	on your instrument matches the one being played by Infinite Etudes.`

	p1 := `If you're just starting out with your instrument, please make
	sure you've had at least some basic instruction in how to hold your
	instrument comfortably with good posture and hand position and how to
	play individual notes cleanly over the full range of your instrument.`

	p2 := `If your instrument is tunable, make sure it's correctly tuned to
	standard pitch (A4=440). You'll also want to adjust the play volume to a level
	that's about the same volume as your instrument.  If you play Bass or any other
	low-pitched instrument make sure to use good speakers or a headset for the output
	from your computer or mobile device.`

	p3 := `If you're not sure about your readiness, test yourself with the
	with the simplest pattern and interval (One Interval, Unison) at the
	slowest tempo (60 BPM) with 3 repeats. Each measure will contain a single
	pitch and the entire measure will repeat 3 times after the first hearing.
	At 60 BPM, you'll have 16 seconds to locate and play the right pitch
	before the etude moves on to a new pitch. If you're not finding the pitch
	at least half the time, you may want to wait until you've become more
	familiar with your instrument.`

	p4 := `Using Infinite Etudes for five minutes every day will serve you
	better than an hour once a week. Your brain and neuromuscular system
	consolidate learning during sleep and, truthfully, the amount of new
	information we can absorb each day is limited. So be patient, please.
	Start with the One Interval pattern working from the smallest intervals
	(Minor2, Major2) up to the largest (Octave). Test yourself regularly with
	the Tonic Intervals pattern and a random Tonal Center. Try to wait until
	you're getting most of the intervals right on the first repeat before moving
	to the Two Intervals pattern.`
	div = Div("",
		A(`name="beginners"`, H3("", "For Beginners")),
		P("", p1),
		P("", p1a),
		P("", p2),
		P("", p3),
		P("", p4),
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

func userInterface() (div *HtmlTree) {
	p1 := `The Pattern selector allows you to choose one of the patterns
	described above. Your choice affects the visibility of the interval and
	tonal center selectors. The interval selectors, (Interval1, Interval2 and
	Interval3), appear according to the number intervals in the chosen
	pattern. The interval choices are labeled by the number of semitones
	(half steps) and the corresponding musical name, e.g. "4 (Minor Third)".
	The Tonal Center selector appears only when the Tonic Intervals pattern
	is selected.`

	p2 := `The Instrument selector provides a choice of common instrument sounds. Your choice also
	determines the range of pitches that can occur within an etude.`

	p2a := `Each etude starts on a randomly selected pitch somewhere between
	the lowest and highest notes that can commonly be played on your chosen
	instrument. The sequence orderings are a random walk constructed to that
	the first pitch of each sequence is "close" to the preceding pitch
	without wandering outside the playable range of your instrument.`

	p3 := `By default the metronome gives an initial 1 measure count-in and
	continues to click on each beat of the etude.  You can control this with the
	Metronome selector. Choose "downbeat" to have it click only on beat 1 of each measure.
	Choose "off" for silence after the count-in.`

	p4 := `Infinite Etudes generates MIDI files in 4/4 time with the tempo
	defaulted to 120 beats per minute. If you need it slower or faster, use
	the Tempo selector to choose a value between 60 and 480 beats per
	minute.`

	p5 := `Use the Repeats selector to change the number of repeats for each sequence. The default is
	3. You can set it to 2 or 1 to increase the challenge. You can also set it to 0, but that's
	not useful unless you want to download an example to import into a score editor.`

	p6 := `The Muting selector allows you silence one or more of the repeated
	measures. The cross mark symbol, &#x2717;, indicates a silent measure and
	the check mark, &#x2713;, indicates an audible one.`

	p7 := `The Play button tells the server to generate and start playing a
	new etude using the settings you've chosen in the the selectors. The Stop
	button stops the playback before the end of the etude. The Download
	button generates a new etude and allows you to save it as a MIDI file.`

	div = Div("",
		A(`name="ui"`, H3("", "User Interface")),
		H4("", "Pattern"),
		P("", p1),
		H4("", "Instrument"),
		P("", p2),
		P("", p2a),
		H4("", "Metronome"),
		P("", p3),
		H4("", "Tempo"),
		P("", p4),
		H4("", "Repeats"),
		P("", p5),
		H4("", "Muting"),
		P("", p6),
		H4("", "Play, Stop, Download"),
		P("", p7),
	)
	return
}
func custom() (div *HtmlTree) {
	p1 := `If you need something beyond the available tempi and instrument sounds, the easiest solution
	is to use the download button to save a local copy of a file and play it with
	a program that allows you finer control of the playback.  I recommend QMidi for Mac. I don't
	know what's good on PC but a little Googling should turn up something appropriate. Downloading
	also allows you to play the files through better equipment for more realistic sound.`

	p2 := `You might also consider installing MuseScore, the excellent open
	source notation editor. Version 3.1 and higher does a very good job
	importing Infinite Etudes midi files. Besides controlling tempo, you can
	print the etude as sheet music or play it back with real-time
	highlighting of each note as it's played.`

	p3 := `A third option, if you have software skills, is to install
	Infinite Etudes on your computer from the source code on <a
	href="https://github.com/Michael-F-Ellis/infinite-etudes">GitHub.</a> and
	adapt the program to your needs.`

	div = Div("",
		H3("", "Customizing"),
		P("", p1),
		P("", p2),
		P("", p3),
	)
	return
}
func variations() (div *HtmlTree) {
	p1 := `As you progress, some patterns will become easy to recognize and
	play before others. When you're nailing a particular pattern correctly and
	confidently on first hearing (hooray!), you can either increase the tempo, decrease the
	number of repeats, or leave the repeats at 3 and put the remaining two bars
	to good use in a variety of ways. Here are a few suggestions, some simple
	and some difficult:`
	var variants = []string{
		`Finger it differently.`,
		`Change the bowing or picking.`,
		`Play it with the other hand (keyboards).`,
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

/*
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
*/
func coda() (div *HtmlTree) {
	p1 := `I wrote Infinite Etudes for two reasons: First, as a tool for my
	own practice on piano, guitar and viola; second as a small project to develop a
	complete application in the Go programming language. I'm happy with it on
	both counts and I hope you find it useful also. The source code is available
	on <a href="https://github.com/Michael-F-Ellis/infinite-etudes">GitHub.</a>
	<br><br>Mike Ellis<br>Burlington, VT<br>Nov. 2020`
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
	  background-color: #000;
	  color: #CFC;
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
        margin-bottom: 5%;
    }
    img.example {
        margin-left: 5%;
        margin-right: 10%;
        width: 85vw;
	}
	label {
		display: inline-block;
		text-align: center;
		font-size: 80%;
	}
    select {
	  display: inline-block;
	  font-size: 125%;
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
	/* */
	.Row {
    display: table;
    width: auto;
    table-layout: auto;
    border-spacing: 10px;
    }
    .Column{
    display: table-cell;
    /* background-color: red; */
    }
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
		  var scaleselect = document.getElementById("scale-select")
		  scaleselect.addEventListener("change", manageInputs)
		  manageInputs()
		}
		// returns true if the selected key is an interval name
		function isIntervalName(name) {
			var inames = ['minor2', 'major2', 'minor3', 'major3', 'perfect4', 'tritone',
			'perfect5', 'minor6', 'major6', 'minor7', 'major7', 'octave']
			return inames.includes(name)
		}
		// manageInputs adjusts the enable status of the key and interval widgets
		// when scale-select value changes
		function manageInputs() {
			var key = document.getElementById("key-div")
			var interval1 = document.getElementById("interval1-div")
			var interval2 = document.getElementById("interval2-div")
			var interval3 = document.getElementById("interval3-div")
			var scalePattern = document.getElementById("scale-select").value
			if (scalePattern == "interval") {
				interval1.style.display=""
				interval2.style.display="none"
				interval3.style.display="none"
				key.style.display="none"
				return
			}
			if (scalePattern == "intervalpair") {
				interval1.style.display=""
				interval2.style.display=""
				interval3.style.display="none"
				key.style.display="none"
				return
			}
			if (scalePattern == "intervaltriple") {
				interval1.style.display=""
				interval2.style.display=""
				interval3.style.display=""
				key.style.display="none"
				return
			}
			// all the other patterns are chosen by key
			interval1.style.display="none"
			interval2.style.display="none"
			interval3.style.display="none"
			key.style.display=""
			return
		}
		// Read the selects and return the URL for the etude to be played or downloaded.
		function etudeURL() {
		  scale = document.getElementById("scale-select").value
		  key = document.getElementById("key-select").value
		  if (scale != "intervals" &&  isIntervalName(key)) {
			  alert(key + " is only valid when the scale pattern is Intervals.")
			  return ""
		  }
		  if (key=="random") {
			  key=randomKey()
			  };
		  interval1 = document.getElementById("interval1-select").value
		  interval2 = document.getElementById("interval2-select").value
		  interval3 = document.getElementById("interval3-select").value
		  sound = document.getElementById("sound-select").value
		  metronome = document.getElementById("metro-select").value
		  tempo = document.getElementById("tempo-select").value
		  repeats = document.getElementById("repeat-select").value
		  silent = document.getElementById("silence-select").value
		  return "/etude/" + key + "/" + scale + "/" + interval1 + "/" + interval2 + "/" + interval3 + "/" + sound + "/" + metronome + "/" + tempo + "/" + repeats + "/" + silent
		}

		// Read the selects and returns a proposed filename for the etude to be downloaded.
		function etudeFileName() {
		  key = document.getElementById("key-select").value
		  if (key=="random") {
			  key=randomKey()
			  };
		  scale = document.getElementById("scale-select").value
		  interval1 = document.getElementById("interval1-select").value
		  interval2 = document.getElementById("interval2-select").value
		  interval3 = document.getElementById("interval3-select").value
		  sound = document.getElementById("sound-select").value
		  metronome = document.getElementById("metro-select").value
		  tempo = document.getElementById("tempo-select").value
		  repeats = document.getElementById("repeat-select").value
		  silent = document.getElementById("silence-select").value
		  if (scale=="interval"){
			  return scale + "_" + interval1 + "_" + sound + "_" + metronome + "_" + tempo + "_" + repeats  + "_"+ silent + ".midi" 
		  }
		  if (scale=="intervalpair"){
			  return scale + "_" + interval1 + "_" + interval2 + "_" + sound + "_" + metronome + "_" + tempo + "_" + repeats  + "_" + silent + ".midi" 
		  }
		  if (scale=="intervaltriple"){
			  return scale + "_" + interval1 + "_" + interval2 + "_"  + interval3 + "_" + sound + "_" + metronome + "_" + tempo + "_" + repeats  + "_" + silent + ".midi" 
		  }
		  // any other scale 
		  return key + "_" + scale + "_" + sound + "_" + metronome + "_" + tempo + "_" + repeats  + "_" + silent + ".midi"
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
			var url = etudeURL()
			if (url != "") {
			  MIDIjs.play(url)
			}
		}

		function playStop() {
		    MIDIjs.stop()
		}
        
		function downloadEtude() {
          var url = etudeURL()
		  if (url == "") {
			  return // bad selection
		  }
		  // adapted from https://stackoverflow.com/a/49917066/426853
		  let a = document.createElement('a')
		  a.href = url
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
