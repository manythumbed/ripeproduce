package scraper

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

type Product struct {
	Title       string  `json:"title"`
	Size        string  `json:"size"`
	UnitPrice   float32 `json:"unit_price"`
	Description string  `json:"description"`
}

type Results struct {
	Results []Product `json:"results"`
	Total   float32   `json:"total"`
}

func (r *Results) add(p Product) {
	r.Results = append(r.Results, p)
	r.Total += p.UnitPrice
}

type Scraper struct {
	Url string
}

func (s Scraper) Scrape() (Results, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Timeout: time.Second * 10, Jar: jar}

	items, err := fetchPage(client, s.Url)
	if err != nil {
		return Results{}, err
	}

	results := Results{}
	for _, item := range items {
		if "" == item.url {
			continue
		}

		size, desc, err := fetchItem(client, item.url)
		if err != nil {
			continue
		}

		results.add(Product{Title: item.name, Size: bytesToKilobytes(size), UnitPrice: item.price, Description: desc})
	}

	return results, nil
}

func bytesToKilobytes(size int64) string {
	return fmt.Sprintf("%3.2fkb", float64(size)/1024)
}

func fetchPage(client *http.Client, url string) ([]partialItem, error) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching page: %v\n", err)
		return []partialItem{}, err
	}

	defer resp.Body.Close()

	return scrapePage(resp.Body)
}

func fetchItem(client *http.Client, url string) (int64, string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return 0, "", err
	}

	defer resp.Body.Close()

	return scrapeItem(resp.Body)
}
