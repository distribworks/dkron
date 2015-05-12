package middleware

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// Nosurf is a wrapper for justinas' csrf protection middleware
func Nosurf() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return nosurf.New(next)
	}
}
