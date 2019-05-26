package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
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

func TestGoodEtudeRequest(t *testing.T) {
	url := "http://" + testhost + "/etude/aflat/pentatonic/trumpet"
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("GET failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}
}

func TestBadEtudeRequest(t *testing.T) {
	url := "http://" + testhost + "/etude/aflat/pentatonic/fromix_horn"
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("GET failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Mkdir("test", 0777) // set up a temporary dir for generated files
	var err error
	html := []byte("<html></html>")
	err = ioutil.WriteFile("test/index.html", html, 0644)
	if err != nil {
		fmt.Printf("could not write test/index.html : %v", err)
		os.Exit(1)
	}

	// Run all tests and clean up
	wd, _ := os.Getwd()
	os.Chdir(path.Join(wd, "test"))
	go serveEtudes(testhost, 1) // max etude age = 1 second so we don't wait forever while testing.
	exitcode := m.Run()
	os.Chdir(wd)
	os.RemoveAll("test") // remove the directory and its contents.
	os.Exit(exitcode)
}
