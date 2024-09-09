package schema

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/isodude/preserves-go/lib/extras"
	"github.com/isodude/preserves-go/lib/preserves/text"

	. "github.com/isodude/preserves-go/lib/preserves"
)

func FromPreserves(file string) (*Schema, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	buf = bufio.NewReaderSize(buf, 200000)
	tp := text.TextParser{}
	_, err = tp.ReadFrom(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return SchemaFromPreserves(text.ToPreserves(tp.GetValue())), nil
}

func FromPreservesSchemaFile(file string) (*Schema, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	buf = bufio.NewReaderSize(buf, 200000)
	tp := text.TextParser{}
	schema := Schema{EmbeddedType: &EmbeddedTypeNameFalse{}}
	var values []Value
	var beforeValues []Value
	var insideDefinition bool
	for err == nil {
		_, err = tp.ReadFrom(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {

			return &schema, nil
		}
		value := text.ToPreserves(tp.GetValue())
	retry:
		switch a := value.(type) {
		case *Annotation:
			if !insideDefinition {
				value = a.AnnotatedValue
				goto retry
			}
			values = append(values, value)
		case *Comment:
			value = a.AnnotatedValue
			goto retry
		case *Symbol:
			if a.String() == "=" {
				insideDefinition = true
			}
			if a.String() == "." {
				insideDefinition = false
				if s := SchemaFromPreservesSchema(schema, values); s == nil {
					return nil, fmt.Errorf("Schema line (%v) could not be parsed, line before: %s", values, beforeValues)
				} else {
					schema = *s
				}
				// empty the slice
				beforeValues = slices.Clone(values)
				values = values[len(values):]
			} else {
				values = append(values, value)
			}
		default:
			values = append(values, value)
		}
	}
	if err != io.EOF {
		return nil, err
	}
	return &schema, nil
}

func AtomKindFromPreservesSchema(values []Value) AtomKind {
	if len(values) != 1 {
		return nil
	}
	u, ok := values[0].(*Symbol)
	if !ok {
		return nil
	}
	switch u.String() {
	case "bool":
		return &AtomKindBoolean{}
	case "double":
		return &AtomKindDouble{}
	case "int":
		return &AtomKindSignedInteger{}
	case "string":
		return &AtomKindString{}
	case "bytes":
		return &AtomKindByteString{}
	case "symbol":
		return &AtomKindSymbol{}
	default:
		return nil
	}
}
func AtomKindToPreservesSchema(p AtomKind, _ string) string {
	switch p.(type) {
	case *AtomKindBoolean:
		return "bool"
	case *AtomKindDouble:
		return "double"
	case *AtomKindSignedInteger:
		return "int"
	case *AtomKindString:
		return "string"
	case *AtomKindByteString:
		return "bytes"
	case *AtomKindSymbol:
		return "symbol"
	}

	return ""
}

func BindingFromPreservesSchema(values []Value) *Binding {
	if len(values) != 1 {
		return nil
	}
	a, ok := values[0].(*Annotation)
	if !ok {
		return nil
	}
	s, ok := a.Value.(*Symbol)
	if !ok {
		return nil
	}
	name := NewSymbol(s.String())
	u := SimplePatternFromPreservesSchema([]Value{a.AnnotatedValue})
	if u == nil {
		return nil
	}
	pattern, ok := u.(SimplePattern)
	if !ok {
		return nil
	}
	return &Binding{Name: *name, Pattern: pattern}
}
func BindingToPreservesSchema(p Binding, indent string) string {
	return fmt.Sprintf("@%s %s", SymbolToPreservesSchema(p.Name, indent), SimplePatternToPreservesSchema(p.Pattern, indent))
}

func BundleFromPreservesSchema(values []Value) *Bundle {
	return nil
}
func BundleToPreservesSchema(p Bundle, indent string) string {
	return fmt.Sprintf("<bundle @modules %s>", ModulesToPreservesSchema(p.Modules, indent))
}

func (d Definitions) Add(k Symbol, v Definition) Definitions {
	d[k] = v
	return d
}

func DefinitionsFromPreservesSchema(d Definitions, values []Value) *Definitions {
	symbol := SymbolFromPreservesSchema(values[0:1])
	if symbol == nil {
		return nil
	}
	definition := DefinitionFromPreservesSchema(values[2:])
	if definition == nil {
		return nil
	}
	d[*symbol] = definition
	return &d
}
func DefinitionsToPreservesSchema(p Definitions, indent string) string {
	var s string
	var keys []Symbol
	for k := range p {
		keys = append(keys, k)
	}
	slices.Sort[[]Symbol](keys)
	for _, k := range keys {
		def := DefinitionToPreservesSchema(p[k], indent)
		space := " "
		if strings.ContainsAny(string(def[len(def)-1]), ">}]") {
			space = ""
		}
		s = fmt.Sprintf("%s\n%s = %s%s.", s, SymbolToPreservesSchema(k, indent), def, space)
	}
	return s
}

func (d DictionaryEntries) Add(value Value, namedSimplePattern NamedSimplePattern) DictionaryEntries {
	d[value] = namedSimplePattern
	return d
}

func DictionaryEntriesFromPreservesSchema(d DictionaryEntries, values []Value) *DictionaryEntries {
	if len(values) != 1 {
		return nil
	}
	dictionary, ok := values[0].(*Dictionary)
	if !ok {
		return nil
	}
	for k, v := range *dictionary {
		key := ValueFromPreservesSchema([]Value{k})
		if key == nil {
			return nil
		}

		var _key Value
		_key = key
		if a, ok := v.(*Annotation); ok {
			_key = a.Value
			v = a.AnnotatedValue
		}
		namedSimplePattern := NamedSimplePatternFromPreservesSchema([]Value{v})
		if namedSimplePattern == nil {
			return nil
		}
		var value NamedSimplePattern
		if simplePattern, ok := namedSimplePattern.(*NamedSimplePatternAnonymous); ok {
			if ref, ok := simplePattern.SimplePattern.(*SimplePatternRef); ok {
				if s, ok := _key.(*Symbol); ok {
					value = &NamedSimplePatternNamed{*NewBinding(*NewSymbol((*s).String()), ref)}
				}
			}
		}

		if value == nil {
			return nil
		}
		d.Add(key, value)
	}
	if len(d) == 0 {
		return nil
	}
	return &d
}
func DictionaryEntriesToPreservesSchema(p DictionaryEntries, indent string) string {
	var s string
	var keys []string
	reverse := make(map[string]Value)
	for k := range p {
		ks := ValueToPreservesSchema(k, "")
		reverse[ks] = k
		keys = append(keys, ks)
	}
	slices.Sort[[]string](keys)
	for _, ks := range keys {
		k := reverse[ks]
		v := p[k]
		s = fmt.Sprintf("%s\n%s  %s: %s", s, indent, ValueToPreservesSchema(k, fmt.Sprintf("%s  ", indent)), NamedSimplePatternToPreservesSchema(v, fmt.Sprintf("%s  ", indent)))
	}
	if len(s) > 0 {
		return fmt.Sprintf("{%s\n%s}", s, indent)
	}
	return "{}"
}

func ModulePathFromPreservesSchema(values []Value) *ModulePath {
	if len(values) == 1 {
		if s, ok := values[0].(*Symbol); ok {
			p := strings.Split(s.String(), ".")
			if len(p) > 0 {
				var symbols []Symbol
				for _, k := range p {
					symbols = append(symbols, *NewSymbol(k))

				}
				m := ModulePath(symbols)
				return &m
			}
		}
	}
	return nil
}
func ModulePathToPreservesSchema(p ModulePath, indent string) (s string) {
	for i, k := range p {
		s += SymbolToPreservesSchema(k, indent)
		if len(p) == i+1 {
			s += "."
		}
	}
	return
}

func (m Modules) Add(modulePath ModulePath, schema Schema) Modules {
	m[modulePath.ToHash()] = schema
	return m
}
func ModulesFromPreservesSchema(m Modules, values []Value) *Modules {
	if len(values) != 1 {
		return nil
	}
	dictionary, ok := values[0].(*Dictionary)
	if !ok {
		return nil
	}
	for k, v := range *dictionary {
		modulePath := ModulePathFromPreservesSchema([]Value{k})
		if modulePath == nil {
			return nil
		}
		schema := SchemaFromPreservesSchema(Schema{}, []Value{v})
		if schema == nil {
			return nil
		}
		m[modulePath.ToHash()] = *schema
	}
	if len(m) == 0 {
		return nil
	}
	return &m
}
func ModulesToPreservesSchema(p Modules, indent string) string {
	var s string
	var keys []string
	reverse := make(map[string]extras.Hash[ModulePath])
	for k := range p {
		var ks string
		reverse[ks] = k
		keys = append(keys, ks)
	}
	slices.Sort[[]string](keys)
	for _, ks := range keys {
		k := reverse[ks]
		v := (p)[k]
		m := k.FromHash()
		s = fmt.Sprintf("%s\n%s%s: %s", s, indent, ModulePathToPreservesSchema(m, fmt.Sprintf("%s  ", indent)), SchemaToPreservesSchema(v, fmt.Sprintf("%s  ", indent)))
	}
	if len(s) > 0 {
		return fmt.Sprintf("{%s\n%s}", s, indent)
	}
	return "{}"
}

func RefFromPreservesSchema(values []Value) *Ref {
	if len(values) != 1 {
		return nil
	}
	s, ok := values[0].(*Symbol)
	if !ok {
		return nil
	}
	var path *Symbol
	var name *Symbol
	if i := strings.LastIndex(s.String(), "."); i > 0 {
		name = NewSymbol(s.String()[i:])
		path = NewSymbol(s.String()[:i])
	} else {
		name = NewSymbol(s.String())
	}
	var modulePath ModulePath
	if path != nil && len(*path) > 0 {
		modulePath = *ModulePathFromPreservesSchema([]Value{path})
		if modulePath == nil {
			return nil
		}
	}
	return &Ref{Module: modulePath, Name: *name}
}
func RefToPreservesSchema(p Ref, indent string) string {
	if len(p.Module) > 0 {
		return fmt.Sprintf("%s%s", ModulePathToPreservesSchema(p.Module, indent), SymbolToPreservesSchema(p.Name, indent))
	}
	return SymbolToPreservesSchema(p.Name, indent)
}

func SchemaFromPreservesSchema(s Schema, values []Value) *Schema {
	if len(values) < 1 {
		return nil
	}
	v, ok := values[0].(*Symbol)
	if !ok {
		return nil
	}
	if len(values) == 2 {
		if v.String() == "version" {
			version := VersionFromPreservesSchema(values)
			if version == nil {
				return nil
			}
			s.Version = *version
			return &s
		}
		if v.String() == "embeddedType" {
			embeddedTypeName := EmbeddedTypeNameFromPreservesSchema(values[1:])
			if embeddedTypeName == nil {
				return nil
			}
			s.EmbeddedType = embeddedTypeName
			return &s
		}
		return nil
	}

	e, ok := values[1].(*Symbol)
	if !ok {
		return nil
	}
	if e.String() != "=" {
		return nil
	}
	if s.Definitions == nil {
		s.Definitions = NewDefinitions()
	}
	d := DefinitionsFromPreservesSchema(s.Definitions, values)
	if d == nil {
		return nil
	}

	return &s

}
func SchemaToPreservesSchema(p Schema, indent string) string {
	return fmt.Sprintf("embeddedType %s .\nversion %s .\n%s",
		EmbeddedTypeNameToPreservesSchema(p.EmbeddedType, indent),
		VersionToPreservesSchema(p.Version, indent),
		DefinitionsToPreservesSchema(p.Definitions, indent))
}

func VersionFromPreservesSchema(values []Value) *Version {
	if u, ok := values[1].(*SignedInteger); ok {
		i := big.Int(*u)
		if i.Int64() == 1 {
			return &Version{}
		}
	}
	return nil
}
func VersionToPreservesSchema(_ Version, _ string) string {
	return "1"
}

func NamedAlternativeFromPreservesSchema(values []Value) *NamedAlternative {
	if len(values) != 1 {
		return nil
	}

	var label Pstring
	if value, ok := values[0].(*Annotation); ok {
		if symbol, ok := value.Value.(*Symbol); ok {
			label = *NewPstring(symbol.String())
			values = []Value{value.AnnotatedValue}
		}
	}

	pattern := PatternFromPreservesSchema(values)
	if pattern == nil {
		return nil
	}
	if label == "" {
		switch p := pattern.(type) {
		case *PatternCompoundPattern:
			switch q := p.CompoundPattern.(type) {
			case *CompoundPatternRec:
				l, ok := q.Label.(*NamedPatternAnonymous).Pattern.(*PatternSimplePattern).SimplePattern.(*SimplePatternLit)
				if !ok {
					return nil
				}
				s, ok := l.Value.(*Symbol)
				if !ok {
					return nil
				}
				label = *NewPstring(string(*s))
			}
		case *PatternSimplePattern:
			switch q := p.SimplePattern.(type) {
			case *SimplePatternRef:
				label = *NewPstring(string(q.Ref.Name))
			case *SimplePatternLit:
				switch o := q.Value.(type) {
				case *Boolean:
					if *o {
						label = "true"
					} else {
						label = "false"
					}
				case *Symbol:
					label = *NewPstring(o.String())
				default:
					label = *NewPstring(SimplePatternLitToPreservesSchema(*q, ""))
				}

			case *SimplePatternAny:
				label = "any"
			case *SimplePatternAtom:
				label = *NewPstring(SimplePatternAtomToPreservesSchema(*q, ""))
			default:
				panic(fmt.Sprintf("%s: %s", p, reflect.TypeOf(q)))
			}
		default:
			panic(reflect.TypeOf(p))
		}
	}
	return &NamedAlternative{VariantLabel: label, Pattern: pattern}
}

func NamedAlternativeToPreservesSchema(p NamedAlternative, indent string) string {
	s := PatternToPreservesSchema(p.Pattern, indent)
	if string(p.VariantLabel) != strings.TrimPrefix(s, "=") && !slices.Contains([]byte{'[', '<', '{'}, s[0]) {
		s = fmt.Sprintf("@%s %s", p.VariantLabel, s)
	}
	return s
}

func SimplePatternFromPreservesSchema(values []Value) SimplePattern {
	if a := SimplePatternAnyFromPreservesSchema(values); a != nil {
		return a
	}
	if a := SimplePatternAtomFromPreservesSchema(values); a != nil {
		return a
	}
	if a := SimplePatternEmbeddedFromPreservesSchema(values); a != nil {
		return a
	}
	if a := SimplePatternLitFromPreservesSchema(values); a != nil {
		return a
	}
	if a := SimplePatternSeqofFromPreservesSchema(values); a != nil {
		return a
	}
	if a := SimplePatternSetofFromPreservesSchema(values); a != nil {
		return a
	}
	if a := SimplePatternDictofFromPreservesSchema(values); a != nil {
		return a
	}
	if a := SimplePatternRefFromPreservesSchema(values); a != nil {
		return a
	}
	return nil
}
func SimplePatternToPreservesSchema(p SimplePattern, indent string) string {
	if a, ok := p.(*SimplePatternAny); ok {
		return SimplePatternAnyToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*SimplePatternAtom); ok {
		return SimplePatternAtomToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*SimplePatternEmbedded); ok {
		return SimplePatternEmbeddedToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*SimplePatternLit); ok {
		return SimplePatternLitToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*SimplePatternSeqof); ok {
		return SimplePatternSeqofToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*SimplePatternSetof); ok {
		return SimplePatternSetofToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*SimplePatternDictof); ok {
		return SimplePatternDictofToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*SimplePatternRef); ok {
		return SimplePatternRefToPreservesSchema(*a, indent)
	}
	return ""
}

func SimplePatternAnyFromPreservesSchema(values []Value) *SimplePatternAny {
	if len(values) == 1 {
		if u, ok := values[0].(*Symbol); ok {
			if u.String() == "any" {
				return &SimplePatternAny{}
			}
		}
	}
	return nil
}
func SimplePatternAnyToPreservesSchema(p SimplePatternAny, indent string) string {
	return "any"
}

func SimplePatternAtomFromPreservesSchema(values []Value) *SimplePatternAtom {
	atomKind := AtomKindFromPreservesSchema(values)
	if atomKind == nil {
		return nil
	}
	return &SimplePatternAtom{AtomKind: atomKind}
}

func SimplePatternAtomToPreservesSchema(p SimplePatternAtom, indent string) string {
	return fmt.Sprintf("%s", AtomKindToPreservesSchema(p.AtomKind, indent))
}

func SimplePatternEmbeddedFromPreservesSchema(values []Value) *SimplePatternEmbedded {
	if len(values) == 1 {
		if v, ok := values[0].(*Embedded); ok {
			if p := SimplePatternFromPreservesSchema([]Value{v.Value}); p == nil {
				if iface, ok := p.(SimplePattern); ok {
					return &SimplePatternEmbedded{Interface: iface}
				}
			}
		}
	}
	return nil
}

func SimplePatternEmbeddedToPreservesSchema(p SimplePatternEmbedded, indent string) string {
	return fmt.Sprintf("<embedded %s>", SimplePatternToPreservesSchema(p.Interface, indent))
}

func SimplePatternLitFromPreservesSchema(values []Value) *SimplePatternLit {
	if len(values) != 1 {
		return nil
	}
	switch u := values[0].(type) {
	case *Symbol:
		if strings.HasPrefix(u.String(), "=") {
			return &SimplePatternLit{Value: NewSymbol(strings.TrimPrefix(u.String(), "="))}
		}
		if u.String() == "symbol" {
			return &SimplePatternLit{Value: NewSymbol(u.String())}
		}
		return nil
	case *SignedInteger:
	case *Pstring:
	case *Double:
	case *Boolean:
	case *ByteString:
	default:
		return nil
	}
	// TODO howto do this?
	//if u := AtomKindFromPreservesSchema(values[0:1]); u != nil {
	//if a, ok := u.(Any); ok {
	//	return &Lit{Value: a}
	//}
	//}
	return &SimplePatternLit{Value: values[0]}
}

func SimplePatternLitToPreservesSchema(p SimplePatternLit, indent string) string {
	var b bytes.Buffer
	text.FromPreserves(p.Value).WriteTo(&b)
	re := regexp.MustCompile("^([A-Z].+|any$)")
	if re.MatchString(b.String()) {
		return fmt.Sprintf("=%s", b.String())
	}
	return b.String()
}

func SimplePatternSeqofFromPreservesSchema(values []Value) *SimplePatternSeqof {
	if len(values) != 1 {
		return nil
	}
	sequence, ok := values[0].(*Sequence)
	if !ok {
		return nil
	}
	if len(*sequence) != 2 {
		return nil
	}
	if dotdotdot, ok := (*sequence)[1].(*Symbol); ok {
		if dotdotdot.String() != "..." {
			return nil
		}
	}
	u := SimplePatternFromPreservesSchema([]Value{(*sequence)[0]})
	if u == nil {
		return nil
	}
	pattern, ok := u.(SimplePattern)
	if !ok {
		return nil
	}
	return &SimplePatternSeqof{Pattern: pattern}
}

func SimplePatternSeqofToPreservesSchema(p SimplePatternSeqof, indent string) string {
	return fmt.Sprintf("[%s ...]", SimplePatternToPreservesSchema(p.Pattern, indent))
}

func SimplePatternSetofFromPreservesSchema(values []Value) *SimplePatternSetof {
	if len(values) != 1 {
		return nil
	}
	set, ok := values[0].(*Set)
	if !ok {
		return nil
	}
	if len(*set) != 1 {
		return nil
	}
	var items []Value
	for k := range *set {
		items = append(items, k)
	}
	u := SimplePatternFromPreservesSchema(items)
	if u == nil {
		return nil
	}
	pattern, ok := u.(SimplePattern)
	if !ok {
		return nil
	}
	return &SimplePatternSetof{Pattern: pattern}
}

func SimplePatternSetofToPreservesSchema(p SimplePatternSetof, indent string) string {
	return fmt.Sprintf("#{%s ...}", SimplePatternToPreservesSchema(p.Pattern, indent))
}

func SimplePatternDictofFromPreservesSchema(values []Value) *SimplePatternDictof {
	if len(values) != 1 {
		return nil
	}
	dictionary, ok := values[0].(*Dictionary)
	if !ok {
		return nil
	}
	if len(*dictionary) != 2 {
		return nil
	}
	if dotdotdot, ok := dictionary.Get(NewSymbol("...")); ok {
		d, ok := dotdotdot.(*Symbol)
		if !ok {
			return nil
		}
		if d.String() != "..." {
			return nil
		}

		dictionary.Delete(NewSymbol("..."))
	} else {
		return nil
	}
	for k, v := range *dictionary {

		s := SimplePatternFromPreservesSchema([]Value{k})
		if s == nil {
			return nil
		}

		key, ok := s.(SimplePattern)
		if !ok {
			return nil
		}

		s = SimplePatternFromPreservesSchema([]Value{v})
		if s == nil {
			return nil
		}
		value, ok := s.(SimplePattern)
		if !ok {
			return nil
		}
		return &SimplePatternDictof{Key: key, Value: value}
	}
	return nil
}

func SimplePatternDictofToPreservesSchema(p SimplePatternDictof, indent string) string {
	return fmt.Sprintf("{ %s: %s ...:... }", SimplePatternToPreservesSchema(p.Key, indent), SimplePatternToPreservesSchema(p.Value, indent))
}

func SimplePatternRefFromPreservesSchema(values []Value) *SimplePatternRef {
	if a := RefFromPreservesSchema(values); a != nil {
		return &SimplePatternRef{*a}
	}
	return nil
}

func SimplePatternRefToPreservesSchema(p SimplePatternRef, indent string) string {
	return RefToPreservesSchema(p.Ref, indent)
}

func NamedSimplePatternFromPreservesSchema(values []Value) NamedSimplePattern {
	if a := NamedSimplePatternNamedFromPreservesSchema(values); a != nil {
		return a
	}
	if a := NamedSimplePatternAnonymousFromPreservesSchema(values); a != nil {
		return a
	}
	return nil
}

func NamedSimplePatternToPreservesSchema(p NamedSimplePattern, indent string) string {
	if a, ok := p.(*NamedSimplePatternNamed); ok {
		return NamedSimplePatternNamedToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*NamedSimplePatternAnonymous); ok {
		return NamedSimplePatternAnonymousToPreservesSchema(*a, indent)
	}
	return ""
}

func NamedSimplePatternNamedFromPreservesSchema(values []Value) *NamedSimplePatternNamed {
	if a := BindingFromPreservesSchema(values); a != nil {
		return &NamedSimplePatternNamed{*a}
	}

	return nil
}
func NamedSimplePatternNamedToPreservesSchema(p NamedSimplePatternNamed, indent string) string {

	return BindingToPreservesSchema(p.Binding, indent)

}
func NamedSimplePatternAnonymousFromPreservesSchema(values []Value) *NamedSimplePatternAnonymous {
	if a := SimplePatternFromPreservesSchema(values); a != nil {
		return &NamedSimplePatternAnonymous{a}
	}

	return nil
}
func NamedSimplePatternAnonymousToPreservesSchema(p NamedSimplePatternAnonymous, indent string) string {
	return SimplePatternToPreservesSchema(p.SimplePattern, indent)
}

func NamedPatternFromPreservesSchema(values []Value) NamedPattern {
	if a := NamedPatternNamedFromPreservesSchema(values); a != nil {
		return a
	}
	if a := NamedPatternAnonymousFromPreservesSchema(values); a != nil {
		return a
	}

	return nil
}
func NamedPatternToPreservesSchema(p NamedPattern, indent string) string {
	if a, ok := p.(*NamedPatternNamed); ok {
		return NamedPatternNamedToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*NamedPatternAnonymous); ok {
		return NamedPatternAnonymousToPreservesSchema(*a, indent)
	}

	return ""
}

func NamedPatternNamedFromPreservesSchema(values []Value) *NamedPatternNamed {
	if a := BindingFromPreservesSchema(values); a != nil {
		return &NamedPatternNamed{*a}
	}

	return nil
}
func NamedPatternNamedToPreservesSchema(p NamedPatternNamed, indent string) string {

	return BindingToPreservesSchema(p.Binding, indent)

}

func NamedPatternAnonymousFromPreservesSchema(values []Value) *NamedPatternAnonymous {
	if a := PatternFromPreservesSchema(values); a != nil {
		return &NamedPatternAnonymous{a}
	}

	return nil
}
func NamedPatternAnonymousToPreservesSchema(p NamedPatternAnonymous, indent string) string {
	return PatternToPreservesSchema(p.Pattern, indent)
}

func EmbeddedTypeNameFromPreservesSchema(values []Value) EmbeddedTypeName {
	if a := EmbeddedTypeNameFalseFromPreservesSchema(values); a != nil {
		return a
	}
	if a := EmbeddedTypeNameRefFromPreservesSchema(values); a != nil {
		return a
	}

	return nil
}

func EmbeddedTypeNameToPreservesSchema(p EmbeddedTypeName, indent string) string {
	if a, ok := p.(*EmbeddedTypeNameFalse); ok {
		return EmbeddedTypeNameFalseToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*EmbeddedTypeNameRef); ok {
		return EmbeddedTypeNameRefToPreservesSchema(*a, indent)
	}

	return ""
}

func EmbeddedTypeNameFalseFromPreservesSchema(values []Value) *EmbeddedTypeNameFalse {
	if len(values) != 1 {
		return nil
	}
	e, ok := values[0].(*Boolean)
	if !ok {
		return nil
	}
	if !(*e) {
		return &EmbeddedTypeNameFalse{}
	}
	return nil
}
func EmbeddedTypeNameFalseToPreservesSchema(p EmbeddedTypeNameFalse, indent string) string {
	return "#f"
}

func EmbeddedTypeNameRefFromPreservesSchema(values []Value) *EmbeddedTypeNameRef {
	if a := RefFromPreservesSchema(values); a != nil {
		return &EmbeddedTypeNameRef{*a}
	}

	return nil
}
func EmbeddedTypeNameRefToPreservesSchema(p EmbeddedTypeNameRef, indent string) string {
	return RefToPreservesSchema(p.Ref, indent)

}

func DefinitionFromPreservesSchema(values []Value) Definition {
	if a := DefinitionOrFromPreservesSchema(values); a != nil {
		return a
	}
	if a := DefinitionAndFromPreservesSchema(values); a != nil {
		return a
	}
	if a := DefinitionPatternFromPreservesSchema(values); a != nil {
		return a
	}
	return nil
}

func DefinitionToPreservesSchema(p Definition, indent string) string {
	if a, ok := p.(*DefinitionOr); ok {
		return DefinitionOrToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*DefinitionAnd); ok {
		return DefinitionAndToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*DefinitionPattern); ok {
		return DefinitionPatternToPreservesSchema(*a, indent)
	}

	return ""
}

func DefinitionOrFromPreservesSchema(values []Value) *DefinitionOr {
	if len(values) < 2 {
		return nil
	}
	sawSlash := true
	var patterns []NamedAlternative
	for _, value := range values {
	retryOr:
		switch v := value.(type) {
		case *Comment:
			value = v.AnnotatedValue
			goto retryOr
		case *Symbol:
			if v.String() == "/" {
				sawSlash = true
				continue
			} else {
				if !sawSlash {
					return nil
				}
				sawSlash = false
				p := NamedAlternativeFromPreservesSchema([]Value{value})
				if p == nil {
					return nil
				}
				patterns = append(patterns, *p)
			}
		default:
			if !sawSlash {
				return nil
			}
			sawSlash = false
			p := NamedAlternativeFromPreservesSchema([]Value{value})
			if p == nil {
				return nil
			}
			patterns = append(patterns, *p)
		}
	}
	if len(patterns) < 2 {
		return nil
	}
	return &DefinitionOr{Pattern0: patterns[0], Pattern1: patterns[1], PatternN: patterns[2:]}
}
func DefinitionOrToPreservesSchema(p DefinitionOr, indent string) string {
	var s string
	for _, v := range p.PatternN {
		s = fmt.Sprintf("%s\n%s  /\n%s  %s", s, indent, indent, NamedAlternativeToPreservesSchema(v, fmt.Sprintf("%s  ", indent)))
	}
	return fmt.Sprintf("\n%s  /\n%s  %s\n%s  /\n%s  %s%s", indent, indent, NamedAlternativeToPreservesSchema(p.Pattern0, fmt.Sprintf("%s  ", indent)), indent, indent, NamedAlternativeToPreservesSchema(p.Pattern1, fmt.Sprintf("%s  ", indent)), s)
}
func DefinitionAndFromPreservesSchema(values []Value) *DefinitionAnd {
	if len(values) < 2 {
		return nil
	}
	sawAnd := true
	var patterns []NamedPattern
	for _, value := range values {
	retryOr:
		switch v := value.(type) {
		case *Comment:
			value = v.AnnotatedValue
			goto retryOr
		case *Symbol:
			if v.String() == "&" {
				sawAnd = true
				continue
			} else {
				if !sawAnd {
					return nil
				}
				sawAnd = false
				p := NamedPatternFromPreservesSchema([]Value{value})
				if p == nil {
					return nil
				}
				patterns = append(patterns, p)
			}
		default:
			if !sawAnd {
				return nil
			}
			sawAnd = false
			p := NamedPatternFromPreservesSchema([]Value{value})
			if p == nil {
				return nil
			}
			patterns = append(patterns, p)
		}
	}
	if len(patterns) < 2 {
		return nil
	}
	return &DefinitionAnd{Pattern0: patterns[0], Pattern1: patterns[1], PatternN: patterns[2:]}

}
func DefinitionAndToPreservesSchema(p DefinitionAnd, indent string) string {
	var s string
	for _, v := range p.PatternN {
		s = fmt.Sprintf("%s\n%s  &\n%s  %s", s, indent, indent, NamedPatternToPreservesSchema(v, fmt.Sprintf("%s  ", indent)))
	}
	return fmt.Sprintf("\n%s  &\n%s%s\n%s  &\n%s%s%s", indent, indent, NamedPatternToPreservesSchema(p.Pattern0, fmt.Sprintf("%s  ", indent)), indent, indent, NamedPatternToPreservesSchema(p.Pattern0, fmt.Sprintf("%s  ", indent)), s)
}

func DefinitionPatternFromPreservesSchema(values []Value) *DefinitionPattern {
	if a := PatternFromPreservesSchema(values); a != nil {
		return &DefinitionPattern{a}
	}

	return nil
}
func DefinitionPatternToPreservesSchema(p DefinitionPattern, indent string) string {
	return PatternToPreservesSchema(p.Pattern, indent)
}

func CompoundPatternFromPreservesSchema(values []Value) CompoundPattern {
	if a := CompoundPatternRecFromPreservesSchema(values); a != nil {
		return a
	}
	if a := CompoundPatternTupleFromPreservesSchema(values); a != nil {
		return a
	}
	if a := CompoundPatternTuplePrefixFromPreservesSchema(values); a != nil {
		return a
	}
	if a := CompoundPatternDictFromPreservesSchema(values); a != nil {
		return a
	}

	return nil
}
func CompoundPatternToPreservesSchema(p CompoundPattern, indent string) string {
	if a, ok := p.(*CompoundPatternRec); ok {
		return CompoundPatternRecToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*CompoundPatternTuple); ok {
		return CompoundPatternTupleToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*CompoundPatternTuplePrefix); ok {
		return CompoundPatternTuplePrefixToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*CompoundPatternDict); ok {
		return CompoundPatternDictToPreservesSchema(*a, indent)
	}

	return ""
}

// From an actual <rec <lit any> <tuple []>>
// r *Rec is a new object
func CompoundPatternRecFromPreservesSchema(values []Value) *CompoundPatternRec {
	if len(values) != 1 {
		return nil
	}
	r, ok := values[0].(*Record)
	if !ok {
		return nil
	}
	s := r.Key.(*Symbol)
	if s == nil {
		return nil
	}
	label := NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol(s.String()))))

	if len(r.Fields) < 1 {
		return nil
	}
	p := Sequence(r.Fields)
	t := CompoundPatternTupleFromPreservesSchema([]Value{&p})
	if t == nil {
		return nil
	}
	// TODO solve skip brackets
	// t.(*Tuple).SkipBrackets = true
	fields := NewNamedPatternAnonymous(NewPatternCompoundPattern(t))

	return &CompoundPatternRec{Label: label, Fields: fields}
}
func CompoundPatternRecToPreservesSchema(p CompoundPatternRec, indent string) string {
	return fmt.Sprintf("<%s %s>", NamedPatternToPreservesSchema(p.Label, fmt.Sprintf("%s  ", indent)), NamedPatternToPreservesSchema(p.Fields, fmt.Sprintf("%s  ", indent)))
}

