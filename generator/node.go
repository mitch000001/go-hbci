package generator

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
)

func nodeToString(node ast.Node, fileSet *token.FileSet) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fileSet, node)
	return buf.String()
}
