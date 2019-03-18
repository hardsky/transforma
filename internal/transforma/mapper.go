package transforma

import (
	"go/ast"
	"go/token"
	"strings"
)

const (
	tag      = "trf"
	skipName = "-"
)

type mapper struct {
	Fil     *ast.File
	Fnc     *ast.FuncDecl
	In      *ast.Field
	Out     *ast.Field
	InType  *ast.StructType
	OutType *ast.StructType
}

func (m *mapper) generate() {
	res := ast.NewIdent("res")
	m.Fnc.Body.List = []ast.Stmt{define(res, m.Out)}

	dst := m.OutType
	src := m.InType
	for _, fl := range dst.Fields.List {
		rightField := findField(src, fl)
		if rightField == nil {
			continue
		}

		left := selectField(res, fl)
		right := selectField(m.In.Names[0], rightField)
		assig := ast.AssignStmt{Lhs: []ast.Expr{left}, Tok: token.ASSIGN, Rhs: []ast.Expr{right}}
		m.Fnc.Body.List = append(m.Fnc.Body.List, &assig)
	}

	ret := &ast.ReturnStmt{Results: []ast.Expr{res}}
	m.Fnc.Body.List = append(m.Fnc.Body.List, ret)
	m.Fnc.Body.Rbrace = ret.End()
}

func selectField(obj *ast.Ident, fl *ast.Field) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: obj, Sel: ast.NewIdent(fl.Names[0].Name)}
}

func findField(src *ast.StructType, left *ast.Field) *ast.Field {
	rightName := nameFromTag(left)
	if rightName == skipName {
		return nil
	}

	if len(rightName) > 0 {
		return findFieldByName(src, rightName)
	}

	leftName := left.Names[0].Name
	fl := findFieldByTag(src, leftName)
	if fl != nil {
		return fl
	}

	rightField := findFieldByName(src, leftName)
	if rightField == nil {
		return nil
	}

	if nameFromTag(rightField) == skipName {
		return nil
	}

	return rightField
}

func findFieldByTag(src *ast.StructType, fieldName string) *ast.Field {
	for _, fl := range src.Fields.List {
		name := nameFromTag(fl)
		if name == fieldName {
			return fl
		}
	}
	return nil
}

func findFieldByName(src *ast.StructType, name string) *ast.Field {
	for _, fl := range src.Fields.List {
		if fl.Names[0].Name == name {
			return fl
		}
	}

	return nil
}

func define(x *ast.Ident, tp *ast.Field) *ast.AssignStmt {
	var right ast.Expr
	switch v := tp.Type.(type) {
	case *ast.StarExpr:
		right = &ast.UnaryExpr{Op: token.AND, X: &ast.CompositeLit{Type: v.X}} // &B{}
	case *ast.Ident:
		right = &ast.CompositeLit{Type: v} // B{}
	}
	return &ast.AssignStmt{Lhs: []ast.Expr{x}, Tok: token.DEFINE, Rhs: []ast.Expr{right}}
}

func nameFromTag(fl *ast.Field) string {
	if fl.Tag == nil {
		return ""
	}

	tags := strings.Split(fl.Tag.Value, " ")
	for _, t := range tags {
		if strings.Contains(t, tag) {
			kv := strings.Split(t, ":")
			if len(kv) != 2 {
				return ""
			}

			return strings.Trim(kv[1], "\"`")
		}
	}

	return ""
}

// func conversion() ast.Stmt {
// 	// Type-check the package.
// 	// We create an empty map for each kind of input
// 	// we're interested in, and Check populates them.
// 	info := types.Info{
// 		Types: make(map[ast.Expr]types.TypeAndValue),
// 		Defs:  make(map[*ast.Ident]types.Object),
// 		Uses:  make(map[*ast.Ident]types.Object),
// 	}
// 	var conf types.Config
// 	pkg, err := conf.Check("fib", fset, []*ast.File{f}, &info)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
