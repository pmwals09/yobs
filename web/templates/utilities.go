package templates

import (
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

func formatApplicationDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func insertIDIntoHref(href string, id uint) templ.SafeURL {
  return templ.SafeURL(insertIDIntoString(href, id))
}

func insertIDIntoString(href string, id uint) string {
  stringifiedId := strconv.FormatUint(uint64(id), 10)
  return strings.Replace(href, "{}", stringifiedId, 1)
}

func UnsafeRawHtml(html string) templ.Component {
  return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
    _, err := io.WriteString(w, html)
    return err
  })
}
