package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	inputFileLangArg, outputFileLangArg string
	inputFileLang, outputFileLang       Language
)

func init() {
	RootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringVarP(&inputFileLangArg, "input", "i", "", "input file type.  "+supportedLangsArg)
	parseCmd.Flags().StringVarP(&outputFileLangArg, "output", "o", "", "output file type. "+supportedLangsArg)
}

var parseCmd = &cobra.Command{
	Use:   "parse [file]",
	Short: "Parse and re-encode files in different formats",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		inputFileLang, err = parseLanguageArg(inputFileLangArg)
		if err != nil {
			logrus.WithError(err).Fatal("Invalid language specified")
		}

		outputFileLang, err = parseLanguageArg(outputFileLangArg)
		if err != nil {
			logrus.WithError(err).Fatal("Invalid language specified")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleErr(runParse(args[0]))
	},
}

func runParse(file string) error {
	contents, lang, err := parse(inputFileLang, file)
	if err != nil {
		return errors.Wrap(err, "failed to parse file")
	}

	if outputFileLang == Any {
		outputFileLang = lang
	}

	logrus.WithFields(logrus.Fields{
		"input":  lang,
		"output": outputFileLang,
	}).Debug("Beginning conversion")

	var result []byte
	for _, parser := range parsers {
		if parser.lang != outputFileLang {
			continue
		}

		logrus.Debug()
		result, err = parser.marshal(contents)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
	}

	if !quiet {
		fmt.Printf("%s\n", result)
	}

	return nil
}
