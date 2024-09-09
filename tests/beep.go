package beep

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	. "github.com/isodude/preserves-go/lib/preserves"
	"github.com/isodude/preserves-go/lib/preserves/text"
)

type EmbeddedTypeName interface {
	IsEmbeddedTypeName()
}

func EmbeddedTypeNameFromPreserves(value Value) EmbeddedTypeName {
	for _, v := range []EmbeddedTypeName{&EmbeddedTypeNameFalse{}, &EmbeddedTypeNameRef{}} {
		switch u := v.(type) {
		case *EmbeddedTypeNameFalse:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *EmbeddedTypeNameRef:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		}
	}
	return nil
}

type EmbeddedTypeNameFalse struct {
}

func (*EmbeddedTypeNameFalse) IsEmbeddedTypeName() {
}
func (*EmbeddedTypeNameFalse) FromPreserves(value Value) *EmbeddedTypeNameFalse {
	if v := BooleanFromPreserves(value); v != nil {
		if *v == false {
			return &EmbeddedTypeNameFalse{}
		}
	}
	return nil
}

type EmbeddedTypeNameRef struct {
	Ref
}

func NewEmbeddedTypeNameRef(obj Ref) *EmbeddedTypeNameRef {
	return &EmbeddedTypeNameRef{Ref: obj}
}
func (*EmbeddedTypeNameRef) IsEmbeddedTypeName() {
}
func (p *EmbeddedTypeNameRef) FromPreserves(value Value) *EmbeddedTypeNameRef {
	if o := (&Ref{}).FromPreserves(value); o != nil {
		return &EmbeddedTypeNameRef{Ref: *o}
	}
	return nil
}

type ModulePath struct {
	s string
}

func (*ModulePath) FromPreserves(value Value) *ModulePath {
	if seq, ok := value.(*Sequence); ok {
		var s []string
		for _, k := range *seq {
			if v := SymbolFromPreserves(k); v != nil {
				s = append(s, v.String())
			}
		}
		return &ModulePath{s: strings.Join(s, " ")}
	}
	return nil
}

type Binding struct {
	Name    Symbol
	Pattern SimplePattern
}

func NewBinding(name Symbol, pattern SimplePattern) *Binding {
	return &Binding{Name: name, Pattern: pattern}
}
func (*Binding) FromPreserves(value Value) *Binding {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "named" {
			if p0 := SymbolFromPreserves(rec.Fields[0]); p0 != nil {
				if p1 := SimplePatternFromPreserves(rec.Fields[1]); p1 != nil {
					return &Binding{Name: *p0, Pattern: p1}
				} else {
					fmt.Printf("(Binding) Failed on %v: %v\n", p0, rec.Fields[1])

				}
			} else {
				fmt.Printf("(Binding) Failed on %v\n", rec.Fields[0])

			}
		}
	}
	return nil
}

type Schema struct {
	Definitions  Definitions
	EmbeddedType EmbeddedTypeName
	Version      Version
}

func NewSchema(definitions Definitions, embeddedType EmbeddedTypeName, version Version) *Schema {
	return &Schema{Definitions: definitions, EmbeddedType: embeddedType, Version: version}
}
func (*Schema) FromPreserves(value Value) *Schema {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "schema" {
			if dict, ok := rec.Fields[0].(*Dictionary); ok {
				var obj Schema
				for dictKey, dictValue := range *dict {
					if dictKey.Equal(NewSymbol("definitions")) {
						if p0 := NewDefinitions().FromPreserves(dictValue); p0 != nil {
							obj.Definitions = *p0
							continue
						}
					}
					if dictKey.Equal(NewSymbol("embeddedType")) {
						if p1 := EmbeddedTypeNameFromPreserves(dictValue); p1 != nil {
							obj.EmbeddedType = p1
							continue
						}
					}
					if dictKey.Equal(NewSymbol("version")) {
						if p2 := (&Version{}).FromPreserves(dictValue); p2 != nil {
							obj.Version = *p2
							continue
						}
					}
					return nil
				}
				return &obj
			}
		}
	}
	return nil
}

