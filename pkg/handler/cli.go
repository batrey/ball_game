package handler

import (
	"flag"
	"fmt"
	"io"
)

// TouchCounter calculates the maximum touch count from an input file.
type TouchCounter interface {
	MaxTouchesFromFile(fileName string) (int, error)
}

// RunCLI handles command-line input, output, and exit codes.
func RunCLI(args []string, stdout io.Writer, stderr io.Writer, touchCounter TouchCounter) int {
	flags := flag.NewFlagSet("ballgame", flag.ContinueOnError)
	flags.SetOutput(stderr)

	fileName := flags.String("file", "", "Path to the input file")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	if *fileName == "" {
		fmt.Fprintln(stderr, "missing required -file flag")
		return 2
	}

	maxTouches, err := touchCounter.MaxTouchesFromFile(*fileName)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}

	fmt.Fprintln(stdout, maxTouches)
	return 0
}
