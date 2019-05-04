package toml

import (
	"testing"

	r "github.com/stretchr/testify/require"
)

const fixture = `title = "TOML Example"

[owner]
  name = "Felix"
`

func TestRoundTrip(t *testing.T) {
	require := r.New(t)

	var contents interface{}

	require.NoError(Unmarshal([]byte(fixture), &contents))
	output, err := Marshal(contents)
	require.NoError(err)
	require.Equal(fixture, string(output))
}
