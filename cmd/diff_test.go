package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	require := require.New(t)

	require.NoError(runDiff("fixtures/test1-a.json", "fixtures/test1-a-y", Any, Any))
	require.NoError(runDiff("fixtures/test1-a-y", "fixtures/test1-a.json", Any, Any))
	require.Equal(errOsExit1, runDiff("fixtures/test1-a.json", "fixtures/test1-b-j", Any, Any))
	require.Equal(errOsExit1, runDiff("fixtures/test1-b-j", "fixtures/test1-a.json", Any, Any))
	require.Equal(errOsExit1, runDiff("fixtures/test1-a-y", "fixtures/test1-b-j", Any, Any))
	require.Equal(errOsExit1, runDiff("fixtures/test1-b-j", "fixtures/test1-a-y", Any, Any))

	require.NoError(runDiff("fixtures/test2-a.json", "fixtures/test2-a.toml", JSON, TOML))
	require.NoError(runDiff("fixtures/test2-a.toml", "fixtures/test2-a.yaml", TOML, YAML))
}
