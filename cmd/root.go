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
	"github.com/szabado/ds/xml"
)

//go:generate stringer -type=Language
type Language int

const supportedLangsArg = "Valid types: [yaml|json|toml|xml]"
const (
	Any Language = iota
	JSON
	TOML
	XML
	YAML
)

type parser struct {
	lang       Language
	unmarshal  func([]byte, interface{}) error
	marshal    func(v interface{}) ([]byte, error)
	cleanInput func(interface{}) interface{}
}

func jsonMarshalPretty(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

var parsers = []*parser{
	{
		lang:       JSON,
		unmarshal:  json.Unmarshal,
		marshal:    jsonMarshalPretty,
		cleanInput: stringKeyedMapCleaner,
	},
	{
		lang:       TOML,
		unmarshal:  toml.Unmarshal,
		marshal:    toml.Marshal,
		cleanInput: stringKeyedMapCleaner,
	},
	{
		lang:       XML,
		unmarshal:  xml.Unmarshal,
		marshal:    xml.Marshal,
		cleanInput: stringKeyedMapCleaner,
	},
	{
		// Yaml parser is the most permissive and will frequently misinterpret other files.
		// Call it last
		lang:       YAML,
		unmarshal:  yaml.Unmarshal,
		marshal:    yaml.Marshal,
		cleanInput: stringKeyedMapCleaner,
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
  - .xml:  XML

If the file extension is unknown, it will try a series of parsers until one
works (in this order):
  1. JSON
  2. TOML
  3. XML
  3. YAML


If it fails, or you want to be extra sure it's using the right parser, you can 
also specify the file type. These will override the file name, and supported
values are:
  - yaml/yml
  - json
  - toml
  - xml
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.FatalLevel)
		}

		if compact {
			for _, parser := range parsers {
				if parser.lang == JSON {
					parser.marshal = json.Marshal
				}
			}
		}
	},
}

var (
	verbose, quiet, compact bool
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"verbose log output")
	RootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false,
		"output nothing. A non-zero exit code indicates failure")
	RootCmd.Flags().BoolVarP(&compact, "compact", "c", false,
		"compress/minify output where possible")
}

func parseLanguageArg(s string) (Language, error) {
	switch strings.ToLower(s) {
	case "yaml", "yml":
		return YAML, nil
	case "json":
		return JSON, nil
	case "toml":
		return TOML, nil
	case "xml":
		return XML, nil
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
