// segmentParser.js
// This module contains the functions used to parse a segment of input text and
// convert it to a MIDI file.

// The input text is a string of simplified TBON notation that represents rhythm
// in a single voice.

// For example "f m m m | f m m m | f m m m | f m m m |" represents a segment of
// 4 bars of 4/4 time 4 quarter notes in each bar, with the first beat of each
// bar accented (forte) and the other beats mezzo-forte.
// 
// Bars are separated by the "|" character.  Beats are separated by spaces.

// Sub-beats are represented by one of the following characters: "fmpr-".  The
// characters 'f', 'm', and 'p' represent click volume, forte, mezzo-forte, and
// piano respectively.  The character 'r' represents a rest, and '-' represents
// a hold (tie), i.e sustaining the previous note.
//
// A few beat examples (assuming quarter note beat length):
// pr: soft eighth note followed by an eighth note rest
// fmmmm: a quintuplet, first note forte, other mezzo.
// f--m: dotted eighth note followed by a sixteenth note.
// fff: a triplet played with each sub-beat played forte
// frf: a triplet with a rest in the middle
// p--pp-: a triplet with a dotted eighth, a sixteenth, and an eighth note, all piano
// Notes can hold into beats, e.g.
// f -f: a dotted quarter and eighth, both forte, spanning two beats.
// fm-- -f: a forte sixteenth note, a mezzo-forte quarter note, and a forte eighth note.
//
// The overriding principle is that each beat is divided into equal sub-beats
// according the number of characters in the beat. Almost any imaginable rhythm
// can be represented in this notation.

// Beat durations are quarter notes by default, but this can be changed by
// placing an integer at the beginning of the segment, followed by space. The
// integer represents the the number of 16th notes in a beat.  For example, a
// measure 6/8 time could be represented as "6 fmm m |" to indicate that the
// beat duration is a dotted quarter note (6 16th notes) and the first beat is 3
// eighth notes long.  Note that same sound would be produced by "4 fm mm - |",
// a measure of 3/4 time.
//
// Changing the beat duration changes all subsequent measures until the next
// change.

