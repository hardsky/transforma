package transforma

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestNameFromTag(t *testing.T) {
	src := `
package src

type A struct{
  F string ` + "`trf:\"testName\"`" +
		`
}
`
	want := "testName"

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Error(err)
		return
	}

	fieldWasFound := false

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Field:
			fieldWasFound = true
			got := nameFromTag(x)
			if want != got {
				t.Errorf("expected %q, but got %q", want, got)
			}
			return false
		}

		return true
	})

	if !fieldWasFound {
		t.Error("did not found struct field for test")
	}
}
