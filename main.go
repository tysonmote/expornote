package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	importPath string
	exportPath string
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: expornote <.enex file>`)
	os.Exit(1)
}

func init() {
	if len(os.Args) != 2 {
		usage()
	}

	importPath = os.Args[1]
	if len(importPath) == 0 {
		usage()
	}

	path := filepath.Dir(importPath)
	name := strings.Split(filepath.Base(importPath), ".")[0]
	exportPath = filepath.Join(path, name)
}

func main() {
	archive, err := os.Open(importPath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	err = os.Mkdir(exportPath, 0777)
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
				n.ExportTo(exportPath)
			}
		}
	}

}
