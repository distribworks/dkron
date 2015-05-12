package main

import (
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
)

func main() {
	middle := interpose.New()

	// Send a header telling the world this is coming from an Interpose Test Server
	// You could imagine setting Content-type application/json or other useful
	// headers in other circumstances.
	middle.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("X-Server-Name", "Interpose Test Server")
	}))

	// In this example, we use a basic router with a catchall route that
	// matches any path. In other examples in this project, we use a
	// more advanced router.
	// The last middleware added is the last executed, so by adding the router
	// last, our other middleware get a chance to modify HTTP headers before our
	// router writes to the HTTP Body
	router := http.NewServeMux()
	middle.UseHandler(router)

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to %s!", req.URL.Path)
	}))

	http.ListenAndServe(":3001", middle)
}
