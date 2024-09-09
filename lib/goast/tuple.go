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

type Tuple struct {
	Name            string
	Fields          []*ast.Field
	Kind            ObjectType
	mapFieldsToType map[string]ObjectType
	identifier      []AST
}

func (t *Tuple) GetObjectType() ObjectType {
	return TupleType
}
func (t *Tuple) GetName() string {
	return t.Name
}
func (t *Tuple) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(t.Name)
}
func (t *Tuple) Under(_ AST)              {}
func (*Tuple) SetKind(o ObjectType)       {}
func (*Tuple) SetStructKind(o ObjectType) {}
func (t *Tuple) AST(above AST) (decl []ast.Decl) {
	name := t.GetTitle()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
	}
	var fields []*ast.Field
	for _, field := range t.Fields {
		if len(field.Names) < 1 {
			continue
		}

		var nameType ast.Expr
		nameType = field.Type
		if ident, ok := nameType.(*ast.Ident); ok {
			if strings.ToLower(ident.String()) == "any" {
				nameType = ast.NewIdent("Value")
			}

			if ident.String() == "String" {
				nameType = ast.NewIdent("Pstring")
			}
		}

		fname := cases.Title(language.English, cases.NoLower).String(field.Names[0].String())
		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(fname)},
			Type:  nameType,
		})
	}
	decl = append(decl, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: fields,
					},
				},
			},
		},
	})
	var smallFields []*ast.Field
	for _, field := range fields {
		if len(field.Names) < 1 {
			continue
		}
		fname := fmt.Sprintf("%s%s", strings.ToLower(field.Names[0].Name[0:1]), field.Names[0].Name[1:])
		switch fname {
		case "interface":
			fallthrough
		case "any":
			fname = fmt.Sprintf("_%s", fname)
		}
		smallFields = append(smallFields, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(fname)},
			Type:  field.Type,
		})
	}
	var smallValues []ast.Expr
	for _, field := range fields {
		if len(field.Names) < 1 {
			continue
		}
		fname := fmt.Sprintf("%s%s", strings.ToLower(field.Names[0].Name[0:1]), field.Names[0].Name[1:])
		switch fname {
		case "interface":
			fallthrough
		case "any":
			fname = fmt.Sprintf("_%s", fname)
		}
		smallValues = append(smallValues, &ast.KeyValueExpr{
			Key:   field.Names[0],
			Value: ast.NewIdent(fname),
		})
	}

	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("New%s", name)),
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Params: &ast.FieldList{
					List: smallFields,
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  &ast.StarExpr{X: ast.NewIdent(name)},
						},
					},
				},
			},
			Body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent(name),
								Elts: smallValues,
							},
						},
					},
				},
			},
			},
		},
	)

	sname := &ast.StarExpr{X: ast.NewIdent(name)}
	if above != nil {
		decl = append(decl,
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  sname,
						}},
				},
				Name: ast.NewIdent(fmt.Sprintf("Is%s", above.GetName())),
				Type: &ast.FuncType{
					Func:   token.Pos(token.FUNC),
					Params: &ast.FieldList{},
				},
				Body: &ast.BlockStmt{},
			})
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
				},
			},
		},
	)

	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "ToPreserves")),
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("d")},
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
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.CompositeLit{
									Type: ast.NewIdent("Sequence"),
									Elts: keyValues,
								},
							},
						},
					}},
			},
		},
	)
	return
}
