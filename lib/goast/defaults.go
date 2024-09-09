package goast

import (
	"go/ast"
)

type Boolean struct{}

func (*Boolean) GetObjectType() ObjectType  { return SimpleBoolType }
func (*Boolean) Under(_ AST)                {}
func (*Boolean) GetName() string            { return "boolean" }
func (*Boolean) GetTitle() string           { return "Boolean" }
func (*Boolean) SetKind(o ObjectType)       {}
func (*Boolean) SetStructKind(o ObjectType) {}
func (*Boolean) AST(above AST) (decl []ast.Decl) {
	/*
		name := "Boolean"
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
									ast.NewIdent("obj"),
									ast.NewIdent("ok"),
								},
								Rhs: []ast.Expr{
									&ast.TypeAssertExpr{
										X:    ast.NewIdent("value"),
										Type: &ast.StarExpr{X: ast.NewIdent("Boolean")},
									},
								},
							},
							Cond: ast.NewIdent("ok"),
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("obj"),
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
		decl = append(decl,
			&ast.FuncDecl{
				Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "ToPreserves")),
				Type: &ast.FuncType{
					Func: token.Pos(token.FUNC),
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("v")},
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
								&ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("v")},
							},
						},
					},
				},
			})
	*/
	return
}

type Symbol struct{}

func (*Symbol) GetObjectType() ObjectType  { return SimpleStringType }
func (*Symbol) Under(_ AST)                {}
func (*Symbol) GetName() string            { return "symbol" }
func (*Symbol) GetTitle() string           { return "Symbol" }
func (*Symbol) SetKind(o ObjectType)       {}
func (*Symbol) SetStructKind(o ObjectType) {}
func (*Symbol) AST(above AST) (decl []ast.Decl) {
	/*
		name := "Symbol"
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
									ast.NewIdent("obj"),
									ast.NewIdent("ok"),
								},
								Rhs: []ast.Expr{
									&ast.TypeAssertExpr{
										X:    ast.NewIdent("value"),
										Type: &ast.StarExpr{X: ast.NewIdent("Symbol")},
									},
								},
							},
							Cond: ast.NewIdent("ok"),
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("obj"),
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
		decl = append(decl,
			&ast.FuncDecl{
				Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "ToPreserves")),
				Type: &ast.FuncType{
					Func: token.Pos(token.FUNC),
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("v")},
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
								&ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("v")},
							},
						},
					},
				},
			})
	*/
	return
}

type SignedInteger struct{}

func (*SignedInteger) GetObjectType() ObjectType  { return SimpleSignedIntegerType }
func (*SignedInteger) Under(_ AST)                {}
func (*SignedInteger) GetName() string            { return "SignedInteger" }
func (*SignedInteger) GetTitle() string           { return "SignedInteger" }
func (*SignedInteger) SetKind(o ObjectType)       {}
func (*SignedInteger) SetStructKind(o ObjectType) {}
func (*SignedInteger) AST(above AST) (decl []ast.Decl) {
	/*
		name := "SignedInteger"
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
									ast.NewIdent("obj"),
									ast.NewIdent("ok"),
								},
								Rhs: []ast.Expr{
									&ast.TypeAssertExpr{
										X:    ast.NewIdent("value"),
										Type: &ast.StarExpr{X: ast.NewIdent("SignedInteger")},
									},
								},
							},
							Cond: ast.NewIdent("ok"),
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("obj"),
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
		decl = append(decl,
			&ast.FuncDecl{
				Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "ToPreserves")),
				Type: &ast.FuncType{
					Func: token.Pos(token.FUNC),
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("v")},
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
								&ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("v")},
							},
						},
					},
				},
			})
	*/
	return
}

type Pstring struct{}

func (*Pstring) GetObjectType() ObjectType  { return SimplePstringType }
func (*Pstring) Under(_ AST)                {}
func (*Pstring) GetName() string            { return "Pstring" }
func (*Pstring) GetTitle() string           { return "Pstring" }
func (*Pstring) SetKind(o ObjectType)       {}
func (*Pstring) SetStructKind(o ObjectType) {}
func (*Pstring) AST(above AST) (decl []ast.Decl) {
	/*
		name := "Pstring"
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
									ast.NewIdent("obj"),
									ast.NewIdent("ok"),
								},
								Rhs: []ast.Expr{
									&ast.TypeAssertExpr{
										X:    ast.NewIdent("value"),
										Type: &ast.StarExpr{X: ast.NewIdent("Pstring")},
									},
								},
							},
							Cond: ast.NewIdent("ok"),
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("obj"),
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
		decl = append(decl,
			&ast.FuncDecl{
				Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "ToPreserves")),
				Type: &ast.FuncType{
					Func: token.Pos(token.FUNC),
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("v")},
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
								&ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("v")},
							},
						},
					},
				},
			})
	*/
	return
}

type Value struct{}

func (*Value) GetObjectType() ObjectType  { return ValueType }
func (*Value) Under(_ AST)                {}
func (*Value) GetName() string            { return "Value" }
func (*Value) GetTitle() string           { return "Value" }
func (*Value) SetKind(o ObjectType)       {}
func (*Value) SetStructKind(o ObjectType) {}
func (*Value) AST(above AST) (decl []ast.Decl) {
	/*name := "Value"
	var stmts []ast.Stmt
	for _, t := range []AST{&Boolean{}, &SignedInteger{}, &Pstring{}, &Symbol{}} {
		stmts = append(stmts, &ast.IfStmt{
			Init: &ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					ast.NewIdent("a"),
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent(fmt.Sprintf("%sFromPreserves", t.GetTitle())),
						Args: []ast.Expr{
							ast.NewIdent("value"),
						},
					},
				},
			},
			Cond: &ast.BinaryExpr{
				Op: token.NEQ,
				X:  ast.NewIdent("a"),
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{ast.NewIdent("a")}}}},
		})
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
							Type:  ast.NewIdent(name),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: append(stmts, &ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				}),
			},
		},
	)
	*/
	/*

	   func ValueToPreserves(v Value) Value {
	   	return v
	   }
	*/
	/*
		decl = append(decl,
			&ast.FuncDecl{
				Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "ToPreserves")),
				Type: &ast.FuncType{
					Func: token.Pos(token.FUNC),
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("v")},
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
								ast.NewIdent("v"),
							},
						},
					},
				},
			})
	*/
	return
}
