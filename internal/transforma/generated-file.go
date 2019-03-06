package transforma

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type generatedFile struct {
	fset    *token.FileSet
	name    string
	tagNode *ast.Comment
	file    *ast.File
	mappers []*mapper
}

func genFileName(file string) string {
	dir, name := filepath.Split(file)
	genName := strings.TrimSuffix(name, ".go") + "_gen.go"
	return filepath.Join(dir, genName)
}

func (p *generatedFile) generate() {
	for _, m := range p.mappers {
		m.generate()
	}

	p.revertBuildTag()
}

func (p *generatedFile) printFile() {
	var buf bytes.Buffer
	if err := format.Node(&buf, p.fset, p.file); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf.Bytes())
}

func (p *generatedFile) saveFile() {
	var buf bytes.Buffer
	if err := format.Node(&buf, p.fset, p.file); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(p.name, buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}

func (p *generatedFile) printast(f interface{}) {
	if err := ast.Print(p.fset, p.file); err != nil {
		panic(err)
	}
}

func (p *generatedFile) revertBuildTag() {
	p.tagNode.Text = "//+build !transforma"
}
