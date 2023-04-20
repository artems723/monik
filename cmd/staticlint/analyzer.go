// Package analyzer checks for os.Exit in main function of main package
package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer is an analyzer that checks for os.Exit in main function of main package
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit in main function and main package",
	Run:  run,
}

// This function is called for main package and main function
// It checks if there is a call to os.Exit in main function
// If there is a call to os.Exit, it reports an error
// It iterates through all top-level declarations in main function
// FuncDecl->ExprStmt->CallExpr->SelectorExpr
func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		// Check if file is main
		if f.Name.Name == "main" {
			// Iterate all top-level declarations
			for _, decl := range f.Decls {
				// If it is a Func, check if it is main
				if funcDecl, ok := decl.(*ast.FuncDecl); ok {
					// If it is main, iterate all elements in main
					if funcDecl.Name.Name == "main" {
						// Iterate all elements in main
						for _, l := range funcDecl.Body.List {
							// Check elements is a ExprStmt
							switch exprStmt := l.(type) {
							case *ast.ExprStmt:
								// Check if ExprStmt is a CallExpr
								if call, ok := exprStmt.X.(*ast.CallExpr); ok {
									// Check if CallExpr is a SelectorExpr
									if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
										// Final get expression
										if first, ok := fun.X.(*ast.Ident); ok {
											result := first.Name + "." + fun.Sel.Name
											if result == "os.Exit" {
												pass.Reportf(first.NamePos, "call os.Exit in main() in package main %+v,%+v,%+v", f.Package, f.Name, f.Name.Name)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil, nil
}
