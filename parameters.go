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

func (p *Parameters) SetDimensions(width, height int) {
	p.Width = width
	p.Height = height
}

func (p *Parameters) SetColor(color string) {
	p.Color = color
}

func (ip *Parameters) AsQueryString() string  {
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