type AtomKind interface {
	IsAtomKind()
}

func AtomKindFromPreserves(value Value) AtomKind {
	for _, v := range []AtomKind{&AtomKindBoolean{}, &AtomKindDouble{}, &AtomKindSignedInteger{}, &AtomKindString{}, &AtomKindByteString{}, &AtomKindSymbol{}} {
		switch u := v.(type) {
		case *AtomKindBoolean:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *AtomKindDouble:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *AtomKindSignedInteger:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *AtomKindString:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *AtomKindByteString:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *AtomKindSymbol:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		}
	}
	return nil
}

type AtomKindBoolean struct {
}

func (*AtomKindBoolean) IsAtomKind() {
}
func (*AtomKindBoolean) FromPreserves(value Value) *AtomKindBoolean {
	if v := SymbolFromPreserves(value); v != nil {
		if *v == "Boolean" {
			return &AtomKindBoolean{}
		}
	}
	return nil
}

type AtomKindDouble struct {
}

func (*AtomKindDouble) IsAtomKind() {
}
func (*AtomKindDouble) FromPreserves(value Value) *AtomKindDouble {
	if v := SymbolFromPreserves(value); v != nil {
		if *v == "Double" {
			return &AtomKindDouble{}
		}
	}
	return nil
}

type AtomKindSignedInteger struct {
}

func (*AtomKindSignedInteger) IsAtomKind() {
}
func (*AtomKindSignedInteger) FromPreserves(value Value) *AtomKindSignedInteger {
	if v := SymbolFromPreserves(value); v != nil {
		if *v == "SignedInteger" {
			return &AtomKindSignedInteger{}
		}
	}
	return nil
}

type AtomKindString struct {
}

func (*AtomKindString) IsAtomKind() {
}
func (*AtomKindString) FromPreserves(value Value) *AtomKindString {
	if v := SymbolFromPreserves(value); v != nil {
		if *v == "String" {
			return &AtomKindString{}
		}
	}
	return nil
}

type AtomKindByteString struct {
}

func (*AtomKindByteString) IsAtomKind() {
}
func (*AtomKindByteString) FromPreserves(value Value) *AtomKindByteString {
	if v := SymbolFromPreserves(value); v != nil {
		if *v == "ByteString" {
			return &AtomKindByteString{}
		}
	}
	return nil
}

type AtomKindSymbol struct {
}

func (*AtomKindSymbol) IsAtomKind() {
}
func (*AtomKindSymbol) FromPreserves(value Value) *AtomKindSymbol {
	if v := SymbolFromPreserves(value); v != nil {
		if *v == "Symbol" {
			return &AtomKindSymbol{}
		}
	}
	return nil
}

type Pattern interface {
	IsPattern()
}

func PatternFromPreserves(value Value) Pattern {
	for _, v := range []Pattern{&PatternSimplePattern{}, &PatternCompoundPattern{}} {
		switch u := v.(type) {
		case *PatternSimplePattern:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *PatternCompoundPattern:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		}
	}
	return nil
}

type PatternSimplePattern struct {
	SimplePattern
}

func NewPatternSimplePattern(obj SimplePattern) *PatternSimplePattern {
	return &PatternSimplePattern{SimplePattern: obj}
}
func (*PatternSimplePattern) IsPattern() {
}
func (p *PatternSimplePattern) FromPreserves(value Value) *PatternSimplePattern {
	if o := SimplePatternFromPreserves(value); o != nil {
		return &PatternSimplePattern{SimplePattern: o}
	}
	return nil
}

type PatternCompoundPattern struct {
	CompoundPattern
}

func NewPatternCompoundPattern(obj CompoundPattern) *PatternCompoundPattern {
	return &PatternCompoundPattern{CompoundPattern: obj}
}
func (*PatternCompoundPattern) IsPattern() {
}
func (p *PatternCompoundPattern) FromPreserves(value Value) *PatternCompoundPattern {
	if o := CompoundPatternFromPreserves(value); o != nil {
		return &PatternCompoundPattern{CompoundPattern: o}
	}
	return nil
}

