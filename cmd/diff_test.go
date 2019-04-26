package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	require := require.New(t)

	require.NoError(runDiff("fixtures/test1-a.json", "fixtures/test1-a.yaml", Any, Any))
	require.NoError(runDiff("fixtures/test1-a.yaml", "fixtures/test1-a.json", Any, Any))
	require.NoError(runDiff("fixtures/test1-a.json", "fixtures/test1-a.yaml", YAML, Any))

	require.Equal(errOsExit1, runDiff("fixtures/test1-a.json", "fixtures/test1-b.json", Any, Any))
	require.Equal(errOsExit1, runDiff("fixtures/test1-b.json", "fixtures/test1-a.json", Any, Any))
	require.Equal(errOsExit1, runDiff("fixtures/test1-a.yaml", "fixtures/test1-b.yaml", Any, Any))
	require.Equal(errOsExit1, runDiff("fixtures/test1-b.yaml", "fixtures/test1-a.yaml", Any, Any))
	require.Equal(errOsExit1, runDiff("fixtures/test1-a.yaml", "fixtures/test1-b.json", Any, Any))
	require.Equal(errOsExit1, runDiff("fixtures/test1-b.yaml", "fixtures/test1-a.json", Any, Any))
}
