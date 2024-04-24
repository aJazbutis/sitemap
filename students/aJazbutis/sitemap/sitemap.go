package sitemap

import (
	"bufio"
	"encoding/xml"
	"github.com/aJazbutis/link/students/ajazbutis/link"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

func MapSite(site string, depth int, toFile bool) {
	URL, err := url.Parse(site)
	if err != nil {
		panic(err)
	}
	visited := make(map[string]bool)
	toVisit := []link.Link{{Href: site}}
	uSet := urlSet{URL: URL}
	crawl(&uSet, visited, &toVisit, depth)
	linksToXml(&uSet, toFile)
}

func sameDomainLinks(resp *http.Response, visited map[string]bool, toVisit *[]link.Link, uSet *urlSet) {
	defer resp.Body.Close()
	links := link.ExtractLinks(resp.Body)
	for _, l := range links {
		href, err := url.Parse(l.Href)
		if err != nil {
			panic(err)
		}
		switch {
		case strings.HasPrefix(l.Href, "/"):
			href.Scheme = uSet.URL.Scheme
			href.Host = uSet.URL.Host
			l.Href = href.String()
		case href.Scheme != uSet.URL.Scheme || href.Hostname() != uSet.URL.Hostname():
			continue
		}
		if !visited[l.Href] {
			*toVisit = append(*toVisit, l)
		}
	}
}

func crawl(uSet *urlSet, visited map[string]bool, toVisit *[]link.Link, depth int) {
	for len(*toVisit) > 0 && depth >= -1 {
		level := len(*toVisit)
		for level > 0 {
			link := (*toVisit)[0]
			*toVisit = (*toVisit)[1:]
			if !visited[link.Href] {
				visited[link.Href] = true
				uSet.add(link)
				resp, err := http.Get(link.Href)
				if err != nil {
					panic(err)
				}
				sameDomainLinks(resp, visited, toVisit, uSet)
			}
			level--
		}
		depth--
	}
}

/*
**	type Link struct {
**		XMLName xml.Name `xml:"url"`
**		Href    string   `xml:"loc"`
**		Text    string   `xml:"-"`
**	}
 */

type urlSet struct {
	Urls  []link.Link
	Xmlns string   `xml:"xmlns,attr"`
	URL   *url.URL `xml:"-"`
}

func (uSet *urlSet) add(link link.Link) {
	uSet.Urls = append(uSet.Urls, link)
}

func outputFile(name string) *os.File {
	var b strings.Builder
	name = strings.Replace(name, ".", "_", -1)
	b.WriteString(name)
	b.WriteRune('_')
	b.WriteString(strconv.Itoa(rand.Intn(9999999999)))
	b.WriteString(".xml")
	file, err := os.OpenFile(b.String(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	return file
}

func linksToXml(uSet *urlSet, toFile bool /*, w io.Writer*/) {
	var w io.Writer
	if toFile {
		file := outputFile(uSet.URL.Hostname())
		defer file.Close()
		w = bufio.NewWriter(file)
	} else {
		w = os.Stdout
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	uSet.Xmlns = xmlns
	w.Write([]byte(xml.Header))
	enc.Encode(*uSet)
}
