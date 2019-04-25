package cmd

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/kylelemons/godebug/pretty"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/szabado/dsdiff/pkg"
)

var RootCmd = &cobra.Command{
	Use:   "dsdiff [old file] [new file]",
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
	firstFileLang, secondFileLang       pkg.Language
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&firstFileLangArg, "first", "1", "",
		"First file type. If unspecified, dsdiff will automatically determine"+
			" what the file type is (specifying could be more efficient). Supported values: yaml, json")
	RootCmd.PersistentFlags().StringVarP(&secondFileLangArg, "second", "2", "",
		"First file type. If unspecified, dsdiff will automatically determine"+
			" what the file type is (specifying could be more efficient). Supported values: yaml, json")
}

func parseLanguageArg(s string) (pkg.Language, error) {
	switch strings.ToLower(s) {
	case "yaml", "yml":
		return pkg.YAML, nil
	case "json":
		return pkg.JSON, nil
	case "":
		return pkg.ANY, nil
	default:
		return pkg.ANY, errors.Errorf("unsupported language %s", s)
	}
}

func runRoot(file1, file2 string, lang1, lang2 pkg.Language) error {
	contents1, err := pkg.Parse(lang1, file1)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %v", file1)
	}

	contents2, err := pkg.Parse(lang2, file2)
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
