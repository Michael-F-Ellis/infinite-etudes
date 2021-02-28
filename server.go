package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Michael-F-Ellis/infinite-etudes/internal/valid"
)

// Bundle our static files with the app
//go:embed assets
var assets embed.FS

// serveEtudes serves etude midi files from the current working directory.
func serveEtudes(hostport string) {
	var err error
	mux := http.NewServeMux()
	mux.HandleFunc("/etude/", etudeHndlr)
	assetSys, err := fs.Sub(assets, "assets")
	if err != nil {
		log.Fatalf("could not create assets subtree: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(assetSys)))
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
	wd, _ := os.Getwd()
	log.Printf("serving on %s\n from %s", hostport, wd)
	switch serveSecure {
	case true:
		if err := http.ListenAndServeTLS(hostport, certpath, certkeypath, mux); err != nil {
			log.Fatalf("Could not listen on port %s : %v", hostport, err)
		}
	default:
		if err := http.ListenAndServe(hostport, mux); err != nil {
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
	case "intervalpair":
		parts = []string{r.pattern, r.interval1, r.interval2, r.instrument, metronomeString(r), r.tempo, repeats, silence}
	case "intervaltriple":
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
	iInfo, _ := valid.InstrumentByName(req.instrument) // already validated. ignore err value
	// fmt.Printf("%v %s\n", iInfo, filename)
	instrument := iInfo.GMNumber - 1
	midilo := iInfo.MidiLo
	midihi := iInfo.MidiHi
	tempo, _ := strconv.Atoi(req.tempo)
	mkRequestedEtude(midilo, midihi, tempo, instrument, req)
}

// validEtudeRequest returns true if the request is correctly formed
// and references a valid etude filename.
func validEtudeRequest(req etudeRequest) (ok bool) {
	if !valid.Pattern(req.pattern) {
		return
	}

	switch req.pattern { // Intervals get special handling
	case "allintervals":
		if !valid.KeyName(req.tonalCenter) {
			return
		}
	case "interval":
		if !valid.IntervalName(req.interval1) {
			return
		}
	case "intervalpair":
		if !valid.IntervalName(req.interval1) || !valid.IntervalName(req.interval2) {
			return
		}
	case "intervaltriple":
		if !valid.IntervalName(req.interval1) ||
			!valid.IntervalName(req.interval2) ||
			!valid.IntervalName(req.interval3) {
			return
		}

	default:
		if !valid.KeyName(req.tonalCenter) {
			return
		}
	}
	if !valid.InstrumentName(req.instrument) {
		return
	}
	if !valid.MetronomePattern(metronomeString(&req)) {
		return
	}
	if !valid.Tempo(req.tempo) {
		return
	}
	ok = true
	return
}

// intervalSizeByName returns the size of name in half-steps
func intervalSizeByName(name string) (sz int) {
	for _, inf := range valid.IntervalInfo {
		if inf.FileName == name {
			sz = inf.Size
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
