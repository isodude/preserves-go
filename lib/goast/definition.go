package goast

import (
	"fmt"
	"go/ast"

	"github.com/isodude/preserves-go/lib/preserves"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Definition struct {
	Name            string
	ASTs            []AST
	Fields          []*ast.Field
	ASTFields       []*Field
	Kind            *ObjectType
	StructKind      *ObjectType
	mapFieldsToType map[string]ObjectType
	MapKeyToField   []preserves.Value
}

func NewDefinition(name string) *Definition {
	return &Definition{
		Name: name,
	}
}

func (d *Definition) GetObjectType() ObjectType {
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

func (d *Definition) Under(a AST) {
	d.ASTs = append(d.ASTs, a)
}
func (d *Definition) AddField(a *Field) {
	d.ASTFields = append(d.ASTFields, a)
}
func (d *Definition) SetKind(o ObjectType) {
	d.Kind = &o
}
func (d *Definition) SetStructKind(o ObjectType) {
	d.StructKind = &o
}
func (d *Definition) GetName() string {
	return d.Name
}
func (d *Definition) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(d.Name)
}
func (d *Definition) AST(above AST) (decl []ast.Decl) {
	if d.StructKind != nil {
		s := &Struct{
			Name:            d.Name,
			Fields:          d.Fields,
			StructKind:      *d.StructKind,
			mapFieldsToType: d.mapFieldsToType,
			MapKeyToField:   d.MapKeyToField,
		}
		if d.Kind != nil {
			s.Kind = *d.Kind
		}
		decl = append(decl, s.AST(above)...)
		return
	}

	name := d.GetName()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetName(), name)
	}
	unionInterface := NewUnionInterface(name)
	unionInterface.ASTs = d.ASTs
	unionInterface.mapFieldsToType = d.mapFieldsToType
	decl = append(decl, unionInterface.AST(above)...)

	return

}
