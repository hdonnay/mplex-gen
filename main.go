// mplex-gen generates code for a channel multiplexer.
//
// The function generated accepts a slice of channels, and spawns n+1 goroutines
// to service the supplied channel.
// The returned channel will be closed if all supplied channels are closed.
// It is the caller's responsibility to service the returned channel. Failure to
// do so will put backpressue on the supplier channels.
//
// This tool is meant to be called via `go generate`:
//
//	//go:generate mplex-gen -o mplex.go main *bytes.Buffer
//
package main

// I'm lazy.
//go:generate sh -c "godoc . | sed '1d;2d;s/^    //' > README.md"

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

var (
	outfile = flag.String("o", "", "file to output to")

	tmpl *template.Template
)

func main() {
	var out *os.File
	var err error
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatalln("invalid arguments")
	}

	pkgName := flag.Arg(0)
	chType := flag.Arg(1)

	if *outfile == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(*outfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if chType == "" {
		log.Fatal("type not specified")
	}

	info := struct {
		PkgName   string
		ChType    string
		PrintName string
	}{
		pkgName,
		chType,
		printableName(chType),
	}

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, info); err != nil {
		log.Fatal(err)
	}
	b, err := imports.Process(*outfile, buf.Bytes(), nil)
	if err != nil {
		log.Fatal(err)
	}
	buf = bytes.NewBuffer(b)

	if _, err := io.Copy(out, buf); err != nil {
		log.Fatal(err)
	}
}

func printableName(n string) string {
	n = strings.Replace(n, "*", "", -1)  // pointer
	n = strings.Replace(n, "[]", "", -1) // slice
	n = strings.Replace(n, ".", " ", -1) // package identifiers
	n = strings.Replace(n, "[", " ", -1) // map
	n = strings.Replace(n, "]", " ", -1) // map
	n = strings.Title(n)                 // change case
	n = strings.Replace(n, " ", "", -1)  // crush
	return n
}
