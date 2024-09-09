package goast

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Map struct {
	Name            string
	Key             string
	Value           string
	Union           string
	mapFieldsToType map[string]ObjectType
}

func (m *Map) GetObjectType() ObjectType {
	return MapObjectType
}
func (m *Map) GetName() string {
	return m.Name
}
func (m *Map) GetTitle() string         { return cases.Title(language.English, cases.NoLower).String(m.Name) }
func (m *Map) Under(_ AST)              {}
func (*Map) SetKind(o ObjectType)       {}
func (*Map) SetStructKind(o ObjectType) {}

/*
	func ModulesFromPreserves(value Value) *Modules {
		if dict, ok := value.(*Dictionary); ok {
			m := NewModules()
			for k, v := range *dict {
				if a := ModulePathFromPreserves(k); a != nil {
					if b := (&Schema{}).FromPreserves(v); b != nil {
						m[a]=*b
						continue
					}
				}
				return nil
			}
		}
		return nil
	}
*/
func (m *Map) AST(above AST) (decl []ast.Decl) {
	name := m.GetTitle()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
	}
	key := m.Key
	if m.Key == "any" {
		key = "Value"
	}
	value := m.Value
	if m.Value == "any" {
		value = "Value"
	}
	var sKey, sValue ast.Expr
	sKey = ast.NewIdent(key)
	if o, ok := m.mapFieldsToType[strings.ToLower(key)]; ok {
		if o == StructSeqofType {
			sKey = &ast.IndexExpr{X: &ast.SelectorExpr{X: ast.NewIdent("extras"), Sel: ast.NewIdent("Hash")}, Index: sKey}
		}
	}
	sValue = ast.NewIdent(value)
	decl = append(decl, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.MapType{
					Key:   sKey,
					Value: sValue,
				},
			},
		},
	})
	/*
		var smallFields []*ast.Field
		for _, field := range []string{m.Key, m.Value} {
			fname := fmt.Sprintf("%s%s", strings.ToLower(field[0:1]), field[1:])
			smallFields = append(smallFields, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fname)},
				Type:  ast.NewIdent(field),
			})
		}
		var smallValues []ast.Expr
		for _, field := range []string{m.Key, m.Value} {
			fname := fmt.Sprintf("%s%s", strings.ToLower(field[0:1]), field[1:])
			smallValues = append(smallValues, &ast.KeyValueExpr{
				Key:   ast.NewIdent(field),
				Value: ast.NewIdent(fname),
			})
		}
	*/
	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("New%s", name)),
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  ast.NewIdent(name),
						},
					},
				},
			},
			Body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun:  ast.NewIdent("make"),
							Args: []ast.Expr{ast.NewIdent(name)},
						},
					},
				},
			},
			},
		},
	)
	if above != nil && above.GetObjectType() == UnionInterfaceObjectType {
		decl = append(decl,
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  ast.NewIdent(name),
						}},
				},
				Name: ast.NewIdent(fmt.Sprintf("Is%s", above.GetTitle())),
				Type: &ast.FuncType{
					Func:   token.Pos(token.FUNC),
					Params: &ast.FieldList{},
				},
				Body: &ast.BlockStmt{},
			})
	}
	callFuncKeyExpr := MapToCallFunc(m.mapFieldsToType, key)
	callFuncValueExpr := MapToCallFunc(m.mapFieldsToType, value)
	if callFuncKeyExpr == nil || callFuncValueExpr == nil {
		return
	}
	o, ok := m.mapFieldsToType[strings.ToLower(key)]
	if !ok {
		return
	}

	var dKey ast.Expr
	var toKey ast.Expr
	dKey = ast.NewIdent("dKey")
	toKey = ast.NewIdent("dKey")
	if o == StructSeqofType {
		dKey = &ast.CallExpr{Fun: &ast.SelectorExpr{X: dKey, Sel: ast.NewIdent("ToHash")}}
		toKey = &ast.CallExpr{Fun: &ast.SelectorExpr{X: toKey, Sel: ast.NewIdent("FromHash")}}
	} else if o != UnionInterfaceObjectType && o != MapObjectType && o != ValueType {
		dKey = &ast.StarExpr{X: dKey}
	}
	o, ok = m.mapFieldsToType[strings.ToLower(m.Value)]
	if !ok {
		return
	}
	var dValue ast.Expr
	dValue = ast.NewIdent("dValue")

	if o != UnionInterfaceObjectType && o != MapObjectType && o != ValueType {
		dValue = &ast.StarExpr{X: dValue}
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
							Type:  &ast.StarExpr{X: ast.NewIdent(name)},
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
								ast.NewIdent("dict"),
								ast.NewIdent("ok"),
							},
							Rhs: []ast.Expr{
								&ast.TypeAssertExpr{
									X:    ast.NewIdent("value"),
									Type: &ast.StarExpr{X: ast.NewIdent("Dictionary")},
								},
							},
						},
						Cond: ast.NewIdent("ok"),
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.AssignStmt{
									Tok: token.DEFINE,
									Lhs: []ast.Expr{ast.NewIdent("obj")},
									Rhs: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent(fmt.Sprintf("New%s", name))}},
								},

								&ast.RangeStmt{
									Tok:   token.DEFINE,
									Key:   ast.NewIdent("dictKey"),
									Value: ast.NewIdent("dictValue"),
									X:     &ast.StarExpr{X: ast.NewIdent("dict")},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.IfStmt{
												Init: &ast.AssignStmt{
													Tok: token.DEFINE,
													Lhs: []ast.Expr{
														ast.NewIdent("dKey"),
													},
													Rhs: []ast.Expr{
														&ast.CallExpr{
															Fun: callFuncKeyExpr,
															Args: []ast.Expr{
																ast.NewIdent("dictKey"),
															},
														},
													},
												},
												Cond: &ast.BinaryExpr{
													Op: token.NEQ,
													X:  ast.NewIdent("dKey"),
													Y:  ast.NewIdent("nil"),
												},
												Body: &ast.BlockStmt{
													List: []ast.Stmt{
														&ast.IfStmt{
															Init: &ast.AssignStmt{
																Tok: token.DEFINE,
																Lhs: []ast.Expr{
																	ast.NewIdent("dValue"),
																},
																Rhs: []ast.Expr{
																	&ast.CallExpr{
																		Fun: callFuncValueExpr,
																		Args: []ast.Expr{
																			ast.NewIdent("dictValue"),
																		},
																	},
																},
															},
															Cond: &ast.BinaryExpr{
																Op: token.NEQ,
																X:  ast.NewIdent("dValue"),
																Y:  ast.NewIdent("nil"),
															},
															Body: &ast.BlockStmt{
																List: []ast.Stmt{
																	&ast.AssignStmt{
																		Tok: token.ASSIGN,
																		Lhs: []ast.Expr{&ast.IndexExpr{
																			X:     ast.NewIdent("obj"),
																			Index: dKey,
																		}},
																		Rhs: []ast.Expr{
																			dValue,
																		},
																	},
																	&ast.ExprStmt{
																		X: &ast.BasicLit{Kind: token.STRING, Value: "continue"},
																	},
																},
															},
														},
													},
												},
											},
											&ast.ReturnStmt{
												Results: []ast.Expr{
													ast.NewIdent("nil"),
												},
											},
										},
									},
								},
								&ast.ReturnStmt{
									Results: []ast.Expr{
										&ast.UnaryExpr{
											Op: token.AND,
											X:  ast.NewIdent("obj"),
										},
									},
								},
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("nil"),
						},
					},
				},
			},
		})
	/*
		func ModulesToPreserves(m Modules) Value {
			d := make(Dictionary)
			for k, v := range m {
				d[ModulePathToPreserves(k)] = SchemaToPreserves(v)
			}
			return &d
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
							Names: []*ast.Ident{ast.NewIdent("m")},
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
					&ast.DeclStmt{Decl: &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{&ast.TypeSpec{Assign: token.Pos(token.ASSIGN), Name: ast.NewIdent("dictionary"), Type: &ast.CallExpr{Fun: ast.NewIdent("make"), Args: []ast.Expr{ast.NewIdent("Dictionary")}}}}}},
					&ast.RangeStmt{
						Key:   ast.NewIdent("dKey"),
						Value: ast.NewIdent("dValue"),
						Tok:   token.DEFINE,
						X:     ast.NewIdent("m"),
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.AssignStmt{
									Lhs: []ast.Expr{
										&ast.IndexExpr{X: ast.NewIdent("dictionary"), Index: &ast.CallExpr{Fun: ast.NewIdent(fmt.Sprintf("%sToPreserves", key)), Args: []ast.Expr{toKey}}},
									},
									Rhs: []ast.Expr{
										&ast.CallExpr{Fun: ast.NewIdent(fmt.Sprintf("%sToPreserves", value)), Args: []ast.Expr{ast.NewIdent("dValue")}},
									},
									Tok: token.ASSIGN,
								},
							},
						},
					},

					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("dictionary")},
						},
					},
				},
			},
		})
	return
}
