package mermaid

import (
	"fmt"
	"strings"

	"github.com/adam-baker/vsdxparser/parser"
)

func GenerateMermaid(shapes []parser.Shape, connections []parser.Connect) (string, error) {
	if len(shapes) == 0 {
		return "", fmt.Errorf("no shapes found in input")
	}

	var b strings.Builder
	b.WriteString("graph TD\n")

	// Keep track of shape IDs and readable node names
	idMap := make(map[string]string)
	nameCount := make(map[string]int)

	for i, shape := range shapes {
		label := strings.TrimSpace(shape.Text)
		if label == "" {
			continue
		}

		// De-dupe duplicate labels for clarity
		nameCount[label]++
		if nameCount[label] > 1 {
			label = fmt.Sprintf("%s (%d)", label, nameCount[label])
		}

		nodeID := fmt.Sprintf("n%d", i+1)
		idMap[shape.ID] = nodeID
		b.WriteString(fmt.Sprintf("    %s[%q]\n", nodeID, label))
	}

	for _, conn := range connections {
		from, okFrom := idMap[conn.FromID]
		to, okTo := idMap[conn.ToID]
		if okFrom && okTo {
			b.WriteString(fmt.Sprintf("    %s --> %s\n", from, to))
		}
	}

	return b.String(), nil
}
