package main

import (
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/middleware"
	"github.com/gorilla/mux"
)

func main() {
	middle := interpose.New()

	// First apply any middleware that will not write output to http body
	// but which may (or may not) modify headers
	middle.Use(middleware.GorillaLog())

	// Now apply any middleware that might modify the http body. This permits the
	// preceding middleware to alter headers
	router := mux.NewRouter()
	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})
	middle.UseHandler(router)

	http.ListenAndServe(":3001", middle)
}
