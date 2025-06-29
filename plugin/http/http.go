package http

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/armon/circbuf"
	types "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	dkplugin "github.com/distribworks/dkron/v4/plugin"
	lru "github.com/hashicorp/golang-lru"
)

const (
	// timeout seconds
	timeout = 30
	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 256000
	// maxClientPoolSize limits the number of HTTP clients cached
	// This prevents unbounded memory growth in the client pool
	maxClientPoolSize = 100
)

// HTTP process http request
type HTTP struct {
	clientPool *lru.Cache
	mu         sync.RWMutex
}

// New
func New() *HTTP {
	cache, err := lru.NewWithEvict(maxClientPoolSize, func(key interface{}, value interface{}) {
		// Optionally close idle connections when evicting clients
		if client, ok := value.(http.Client); ok {
			if transport, ok := client.Transport.(*http.Transport); ok {
				transport.CloseIdleConnections()
			}
		}
	})
	if err != nil {
		// Fallback to a smaller cache if creation fails
		cache, _ = lru.New(10)
	}
	return &HTTP{
		clientPool: cache,
	}
}

// Execute Process method of the plugin
// "executor": "http",
//
//	"executor_config": {
//	    "method": "GET",             // Request method in uppercase
//	    "url": "http://example.com", // Request url
//	    "headers": "[]"              // Json string, such as "[\"Content-Type: application/json\"]"
//	    "body": "",                  // POST body
//	    "timeout": "30",             // Request timeout, unit seconds
//	    "expectCode": "200",         // Expect response code, such as 200,206
//	    "expectBody": "",            // Expect response body, support regexp, such as /success/
//	    "debug": "true"              // Debug option, will log everything when this option is not empty
//	}
func (s *HTTP) Execute(args *types.ExecuteRequest, cb dkplugin.StatusHelper) (*types.ExecuteResponse, error) {
	out, err := s.ExecuteImpl(args)
	resp := &types.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// ExecuteImpl do http request
func (s *HTTP) ExecuteImpl(args *types.ExecuteRequest) ([]byte, error) {
	output, _ := circbuf.NewBuffer(maxBufSize)
	var debug bool
	if args.Config["debug"] != "" {
		debug = true
		log.Printf("config  %#v\n\n", args.Config)
		output.Write([]byte(fmt.Sprintf("Config: %#v\n", args.Config)))
	}

	if args.Config["url"] == "" {
		return output.Bytes(), errors.New("url is empty")
	}

	if args.Config["method"] == "" {
		return output.Bytes(), errors.New("method is empty")
	}

	req, err := http.NewRequest(args.Config["method"], args.Config["url"], bytes.NewBuffer([]byte(args.Config["body"])))
	if err != nil {
		return output.Bytes(), err
	}
	req.Close = true

	var headers []string
	if args.Config["headers"] != "" {
		if err := json.Unmarshal([]byte(args.Config["headers"]), &headers); err != nil {
			output.Write([]byte("Error: parsing headers failed\n"))
		}
	}

	for _, h := range headers {
		if h != "" {
			kv := strings.Split(h, ":")
			req.Header.Set(kv[0], strings.TrimSpace(kv[1]))
		}
	}
	if debug {
		log.Printf("request  %#v\n\n", req)
	}

	// get client from pool
	var (
		client http.Client
		ok     bool
	)

	clientKey := s.generateClientKey(args.Config)

	s.mu.RLock()
	if cachedClient, found := s.clientPool.Get(clientKey); found {
		client = cachedClient.(http.Client)
		ok = true
	}
	s.mu.RUnlock()

	if !ok {
		var warns []error
		client, warns = createClient(args.Config)
		for _, warn := range warns {
			_, _ = output.Write([]byte(fmt.Sprintf("Warning: %s.\n", warn.Error())))
		}
		s.mu.Lock()
		s.clientPool.Add(clientKey, client)
		s.mu.Unlock()
	}

	// do request
	resp, err := client.Do(req)
	if err != nil {
		return output.Bytes(), err
	}

	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return output.Bytes(), err
	}

	if debug {
		log.Printf("response  %#v\n\n", resp)
		log.Printf("response body  %#v\n\n", string(out))
	}

	// write the response to output
	_, err = output.Write(out)
	if err != nil {
		return output.Bytes(), err
	}

	// match response code
	if args.Config["expectCode"] != "" && !strings.Contains(args.Config["expectCode"]+",", fmt.Sprintf("%d,", resp.StatusCode)) {
		return output.Bytes(), errors.New("received response code does not match the expected code")
	}

	// match response
	if args.Config["expectBody"] != "" {
		if m, _ := regexp.MatchString(args.Config["expectBody"], string(out)); !m {
			return output.Bytes(), errors.New("received response body did not match the expected body")
		}
	}

	// Warn if buffer is overritten
	if output.TotalWritten() > output.Size() {
		log.Printf("'%s %s': generated %d bytes of output, truncated to %d",
			args.Config["method"], args.Config["url"],
			output.TotalWritten(), output.Size())
	}

	return output.Bytes(), nil
}

