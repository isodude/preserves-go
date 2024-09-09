package schema

import (
	"fmt"
	"go/ast"
	"reflect"
	"slices"
	"strings"

	"github.com/isodude/preserves-go/lib/goast"
	. "github.com/isodude/preserves-go/lib/preserves"
)

func (d DictionaryEntries) EncodeToGoASTFields(obj goast.AST) (fields []*ast.Field) {
	var keys []Value
	for k := range d {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(e1, e2 Value) int {
		return e1.Cmp(e2)
	})
	for _, k := range keys {
		v := d[k]
		if f, ok := v.(goast.Fields); ok {
			_f := f.EncodeToGoASTFields(obj)
			if len(_f) > 0 {
				fields = append(fields, _f...)
				if r, ok := obj.(*goast.Struct); ok {
					r.MapKeyToField = append(r.MapKeyToField, k)
				}
			}
		}
	}
	return
}

func (r *Ref) ASTString() string {
	sname := string(r.Name)
	var (
		smodule string
	)
	for _, k := range r.Module {
		smodule = fmt.Sprintf("%s%s.", smodule, string(k))
	}
	sname = fmt.Sprintf("%s%s", smodule, sname)
	return sname
}
func (r *Ref) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	sname := string(r.Name)
	var (
		smodule string
	)
	for _, k := range r.Module {
		smodule = fmt.Sprintf("%s%s.", smodule, string(k))
	}
	sname = fmt.Sprintf("%s%s", smodule, sname)
	return []*ast.Field{{
		Names: []*ast.Ident{},
		Type:  ast.NewIdent(sname),
	}}
}
func (r *Ref) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	if above != nil {
		sname := string(r.Name)
		var (
			smodule string
		)
		for _, k := range r.Module {
			smodule = fmt.Sprintf("%s%s.", smodule, string(k))
		}
		sname = fmt.Sprintf("%s%s", smodule, sname)

		p := &goast.Passthrough{Name: name, Object: sname}
		above.Under(p)
	}
	return
}

func (s *Schema) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	var keys []Symbol
	for k := range s.Definitions {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		v := s.Definitions[k]
		if g, ok := v.(goast.Encoder); ok {
			asts = append(asts, g.EncodeToGoAST(above, string(k))...)
		}
	}
	return
}

func (n *NamedAlternative) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	if encoder, ok := n.Pattern.(goast.Encoder); ok {
		asts = append(asts, encoder.EncodeToGoAST(above, string(n.VariantLabel))...)
	}
	return
}

func (b *Binding) ASTString() string { return string(b.Name) }
func (b *Binding) EncodeToGoASTFields(obj goast.AST) (fields []*ast.Field) {
	var d []*ast.Field
	if s, ok := b.Pattern.(goast.Fields); ok {
		d = s.EncodeToGoASTFields(obj)
	}
	for _, f := range d {
		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(string(b.Name))},
			Type:  f.Type,
		})
	}
	return
}
func (*Binding) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	return
}

func (s *SimplePatternAny) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	return []*ast.Field{{
		Names: []*ast.Ident{},
		Type:  ast.NewIdent("Any"),
	}}
}

func (a *SimplePatternAtom) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	typ := reflect.TypeOf(a.AtomKind).String()
	typ = strings.TrimPrefix(typ, "*schema.AtomKind")
	return []*ast.Field{{
		Names: []*ast.Ident{},
		Type:  ast.NewIdent(typ),
	}}
}

func (l *SimplePatternLit) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	//if len(funcs) > 0 {
	//	name = fmt.Sprintf("%s%s", strToCamelCase(funcs[len(funcs)-1]), name)
	//}

	u := &goast.Lit{Name: name, Type: l.Value}

	if above != nil {
		above.Under(u)
	} else {
		asts = append(asts, u)
	}
	return
}

