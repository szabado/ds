package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
	"github.com/kylelemons/godebug/pretty"
	"github.com/logrusorgru/aurora"
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
	Use:   "dsdiff [file 1] [file 2]",
	Short: "Semantic diffs of data serialization languages",
	Args:  cobra.ExactArgs(2),
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
	Run: func(cmd *cobra.Command, args []string) {
		if err := runRoot(args[0], args[1], firstFileLang, secondFileLang); err != nil {
			logrus.WithError(err).Fatal("Fatal error")
		}
	},
}

var (
	firstFileLangArg, secondFileLangArg string
	firstFileLang, secondFileLang       language
	verbose                             bool
)

func init() {
	const fileTypeDesc = "Valid values: [yaml|json|toml]. If unspecified, dsdiff will try to " +
		"automatically determine what the file type is (specifying could be more efficient)."

	RootCmd.PersistentFlags().StringVarP(&firstFileLangArg, "file1type", "1", "", "First file type. "+fileTypeDesc)
	RootCmd.PersistentFlags().StringVarP(&secondFileLangArg, "file2type", "2", "", "Second file type. "+fileTypeDesc)
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

func runRoot(file1, file2 string, lang1, lang2 language) error {
	contents1, err := parse(lang1, file1)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %v", file1)
	}

	contents2, err := parse(lang2, file2)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %v", file2)
	}

	logrus.Debug("Calculating diff")
	output := pretty.Compare(contents1, contents2)
	if output == "" {
		logrus.Debug("Files are identical")
	} else {
		logrus.Debug("Pretty printing output")
		prettyPrint(output)
	}

	return nil
}

func prettyPrint(s string) {
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		if line := sc.Text(); len(line) == 0 {
			fmt.Println(line)
		} else if line[0] == '+' {
			fmt.Println(aurora.Green(line))
		} else if line[0] == '-' {
			fmt.Println(aurora.Red(line))
		} else {
			fmt.Println(line)
		}
	}
}

func parse(lang language, path string) (interface{}, error) {
	logrus := logrus.WithField("path", path)

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file")
	}
	logrus.Debugf("Read %v bytes successfully", len(contents))

	extensionLang := getFileExtLang(path)

	var value interface{}
	switch {
	case lang == langAny && extensionLang == langJson:
		logrus.Debugf("File has JSON extension, assuming the contents are JSON")
		fallthrough
	case lang == langJson:
		logrus.Debug("Calling JSON parser")
		return value, json.Unmarshal(contents, &value)

	case lang == langAny && extensionLang == langYaml:
		logrus.Debug("File has YAML extension, assuming the contents are YAML")
		fallthrough
	case lang == langYaml:
		logrus.Debug("Calling YAML parser")
		return value, yaml.Unmarshal(contents, &value)

	case lang == langAny && extensionLang == langToml:
		logrus.Debug("File has TOML extension, assuming the contents are TOML")
		fallthrough
	case lang == langToml:
		logrus.Debug("Calling TOML parser")
		return value, toml.Unmarshal(contents, &value)

	default:
		logrus.Debug("Unknown file extension and language wasn't specified")

		logrus.Debug("Attempting to use JSON parser")
		err := json.Unmarshal(contents, &value)
		if err == nil {
			logrus.Debug("JSON parser succeeded")
			return value, nil
		}

		logrus.Debug("Attempting to use TOML parser")
		err = toml.Unmarshal(contents, &value)
		if err == nil {
			logrus.Debug("TOML parser succeeded")
			return value, nil
		}

		logrus.Debug("Attempting to use YAML parser")
		err = yaml.Unmarshal(contents, &value)
		if err == nil {
			logrus.Debug("YAML parser succeeded")
			return value, nil
		}

		return nil, errors.Errorf("unable to parse file")
	}
}

func getFileExtLang(file string) language {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".yaml", ".yml":
		return langYaml
	case ".json":
		return langJson
	case ".toml":
		return langToml
	default:
		return langAny
	}
}
