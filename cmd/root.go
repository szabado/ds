package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//go:generate stringer -type=Language
type Language int

const (
	Any Language = iota
	JSON
	YAML
	TOML
)

type parser struct {
	lang      Language
	unmarshal func([]byte, interface{}) error
}

var parsers = []parser{
	{
		lang:      JSON,
		unmarshal: json.Unmarshal,
	},
	{
		lang:      TOML,
		unmarshal: toml.Unmarshal,
	},
	{
		// Yaml parser is the most permissive and will frequently misinterpret other files.
		// Call it last
		lang:      YAML,
		unmarshal: yaml.Unmarshal,
	},
}

var errOsExit1 = errors.New("ds should os.Exit(1)")

var RootCmd = &cobra.Command{
	Use:   "ds",
	Short: "A swiss army tool for markup languages like json, yaml, and toml.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		firstFileLang, err = parseLanguageArg(firstFileLangArg)
		if err != nil {
			logrus.WithError(err).Fatal("Invalid Language specified")
		}

		secondFileLang, err = parseLanguageArg(secondFileLangArg)
		if err != nil {
			logrus.WithError(err).Fatal("Invalid Language specified")
		}

		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.FatalLevel)
		}
	},
}

var (
	firstFileLangArg, secondFileLangArg string
	firstFileLang, secondFileLang       Language
	verbose, quiet                      bool
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
