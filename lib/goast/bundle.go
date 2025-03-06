package goast

import (
	"go/ast"
)

type Bundle struct {
	Name    string
	Paths   [][]string
	Schemas []*Schema
}

func NewBundle(name string) *Bundle {
	return &Bundle{
		Name: name,
	}
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

func (b *Bundle) AST() (r []*ast.Decl) {
	return
}

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
