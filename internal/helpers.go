package helpers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/pmwals09/yobs/internal/models/contact"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
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

type FormData struct {
	Errors map[string][]string
	Values map[string]string
}

func (fd *FormData) AddError(key, val string) {
	if fd.Errors == nil {
		fd.Errors = make(map[string][]string)
	}
	if fd.Errors[key] == nil {
		fd.Errors[key] = make([]string, 0)
	}
	fd.Errors[key] = append(fd.Errors[key], val)
}

func (fd *FormData) AddValue(key, val string) {
	if fd.Values == nil {
		fd.Values = make(map[string]string)
	}
	fd.Values[key] = val
}

type ProfileArgs struct {
	Username string
	Email    string
}

type OpptyDetails struct {
	Oppty     opportunity.Opportunity
	Documents []document.Document
	Contacts  []contact.Contact
}

func logError(msg string, logger slog.Logger) {
	logger.Error(msg)

}
