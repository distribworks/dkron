package main

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/middleware"
	"github.com/gorilla/mux"
)

func main() {
	middle := interpose.New()

	// First apply any middleware that will not write output to http body

	// Log to stdout. Taken from Gorilla
	middle.Use(middleware.GorillaLog())

	// Gzip output that follows. Taken from Negroni
	middle.Use(middleware.NegroniGzip(gzip.DefaultCompression))

	// Now apply any middleware that modify the http body.
	router := mux.NewRouter()
	middle.UseHandler(router)

	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})

	http.ListenAndServe(":3001", middle)
}
