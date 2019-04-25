package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

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
	const fileTypeDesc = "Valid values: [yaml|json]. If unspecified, dsdiff will try to " +
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

	output := pretty.Compare(contents1, contents2)
	prettyPrint(output)
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
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file")
	}

	var value interface{}
	switch {
	case lang == langJson || lang == langAny && hasFileExt(langJson, path):
		return value, json.Unmarshal(contents, &value)
	case lang == langYaml || lang == langAny && hasFileExt(langYaml, path):
		return value, yaml.Unmarshal(contents, &value)
	default: // lang == langAny && no known file extension
		if err := json.Unmarshal(contents, &value); err == nil {
			return value, nil
		} else if err = yaml.Unmarshal(contents, &value); err == nil {
			return value, nil
		} else {
			return nil, errors.Errorf("unable to parse file")
		}
	}
}

func hasFileExt(lang language, file string) bool {
	switch strings.ToLower(filepath.Ext(file)) {
	case "yaml", "yml":
		return lang == langYaml
	case "json":
		return lang == langJson
	default:
		return false
	}
}
