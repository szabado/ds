package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	require := require.New(t)

	require.NoError(runParse("fixtures/test1-a.json", JSON, JSON))
	require.NoError(runParse("fixtures/test1-a-y", YAML, YAML))
	require.Error(runParse("fixtures/invalid.json", JSON, JSON))

	require.NoError(runParse("fixtures/test2-a.json", JSON, JSON))
	require.NoError(runParse("fixtures/test2-a.json", JSON, YAML))
	require.NoError(runParse("fixtures/test2-a.json", JSON, TOML))

	require.NoError(runParse("fixtures/test2-a.json", YAML, JSON))
	require.NoError(runParse("fixtures/test2-a.json", YAML, YAML))
	require.NoError(runParse("fixtures/test2-a.json", YAML, TOML))

	require.NoError(runParse("fixtures/test2-a.toml", TOML, JSON))
	require.NoError(runParse("fixtures/test2-a.toml", TOML, YAML))
	require.NoError(runParse("fixtures/test2-a.toml", TOML, TOML))

	require.NoError(runParse("fixtures/test2-a.yaml", YAML, JSON))
	require.NoError(runParse("fixtures/test2-a.yaml", YAML, YAML))
	require.NoError(runParse("fixtures/test2-a.yaml", YAML, TOML))
}
