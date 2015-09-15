package boltdb

import (
	"os"
	"testing"

	"github.com/docker/libkv/store"
	"github.com/docker/libkv/testutils"
)

func makeBoltDBClient(t *testing.T) store.Store {
	kv, err := New([]string{"/tmp/not_exist_dir/__boltdbtest"}, &store.Config{Bucket: "boltDBTest"})

	if err != nil {
		t.Fatalf("cannot create store: %v", err)
	}

	return kv
}

func TestBoldDBStore(t *testing.T) {
	kv := makeBoltDBClient(t)

	testutils.RunTestCommon(t, kv)
	testutils.RunTestAtomic(t, kv)

	_ = os.Remove("/tmp/not_exist_dir/__boltdbtest")
}
