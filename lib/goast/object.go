package goast

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
)

//	type Object struct {
//	  Field Field
//	  Interface Interface
//	}
func objectTypeSpec(prefix *title, name *title, fields []*Field) ast.Decl {
	var (
		list []*ast.Field
	)

	decl := &ast.GenDecl{Tok: token.TYPE, Specs: []ast.Spec{&ast.TypeSpec{
		Name: ast.NewIdent(name.GetPrefixTitle(prefix)),
		Type: &ast.StructType{Fields: &ast.FieldList{List: list}},
	}}}

	for _, field := range fields {
		list = append(list, &ast.Field{Names: []*ast.Ident{field.Ident(nil, false)}, Type: field.Expr()})
	}

	return decl
}

//	func NewObject(field Field, _interface Interface) {
//	  return &Object{Field: field, Interface: _interface }
//	}
func objectFuncDeclNew(prefix *title, name *title, fields []*Field) ast.Decl {
	var (
		params []*ast.Field
		elts   []ast.Expr
	)

	funcDecl := &ast.FuncDecl{Name: name.Ident(prefix.PrefixTitle("New"), false),
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: params},
			Results: &ast.FieldList{List: []*ast.Field{
				{Names: []*ast.Ident{}, Type: &ast.StarExpr{X: name.Ident(prefix, false)}},
			},
			}},
		Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{&ast.UnaryExpr{
			Op: token.AND,
			X:  &ast.CompositeLit{Type: name.Ident(nil, false), Elts: elts},
		}}}}},
	}

	for _, field := range fields {
		f := &ast.Field{Type: field.Expr()}
		kv := &ast.KeyValueExpr{
			Key: field.Ident(nil, false),
		}
		if slices.Contains([]string{"interface", "any"}, field.GetFieldTypeName()) {
			f.Names = []*ast.Ident{field.Ident(&title{name: "_"}, true)}
			kv.Value = field.Ident(&title{name: "_"}, true)
		} else {
			f.Names = []*ast.Ident{field.Ident(nil, true)}
			kv.Value = field.Ident(nil, true)
		}
		params = append(params, f)
		elts = append(elts, kv)
	}

	return funcDecl
}

// func (*Object) IsInterface()
func objectFuncDeclIs(prefix *title, name *title, target *title) ast.Decl {
	return &ast.FuncDecl{
		Body: &ast.BlockStmt{},
		Name: target.Ident(&title{name: "Is"}, false),
		Recv: &ast.FieldList{List: []*ast.Field{
			{Names: []*ast.Ident{}, Type: &ast.StarExpr{X: name.Ident(prefix, false)}},
		}},
		Type: &ast.FuncType{Params: &ast.FieldList{}},
	}
}

// func ObjectFromPreserves(value Value) *Object { }
func objectFuncDeclFromPreserves(prefix *title, name *title, fields []*Field) ast.Decl {
	var (
		list []ast.Stmt
	)
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("%s%s", name.GetPrefixTitle(prefix), "FromPreserves")),
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: []*ast.Field{
				{Names: []*ast.Ident{ast.NewIdent("value")}, Type: ast.NewIdent("Value")},
			}},
			Results: &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{}, Type: &ast.StarExpr{X: ast.NewIdent(name.GetPrefixTitle(prefix))}}}},
		},
		Body: &ast.BlockStmt{List: list},
	}
	return funcDecl
}

// func ObjectToPreserves(d Object) Value { }
func objectFuncDeclToPreserves(prefix *title, name *title, fields []*Field) ast.Decl {
	var (
		list []ast.Stmt
	)

	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("%s%s", name.GetPrefixTitle(prefix), "ToPreserves")),
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: []*ast.Field{
				{Names: []*ast.Ident{ast.NewIdent("d")}, Type: ast.NewIdent(name.GetPrefixTitle(prefix))},
			}},
			Results: &ast.FieldList{List: []*ast.Field{
				{Names: []*ast.Ident{}, Type: ast.NewIdent("Value")}},
			}},
		Body: &ast.BlockStmt{List: list},
	}

	return funcDecl
}

// return &Sequence{}
func returnSequence(elts []ast.Expr) *ast.ReturnStmt {
	return &ast.ReturnStmt{Results: []ast.Expr{&ast.UnaryExpr{
		Op: token.AND,
		X:  &ast.CompositeLit{Type: ast.NewIdent("Sequence"), Elts: elts},
	}}}
}

/*
			&ast.IfStmt{
				Init: &ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						ast.NewIdent("seq"),
						ast.NewIdent("ok"),
					},
					Rhs: []ast.Expr{
						&ast.TypeAssertExpr{
							X:    ast.NewIdent("value"),
							Type: &ast.StarExpr{X: ast.NewIdent("Sequence")},
						},
					},
				},
				Cond: &ast.BinaryExpr{
					Op: token.LAND,
					X:  ast.NewIdent("ok"),
					Y: &ast.BinaryExpr{
						Op: token.EQL,
						X: &ast.CallExpr{
							Fun:  ast.NewIdent("len"),
							Args: []ast.Expr{&ast.StarExpr{X: ast.NewIdent("seq")}},
						},
						Y: ast.NewIdent(strconv.Itoa(len(t.Fields))),
					},
				},
				Body: body,
			},
			&ast.ReturnStmt{
				Results: []ast.Expr{ast.NewIdent("nil")},
			},
		}},
	}
}

	var returnElts []ast.Expr
	var firstStmt, curStmt ast.Stmt
	var keyValues []ast.Expr
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
		if strings.ToLower(fieldType) == "any" && !arrayType {
			stmt = &ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X:  &ast.CompositeLit{Type: ast.NewIdent(name), Elts: []ast.Expr{&ast.KeyValueExpr{Key: ast.NewIdent(fieldName), Value: ast.NewIdent("value")}}},
					},
				},
			}
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
		varNameItemParsed := ast.NewIdent("itemParsed")
		var dVarName, dVarNameItemParsed ast.Expr
		if o == UnionInterfaceObjectType {
			dVarName = varName
			dVarNameItemParsed = varName
		} else if o == InterfaceObjectType {
			dVarName = varName
			dVarNameItemParsed = varName
		} else if o == MapObjectType {
			dVarName = varName
			dVarNameItemParsed = varName
		} else {
			dVarName = &ast.StarExpr{X: varName}
			dVarNameItemParsed = &ast.StarExpr{X: varNameItemParsed}
		}

		returnElts = append(returnElts, &ast.KeyValueExpr{
			Key:   ast.NewIdent(fieldName),
			Value: dVarName,
		})

		if arrayType {
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
																dVarNameItemParsed,
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
			keyValues = append(keyValues, &ast.CallExpr{Fun: ast.NewIdent(fmt.Sprintf("%sToPreserves", fieldType)), Args: []ast.Expr{&ast.SelectorExpr{X: ast.NewIdent("d"), Sel: ast.NewIdent(fieldName)}}})
			stmt = &ast.IfStmt{
				Init: &ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{varName},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: callExprFun,
							Args: []ast.Expr{
								&ast.IndexExpr{
									X: &ast.ParenExpr{
										X: &ast.StarExpr{X: ast.NewIdent("seq")},
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
	body := &ast.BlockStmt{}
	if firstStmt != nil {
		body.List = []ast.Stmt{firstStmt}
	} else {
		body.List = []ast.Stmt{&ast.ReturnStmt{
			Results: []ast.Expr{ast.NewIdent("hello")},
		}}
	}
	decl = append(decl,
		,
	)

	decl = append(decl,
*/
