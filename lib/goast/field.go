package goast

import (
	"go/ast"

	"github.com/isodude/preserves-go/lib/preserves"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Field struct {
	Name            string
	ASTs            []AST
	Type            string
	Value           preserves.Value
	mapFieldsToType map[string]ObjectType
}

func NewField(name string) *Field {
	return &Field{
		Name: name,
	}
}

func (f *Field) SetValue(v preserves.Value) {
	f.Value = v
}
func (f *Field) GetObjectType() ObjectType {
	/*allLit := true
	oneLit := false
	for _, a := range u.ASTs {
		if _, ok := a.(*Lit); !ok {
			allLit = false
		} else {
			oneLit = true
		}
	}
	if allLit {
		return UnionConstObjectType
	}
	if oneLit {
		return UnionVariantObjectType
	}*/
	return UnionInterfaceObjectType
}

func (f *Field) Under(a AST) {
	f.ASTs = append(f.ASTs, a)
}
func (*Field) SetKind(o ObjectType) {}
func (f *Field) SetType(s string) {
	f.Type = s
}
func (*Field) SetStructKind(o ObjectType) {}
func (f *Field) GetName() string {
	return f.Name
}
func (f *Field) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(f.Name)
}
func (f *Field) AST(_ AST) (decl []ast.Decl) { return }
func (f *Field) ASTField() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(f.Name)},
		Type:  ast.NewIdent(f.Type),
	}
}
