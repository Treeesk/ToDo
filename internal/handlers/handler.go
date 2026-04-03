package handlers

import (
	"fmt"
	"net/http"
)

func HandlerStart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
