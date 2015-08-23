package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type countingReader struct {
	io.Reader
	total int64
}

func (r *countingReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.total += int64(n)

	return n, err
}

var floatMatcher = regexp.MustCompile("(.)(\\d*\\.\\d*)(/unit)")

func price(s string) (float32, error) {
	matches := floatMatcher.FindStringSubmatch(s)
	if len(matches) != 4 {
		return 0, fmt.Errorf("No string value found in %s", s)
	}

	price, err := strconv.ParseFloat(matches[2], 32)
	if err != nil {
		return 0, err
	}

	return float32(price), nil
}

func scrapeItem(r io.Reader) (int64, string, error) {
	reader := &countingReader{r, 0}
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return reader.total, "", nil
	}

	return reader.total, strings.TrimSpace(doc.Find(".productText").First().Text()), nil
}

type partialItem struct {
	name  string
	price float32
	url   string
}

func scrapePage(r io.Reader) ([]partialItem, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scraping page: %v\n", err)
		return nil, err
	}

	items := []partialItem{}
	doc.Find(".product").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".productInfo h3").Text())
		price, err := price(strings.TrimSpace(s.Find(".pricePerUnit").Text()))
		url, ok := s.Find(".productInfo h3 a").Attr("href")

		if err == nil && ok {
			items = append(items, partialItem{name: name, price: price, url: url})
		}
	})

	return items, nil
}
