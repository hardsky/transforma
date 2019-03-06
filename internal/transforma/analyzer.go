package transforma

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

type analyzer struct {
	pkgs []*packages.Package
}

func (p *analyzer) analyze() []*generatedFile {
	res := make([]*generatedFile, 0)
	for _, pkg := range p.pkgs {
		for i, f := range pkg.Syntax {

			tagNode := p.findTagComment(f) // search +build transforma
			if tagNode == nil {
				continue
			}

			mappers := p.findMappers(pkg.Fset, f)

			res = append(res, &generatedFile{
				fset:    pkg.Fset,
				name:    genFileName(pkg.CompiledGoFiles[i]), // suggestion, that Syntax related to CompiledGoFiles
				tagNode: tagNode,
				file:    f,
				mappers: mappers,
			})
		}
	}

	return res
}

func (p *analyzer) findTagComment(f *ast.File) *ast.Comment {
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if strings.Contains(c.Text, "+build") &&
				strings.Contains(c.Text, "transforma") {
				return c
			}
		}
	}
	return nil
}

func (p *analyzer) findMappers(fset *token.FileSet, f *ast.File) []*mapper {

	res := make([]*mapper, 0)
	for _, decl := range f.Decls {
		if fun, ok := decl.(*ast.FuncDecl); ok {
			res = append(res, &mapper{
				f,
				fun,
				fun.Type.Params.List[0],
				fun.Type.Results.List[0],
				p.findType(fun.Type.Params.List[0]),
				p.findType(fun.Type.Results.List[0]),
			})
		}
	}

	return res
}

func (p *analyzer) findType(fl *ast.Field) *ast.StructType {

	name := typeName(fl)

	for _, pkg := range p.pkgs {
		for _, f := range pkg.Syntax {
			// TODO: use f.Scope
			for _, dcl := range f.Decls {
				if t, ok := dcl.(*ast.GenDecl); ok && t.Tok == token.TYPE {

					if tp, ok := t.Specs[0].(*ast.TypeSpec); ok && tp.Name.Name == name.Name {
						return tp.Type.(*ast.StructType)
					}
				}

			}
		}
	}

	return nil
}

func typeName(tp *ast.Field) *ast.Ident {
	switch v := tp.Type.(type) {
	case *ast.StarExpr: // *B
		return v.X.(*ast.Ident)
	case *ast.Ident: // B
		return v
	}

	return &ast.Ident{}
}
