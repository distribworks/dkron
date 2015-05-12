package adaptors

import (
	"net/http"

	"github.com/go-martini/martini"
)

func FromMartini(handler martini.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		m := martini.New()
		m.Use(handler)
		m.Use(next)
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			m.ServeHTTP(rw, req)
		})
	}
}

func HandlerFromMartini(handler martini.Handler) http.Handler {
	m := martini.New()
	m.Use(handler)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		m.ServeHTTP(rw, req)
	})
}
