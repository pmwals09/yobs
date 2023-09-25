package helpers

import (
	"fmt"
	"net/http"
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
