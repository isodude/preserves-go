package goast

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type StructTuple struct {
	Name            string
	Fields          []*ast.Field
	Kind            ObjectType
	mapFieldsToType map[string]ObjectType
	identifier      []AST
}

func (t *StructTuple) GetObjectType() ObjectType {
	return StructTupleType
}
func (t *StructTuple) GetName() string {
	return t.Name
}
func (t *StructTuple) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(t.Name)
}
func (t *StructTuple) Under(_ AST)              {}
func (*StructTuple) SetKind(o ObjectType)       {}
func (*StructTuple) SetStructKind(o ObjectType) {}
func (t *StructTuple) AST(above AST) (decl []ast.Decl) {
	name := t.GetTitle()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
	}
	sname := &ast.StarExpr{X: ast.NewIdent(name)}
	var returnElts []ast.Expr
	var firstStmt, curStmt ast.Stmt

	keyValues := []ast.Expr{
		&ast.KeyValueExpr{
			Key: ast.NewIdent("Key"),
			Value: &ast.CallExpr{
				Fun:  ast.NewIdent("NewSymbol"),
				Args: []ast.Expr{ast.NewIdent(strconv.Quote(t.GetName()))},
			},
		},
	}
	if len(t.identifier) > 0 {
		if stmter, ok := t.identifier[0].(ToStmt); ok {
			keyValues = stmter.ToStmt(ast.NewIdent("Key"))
		}
	}
	fieldValues := []ast.Expr{}
	iterateValues := []ast.Stmt{}
	for i, field := range t.Fields {
		var arrayType bool
		if len(field.Names) != 1 {
			continue
		}

		var stmt ast.Stmt
		var callExprFun *ast.Ident
		fieldName := field.Names[0].String()
		fieldName = fmt.Sprintf("%s%s", strings.ToUpper(string(fieldName[0])), fieldName[1:])
		var fieldType string
		switch t := field.Type.(type) {
		case *ast.ArrayType:
			switch u := t.Elt.(type) {
			case *ast.Ident:
				fieldType = u.String()
				arrayType = true
			default:
				panic(fmt.Sprintf("beep: %s", reflect.TypeOf(field.Type)))
			}
		case *ast.Ident:
			fieldType = t.String()
		default:
			panic(fmt.Sprintf("beep: %s", reflect.TypeOf(field.Type)))
		}
		if fieldType == "String" {
			fieldType = "Pstring"
		}

		var toCallExprFun ast.Expr

		toCallExprFun = ToMapToCallFunc(t.mapFieldsToType, fieldType)

		if strings.ToLower(fieldType) == "any" && !arrayType {
			stmt = &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X:  &ast.CompositeLit{Type: ast.NewIdent(name), Elts: []ast.Expr{&ast.KeyValueExpr{Key: ast.NewIdent(fieldName), Value: &ast.IndexExpr{Index: ast.NewIdent("0"), X: &ast.SelectorExpr{X: ast.NewIdent("rec"), Sel: ast.NewIdent("Fields")}}}}},
					},
				},
			}
			fieldValues = append(fieldValues, &ast.SelectorExpr{
				X:   ast.NewIdent(`s`),
				Sel: ast.NewIdent(fieldName),
			})
			if firstStmt == nil {
				firstStmt = stmt
				curStmt = stmt
			} else {
				curStmt.(*ast.IfStmt).Body.List = append(curStmt.(*ast.IfStmt).Body.List, stmt)
				curStmt = stmt
			}
			continue
		}
		o, ok := t.mapFieldsToType[strings.ToLower(fieldType)]
		if !ok {
			continue
		}

		if callExprFun = MapToCallFunc(t.mapFieldsToType, fieldType); callExprFun == nil {
			continue
		}

		varName := ast.NewIdent(fmt.Sprintf("p%d", i))

		var dVarName ast.Expr
		if o == UnionInterfaceObjectType {
			dVarName = varName
		} else if o == ValueType {
			dVarName = varName
		} else {
			dVarName = &ast.StarExpr{X: varName}
		}
		returnElts = append(returnElts, &ast.KeyValueExpr{
			Key:   ast.NewIdent(fieldName),
			Value: dVarName,
		})
		if arrayType {
			fieldValues = append(fieldValues, varName)
			d := &ast.DeclStmt{Decl: &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{&ast.TypeSpec{Assign: token.Pos(token.ASSIGN), Name: varName, Type: &ast.UnaryExpr{Op: token.AND, X: &ast.CompositeLit{Type: ast.NewIdent("Sequence")}}}}}}
			r := &ast.RangeStmt{Key: ast.NewIdent("_"), Value: ast.NewIdent("k"), Tok: token.DEFINE, X: &ast.SelectorExpr{
				X:   ast.NewIdent(`s`),
				Sel: ast.NewIdent(fieldName),
			}, Body: &ast.BlockStmt{List: []ast.Stmt{&ast.AssignStmt{Tok: token.ASSIGN, Lhs: []ast.Expr{&ast.StarExpr{X: varName}}, Rhs: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("append"), Args: []ast.Expr{&ast.StarExpr{X: varName}, &ast.CallExpr{Fun: toCallExprFun, Args: []ast.Expr{ast.NewIdent("k")}}}}}}}}}
			iterateValues = append(iterateValues, d, r)
			stmt = &ast.IfStmt{
				Init: &ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						ast.NewIdent("seq"),
						ast.NewIdent("ok"),
					},
					Rhs: []ast.Expr{
						&ast.TypeAssertExpr{
							X: &ast.IndexExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("rec"),
									Sel: ast.NewIdent("Fields"),
								},
								Index: ast.NewIdent(strconv.Itoa(i)),
							},
							Type: &ast.StarExpr{X: ast.NewIdent("Sequence")},
						},
					},
				},
				Cond: ast.NewIdent("ok"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.DeclStmt{
							Decl: &ast.GenDecl{
								Tok: token.VAR,
								Specs: []ast.Spec{
									&ast.TypeSpec{
										Name: varName,
										Type: &ast.ArrayType{
											Elt: ast.NewIdent(fieldType),
										},
									},
								},
							},
						},
						&ast.RangeStmt{
							Tok:   token.DEFINE,
							Key:   ast.NewIdent("_"),
							Value: ast.NewIdent("item"),
							X:     &ast.StarExpr{X: ast.NewIdent("seq")},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{

									&ast.IfStmt{
										Init: &ast.AssignStmt{
											Tok: token.DEFINE,
											Lhs: []ast.Expr{ast.NewIdent("itemParsed")},
											Rhs: []ast.Expr{
												&ast.CallExpr{
													Fun: callExprFun,
													Args: []ast.Expr{
														ast.NewIdent("item"),
													},
												},
											},
										},
										Cond: &ast.BinaryExpr{
											Op: token.NEQ,
											X:  ast.NewIdent("itemParsed"),
											Y:  ast.NewIdent("nil"),
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.AssignStmt{
													Tok: token.ASSIGN,
													Lhs: []ast.Expr{varName},
													Rhs: []ast.Expr{
														&ast.CallExpr{
															Fun: ast.NewIdent("append"),
															Args: []ast.Expr{
																varName,
																ast.NewIdent("itemParsed"),
															},
														},
													},
												},
											},
										},
										Else: &ast.ReturnStmt{
											Results: []ast.Expr{ast.NewIdent("nil")},
										},
									},
								},
							},
						},
					},
				},
			}
		} else {
			fieldValues = append(fieldValues, &ast.CallExpr{Fun: toCallExprFun, Args: []ast.Expr{&ast.SelectorExpr{
				X:   ast.NewIdent(`s`),
				Sel: ast.NewIdent(fieldName),
			}}})
			stmt = &ast.IfStmt{
				Init: &ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{varName},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: callExprFun,
							Args: []ast.Expr{
								&ast.IndexExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("rec"),
										Sel: ast.NewIdent("Fields"),
									},
									Index: ast.NewIdent(strconv.Itoa(i)),
								},
							},
						},
					},
				},
				Cond: &ast.BinaryExpr{
					Op: token.NEQ,
					X:  varName,
					Y:  ast.NewIdent("nil"),
				},
				Body: &ast.BlockStmt{},
			}
		}
		if firstStmt == nil {
			firstStmt = stmt
			curStmt = stmt
		} else {
			curStmt.(*ast.IfStmt).Body.List = append(curStmt.(*ast.IfStmt).Body.List, stmt)
			curStmt = stmt
		}
		if len(t.Fields) == i+1 {
			returnStmt := &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X: &ast.CompositeLit{
							Type: ast.NewIdent(name),
							Elts: returnElts,
						},
					},
				},
			}

			stmt.(*ast.IfStmt).Body.List = append(stmt.(*ast.IfStmt).Body.List, returnStmt)
		}

	}

	keyValues = append(keyValues, &ast.KeyValueExpr{Key: ast.NewIdent("Fields"), Value: &ast.CompositeLit{
		Type: &ast.ArrayType{
			Elt: ast.NewIdent("Value"),
		},
		Elts: fieldValues,
	}})

	body := &ast.BlockStmt{}
	if firstStmt != nil {
		body.List = []ast.Stmt{firstStmt}
	} else {
		body.List = []ast.Stmt{&ast.ReturnStmt{
			Results: []ast.Expr{ast.NewIdent("hello")},
		}}
	}
	var ifStmt ast.Stmt
	ifStmt = &ast.IfStmt{
		Init: &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				ast.NewIdent("sym"),
				ast.NewIdent("ok"),
			},
			Rhs: []ast.Expr{
				&ast.TypeAssertExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("rec"),
						Sel: ast.NewIdent("Key"),
					},
					Type: &ast.StarExpr{X: ast.NewIdent("Symbol")},
				},
			},
		},
		Cond: &ast.BinaryExpr{
			Op: token.LAND,
			X:  ast.NewIdent("ok"),
			Y: &ast.BinaryExpr{
				Op: token.EQL,
				X: &ast.CallExpr{
					Fun: ast.NewIdent("sym.String"),
				},
				Y: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", t.Name),
				},
			},
		},
		Body: body,
	}
	if len(t.identifier) > 0 {
		if stmter, ok := t.identifier[0].(Stmt); ok {
			ifStmt = stmter.Stmt(&ast.SelectorExpr{
				X:   ast.NewIdent("rec"),
				Sel: ast.NewIdent("Key"),
			}, body.List)
		}
	}
	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "FromPreserves")),
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("value")},
							Type:  ast.NewIdent("Value"),
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  sname,
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.IfStmt{
						Init: &ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{
								ast.NewIdent("rec"),
								ast.NewIdent("ok"),
							},
							Rhs: []ast.Expr{
								&ast.TypeAssertExpr{
									X:    ast.NewIdent("value"),
									Type: &ast.StarExpr{X: ast.NewIdent("Record")},
								},
							},
						},
						Cond: &ast.BinaryExpr{
							Op: token.LAND,
							X:  ast.NewIdent("ok"),
							Y: &ast.BinaryExpr{
								Op: token.EQL,
								X: &ast.CallExpr{
									Fun: ast.NewIdent("len"),
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("rec"),
											Sel: ast.NewIdent("Fields"),
										},
									},
								},
								Y: ast.NewIdent(strconv.Itoa(len(t.Fields))),
							},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								ifStmt,
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{ast.NewIdent("nil")},
					},
				},
			},
		},
	)
	/*
	   func RefToPreserves(r Ref) Value {
	   	return &Record{Key: NewSymbol("ref"), Fields: []Value{ModulePathToPreserves(r.Module), SymbolToPreserves(r.Name)}}
	   }
	*/

	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "ToPreserves")),
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("s")},
							Type:  ast.NewIdent(name),
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  ast.NewIdent("Value"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: append(iterateValues,
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.CompositeLit{
									Type: ast.NewIdent("Record"),
									Elts: keyValues,
								},
							},
						},
					}),
			},
		},
	)
	return
}
