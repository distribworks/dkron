package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/distribworks/dkron/v4/types"
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
				in, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				r.Body.Close()
				_, _ = w.Write(in)
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
			http := New()
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
	http := New()
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
	http := New()
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
	http := New()
	output, _ := http.Execute(pa, nil)
	fmt.Println(string(output.Output))
	fmt.Println(output.Error)
	assert.Equal(t, "", output.Error)
}

// TestClientPoolMemoryLeak tests that the client pool doesn't grow unbounded
func TestClientPoolMemoryLeak(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	httpExecutor := New()
	
	// Get initial memory stats
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	// Simulate many different configurations to fill the client pool
	// Use different combinations of TLS configs to create unique client pool keys
	for i := 0; i < 1000; i++ {
		config := map[string]string{
			"method":              "GET",
			"url":                 fmt.Sprintf("%s/200", ts.URL),
			"expectCode":          "200",
			"timeout":             fmt.Sprintf("%d", 30+i%100),                 // Different timeout values
			"tlsRootCAsFile":      fmt.Sprintf("/tmp/fake%d.crt", i%100),        // Different fake TLS files
			"tlsCertificateFile":  fmt.Sprintf("/tmp/cert%d.pem", i%100),        // Different fake cert files
			"tlsCertificateKeyFile": fmt.Sprintf("/tmp/key%d.pem", i%100),       // Different fake key files
		}
		
		pa := &types.ExecuteRequest{
			JobName: fmt.Sprintf("testJob%d", i),
			Config:  config,
		}
		
		// We expect this to fail due to file not found, but the client should still be cached
		httpExecutor.ExecuteImpl(pa)
	}
	
	// Check memory usage after creating many clients
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	// Check that client pool size is reasonable
	poolSize := httpExecutor.clientPool.Len()
	t.Logf("Client pool size: %d", poolSize)
	assert.LessOrEqual(t, poolSize, maxClientPoolSize, "Client pool should be capped at max size")
	assert.Greater(t, poolSize, 10, "Client pool should have many entries for diverse configurations")
}

// TestClientPoolKeyCollisions tests for unintended client sharing due to poor key generation
func TestClientPoolKeyCollisions(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	httpExecutor := New()
	
	// These two configurations should create different clients but currently don't
	// due to poor key generation (simple string concatenation)
	config1 := map[string]string{
		"method":     "GET",
		"url":        fmt.Sprintf("%s/200", ts.URL),
		"expectCode": "200",
		"timeout":    "30", // timeout = "30"
		"tlsRootCAsFile": "file.crt", // tlsRootCAsFile = "file.crt"
	}
	
	config2 := map[string]string{
		"method":     "GET", 
		"url":        fmt.Sprintf("%s/200", ts.URL),
		"expectCode": "200",
		"timeout":    "3", // timeout = "3"
		"tlsRootCAsFile": "0file.crt", // tlsRootCAsFile = "0file.crt"
	}
	
	// The concatenated keys would be:
	// config1: "30" + "file.crt" + "" + "" = "30file.crt"
	// config2: "3" + "0file.crt" + "" + "" = "30file.crt"
	// These are identical! This is a collision bug.
	
	pa1 := &types.ExecuteRequest{JobName: "test1", Config: config1}
	pa2 := &types.ExecuteRequest{JobName: "test2", Config: config2}
	
	// Execute both
	httpExecutor.ExecuteImpl(pa1)
	httpExecutor.ExecuteImpl(pa2)
	
	// Due to the collision, only one client should be in the pool
	poolSize := httpExecutor.clientPool.Len()
	t.Logf("Client pool size after potential collision: %d", poolSize)
	
	// With the fix, different configs should create different clients
	assert.Equal(t, 2, poolSize, "Proper key generation should prevent unintended client sharing")
}
