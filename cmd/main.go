package main

import (
	"ProjectGo/internal/transport"
	"log"
	"net/http"
)

func main() {
	transport.Setuprouter()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
