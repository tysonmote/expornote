package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Some regexes to clean up the HTML a bit before sending it to pandoc. Some
	// of this could probably be handled by pandoc, but I'm lazy.
	stripFromHTML = []*regexp.Regexp{
		regexp.MustCompile(`<\?xml .+>`),
		regexp.MustCompile(`<!DOCTYPE .+>`),
		regexp.MustCompile(`<en-note>`),
		regexp.MustCompile(`</en-note>`),
		regexp.MustCompile(`<en-media .+/>`),
		regexp.MustCompile(`<br clear="none"/>`),
		regexp.MustCompile(` style="[^"]+"`),
		regexp.MustCompile(`<font [^>]+>`),
		regexp.MustCompile(`</font>`),
		regexp.MustCompile(`<div></div>`),
	}
)

type Note struct {
	Title     string      `xml:"title"`
	Content   []byte      `xml:"content"`
	Resources []*Resource `xml:"resource"`
}

func (n *Note) ExportTo(path string) {
	var counter int

	if markdown := n.ContentMarkdown(); len(markdown) > 0 {
		counter++

		destPath := fullFilePath(path, n.Title, ".md")
		file, err := os.Create(destPath)
		if err != nil {
			log.Printf("ERROR: Couldn't create %s: %s", destPath, err)
		} else {
			io.Copy(file, bytes.NewReader(markdown))
			file.Close()
			log.Printf("Created %s", destPath)
		}
	}

	for _, resource := range n.Resources {
		title := n.Title
		if counter > 0 {
			title = fmt.Sprintf("%s (%d)", title, counter)
		}

		resource.CopyTo(path, title)

		counter++
	}
}

func (n *Note) ContentMarkdown() []byte {
	content := n.Content

	for _, re := range stripFromHTML {
		content = re.ReplaceAll(content, []byte{})
	}

	content = bytes.TrimSpace(content)

	if len(content) == 0 {
		return []byte{}
	}

	// Use pandoc to convert to Markdown
	pandoc := exec.Command("pandoc", "--from", "html", "--to", "markdown_github-raw_html-hard_line_breaks+grid_tables", "--normalize", "--reference-links", "--reference-location", "block")
	pandoc.Stdin = bytes.NewReader(content)
	markdown, err := pandoc.Output()
	if err != nil {
		log.Printf("ERROR: Couldn't convert HTML to Markdown: %s", err)
	}

	return bytes.TrimSpace(markdown)
}

type Resource struct {
	Name string `xml:"resource-attributes>file-name"`
	Data []byte `xml:"data"`
}

func (r *Resource) CopyTo(path string, title string) {
	destPath := fullFilePath(path, title, filepath.Ext(r.Name))

	file, err := os.Create(destPath)
	if err != nil {
		log.Printf("ERROR: Couldn't create %s: %s", destPath, err)
		return
	}
	defer file.Close()

	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(r.Data))
	_, err = io.Copy(file, decoder)
	if err != nil {
		log.Printf("ERROR: Couldn't write to %s: %s", destPath, err)
		return
	}

	log.Printf("Created %s", destPath)
}

func fullFilePath(path, title, ext string) string {
	filename := fmt.Sprintf("%s%s", title, ext)
	filename = strings.Replace(filename, string(filepath.Separator), "-", -1)
	return filepath.Join(path, filename)
}