// generateClientKey creates a unique key for the client pool based on configuration
// This fixes the key collision issue by using proper hashing instead of string concatenation
func (s *HTTP) generateClientKey(config map[string]string) string {
	// Only include configuration that affects the HTTP client behavior
	relevantKeys := []string{
		"timeout",
		"tlsNoVerifyPeer",
		"tlsRootCAsFile",
		"tlsCertificateFile",
		"tlsCertificateKeyFile",
	}

	// Create a sorted map to ensure consistent key generation
	var keyParts []string
	for _, key := range relevantKeys {
		if value, exists := config[key]; exists && value != "" {
			keyParts = append(keyParts, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Sort to ensure consistent ordering
	sort.Strings(keyParts)

	// Create hash of the configuration
	hasher := sha256.New()
	hasher.Write([]byte(strings.Join(keyParts, "|")))
	return hex.EncodeToString(hasher.Sum(nil))
}

// createClient always returns a new http client. Any errors returned are
// errors in the configuration.
func createClient(config map[string]string) (http.Client, []error) {
	var errs []error

	_timeout, err := atoiOrDefault(config["timeout"], timeout)
	if config["timeout"] != "" && err != nil {
		errs = append(errs, fmt.Errorf("invalid timeout value: %s", err.Error()))
	}

	tlsconf := &tls.Config{}
	tlsconf.InsecureSkipVerify, err = strconv.ParseBool(config["tlsNoVerifyPeer"])
	if config["tlsNoVerifyPeer"] != "" && err != nil {
		errs = append(errs, fmt.Errorf("not disabling certificate validation: %s", err.Error()))
	}

	if config["tlsCertificateFile"] != "" {
		cert, err := tls.LoadX509KeyPair(config["tlsCertificateFile"], config["tlsCertificateKeyFile"])
		if err == nil {
			tlsconf.Certificates = append(tlsconf.Certificates, cert)
		} else {
			errs = append(errs, fmt.Errorf("not using client certificate: %s", err.Error()))
		}
	}

	if config["tlsRootCAsFile"] != "" {
		tlsconf.RootCAs, err = loadCertPool(config["tlsRootCAsFile"])
		if err != nil {
			errs = append(errs, fmt.Errorf("using system root CAs instead of configured CAs: %s", err.Error()))
		}
	}

	// Create transport with proper connection pooling settings
	transport := &http.Transport{
		TLSClientConfig: tlsconf,
		// Set reasonable connection pool limits to prevent resource leaks
		MaxIdleConns:          10,
		MaxIdleConnsPerHost:   5,
		MaxConnsPerHost:       20,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return http.Client{
		Transport: transport,
		Timeout:   time.Duration(_timeout) * time.Second,
	}, errs
}

// loadCertPool creates a CertPool using the given file
func loadCertPool(filename string) (*x509.CertPool, error) {
	certsFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(certsFile) {
		return nil, fmt.Errorf("no certificates in file")
	}

	return roots, nil
}

// atoiOrDefault returns the integer value of s, or a default value
// if s could not be converted, along with an error.
func atoiOrDefault(s string, _default int) (int, error) {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i, nil
	}
	return _default, fmt.Errorf("\"%s\" not understood (%s), using default value of %d", s, err, _default)
}