func CompoundPatternTupleFromPreservesSchema(values []Value) *CompoundPatternTuple {
	if len(values) != 1 {
		return nil
	}
	u, ok := values[0].(*Sequence)
	if !ok {
		return nil
	}
	if len(*u) == 0 {
		return nil
	}
	if symbol, ok := (*u)[len(*u)-1].(*Symbol); ok {
		if symbol.String() == "..." {
			return nil
		}
	}

	var patterns []NamedPattern
	for _, v := range *u {
		t := NamedPatternFromPreservesSchema([]Value{v})
		if t == nil {
			return nil
		}
		namedPattern, ok := t.(NamedPattern)
		if !ok {
			return nil
		}
		patterns = append(patterns, namedPattern)
	}
	if len(patterns) == 0 {
		return nil
	}
	return &CompoundPatternTuple{Patterns: patterns}
}
func CompoundPatternTupleToPreservesSchema(p CompoundPatternTuple, indent string) string {
	var s string
	for _, a := range p.Patterns {
		if len(s) > 0 {
			s = fmt.Sprintf("%s ", s)
		}

		s = fmt.Sprintf("%s%s", s, NamedPatternToPreservesSchema(a, fmt.Sprintf("%s  ", indent)))

	}
	// TODO: solve skip brackets
	// if len(s) > 0 {
	// 	if t.SkipBrackets {
	// 		return s
	// 	}
	//	return fmt.Sprintf("[%s]", s)
	// }
	if len(indent) == 0 {
		return fmt.Sprintf("[%s]", s)
	} else {
		return s
	}
}