/***************************************************************************
CLASS DEFINITIONS
***************************************************************************/
// The MIDIParameters class is used to store the parameters needed to convert an
// array of subBeat objects to a MIDI file
class MIDIParameters {
	constructor(qpm = 90, channel = 10, fNote = 76,
		mNote = 77, pNote = 77, rNote = 36,
		fVelocity = 110, mVelocity = 80,
		pVelocity = 40, rVelocity = 1) {

		this.qpm = qpm
		this.channel = channel
		this.fNote = fNote
		this.mNote = mNote
		this.pNote = pNote
		this.rNote = rNote
		this.fVelocity = fVelocity
		this.mVelocity = mVelocity
		this.pVelocity = pVelocity
		this.rVelocity = rVelocity
	}
	log() {
		console.log("MIDI Parameters:")
		console.log(" . qpm: " + this.qpm)
		console.log(" . channel: " + this.channel)
		console.log(" . fNote: " + this.fNote)
		console.log(" . mNote: " + this.mNote)
		console.log(" . pNote: " + this.pNote)
		console.log(" . rNote: " + this.rNote)
		console.log(" . fVelocity: " + this.fVelocity)
		console.log(" . mVelocity: " + this.mVelocity)
		console.log(" . pVelocity: " + this.pVelocity)
		console.log(" . rVelocity: " + this.rVelocity)

	}
	// args returns the arguments to the constructor as an array
	// suitable for serialization.
	args() {
		return [this.qpm, this.channel, this.fNote,
		this.mNote, this.pNote, this.rNote, this.fVelocity,
		this.mVelocity, this.pVelocity,
		this.rVelocity]
	}
}
// The Beat class is used to store the properties of a beat.
class Beat {
	constructor(beatText, barNumber, beatNumber, beatLength) {
		this.beatText = beatText
		this.barNumber = barNumber
		this.beatNumber = beatNumber
		this.beatLength = beatLength // the number of 16th notes in the beat
	}
	log() {
		console.log("Beat:")
		console.log(" . beatText: " + this.beatText)
		console.log(" . barNumber: " + this.barNumber)
		console.log(" . beatNumber: " + this.beatNumber)
		console.log(" . beatLength: " + this.beatLength)
	}
}
// The SubBeat class is used to store the properties of a sub-beat.
class SubBeat {
	constructor(subBeatText = "", barNumber = 0, beatNumber = 0, beatLength = 0,
		subBeatNumber = 0, isDownbeat = false, isAttack = false, isRest = false, isHold = false) {
		this.subBeatText = subBeatText
		this.barNumber = barNumber
		this.beatNumber = beatNumber
		this.beatLength = beatLength // the number of 16th notes in the beat
		this.subBeatNumber = subBeatNumber
		this.isDownbeat = isDownbeat
		this.isAttack = isAttack
		this.isRest = isRest
		this.isHold = isHold
		// the following properties are computed later
		this.Length = 0; // subBeat length in quarter notes
		this.Note = 0; // MIDI note number
		this.Velocity = 0; // MIDI velocity
		this.Start = 0; // start time in quarter notes
	}
	log() {
		console.log("SubBeat:")
		console.log(" . subBeatText: " + this.subBeatText)
		console.log(" . barNumber: " + this.barNumber)
		console.log(" . beatNumber: " + this.beatNumber)
		console.log(" . beatLength: " + this.beatLength)
		console.log(" . subBeatNumber: " + this.subBeatNumber)
		console.log(" . isDownbeat: " + this.isDownbeat)
		console.log(" . isAttack: " + this.isAttack)
		console.log(" . isRest: " + this.isRest)
		console.log(" . isHold: " + this.isHold)
		console.log(" . Length: " + this.Length)
		console.log(" . Note: " + this.Note)
		console.log(" . Velocity: " + this.Velocity)
		console.log(" . Start: " + this.Start)
	}
}
// The ParsingStages class is used to store the intermediate results of parsing
// the input text.  It has methods to parse the input text into bars, beats, and
// sub-beats, and to compute the length and start time of each sub-beat.  It also
// has a method to create a MIDI file from the sub-beats.
class ParsingStages {
	constructor(pattern, midiparms, nrepeats, preroll, description = "describe me") {
		this.pattern = pattern // 
		this.preroll = preroll // one or more bars to play before the segment
		this.nrepeats = nrepeats // the number of times to repeat the segment
		this.midiparms = midiparms // an object of class MIDIParameters
		this.description = description // a description of the pattern
		this.segment = ""	// the preroll and repeated pattern
		this.bars = { barTexts: [], beatlengths: [] }
		this.beats = [] // an array of beat objects
		this.subBeats = [] // an array of subBeat objects
		this.qlength = 0 // the total duration of the track in quarter notes
		this.qpreroll = 0 // the duration of the preroll in quarter notes
		this.ErrorMessages = []
	}
	/***************************************************************
	METHODS
	***************************************************************/
	// log() is a utility function that logs the properties of the ParsingStages
	// object to the console, invoking the log() method of objects that compose
	// the ParsingStages object.
	log() {
		console.log("Parsing Stages:")
		console.log(" . pattern: " + this.segment)
		console.log(" . description: " + this.description)
		console.log(" . preroll: " + this.preroll)
		console.log(" . nrepeats: " + this.nrepeats)
		console.log(" . bars: " + this.bars.barTexts)
		console.log(" . beat lengths: " + this.bars.beatlengths)
		// log each object in the beats array
		for (let i = 0; i < this.beats.length; i++) {
			this.beats[i].log()
		}
		// log each object in the subBeats array
		for (let i = 0; i < this.subBeats.length; i++) {
			this.subBeats[i].log()
		}
		// log the MIDI parameters
		this.midiparms.log()

		// log the total duration of the track
		console.log("Track duration (in quarter-notes): " + this.qlength)

		// log any error messages
		if (this.ErrorMessages.length > 0) {
			console.log("Error Messages:")
			for (let i = 0; i < this.ErrorMessages.length; i++) {
				console.log(" . " + this.ErrorMessages[i])
			}
		} else {
			console.log("No error messages")
		}
	}
	// args returns the arguments to the constructor as an array
	// suitable for serialization.
	args() {
		return [this.pattern, this.midiparms.args(), this.nrepeats, this.preroll, this.description]
	}

