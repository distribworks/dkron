package cmd

import (
	"os"
	"testing"
)

var (
	logLevel = "error"
	etcdAddr = getEnvWithDefault()
)

func getEnvWithDefault() string {
	ea := os.Getenv("DKRON_BACKEND_MACHINE")
	if ea == "" {
		return "127.0.0.1:2379"
	}
	return ea
}

func Test_unmarshalTags(t *testing.T) {
	tagPairs := []string{
		"tag1=val1",
		"tag2=val2",
	}

	tags, err := unmarshalTags(tagPairs)

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if v, ok := tags["tag1"]; !ok || v != "val1" {
		t.Fatalf("bad: %v", tags)
	}
	if v, ok := tags["tag2"]; !ok || v != "val2" {
		t.Fatalf("bad: %v", tags)
	}
}
