package schema

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/isodude/preserves-go/lib/goast"
	. "github.com/isodude/preserves-go/lib/preserves"
)

func BundleGenerator(name string, b Bundle) *goast.Bundle {
	g := goast.NewBundle(name)
	ModulesGenerator(g, b.Modules)
	return g
}

func ModulesGenerator(b *goast.Bundle, m Modules) {
	for path, schema := range m {
		SchemaGenerator(b, path.FromHash(), schema)
	}
}

func SchemaGenerator(b *goast.Bundle, mp ModulePath, m Schema) {
	g := goast.NewSchema()
	DefinitionsGenerator(g, m.Definitions)
	VersionGenerator(g, m.Version)
	EmbeddedTypeGenerator(g, m.EmbeddedType)

	var path []string
	for _, s := range mp {
		path = append(path, string(s))
	}
	b.AddSchema(path, g)
}

func VersionGenerator(s *goast.Schema, v Version) {
	f := goast.NewField("")
	f.SetValue(VersionToPreserves(v))
	s.Version = f
}

func EmbeddedTypeGenerator(s *goast.Schema, e EmbeddedTypeName) {
	f := goast.NewField("")
	switch u := e.(type) {
	case *EmbeddedTypeNameFalse:
		f.SetValue(NewSymbol("false"))
	case *EmbeddedTypeNameRef:
		RefGenerator(f, &u.Ref)
	default:
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(e)))
	}
	s.EmbeddedType = f
}

func DefinitionsGenerator(s *goast.Schema, d Definitions) {
	for name, definition := range d {
		s.AddDefinition(DefinitionGenerator(string(name), definition))
	}
}

func DefinitionGenerator(name string, d Definition) *goast.Definition {
	g := goast.NewDefinition(name)
	switch u := d.(type) {
	case *DefinitionOr:
		DefinitionOrGenerator(g, u)
	case *DefinitionAnd:
		DefinitionAndGenerator(g, u)
	case *DefinitionPattern:
		DefinitionPatternGenerator(g, u)
	default:
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(d)))
	}

	return g
}

func DefinitionOrGenerator(g goast.AST, d *DefinitionOr) {
	g.SetKind(goast.InterfaceObjectType)
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
	s := goast.NewDefinition(string(n.VariantLabel))
	g.Under(s)
	s.SetKind(goast.StructObjectType)
	if u, ok := n.Pattern.(*DefinitionPattern); ok {
		PatternGenerator(s, u.Pattern)
	} else if u, ok := n.Pattern.(*NamedPatternAnonymous); ok {
		PatternGenerator(s, u.Pattern)
	} else {
		PatternGenerator(s, n.Pattern)
	}
}

func NamedPatternGenerator(g goast.AST, n NamedPattern) {
	switch u := n.(type) {
	case *NamedPatternNamed:
		NamedPatternNamedGenerator(g, u)
	case *NamedPatternAnonymous:
		NamedPatternAnonymousGenerator(g, u)
	default:
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(n)))
	}
}

func NamedSimplePatternGenerator(g goast.AST, n NamedSimplePattern) {
	switch u := n.(type) {
	case *NamedSimplePatternNamed:
		NamedSimplePatternNamedGenerator(g, u)
	case *NamedSimplePatternAnonymous:
		NamedSimplePatternAnonymousGenerator(g, u)
	default:
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(n)))
	}
}

func NamedSimplePatternNamedGenerator(g goast.AST, n *NamedSimplePatternNamed) {
	BindingGenerator(g, n.Binding)
}
func NamedPatternNamedGenerator(g goast.AST, n *NamedPatternNamed) {
	if u, ok := g.(*goast.Field); ok {
		BindingGenerator(u, n.Binding)
		return
	}
	f := &goast.Field{}
	BindingGenerator(f, n.Binding)

	if u, ok := g.(*goast.Definition); ok {
		u.AddField(f)
	} else {
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
	}
}
func NamedSimplePatternAnonymousGenerator(g goast.AST, n *NamedSimplePatternAnonymous) {
	SimplePatternGenerator(g, n.SimplePattern)
}
func NamedPatternAnonymousGenerator(g goast.AST, n *NamedPatternAnonymous) {
	PatternGenerator(g, n.Pattern)
}

