package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	require := require.New(t)

	require.NoError(runParse("fixtures/test1-a.json"))
	require.NoError(runParse("fixtures/test1-a.yaml"))
	require.Error(runParse("fixtures/invalid.json"))
}
