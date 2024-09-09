package goast

import (
	"go/ast"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Rec struct {
	Name            string
	Fields          []*ast.Field
	Kind            ObjectType
	StructKind      ObjectType
	mapFieldsToType map[string]ObjectType
	identifier      []AST
}

func (r *Rec) GetObjectType() ObjectType {
	return StructObjectType
}
func (r *Rec) Under(_ AST) {}
func (r *Rec) GetName() string {
	return r.Name
}
func (r *Rec) GetTitle() string { return cases.Title(language.English, cases.NoLower).String(r.Name) }
func (r *Rec) SetKind(o ObjectType) {
	r.Kind = o
}
func (r *Rec) SetStructKind(o ObjectType) {
	r.StructKind = o
}

/*
func (s *Struct) FromPreserves(Value) *Struct {

}
*/
func (r *Rec) AST(above AST) (decl []ast.Decl) {

	//	name := r.GetTitle()
	/*var nameType ast.Expr
	nameType = ast.NewIdent(name)
	if strings.ToLower(name) == "any" {
		nameType = ast.NewIdent("Value")
	}*/

	//	sname := &ast.StarExpr{X: ast.NewIdent(name)}

	/*
		func (*Lit) FromPreserves(value Value) *Lit {
			if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
				if s, ok := rec.Key.(*Symbol); ok && s.String() == "lit" {
					if a := AnyFromPreserves(rec.Fields[0]); a != nil {
						return &Lit{Value: a}
					}
				}
			}
			return nil
		}
	*/
	/*
		var returnElts []ast.Expr
		var firstStmt, curStmt ast.Stmt
		for i, field := range r.Fields {
			var arrayType bool
			if len(field.Names) != 1 {
				continue
			}

			varName := ast.NewIdent(fmt.Sprintf("p%d", i))
			var callExprFun *ast.Ident
			fieldName := field.Names[0].String()
			if fieldName == "interface" {
				fieldName = "Interface"
			}
			if fieldName == "any" {
				fieldName = "Any"
			}
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
			if callExprFun = MapToCallFunc(r.mapFieldsToType, fieldType); callExprFun == nil {
				continue
			}
			o, ok := r.mapFieldsToType[strings.ToLower(fieldType)]
			if !ok {
				continue
			}
			var dVarName ast.Expr
			dVarName = varName
			if o != UnionInterfaceObjectType && o != ValueType {
				dVarName = &ast.StarExpr{X: varName}
			}
			if r.GetTitle() == "Bundle" {
				fmt.Printf("hello")
			}
			returnElts = append(returnElts, &ast.KeyValueExpr{
				Key:   ast.NewIdent(fieldName),
				Value: dVarName,
			})
			var stmt ast.Stmt
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
											Type: ast.NewIdent(fieldType),
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

			if len(r.Fields) == i+1 {
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
				if arrayType {
					curStmt.(*ast.IfStmt).Body.List = append(curStmt.(*ast.IfStmt).Body.List, returnStmt)
				} else {
					curStmt.(*ast.IfStmt).Body.List = append(curStmt.(*ast.IfStmt).Body.List, returnStmt)
				}
			}
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
						Value: strconv.Quote(r.GetName()),
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{firstStmt},
			},
		}
		if len(r.identifier) > 0 {
			if stmter, ok := r.identifier[0].(Stmt); ok {
				ifStmt = stmter.Stmt(&ast.SelectorExpr{
					X:   ast.NewIdent("rec"),
					Sel: ast.NewIdent("Key"),
				}, []ast.Stmt{firstStmt})
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
									Y: ast.NewIdent(strconv.Itoa(len(r.Fields))),
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{ifStmt},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent("nil")},
						},
					},
				},
			},
		)
	*/
	/*
	   func RefToPreserves(r Ref) Value {
	   	return &Record{Key: NewSymbol("ref"), Fields: []Value{ModulePathToPreserves(r.Module), SymbolToPreserves(r.Name)}}
	   }
	*/
	/*
		keyValues := []ast.Expr{
			&ast.KeyValueExpr{
				Key: ast.NewIdent("Key"),
				Value: &ast.CallExpr{
					Fun:  ast.NewIdent("NewSymbol"),
					Args: []ast.Expr{ast.NewIdent(strconv.Quote(""))},
				},
			},
		}
		//if len(r.identifier) > 0 {
		//if stmter, ok := r.identifier[0].(Stmt); ok {
		// = stmter.ToStmt()
		//}
		//}
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
					List: []ast.Stmt{
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
						},
					},
				},
			},
		)
	*/
	return
}