func BindingGenerator(g goast.AST, b Binding) {
	if u, ok := g.(*goast.Definition); ok {
		f := goast.NewField(string(b.Name))
		SimplePatternGenerator(f, b.Pattern)
		u.AddField(f)
	} else if u, ok := g.(*goast.Field); ok {
		u.SetName(string(b.Name))
		SimplePatternGenerator(u, b.Pattern)
	} else {
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
	}
}

func SimplePatternGenerator(g goast.AST, s SimplePattern) {
	switch u := s.(type) {
	case *SimplePatternAny:
		SimplePatternAnyGenerator(g, *u)
	case *SimplePatternAtom:
		AtomKindGenerator(g, u.AtomKind)
	case *SimplePatternEmbedded:
		SimplePatternEmbeddedGenerator(g, u.Interface)
	case *SimplePatternLit:
		SimplePatternLitGenerator(g, u)
	case *SimplePatternSeqof:
		SimplePatternSeqofGenerator(g, u.Pattern)
	case *SimplePatternSetof:
		SimplePatternSetofGenerator(g, u.Pattern)
	case *SimplePatternDictof:
		SimplePatternDictofGenerator(g, *u)
	case *SimplePatternRef:
		RefGenerator(g, &u.Ref)
	default:
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(s)))
	}
}

func SimplePatternAnyGenerator(g goast.AST, s SimplePatternAny) {
	if u, ok := g.(*goast.Field); ok {
		u.SetType("Any")
	} else {
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
	}
}

func SimplePatternEmbeddedGenerator(g goast.AST, s SimplePattern) {
	g.SetStructKind(goast.StructRecType)
	panic("not implemented")
	// TODO
}

func SimplePatternSeqofGenerator(g goast.AST, s SimplePattern) {
	g.SetKind(goast.StructObjectType)
	g.SetStructKind(goast.StructSeqofType)
	if u, ok := g.(*goast.Field); ok {
		u.SetArray()
		SimplePatternGenerator(g, s)
	} else if u, ok := g.(*goast.Definition); ok {
		f := &goast.Field{}
		f.SetArray()
		SimplePatternGenerator(f, s)
		u.AddField(f)
	} else {
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
	}
}
func SimplePatternSetofGenerator(g goast.AST, s SimplePattern) {
	// TODO
	panic("not implemented")
	g.SetStructKind(goast.StructDictType)
	if u, ok := g.(*goast.Field); ok {
		u.SetKind(goast.SetField)
		SimplePatternGenerator(g, s)
	} else {
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
	}
}
func SimplePatternDictofGenerator(g goast.AST, s SimplePatternDictof) {
	g.SetStructKind(goast.StructRecType)
	if u, ok := g.(*goast.Field); ok {
		u.SetKind(goast.DictField)
		key := &goast.Field{}
		SimplePatternGenerator(key, s.Key)
		value := &goast.Field{}
		SimplePatternGenerator(value, s.Value)
		u.SetDict(key, value)
	} else if u, ok := g.(*goast.Definition); ok {
		f := &goast.Field{}
		f.SetKind(goast.DictField)
		key := &goast.Field{}
		SimplePatternGenerator(key, s.Key)
		value := &goast.Field{}
		SimplePatternGenerator(value, s.Value)
		f.SetDict(key, value)
		u.AddField(f)
	} else {
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
	}
}
func AtomKindGenerator(g goast.AST, a AtomKind) {
	if u, ok := g.(*goast.Field); ok {
		s := strings.Split(fmt.Sprintf("%s", reflect.TypeOf(a)), ".AtomKind")
		u.SetType(s[len(s)-1])
	} else {
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
	}
}
func SimplePatternLitGenerator(g goast.AST, l *SimplePatternLit) {
	if u, ok := g.(*goast.Field); ok {
		u.SetValue(l.Value)
	} else if u, ok := g.(*goast.Label); ok {
		u.SetStmt(goast.NewLit(l.Value))
	} else if u, ok := g.(*goast.Definition); ok {
		u.SetKind(goast.StructObjectType)
		u.SetValue(l.Value)
	} else {
		panic(fmt.Sprintf("unknown type %v, %v", reflect.TypeOf(g), g))
	}
}
func RefGenerator(g goast.AST, r *Ref) {
	if u, ok := g.(*goast.Field); ok {
		module := ""
		if len(r.Module) > 0 {
			module = fmt.Sprintf("%s.", r.Module[len(r.Module)-1])
		}

		u.SetType(fmt.Sprintf("%s%s", module, r.Name))
	} else if u, ok := g.(*goast.Definition); ok {
		u.SetKind(goast.PassthroughObjectType)
		u.SetStructKind(goast.StructTupleType)
		f := goast.NewField("")
		module := ""
		if len(r.Module) > 0 {
			module = fmt.Sprintf("%s.", r.Module[len(r.Module)-1])
		}

		f.SetType(fmt.Sprintf("%s%s", module, r.Name))
		u.AddField(f)
	} else {
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
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
		panic(fmt.Sprintf("debug: %v\n", reflect.TypeOf(c)))
	}
}

