package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
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

func (pl *packageLoader) RecursionParsePkg(pkg *packages.Package, pkgName string, pkgMap map[string]*packages.Package) {
	var returnTag bool
	for _, filter := range pl.PkgFilter {
		if !filter(pkg.PkgPath) {
			returnTag = true
		}
	}
	if returnTag {
		fmt.Println(pkg.PkgPath)
		return
	}

	if _, ok := pkgMap[pkg.ID]; !ok {
		pkgMap[pkg.ID] = pkg
	} else {
		fmt.Println(pkg.PkgPath)
		return
	}

	for _, imp := range pkg.Imports {
		pl.RecursionParsePkg(imp, imp.ID, pkgMap)
	}
}

func main() {
	PACKAGE_NAME := "github.com/noobone/go-ast-book/demo"

	loadMode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		// packages.NeedDeps |
		// packages.NeedExportsFile |
		// packages.NeedTypes |
		// packages.NeedTypesInfo |
		// packages.NeedTypesSizes |
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

	pkgs, err := packages.Load(pl.Config, PACKAGE_NAME)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, pkg := range pkgs {
		pl.RecursionParsePkg(pkg, PACKAGE_NAME, pkgMap)
	}

	for pkgID, pkg := range pkgMap {
		fmt.Println("package name:", pkgID)
		for _, file := range pkg.Syntax {
			ast.Print(pkg.Fset, file)
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

// func main() {
// 	loadMode := packages.NeedName |
// 		packages.NeedFiles |
// 		packages.NeedCompiledGoFiles |
// 		packages.NeedImports |
// 		packages.NeedSyntax |
// 		packages.NeedModule

// 	cfg := &packages.Config{
// 		Mode: loadMode,
// 		// BuildFlags: build.Default.BuildTags,
// 	}
// 	pkgs, err := packages.Load(cfg, "github.com/noobone/go-ast-book")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	demo.Demo(1, "1")

// 	bPkgs, _ := json.Marshal(pkgs)
// 	fmt.Println(string(bPkgs))

// 	// fs := token.NewFileSet()
// 	for _, pkg := range pkgs {
// 		for _, file := range pkg.Syntax {
// 			ast.Inspect(file, func(n ast.Node) bool {
// 				if parsedSelectorExpr, ok := n.(*ast.SelectorExpr); ok {
// 					if parsedIndet, ok := parsedSelectorExpr.X.(*ast.Ident); ok {
// 						if parsedIndet.Name == "log" && parsedSelectorExpr.Sel.Name == "Print" {
// 							log.Print("妙啊")
// 						}
// 					}
// 				}
// 				return true
// 			})
// 		}
// 	}
// }

func main1() {
	fs := token.NewFileSet()
	pkgs, _ := parser.ParseDir(fs, "demo", nil, parser.AllErrors)
	// ast.Print(fs, pkgs)
	ast.Print(fs, pkgs["demo"].Files["demo\\demo.go"])
	for pkgName, pkg := range pkgs {
		fmt.Println(pkgName)
		for fileName, file := range pkg.Files {
			fmt.Println(fileName)

			ast.Inspect(file, func(n ast.Node) bool {
				switch res := n.(type) {
				// Find Return Statements
				case *ast.ReturnStmt:
					fmt.Printf("return statement found on line %d\n", fs.Position(res.Pos()).Line)
					return true
				// Find Functions
				case *ast.FuncDecl:
					var exported string
					if res.Name.IsExported() {
						exported = "exported "
					}
					fmt.Printf("%sfunction declaration found on line %d: %s\n", exported, fs.Position(res.Pos()).Line, res.Name.Name)
					return true
				case *ast.SelectorExpr:
					if parsedIndet, ok := res.X.(*ast.Ident); ok {
						sel := res.Sel
						fmt.Printf("pkg name: %s, pkg imported position: %s\n", parsedIndet.Name, fs.Position(parsedIndet.NamePos).String())
						fmt.Printf("select func name: %s, select func called position: %s\n", sel.Name, fs.Position(sel.NamePos).String())
					}
					return true
				default:
					return true
				}
			})
		}
	}
}

// func main4() {
// 	fs := token.NewFileSet()
// 	pkgs, _ := parser.ParseDir(fs, "demo", nil, parser.AllErrors)
// 	ast.Print(fs, pkgs)
// 	ast.Print(fs, pkgs["demo"].Files["demo\\demo.go"])
// 	for pkgName, pkg := range pkgs {
// 		fmt.Println(pkgName)
// 		for fileName, file := range pkg.Files {
// 			fmt.Println(fileName)
// 			for _, decl := range file.Decls {
// 				if parsedDecl, ok := decl.(*ast.FuncDecl); ok {
// 					fmt.Println(parsedDecl.Name)
// 					for _, stmt := range parsedDecl.Body.List {
// 						// fmt.Println(stmt.Pos())
// 						if parsedStmt, ok := stmt.(*ast.ExprStmt); ok {
// 							fmt.Println("parsedStmt")
// 							if parsedExpr, ok := parsedStmt.X.(*ast.CallExpr); ok {
// 								fmt.Println("parsedExpr")
// 								if parsedSelectorExpr, ok := parsedExpr.Fun.(*ast.SelectorExpr); ok {
// 									if parsedIndet, ok := parsedSelectorExpr.X.(*ast.Ident); ok {
// 										sel := parsedSelectorExpr.Sel
// 										fmt.Printf("pkg name: %s, pkg imported position: %s\n", parsedIndet.Name, fs.Position(parsedIndet.NamePos).String())
// 										fmt.Printf("select func name: %s, select func called position: %s\n", sel.Name, fs.Position(sel.NamePos).String())
// 									}
// 								}
// 							}

// 						}
// 					}
// 				}

// 			}
// 		}
// 	}
// }

// func main3() {
// 	expr, _ := parser.ParseExpr(`1+2*3`)
// 	fmt.Println(Eval(expr))
// }

// func Eval(exp ast.Expr) float64 {
// 	switch exp := exp.(type) {
// 	case *ast.BinaryExpr:
// 		return EvalBinaryExpr(exp)
// 	case *ast.BasicLit:
// 		f, _ := strconv.ParseFloat(exp.Value, 64)
// 		return f
// 	}
// 	return 0
// }

// func EvalBinaryExpr(exp *ast.BinaryExpr) float64 {
// 	switch exp.Op {
// 	case token.ADD:
// 		return Eval(exp.X) + Eval(exp.Y)
// 	case token.MUL:
// 		return Eval(exp.X) * Eval(exp.Y)
// 	}
// 	return 0
// }

// func main22() {
// 	expr, _ := parser.ParseExpr(`9527`)
// 	ast.Print(nil, expr)
// }

// func main21() {
// 	var lit9527 = &ast.BasicLit{
// 		Kind:  token.FLOAT,
// 		Value: "9527",
// 	}
// 	ast.Print(nil, lit9527)
// }

// func main1() {
// 	var src = []byte(`println("你好，世界")`)

// 	var fset = token.NewFileSet()
// 	var file = fset.AddFile("hello.go", fset.Base(), len(src))

// 	var s scanner.Scanner
// 	s.Init(file, src, nil, scanner.ScanComments)

// 	for {
// 		pos, tok, lit := s.Scan()
// 		if tok == token.EOF {
// 			break
// 		}
// 		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
// 	}
// }
