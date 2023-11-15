package templicons

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/a-h/templ"
)


type Cache struct {
	Errors   chan error
	API      string

	cache   map[string][]byte
	mutex   sync.Mutex
	client  Client

	AllowFallback bool
}

func NewCache(api string, allowFallback bool, client Client) *Cache {
	return &Cache{
		Errors:        make(chan error),
		API:           api,
		cache:         make(map[string][]byte),
		client:        client,
		AllowFallback: allowFallback,
	}
}

var cache *Cache

func init() {
	client := NewIconifyClient()
	cache = NewCache("https://api.iconify.design", false, client)
}

func Icon(name string) templ.Component {
	return cache.Icon(name, nil)
}

func (c *Cache) Icon(name string, p *Parameters) templ.Component {
	return c.IconWithFallback(name, "", p)
}

func IconWithFallback(name string, fallback string) templ.Component {
	return cache.IconWithFallback(name, fallback, nil)
}

func (c *Cache) IconWithFallback(name string, fallback string, p *Parameters) templ.Component {
	parsed := strings.Split(name, ":")
	if len(parsed) != 2 {
		panic("invalid icon name")
	}

	set := parsed[0]
	icon := parsed[1]

	url := iconURL(c.API, set, icon, p)

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		c.mutex.Lock()
		if cached, ok := c.cache[url]; ok {
			_, err := w.Write(cached)

			c.mutex.Unlock()
			return err
		}
		c.mutex.Unlock()

		if !c.AllowFallback {
			svg, err := c.fetchAndSave(url)
			if err != nil {
				return err
			}

			_, err = w.Write(svg)
			return err
		}

		return nil
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
