package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildCmd(t *testing.T) {

	// test shell command
	cmd, err := buildCmd("echo 'test1' && echo 'success'", true, []string{}, "")
	assert.NoError(t, err)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "test1\nsuccess\n", string(out))

	// test not shell command
	cmd, err = buildCmd("date && echo 'success'", false, []string{}, "")
	assert.NoError(t, err)
	assert.Len(t, cmd.Args, 1)
}

func Test_buildCmdWithCustomEnvironmentVariables(t *testing.T) {
	defer func() {
		os.Setenv("Foo", "")
	}()
	os.Setenv("Foo", "Bar")

	cmd, err := buildCmd("echo $Foo && echo $He", true, []string{"Foo=Toto", "He=Ho"}, "")
	assert.NoError(t, err)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "Toto\nHo\n", string(out))

}