type Ref struct {
	Module ModulePath
	Name   Symbol
}

func NewRef(module ModulePath, name Symbol) *Ref {
	return &Ref{Module: module, Name: name}
}
func (*Ref) FromPreserves(value Value) *Ref {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "ref" {
			if p0 := (&ModulePath{}).FromPreserves(rec.Fields[0]); p0 != nil {
				if p1 := SymbolFromPreserves(rec.Fields[1]); p1 != nil {
					return &Ref{Module: *p0, Name: *p1}
				} else {
					fmt.Printf("(Ref) Failed on %v: %v\n", p0, rec.Fields[1])
				}
			} else {
				fmt.Printf("(tuple) Failed on %v\n", rec.Fields[0])

			}
		}
	}
	return nil
}

type NamedAlternative struct {
	VariantLabel Pstring
	Pattern      Pattern
}

func NewNamedAlternative(variantLabel Pstring, pattern Pattern) *NamedAlternative {
	return &NamedAlternative{VariantLabel: variantLabel, Pattern: pattern}
}
func (*NamedAlternative) FromPreserves(value Value) *NamedAlternative {
	if seq, ok := value.(*Sequence); ok && len(*seq) == 2 {
		if p0 := PstringFromPreserves((*seq)[0]); p0 != nil {
			if p1 := PatternFromPreserves((*seq)[1]); p1 != nil {
				return &NamedAlternative{VariantLabel: *p0, Pattern: p1}
			} else {
				fmt.Printf("(NamedAlternative) Failed on %v: %v\n", p0, (*seq)[1])
			}
		} else {

			fmt.Printf("(NamedAlternative) Failed on %v\n", (*seq)[0])
		}
	}

	return nil
}

type DictionaryEntries map[Value]NamedSimplePattern

func NewDictionaryEntries() DictionaryEntries {
	return make(DictionaryEntries)
}

func (DictionaryEntries) FromPreserves(value Value) *DictionaryEntries {
	if dict, ok := value.(*Dictionary); ok {
		obj := NewDictionaryEntries()
		for dictKey, dictValue := range *dict {
			if dKey := ValueFromPreserves(dictKey); dKey != nil {
				if dValue := NamedSimplePatternFromPreserves(dictValue); dValue != nil {
					obj[dKey] = dValue
					continue
				}
			}
			return nil
		}
		return &obj
	}
	return nil
}

type Modules map[ModulePath]Schema

func NewModules() Modules {
	return make(Modules)
}
func (Modules) FromPreserves(value Value) *Modules {
	if dict, ok := value.(*Dictionary); ok {
		obj := NewModules()
		for dictKey, dictValue := range *dict {
			if dKey := (&ModulePath{}).FromPreserves(dictKey); dKey != nil {
				if dValue := (&Schema{}).FromPreserves(dictValue); dValue != nil {
					obj[*dKey] = *dValue
					continue
				} else {
					fmt.Printf("(Modules) Failed on %v: %v\n", dKey, dictValue)
				}
			} else {
				fmt.Printf("(Modules) Failed on %v\n", dictKey)
			}
			return nil
		}
		return &obj
	}
	return nil
}

type NamedPattern interface {
	IsNamedPattern()
}

func NamedPatternFromPreserves(value Value) NamedPattern {
	for _, v := range []NamedPattern{&NamedPatternNamed{}, &NamedPatternAnonymous{}} {
		switch u := v.(type) {
		case *NamedPatternNamed:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *NamedPatternAnonymous:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		}
	}
	return nil
}

type NamedPatternNamed struct {
	Binding
}

func NewNamedPatternNamed(obj Binding) *NamedPatternNamed {
	return &NamedPatternNamed{Binding: obj}
}
func (*NamedPatternNamed) IsNamedPattern() {
}
func (p *NamedPatternNamed) FromPreserves(value Value) *NamedPatternNamed {
	if o := (&Binding{}).FromPreserves(value); o != nil {
		return &NamedPatternNamed{Binding: *o}
	}
	return nil
}

