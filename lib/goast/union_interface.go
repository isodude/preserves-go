package goast

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type UnionInterface struct {
	Name            string
	ASTs            []AST
	mapFieldsToType map[string]ObjectType
}

func NewUnionInterface(name string) *UnionInterface {
	return &UnionInterface{
		Name: name,
	}
}
func (u *UnionInterface) GetObjectType() ObjectType {
	return UnionInterfaceObjectType
}
func (u *UnionInterface) Under(a AST) {
	u.ASTs = append(u.ASTs, a)
}
func (*UnionInterface) SetKind(o ObjectType)       {}
func (*UnionInterface) SetStructKind(o ObjectType) {}
func (u *UnionInterface) GetName() string {
	return u.Name
}
func (u *UnionInterface) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(u.Name)
}
func (u *UnionInterface) AST(above AST) (decl []ast.Decl) {
	fields := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("Is%s", u.Name))},
			Type: &ast.FuncType{
				Func:   token.NoPos,
				Params: &ast.FieldList{},
			},
		},
	}
	/*
		for _, union := range u.ASTs {
			fields = append(fields, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(union.GetName())},
				Type:  &ast.FuncType{},
			})
		}*/
	h := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(u.Name),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: fields,
					},
				},
			},
		},
	}

	decl = append(decl, h)

	var body []ast.Stmt
	var elements []ast.Expr
	var toCase []ast.Stmt
	for _, s := range u.ASTs {
		aName := fmt.Sprintf("%s%s", u.Name, s.GetTitle())
		astName := ast.NewIdent(aName)

		elements = append(elements, &ast.UnaryExpr{
			Op: token.AND,
			X: &ast.CompositeLit{
				Type: astName,
				Elts: []ast.Expr{},
			},
		})
		varName := ast.NewIdent("o")
		var dVarName ast.Expr
		switch s.GetObjectType() {
		case StructObjectType:
			dVarName = varName
		default:
			dVarName = varName
		}
		toCase = append(toCase, &ast.CaseClause{
			List: []ast.Expr{&ast.StarExpr{X: astName}},
			Body: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(
								fmt.Sprintf("%sToPreserves", aName)),
							Args: []ast.Expr{
								&ast.StarExpr{X: ast.NewIdent(`u`)},
							}}}},
			}})
		body = append(body, &ast.IfStmt{

			Init: &ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{varName},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent(fmt.Sprintf("%s%s", aName, "FromPreserves")),
						Args: []ast.Expr{
							ast.NewIdent("value"),
						},
					},
				},
			},
			Cond: &ast.BinaryExpr{
				Op: token.NEQ,
				X:  varName,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							dVarName,
						},
					},
				},
			},
		})
	}
	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%sFromPreserves", u.Name)),
			Type: &ast.FuncType{

				Func: token.Pos(token.FUNC),
				Params: &ast.FieldList{

					List: []*ast.Field{{
						Names: []*ast.Ident{ast.NewIdent("value")},
						Type:  ast.NewIdent("Value"),
					}},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  ast.NewIdent(u.Name),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: append(body, &ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				}),
			},
		},
	)

	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", u.GetTitle(), "ToPreserves")),
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("s")},
							Type:  ast.NewIdent(u.GetTitle()),
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
				List: []ast.Stmt{&ast.TypeSwitchStmt{
					Assign: &ast.AssignStmt{
						Tok: token.DEFINE,
						Lhs: []ast.Expr{ast.NewIdent(`u`)},
						Rhs: []ast.Expr{
							&ast.TypeAssertExpr{
								X:    ast.NewIdent(`s`),
								Type: ast.NewIdent("type"),
							},
						},
					},
					Body: &ast.BlockStmt{
						List: toCase,
					},
				},
					&ast.ReturnStmt{
						Results: []ast.Expr{ast.NewIdent("nil")},
					},
				},
			},
		},
	)
	for _, a := range u.ASTs {
		decl = append(decl, a.AST(u)...)
	}
	return
}
