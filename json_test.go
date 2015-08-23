package scraper

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestProductSerialisesAsJson(t *testing.T) {
	tests := []struct {
		item Product
		out  string
	}{
		{Product{}, `{"title":"","size":"","unit_price":0,"description":""}`},
		{Product{Title: "title", Size: "120kb", UnitPrice: 12.34, Description: "A test product"}, `{"title":"title","size":"120kb","unit_price":12.34,"description":"A test product"}`},
	}

	for _, test := range tests {
		actual, err := json.Marshal(test.item)
		if err != nil {
			t.Errorf("json.Marshal(%v); want %v, error produced %v", test.item, test.out, err)
		}
		if !bytes.Equal(actual, []byte(test.out)) {
			t.Errorf("json.Marshal(%v) = %v; want %v", test.item, string(actual[:]), test.out)
		}
	}
}

func TestResultsSerialisesAsJson(t *testing.T) {
	tests := []struct {
		result Results
		out    string
	}{
		{Results{}, `{"results":null,"total":0}`},
		{Results{Results: []Product{{Title: "title", Size: "120kb", UnitPrice: 12.34, Description: "A test product"}}, Total: 12.34},
			`{"results":[{"title":"title","size":"120kb","unit_price":12.34,"description":"A test product"}],"total":12.34}`},
	}

	for _, test := range tests {
		actual, err := json.Marshal(test.result)
		if err != nil {
			t.Errorf("json.Marshal(%v); want %v, error produced %v", test.result, test.out, err)
		}
		if !bytes.Equal(actual, []byte(test.out)) {
			t.Errorf("json.Marshal(%v) = %v; want %v", test.result, string(actual[:]), test.out)
		}
	}
}
