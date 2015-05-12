package main

import (
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type key int

const (
	CountKey key = iota
)

func main() {
	middle := interpose.New()

	// Create a middleware that yields a global counter that increments until
	// the server is shut down. Note that this would actually require a mutex
	// or a channel to be safe for concurrent use. Therefore, this example is
	// unsafe.
	middle.Use(context.ClearHandler)
	middle.Use(func() func(http.Handler) http.Handler {
		c := 0

		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				c++
				context.Set(req, CountKey, c)
				next.ServeHTTP(w, req)
			})
		}
	}())

	// Apply the router.
	router := mux.NewRouter()
	router.HandleFunc("/test/{user}", func(w http.ResponseWriter, req *http.Request) {
		c, ok := context.GetOk(req, CountKey)
		if !ok {
			fmt.Println("Context not ok")
		}

		fmt.Fprintf(w, "Hi %s, this is visit #%d to the site since the server was last rebooted.", mux.Vars(req)["user"], c)

	})
	middle.UseHandler(router)

	http.ListenAndServe(":3001", middle)
}