func CompoundPatternTuplePrefixFromPreservesSchema(values []Value) *CompoundPatternTuplePrefix {
	if len(values) != 1 {
		return nil
	}
	sequence, ok := values[0].(*Sequence)
	if !ok {
		return nil
	}
	if len(*sequence) < 3 {
		return nil
	}
	i := 2
	if dotdotdot, ok := (*sequence)[len(*sequence)-1].(*Symbol); ok {
		if dotdotdot.String() != "..." {
			return nil
		}
		if anno, ok := (*sequence)[len(*sequence)-2].(*Annotation); ok {
			newValue := Sequence([]Value{anno.AnnotatedValue, dotdotdot})
			anno.AnnotatedValue = &newValue
			*sequence = (*sequence)[:len(*sequence)-1]
			i = 1
		}
	}

	var fixed []NamedPattern
	for _, v := range (*sequence)[:len(*sequence)-i] {
		u := NamedPatternFromPreservesSchema([]Value{v})
		if u == nil {
			return nil
		}
		f, ok := u.(NamedPattern)
		if !ok {
			return nil
		}
		fixed = append(fixed, f)
	}
	seq := (*sequence)[len(fixed):]
	var u NamedSimplePattern
	if len(seq) == 1 {
		u = NamedSimplePatternFromPreservesSchema([]Value{seq[0]})
	} else {
		u = NamedSimplePatternFromPreservesSchema([]Value{&seq})
	}
	if u == nil {
		return nil
	}
	v, ok := u.(NamedSimplePattern)
	if !ok {
		return nil
	}
	return &CompoundPatternTuplePrefix{Fixed: fixed, Variable: v}

}
func CompoundPatternTuplePrefixToPreservesSchema(p CompoundPatternTuplePrefix, indent string) string {
	var s string
	for _, f := range p.Fixed {
		s = fmt.Sprintf("%s\n%s  %s", s, indent, NamedPatternToPreservesSchema(f, indent))
	}
	t := NamedSimplePatternToPreservesSchema(p.Variable, indent)
	if strings.HasSuffix(t, "]") {
		t = strings.TrimSuffix(t, "]")
		t = strings.Replace(t, "[", "", 1)
	}
	return fmt.Sprintf("[%s %s]", s, t)
}

