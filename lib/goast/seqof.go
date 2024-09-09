package goast

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Seqof struct {
	Name            string
	Fields          []*ast.Field
	Kind            ObjectType
	mapFieldsToType map[string]ObjectType
}

func (s *Seqof) GetObjectType() ObjectType {
	return StructSeqofType
}
func (s *Seqof) GetName() string {
	return s.Name
}
func (s *Seqof) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(s.Name)
}
func (s *Seqof) Under(_ AST)              {}
func (*Seqof) SetKind(o ObjectType)       {}
func (*Seqof) SetStructKind(o ObjectType) {}
func (s *Seqof) AST(above AST) (decl []ast.Decl) {
	if len(s.Fields) != 1 {
		return
	}
	name := s.GetTitle()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
	}
	elt := s.Fields[0].Type.(*ast.ArrayType).Elt
	decl = append(decl, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: s.Fields[0].Type,
			},
		},
	})

	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("New%s", name)),
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("items")},
							Type:  &ast.Ellipsis{Elt: elt},
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
			Body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						ast.NewIdent("_items"),
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun:  ast.NewIdent(name),
							Args: []ast.Expr{ast.NewIdent("items")},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent("_items"),
						},
					},
				},
			},
			},
		},
	)

	varName := ast.NewIdent("itemParsed")
	var dVarName ast.Expr
	dVarName = varName
	if o, ok := s.mapFieldsToType[strings.ToLower(elt.(*ast.Ident).String())]; ok {
		if o != UnionInterfaceObjectType && o != MapObjectType {
			dVarName = &ast.StarExpr{X: varName}
		}
	}

	sname := &ast.StarExpr{X: ast.NewIdent(name)}
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
						Cond: ast.NewIdent("ok"),
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.DeclStmt{
									Decl: &ast.GenDecl{
										Tok: token.VAR,
										Specs: []ast.Spec{
											&ast.TypeSpec{
												Name: ast.NewIdent("items"),
												Type: s.Fields[0].Type,
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
															Fun: ast.NewIdent(fmt.Sprintf("%sFromPreserves", elt.(*ast.Ident).String())),
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
															Lhs: []ast.Expr{ast.NewIdent("items")},
															Rhs: []ast.Expr{
																&ast.CallExpr{
																	Fun: ast.NewIdent("append"),
																	Args: []ast.Expr{
																		ast.NewIdent("items"),
																		dVarName,
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
								&ast.AssignStmt{
									Tok: token.DEFINE,
									Lhs: []ast.Expr{
										ast.NewIdent("_items"),
									},
									Rhs: []ast.Expr{
										&ast.CallExpr{
											Fun:  ast.NewIdent(name),
											Args: []ast.Expr{ast.NewIdent("items")},
										},
									},
								},
								&ast.ReturnStmt{
									Results: []ast.Expr{
										&ast.UnaryExpr{
											Op: token.AND,
											X:  ast.NewIdent("_items"),
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
		})

	/*
		func ModulesPathToPreserves(d ModulePath) Value {
			var sequence []Value
			for _, k := range d {
				sequence = append(sequence, SymbolToPreserves(k))
			}
			return &Sequence(sequence)
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
					&ast.DeclStmt{Decl: &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{&ast.TypeSpec{Name: ast.NewIdent("values"), Type: &ast.ArrayType{Elt: ast.NewIdent("Value")}}}}},
					&ast.RangeStmt{
						Key:   ast.NewIdent("_"),
						Value: ast.NewIdent("v"),
						Tok:   token.DEFINE,
						X:     ast.NewIdent("d"),
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.AssignStmt{
									Tok: token.ASSIGN,
									Lhs: []ast.Expr{ast.NewIdent("values")},
									Rhs: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("append"), Args: []ast.Expr{ast.NewIdent("values"), &ast.CallExpr{Fun: ast.NewIdent(fmt.Sprintf("%sToPreserves", elt.(*ast.Ident).String())), Args: []ast.Expr{ast.NewIdent("v")}}}}},
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("s"),
						},
						Rhs: []ast.Expr{&ast.CallExpr{
							Fun:  ast.NewIdent("Sequence"),
							Args: []ast.Expr{ast.NewIdent("values")},
						}},
						Tok: token.DEFINE,
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{&ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("s")}},
					}},
			},
		},
	)

	/*
			func (m ModulePath) Hash() (s string) {
			for i, k := range m {
				s += k.String()
				if len(n) != i+1 {
					s += "."
				}
			}
			return
		}
	*/
	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent("Hash"),
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("m")},
						Type:  ast.NewIdent(name),
					},
				},
			},
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("s")},
							Type:  ast.NewIdent("string"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.RangeStmt{
						Key:   ast.NewIdent("i"),
						Value: ast.NewIdent("v"),
						Tok:   token.DEFINE,
						X:     ast.NewIdent("m"),
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.AssignStmt{
									Tok: token.ADD_ASSIGN,
									Lhs: []ast.Expr{ast.NewIdent("s")},
									Rhs: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("v.String")}},
								},
								&ast.IfStmt{
									Cond: &ast.BinaryExpr{Op: token.NEQ,
										X: &ast.CallExpr{
											Fun:  ast.NewIdent("len"),
											Args: []ast.Expr{ast.NewIdent("m")},
										},
										Y: &ast.BinaryExpr{Op: token.ADD, X: ast.NewIdent("i"), Y: ast.NewIdent("1")},
									},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.AssignStmt{
												Tok: token.ADD_ASSIGN,
												Lhs: []ast.Expr{
													ast.NewIdent("s"),
												},
												Rhs: []ast.Expr{
													ast.NewIdent(strconv.Quote(".")),
												},
											},
										},
									},
								},
							},
						},
					},
					&ast.ReturnStmt{},
				},
			},
		},
	)

	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent("ToHash"),
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("m")},
						Type:  &ast.StarExpr{X: ast.NewIdent(name)},
					},
				},
			},
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  &ast.IndexExpr{X: &ast.SelectorExpr{X: ast.NewIdent("extras"), Sel: ast.NewIdent("Hash")}, Index: ast.NewIdent(name)},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CallExpr{Fun: &ast.SelectorExpr{X: ast.NewIdent("extras"), Sel: ast.NewIdent("NewHash")}, Args: []ast.Expr{&ast.StarExpr{X: ast.NewIdent("m")}}},
						},
					},
				},
			},
		},
	)
	return
}
