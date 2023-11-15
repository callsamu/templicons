package templicons

import (
	"testing"
	"net/url"
)

func TestHasNoColorIfNotSet(t *testing.T) {
	p := &Parameters{}

	p.SetDimensions(10, 20)
	
	query := p.asQueryString()
	values, err := url.ParseQuery(query)
	if err != nil {
		t.Fatal(err)
	}

	if values.Get("width") != "10" {
		t.Fatalf("Expected 10, got %s", values.Get("width"))
	}

	if values.Get("height") != "20" {
		t.Fatalf("Expected 20, got %s", values.Get("height"))
	}

	if values.Get("color") != "" {
		t.Fatalf("Expected empty string, got %s", values.Get("color"))
	}
}

func TestHasNoDimensionsIfNotSet(t *testing.T) {
	p := &Parameters{}

	p.SetColor("#234891")
	
	query := p.asQueryString()
	values, err := url.ParseQuery(query)
	if err != nil {
		t.Fatal(err)
	}

	if values.Get("width") != "" {
		t.Fatalf("Expected empty string, got %s", values.Get("width"))
	}

	if values.Get("height") != "" {
		t.Fatalf("Expected empty string, got %s", values.Get("height"))
	}

	if values.Get("color") != "#234891" {
		t.Fatalf("Expected #234891, got %s", values.Get("color"))
	}
}
