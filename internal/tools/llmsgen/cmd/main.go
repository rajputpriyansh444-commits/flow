package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rajputpriyansh444-commits/flow/internal/tools/llmsgen"
)

func main() {
	sourceDir := flag.String("source", filepath.Join("skills", "flow"), "directory containing the bundled Flow skill")
	sourcePrefix := flag.String("source-prefix", filepath.ToSlash(filepath.Join("skills", "flow")), "source path prefix written into llms.txt")
	outputPath := flag.String("output", "llms.txt", "path to write")
	flag.Parse()

	err := llmsgen.WriteFile(*outputPath, llmsgen.Options{
		SourceDir:    *sourceDir,
		SourcePrefix: *sourcePrefix,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
