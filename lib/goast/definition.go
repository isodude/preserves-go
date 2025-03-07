package goast

import (
	"fmt"
	"go/ast"
	"strings"

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
	Value           preserves.Value
	Identifier      []AST
	mapFieldsToType map[string]ObjectType
	MapKeyToField   []preserves.Value
}

func NewDefinition(name string) *Definition {
	return &Definition{
		Name: name,
	}
}

func (d *Definition) SetValue(v preserves.Value) {
	d.StructKind = &([]ObjectType{LitType}[0])
	d.Value = v
}
func (d *Definition) GetObjectType() ObjectType {
	if len(d.ASTs) == 1 {
		return d.ASTs[0].GetObjectType()
	}
	if d.Kind != nil {
		if *d.Kind == StructObjectType {
			return *d.StructKind
		}
		return *d.Kind
	}
	return DefinitionObjectType
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
	name := d.GetName()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetName(), name)
	}
	if d.Kind == nil {
		panic(fmt.Sprintf("d.Kind for %s is nil", name))
	}
	switch *d.Kind {
	case PassthroughObjectType:
		u := NewPassthrough(name)
		u.ASTs = d.ASTs
		for _, f := range d.ASTFields {
			u.Object = f.Type
			if t, ok := d.mapFieldsToType[strings.ToLower(f.Type)]; ok {
				u.ObjectType = t
			} else {
				panic(fmt.Sprintf("could not find type %s in %v", strings.ToLower(f.Type), d.mapFieldsToType))
			}
		}
		if d.Kind != nil {
			u.ObjectType = *d.Kind
		}
		u.mapFieldsToType = d.mapFieldsToType
		decl = append(decl, u.AST(nil)...)
	case InterfaceObjectType:
		unionInterface := NewUnionInterface(name)
		unionInterface.ASTs = d.ASTs
		unionInterface.mapFieldsToType = d.mapFieldsToType
		decl = append(decl, unionInterface.AST(nil)...)
	case StructObjectType:
		if d.StructKind == nil {
			panic(fmt.Sprintf("d.StructKind for %s is nil", name))
		}
		s := &Struct{
			Name:            name,
			Fields:          d.Fields,
			StructKind:      *d.StructKind,
			mapFieldsToType: d.mapFieldsToType,
			MapKeyToField:   d.MapKeyToField,
			Identifier:      d.Identifier,
			Value:           d.Value,
		}
		s.Kind = *d.Kind
		decl = append(decl, s.AST(above)...)
	default:
		panic(fmt.Sprintf("kind %d is unknown", *d.Kind))
	}

	return

}
