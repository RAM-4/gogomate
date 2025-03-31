package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type getter interface {
	Get(url string) (*http.Response, error)
}

type Scraper struct {
	http getter
}

func NewScraper(httpClient getter) *Scraper {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Scraper{
		http: httpClient,
	}
}

func (s *Scraper) Content(url string) (string, error) {
	doc, err := s.document(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch and parse document: %w", err)
	}

	clean(doc)
	content := extract(doc)

	return strings.Join(content, " "), nil
}

func (s *Scraper) document(url string) (*goquery.Document, error) {
	response, err := s.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

func clean(doc *goquery.Document) {
	doc.Find("script, style, noscript").Remove()
}

func extract(doc *goquery.Document) []string {
	var content []string
	doc.Find("section, article").Each(func(_ int, selection *goquery.Selection) {
		if text := cleanText(selection); text != "" {
			content = append(content, text)
		}
	})
	return content
}

func cleanText(selection *goquery.Selection) string {
	return strings.TrimSpace(selection.Text())
}
