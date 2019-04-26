package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/szabado/ds/toml"
)

//go:generate stringer -type=Language
type Language int

const supportedLangsArg = "Valid types: [yaml|json|toml]"
const (
	Any Language = iota
	JSON
	YAML
	TOML
)

type parser struct {
	lang      Language
	unmarshal func([]byte, interface{}) error
	marshal   func(v interface{}) ([]byte, error)
}

var parsers = []parser{
	{
		lang:      JSON,
		unmarshal: json.Unmarshal,
		marshal:   json.Marshal,
	},
	{
		lang:      TOML,
		unmarshal: toml.Unmarshal,
		marshal:   toml.Marshal,
	},
	{
		// Yaml parser is the most permissive and will frequently misinterpret other files.
		// Call it last
		lang:      YAML,
		unmarshal: yaml.Unmarshal,
		marshal:   yaml.Marshal,
	},
}

var errOsExit1 = errors.New("ds should os.Exit(1)")

var RootCmd = &cobra.Command{
	Use:   "ds",
	Short: "A swiss army tool for markup languages like json, yaml, and toml.",
	Long: `A swiss army tool for markup languages like json, yaml, and toml.

ds will try to figure out automagically what the type of any file passed in
is. It uses a couple methods to do this. First, the file extension:
  - .yaml/.yml: YAML
  - .toml: TOML
  - .json: JSON

If the file extension is unknown, it will try a series of parsers until one
works (in this order):
  1. JSON
  2. TOML
  3. YAML


If it fails, or you want to be extra sure it's using the right parser, you can 
also specify the file type. These will override the file name, and supported
values are:
  - yaml/yml
  - json
  - toml
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.FatalLevel)
		}
	},
}

var (
	verbose, quiet bool
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose log output")
	RootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "output nothing."+
		"an nonzero exit code indicates failure")
}

func parseLanguageArg(s string) (Language, error) {
	switch strings.ToLower(s) {
	case "yaml", "yml":
		return YAML, nil
	case "json":
		return JSON, nil
	case "toml":
		return TOML, nil
	case "":
		return Any, nil
	default:
		return Any, errors.Errorf("unsupported Language %s", s)
	}
}

func getFileExtLang(file string) Language {
	ext := filepath.Ext(file)
	lang, _ := parseLanguageArg(strings.TrimPrefix(strings.ToLower(ext), "."))
	if lang == Any {
		logrus.WithField("ext", ext).Debug("Unknown file extension")
	}

	return lang
}

// handleErr either returns successfully if there's no error, or calls os.Exit(). This
// should only be used in the Run section of cobra.Command.
func handleErr(err error) {
	if err == nil {
		return
	}

	if err == errOsExit1 {
		os.Exit(1)
	} else if quiet {
		os.Exit(1)
	} else {
		logrus.WithError(err).Fatal("Fatal error")
	}
}
