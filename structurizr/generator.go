package structurizr

import (
	"fmt"
	"github.com/adam-baker/vdsxparser/parser"
	"strings"
)

func GenerateDSL(shapes []parser.Shape, connections []parser.Connect) (string, error) {
	if len(shapes) == 0 {
		return "", fmt.Errorf("no shapes found in input")
	}

	var b strings.Builder
	b.WriteString("workspace {\n\n  model {\n")

	idMap := make(map[string]string)
	for i, shape := range shapes {
		ref := fmt.Sprintf("s%d", i+1)
		idMap[shape.ID] = ref
		b.WriteString(fmt.Sprintf("    %s = container \"%s\"\n", ref, shape.Text))
	}

	for _, c := range connections {
		fromRef, fromOK := idMap[c.FromID]
		toRef, toOK := idMap[c.ToID]
		if fromOK && toOK {
			b.WriteString(fmt.Sprintf("    %s -> %s\n", fromRef, toRef))
		}
	}

	b.WriteString("  }\n\n  views {\n    systemContext * {\n      include *\n      autolayout lr\n    }\n  }\n\n}\n")
	return b.String(), nil
}
