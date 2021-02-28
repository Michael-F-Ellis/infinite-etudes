// +build mage

package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Project directory tree. Values populated initPaths()
var (
	MageRoot     string // location of this file
	GoRoot       string // path to go installation
	AssetsPath   string // assets subdir
	InternalPath string // cmd/internal subdir
	// CommonPath   string // common subdir
	// ServerPath   string // server subdir
	// WasmPath     string // wasm subdir
)

func initPaths() {
	must := func(_err error) {
		if _err != nil {
			log.Fatal(_err)
		}
	}
	var err error
	GoRoot, err = sh.Output("go", "env", "GOROOT")
	must(err)
	MageRoot, err = os.Getwd()
	must(err)
	fmt.Println(MageRoot)
	AssetsPath = path.Join(MageRoot, "assets")
	InternalPath = path.Join(MageRoot, "internal")
	// CommonPath = path.Join(InternalPath, "common")
	// ServerPath = path.Join(MageRoot, "server")
	// WasmPath = path.Join(MageRoot, "wasm")
}

var Default = Build

func Build() {
	initPaths()
	must := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	must(mkWebPages())
	must(sh.Run("go", "build"))
}

func Run() {
	must := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	mg.Deps(Build)
	// launch the server
	must(sh.Run(path.Join(MageRoot, "infinite-etudes"), "-p", ":8081"))
}

func Clean() {
	initPaths()
	must := func(_err error) {
		if _err != nil {
			log.Fatal(_err)
		}
	}
	must(os.Remove(path.Join(MageRoot, "infinite-etudes")))
	must(os.Remove(path.Join(AssetsPath, "index.html")))
}
