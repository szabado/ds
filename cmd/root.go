package cmd

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//go:generate stringer -type=language
type language int

const (
	langAny language = iota
	langJson
	langYaml
	langToml
)

var RootCmd = &cobra.Command{
	Use:   "ds",
	Short: "A swiss army tool for markup languages like json, yaml, and toml.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		firstFileLang, err = parseLanguageArg(firstFileLangArg)
		if err != nil {
			logrus.WithError(err).Fatal("Invalid language specified")
		}

		secondFileLang, err = parseLanguageArg(secondFileLangArg)
		if err != nil {
			logrus.WithError(err).Fatal("Invalid language specified")
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
	firstFileLang, secondFileLang       language
	verbose                             bool
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose log output.")
}

func parseLanguageArg(s string) (language, error) {
	switch strings.ToLower(s) {
	case "yaml", "yml":
		return langYaml, nil
	case "json":
		return langJson, nil
	case "toml":
		return langToml, nil
	case "":
		return langAny, nil
	default:
		return langAny, errors.Errorf("unsupported language %s", s)
	}
}

func getFileExtLang(file string) language {
	switch ext := strings.ToLower(filepath.Ext(file)); ext {
	case ".yaml", ".yml":
		return langYaml
	case ".json":
		return langJson
	case ".toml":
		return langToml
	default:
		logrus.WithField("extension", ext).Debug("Unknown file extension")
		return langAny
	}
}
