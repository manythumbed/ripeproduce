package main

import (
	"encoding/json"
	"fmt"
	"github.com/manythumbed/ripeproduce"
	"os"
	"flag"
)

const (
	defaultUrl = "http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?listView=true&orderBy=FAVOURITES_FIRST&parent_category_rn=12518&top_category=12518&langId=44&beginIndex=0&pageSize=20&catalogId=10137&searchTerm=&categoryId=185749&listId=&storeId=10151&promotionId=#langId=44&storeId=10151&catalogId=10137&categoryId=185749&parent_category_rn=12518&top_category=12518&pageSize=20&orderBy=FAVOURITES_FIRST&searchTerm=&beginIndex=0&hideFilters=true"
)

var url = flag.String("url", defaultUrl, "URL to scrape")

func main() {
	flag.Parse()

	scraper := scraper.Scraper{Url: *url}

	results, err := scraper.Scrape()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem scraping: %v", err)
		return
	}

	encoder := json.NewEncoder(os.Stdout)
	if encoder.Encode(results) != nil {
		fmt.Fprintf(os.Stderr, "Unable to encode results as json: %v", err)
	}
}
