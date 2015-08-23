package scraper

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCountingReaderWrapsReaderAndCounts(t *testing.T) {
	tests := []struct {
		in  string
		out int64
	}{
		{"", 0},
		{"1234", 4},
	}

	for _, test := range tests {
		reader := &countingReader{strings.NewReader(test.in), 0}
		read, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Errorf("Unable to read from counting reader with value %v; received error %v", test.in, err)
		}
		if string(read[:]) != test.in {
			t.Errorf("ioUtil.ReadAll(%v) = %v; want %v", test.in, string(read[:]), test.in)
		}
		if reader.total != test.out {
			t.Errorf("reader.total = %v for %v; want %v", reader.total, test.in, test.out)
		}
	}
}

func TestParsesPage(t *testing.T) {
	file, err := os.Open("page.html")
	if err != nil {
		t.Errorf("Unable to open file: %v", err)
	}
	defer file.Close()

	items, err := scrapePage(file)

	if err != nil {
		t.Errorf("Error occured when scraping page: %v", err)
	}

	tests := []partialItem{
		{"Sainsbury's Apricot Ripe & Ready x5", 3.0, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-apricot-ripe---ready-320g"},
		{"Sainsbury's Avocado Ripe & Ready XL Loose 300g", 1.5, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-avocado-xl-pinkerton-loose-300g"},
		{"Sainsbury's Avocado, Ripe & Ready x2", 1.8, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-avocado--ripe---ready-x2"},
		{"Sainsbury's Avocados, Ripe & Ready x4", 3.2, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-avocados--ripe---ready-x4"},
		{"Sainsbury's Conference Pears, Ripe & Ready x4 (minimum)", 2.0, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-conference-pears--ripe---ready-x4-%28minimum%29"},
		{"Sainsbury's Kiwi Fruit, Ripe & Ready x4", 1.8, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-kiwi-fruit--ripe---ready-x4"},
		{"Sainsbury's Mango, Ripe & Ready x2", 2.0, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-mango--ripe---ready-x2"},
		{"Sainsbury's Nectarines, Ripe & Ready x4", 2.0, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-nectarines--ripe---ready-x4"},
		{"Sainsbury's Peaches Ripe & Ready x4", 2.0, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-peaches-ripe---ready-x4"},
		{"Sainsbury's Pears, Ripe & Ready x4 (minimum)", 2.0, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-pears--ripe---ready-x4-%28minimum%29"},
		{"Sainsbury's Plums Ripe & Ready x5", 2.5, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-plums--firm---sweet-x4-%28minimum%29"},
		{"Sainsbury's Ripe & Ready Golden Plums x6", 2.5, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-ripe---ready-golden-plums-x6"},
		{"Sainsbury's White Flesh Nectarines, Ripe & Ready x4", 2.0, "http://www.sainsburys.co.uk/shop/gb/groceries/ripe---ready/sainsburys-white-flesh-nectarines--ripe---ready-x4"},
	}

	if len(items) != len(tests) {
		t.Errorf("len(items) == %d, want %d", len(items), len(tests))
	}

	for index, test := range tests {
		if items[index] != test {
			t.Errorf("items[%d] = %v, want %v", index, items[index], test)
		}
	}
}

func TestParseItem(t *testing.T) {
	file, err := os.Open("item.html")
	if err != nil {
		t.Errorf("Unable to open file: %v", err)
	}
	defer file.Close()

	total, desc, err := scrapeItem(file)

	if err != nil {
		t.Error(err)
	}

	fileSize := int64(49716)
	if total != fileSize {
		t.Errorf("Page size for item = %d, want %d", total, fileSize)
	}
	description := "Apricots"
	if desc != description {
		t.Errorf("Description for item = %s, want %s", desc, description)
	}
}

func TestExtractsPriceFromString(t *testing.T) {
	tests := []struct {
		text string
		err  string
		out  float32
	}{
		{"abc", `No string value found in abc`, 0},
		{"£1.80/unit", "", 1.80},
		{"£1001.23/unit", "", 1001.23},
	}

	for _, test := range tests {
		actual, err := price(test.text)
		if err != nil {
			if err.Error() != test.err {
				t.Errorf("price(%v) produced error: %v, expected %v", test.text, err.Error(), test.err)
			}
		}
		if actual != test.out {
			t.Errorf("price(%v) = %v; wanted %v", test.text, actual, test.out)
		}
	}
}
