package schema

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/isodude/preserves-go/lib/goast"
)

func DefinitionGenerator(name string, d Definition) goast.AST {
	g := goast.NewDefinition(name)
	switch u := d.(type) {
	case *DefinitionOr:
		DefinitionOrGenerator(g, u)
	case *DefinitionAnd:
		DefinitionAndGenerator(g, u)
	case *DefinitionPattern:
		DefinitionPatternGenerator(g, u)
	}

	return g
}

func DefinitionOrGenerator(g goast.AST, d *DefinitionOr) {
	NamedAlternativeGenerator(g, d.Pattern0)
	NamedAlternativeGenerator(g, d.Pattern1)
	for _, pattern := range d.PatternN {
		NamedAlternativeGenerator(g, pattern)
	}
	return
}

func DefinitionAndGenerator(g goast.AST, d *DefinitionAnd) {
	g.SetKind(goast.StructObjectType)
	NamedPatternGenerator(g, d.Pattern0)
	NamedPatternGenerator(g, d.Pattern1)
	for _, pattern := range d.PatternN {
		NamedPatternGenerator(g, pattern)
	}
}

func DefinitionPatternGenerator(g goast.AST, d *DefinitionPattern) {
	PatternGenerator(g, d.Pattern)
}

func NamedAlternativeGenerator(g goast.AST, n NamedAlternative) {
	s := goast.NewStruct(string(n.VariantLabel))
	g.Under(s)
	s.SetKind(goast.StructObjectType)
	PatternGenerator(s, n.Pattern)
}

func NamedPatternGenerator(g goast.AST, n NamedPattern) {
	switch u := n.(type) {
	case *NamedPatternNamed:
		NamedPatternNamedGenerator(g, u)
	case *NamedPatternAnonymous:
		NamedPatternAnonymousGenerator(g, u)
	}
}

func NamedPatternNamedGenerator(g goast.AST, n *NamedPatternNamed) {
	BindingGenerator(g, n.Binding)
}
func NamedPatternAnonymousGenerator(g goast.AST, n *NamedPatternAnonymous) {
	PatternGenerator(g, n.Pattern)
}

func BindingGenerator(g goast.AST, b Binding) {
	if u, ok := g.(*goast.Struct); ok {
		f := goast.NewField(string(b.Name))
		SimplePatternGenerator(f, b.Pattern)
		u.AddField(f)
	}
	if u, ok := g.(*goast.Definition); ok {
		f := goast.NewField(string(b.Name))
		SimplePatternGenerator(f, b.Pattern)
		u.AddField(f)
	}
}

func SimplePatternGenerator(g goast.AST, s SimplePattern) {
	g.SetStructKind(goast.StructRecType)
	switch u := s.(type) {
	case *SimplePatternAtom:
		AtomKindGenerator(g, u.AtomKind)
	case *SimplePatternLit:
		SimplePatternLitGenerator(g, u)
	case *SimplePatternRef:
		RefGenerator(g, &u.Ref)
	default:
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(s)))
	}
}

func AtomKindGenerator(g goast.AST, a AtomKind) {
	if u, ok := g.(*goast.Field); ok {
		s := strings.Split(fmt.Sprintf("%s", reflect.TypeOf(a)), ".")
		u.SetType(s[len(s)-1])
	}
}
func SimplePatternLitGenerator(g goast.AST, l *SimplePatternLit) {
	if u, ok := g.(*goast.Field); ok {
		u.SetValue(l.Value)
	}
	if u, ok := g.(*goast.Struct); ok {
		u.SetValue(l.Value)
	}
}
func RefGenerator(g goast.AST, r *Ref) {
	if u, ok := g.(*goast.Field); ok {
		module := ""
		if len(r.Module) > 0 {
			module = fmt.Sprintf("%s.", r.Module[len(r.Module)-1])
		}

		u.SetType(fmt.Sprintf("%s%s", module, r.Name))
	}
}

func CompoundPatternGenerator(g goast.AST, c CompoundPattern) {
	switch u := c.(type) {
	case *CompoundPatternRec:
		CompoundPatternRecGenerator(g, *u)
	case *CompoundPatternTuple:
		CompoundPatternTupleGenerator(g, *u)
	case *CompoundPatternTuplePrefix:
		CompoundPatternTuplePrefixGenerator(g, *u)
	case *CompoundPatternDict:
		CompoundPatternDictGenerator(g, *u)
	default:
		fmt.Printf("debug: %v\n", reflect.TypeOf(u))
	}
}

func CompoundPatternRecGenerator(g goast.AST, c CompoundPatternRec) {}
func CompoundPatternTupleGenerator(g goast.AST, c CompoundPatternTuple) {
	g.SetStructKind(goast.StructTupleType)
	for _, pattern := range c.Patterns {
		NamedPatternGenerator(g, pattern)
	}
}
func CompoundPatternTuplePrefixGenerator(g goast.AST, c CompoundPatternTuplePrefix) {}
func CompoundPatternDictGenerator(g goast.AST, c CompoundPatternDict)               {}

func PatternGenerator(g goast.AST, p Pattern) {
	switch u := p.(type) {
	case *PatternSimplePattern:
		SimplePatternGenerator(g, u.SimplePattern)
	case *PatternCompoundPattern:
		CompoundPatternGenerator(g, u.CompoundPattern)
	default:
		fmt.Printf("debug: %v\n", reflect.TypeOf(u))
	}
}
