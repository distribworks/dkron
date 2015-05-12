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

	// First call middleware that may manipulate HTTP headers, since
	// they must be called before the HTTP body is manipulated

	// Using Gorilla framework's combined logger
	middle.Use(middleware.GorillaLog())

	//Using Negroni's Gzip functionality
	middle.Use(middleware.NegroniGzip(gzip.DefaultCompression))

	// Now call middleware that can manipulate the HTTP body

	// Define the router. Note that we can add the router to our
	// middleware stack before we define the routes, if we want.
	router := mux.NewRouter()
	middle.UseHandler(router)

	// Configure our router
	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})

	// Define middleware that will apply to only some routes
	greenMiddle := interpose.New()

	// Tell the main router to send /green requests to our subrouter.
	// Again, we can do this before defining the full middleware stack.
	router.Methods("GET").PathPrefix("/green").Handler(greenMiddle)

	// Within the secondary middleware, just like above, we want to call anything that
	// will modify the HTTP headers before anything that will modify the body
	greenMiddle.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("X-Favorite-Color", "green")
	}))

	// Finally, define a sub-router based on our love of the color green
	// When you call any url such as http://localhost:3001/green/man , you will
	// also see an HTTP header sent called X-Favorite-Color with value "green"
	greenRouter := mux.NewRouter().Methods("GET").PathPrefix("/green").Subrouter() //Headers("Accept", "application/json")
	greenMiddle.UseHandler(greenRouter)

	greenRouter.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, green %s!", mux.Vars(req)["user"])
	})

	http.ListenAndServe(":3001", middle)
}
