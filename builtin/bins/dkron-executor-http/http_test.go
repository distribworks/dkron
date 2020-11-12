package main

import (
	"fmt"
	"testing"

	"github.com/distribworks/dkron/v3/plugin/types"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	pa := &types.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"method":     "GET",
			"url":        "https://httpbin.org/get",
			"expectCode": "200",
			"debug":      "true",
		},
	}
	http := &HTTP{}
	output, err := http.Execute(pa, nil)
	fmt.Println(string(output.Output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecutePost(t *testing.T) {
	pa := &types.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"method":     "POST",
			"url":        "https://httpbin.org/post",
			"body":       "{\"hello\":11}",
			"headers":    "[\"Content-Type:application/json\"]",
			"expectCode": "200",
			"expectBody": "",
			"debug":      "true",
		},
	}
	http := &HTTP{}
	output, err := http.Execute(pa, nil)
	fmt.Println(string(output.Output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
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
