package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Michael-F-Ellis/infinite-etudes/internal/miditempo"
)

// randString returns a random string of length n chosen from chars.
func randString(chars []rune, n uint) (out string) {
	var outslice []rune
	for i := 0; i < int(n); i++ {
		outslice = append(outslice, chars[rand.Intn(int(n))])
	}
	out = string(outslice)
	return
}

// serveEtudes serves etude midi files from the current working directory.
func serveEtudes(hostport string, maxAgeSeconds int, midijsPath string) {
	err := mkWebPages()
	if err != nil {
		log.Fatalf("could not write web pages: %v", err)
	}
	err = validDirPath(midijsPath)
	if err != nil {
		log.Fatalf("invalid midijs path: %v", err)
	}
	os.Setenv("MIDIJS", midijsPath)
	defer os.Unsetenv("MIDIJS")
	os.Setenv("ETUDE_MAX_AGE", fmt.Sprintf("%d", maxAgeSeconds))
	defer os.Unsetenv("ETUDE_MAX_AGE")
	http.Handle("/", http.HandlerFunc(indexHndlr))
	http.Handle("/etude/", http.HandlerFunc(etudeHndlr))
	http.Handle("/midijs/", http.HandlerFunc(midijsHndlr))
	log.Printf("midijs path is %s", os.Getenv("MIDIJS"))
	var serveSecure bool
	var certpath, certkeypath string
	if hostport == ":443" {
		certpath, certkeypath, err = getCertPaths()
		if err != nil {
			log.Printf("Can't find SSL certificates: %v", err)
			hostport = ":80"
		}
		serveSecure = true
	}
	log.Printf("serving on %s\n", hostport)
	switch serveSecure {
	case true:
		if err := http.ListenAndServeTLS(hostport, certpath, certkeypath, nil); err != nil {
			log.Fatalf("Could not listen on port %s : %v", hostport, err)
		}
	default:
		if err := http.ListenAndServe(hostport, nil); err != nil {
			log.Fatalf("Could not listen on port %s : %v", hostport, err)
		}
	}
}

// getCert attempts to retrieve a certficate and key for use with
// ListenAndServeTLS. It returns an error if either item cannot be found but
// does not otherwise attempt to validate them. That is left up to
// ListenAndServeTLS.
func getCertPaths() (certpath string, keypath string, err error) {
	certpath = os.Getenv("IETUDE_CERT_PATH")
	if certpath == "" {
		err = fmt.Errorf("no environment variable IETUDE_CERT_PATH")
		return
	}
	keypath = os.Getenv("IETUDE_CERTKEY_PATH")
	if keypath == "" {
		err = fmt.Errorf("no environment variable IETUDE_CERTKEY_PATH")
		return
	}
	return
}

