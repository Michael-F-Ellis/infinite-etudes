package main

import (
	"bytes"
	"io/ioutil"

	. "github.com/Michael-F-Ellis/infinite-etudes/internal/ht"
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
	viewport := Meta(`name="viewport" content="width=device-width, initial-scale=1"`)
	description := Meta(`name="description", content="Infinite Etudes demo"`)
	keywords := Meta(`name="keywords", content="music,notation,midi,tbon"`)
	head := Head("", viewport, description, keywords, indexCSS())

	// <body>
	body := Body("", SC("Watch this space for new developments coming very soon!"))

	// <html>
	page := Html("", head, body)
	page.Render(&buf, 0)
	err = ioutil.WriteFile("index.html", buf.Bytes(), 0644)
	return
}

func indexCSS() *ElementTree {
	return Style("", SC(`
    body {margin: 0; height: 100%; overflow: hidden}
    h1 {font-size: 300%; margin-bottom: 1vh}
    h2 {font-size: 200%}
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
    select {margin-left: 5%; margin-bottom: 1%}
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
	`))
}
