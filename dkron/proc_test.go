package dkron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildCmd(t *testing.T) {
	testJob1 := &Job{
		Command: "echo 'test1' && echo 'success'",
		Shell:   true,
	}

	cmd := buildCmd(testJob1)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "test1\nsuccess\n", string(out))

	testJob2 := &Job{
		Command: "date && echo 'success'",
		Shell:   false,
	}
	cmd = buildCmd(testJob2)
	out, err = cmd.CombinedOutput()
	assert.Error(t, err)
}