	// appendErrorMessage() is a utility function that appends an error message to
	// the ErrorMessages array.
	appendErrorMessage(msg) {
		this.ErrorMessages.push(msg)
	}
	// clearErrorMessages() is a utility function that clears the ErrorMessages array.
	clearErrorMessages() {
		this.ErrorMessages = []
	}
	// parse() is the main function that parses the input text and computes the
	// properties of the sub-beats, and creates a midi track and returns a playable URI.
	// If an error is encountered, the function appends an error message to the
	// ErrorMessages array and returns null.  
	parse() {
		// Trim any leading or trailing whitespace from the input text.
		this.segment = this.pattern.trim()
		// Catenate the input text nrepeats times. Before doing this, we need to
		// ensure that the input text ends with a "|".  If it doesn't, we append
		// one.  Similarly, we need to ensure that the input text does not start
		// with a "|".  If it does, we remove it.
		if (this.segment.endsWith("|") == false) {
			this.segment += " |"
		}
		if (this.segment.startsWith("|")) {
			this.segment = this.segment.slice(1)
		}
		// catenate the input text nrepeats times. If repeats is <1, we set it
		// to 1.
		this.nrepeats = (this.nrepeats < 1) ? 1 : this.nrepeats
		this.segment = this.segment.repeat(this.nrepeats)

		// If the preroll is not empty, we prepend it to the input text after
		// ensuring that the preroll ends with a "|".  If it doesn't, we append
		// one.
		this.preroll = this.preroll.trim()
		if (this.preroll != "") {
			if (this.preroll.endsWith("|") == false) {
				this.preroll += " |"
			}
			this.segment = this.preroll + this.segment
		}

		// The first step is to parse the input text into an array of bars
		this.parseBars()
		if (this.ErrorMessages.length > 0) {
			return null
		}
		this.parseBeats()
		if (this.ErrorMessages.length > 0) {
			return
		}
		this.parseSubBeats()
		if (this.ErrorMessages.length > 0) {
			return null
		}
		this.computeSubBeatLengths()
		if (this.ErrorMessages.length > 0) {
			return null
		}
		this.computeSubBeatStarts()
		if (this.ErrorMessages.length > 0) {
			return null
		}
		this.setSubBeatNotes()
		if (this.ErrorMessages.length > 0) {
			return null
		}
		this.computeTotalDuration()
		if (this.ErrorMessages.length > 0) {
			return null
		}
		// Generate the MIDI file as a data URI
		let uri = this.createMidi()
		if (this.ErrorMessages.length > 0) {
			return null
		} else {
			return uri
		}
	}

