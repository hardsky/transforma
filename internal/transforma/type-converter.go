package transforma

import (
	"go/ast"
	"go/token"
)

var (
	defaultConverter = &nonConverter{}

	// src type -> set of dst types -> converter from 'src' to 'dst'
	converters = map[string]map[string]typeConverter{
		"string": map[string]typeConverter{
			"int": &funConverter{
				"strconv",
				"Atoi",
				true,
				nil,
			},
		},
	}
)

type typeConverter interface {
	usedImport() string
	ableError() bool
	convert(arg ast.Expr) ast.Expr
}

type funConverter struct {
	Pkg            string
	Fun            string
	ReturnErr      bool
	AdditionalArgs []string
}

func (p *funConverter) fun() ast.Expr {
	return &ast.SelectorExpr{ast.NewIdent(p.Pkg), ast.NewIdent(p.Fun)}
}

func (p *funConverter) convert(arg ast.Expr) ast.Expr {
	return &ast.CallExpr{Fun: p.fun(), Args: []ast.Expr{arg}}
}
func (p *funConverter) ableError() bool {
	return p.ReturnErr
}

func (p *funConverter) usedImport() string {
	return `"` + p.Pkg + `"`
}

type nonConverter struct {
}

func (p *nonConverter) convert(arg ast.Expr) ast.Expr {
	return arg
}
func (p *nonConverter) ableError() bool {
	return false
}
func (p *nonConverter) usedImport() string {
	return ""
}

type objectConverter struct {
	leftObj  *ast.Ident
	rightObj *ast.Ident
	imp      string
}

func (p *objectConverter) convert(left *ast.Field, right *ast.Field) *ast.AssignStmt {
	conv := findTypeConverter(fieldType(right), fieldType(left))
	p.imp = conv.usedImport()

	lhs := []ast.Expr{selectField(p.leftObj, left)}

	if conv.ableError() {
		lhs = append(lhs, ast.NewIdent("_")) // skip error check
	}

	rhs := []ast.Expr{conv.convert(selectField(p.rightObj, right))}

	return &ast.AssignStmt{Lhs: lhs, Rhs: rhs, Tok: token.ASSIGN}
}

func (p *objectConverter) usedImport() string {
	return p.imp
}

func findTypeConverter(srcType string, dstType string) typeConverter {
	dsts := converters[srcType]
	if dsts == nil {
		return defaultConverter
	}

	conv := dsts[dstType]
	if conv == nil {
		return defaultConverter
	}

	return conv
}

func fieldType(a *ast.Field) string {
	ident, ok := a.Type.(*ast.Ident)
	if !ok {
		return ""
	}

	return ident.Name
}
