package adaptors

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

func FromNegroni(handler negroni.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		n := negroni.New()
		n.Use(handler)
		n.UseHandler(next)
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			n.ServeHTTP(rw, req)
		})
	}
}

func HandlerFromNegroni(handler negroni.Handler) http.Handler {
	n := negroni.New()
	n.Use(handler)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		n.ServeHTTP(rw, req)
	})
}
