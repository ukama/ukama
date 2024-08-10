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
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func walkAndParse(dir string, out io.Writer) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".go" {
			err := parseFile(path, out)

			if err != nil {
				fmt.Fprintf(os.Stderr,
					"Error parsing file %s: %v\n", path, err)
			}
		}
		return nil
	})
}

func parseFile(filename string, output io.Writer) error {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	listen(node, output)

	return nil
}

func listen(node *ast.File, output io.Writer) {
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "NewConfig" {
			ast.Inspect(fn, func(n ast.Node) bool {
				if ret, ok := n.(*ast.ReturnStmt); ok {
					for _, retVal := range ret.Results {
						if unaryExpr, ok := retVal.(*ast.UnaryExpr); ok {
							if compositeLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
								for _, elt := range compositeLit.Elts {
									if keyValue, ok := elt.(*ast.KeyValueExpr); ok {
										if key, ok := keyValue.Key.(*ast.Ident); ok {
											if key.Name == "MsgClient" {
												if unaryExpr2, ok := keyValue.Value.(*ast.UnaryExpr); ok {
													if compositeLit2, ok := unaryExpr2.X.(*ast.CompositeLit); ok {
														for _, elt := range compositeLit2.Elts {
															if keyValue, ok := elt.(*ast.KeyValueExpr); ok {
																if key, ok := keyValue.Key.(*ast.Ident); ok {
																	if key.Name == "ListenerRoutes" {
																		if arrayLit, ok := keyValue.Value.(*ast.CompositeLit); ok {
																			for _, route := range arrayLit.Elts {
																				if basicLit, ok := route.(*ast.BasicLit); ok {
																					fmt.Fprintln(output, basicLit.Value)
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
	})
}
