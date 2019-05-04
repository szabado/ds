package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
	r "github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	var writer io.Writer = &NoopWriter{}
	logrus.SetLevel(logrus.FatalLevel)
	writer = os.Stdout

	RunSuite(t, func(file1, file2 string) error {
		return runDiff(file1, file2, Any, Any, writer)
	})
}

const (
	fixturesDir = `fixtures`
)

var (
	testFilesRegexes = regexp.MustCompile(`^test(\d+)-(\w+)[-.]\w+$`)
)

func RunSuite(t *testing.T, f func(file1, file2 string) error) {
	require := r.New(t)

	files, err := ioutil.ReadDir(fixturesDir)
	require.NoError(err)

	testFiles := make(map[string]map[string][]string)
	for _, file := range files {
		submatches := testFilesRegexes.FindStringSubmatch(file.Name())
		if len(submatches) != 3 {
			t.Logf("Skipping test: %s", file.Name())
			continue
		}

		testNum := submatches[1]
		inputNum := submatches[2]

		if _, ok := testFiles[testNum]; !ok {
			testFiles[testNum] = make(map[string][]string)
		}

		testFiles[testNum][inputNum] = append(testFiles[testNum][inputNum], file.Name())
	}

	for _, suite := range testFiles {
		for key, tests := range suite {
			t.Log("Asserting files are the same")
			for i := 0; i < len(tests); i++ {
				for j := i; j < len(tests); j++ {
					testI := tests[i]
					testJ := tests[j]

					file1 := path.Join(fixturesDir, testI)
					file2 := path.Join(fixturesDir, testJ)

					t.Run(fmt.Sprintf("%s__%s", testI, testJ), func(t *testing.T) {
						r.NoError(t, f(file1, file2))
					})

					t.Run(fmt.Sprintf("%s__%s", testJ, testI), func(t *testing.T) {
						r.NoError(t, f(file2, file1))
					})
				}
			}

			t.Log("Asserting files are different")
			for _, test := range tests {
				for altKey, altTests := range suite {
					if altKey == key {
						continue
					}

					for _, altTest := range altTests {
						file1 := path.Join(fixturesDir, test)
						file2 := path.Join(fixturesDir, altTest)

						t.Run(fmt.Sprintf("%s__%s", test, altTest), func(t *testing.T) {
							r.Error(t, f(file1, file2))
						})

						t.Run(fmt.Sprintf("%s__%s", altTest, test), func(t *testing.T) {
							r.Error(t, f(file2, file1))
						})
					}
				}
			}
		}
	}
}
