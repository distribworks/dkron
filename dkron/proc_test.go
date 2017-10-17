package dkron

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildCmd(t *testing.T) {

	// test shell command
	testJob1 := &Job{
		Command: "echo 'test1' && echo 'success'",
		Shell:   true,
	}

	cmd := buildCmd(testJob1)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "test1\nsuccess\n", string(out))

	// test not shell command
	testJob2 := &Job{
		Command: "date && echo 'success'",
		Shell:   false,
	}
	cmd = buildCmd(testJob2)
	out, err = cmd.CombinedOutput()
	assert.Error(t, err)
}

func Test_buildCmdWithCustomEnvironmentVariables(t *testing.T) {
	defer func() {
		os.Setenv("Foo", "")
	}()
	os.Setenv("Foo", "Bar")
	testJob := &Job{
		Command:              "echo $Foo && echo $He",
		EnvironmentVariables: []string{"Foo=Toto", "He=Ho"},
		Shell:                true,
	}

	cmd := buildCmd(testJob)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "Toto\nHo\n", string(out))

}
