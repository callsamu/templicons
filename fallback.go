package templicons

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

type Fallback func(span, iconURL string) templ.Component

func DefaultFallback(span, iconURL string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		html := templ.EscapeString("<span>" + span + "</span>")
		_, err := w.Write([]byte(html))
		return err
	})
}
