package middleware

import (
	"net/http"

	"github.com/goods/httpbuf"
)

/*
Middleware that buffers all http output. This permits
output to be written before headers are sent. Downside:
no output is sent until it's all ready to be sent, so
this breaks streaming.

Note: currently ignores errors
*/
func Buffer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bw := new(httpbuf.Buffer)
			next.ServeHTTP(bw, r)
			bw.Apply(w)
		})
	}
}
