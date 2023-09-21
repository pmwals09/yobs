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
