package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
)

func main() {
	middle := interpose.New()

	// Tell the browser which server this came from.
	// This modifies headers, so we want this to be called before
	// any middleware which might modify the body (in HTTP, the headers cannot be
	// modified after the body is modified)
	middle.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("X-Server-Name", "Interpose Test Server")
			next.ServeHTTP(rw, req)
		})
	})

	// Apply the router. By adding it last, all of our other middleware will be
	// executed before the router, allowing us to modify headers before any
	// output has been generated.
	router := mux.NewRouter()
	middle.UseHandler(router)

	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})

	// Launch and permit graceful shutdown, allowing up to 10 seconds for existing
	// connections to end
	graceful.Run(":3001", 10*time.Second, middle)
}
