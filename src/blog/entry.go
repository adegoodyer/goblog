// Package blog contains blog content and entries
package blog

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"github.com/mariomac/goblog/src/logr"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	nethtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var log = logr.Get()

// Entry holds the information of a blog entry or page
type Entry struct {
	FilePath string
	Title    string
	HTML     template.HTML
	Preview  template.HTML
	Time     time.Time // may be zero, if it is a page
}

// YYYYMMDDHHMMsome-text_here.md
var entryFormat = regexp.MustCompile(`[0-9]{12}[_\-a-zA-Z0-9]+\.md$`)

// LoadEntry loads and renders a blog entry given a file path
func LoadEntry(file string) (*Entry, error) {
	llog := log.WithField("file", file)
	llog.Debug("loading blog Entry")

	fileBody, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	// TODO: support parsing
	timestamped := entryFormat.FindString(file)
	var timestamp time.Time
	if len(timestamped) > 0 {
		timestamp = extractTime(timestamped)
	}
	title, html, preview := getTitleBodyAndPreview(fileBody)
	return &Entry{
		Time:     timestamp,
		Title:    title,
		FilePath: file,
		HTML:     html,
		Preview:  preview,
	}, nil
}

// TODO: configure by env
var location, _ = time.LoadLocation("Europe/Madrid")

// Extracts a Time from a string beggining with a timestamp in the format YYYYMMDDHHMM...
func extractTime(timestr string) time.Time {
	year, _ := strconv.Atoi(timestr[:4])
	month, _ := strconv.Atoi(timestr[4:6])
	day, _ := strconv.Atoi(timestr[6:8])
	hour, _ := strconv.Atoi(timestr[8:10])
	minute, _ := strconv.Atoi(timestr[10:12])
	parsedTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, location)
	return parsedTime
}

func getTitleBodyAndPreview(mdBytes []byte) (string, template.HTML, template.HTML) {
	// TODO: proper caching of goldmark
	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("tango"),
			),
		),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	htmlBytes := bytes.Buffer{}
	if err := markdown.Convert(mdBytes, &htmlBytes); err != nil {
		// TODO: properly log/manage blogerr
		htmlBytes = bytes.Buffer{}
		htmlBytes.WriteString(`<h1>Error parsing markdown</h1><p>` + err.Error() + `</p>`)
	}
	htmlNode, err := nethtml.Parse(bytes.NewReader(htmlBytes.Bytes()))

	// TODO: properly handle error
	if err != nil {
		return err.Error(), "", ""
	}

	firstParagraph := getFirstParagraph(htmlNode)

	h1 := removeFirstH1(htmlNode)
	title, _ := getText(h1)
	log.Debugf("Parsed title: %s", title)

	bodyBuf := new(bytes.Buffer)
	nethtml.Render(bodyBuf, htmlNode)
	body := template.HTML(bodyBuf.String())

	var preview template.HTML
	if firstParagraph != nil {
		previewBuf := new(bytes.Buffer)
		nethtml.Render(previewBuf, firstParagraph)
		preview = template.HTML(previewBuf.String())
	} else {
		preview = ""
	}

	return title, body, preview
}

// Parameter, parent node. Return type, removed node
func removeFirstH1(parent *nethtml.Node) *nethtml.Node {
	child := parent.FirstChild
	for child != nil {
		if child.DataAtom == atom.H1 {
			parent.RemoveChild(child)
			return child
		}
		removedH1 := removeFirstH1(child)
		if removedH1 != nil {
			return removedH1
		}
		child = child.NextSibling
	}
	return nil
}

func getFirstParagraph(parent *nethtml.Node) *nethtml.Node {
	child := parent.FirstChild
	for child != nil {
		if child.DataAtom == atom.P {
			return child
		}
		paragraph := getFirstParagraph(child)
		if paragraph != nil {
			return paragraph
		}
		child = child.NextSibling
	}
	return nil
}

func getText(parent *nethtml.Node) (string, bool) {
	child := parent.FirstChild
	for child != nil {
		if child.Type == nethtml.TextNode && child.FirstChild == nil {
			return child.Data, true
		}
		text, found := getText(child)
		if found {
			return text, true
		}
		child = child.NextSibling
	}
	return "", false
}
