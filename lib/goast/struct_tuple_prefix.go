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

type TuplePrefix struct {
	Name            string
	Fields          []*ast.Field
	Kind            ObjectType
	mapFieldsToType map[string]ObjectType
	identifier      []AST
}

func (t *TuplePrefix) GetObjectType() ObjectType {
	return StructTuplePrefixType
}
func (t *TuplePrefix) GetName() string {
	return t.Name
}
func (t *TuplePrefix) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(t.Name)
}
func (t *TuplePrefix) Under(_ AST)              {}
func (*TuplePrefix) SetKind(o ObjectType)       {}
func (*TuplePrefix) SetStructKind(o ObjectType) {}
func (t *TuplePrefix) AST(above AST) (decl []ast.Decl) {
	name := t.GetTitle()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
	}
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
	sname := &ast.StarExpr{X: ast.NewIdent(name)}
	var fieldType string
	switch _t := t.Fields[0].Type.(type) {
	case *ast.ArrayType:
		switch u := _t.Elt.(type) {
		case *ast.Ident:
			fieldType = u.String()
		default:
			panic(fmt.Sprintf("beep: %s", reflect.TypeOf(t.Fields[0].Type)))
		}
	case *ast.Ident:
		fieldType = _t.String()
	default:
		panic(fmt.Sprintf("beep: %s", reflect.TypeOf(t.Fields[0].Type)))
	}

	if fieldType == "String" {
		fieldType = "Pstring"
	}
	var callExprFun *ast.Ident
	fieldTypeName := ast.NewIdent(fieldType)
	o, ok := t.mapFieldsToType[strings.ToLower(fieldType)]
	if !ok {
		return
	}

	if callExprFun = MapToCallFunc(t.mapFieldsToType, fieldType); callExprFun == nil {
		return
	}
	var sItemParsed ast.Expr
	sItemParsed = ast.NewIdent("itemParsed")
	if o != UnionInterfaceObjectType && o != MapObjectType {
		sItemParsed = &ast.StarExpr{X: sItemParsed}
	}
	var returnElts []ast.Expr
	iterateValues := []ast.Stmt{
		&ast.DeclStmt{Decl: &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{&ast.TypeSpec{Name: ast.NewIdent("values"), Type: &ast.ArrayType{Elt: ast.NewIdent("Value")}}}}},
	}
	for i, field := range t.Fields {
		if len(field.Names) != 1 {
			continue
		}

		fieldName := field.Names[0].String()
		fieldName = fmt.Sprintf("%s%s", strings.ToUpper(string(fieldName[0])), fieldName[1:])
		var arrayType bool
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

		var varName ast.Expr
		var fStmt ast.Stmt
		fVarName := &ast.SelectorExpr{X: ast.NewIdent("s"), Sel: ast.NewIdent(fieldName)}
		if arrayType {
			varName = &ast.SliceExpr{
				X:   ast.NewIdent("patterns"),
				Low: ast.NewIdent(strconv.Itoa(i)),
			}
			fStmt = &ast.RangeStmt{
				Key:   ast.NewIdent("_"),
				Value: ast.NewIdent("v"),
				Tok:   token.DEFINE,
				X:     fVarName,
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Tok: token.ASSIGN,
							Lhs: []ast.Expr{ast.NewIdent("values")},
							Rhs: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("append"), Args: []ast.Expr{ast.NewIdent("values"), &ast.CallExpr{Fun: ast.NewIdent(fmt.Sprintf("%sToPreserves", fieldType)), Args: []ast.Expr{ast.NewIdent("v")}}}}},
						},
					},
				},
			}
		} else {
			varName = &ast.IndexExpr{
				X:     ast.NewIdent("patterns"),
				Index: ast.NewIdent(strconv.Itoa(i)),
			}
			fVarName = &ast.SelectorExpr{X: ast.NewIdent("s"), Sel: ast.NewIdent(fieldName)}
			fStmt = &ast.AssignStmt{
				Tok: token.ASSIGN,
				Lhs: []ast.Expr{ast.NewIdent("values")},
				Rhs: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("append"), Args: []ast.Expr{ast.NewIdent("values"), &ast.CallExpr{Fun: ast.NewIdent(fmt.Sprintf("%sToPreserves", fieldType)), Args: []ast.Expr{fVarName}}}}}}

		}
		returnElts = append(returnElts, &ast.KeyValueExpr{
			Key:   ast.NewIdent(fieldName),
			Value: varName,
		})
		iterateValues = append(iterateValues, fStmt)
	}
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
	keyValues = append(keyValues, &ast.KeyValueExpr{Key: ast.NewIdent("Fields"), Value: &ast.ArrayType{Elt: &ast.CompositeLit{Type: ast.NewIdent("Value"), Elts: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("extras.Reference"), Args: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("Sequence"), Args: []ast.Expr{ast.NewIdent("values")}}}}}}}})
	seqStmt := &ast.IfStmt{
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
						Index: ast.NewIdent(strconv.Itoa(0)),
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
								Name: ast.NewIdent("patterns"),
								Type: &ast.ArrayType{
									Elt: fieldTypeName,
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
											Lhs: []ast.Expr{ast.NewIdent("patterns")},
											Rhs: []ast.Expr{
												&ast.CallExpr{
													Fun: ast.NewIdent("append"),
													Args: []ast.Expr{
														ast.NewIdent("patterns"),
														sItemParsed,
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
				returnStmt,
			},
		},
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
								Y: ast.NewIdent(strconv.Itoa(1)),
							},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.IfStmt{
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
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											seqStmt,
										},
									},
								},
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