type NamedPatternAnonymous struct {
	Pattern
}

func NewNamedPatternAnonymous(obj Pattern) *NamedPatternAnonymous {
	return &NamedPatternAnonymous{Pattern: obj}
}
func (*NamedPatternAnonymous) IsNamedPattern() {
}
func (p *NamedPatternAnonymous) FromPreserves(value Value) *NamedPatternAnonymous {
	if o := PatternFromPreserves(value); o != nil {
		return &NamedPatternAnonymous{Pattern: o}
	}
	return nil
}

type Definition interface {
	IsDefinition()
}

func DefinitionFromPreserves(value Value) Definition {
	for _, v := range []Definition{&DefinitionOr{}, &DefinitionAnd{}, &DefinitionPattern{}} {
		switch u := v.(type) {
		case *DefinitionOr:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *DefinitionAnd:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *DefinitionPattern:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		}
	}
	return nil
}

type DefinitionOr struct {
	Pattern0 NamedAlternative
	Pattern1 NamedAlternative
	PatternN []NamedAlternative
}

func NewDefinitionOr(pattern0 NamedAlternative, pattern1 NamedAlternative, patternN []NamedAlternative) *DefinitionOr {
	return &DefinitionOr{Pattern0: pattern0, Pattern1: pattern1, PatternN: patternN}
}
func (*DefinitionOr) IsDefinition() {
}
func (*DefinitionOr) FromPreserves(value Value) *DefinitionOr {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "or" {
			if seq, ok := rec.Fields[0].(*Sequence); ok {
				var patterns []NamedAlternative
				for _, item := range *seq {
					if itemParsed := (&NamedAlternative{}).FromPreserves(item); itemParsed != nil {
						patterns = append(patterns, *itemParsed)
					} else {
						var b bytes.Buffer
						text.FromPreserves(item).WriteTo(&b)
						fmt.Printf("(DefinitionOr) Failed on %v\n", b.String())
						return nil
					}
				}
				return &DefinitionOr{Pattern0: patterns[0], Pattern1: patterns[1], PatternN: patterns[2:]}
			}
		}
	}
	return nil
}

type DefinitionAnd struct {
	Pattern0 NamedPattern
	Pattern1 NamedPattern
	PatternN []NamedPattern
}

func NewDefinitionAnd(pattern0 NamedPattern, pattern1 NamedPattern, patternN []NamedPattern) *DefinitionAnd {
	return &DefinitionAnd{Pattern0: pattern0, Pattern1: pattern1, PatternN: patternN}
}
func (*DefinitionAnd) IsDefinition() {
}
func (*DefinitionAnd) FromPreserves(value Value) *DefinitionAnd {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "and" {
			if seq, ok := rec.Fields[0].(*Sequence); ok {
				var patterns []NamedPattern
				for _, item := range *seq {
					if itemParsed := NamedPatternFromPreserves(item); itemParsed != nil {
						patterns = append(patterns, itemParsed)
					} else {
						return nil
					}
				}
				return &DefinitionAnd{Pattern0: patterns[0], Pattern1: patterns[1], PatternN: patterns[2:]}
			}
		}
	}
	return nil
}

type DefinitionPattern struct {
	Pattern
}

func NewDefinitionPattern(obj Pattern) *DefinitionPattern {
	return &DefinitionPattern{Pattern: obj}
}
func (*DefinitionPattern) IsDefinition() {
}
func (p *DefinitionPattern) FromPreserves(value Value) *DefinitionPattern {
	if o := PatternFromPreserves(value); o != nil {
		return &DefinitionPattern{Pattern: o}
	}
	return nil
}

type Bundle struct {
	Modules Modules
}

func NewBundle(modules Modules) *Bundle {
	return &Bundle{Modules: modules}
}
func (*Bundle) FromPreserves(value Value) *Bundle {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "bundle" {
			if p0 := NewModules().FromPreserves(rec.Fields[0]); p0 != nil {
				return &Bundle{Modules: *p0}
			}
		}
	}
	return nil
}

