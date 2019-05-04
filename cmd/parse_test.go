package cmd

import (
	"io"
	"testing"

	r "github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	require := r.New(t)

	require.NoError(runParse("fixtures/test1-a.json", JSON, JSON, &NoopWriter{}))
	require.NoError(runParse("fixtures/test1-a-y", YAML, YAML, &NoopWriter{}))
	require.Error(runParse("fixtures/invalid.json", JSON, JSON, &NoopWriter{}))

	require.NoError(runParse("fixtures/test2-a.json", JSON, JSON, &NoopWriter{}))
	require.NoError(runParse("fixtures/test2-a.json", JSON, YAML, &NoopWriter{}))
	require.NoError(runParse("fixtures/test2-a.json", JSON, TOML, &NoopWriter{}))

	require.NoError(runParse("fixtures/test2-a.json", YAML, JSON, &NoopWriter{}))
	require.NoError(runParse("fixtures/test2-a.json", YAML, YAML, &NoopWriter{}))
	require.NoError(runParse("fixtures/test2-a.json", YAML, TOML, &NoopWriter{}))

	require.NoError(runParse("fixtures/test2-a.toml", TOML, JSON, &NoopWriter{}))
	require.NoError(runParse("fixtures/test2-a.toml", TOML, YAML, &NoopWriter{}))
	require.NoError(runParse("fixtures/test2-a.toml", TOML, TOML, &NoopWriter{}))

	require.NoError(runParse("fixtures/test2-a.yaml", YAML, JSON, &NoopWriter{}))
	require.NoError(runParse("fixtures/test2-a.yaml", YAML, YAML, &NoopWriter{}))
	require.NoError(runParse("fixtures/test2-a.yaml", YAML, TOML, &NoopWriter{}))

	require.Error(runParse("fixtures/test3-a.yaml", Any, TOML, &NoopWriter{}))
	require.Error(runParse("fixtures/test3-a.yaml", Any, JSON, &NoopWriter{}))
	require.NoError(runParse("fixtures/test3-a.yaml", Any, YAML, &NoopWriter{}))
}

var _ io.Writer = (*NoopWriter)(nil)

type NoopWriter struct{}

func (n *NoopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
