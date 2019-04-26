package toml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const fixture = `title = "TOML Example"

[owner]
  name = "Felix"
`

func TestRoundTrip(t *testing.T) {
	require := require.New(t)

	var contents interface{}

	require.NoError(Unmarshal([]byte(fixture), &contents))
	output, err := Marshal(contents)
	require.NoError(err)
	require.Equal(fixture, string(output))
}
