package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/kylelemons/godebug/pretty"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	firstFileLangArg, secondFileLangArg string
	firstFileLang, secondFileLang       Language
	prettyDiff bool
)

func init() {
	RootCmd.AddCommand(diffCmd)

	diffCmd.PersistentFlags().StringVarP(&firstFileLangArg, "file1type", "1", "", "first file type.  "+supportedLangsArg)
	diffCmd.PersistentFlags().StringVarP(&secondFileLangArg, "file2type", "2", "", "second file type. "+supportedLangsArg)
	diffCmd.PersistentFlags().BoolVarP(&prettyDiff, "pretty", "p", false, "pretty printed diff. " +
		"May result in unexpected type coercion.")
}

var diffCmd = &cobra.Command{
	Use:   "diff [file 1] [file 2]",
	Short: "Take semantic diffs of different markup files.",
	Args:  cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
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
		handleErr(runDiff(args[0], args[1], firstFileLang, secondFileLang, os.Stdout))
	},
}

func runDiff(file1, file2 string, lang1, lang2 Language, output io.Writer) error {
	contents1, _, err := parse(lang1, file1)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %v", file1)
	}

	contents2, _, err := parse(lang2, file2)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %v", file2)
	}

	logrus.Debug("Calculating diff")
	diff := cmp.Diff(contents1, contents2)
	if prettyDiff {
		pDiff := pretty.Compare(contents1, contents2)
		if pDiff == "" && diff != "" {
			logrus.Debug("No pretty diff found but a diff was present")
		} else {
			diff = pDiff
		}
	}

	if diff == "" {
		logrus.Debug("Files are identical")
		return nil
	}

	if !quiet {
		logrus.Debug("Pretty printing output")
		prettyPrint(diff, output)
	} else {
		logrus.Debug("Quiet output")
	}

	return errOsExit1
}

func prettyPrint(s string, output io.Writer) {
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		if line := sc.Text(); len(line) == 0 {
			fmt.Fprintln(output, line)
		} else if line[0] == '+' {
			fmt.Fprintln(output, aurora.Green(line))
		} else if line[0] == '-' {
			fmt.Fprintln(output, aurora.Red(line))
		} else {
			fmt.Fprintln(output, line)
		}
	}
}

func parse(lang Language, path string) (interface{}, Language, error) {
	logrus := logrus.WithField("path", path)

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, Any, errors.Wrapf(err, "failed to read file")
	}
	logrus.Debugf("Read %v bytes successfully", len(contents))

	extensionLang := getFileExtLang(path)

	var value interface{}
	for _, parser := range parsers {
		if lang == Any && extensionLang == parser.lang {
			logrus.Debugf("File has %[1]s extension, assuming the contents are %[1]s", extensionLang)
		} else if lang != parser.lang {
			// This is not the right parser
			continue
		}

		err = parser.unmarshal(contents, &value)
		return parser.cleanInput(value), parser.lang, err
	}

	logrus.Debug("Unknown file extension and language wasn't specified")

	for _, parser := range parsers {
		logrus.Debugf("Attempting to use %s parser", parser.lang)
		err = parser.unmarshal(contents, &value)
		if err == nil {
			logrus.Debugf("%s parser succeeded", parser.lang)
			return parser.cleanInput(value), parser.lang, nil
		}
	}

	return nil, Any, errors.Errorf("unable to parse file")
}