func CompoundPatternDictFromPreservesSchema(values []Value) *CompoundPatternDict {
	d := DictionaryEntriesFromPreservesSchema(make(DictionaryEntries), values)
	if d == nil {
		return nil
	}
	if len(*d) == 0 {
		return nil
	}
	return &CompoundPatternDict{Entries: *d}
}
func CompoundPatternDictToPreservesSchema(p CompoundPatternDict, indent string) string {
	return fmt.Sprintf("%s", DictionaryEntriesToPreservesSchema(p.Entries, indent))
}

func BooleanFromPreservesSchema(values []Value) *Boolean {
	if len(values) != 1 {
		return nil
	}
	u, ok := values[0].(*Boolean)
	if !ok {
		return nil
	}
	b := Boolean(bool(*u))
	return &b
}
func BooleanToPreservesSchema(p Boolean, _ string) string {
	if p {
		return "#t"
	}
	return "#f"
}
func DoubleFromPreservesSchema(values []Value) *Double {
	if len(values) != 1 {
		return nil
	}
	u, ok := values[0].(*Double)
	if !ok {
		return nil
	}
	d := Double(float64(*u))
	return &d
}
func DoubleToPreservesSchema(p Double, _ string) string { return fmt.Sprintf("%f", p) }

func SignedIntegerFromPreservesSchema(values []Value) *SignedInteger {
	if len(values) != 1 {
		return nil
	}
	u, ok := values[0].(*SignedInteger)
	if !ok {
		return nil
	}
	s := SignedInteger(big.Int(*u))
	return &s
}
func SignedIntegerToPreservesSchema(p SignedInteger, _ string) string {
	return big.NewInt(0).Set(&([]big.Int{big.Int(p)}[0])).String()
}

