package goast

import (
	"go/ast"
)

type Bundle struct {
	Paths   [][]string
	Schemas []*Schema
	title
}

func NewBundle(name string) *Bundle {
	return &Bundle{title: title{name: name}}
}

func (b *Bundle) AddSchema(path []string, schema *Schema) {
	b.Paths = append(b.Paths, path)
	b.Schemas = append(b.Schemas, schema)
}

func (b *Bundle) Definitions() (d []*Definition) {
	for _, schema := range b.Schemas {
		d = append(d, schema.Definitions...)
	}
	return
}

func (b *Bundle) AST(_ AST) (r []ast.Decl) {
	return
}
func (b *Bundle) GetObjectType() ObjectType  { return InvalidObjectType }
func (b *Bundle) SetKind(_ ObjectType)       { return }
func (b *Bundle) SetStructKind(_ ObjectType) { return }
func (b *Bundle) Under(_ AST)                { return }

type Schema struct {
	Version      *Field
	Definitions  []*Definition
	EmbeddedType *Field
}

func NewSchema() *Schema {
	return &Schema{}
}

func (s *Schema) AddDefinition(a *Definition) {
	s.Definitions = append(s.Definitions, a)
}
