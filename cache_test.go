package templicons

import (
	"bytes"
	"context"
	"testing"
)

func TestCacheWithNoFallback(t *testing.T) {
	api := "https://api.iconify.design"
	icon := "mdi:home"

	client := MockClient{}
	cache := NewCache(api, client)

	c := cache.Icon(icon, nil)

	var buffer bytes.Buffer
	err := c.Render(context.Background(), &buffer)
	if err != nil {
		t.Fatal(err)
	}

	result := buffer.String()
	if result != "<svg></svg>" {
		t.Fatalf("Expected <svg></svg>, got %s", result)
	}

	url := "https://api.iconify.design/mdi/home"
	svg, ok := cache.cache[url]
	if !ok {
		t.Fatalf("Expected %s in cache, got nil", url)
	}

	if string(svg) != "<svg></svg>" {
		t.Fatalf("Expected <svg></svg>, got %s", string(svg))
	}
}

func TestCacheWithFallback(t *testing.T) {
	api := "https://api.iconify.design"
	icon := "mdi:home"

	client := MockClient{}
	cache := NewCache(api, client)

	c := cache.IconWithFallback(icon, "Home", nil)

	var buffer bytes.Buffer
	err := c.Render(context.Background(), &buffer)
	if err != nil {
		t.Fatal(err)
	}

	result := buffer.String()
	if result != "<span>Home</span>" {
		t.Fatalf("Expected <span>Home</span>, got %s", result)
	}

	log := <-cache.Logs
	if log != "Cached: https://api.iconify.design/mdi/home" {
		t.Fatalf("Expected 'Cached: https://api.iconify.design/mdi/home', got %s", log)
	}
}
