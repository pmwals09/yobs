package helpers

import (
	"fmt"
  "html/template"
	"net/http"
	"time"
)

func WriteError(w http.ResponseWriter, err error) {
	w.Write([]byte(fmt.Sprintf(
		"<p>An error has occurred: %s</p>",
		err.Error(),
	)))
}

func HTMXRedirect(w http.ResponseWriter, route string, code int) {
	w.Header().Add("HX-Redirect", route)
	w.WriteHeader(code)
	w.Write([]byte(""))
}

func GetListFuncMap() template.FuncMap {
	return template.FuncMap{
		"FormatApplicationDate": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02")
		},
	}
}
