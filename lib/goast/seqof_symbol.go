package goast

import (
	"go/ast"
	"go/token"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type SeqofSymbol struct {
	Name string
}

func (s *SeqofSymbol) GetObjectType() ObjectType {
	return SeqofSymbolObjectType
}
func (s *SeqofSymbol) GetName() string {
	return s.Name
}
func (s *SeqofSymbol) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(s.Name)
}
func (s *SeqofSymbol) Under(_ AST)              {}
func (*SeqofSymbol) SetKind(o ObjectType)       {}
func (*SeqofSymbol) SetStructKind(o ObjectType) {}
func (s *SeqofSymbol) AST(above AST) (decl []ast.Decl) {
	decl = append(decl, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(s.Name),
				Type: ast.NewIdent("string"),
			},
		},
	})
	return
}
