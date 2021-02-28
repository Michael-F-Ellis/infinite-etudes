package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var testhost = "localhost:8080"

func TestMidijsRequest(t *testing.T) {
	url := "http://" + testhost + "/midijs/pat/arachno-0.pat"
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("GET failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	fmt.Println(len(body))

	//exp := 0
	//got, err := what()
	//if err != nil {
	//	t.Errorf("describe: %v", err)
	//}
	//if got != exp {
	//	t.Errorf("\nexp: %v\ngot: %v", exp, got)
	//}
}

func TestGoodEtudeRequest(t *testing.T) {
	type testcase struct {
		url      string
		filename string
	}
	testTable := []testcase{
		{
			url:      "http://" + testhost + "/etude/aflat/allintervals/minor2/minor2/minor2/trumpet/on/120/3/0",
			filename: "aflat_allintervals_trumpet_on_120_3_0.mid",
		},
		{
			url:      "http://" + testhost + "/etude/aflat/intervalpair/minor2/minor2/minor2/trumpet/on/120/1/0",
			filename: "intervalpair_minor2_minor2_trumpet_on_120_1_0.mid",
		},
	}
	for _, tcase := range testTable {
		resp, err := http.Get(tcase.url)
		if err != nil {
			t.Errorf("GET failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
		}
		exp, _ := ioutil.ReadFile(tcase.filename)
		got, _ := ioutil.ReadAll(resp.Body)
		if !bytes.Equal(got, exp) {
			t.Errorf("response didn't match the file content")
		}
		// now test the age check
		time.Sleep(time.Duration(expireSeconds) * time.Second)
		resp2, err := http.Get(tcase.url)
		if err != nil {
			t.Errorf("GET failed: %v", err)
		}
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, resp2.StatusCode)
		}
		got, _ = ioutil.ReadAll(resp2.Body)
		if bytes.Equal(got, exp) { // exp is unchanged and should not match got.
			t.Errorf("file did not update")
		}
	}
}

func TestVocalEtudeRequest(t *testing.T) {
	// because multiple vocal parts are mapped to the same midi number
	var err error
	url := "http://" + testhost + "/etude/aflat/allintervals/minor2/minor2/minor2/choir_aahs_tenor/off/120/3/0"
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("GET failed: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}
	exp, _ := ioutil.ReadFile("aflat_allintervals_choir_aahs_tenor_off_120_3_0.mid")
	got, _ := ioutil.ReadAll(resp.Body)
	if !bytes.Equal(got, exp) {
		t.Errorf("response didn't match the file content")
	}
}
func TestValidEtudeRequest(t *testing.T) {
	badRequests := []etudeRequest{
		{tonalCenter: "hsharp", pattern: "pentatonic", instrument: "trumpet", tempo: "120"},
	}
	for _, req := range badRequests {
		ok := validEtudeRequest(req)
		if ok {
			t.Errorf("request should not have succeeded:\n%v", req)
		}
	}
	goodRequests := []etudeRequest{
		{tonalCenter: "", pattern: "intervalpair", interval1: "minor3", interval2: "major3", instrument: "trumpet", metronome: metronomeDownbeatOnly, tempo: "120"},
		{tonalCenter: "", pattern: "intervaltriple", interval1: "minor3", interval2: "major3", interval3: "minor3", instrument: "trumpet", metronome: metronomeOff, tempo: "120"},
	}
	for _, req := range goodRequests {
		ok := validEtudeRequest(req)
		if !ok {
			t.Errorf("request should have succeeded:\n%v", req)
		}
	}

}
func TestBadEtudeRequest(t *testing.T) {
	badRequests := []string{
		"/etude/c/pentatonic/minor2/minor2/minor2/trumpet/on/120",            // no repeat count
		"/etude/hsharp/pentatonic/minor2/minor2/minor2/trumpet/on/120/3",     // bad tonal center
		"/etude/c/schizotonic/minor2/minor2/minor2/trumpet/on/120/3",         // bad pattern
		"/etude/c/interval/fermented2/minor2/minor2/trumpet/on/120/3",        // bad interval1
		"/etude/c/intervalpairs/minor2/minor2/toxic2/trumpet/on/120/3",       // bad interval2
		"/etude/c/pentatonic/minor2/minor2/toxic2/fromixhorn/on/120/3",       // bad instrument
		"/etude/c/pentatonic/minor2/minor2/minor2/trumpet/jittery/120/3",     // bad rhythm
		"/etude/c/pentatonic/minor2/minor2/minor2/trumpet/on/allaregretto/3", // bad tempo
	}
	for _, path := range badRequests {
		url := "http://" + testhost + path
		resp, err := http.Get(url) // TODO #1 This is returning a nil response.
		if err != nil {
			t.Errorf("GET failed: %v", err)
			continue
		} else {
			defer resp.Body.Close()
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("%s : xpected status code %v, got %v",
				path, http.StatusBadRequest, resp.StatusCode)
		}
	}
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	// set up a temporary dir for generated files
	_ = os.RemoveAll("test") // in case prior run crashed
	err := os.Mkdir("test", 0777)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(-1)
	}
	//html := []byte("<html></html>")
	//err = ioutil.WriteFile("test/index.html", html, 0644)

	// Run all tests and clean up
	wd, _ := os.Getwd()
	err = os.Chdir(filepath.Join(wd, "test"))
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(-1)
	}
	expireSeconds = 1
	go serveEtudes(testhost) // max etude age = 1 second so we don't wait forever while testing.
	exitcode := m.Run()
	err = os.Chdir(wd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(-1)
	}
	os.RemoveAll("test") // remove the directory and its contents.
	os.Exit(exitcode)
}

func BenchmarkMkAllEtudes(b *testing.B) {
	req := etudeRequest{instrument: "viola", pattern: "allintervals", tonalCenter: "c"}
	for i := 0; i < b.N; i++ {
		mkRequestedEtude(48, 84, 120, 15, req)
	}
}
