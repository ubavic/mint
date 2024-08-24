package writer

import (
	"fmt"
	"strings"

	"github.com/ubavic/mint/parser"
	"github.com/ubavic/mint/schema"
)

// TODO: use io.Writer for output
// TODO: cache commandExpression in map
func Write(targetSchema *schema.Target, element parser.Element) string {
	switch v := element.(type) {
	case *parser.TextContent:
		return v.String()

	case *parser.Block:
		result := ""
		for _, e := range v.Content() {
			result += Write(targetSchema, e)
		}
		return result
	case *parser.Command:
		var commandExpression *string
		for _, c := range targetSchema.Commands {
			if c.Command == v.Name {
				commandExpression = &c.Expression
			}
		}

		if commandExpression == nil {
			panic("Command not found: " + v.Name)
		}

		result := *commandExpression
		for i, a := range v.Arguments {
			r := Write(targetSchema, a)
			result = strings.ReplaceAll(result, fmt.Sprintf("$%d", i+1), r)
		}

		return result
	default:
		return element.String()
	}
}
