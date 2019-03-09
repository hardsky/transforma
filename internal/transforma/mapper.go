package transforma

import (
	"go/ast"
	"go/token"
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

	m.Fnc.Body.List = append(m.Fnc.Body.List, &ast.ReturnStmt{Results: []ast.Expr{res}})
}

func selectField(obj *ast.Ident, fl *ast.Field) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: obj, Sel: ast.NewIdent(fl.Names[0].Name)}
}

func findField(src *ast.StructType, same *ast.Field) *ast.Field {
	for _, fl := range src.Fields.List {
		if same.Names[0].Name == fl.Names[0].Name {
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
