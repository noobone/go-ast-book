package demo

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

type packageLoader struct {
	*packages.Config
	PkgFilter []func(pkgPath string) bool
}

func NewPackageLoader() *packageLoader {
	return &packageLoader{}
}

func (pl *packageLoader) RecursionParsePkg(pkgName string, pkgMap map[string]*packages.Package) {
	pkgs, err := packages.Load(pl.Config, pkgName)
	if err != nil {
		fmt.Println(err)
		return
	}
	var continueFlag bool
	for _, pkg := range pkgs {
		fmt.Println(pkg.PkgPath)
		continueFlag = false
		for _, filter := range pl.PkgFilter {
			if !filter(pkg.PkgPath) {
				continueFlag = true
			}
		}
		if continueFlag {
			continue
		}

		if _, ok := pkgMap[pkg.ID]; !ok {
			pkgMap[pkg.ID] = pkg
		} else {
			continue
		}

		for importPkgName, _ := range pkg.Imports {
			pl.RecursionParsePkg(importPkgName, pkgMap)
		}
	}
}

func main() {
	PACKAGE_NAME := ""

	loadMode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedExportsFile |
		packages.NeedTypes |
		packages.NeedTypesInfo |
		packages.NeedTypesSizes |
		packages.NeedSyntax |
		packages.NeedModule

	cfg := &packages.Config{
		Mode:       loadMode,
		BuildFlags: build.Default.BuildTags,
		Dir:        "",
	}
	pl := NewPackageLoader()
	pl.Config = cfg
	pl.PkgFilter = []func(pkgPath string) bool{
		func(pkgPath string) bool {
			if strings.HasPrefix(pkgPath, PACKAGE_NAME) {
				return true
			}
			return false
		},
	}
	pkgMap := map[string]*packages.Package{}
	pl.RecursionParsePkg(PACKAGE_NAME, pkgMap)

	fs := token.NewFileSet()
	for pkgID, pkg := range pkgMap {
		fmt.Println("package name:", pkgID)
		for _, file := range pkg.Syntax {

			fmt.Println("file name:", file.Name.Name)
			ast.Print(fs, file)
			return
			// ast.Inspect(file, func(n ast.Node) bool {
			// 	if parsedSelectorExpr, ok := n.(*ast.SelectorExpr); ok {
			// 		if parsedIndet, ok := parsedSelectorExpr.X.(*ast.Ident); ok {
			// 			if parsedIndet.Name == "logging" && parsedSelectorExpr.Sel.Name == "Infof" {
			// 				ast.Print(fs, parsedIndet)
			// 				ast.Print(fs, parsedSelectorExpr.Sel)
			// 			}
			// 		}
			// 	}
			// 	return true
			// })
		}
	}
}
