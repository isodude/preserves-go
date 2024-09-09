package goast

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type UnionConst struct {
	Name            string
	ASTs            []AST
	mapFieldsToType map[string]ObjectType
}

func NewUnionConst(name string) *UnionConst {
	return &UnionConst{
		Name: name,
	}
}
func (u *UnionConst) GetObjectType() ObjectType {
	return UnionConstObjectType
}
func (u *UnionConst) Under(a AST) {
	u.ASTs = append(u.ASTs, a)
}
func (*UnionConst) SetKind(o ObjectType)       {}
func (*UnionConst) SetStructKind(o ObjectType) {}
func (u *UnionConst) GetName() string {
	return u.Name
}
func (u *UnionConst) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(u.Name)
}

/*
type AtomKind string

const AtomKindBoolean AtomKind = "Boolean"
const AtomKindDouble AtomKind = "Double"
const AtomKindSignedInteger AtomKind = "SignedInteger"
const AtomKindString AtomKind = "String"
const AtomKindByteString AtomKind = "ByteString"
const AtomKindSymbol AtomKind = "Symbol"

	func AtomKindFromPreserves(value Value) *AtomKind {
		if symbol, ok := value.(*Symbol); ok {
			switch string(*symbol) {
			case string(AtomKindBoolean):
				return &([]AtomKind{AtomKindBoolean}[0])
			case string(AtomKindDouble):
				return &([]AtomKind{AtomKindDouble}[0])
			case string(AtomKindSignedInteger):
				return &([]AtomKind{AtomKindSignedInteger}[0])
			case string(AtomKindString):
				return &([]AtomKind{AtomKindString}[0])
			case string(AtomKindByteString):
				return &([]AtomKind{AtomKindByteString}[0])
			case string(AtomKindSymbol):
				return &([]AtomKind{AtomKindSymbol}[0])
			}
		}

		return nil
	}
*/
func (u *UnionConst) AST(above AST) (decl []ast.Decl) {
	name := u.GetTitle()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
	}
	h := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(u.Name),
				Type: &ast.BasicLit{Kind: token.STRING, Value: "string"},
			},
		},
	}
	for _, a := range u.ASTs {
		if lit, ok := a.(*Lit); ok {
			l := &ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names:  []*ast.Ident{ast.NewIdent(fmt.Sprintf("%s%s", u.Name, a.GetName()))},
						Type:   &ast.BasicLit{Kind: token.STRING, Value: u.Name},
						Values: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", lit.Type)}},
					},
				},
			}
			decl = append(decl, l)
		}
	}
	decl = append(decl, h)
	var caseClauses []ast.Stmt
	for _, a := range u.ASTs {
		caseClauses = append(caseClauses, &ast.CaseClause{
			List: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("string"), Args: []ast.Expr{ast.NewIdent(fmt.Sprintf("%s%s", u.GetName(), a.GetName()))}}},
			Body: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.ParenExpr{
								X: &ast.IndexExpr{
									Index: ast.NewIdent("0"),
									X: &ast.CompositeLit{
										Type: &ast.ArrayType{Elt: ast.NewIdent(u.GetName())},
										Elts: []ast.Expr{ast.NewIdent(fmt.Sprintf("%s%s", u.GetName(), a.GetName()))},
									},
								},
							},
						},
					},
				},
			},
		})
	}
	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "FromPreserves")),
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
							Type:  &ast.StarExpr{X: ast.NewIdent(u.Name)},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.IfStmt{
						Init: &ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{ast.NewIdent("symbol"),
								ast.NewIdent("ok")},
							Rhs: []ast.Expr{
								&ast.TypeAssertExpr{
									X:    ast.NewIdent("value"),
									Type: &ast.StarExpr{X: ast.NewIdent("Symbol")},
								}},
						},
						Cond: ast.NewIdent("ok"),
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.SwitchStmt{
									Tag: &ast.CallExpr{
										Fun:  ast.NewIdent("string"),
										Args: []ast.Expr{&ast.StarExpr{X: ast.NewIdent("symbol")}},
									},
									Body: &ast.BlockStmt{
										List: caseClauses,
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
	)
	return
}
