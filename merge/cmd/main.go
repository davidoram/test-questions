package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/davidoram/test-questions/merge"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "merge: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout *os.File) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: merge file1 file2 [file...]")
	}

	files := make([]*os.File, 0, len(args))
	for _, name := range args {
		file, err := os.Open(name)
		if err != nil {
			closeFiles(files)
			return fmt.Errorf("open %q: %w", name, err)
		}
		files = append(files, file)
	}
	defer closeFiles(files)

	inputs := make([]*bufio.Scanner, len(files))
	for i, file := range files {
		inputs[i] = bufio.NewScanner(file)
	}

	out := bufio.NewWriter(stdout)
	if err := merge.MergeStreams(out, strings.Compare, inputs...); err != nil {
		return err
	}
	if err := out.Flush(); err != nil {
		return err
	}

	return nil
}

func closeFiles(files []*os.File) {
	for _, file := range files {
		_ = file.Close()
	}
}
