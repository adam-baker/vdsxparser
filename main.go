package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adam-baker/vdsxparser/parser"
	"github.com/adam-baker/vdsxparser/structurizr"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: vsdx_to_structurizr <input.vsdx>")
		os.Exit(1)
	}

	input := os.Args[1]

	shapes, connectors, err := parser.ExtractVSDX(input)
	if err != nil {
		handleError("extracting VSDX", err)
	}

	dsl, err := structurizr.GenerateDSL(shapes, connectors)
	if err != nil {
		handleError("generating Structurizr DSL", err)
	}

	outFile := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input)) + ".dsl"
	err = os.WriteFile(outFile, []byte(dsl), 0644)
	if err != nil {
		handleError("writing output file", err)
	}

	fmt.Printf("✅ Structurizr DSL generated: %s\n", outFile)
}

func handleError(context string, err error) {
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", context, err)
	os.Exit(1)
}
