package cmd

import (
	"errors"
	"testing"

	typesv1 "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	dkplugin "github.com/distribworks/dkron/v4/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type schemaExecutor struct {
	schema string
	err    error
}

func (e *schemaExecutor) Execute(args *typesv1.ExecuteRequest, cb dkplugin.StatusHelper) (*typesv1.ExecuteResponse, error) {
	return &typesv1.ExecuteResponse{}, nil
}

func (e *schemaExecutor) ConfigSchema() (string, error) {
	return e.schema, e.err
}

func TestRegisterExecutorSchemaStoresValidSchema(t *testing.T) {
	plugins := &Plugins{ExecutorSchemas: map[string]string{}}

	plugins.registerExecutorSchema("custom", &schemaExecutor{schema: `{"type":"object"}`})

	require.Contains(t, plugins.ExecutorSchemas, "custom")
	assert.Equal(t, `{"type":"object"}`, plugins.ExecutorSchemas["custom"])
}

func TestRegisterExecutorSchemaKeepsExecutorAvailableWithoutSchema(t *testing.T) {
	plugins := &Plugins{
		Executors:       map[string]dkplugin.Executor{},
		ExecutorSchemas: map[string]string{},
	}
	executor := &schemaExecutor{err: errors.New("schema timeout")}

	plugins.Executors["custom"] = executor
	plugins.registerExecutorSchema("custom", executor)

	require.Contains(t, plugins.Executors, "custom")
	assert.NotContains(t, plugins.ExecutorSchemas, "custom")
}
