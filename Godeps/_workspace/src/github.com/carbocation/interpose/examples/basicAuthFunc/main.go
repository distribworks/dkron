package main

import (
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/middleware"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type ctxKey int

const (
	authKey ctxKey = iota
)

func main() {
	middle := interpose.New()

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, `<h1>Welcome to the public page!</h1><p><a href="/protected/">Rabbit hole</a></p>`)
	})

	middle.UseHandler(router)

	// Now we will define a sub-router that uses the BasicAuth middleware
	// When you call any url starting with the path /protected, you will need to authenticate
	protectedRouter := mux.NewRouter().Methods("GET").PathPrefix("/protected").Subrouter()
	protectedRouter.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the protected page, %s!", context.Get(req, authKey))
	})

	protectedMiddlew := interpose.New()
	protectedMiddlew.Use(middleware.BasicAuthFunc(func(user, pass string, req *http.Request) bool {
		if middleware.SecureCompare(user, "admin") && middleware.SecureCompare(pass, "guessme") {
			context.Set(req, authKey, user)
			return true
		} else {
			return false
		}
	}))
	protectedMiddlew.Use(context.ClearHandler)
	protectedMiddlew.UseHandler(protectedRouter)

	router.Methods("GET").PathPrefix("/protected").Handler(protectedMiddlew)

	http.ListenAndServe(":3001", middle)
}
