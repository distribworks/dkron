package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"
)

// User is the authenticated username that was extracted from the request.
type User string

// BasicRealm is used when setting the WWW-Authenticate response header.
var BasicRealm = "Authorization Required"

// Basic returns a Handler that authenticates via Basic Auth. Writes a http.StatusUnauthorized
// if authentication fails
func BasicAuth(username string, password string) func(http.Handler) http.Handler {
	var siteAuth = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			auth := req.Header.Get("Authorization")
			if !SecureCompare(auth, "Basic "+siteAuth) {
				unauthorized(res)
				return
			}
			next.ServeHTTP(res, req)
		})
	}
}

// BasicAuthFunc returns a Handler that authenticates via Basic Auth using the provided function.
// The function should return true for a valid username/password combination.
func BasicAuthFunc(authfn func(string, string, *http.Request) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			auth := req.Header.Get("Authorization")
			if len(auth) < 6 || auth[:6] != "Basic " {
				unauthorized(res)
				return
			}
			b, err := base64.StdEncoding.DecodeString(auth[6:])
			if err != nil {
				unauthorized(res)
				return
			}
			tokens := strings.SplitN(string(b), ":", 2)
			if len(tokens) != 2 || !authfn(tokens[0], tokens[1], req) {
				unauthorized(res)
				return
			}
			next.ServeHTTP(res, req)
		})
	}
}

// SecureCompare performs a constant time compare of two strings to limit timing attacks.
func SecureCompare(given string, actual string) bool {
	givenSha := sha256.Sum256([]byte(given))
	actualSha := sha256.Sum256([]byte(actual))

	return subtle.ConstantTimeCompare(givenSha[:], actualSha[:]) == 1
}

func unauthorized(res http.ResponseWriter) {
	res.Header().Set("WWW-Authenticate", "Basic realm=\""+BasicRealm+"\"")
	http.Error(res, "Not Authorized", http.StatusUnauthorized)
}
