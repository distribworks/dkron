package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
	"github.com/unrolled/secure"
)

func main() {
	mw := interpose.New()

	// Use unrolled's secure framework
	// If you inspect the headers, you will see X-Frame-Options set to DENY
	// Must be called before the router because it modifies HTTP headers
	secureMiddleware := secure.New(secure.Options{
		FrameDeny: true,
	})
	mw.Use(secureMiddleware.Handler)

	// Apply the router. By adding it first, all of our other middleware will be
	// executed before the router, allowing us to modify headers before any
	// output has been generated.
	router := mux.NewRouter()
	mw.UseHandler(router)

	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})

	// Launch and permit graceful shutdown, allowing up to 10 seconds for existing
	// connections to end
	graceful.Run(":3001", 10*time.Second, mw)
}
