package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// serveEtudes serves etude midi files from the current working directory.
func serveEtudes(hostport string, maxAgeSeconds int64) {
	os.Setenv("ETUDE_MAX_AGE", fmt.Sprintf("%d", maxAgeSeconds))
	defer os.Unsetenv("ETUDE_MAX_AGE")
	http.Handle("/", http.HandlerFunc(indexHndlr))
	http.Handle("/etude/", http.HandlerFunc(etudeHndlr))
	if err := http.ListenAndServe(hostport, nil); err != nil {
		log.Fatalf("Could not listen on port %s", hostport)
	}
}

// indexHndlr returns index.html
func indexHndlr(w http.ResponseWriter, r *http.Request) {
	text, err := ioutil.ReadFile("index.html")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(text)
}

// etudeHndlr returns a midi file that matches the get request or a 404 for
// incorrectly specified etudes. The pattern is
// /etude/<key>/<scale>/<instrument> where key is a pitchname like "c" or
// "aflat", <scale> is a scalename like "pentatonic", and instrument is a
// formatted General Midi instrument name like "acoustic_grand_piano".  If any
// of the forgoing are unknown or unsupported by this app, etudeHndlr gives a
// 400 response (StatusBadRequest). If the request matches a valid filename,
// the file will be returned in the response body if it exists and is younger
// than the maximum age imposed by this service. Otherwise the app will
// generate it so it can be returned.
func etudeHndlr(w http.ResponseWriter, r *http.Request) {
	what := strings.Split(r.URL.Path, "/")
	if what[1] != "etude" {
		panic("programming error. got request path that didn't start with 'etude'")
	}
	if !validEtudeRequest(what[2:]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filename := strings.Join(what[2:], "_") + ".mid"
	makeEtudesIfNeeded(filename, what[2])
	http.ServeFile(w, r, filename)
}

// makeEtudesIfNeeded generates a full set of etudes in the current
// working directory if the requested file doesn't exist or is older
// than the age limit set by serveEtudes in os.Environ. Otherwise
// it does nothing.
func makeEtudesIfNeeded(filename, instrument string) {
	var exists = true // initial assumption
	finfo, err := os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			// something's really wrong
			panic(fmt.Sprintf("error statting %s: %v", filename, err))
		}
		exists = false
	}
	// if file exists, we need to check its age
	if exists {
		maxage, e := strconv.Atoi(os.Getenv("ETUDE_MAX_AGE"))
		if e != nil {
			panic("programming error. can't convert ETUDE_MAX_AGE to integer")
		}

		maxduration := time.Duration(maxage) * time.Second
		modtime := finfo.ModTime()
		if time.Now().Sub(modtime) < maxduration {
			// nothing to do
			return
		}
	}
	// need to generate if we get to here
	// stubbed for testing ...
	txt := []byte("PHONY")
	err = ioutil.WriteFile(filename, txt, 0644)

}

// validEtudeRequest returns true if the request is correctly formed
// and references a valid etude file.
func validEtudeRequest(ksi []string) (ok bool) {
	if len(ksi) != 3 {
		return
	}
	if !validKeyName(ksi[0]) {
		return
	}
	if !validScaleName(ksi[1]) {
		return
	}
	if !validInstrumentName(ksi[2]) {
		return
	}

	ok = true
	return
}

// validKeyName returns true if the key name is in the ones we support.
func validKeyName(name string) (ok bool) {
	for _, k := range keyNames {
		if k == name {
			ok = true
			break
		}
	}
	return
}

// validScaleName returns true if the scale name is in the ones we support.
func validScaleName(name string) (ok bool) {
	scaleNames := []string{"final", "pentatonic"}
	for _, s := range scaleNames {
		if s == name {
			ok = true
			break
		}
	}
	return
}

// validInstrumentName returns true if the instrument name is in the ones we
// support.
func validInstrumentName(name string) (ok bool) {
	InstrumentNames := []string{"acoustic_grand_piano", "trumpet"}
	for _, i := range InstrumentNames {
		if i == name {
			ok = true
			break
		}
	}
	return
}
