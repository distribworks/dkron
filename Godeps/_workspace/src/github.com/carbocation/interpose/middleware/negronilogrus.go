package middleware

import (
	"net/http"

	"github.com/carbocation/interpose/adaptors"
	"github.com/meatballhat/negroni-logrus"
)

func NegroniLogrus() func(http.Handler) http.Handler {
	return adaptors.FromNegroni(negronilogrus.NewMiddleware())
}