	// parseBars splits the input text into an array of bars
	parseBars() {
		let rawbars = this.segment.trim().split("|")
		// remove any empty bars
		for (let i = 0; i < rawbars.length; i++) {
			if (rawbars[i].trim() == "") {
				rawbars.splice(i, 1)
			}
		}
		if (rawbars.length == 0) {
			this.appendErrorMessage("No bars found in input text: " + this.segment)
		}
		this.bars.barTexts = rawbars
	}
	// parseBeats splits the array of bars into an array of beats. 
	parseBeats() {
		for (let i = 0; i < this.bars.barTexts.length; i++) {
			let bar = this.bars.barTexts[i].trim()
			let beatArray = bar.split(/\s/) // split on whitespace
			// If the first element of the beatArray is an integer >= 1, it
			// represents the number of 16th notes in each subsequent beat and we append it
			// to the beatlengths array in the parsingStages object. Otherwise, we append the
			// prior beatlength to the beatlengths array. The default beatlength is 4.
			let priorLength = 4
			if (beatArray[0].match(/^\d+$/)) {
				priorLength = parseInt(beatArray[0])
				beatArray.shift() // remove the first element from the beatArray
			}
			this.bars.beatlengths.push(priorLength)

			// append the beat to the beats array as a Beat object
			for (let j = 0; j < beatArray.length; j++) {
				this.beats.push(new Beat(
					beatArray[j], 				// beat text
					i + 1,        				// bar number
					j + 1,       				// beat number
					this.bars.beatlengths[i] // beat length
				),
				)
			}
		}
	}
	// parseSubBeats splits the array of beats into an array of sub-beat objects.
	// Sub-beats are numbered starting from 1. 
	parseSubBeats() {
		for (let i = 0; i < this.beats.length; i++) {
			let beat = this.beats[i].beatText.trim()
			let subBeatArray = beat.split("") // split the beat into an array of characters	
			let nGood = 0 // count the number of good sub-beats	
			for (let j = 0; j < subBeatArray.length; j++) {
				// valid characters are f, m, p, r and -
				let c = subBeatArray[j]
				if (['f', 'm', 'p', 'r', '-'].includes(c)) {
					let sb = new SubBeat(c,
						this.beats[i].barNumber,
						this.beats[i].beatNumber,
						this.beats[i].beatLength, // sixteenth notes
						nGood + 1, // sub-beat number
						false, // isDownbeat (we'll set this later)
						(c == 'f' || c == 'm' || c == 'p') ? true : false, // isAttack
						(c == 'r') ? true : false, // isRest
						(c == '-') ? true : false // isHold
					)
					// set the isDownbeat property to true if the sub-beat is the
					// first sub-beat in the first beat of the bar.
					if (this.beats[i].beatNumber == 1 && nGood == 0) {
						sb.isDownbeat = true
					}
					this.subBeats.push(sb)
					nGood++
				} else {
					this.appendErrorMessage("Invalid character '" + c + "' in beat " + beat)
				}
			}
		}
	}
	// We need to compute the length of each sub-beat in quarter
	// notes.  In musical notation, the length of a beat is determined by the denominator of the
	// time signature.  For example, in 4/4 time, the length of a beat is 1,
	// in 2/2 time the length is 2, and in 6/8 time the length is 3/2. The
	// existence of compound time signatures like 6/8 this
	// function more complex than it would be otherwise.
	//
	// To simplify the problem, we represent beat length as the number of
	// 16th notes in a beat.  Thus, in 4/4 time, the beat length is 4, in
	// 2/2 time the beat length is 8, and in 6/8 time the beat length is 6.
	//
	// The length of a sub-beat is determined by the length of the beat and
	// the number of sub-beats in the beat. We can compute this as:
	// 		qlength = (beatLength/nSubbeats)/4 
	// Diving by 4 converts the length from sixteenths to quarters.
	//
	// The next step is to compute the length of each sub-beat in quarter notes.
	// We do this by stepping through the array of sub-beats in reverse order.
	// The sub-beat number of the last sub-beat in in each beat is the
	// number of sub-beats in the beat.  We will use this number to compute the
	// length of all the sub-beats in the beat.
	computeSubBeatLengths() {
		let nSubBeats = 0
		let beatLen = 0
		let newBeat = true // we set this to false while processing a beat, and to true when we start a new beat 

		for (let i = this.subBeats.length - 1; i >= 0; i--) {
			if (newBeat) {
				beatLen = this.subBeats[i].beatLength
				nSubBeats = this.subBeats[i].subBeatNumber
				newBeat = false
			}
			// compute and append the length of the sub-beat
			var qlength = (beatLen / nSubBeats) / 4
			this.subBeats[i].Length = qlength

			// if this is the first sub-beat in the beat, set newBeat to true
			if (this.subBeats[i].subBeatNumber == 1) { // 3rd element of the tuple is the sub-beat number
				newBeat = true
			}
		}
	}
	// The next step is computing the start time of each sub-beat in quarter note
	// lengths.  We do this by stepping through the array of sub-beats in forward
	// order, summing the lengths of the sub-beats. 
	// (Note: the durations are all we need to create a midi file, but the start
	// times are needed for marking events on the clock circle.)
	computeSubBeatStarts() {
		let start = 0
		for (let i = 0; i < this.subBeats.length; i++) {
			this.subBeats[i].Start = start
			start += this.subBeats[i].Length
		}
	}
	// Next we use the MIDI parameters to set the note and velocity of each
	// sub-beat.  Attack sub-beats are set to the note and velocity corresponding
	// the sub-beat text.  Rest and Hold sub-beats are allowed to remain at 0.
	setSubBeatNotes() {
		for (let i = 0; i < this.subBeats.length; i++) {
			let c = this.subBeats[i].subBeatText
			if (c == 'f') {
				this.subBeats[i].Note = this.midiparms.fNote
				this.subBeats[i].Velocity = this.midiparms.fVelocity
			} else if (c == 'm') {
				this.subBeats[i].Note = this.midiparms.mNote
				this.subBeats[i].Velocity = this.midiparms.mVelocity
			} else if (c == 'p') {
				this.subBeats[i].Note = this.midiparms.pNote
				this.subBeats[i].Velocity = this.midiparms.pVelocity
			}
		}
	}
	// It's useful to have a function to compute the total duration of a track.
	// This is the sum of the lengths of all the sub-beats. The returned value
	// is the duration in quarter notes. This function also computes the quarter
	// note length of the preroll and assigns it to the qpreroll property of the
	// parsingStages object.
	computeTotalDuration() {
		let total = 0
		for (let i = 0; i < this.subBeats.length; i++) {
			total += this.subBeats[i].Length
		}
		this.qlength = total
		// compute the length of the preroll
		let pr = this.preroll.trim().split("|")
		// remove any empty bars
		for (let i = 0; i < pr.length; i++) {
			if (pr[i].trim() == "") {
				pr.splice(i, 1)
			}
		}
		let prLength = 0
		for (let i = 0; i < pr.length; i++) {
			let beatArray = pr[i].trim().split(/\s/)
			let beatLength = 4 // start with default beat length = quarter note
			// If the first element of the beatArray is an integer >= 1, it
			// represents the number of 16th notes in each subsequent beat.`
			if (beatArray[0].match(/^\d+$/)) {
				// change the beat length
				beatLength = parseInt(beatArray[0])
				beatArray.shift() // remove the first element from the beatArray
			}
			// increment the prLength by the length of the beat
			for (let j = 0; j < beatArray.length; j++) {
				prLength += beatLength
			}
		}
		this.qpreroll = prLength / 4 // convert to quarter notes
	}
	// createMidi generates a midi file from the sub-beats. We step through the
	// sub-beats. At each Attack event, we step through any holds that follow
	// it, summing durations, until we reach a rest or a new attack. We create a
	// new MidiWriter.NoteEvent using the summed durations and add it to the
	// track. For rests we follow a similar process, except that we add a
	// MidiWriter.WaitEvent to the track.
	//
	// Note: MidiWriter uses ticks for durations and waits. One beat is 128
	// ticks.  so we multiply the duration by 128 to get the number of ticks,
	// round it to the nearest integer, and convert it to a string with a 'T'
	// prefix as required by MidiWriter.	
	createMidi() {
		var track = new MidiWriter.Track();
		// Set the tempo specified in the MIDI parameters
		track.setTempo(this.midiparms.qpm)
		// Step through the sub-beats
		for (let i = 0; i < this.subBeats.length; i++) {
			let subBeat = this.subBeats[i]
			if (subBeat.isHold == false) {
				let j = i + 1
				let duration = subBeat.Length
				// sum the durations of any holds that follow the attack
				while (j < this.subBeats.length && this.subBeats[j].isHold == true) {
					duration += this.subBeats[j].Length
					j++
				}
				// get the pitch and velocity
				if (subBeat.subBeatText == "f") {
					subBeat.Note = this.midiparms.fNote
					subBeat.Velocity = this.midiparms.fVelocity
				} else if (subBeat.subBeatText == "m") {
					subBeat.Note = this.midiparms.mNote
					subBeat.Velocity = this.midiparms.mVelocity
				} else if (subBeat.subBeatText == "p") {
					subBeat.Note = this.midiparms.pNote
					subBeat.Velocity = this.midiparms.pVelocity
				} else if (subBeat.subBeatText == "r") {
					subBeat.Note = this.midiparms.rNote
					subBeat.Velocity = this.midiparms.rVelocity // a velocity of 1 is a rest
				}
				// create a new note event and add it to the track
				let note = new MidiWriter.NoteEvent({
					pitch: subBeat.Note,
					velocity: subBeat.Velocity,
					duration: 'T' + Math.round(duration * 128),
					channel: this.midiparms.channel
				})

				track.addEvent(note)

				// set i to the last sub-beat that was part of the attack so the
				// next iteration of the outer loop will start at the next
				// attack or rest
				i = j - 1

			}
		}

		// Generate the MIDI file as a data URI
		var write = new MidiWriter.Writer([track]);
		return write.dataUri();

	}
}
// serializeParsingStages and deserializeParsingStages are utility functions
// to support storing and retrieving ParsingStages objects in persistent storage.
function serializeParsingStages(ps) {
	return JSON.stringify(ps.args())
}
function deserializeParsingStages(s) {
	let args = JSON.parse(s)
	midiparms = new MIDIParameters(...args[1])
	return new ParsingStages(args[0], midiparms, args[2], args[3], args[4])
}

