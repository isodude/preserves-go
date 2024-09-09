package goast

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type UnionVariant struct {
	Name            string
	ASTs            []AST
	mapFieldsToType map[string]ObjectType
}

func NewUnionVariant(name string) *UnionVariant {
	return &UnionVariant{
		Name: name,
	}
}
func (u *UnionVariant) GetObjectType() ObjectType {
	return UnionVariantObjectType
}
func (u *UnionVariant) Under(a AST) {
	u.ASTs = append(u.ASTs, a)
}
func (*UnionVariant) SetKind(o ObjectType)       {}
func (*UnionVariant) SetStructKind(o ObjectType) {}
func (u *UnionVariant) GetName() string {
	return u.Name
}
func (u *UnionVariant) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(u.Name)
}

/*
type EmbeddedTypeNameVariant string

const (

	EmbeddedTypeNameVariantFalse EmbeddedTypeNameVariant = "false"
	EmbeddedTypeNameVariantRef   EmbeddedTypeNameVariant = "Ref"

)

	type EmbeddedTypeName struct {
		_variant EmbeddedTypeNameVariant
		False    bool
		Ref      Ref
	}

	func (e *EmbeddedTypeName) Value() any {
		switch e._variant {
		case EmbeddedTypeNameVariantFalse:
			return &([]bool{e.False}[0])
		case EmbeddedTypeNameVariantRef:
			return e.Ref
		default:
			return nil
		}
	}

func NewEmbeddedTypeName() *EmbeddedTypeName {

}

	func (*EmbeddedTypeName) FromPreserves(value Value) *EmbeddedTypeName {
		if b, ok := value.(*Boolean); ok && !(bool(*b)) {
			return &EmbeddedTypeName{_variant: EmbeddedTypeNameVariantFalse}
		}
		if ref := (&Ref{}).FromPreserves(value); ref != nil {
			return &EmbeddedTypeName{_variant: EmbeddedTypeNameVariantRef, Ref: *ref}
		}
		return nil
	}
*/
func (u *UnionVariant) AST(above AST) (decl []ast.Decl) {
	name := u.GetTitle()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
	}
	astVariant := ast.NewIdent(fmt.Sprintf("%sVariant", u.Name))
	h := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: astVariant,
				Type: &ast.BasicLit{Kind: token.STRING, Value: "string"},
			},
		},
	}
	decl = append(decl, h)
	fields := &ast.FieldList{
		List: []*ast.Field{
			{Names: []*ast.Ident{ast.NewIdent("_variant")}, Type: astVariant},
		},
	}
	var stmts []ast.Stmt
	for _, a := range u.ASTs {
		decl = append(decl, a.AST(u)...)
		typ := fmt.Sprintf("%s%s", u.GetTitle(), a.GetTitle())
		var typExpr ast.Expr
		typExpr = ast.NewIdent(typ)
		l := &ast.GenDecl{
			Tok: token.CONST,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names:  []*ast.Ident{ast.NewIdent(fmt.Sprintf("%sVariant%s", u.Name, a.GetTitle()))},
					Type:   astVariant,
					Values: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", a.GetName())}},
				},
			},
		}
		decl = append(decl, l)
		field := &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(a.GetTitle())},
			Type:  &ast.StarExpr{X: typExpr},
		}
		fields.List = append(fields.List, field)
		var callExprFun *ast.Ident
		if callExprFun = MapToCallFunc(u.mapFieldsToType, typ); callExprFun == nil {
			continue
		}

		stmts = append(stmts, &ast.IfStmt{

			Init: &ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("o")},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: callExprFun,
						Args: []ast.Expr{
							ast.NewIdent("value"),
						},
					},
				},
			},
			Cond: &ast.BinaryExpr{
				Op: token.NEQ,
				X:  ast.NewIdent("o"),
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{Op: token.AND, X: &ast.CompositeLit{Type: ast.NewIdent(name), Elts: []ast.Expr{&ast.KeyValueExpr{Key: ast.NewIdent(a.GetTitle()), Value: ast.NewIdent("o")}}}},
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
				List: append(stmts,
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("nil"),
						},
					}),
			},
		},
	)
	decl = append(decl, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(u.GetName()),
				Type: &ast.StructType{
					Fields: fields,
				},
			},
		},
	})
	return
}
