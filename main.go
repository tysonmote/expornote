package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: expornote <.enex file> [<.enex file>]`)
	os.Exit(1)
}

func init() {
	if len(os.Args) < 2 {
		usage()
	}
}

func main() {
	for _, inPath := range os.Args[1:] {
		if len(inPath) == 0 {
			usage()
		}

		name := strings.Split(filepath.Base(inPath), ".")[0]
		outPath := filepath.Join(filepath.Dir(inPath), name)

		export(inPath, outPath)
	}
}

func export(inPath, outPath string) {
	archive, err := os.Open(inPath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	err = os.Mkdir(outPath, 0777)
	if err != nil {
		panic(err)
	}

	// Find all <note> elements, read them, and extract their contents.
	decoder := xml.NewDecoder(archive)
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		if token == nil {
			break
		}

		switch element := token.(type) {
		case xml.StartElement:
			if element.Name.Local == "note" {
				var n Note
				decoder.DecodeElement(&n, &element)
				n.ExportTo(outPath)
			}
		}
	}
}