// Class Library creates and manages an array of objects consisting of a name, a
// ParsingStages object, and a timestamp.  The Library class has methods to add,
// replace and remove objects and to store, sort and retrieve the library from
// local storage.
class Library {
	constructor(key = "grooveClockLibrary") {
		this.key = key
		this.library = []
		this.load()
	}
	// has(name) returns true if the library contains an object with the given name
	has(name) {
		for (let i = 0; i < this.library.length; i++) {
			if (this.library[i].name == name) {
				return true
			}
		}
		return false
	}

	// get a ParsingStages object from the library by name
	getByName(name) {
		for (let i = 0; i < this.library.length; i++) {
			if (this.library[i].name == name) {
				return deserializeParsingStages(this.library[i].ps)
			}
		}
	}
	// add a new ParsingStages object to the library
	add(name, ps) {
		let timestamp = new Date().getTime()
		this.library.push({
			name: name,
			ps: serializeParsingStages(ps),
			timestamp: timestamp
		})
		this.save()
	}
	// remove an object from the library
	remove(name) {
		for (let i = 0; i < this.library.length; i++) {
			if (this.library[i].name == name) {
				this.library.splice(i, 1)
				this.save()
				break
			}
		}
	}
	// replace an object in the library
	replace(name, ps) {
		for (let i = 0; i < this.library.length; i++) {
			if (this.library[i].name == name) {
				this.library[i].ps = serializeParsingStages(ps)
				this.library[i].timestamp = new Date().getTime()
				this.save()
				break
			}
		}
	}
	// save the library to local storage
	save() {
		localStorage.setItem(this.key, JSON.stringify(this.library))
	}
	// load the library from local storage. If it doesn't exist, create a new
	// empty library and a default groove.
	load() {
		let lib = localStorage.getItem(this.key)
		if (lib != null) {
			this.library = JSON.parse(lib)
		}
		if (this.getByName("Default") == undefined) {
			// Add a detault groove to the library and save it.
			let ps = new ParsingStages(
				"f p p p | f p p p | f p p p | f p p p |",
				new MIDIParameters(), 8,
				'p p p p',
				"Default: 4 bars of 4/4 time with the first beat of each " +
				"bar accented (forte) and the other beats softer (piano). " +
				"Repeats 8 times with a 1 bar preroll.")
			this.add("Default", ps)
			this.save()
		}
	}
	// get the library, sorted by name or by timestamp
	getLibrary(sortbyName = true) {
		if (sortbyName) {
			this.library.sort((a, b) => (a.name > b.name) ? 1 : -1)
		} else {
			this.library.sort((a, b) => (a.timestamp < b.timestamp) ? 1 : -1)
		}
		return this.library
	}
	// fields(btnclass, nameclass, descclass, usecb) returns an html grid of the
	// library. Each row contains the following items: 
	//	1. An html button labeled "Use" with class btnclass and an 
	// .   onclick event that calls usecb with the name of the library 
	// .   item as an argument.
	//  2. An html button labeled "Del." with class btnclass and an onclick
	//     event that calls remove with the name of the library item as an
	//     argument after prompting the user to confirm the deletion.
	//  3. The library item name as a span with class nameclass
	//  4. The library item description as a span with class descclass
	fields(usecb, delcb) {
		let html = "<table id='library' style='border-collapse: separate; border-spacing: 10px;' > "
		for (let i = 0; i < this.library.length; i++) {
			let name = this.library[i].name
			let ps = deserializeParsingStages(this.library[i].ps)
			let desc = ps.description
			html += "<tr>"
			html += "<td><button onclick=" + usecb + '("' + name + '")>Use</button></td>'
			html += "<td><button onclick=" + delcb + '("' + name + '")>Del</button></td>'
			html += "<td>" + name + "</td>"
			html += "<td>" + desc + "</td>"
			html += "</tr>" // close the row
		}
		html += "</table>"
		return html
	}
}
// TwoWayMap is a class that implements a 2-way lookup table. It has methods to
// add, remove, and retrieve items by name or value, and to return the names
// sorted in alphabetic order and the values sorted in ascending order. The constructor
// takes an array of name-value pairs and uses them to initialize the the map.
class TwoWayMap {
	constructor(arrnv) {
		this.nameToValue = new Map();
		this.valueToName = new Map();
		// loop through arrnv adding each pair to each map
		arrnv.forEach(pair => {
			this.add(...pair)
		});
	}

