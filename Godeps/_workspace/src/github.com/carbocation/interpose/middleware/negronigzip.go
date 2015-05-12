package middleware

import (
	"net/http"

	"github.com/carbocation/interpose/adaptors"
	"github.com/phyber/negroni-gzip/gzip"
)

func NegroniGzip(compression int) func(http.Handler) http.Handler {
	return adaptors.FromNegroni(gzip.Gzip(compression))
}
