package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/go.tools/astutil"
)

func warpReader(name ast.Expr, pos token.Pos) ast.Stmt {
	call := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "warped"},
			Sel: &ast.Ident{Name: "Reader"},
		},
		Args: []ast.Expr{name},
	}

	assign := &ast.AssignStmt{
		Lhs:    []ast.Expr{name},
		TokPos: pos,
		Tok:    token.ASSIGN,
		Rhs:    []ast.Expr{call},
	}

	return assign
}

func main() {
	wd, err := os.Getwd()
	if len(os.Args) < 2 {
		fmt.Println("usage: warp [packages]. warp will OVERWRITE your packages. Use source control.")
		os.Exit(2)
	}
	if err != nil {
		fatal(err)
	}
	var files []string
	for _, path := range os.Args[1:] {
		pkg, err := build.Import(path, wd, 0)
		if err != nil {
			fatal(err)
		}
		for _, file := range pkg.GoFiles {
			files = append(files, filepath.Join(pkg.Dir, file))
		}
	}

	for _, fname := range files {
		if fname == "--" {
			continue
		}
		if !strings.HasSuffix(fname, ".go") || strings.HasSuffix(fname, "_test.go") {
			continue
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
		if err != nil {
			continue
		}

		var altered bool

		for _, d := range f.Decls {
			fun, ok := d.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if !ast.IsExported(fun.Name.String()) {
				continue
			}

			// Assume that "io" is imported as "io" if at all.
			// Spares a long diversion through go.types and/or ugly import processing.
			for _, arg := range fun.Type.Params.List {
				sel, ok := arg.Type.(*ast.SelectorExpr)
				if !ok || sel.Sel.String() != "Reader" || len(arg.Names) != 1 {
					continue
				}
				id, ok := sel.X.(*ast.Ident)
				if !ok || id.Name != "io" {
					continue
				}
				name := arg.Names[0]
				assign := warpReader(name, fun.Body.Pos())
				body := []ast.Stmt{assign}
				fun.Body.List = append(body, fun.Body.List...)
				altered = true
			}
		}

		if altered {
			astutil.AddImport(fset, f, "github.com/josharian/warp/warped")

			c, err := os.Create(fname)
			if err != nil {
				fatal(err)
			}

			printer.Fprint(c, fset, f)
		}
	}
}

func fatal(msg interface{}) {
	fmt.Println(msg)
	os.Exit(1)
}
