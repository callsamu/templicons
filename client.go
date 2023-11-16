package templicons

import (
	"io"
	"net/http"
)

type Client interface {
	Fetch(path string) ([]byte, error)
	SetInstances(instances ...string)
}

type IconifyClient struct {
	*http.Client

	Instances []string
}

func NewIconifyClient(instances ...string) *IconifyClient {
	return &IconifyClient{
		Client: &http.Client{},
		Instances: instances,
	}
}

func (c *IconifyClient) SetInstances(instances ...string) {
	c.Instances = instances
}

func (c IconifyClient) Fetch(path string) ([]byte, error) {
	instance := c.Instances[0]

	resp, err := c.Get(instance + path)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	return body, err
}