type Version struct {
}

func (*Version) FromPreserves(value Value) *Version {
	if v := SignedIntegerFromPreserves(value); v != nil {
		a := big.Int(*v)
		if big.NewInt(1).Cmp(&a) == 0 {
			return &Version{}
		}
	}
	return nil
}

type Definitions map[Symbol]Definition

func NewDefinitions() Definitions {
	return make(Definitions)
}
func (Definitions) FromPreserves(value Value) *Definitions {
	if dict, ok := value.(*Dictionary); ok {
		obj := NewDefinitions()
		for dictKey, dictValue := range *dict {
			if dKey := SymbolFromPreserves(dictKey); dKey != nil {
				if dValue := DefinitionFromPreserves(dictValue); dValue != nil {
					obj[*dKey] = dValue
					continue
				} else {
					var b bytes.Buffer
					text.FromPreserves(dictValue).WriteTo(&b)
					fmt.Printf("failed on %v:%v\n", dKey, b.String())
				}
			} else {
				fmt.Printf("failed on %v\n", dictKey)
			}
			return nil
		}
		return &obj
	}
	return nil
}

type SimplePattern interface {
	IsSimplePattern()
}

func SimplePatternFromPreserves(value Value) SimplePattern {
	for _, v := range []SimplePattern{&SimplePatternAny{}, &SimplePatternAtom{}, &SimplePatternEmbedded{}, &SimplePatternLit{}, &SimplePatternSeqof{}, &SimplePatternSetof{}, &SimplePatternDictof{}, &SimplePatternRef{}} {
		switch u := v.(type) {
		case *SimplePatternAny:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *SimplePatternAtom:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *SimplePatternEmbedded:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *SimplePatternLit:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *SimplePatternSeqof:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *SimplePatternSetof:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *SimplePatternDictof:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *SimplePatternRef:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		}
	}
	return nil
}

type SimplePatternAny struct {
}

func (*SimplePatternAny) IsSimplePattern() {
}
func (*SimplePatternAny) FromPreserves(value Value) *SimplePatternAny {
	if v := SymbolFromPreserves(value); v != nil {
		if *v == "any" {
			return &SimplePatternAny{}
		}
	}
	return nil
}

type SimplePatternAtom struct {
	AtomKind AtomKind
}

func NewSimplePatternAtom(atomKind AtomKind) *SimplePatternAtom {
	return &SimplePatternAtom{AtomKind: atomKind}
}
func (*SimplePatternAtom) IsSimplePattern() {
}
func (*SimplePatternAtom) FromPreserves(value Value) *SimplePatternAtom {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "atom" {
			if p0 := AtomKindFromPreserves(rec.Fields[0]); p0 != nil {
				return &SimplePatternAtom{AtomKind: p0}
			}
		}
	}
	return nil
}

type SimplePatternEmbedded struct {
	Interface SimplePattern
}

func NewSimplePatternEmbedded(_interface SimplePattern) *SimplePatternEmbedded {
	return &SimplePatternEmbedded{Interface: _interface}
}
func (*SimplePatternEmbedded) IsSimplePattern() {
}
func (*SimplePatternEmbedded) FromPreserves(value Value) *SimplePatternEmbedded {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "embedded" {
			if p0 := SimplePatternFromPreserves(rec.Fields[0]); p0 != nil {
				return &SimplePatternEmbedded{Interface: p0}
			}
		}
	}
	return nil
}

type SimplePatternLit struct {
	Value Value
}

func NewSimplePatternLit(value Value) *SimplePatternLit {
	return &SimplePatternLit{Value: value}
}
func (*SimplePatternLit) IsSimplePattern() {
}
func (*SimplePatternLit) FromPreserves(value Value) *SimplePatternLit {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "lit" {
			return &SimplePatternLit{Value: value}
		}
	}
	return nil
}

type SimplePatternSeqof struct {
	Pattern SimplePattern
}

