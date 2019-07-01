package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

var testhost = "localhost:8080"

func TestGETIndex(t *testing.T) {
	expbytes, _ := ioutil.ReadFile("index.html")
	exp := string(expbytes)
	url := "http://" + testhost
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("GET failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	got := string(body)
	if got != exp {
		t.Errorf("\nexp: %v\ngot: %v", exp, got)
	}

}

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
	var err error
	url := "http://" + testhost + "/etude/aflat/pentatonic/trumpet/steady"
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("GET failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}
	exp, _ := ioutil.ReadFile("aflat_pentatonic_trumpet_steady.mid")
	got, _ := ioutil.ReadAll(resp.Body)
	if !bytes.Equal(got, exp) {
		t.Errorf("response didn't match the file content")
	}
	// now test the age check
	maxage, _ := strconv.Atoi(os.Getenv("ETUDE_MAX_AGE"))
	maxduration := time.Duration(maxage) * time.Second
	time.Sleep(maxduration)
	resp2, err := http.Get(url)
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

func TestVocalEtudeRequest(t *testing.T) {
	// because multiple vocal parts are mapped to the same midi number
	var err error
	url := "http://" + testhost + "/etude/aflat/pentatonic/choir_aahs_tenor/advancing"
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("GET failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}
	exp, _ := ioutil.ReadFile("aflat_pentatonic_choir_aahs_tenor_advancing.mid")
	got, _ := ioutil.ReadAll(resp.Body)
	if !bytes.Equal(got, exp) {
		t.Errorf("response didn't match the file content")
	}
}

func TestBadEtudeRequest(t *testing.T) {
	badRequests := []string{
		"/etude/hsharp/pentatonic/trumpet",
		"/etude/aflat/schizotonic/trumpet",
		"/etude/aflat/pentatonic/fromix_horn",
	}
	for _, path := range badRequests {
		url := "http://" + testhost + path
		resp, err := http.Get(url)
		if err != nil {
			t.Errorf("GET failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("%s : xpected status code %v, got %v",
				path, http.StatusBadRequest, resp.StatusCode)
		}
	}
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Mkdir("test", 0777) // set up a temporary dir for generated files
	//var err error
	//html := []byte("<html></html>")
	//err = ioutil.WriteFile("test/index.html", html, 0644)

	// Run all tests and clean up
	wd, _ := os.Getwd()
	midijspath := filepath.Join(wd, "midijs")
	os.Chdir(filepath.Join(wd, "test"))
	go serveEtudes(testhost, 1, midijspath) // max etude age = 1 second so we don't wait forever while testing.
	exitcode := m.Run()
	os.Chdir(wd)
	os.RemoveAll("test") // remove the directory and its contents.
	os.Exit(exitcode)
}
