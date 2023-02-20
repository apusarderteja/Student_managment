package handler

import (
	"fmt"
	"net/http"
)

func (h Handler) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>hello</h1>")
}