func (s *SimplePatternSeqof) EncodeToGoASTFields(obj goast.AST) (fields []*ast.Field) {

	var d []*ast.Field
	if f, ok := s.Pattern.(goast.Fields); ok {
		d = append(d, f.EncodeToGoASTFields(obj)...)
	}
	for _, f := range d {
		fields = append(fields, &ast.Field{
			Names: f.Names,
			Type:  &ast.ArrayType{Elt: f.Type},
		})
	}
	obj.SetKind(goast.StructSeqofType)
	return
}
func (s *SimplePatternSeqof) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	gs := &goast.Seqof{Name: name}
	fields := []*ast.Field{}
	d := []*ast.Field{}
	if f, ok := s.Pattern.(goast.Fields); ok {
		d = append(fields, f.EncodeToGoASTFields(gs)...)
	}
	for _, f := range d {
		fields = append(fields, &ast.Field{
			Names: f.Names,
			Type:  &ast.ArrayType{Elt: f.Type},
		})
	}
	gs.Fields = fields

	if above != nil {
		above.Under(gs)
	} else {
		asts = append(asts, gs)
	}
	return
}

func (d *SimplePatternDictof) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	f := func(s SimplePattern) string {
		switch a := s.(type) {
		case *SimplePatternAny:
			return "any"
		case *SimplePatternAtom:
			if _, ok := a.AtomKind.(*AtomKindSymbol); ok {
				return "Symbol"
			} else {
				b := reflect.TypeOf(s).String()
				return strings.TrimPrefix(b, "*schema.")
			}
		case *SimplePatternLit:
		case *SimplePatternEmbedded:
		case *SimplePatternSeqof:
		case *SimplePatternSetof:
		case *SimplePatternDictof:
		case *SimplePatternRef:
			return SimplePatternRefToPreservesSchema(*a, "")
		}
		return "missing"
	}
	m := &goast.Map{Name: name, Key: f(d.Key), Value: f(d.Value)}
	if above != nil {
		above.Under(m)
	} else {
		asts = append(asts, m)
	}
	return
}

func (s *SimplePatternRef) ASTString() string {
	return s.Ref.ASTString()
}
func (s *SimplePatternRef) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	return s.Ref.EncodeToGoASTFields(obj)
}
func (s *SimplePatternRef) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	return s.Ref.EncodeToGoAST(above, name)
}

func (n *NamedSimplePatternNamed) ASTString() string {
	return n.Binding.ASTString()
}
func (n *NamedSimplePatternNamed) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	return n.Binding.EncodeToGoASTFields(obj)
}

func (n *NamedSimplePatternAnonymous) ASTString() string {
	if f, ok := n.SimplePattern.(goast.String); ok {
		return f.ASTString()
	}
	return ""
}
func (n *NamedSimplePatternAnonymous) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	if e, ok := n.SimplePattern.(goast.Fields); ok {
		return e.EncodeToGoASTFields(obj)
	}
	return []*ast.Field{}
}
func (n *NamedSimplePatternAnonymous) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	if e, ok := n.SimplePattern.(goast.Encoder); ok {
		asts = append(asts, e.EncodeToGoAST(above, name)...)
	}
	return
}

func (n *NamedPatternNamed) ASTString() string {
	return n.Binding.ASTString()
}
func (n *NamedPatternNamed) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	return n.Binding.EncodeToGoASTFields(obj)
}
func (n *NamedPatternNamed) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	u := &goast.Union{Name: name}
	if above != nil {
		above.Under(u)
	} else {
		asts = append(asts, u)
	}
	return n.Binding.EncodeToGoAST(u, name)
}

func (n *NamedPatternAnonymous) ASTString() string {
	if f, ok := n.Pattern.(goast.String); ok {
		return f.ASTString()
	}
	return ""
}
func (n *NamedPatternAnonymous) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	if e, ok := n.Pattern.(goast.Fields); ok {
		return e.EncodeToGoASTFields(obj)
	}
	return []*ast.Field{}
}

