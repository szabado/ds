package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/szabado/ds/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		logrus.WithError(err).Trace("Fatal error encountered")

		os.Exit(1)
	}
}
