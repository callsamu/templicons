package templicons

import "errors"

type MockClient struct {}

func (c MockClient) Fetch(url string) ([]byte, error) {
	if url != "" {
		response := []byte("<svg></svg>")
		return response, nil
	}

	err := errors.New("no url provided")
	return nil, err
}
