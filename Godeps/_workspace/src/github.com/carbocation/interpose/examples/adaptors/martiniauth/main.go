package main

import (
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/adaptors"
	"github.com/gorilla/mux"
	"github.com/martini-contrib/auth"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})

	middle := interpose.New()

	// Use Martini-Auth, a Martini library for basic HTTP authentication
	middle.Use(adaptors.FromMartini(auth.Basic("user", "basic")))

	// Finally, add the router
	middle.UseHandler(router)

	// Now visit http://localhost:3001/guest and enter username "user"
	// and password "basic"
	http.ListenAndServe(":3001", middle)
}
