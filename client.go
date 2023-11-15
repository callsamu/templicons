package templicons

import (
	"io"
	"net/http"
)

type Client interface {
	Fetch(url string) ([]byte, error)
}

type IconifyClient struct {
	*http.Client
}

func NewIconifyClient() *IconifyClient {
	return &IconifyClient{&http.Client{}}
}

func (c IconifyClient) Fetch(url string) ([]byte, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	return body, err
}
