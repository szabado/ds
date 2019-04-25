package pkg

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

//go:generate stringer -type=Language
type Language int

const (
	ANY Language = iota
	JSON
	YAML
)

func Parse(lang Language, path string) (interface{}, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file")
	}

	var value interface{}
	switch {
	case lang == JSON || lang == ANY && hasExt(JSON, path):
		return value, json.Unmarshal(contents, &value)
	case lang == YAML || lang == ANY && hasExt(YAML, path):
		return value, yaml.Unmarshal(contents, &value)
	default: // lang == ANY && no known file extension
		if err := json.Unmarshal(contents, &value); err == nil {
			return value, nil
		} else if err = yaml.Unmarshal(contents, &value); err == nil {
			return value, nil
		} else {
			return nil, errors.Errorf("no known parser succeeded", path)
		}
	}
}

func hasExt(lang Language, file string) bool {
	fileExt := filepath.Ext(file)

	switch lang {
	case YAML:
		return strings.EqualFold(YAML.String(), fileExt) || strings.EqualFold("yml", fileExt)
	case ANY:
		return false
	default:
		return strings.EqualFold(lang.String(), fileExt)
	}
}
