package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Embedded static assets.
//
//go:embed img grooveclock midijs ytshed
var Static embed.FS

// serveEtudes serves etude midi files from the current working directory.
func serveEtudes(hostport string) {
	err := mkWebPages()
	if err != nil {
		log.Fatalf("could not write web pages: %v", err)
	}

	http.Handle("/etude/", http.HandlerFunc(etudeHndlr))
	http.Handle("/grooveclock/", http.HandlerFunc(grooveclockHndlr))
	http.Handle("/ytshed/", http.HandlerFunc(ytshedHndlr))
	http.Handle("/img/", http.HandlerFunc(imgHndlr))
	http.Handle("/midijs/", http.HandlerFunc(midijsHndlr))
	log.Printf("midijs path is %s", os.Getenv("MIDIJS"))
	http.Handle("/", http.HandlerFunc(indexHndlr))
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
func imgHndlr(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Split(r.URL.Path, "/")
	if urlPath[1] != "img" {
		log.Fatalf("programming error. got img asset request path that didn't start with 'img': %s", r.URL.Path)
	}
	subpath := path.Join(urlPath[1:]...)
	content, err := Static.ReadFile(subpath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%v", err)
		return
	}
	setContentType(w, subpath)
	_, err = w.Write(content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v", err)
	}
}

// setContentType sets the content type of the response based on the filename.
func setContentType(w http.ResponseWriter, filename string) {
	switch {
	case strings.HasSuffix(filename, ".html"):
		w.Header().Set("Content-Type", "text/html")
	case strings.HasSuffix(filename, ".css"):
		w.Header().Set("Content-Type", "text/css")
	case strings.HasSuffix(filename, ".js"):
		w.Header().Set("Content-Type", "application/javascript")
	case strings.HasSuffix(filename, ".mid"):
		w.Header().Set("Content-Type", "audio/midi")
	case strings.HasSuffix(filename, ".png"):
		w.Header().Set("Content-Type", "image/png")
	case strings.HasSuffix(filename, ".ico"):
		w.Header().Set("Content-Type", "image/x-icon")
	default:
		w.Header().Set("Content-Type", "text/plain")
	}
}

// grooveclockHndlr returns files from the MIDIJS directory.
func grooveclockHndlr(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Split(r.URL.Path, "/")
	if urlPath[1] != "grooveclock" {
		log.Fatalf("programming error. got grooveclock asset request path that didn't start with 'grooveclock': %s", r.URL.Path)
	}
	subpath := path.Join(urlPath[1:]...)
	content, err := Static.ReadFile(subpath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("subpath %s: %v", subpath, err)
		return
	}
	setContentType(w, subpath)
	_, err = w.Write(content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v", err)
	}
}

// midijsHndlr returns files from the MIDIJS directory.
func midijsHndlr(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Split(r.URL.Path, "/")
	if urlPath[1] != "midijs" {
		log.Fatalf("programming error. got midijs asset request path that didn't start with 'midijs': %s", r.URL.Path)
	}
	subpath := path.Join(urlPath[1:]...)
	content, err := Static.ReadFile(subpath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("subpath %s: %v", subpath, err)
		return
	}
	setContentType(w, subpath)
	_, err = w.Write(content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v", err)
	}
}

type etudeRequest struct {
	tonalCenter string
	pattern     string
	interval1   string
	interval2   string
	interval3   string
	instrument  string
	tempo       string // beats per minute
	repeats     int    // number of repeats (0-3)
	metronome   int    // On, DownbeatOnly, Off
	silent      int    // true indicated the corresponding repeat should be silent

}

const (
	metronomeOn int = iota
	metronomeDownbeatOnly
	metronomeOff
)

// metronomeString returns a string representation of the metronome integer value.
func metronomeString(req *etudeRequest) (s string) {
	switch req.metronome {
	case metronomeOn:
		s = "on"
	case metronomeDownbeatOnly:
		s = "downbeat"
	case metronomeOff:
		s = "off"
	default:
		s = "invalid"
	}
	return
}
func (r *etudeRequest) midiFilename() (f string) {
	var parts []string
	repeats := fmt.Sprintf("%d", r.repeats)
	silence := fmt.Sprintf("%d", r.silent)

	switch r.pattern {
	case "interval":
		parts = []string{r.pattern, r.interval1, r.instrument, metronomeString(r), r.tempo, repeats, silence}
	case "intervalpair", "intervalpair_ud":
		parts = []string{r.pattern, r.interval1, r.interval2, r.instrument, metronomeString(r), r.tempo, repeats, silence}
	case "intervaltriple", "intervaltriple_ud":
		parts = []string{r.pattern, r.interval1, r.interval2, r.interval3, r.instrument, metronomeString(r), r.tempo, repeats, silence}
	default:
		parts = []string{r.tonalCenter, r.pattern, r.instrument, metronomeString(r), r.tempo, repeats, silence}
	}
	f = strings.Join(parts, "_") + ".mid"
	return
}

// etudeHndlr returns a midi file that matches the get request or a 404 for
// incorrectly specified etudes. If the request is valid and the file exists
// already, it will be returned in the response body if it is younger than the
// maximum age imposed by this service. Otherwise the app will generate it so it
// can be returned.
func etudeHndlr(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 12 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
	req.interval3 = path[6]
	req.instrument = path[7]
	switch path[8] {
	case "on":
		req.metronome = metronomeOn
	case "downbeat":
		req.metronome = metronomeDownbeatOnly
	case "off":
		req.metronome = metronomeOff
	default:
		req.metronome = 4 // invalid
	}
	req.tempo = path[9]
	repeats, err := strconv.Atoi(path[10])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(`can't convert "%s" to repeat count: %v`, path[10], err)
		return
	}
	req.repeats = repeats
	req.silent, err = strconv.Atoi(path[11])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !validEtudeRequest(req) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filename := (&req).midiFilename()
	log.Printf("%s requested", filename)
	makeEtudesIfNeeded(filename, req)
	http.ServeFile(w, r, filename)
	// log the request in format that's convenient for analysis
	log.Printf("%s %s served\n", r.RemoteAddr, filename)
}

