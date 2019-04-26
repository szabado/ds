package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	inputFileLangArg, outputFileLangArg string
	inputFileLang, outputFileLang       Language
	compact                             bool
)

func init() {
	RootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringVarP(&inputFileLangArg, "input", "i", "", "input file type.  "+supportedLangsArg)
	parseCmd.Flags().StringVarP(&outputFileLangArg, "output", "o", "", "output file type. "+supportedLangsArg)
	parseCmd.Flags().BoolVarP(&compact, "compact", "c", false, "compress/minify output where possible")
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

		if compact {
			for _, parser := range parsers {
				if parser.lang == JSON {
					parser.marshal = json.Marshal
				}
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleErr(runParse(args[0], inputFileLang, outputFileLang))
	},
}

func runParse(file string, inputLang, outputLang Language) error {
	contents, lang, err := parse(inputLang, file)
	if err != nil {
		return errors.Wrap(err, "failed to parse file")
	}

	if outputLang == Any {
		outputLang = lang
	}

	logrus.WithFields(logrus.Fields{
		"input":  lang,
		"output": outputLang,
	}).Debug("Beginning conversion")

	var result []byte
	for _, parser := range parsers {
		if parser.lang != outputLang {
			continue
		}

		result, err = parser.marshal(parser.cleanInput(contents))
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
	}

	if !quiet {
		fmt.Printf("%s\n", result)
	}

	return nil
}
