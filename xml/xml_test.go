package xml

import (
	"fmt"
	"testing"

	r "github.com/stretchr/testify/require"
)

const fixture = `<?xml version="1.0" encoding="UTF-8"?>
<alpha>
    <beta>
        <gamma>delta</gamma>
    </beta>
</alpha>
`

const output = `<doc><alpha><beta><gamma>delta</gamma></beta></alpha></doc>`

func TestRoundTrip(t *testing.T) {
	require := r.New(t)

	outputStruct := map[string]interface{}{
		"alpha": map[string]interface{}{
			"beta": map[string]interface{}{
				"gamma": "delta",
			},
		},
	}

	var value1 interface{}
	require.NoError(Unmarshal([]byte(fixture), &value1))

	require.Equal(outputStruct, value1)

	var value2 map[string]interface{}
	require.NoError(Unmarshal([]byte(fixture), &value2))

	require.Equal(outputStruct, value2)

	output1, err := Marshal(value1)
	require.NoError(err)
	require.Equal(output, fmt.Sprintf("%s", output1))

	output2, err := Marshal(value2)
	require.NoError(err)
	require.Equal(output, fmt.Sprintf("%s", output2))
}
