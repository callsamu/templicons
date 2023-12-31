package templicons

import (
	"context"
	"io"
	"sync"

	"github.com/a-h/templ"
)


type Cache struct {
	Errors   chan error

	cache   map[string][]byte
	mutex   sync.Mutex
	wg      sync.WaitGroup

	client   Client
	Fallback Fallback
}

func NewCache(client Client, fallback Fallback) *Cache {
	return &Cache{
		Errors:   make(chan error),
		cache:    make(map[string][]byte),
		client:   client,
		Fallback: fallback,
	}
}

var cache *Cache

func init() {
	client := NewIconifyClient("https://api.iconify.design")
	cache = NewCache(client, DefaultFallback)
}

func Errors() chan error {
	return cache.Errors
}

func SetInstances(urls ...string) {
	cache.client.SetInstances(urls...)
}

func SetFallback(fallback Fallback) {
	cache.Fallback = fallback
}

func Icon(name string, p *Parameters) templ.Component {
	return cache.Icon(name, p)
}

func (c *Cache) Icon(name string, p *Parameters) templ.Component {
	set, icon := parseName(name)
	url := iconPath(set, icon, p)

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		c.mutex.Lock()
		if cached, ok := c.cache[url]; ok {
			_, err := w.Write(cached)

			c.mutex.Unlock()
			return err
		}
		c.mutex.Unlock()

		svg, err := c.fetchAndSave(url)
		if err != nil {
			return err
		}

		_, err = w.Write(svg) 
		return err
	})
}

func IconWithFallback(name string, fallback string, p *Parameters) templ.Component {
	return cache.IconWithFallback(name, fallback, p)
}

func (c *Cache) IconWithFallback(name string, fallback string, p *Parameters) templ.Component {
	set, icon := parseName(name)
	url := iconPath(set, icon, p)

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if cached, ok := c.cache[url]; ok {
			_, err := w.Write(cached)
			return err
		}

		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			c.fetchAndSave(url)
		}()

		html := "<span>" + fallback + "</span>"
		_, err := w.Write([]byte(html))
		return err
	})
}

func (c *Cache) fetchAndSave(url string) ([]byte, error) {
	svg, err := c.client.Fetch(url)
	if err != nil {
		return nil, err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[url] = svg
	return svg, nil
}
