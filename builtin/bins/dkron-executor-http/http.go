package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/armon/circbuf"
	"github.com/victorcoder/dkron/dkron"
)

const (
	// timeout seconds
	timeout = 30
	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 256000
)

// HTTP process http request
type HTTP struct {
}

// Execute Process method of the plugin
// "executor": "http",
// "executor_config": {
//     "method": "GET",             // Request method in uppercase
//     "url": "http://example.com", // Request url
//     "headers": "[]"              // Json string, such as "[\"Content-Type: application/json\"]"
//     "body": "",                  // POST body
//     "timeout": "30",             // Request timeout, unit seconds
//     "expectCode": "200",         // Expect response code, such as 200,206
//     "expectBody": "",            // Expect response body, support regexp, such as /success/
//     "debug": "true"              // Debug option, will log everything when this option is not empty
// }
func (s *HTTP) Execute(args *dkron.ExecuteRequest) (*dkron.ExecuteResponse, error) {
	out, err := s.ExecuteImpl(args)
	resp := &dkron.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// ExecuteImpl do http request
func (s *HTTP) ExecuteImpl(args *dkron.ExecuteRequest) ([]byte, error) {
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

	_timeout := timeout
	if args.Config["timeout"] != "" {
		_timeout, _ = strconv.Atoi(args.Config["timeout"])
	}

	client := &http.Client{Timeout: time.Duration(_timeout) * time.Second}
	req, err := http.NewRequest(args.Config["method"], args.Config["url"], bytes.NewBuffer([]byte(args.Config["body"])))
	if err != nil {
		return output.Bytes(), err
	}

	var headers []string
	if args.Config["headers"] != "" {
		if err := json.Unmarshal([]byte(args.Config["headers"]), &headers); err != nil {
			output.Write([]byte("Error: headers parse fail\n"))
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

	resp, err := client.Do(req)
	if err != nil {
		return output.Bytes(), err
	}

	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
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
		return output.Bytes(), errors.New("Not reach the expected code")
	}

	// match response
	if args.Config["expectBody"] != "" {
		if m, _ := regexp.MatchString(args.Config["expectBody"], string(out)); !m {
			return output.Bytes(), errors.New("Not match the expected body")
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
