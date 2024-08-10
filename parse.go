/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
)

func walkAndParse(dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".go" {
			fmt.Printf("Processing file: %s\n", path)

			err := parseFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"Error parsing file %s: %v\n", path, err)
			}
		}

		return nil
	})
}

func parseFile(filename string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(node, listen)

	return nil
}

func listen(n ast.Node) bool {
	// Look for function declarations
	if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "NewConfig" {
		// Look for return statements
		ast.Inspect(fn, func(n ast.Node) bool {
			if ret, ok := n.(*ast.ReturnStmt); ok {
				for _, retVal := range ret.Results {
					if compositeLit, ok := retVal.(*ast.CompositeLit); ok {
						for _, elt := range compositeLit.Elts {
							if keyValue, ok := elt.(*ast.KeyValueExpr); ok {
								if key, ok := keyValue.Key.(*ast.Ident); ok && key.Name == "ListenerRoutes" {
									if arrayLit, ok := keyValue.Value.(*ast.CompositeLit); ok {
										for _, route := range arrayLit.Elts {
											if basicLit, ok := route.(*ast.BasicLit); ok {
												fmt.Println(basicLit.Value)
											}
										}
									}
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return true
}
