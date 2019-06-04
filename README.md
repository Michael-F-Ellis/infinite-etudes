# Infinite Etudes
*infinite-etudes* generates ear training exercises for instrumentalists.

You can run it from the command line (cli mode) or as a web server (server mode).

In cli mode, *infinite-etudes* generates a set of 7 midi files for each of 12 key
signatures. Each set covers all possible combinations of 3 pitches within the
key. The files are generated in the current working directory.

In server mode, *infinite-etudes* is a high-performance self-contained web server
that provides a simple user interface that allows the user to choose a key, a
scale pattern and an instrument sound and play a freshly-generated etude in
the web browser. A public demo instance is running at 

https://etudes.ellisandgrant.com

## Installation
You need to have Go installed to build and test infinite-etudes. Get it from https://golang.org/dl/ .

After installing Go, do

```
  go get github.com/Michael-F-Ellis/infinite-etudes
  cd $GOPATH/src/github.com/Michael-F-Ellis/infinite-etudes
  go test
  go install
```

Then run `infinite-etudes -h` for options and usage instructions.
