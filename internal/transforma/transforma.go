package transforma

import (
	"context"
	"os"

	"golang.org/x/tools/go/packages"
)

func Generate(pkg string) error {

	pkgs, err := loadPackages(pkg)
	if err != nil {
		return err
	}

	a := &analyzer{
		pkgs: pkgs,
	}

	files := a.analyze()
	for _, f := range files {
		f.generate()
		f.saveFile()
	}

	return nil
}

func loadPackages(pkg string) ([]*packages.Package, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	env := os.Environ()
	cfg := &packages.Config{
		Context:    context.Background(),
		Mode:       packages.LoadAllSyntax,
		Dir:        wd,
		Env:        env,
		BuildFlags: []string{"-tags=transforma"},
	}
	return packages.Load(cfg, "pattern="+pkg)
}
