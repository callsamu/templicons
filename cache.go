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

	API      string
	client   Client
}

func NewCache(api string, client Client) *Cache {
	return &Cache{
		Errors: make(chan error),
		API:    api,
		cache:  make(map[string][]byte),
		client: client,
	}
}

var cache *Cache

func init() {
	client := NewIconifyClient()
	cache = NewCache("https://api.iconify.design", client)
}

func Icon(name string) templ.Component {
	return cache.Icon(name, nil)
}

func (c *Cache) Icon(name string, p *Parameters) templ.Component {
	set, icon := parseName(name)
	url := iconURL(c.API, set, icon, p)

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
	return cache.IconWithFallback(name, fallback, nil)
}

func (c *Cache) IconWithFallback(name string, fallback string, p *Parameters) templ.Component {
	set, icon := parseName(name)
	url := iconURL(c.API, set, icon, p)

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
