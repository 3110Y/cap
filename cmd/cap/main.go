package main

import (
	"fmt"
	"os"

	"github.com/3110Y/cap/internal/integration/di"
)

func main() {
	rootCmd, err := di.InitializeCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize CLI: %v\n", err)
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
