package main

import (
	"github.com/sirupsen/logrus"
	"os"

	"github.com/szabado/dsdiff/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		logrus.WithError(err).Trace("Fatal error encountered")

		os.Exit(1)
	}
}
