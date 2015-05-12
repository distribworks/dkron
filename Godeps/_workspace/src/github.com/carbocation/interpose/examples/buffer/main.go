package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/middleware"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
)

func main() {
	mw := interpose.New()

	// Turn on output buffering. It's added first because it encapsulates all other output.
	// This enables headers to be written after data is sent, because they're all stored
	// in a buffer, so the user's browser will see everything in the order it expects.
	mw.Use(middleware.Buffer())

	// Tell the browser our output will be JSON. Note that because this is added
	// before the router, we will write JSON headers AFTER the router starts
	// writing output to the browser! Usually this would mean that the header would
	// not get set correctly. However, because we are buffering (see below),
	// nothing will get sent to the browser until all output is finished rendering.
	// This means that the header will get set correctly.
	mw.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	})

	// Apply the router.
	router := mux.NewRouter()
	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})
	mw.UseHandler(router)

	// Launch and permit graceful shutdown, allowing up to 10 seconds for existing
	// connections to end
	graceful.Run(":3001", 10*time.Second, mw)
}
