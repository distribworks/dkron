package interpose

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/interpose/middleware"
)

func BasicMiddleware() *Middleware {
	middle := New()

	middle.Use(middleware.Json())
	middle.Use(middleware.Buffer())

	middle.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			fmt.Fprint(rw, "0")
			next.ServeHTTP(rw, req)
			fmt.Fprint(rw, "0")
		})
	})

	middle.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			fmt.Fprint(rw, "1")
			next.ServeHTTP(rw, req)
			fmt.Fprint(rw, "1")
		})
	})

	middle.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "2")
	}))

	return middle
}

func TestCompiledMiddleware(t *testing.T) {
	response := httptest.NewRecorder()

	middle := BasicMiddleware().Handler()

	middle.ServeHTTP(response, (*http.Request)(nil))
	out, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	expect(t, string(out), "01210")
}

func TestServeHTTP(t *testing.T) {
	response := httptest.NewRecorder()

	middle := BasicMiddleware()

	middle.ServeHTTP(response, (*http.Request)(nil))
	out, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	expect(t, string(out), "01210")
}

func TestEmptyMiddleware(t *testing.T) {
	response := httptest.NewRecorder()

	middle := New()

	middle.ServeHTTP(response, (*http.Request)(nil))
	out, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	expect(t, string(out), "")
}

func BenchmarkCompiled(b *testing.B) {
	response := httptest.NewRecorder()

	middle := BasicMiddleware().Handler()

	for i := 0; i < b.N; i++ {
		middle.ServeHTTP(response, (*http.Request)(nil))
	}
}

func BenchmarkUncompiled(b *testing.B) {
	response := httptest.NewRecorder()

	middle := BasicMiddleware()

	for i := 0; i < b.N; i++ {
		middle.ServeHTTP(response, (*http.Request)(nil))
	}
}

func BenchmarkEmpty(b *testing.B) {
	response := httptest.NewRecorder()

	middle := New()
	middle.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		return
	}))

	for i := 0; i < b.N; i++ {
		middle.ServeHTTP(response, (*http.Request)(nil))
	}
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