	add(name, value) {
		this.nameToValue.set(name, value);
		this.valueToName.set(value, name);
	}

	removeByName(name) {
		const value = this.nameToValue.get(name);
		this.nameToValue.delete(name);
		this.valueToName.delete(value);
	}

	removeByValue(value) {
		const name = this.valueToName.get(value);
		this.nameToValue.delete(name);
		this.valueToName.delete(value);
	}

	getName(value) {
		return this.valueToName.get(value);
	}

	getValue(name) {
		return this.nameToValue.get(name);
	}
	// alphaSortedNames returns an array of names sorted in alphabetic
	// order
	alphaSortedNames() {
		return [...this.nameToValue.keys()].sort()
	}
	// sortedNumericValues returns an array of numeric values sorted in ascending
	// order
	sortedNumericValues() {
		return [...this.valueToName.keys()].sort((a, b) => a - b)
	}

}
// A two way map of all General MIDI percussion sound names and numbers
const gmPercussion = new TwoWayMap([
	["Acoustic Bass Drum", 35],
	["Bass Drum 1", 36],
	["Side Stick", 37],
	["Acoustic Snare", 38],
	["Hand Clap", 39],
	["Electric Snare", 40],
	["Low Floor Tom", 41],
	["Closed Hi Hat", 42],
	["High Floor Tom", 43],
	["Pedal Hi-Hat", 44],
	["Low Tom", 45],
	["Open Hi-Hat", 46],
	["Low-Mid Tom", 47],
	["Hi-Mid Tom", 48],
	["Crash Cymbal 1", 49],
	["High Tom", 50],
	["Ride Cymbal 1", 51],
	["Chinese Cymbal", 52],
	["Ride Bell", 53],
	["Tambourine", 54],
	["Splash Cymbal", 55],
	["Cowbell", 56],
	["Crash Cymbal 2", 57],
	["Vibraslap", 58],
	["Ride Cymbal 2", 59],
	["Hi Bongo", 60],
	["Low Bongo", 61],
	["Mute Hi Conga", 62],
	["Open Hi Conga", 63],
	["Low Conga", 64],
	["High Timbale", 65],
	["Low Timbale", 66],
	["High Agogo", 67],
	["Low Agogo", 68],
	["Cabasa", 69],
	["Maracas", 70],
	["Short Whistle", 71],
	["Long Whistle", 72],
	["Short Guiro", 73],
	["Long Guiro", 74],
	["Claves", 75],
	["Hi Wood Block", 76],
	["Low Wood Block", 77],
	["Mute Cuica", 78],
	["Open Cuica", 79],
	["Mute Triangle", 80],
	["Open Triangle", 81],
]
)