func NewSimplePatternSeqof(pattern SimplePattern) *SimplePatternSeqof {
	return &SimplePatternSeqof{Pattern: pattern}
}
func (*SimplePatternSeqof) IsSimplePattern() {
}
func (*SimplePatternSeqof) FromPreserves(value Value) *SimplePatternSeqof {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "seqof" {
			if p0 := SimplePatternFromPreserves(rec.Fields[0]); p0 != nil {
				return &SimplePatternSeqof{Pattern: p0}
			}
		}
	}
	return nil
}

type SimplePatternSetof struct {
	Pattern SimplePattern
}

func NewSimplePatternSetof(pattern SimplePattern) *SimplePatternSetof {
	return &SimplePatternSetof{Pattern: pattern}
}
func (*SimplePatternSetof) IsSimplePattern() {
}
func (*SimplePatternSetof) FromPreserves(value Value) *SimplePatternSetof {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "setof" {
			if p0 := SimplePatternFromPreserves(rec.Fields[0]); p0 != nil {
				return &SimplePatternSetof{Pattern: p0}
			}
		}
	}
	return nil
}

type SimplePatternDictof struct {
	Key   SimplePattern
	Value SimplePattern
}

func NewSimplePatternDictof(key SimplePattern, value SimplePattern) *SimplePatternDictof {
	return &SimplePatternDictof{Key: key, Value: value}
}
func (*SimplePatternDictof) IsSimplePattern() {
}
func (*SimplePatternDictof) FromPreserves(value Value) *SimplePatternDictof {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "dictof" {
			if p0 := SimplePatternFromPreserves(rec.Fields[0]); p0 != nil {
				if p1 := SimplePatternFromPreserves(rec.Fields[1]); p1 != nil {
					return &SimplePatternDictof{Key: p0, Value: p1}
				}
			}
		}
	}
	return nil
}

type SimplePatternRef struct {
	Ref
}

func NewSimplePatternRef(obj Ref) *SimplePatternRef {
	return &SimplePatternRef{Ref: obj}
}
func (*SimplePatternRef) IsSimplePattern() {
}
func (p *SimplePatternRef) FromPreserves(value Value) *SimplePatternRef {
	if o := (&Ref{}).FromPreserves(value); o != nil {
		return &SimplePatternRef{Ref: *o}
	}
	return nil
}

type CompoundPattern interface {
	IsCompoundPattern()
}

func CompoundPatternFromPreserves(value Value) CompoundPattern {
	for _, v := range []CompoundPattern{&CompoundPatternRec{}, &CompoundPatternTuple{}, &CompoundPatternTuplePrefix{}, &CompoundPatternDict{}} {
		switch u := v.(type) {
		case *CompoundPatternRec:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *CompoundPatternTuple:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *CompoundPatternTuplePrefix:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *CompoundPatternDict:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		}
	}
	return nil
}

type CompoundPatternRec struct {
	Label  NamedPattern
	Fields NamedPattern
}

func NewCompoundPatternRec(label NamedPattern, fields NamedPattern) *CompoundPatternRec {
	return &CompoundPatternRec{Label: label, Fields: fields}
}
func (*CompoundPatternRec) IsCompoundPattern() {
}
func (*CompoundPatternRec) FromPreserves(value Value) *CompoundPatternRec {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "rec" {
			if p0 := NamedPatternFromPreserves(rec.Fields[0]); p0 != nil {
				if p1 := NamedPatternFromPreserves(rec.Fields[1]); p1 != nil {
					return &CompoundPatternRec{Label: p0, Fields: p1}
				} else {
					fmt.Printf("(rec) Failed on %v:%v\n", p0, rec.Fields[1])
				}
			} else {
				fmt.Printf("(rec) Failed on %v\n", rec.Fields[0])
			}
		}
	}
	return nil
}

type CompoundPatternTuple struct {
	Patterns []NamedPattern
}

