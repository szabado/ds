package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
	"github.com/kylelemons/godebug/pretty"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	const fileTypeDesc = "Valid values: [yaml|json|toml]."

	RootCmd.AddCommand(diffCmd)

	diffCmd.PersistentFlags().StringVarP(&firstFileLangArg, "file1type", "1", "", "First file type. "+fileTypeDesc)
	diffCmd.PersistentFlags().StringVarP(&secondFileLangArg, "file2type", "2", "", "Second file type. "+fileTypeDesc)

}

var diffCmd = &cobra.Command{
	Use:   "diff [file 1] [file 2]",
	Short: "Take semantic diffs of different markup files.",
	Long: `Take semantic diffs of different markup files.

ds diff will try to figure out what type each file is based on the file extension:
  - .yaml/.yml: YAML
  - .toml: TOML
  - .json: JSON

If one of these extensions is matched, ds diff requires the contents of the file
to be of that type.


If the file extension is unknown, it will try a series of parsers until one
works:
  1. JSON
  2. TOML
  3. YAML


You can also specify the file type. These will override the file name, and supported
values are:
  - yaml/yml
  - json
  - toml
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runDiff(args[0], args[1], firstFileLang, secondFileLang); err != nil {
			logrus.WithError(err).Fatal("Fatal error")
		}
	},
}

func runDiff(file1, file2 string, lang1, lang2 language) error {
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

		// Yaml parser is the most permissive and will frequently misinterpret other files.
		// Call it last
		logrus.Debug("Attempting to use YAML parser")
		err = yaml.Unmarshal(contents, &value)
		if err == nil {
			logrus.Debug("YAML parser succeeded")
			return value, nil
		}

		return nil, errors.Errorf("unable to parse file")
	}
}