// ytshedHndlr returns ytshed/index.html
func ytshedHndlr(w http.ResponseWriter, r *http.Request) {
	var err error
	log.Printf("about to serve %s", "ytshed/index.html")
	bytes, err := Static.ReadFile("ytshed/index.html")
	if err != nil {
		log.Printf("Could not read index.html: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write(bytes)
}

// removeExpiredMidiFiles deletes midi files in the current working
// directory that are older than expireSeconds
func removeExpiredMidiFiles() {
	fnames, _ := filepath.Glob("*.mid")
	for _, fname := range fnames {
		finfo, err := os.Stat(fname)
		if err != nil {
			// something's really wrong
			log.Fatalf("error statting %s: %v", fname, err)
		}
		if time.Since(finfo.ModTime()) > time.Duration(expireSeconds)*time.Second {
			// file expired, remove it
			err = os.Remove(fname)
			if err != nil {
				// something's really wrong
				log.Fatalf("error removing %s: %v", fname, err)
			}
		}
	}

}

var etudeMutex sync.Mutex

// makeEtudesIfNeeded generates a full set of etudes in the current
// working directory if the requested file doesn't exist or is older
// than the age limit set by serveEtudes in os.Environ. Otherwise
// it does nothing.
func makeEtudesIfNeeded(filename string, req etudeRequest) {
	// use the mutex to ensure that multiple requests can't
	// create or delete files while a request is in process
	etudeMutex.Lock()
	defer etudeMutex.Unlock()
	// Remove all expired files
	removeExpiredMidiFiles()

	// See if file exists
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) { // See https://gist.github.com/mattes/d13e273314c3b3ade33f to understand this slightly weird test
		// file exists, nothing to do
		return
	}
	// need to generate if we get to here
	iInfo, _ := getSupportedInstrumentByName(req.instrument) // already validated. ignore err value
	// fmt.Printf("%v %s\n", iInfo, filename)
	instrument := iInfo.gmnumber - 1
	midilo := iInfo.midilo
	midihi := iInfo.midihi
	tempo, _ := strconv.Atoi(req.tempo)
	mkRequestedEtude(midilo, midihi, tempo, instrument, req)
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
	case "intervalpair", "intervalpair_ud":
		if !validIntervalName(req.interval1) || !validIntervalName(req.interval2) {
			return
		}
	case "intervaltriple", "intervaltriple_ud":
		if !validIntervalName(req.interval1) ||
			!validIntervalName(req.interval2) ||
			!validIntervalName(req.interval3) {
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
	if !validMetronomePattern(metronomeString(&req)) {
		return
	}
	if !validTempo(req.tempo) {
		return
	}

	ok = true
	return
}
func rangeCheck(req etudeRequest) (ok bool, err error) {
	// Now check that the midi range is large enough to contain the etude.
	var minsize int
	switch req.pattern {
	case "allintervals":
		// need at least 2 octaves
		minsize = 24
	case "interval":
		// need at least 1 octave plus the size of the interval
		minsize = 12 + intervalSizeByName(req.interval1)
	case "intervalpair", "intervalpair_ud":
		// need at least 1 octave plus the sum of the intervals
		minsize = 12 +
			intervalSizeByName(req.interval1) +
			intervalSizeByName(req.interval2)
	case "intervaltriple", "intervaltriple_ud":
		// need at least 1 octave plus the sum of the intervals
		minsize = 12 +
			intervalSizeByName(req.interval1) +
			intervalSizeByName(req.interval2) +
			intervalSizeByName(req.interval3)
	}
	iInfo, err := getSupportedInstrumentByName(req.instrument)
	irange := iInfo.midihi - iInfo.midilo
	ok = irange >= minsize
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

// intervalSizeByName returns the size of name in half-steps
func intervalSizeByName(name string) (sz int) {
	for _, inf := range intervalInfo {
		if inf.fileName == name {
			sz = inf.size
			break
		}
	}
	return
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

var patternInfo = []nameInfo{
	{"interval", "One Interval", "One Interval", 0},
	{"allintervals", "Tonic Intervals", "Tonic Intervals", 0},
	{"intervalpair", "Two Intervals", "Two Intervals", 0},
	{"intervalpair_ud", "Two Intervals Up/Down", "Two Intervals", 0},
	{"intervaltriple", "Three Intervals", "Three Intervals", 0},
	{"intervaltriple_ud", "Three Intervals Up/Down", "Three Intervals", 0},
}

// validPattern returns true if the scale name is in the ones we support.
func validPattern(name string) (ok bool) {
	for _, s := range patternInfo {
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

func validMetronomePattern(name string) (ok bool) {
	switch name {
	case "on", "downbeat", "off":
		ok = true
	}
	return
}
func validTempo(ts string) (ok bool) {
	_, err := strconv.Atoi(ts)
	if err == nil {
		ok = true
	}
	return
}