func PstringFromPreservesSchema(values []Value) *Pstring {
	if len(values) != 1 {
		return nil
	}
	u, ok := values[0].(*Pstring)
	if !ok {
		return nil
	}
	return u
}
func PstringToPreservesSchema(s Pstring, _ string) string { return strconv.Quote(string(s)) }

func ByteStringFromPreservesSchema(values []Value) *ByteString {
	if len(values) != 1 {
		return nil
	}
	u, ok := values[0].(*ByteString)
	if !ok {
		return nil
	}
	return u
}
func ByteStringToPreservesSchema(b ByteString, _ string) string { return string(b) }

func SymbolFromPreservesSchema(values []Value) *Symbol {
	if len(values) != 1 {
		return nil
	}
	v, ok := values[0].(*Symbol)
	if !ok {
		return nil
	}
	return v
}
func SymbolToPreservesSchema(p Symbol, _ string) string { return p.String() }
func PatternFromPreservesSchema(values []Value) Pattern {
	if a := PatternSimplePatternFromPreservesSchema(values); a != nil {
		return a
	}
	if a := PatternCompoundPatternFromPreservesSchema(values); a != nil {
		return a
	}

	return nil
}

func PatternToPreservesSchema(p Pattern, indent string) string {
	if a, ok := p.(*PatternSimplePattern); ok {
		return PatternSimplePatternToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*PatternCompoundPattern); ok {
		return PatternCompoundPatternToPreservesSchema(*a, indent)
	}

	return ""
}
func PatternSimplePatternFromPreservesSchema(values []Value) *PatternSimplePattern {
	if p := SimplePatternFromPreservesSchema(values); p != nil {
		return &PatternSimplePattern{p}
	}

	return nil
}
func PatternSimplePatternToPreservesSchema(p PatternSimplePattern, indent string) string {
	return SimplePatternToPreservesSchema(p.SimplePattern, indent)
}
func PatternCompoundPatternFromPreservesSchema(values []Value) *PatternCompoundPattern {
	if p := CompoundPatternFromPreservesSchema(values); p != nil {
		return &PatternCompoundPattern{p}
	}

	return nil
}
func PatternCompoundPatternToPreservesSchema(p PatternCompoundPattern, indent string) string {
	return CompoundPatternToPreservesSchema(p.CompoundPattern, indent)
}

func ValueFromPreservesSchema(values []Value) Value {
	if a := BooleanFromPreservesSchema(values); a != nil {
		return a
	}
	if a := DoubleFromPreservesSchema(values); a != nil {
		return a
	}
	if a := SignedIntegerFromPreservesSchema(values); a != nil {
		return a
	}
	if a := PstringFromPreservesSchema(values); a != nil {
		return a
	}
	if a := ByteStringFromPreservesSchema(values); a != nil {
		return a
	}
	if a := SymbolFromPreservesSchema(values); a != nil {
		return a
	}

	return nil
}
func ValueToPreservesSchema(p Value, indent string) string {
	if a, ok := p.(*Boolean); ok {
		return BooleanToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*Double); ok {
		return DoubleToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*SignedInteger); ok {
		return SignedIntegerToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*Pstring); ok {
		return PstringToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*ByteString); ok {
		return ByteStringToPreservesSchema(*a, indent)
	}
	if a, ok := p.(*Symbol); ok {
		return SymbolToPreservesSchema(*a, indent)
	}

	return ""
}