func (n *NamedPatternAnonymous) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	if f, ok := n.Pattern.(goast.Encoder); ok {
		return f.EncodeToGoAST(above, name)
	}
	return
}

func (o *DefinitionOr) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	u := &goast.Union{Name: name}

	for _, pattern := range append([]NamedAlternative{o.Pattern0, o.Pattern1}, o.PatternN...) {
		asts = append(asts, pattern.EncodeToGoAST(u, "")...)
	}

	if above != nil {
		above.Under(u)
	} else {
		asts = append(asts, u)
	}
	return
}

func (d *DefinitionPattern) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	if f, ok := d.Pattern.(goast.Encoder); ok {
		asts = append(asts, f.EncodeToGoAST(above, name)...)
	}
	return
}

func (r *CompoundPatternRec) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {

	s := &goast.Struct{Name: name}
	s.SetKind(goast.StructObjectType)
	s.SetStructKind(goast.StructRecType)
	// equals struct
	if enc, ok := r.Label.(goast.Encoder); ok {
		s.Identifier = enc.EncodeToGoAST(nil, "")
	}

	fields := []*ast.Field{}
	if f, ok := r.Fields.(goast.Fields); ok {
		fields = append(fields, f.EncodeToGoASTFields(s)...)
	}
	s.Fields = fields

	if above != nil {
		above.Under(s)
	} else {
		asts = append(asts, s)
	}
	return
}

func (t *CompoundPatternTuple) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {

	s := &goast.Tuple{Name: name}
	s.SetStructKind(goast.TupleType)
	// equals struct
	fields := []*ast.Field{}
	for _, p := range t.Patterns {
		if f, ok := p.(goast.Fields); ok {
			fields = append(fields, f.EncodeToGoASTFields(s)...)
		}
	}
	s.Fields = fields

	if above != nil {
		above.Under(s)
	} else {
		asts = append(asts, s)
	}
	return
}
func (t *CompoundPatternTuple) EncodeToGoASTFields(obj goast.AST) (fields []*ast.Field) {
	obj.SetStructKind(goast.StructTupleType)
	for _, pattern := range t.Patterns {
		if f, ok := pattern.(goast.Fields); ok {
			fields = append(fields, f.EncodeToGoASTFields(obj)...)
		}
	}
	return
}

func (t *CompoundPatternTuplePrefix) EncodeToGoASTFields(obj goast.AST) (fields []*ast.Field) {
	obj.SetStructKind(goast.StructTuplePrefixType)
	for _, pattern := range t.Fixed {
		if f, ok := pattern.(goast.Fields); ok {
			fields = append(fields, f.EncodeToGoASTFields(obj)...)
		}
	}
	if f, ok := t.Variable.(goast.Fields); ok {
		fields = append(fields, f.EncodeToGoASTFields(obj)...)
	}
	return
}
func (d *CompoundPatternDict) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	obj.SetStructKind(goast.StructDictType)
	return d.Entries.EncodeToGoASTFields(obj)
}

func (p *PatternSimplePattern) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	if e, ok := p.SimplePattern.(goast.Fields); ok {
		return e.EncodeToGoASTFields(obj)
	}
	return []*ast.Field{}
}
func (p *PatternSimplePattern) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	if f, ok := p.SimplePattern.(goast.Encoder); ok {
		return f.EncodeToGoAST(above, name)
	}
	return
}

func (p *PatternCompoundPattern) EncodeToGoASTFields(obj goast.AST) []*ast.Field {
	if e, ok := p.CompoundPattern.(goast.Fields); ok {
		return e.EncodeToGoASTFields(obj)
	}
	return []*ast.Field{}
}
func (p *PatternCompoundPattern) EncodeToGoAST(above goast.AST, name string) (asts []goast.AST) {
	if f, ok := p.CompoundPattern.(goast.Encoder); ok {
		asts = append(asts, f.EncodeToGoAST(above, name)...)
	}
	return
}
