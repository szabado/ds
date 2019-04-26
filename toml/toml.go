package toml

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

func Unmarshal(p []byte, v interface{}) error {
	return toml.Unmarshal(p, v)
}

func Marshal(v interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := toml.NewEncoder(buffer)
	err := encoder.Encode(v)
	return buffer.Bytes(), err
}
