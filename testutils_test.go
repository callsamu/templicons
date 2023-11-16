package templicons

import "errors"

type MockClient struct {
	calls int
}

func (c MockClient) SetInstances(instances ...string) {}

func (c MockClient) Fetch(url string) ([]byte, error) {
	c.calls++

	if url != "" {
		response := []byte("<svg></svg>")
		return response, nil
	}

	err := errors.New("no url provided")
	return nil, err
}
