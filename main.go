package main

import (
	"fmt"
	"os"

	"github.com/szabado/dsdiff/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Printf("Fatal error: %s", err)
		os.Exit(1)
	}
}
