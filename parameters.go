package templicons

import (
	"net/url"
	"strconv"
)

type Parameters struct {
	Width  int
	Height int
	Color  string
}

func (p *Parameters) SetDimensions(width, height int) *Parameters {
	p.Width = width
	p.Height = height

	return p
}

func (p *Parameters) SetColor(color string) *Parameters {
	p.Color = color
	return p
}

func (ip *Parameters) asQueryString() string  {
	q := url.Values{}

	if ip.Width > 0 {
		q.Add("width", strconv.Itoa(ip.Width))
	}

	if ip.Height > 0 {
		q.Add("height", strconv.Itoa(ip.Height))
	}

	if ip.Color != "" {
		q.Add("color", ip.Color)
	}

	return q.Encode()
}
