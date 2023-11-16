package templicons

import (
	"bytes"
	"context"
	"testing"
)

func TestCacheWithoutFallback(t *testing.T) {
	api := "https://api.iconify.design"
	icon := "mdi:home"

	client := MockClient{}
	cache := NewCache(api, client, DefaultFallback)

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

	url := "https://api.iconify.design/mdi/home.svg"
	svg, ok := cache.cache[url]
	if !ok {
		t.Fatalf("Expected %s in cache, got nil", url)
	}

	if string(svg) != "<svg></svg>" {
		t.Fatalf("Expected <svg></svg>, got %s", string(svg))
	}

	var buffer2 bytes.Buffer
	err = c.Render(context.Background(), &buffer2)
	if err != nil {
		t.Fatal(err)
	}

	result2 := buffer2.String()
	if result2 != "<svg></svg>" {
		t.Fatalf("Expected <svg></svg>, got %s", result2)
	}

	if client.calls > 1 {
		t.Fatalf("Expected a single client fetch call, got %d", client.calls)
	}
}

func TestCacheWithFallback(t *testing.T) {
	api := "https://api.iconify.design"
	icon := "mdi:home"

	client := MockClient{}
	cache := NewCache(api, client, DefaultFallback)

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

	cache.wg.Wait()
	url := "https://api.iconify.design/mdi/home.svg"
	if cache.cache[url] == nil {
		t.Fatalf("Expected %s in cache, got nil", url)
	}

	var buffer2 bytes.Buffer
	err = c.Render(context.Background(), &buffer2)
	if err != nil {
		t.Fatal(err)
	}

	result2 := buffer2.String()
	if result2 != "<svg></svg>" {
		t.Fatalf("Expected <svg></svg>, got %s", result2)
	}

	if client.calls > 1 {
		t.Fatalf("Expected a single client fetch call, got %d", client.calls)
	}
}
