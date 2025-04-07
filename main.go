package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adam-baker/vdsxparser/mermaid"
	"github.com/adam-baker/vdsxparser/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: vsdx_to_mermaid <input.vsdx>")
		os.Exit(1)
	}

	input := os.Args[1]

	shapes, connectors, err := parser.ExtractVSDX(input)
	if err != nil {
		handleError("extracting VSDX", err)
	}

	diagram, err := mermaid.GenerateMermaid(shapes, connectors)
	if err != nil {
		handleError("generating Mermaid diagram", err)
	}

	outFile := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input)) + ".mmd"
	err = os.WriteFile(outFile, []byte(diagram), 0644)
	if err != nil {
		handleError("writing output file", err)
	}

	fmt.Printf("✅ Mermaid diagram generated: %s\n", outFile)
}

func handleError(context string, err error) {
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", context, err)
	os.Exit(1)
}
