package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/armon/circbuf"
	"github.com/victorcoder/dkron/dkron"
)

const (
	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 256000
)

type Http struct {
}

// Process method of the plugin
func (s *Http) Execute(args *dkron.ExecuteRequest) ([]byte, error) {
	output, _ := circbuf.NewBuffer(maxBufSize)

	if args.Config["url"] == "" {
		return nil, errors.New("url is empty")
	}

	if args.Config["method"] == "" {
		return nil, errors.New("method is empty")
	}

	body, err := base64.StdEncoding.DecodeString(args.Config["body"])
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(args.Config["method"], args.Config["url"], bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", args.Config["content_type"])

	log.Printf("%s %s", args.Config["method"], args.Config["url"])

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	_, err = output.Write(out)
	if err != nil {
		return nil, err
	}

	// Warn if buffer is overritten
	if output.TotalWritten() > output.Size() {
		log.Printf("'%s %s': generated %d bytes of output, truncated to %d",
			args.Config["method"], args.Config["url"],
			output.TotalWritten(), output.Size())
	}

	return output.Bytes(), nil
}