func CompoundPatternRecGenerator(g goast.AST, c CompoundPatternRec) {
	g.SetKind(goast.StructObjectType)
	g.SetStructKind(goast.StructRecType)
	l := goast.NewLabel()
	NamedPatternGenerator(l, c.Label)
	if u, ok := g.(*goast.Definition); ok {
		u.Identifier = l
	} else {
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
	}
	NamedPatternGenerator(g, c.Fields)
}
func CompoundPatternTupleGenerator(g goast.AST, c CompoundPatternTuple) {
	g.SetKind(goast.StructObjectType)
	if u, ok := g.(*goast.Definition); ok {
		if u.StructKind != nil && *u.StructKind == goast.StructRecType {
			u.SetStructKind(goast.StructTupleType)
		} else {
			u.SetStructKind(goast.TupleType)
		}
	} else {
		panic("should not reach here")
	}
	for _, pattern := range c.Patterns {
		NamedPatternGenerator(g, pattern)
	}
}
func CompoundPatternTuplePrefixGenerator(g goast.AST, c CompoundPatternTuplePrefix) {
	g.SetKind(goast.StructObjectType)
	g.SetStructKind(goast.StructTuplePrefixType)
	for _, pattern := range c.Fixed {
		if u, ok := g.(*goast.Field); ok {
			NamedPatternGenerator(u, pattern)
			continue
		}
		f := &goast.Field{}
		NamedPatternGenerator(f, pattern)
		if u, ok := g.(*goast.Definition); ok {
			u.AddField(f)
		} else {
			panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
		}
	}
	if u, ok := g.(*goast.Field); ok {
		NamedSimplePatternGenerator(u, c.Variable)
		return
	}
	f := &goast.Field{}
	f.SetArray()
	NamedSimplePatternGenerator(f, c.Variable)
	if u, ok := g.(*goast.Definition); ok {
		u.AddField(f)
	} else {
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
	}
}
func CompoundPatternDictGenerator(g goast.AST, c CompoundPatternDict) {
	g.SetKind(goast.StructObjectType)
	g.SetStructKind(goast.StructDictType)
	DictionaryEntriesGenerator(g, c.Entries)
}
func DictionaryEntriesGenerator(g goast.AST, d DictionaryEntries) {
	var keys []Value
	for k := range d {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(e1, e2 Value) int {
		return e1.Cmp(e2)
	})
	for _, k := range keys {
		v := d[k]
		if u, ok := g.(*goast.Field); ok {
			NamedSimplePatternGenerator(u, v)
			fmt.Printf("%v\n", g)
			panic("dump")
			// TODO: WRONG
			continue
		}
		f := &goast.Field{}
		NamedSimplePatternGenerator(f, v)
		if u, ok := g.(*goast.Definition); ok {
			u.AddField(f)
			u.MapKeyToField = append(u.MapKeyToField, k)
		} else {
			panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(g)))
		}

	}

}

func PatternGenerator(g goast.AST, p Pattern) {
	switch u := p.(type) {
	case *PatternSimplePattern:
		SimplePatternGenerator(g, u.SimplePattern)
	case *PatternCompoundPattern:
		CompoundPatternGenerator(g, u.CompoundPattern)
	default:
		panic(fmt.Sprintf("debug: %v\n", reflect.TypeOf(p)))
	}
}
