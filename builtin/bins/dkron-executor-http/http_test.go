package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/distribworks/dkron/v3/plugin/types"
	"github.com/stretchr/testify/assert"
)

func newTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Path: %s", r.URL.Path)
		switch r.URL.Path {

		// Handlers to control known HTTP response codes
		case "/200":
			w.WriteHeader(200)
			return
		case "/400":
			w.WriteHeader(400)
			return
		case "/401":
			w.WriteHeader(401)
			return
		case "/404":
			w.WriteHeader(404)
			return
		case "/500":
			w.WriteHeader(500)
			return

		// Return a predefined string as the response body
		case "/hello":
			w.Write([]byte("hello"))
			return

		// Echo POST body back to request
		case "/echo":
			if r.Method == http.MethodPost {
				in, err := ioutil.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				r.Body.Close()
				w.Write(in)
				w.WriteHeader(200)
				return
			}
		}

	}))
}

func TestExecute(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	tests := []struct {
		name    string
		config  map[string]string
		want    []byte
		wantErr bool
	}{
		{"Expected 200", map[string]string{"method": "GET", "url": fmt.Sprintf("%s/200", ts.URL), "expectCode": "200"}, []byte{}, false},
		{"Expected 400", map[string]string{"method": "GET", "url": fmt.Sprintf("%s/400", ts.URL), "expectCode": "400"}, []byte{}, false},
		{"Expected 404", map[string]string{"method": "GET", "url": fmt.Sprintf("%s/404", ts.URL), "expectCode": "404"}, []byte{}, false},
		{"Unexpected 400 is error", map[string]string{"method": "GET", "url": fmt.Sprintf("%s/400", ts.URL), "expectCode": "200"}, []byte{}, true},
		{"Empty URL is error", map[string]string{"method": "GET", "url": "", "expectCode": "200"}, []byte{}, true},
		{"Empty method is error", map[string]string{"method": "", "url": fmt.Sprintf("%s/200", ts.URL), "expectCode": "200"}, []byte{}, true},
		{"Expected GET Response", map[string]string{"method": "GET", "url": fmt.Sprintf("%s/hello", ts.URL), "expectCode": "200"}, []byte("hello"), false},
		{"Expected POST Response", map[string]string{"method": "POST", "url": fmt.Sprintf("%s/echo", ts.URL), "body": "this is a post body", "expectBody": "this is a post body", "expectCode": "200"}, []byte("this is a post body"), false},
		{"Unexpected POST Response is error", map[string]string{"method": "POST", "url": fmt.Sprintf("%s/echo", ts.URL), "body": "this is a post body", "expectBody": "not this", "expectCode": "200"}, []byte("this is a post body"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			http := &HTTP{}
			pa := &types.ExecuteRequest{
				JobName: tt.name,
				Config:  tt.config,
			}
			// Err is always nil from http.Execute()
			got, _ := http.Execute(pa, nil)
			if (got.Error != "") != tt.wantErr {
				t.Errorf("HTTP.Execute().Error = %v, wantErr %v", got.Error, tt.wantErr)
				return
			}

			if !bytes.Equal(got.Output, tt.want) {
				t.Errorf("HTTP.Execute().Output = %s, want %s", got.Output, tt.want)
			}

		})
	}
}

// Note: badssl.com was meant for _manual_ testing. Maybe these tests should be disabled by default.
func TestNoVerifyPeer(t *testing.T) {
	pa := &types.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"method":          "GET",
			"url":             "https://self-signed.badssl.com/",
			"expectCode":      "200",
			"debug":           "true",
			"tlsNoVerifyPeer": "true",
		},
	}
	http := &HTTP{}
	output, _ := http.Execute(pa, nil)
	fmt.Println(string(output.Output))
	fmt.Println(output.Error)
	assert.Equal(t, "", output.Error)
}

func TestClientSSLCert(t *testing.T) {
	// client certs: https://badssl.com/download/
	pa := &types.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"method":                "GET",
			"url":                   "https://client.badssl.com/",
			"expectCode":            "200",
			"debug":                 "true",
			"tlsCertificateFile":    "testdata/badssl.com-client.pem",
			"tlsCertificateKeyFile": "testdata/badssl.com-client-key-decrypted.pem",
		},
	}
	http := &HTTP{}
	output, _ := http.Execute(pa, nil)
	fmt.Println(string(output.Output))
	fmt.Println(output.Error)
	assert.Equal(t, "", output.Error)
}

func TestRootCA(t *testing.T) {
	// untrusted root ca cert: https://badssl.com/certs/ca-untrusted-root.crt
	pa := &types.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"method":         "GET",
			"url":            "https://untrusted-root.badssl.com/",
			"expectCode":     "200",
			"debug":          "true",
			"tlsRootCAsFile": "testdata/badssl-ca-untrusted-root.crt",
		},
	}
	http := &HTTP{}
	output, _ := http.Execute(pa, nil)
	fmt.Println(string(output.Output))
	fmt.Println(output.Error)
	assert.Equal(t, "", output.Error)
}