// indexHndlr returns index.html
func indexHndlr(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// midijsHndlr returns files from the MIDIJS directory.
func midijsHndlr(w http.ResponseWriter, r *http.Request) {
	what := strings.Split(r.URL.Path, "/")
	if what[1] != "midijs" {
		log.Fatalf("programming error. got request path that didn't start with 'midijs': %s", r.URL.Path)
	}

	dir := os.Getenv("MIDIJS")
	pathelements := append([]string{dir}, what[2:]...)
	path := filepath.Join(pathelements...)
	http.ServeFile(w, r, path)

}

type etudeRequest struct {
	tonalCenter string
	pattern     string
	interval1   string
	interval2   string
	instrument  string
	rhythm      string
	tempo       string
}

// etudeHndlr returns a midi file that matches the get request or a 404 for
// incorrectly specified etudes. The pattern is
// /etude/<key>/<scale>/<instrument>/<advancing> where <key> is a pitchname like
// "c" or "aflat", <scale> is a scalename like "pentatonic", instrument is a
// formatted General Midi instrument name like "acoustic_grand_piano" and
// advancing is one of 'steady' or 'advancing' indicating the rhythm pattern to
// use. An optional final component, <tempo> allows specifying a tempo in beats
// per minute. Integer values between 20 and 600 are supported. If any of the
// foregoing pattern components are unknown or unsupported by this app,
// etudeHndlr gives a 400 response (StatusBadRequest). If the request matches a
// valid filename, the file will be returned in the response body if it exists
// and is younger than the maximum age imposed by this service. Otherwise the
// app will generate it so it can be returned.
func etudeHndlr(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	// Note first element of what is an empty string
	// log.Println(what)
	if path[1] != "etude" {
		log.Fatalf("programming error. got request path that didn't start with 'etude': %s", r.URL.Path)
	}
	var req etudeRequest
	req.tonalCenter = path[2]
	req.pattern = path[3]
	req.interval1 = path[4]
	req.interval2 = path[5]
	req.instrument = path[6]
	req.rhythm = path[7]
	req.tempo = path[8]
	if !validEtudeRequest(req) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filename := strings.Join(path[2:6], "_") + ".mid"
	log.Println(filename)
	advancing := false
	if path[5] == "advancing" {
		advancing = true
	}
	var tempo int
	if len(path) == 7 {
		log.Println("tempo requested")
		tmpo, err := strconv.Atoi(path[6])
		if err != nil {
			log.Printf("bad tempo requested, %s: %v", path[6], err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if tmpo < 20 || tmpo > 600 {
			log.Printf("tempo %d out of range. must be between 20 and 600", tmpo)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tempo = tmpo
	}
	makeEtudesIfNeeded(filename, path[4], advancing)
	if tempo != 0 {
		// log.Println("need to serve an altered file")
		µs := uint(60000000 / tempo) // microseconds per beat
		bytes, err := miditempo.SetTempo(filename, µs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v", err)
			return
		}
		unique := randString([]rune("abcedfghijklmnopqrstuvwxyz"), 4) // used to make deferred removal safe
		filename = fmt.Sprintf("new filename: %s_%d_%s.mid", filename[0:len(filename)-4], tempo, unique)
		// log.Println(filename)
		err = ioutil.WriteFile(filename, bytes, 0644)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v", err)
			return
		}
		// defer os.Remove(filename) // avoid proliferation
	}
	http.ServeFile(w, r, filename)
	// log the request in format that's convenient for analysis
	log.Printf("%s %s %s %s %s served\n", r.RemoteAddr, path[2], path[3], path[4], filename)
}

// makeEtudesIfNeeded generates a full set of etudes in the current
// working directory if the requested file doesn't exist or is older
// than the age limit set by serveEtudes in os.Environ. Otherwise
// it does nothing.
func makeEtudesIfNeeded(filename, instrumentName string, advancing bool) {
	var exists = true // initial assumption
	finfo, err := os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			// something's really wrong
			log.Fatalf("error statting %s: %v", filename, err)
		}
		exists = false
	}
	// if file exists, we need to check its age
	if exists {
		maxage, e := strconv.Atoi(os.Getenv("ETUDE_MAX_AGE"))
		if e != nil {
			log.Fatalf("programming error. can't convert ETUDE_MAX_AGE '%s' to integer", os.Getenv("ETUDE_MAX_AGE"))
		}

		maxduration := time.Duration(maxage) * time.Second
		modtime := finfo.ModTime()
		if time.Since(modtime) < maxduration {
			// nothing to do
			return
		}
	}
	// need to generate if we get to here
	iInfo, _ := getSupportedInstrumentByName(instrumentName) // already validated. ignore err value
	// fmt.Printf("%v %s\n", iInfo, filename)
	instrument := iInfo.gmnumber - 1
	midilo := iInfo.midilo
	midihi := iInfo.midihi
	tempo := 120
	mkAllEtudes(midilo, midihi, tempo, instrument, iInfo.name, advancing)
}

// validEtudeRequest returns true if the request is correctly formed
// and references a valid etude filename.
func validEtudeRequest(req etudeRequest) (ok bool) {
	if !validPattern(req.pattern) {
		return
	}

	switch req.pattern { // Intervals get special handling
	case "allintervals":
		if !validKeyName(req.tonalCenter) {
			return
		}
	case "interval":
		if !validIntervalName(req.interval1) {
			return
		}
	case "intervalpair":
		if !validIntervalName(req.interval1) || !validIntervalName(req.interval2) {
			return
		}

	default:
		if !validKeyName(req.tonalCenter) {
			return
		}
	}
	if !validInstrumentName(req.instrument) {
		return
	}
	if !validRhythmPattern(req.rhythm) {
		return
	}
	ok = true
	return
}

type nameInfo struct {
	fileName string
	uiName   string
	uiAria   string // alternate text for screen readers
	size     int    // interval size in half steps. Not meaningful for other parameters.
}

var keyInfo = []nameInfo{
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

var intervalInfo = []nameInfo{
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

// validIntervalName returns true if the interval name is in the ones we support.
func validIntervalName(name string) (ok bool) {
	for _, k := range intervalInfo {
		if k.fileName == name {
			ok = true
			break
		}
	}
	return
}

// extractIntervalPair returns two interval sizes from an interval pair string
// of the form "N-M" where N and M are interval sizes in half steps
func extractIntervalPair(s string) (i1 int, i2 int, err error) {
	intervals := strings.Split(s, "-")
	if len(intervals) != 2 {
		err = fmt.Errorf("Expected 2 intervals, got %d", len(intervals))
		return
	}
	i1, err = strconv.Atoi(intervals[0])
	if err != nil {
		err = fmt.Errorf("bad string for first interval: %v", err)
		return
	}
	if i1 < 1 || i1 > 12 {
		err = fmt.Errorf("bad value for first interval: %d", i1)
		return
	}
	i2, err = strconv.Atoi(intervals[1])
	if err != nil {
		err = fmt.Errorf("bad string for second interval: %v", err)
		return
	}
	if i2 < 1 || i2 > 12 {
		err = fmt.Errorf("bad value for second interval: %d", i2)
		return
	}
	return
}

// validKeyName returns true if the key name is in the ones we support.
func validKeyName(name string) (ok bool) {
	for _, k := range keyInfo {
		if k.fileName == name {
			ok = true
			break
		}
	}
	return
}

var scaleInfo = []nameInfo{
	{"allintervals", "All Intervals", "All Intervals", 0},
	{"interval", "Interval", "Interval", 0},
	{"intervalpair", "Interval Pairs", "Interval Pairs", 0},
	{"pentatonic", "Pentatonic", "Pentatonic", 0},
	{"final", "Chromatic Final", "Chromatic Final", 0},
	{"plus_four", "Plus Four", "Plus Four", 0},
	{"plus_seven", "Plus Seven", "Plus Seven", 0},
	{"four_and_seven", "Four and Seven", "Four and Seven", 0},
	{"raised_five", "Harmonic Minor 1", "Harmonic Minor 1", 0},
	{"raised_five_with_four_or_seven", "Harmonic Minor 2", "Harmonic Minor 2", 0},
}

// validPattern returns true if the scale name is in the ones we support.
func validPattern(name string) (ok bool) {
	for _, s := range scaleInfo {
		if s.fileName == name {
			ok = true
			break
		}
	}
	return
}

// validInstrumentName returns true if the instrument name is in the ones we
// support.
func validInstrumentName(name string) (ok bool) {
	_, err := getSupportedInstrumentByName(name)
	if err == nil {
		ok = true
	}
	return
}

func validRhythmPattern(name string) (ok bool) {
	switch name {
	case "steady":
		ok = true
	case "advancing":
		ok = true
	}
	return
}