func NewCompoundPatternTuple(patterns []NamedPattern) *CompoundPatternTuple {
	return &CompoundPatternTuple{Patterns: patterns}
}
func (*CompoundPatternTuple) IsCompoundPattern() {
}
func (*CompoundPatternTuple) FromPreserves(value Value) *CompoundPatternTuple {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "tuple" {
			if seq, ok := rec.Fields[0].(*Sequence); ok {
				var p0 []NamedPattern
				for _, item := range *seq {
					if itemParsed := NamedPatternFromPreserves(item); itemParsed != nil {
						p0 = append(p0, itemParsed)
					} else {
						fmt.Printf("(tuple) Failed on %v: %v\n", sym, item)
						return nil
					}
				}
				return &CompoundPatternTuple{Patterns: p0}
			}
		}
	}
	return nil
}

type CompoundPatternTuplePrefix struct {
	Fixed    []NamedPattern
	Variable NamedSimplePattern
}

func NewCompoundPatternTuplePrefix(fixed []NamedPattern, variable NamedSimplePattern) *CompoundPatternTuplePrefix {
	return &CompoundPatternTuplePrefix{Fixed: fixed, Variable: variable}
}
func (*CompoundPatternTuplePrefix) IsCompoundPattern() {
}
func (*CompoundPatternTuplePrefix) FromPreserves(value Value) *CompoundPatternTuplePrefix {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "tuplePrefix" {
			if seq, ok := rec.Fields[0].(*Sequence); ok {
				var p0 []NamedPattern
				for _, item := range *seq {
					if itemParsed := NamedPatternFromPreserves(item); itemParsed != nil {
						p0 = append(p0, itemParsed)
					} else {
						return nil
					}
				}
				if p1 := NamedSimplePatternFromPreserves(rec.Fields[1]); p1 != nil {
					return &CompoundPatternTuplePrefix{Fixed: p0, Variable: p1}
				}
			}
		}
	}
	return nil
}

type CompoundPatternDict struct {
	Entries DictionaryEntries
}

func NewCompoundPatternDict(entries DictionaryEntries) *CompoundPatternDict {
	return &CompoundPatternDict{Entries: entries}
}
func (*CompoundPatternDict) IsCompoundPattern() {
}
func (*CompoundPatternDict) FromPreserves(value Value) *CompoundPatternDict {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "dict" {
			if p0 := NewDictionaryEntries().FromPreserves(rec.Fields[0]); p0 != nil {
				return &CompoundPatternDict{Entries: *p0}
			}
		}
	}
	return nil
}

type NamedSimplePattern interface {
	IsNamedSimplePattern()
}

func NamedSimplePatternFromPreserves(value Value) NamedSimplePattern {
	for _, v := range []NamedSimplePattern{&NamedSimplePatternNamed{}, &NamedSimplePatternAnonymous{}} {
		switch u := v.(type) {
		case *NamedSimplePatternNamed:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		case *NamedSimplePatternAnonymous:
			if o := u.FromPreserves(value); o != nil {
				return o
			}
		}
	}
	return nil
}

type NamedSimplePatternNamed struct {
	Binding
}

func NewNamedSimplePatternNamed(obj Binding) *NamedSimplePatternNamed {
	return &NamedSimplePatternNamed{Binding: obj}
}
func (*NamedSimplePatternNamed) IsNamedSimplePattern() {
}
func (p *NamedSimplePatternNamed) FromPreserves(value Value) *NamedSimplePatternNamed {
	if o := (&Binding{}).FromPreserves(value); o != nil {
		return &NamedSimplePatternNamed{Binding: *o}
	}
	return nil
}

type NamedSimplePatternAnonymous struct {
	SimplePattern
}

func NewNamedSimplePatternAnonymous(obj SimplePattern) *NamedSimplePatternAnonymous {
	return &NamedSimplePatternAnonymous{SimplePattern: obj}
}
func (*NamedSimplePatternAnonymous) IsNamedSimplePattern() {
}
func (p *NamedSimplePatternAnonymous) FromPreserves(value Value) *NamedSimplePatternAnonymous {
	if o := SimplePatternFromPreserves(value); o != nil {
		return &NamedSimplePatternAnonymous{SimplePattern: o}
	}
	return nil
}
