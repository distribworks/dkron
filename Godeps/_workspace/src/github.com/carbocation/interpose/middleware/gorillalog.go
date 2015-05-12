package middleware

import (
	"net/http"
	"os"

	"github.com/carbocation/handlers"
)

/*
Wraps the Gorilla Logger
*/
func GorillaLog() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.CombinedLoggingHandler(os.Stdout, next).ServeHTTP(w, r)
		})
	}
}
