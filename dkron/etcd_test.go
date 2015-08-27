package dkron

import (
	"testing"
)

func TestEtcdClient(t *testing.T) {
	etcd := NewEtcdClient([]string{}, nil, "dkron-test")

	testJob := &Job{
		Name:     "test",
		Schedule: "@every 2s",
		Disabled: true,
	}

	if err := etcd.SetJob(testJob); err != nil {
		t.Fatalf("error creating job: %s", err)
	}

	if err := etcd.DeleteJob("test"); err != nil {
		t.Fatalf("error deleting job: %s", err)
	}

	if err := etcd.DeleteJob("test"); err == nil {
		t.Fatalf("error job deletion should fail: %s", err)
	}
}
