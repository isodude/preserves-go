package schema

import (
	"github.com/isodude/preserves-go/lib/extras"
	. "github.com/isodude/preserves-go/lib/preserves"
)

type AtomKind interface {
	IsAtomKind()
}

func AtomKindFromPreserves(value Value) AtomKind {
	if o := AtomKindBooleanFromPreserves(value); o != nil {
		return o
	}
	if o := AtomKindDoubleFromPreserves(value); o != nil {
		return o
	}
	if o := AtomKindSignedIntegerFromPreserves(value); o != nil {
		return o
	}
	if o := AtomKindStringFromPreserves(value); o != nil {
		return o
	}
	if o := AtomKindByteStringFromPreserves(value); o != nil {
		return o
	}
	if o := AtomKindSymbolFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func AtomKindToPreserves(s AtomKind) Value {
	switch u := s.(type) {
	case *AtomKindBoolean:
		return AtomKindBooleanToPreserves(*u)
	case *AtomKindDouble:
		return AtomKindDoubleToPreserves(*u)
	case *AtomKindSignedInteger:
		return AtomKindSignedIntegerToPreserves(*u)
	case *AtomKindString:
		return AtomKindStringToPreserves(*u)
	case *AtomKindByteString:
		return AtomKindByteStringToPreserves(*u)
	case *AtomKindSymbol:
		return AtomKindSymbolToPreserves(*u)
	}
	return nil
}

type AtomKindBoolean struct {
}

func (*AtomKindBoolean) IsAtomKind() {
}
func AtomKindBooleanFromPreserves(value Value) *AtomKindBoolean {
	if v := SymbolFromPreserves(value); v != nil {
		if v.Equal(NewSymbol("Boolean")) {
			return &AtomKindBoolean{}
		}
	}
	return nil
}
func AtomKindBooleanToPreserves(l AtomKindBoolean) Value {
	return NewSymbol("Boolean")
}

type AtomKindDouble struct {
}

func (*AtomKindDouble) IsAtomKind() {
}
func AtomKindDoubleFromPreserves(value Value) *AtomKindDouble {
	if v := SymbolFromPreserves(value); v != nil {
		if v.Equal(NewSymbol("Double")) {
			return &AtomKindDouble{}
		}
	}
	return nil
}
func AtomKindDoubleToPreserves(l AtomKindDouble) Value {
	return NewSymbol("Double")
}

type AtomKindSignedInteger struct {
}

func (*AtomKindSignedInteger) IsAtomKind() {
}
func AtomKindSignedIntegerFromPreserves(value Value) *AtomKindSignedInteger {
	if v := SymbolFromPreserves(value); v != nil {
		if v.Equal(NewSymbol("SignedInteger")) {
			return &AtomKindSignedInteger{}
		}
	}
	return nil
}
func AtomKindSignedIntegerToPreserves(l AtomKindSignedInteger) Value {
	return NewSymbol("SignedInteger")
}

type AtomKindString struct {
}

func (*AtomKindString) IsAtomKind() {
}
func AtomKindStringFromPreserves(value Value) *AtomKindString {
	if v := SymbolFromPreserves(value); v != nil {
		if v.Equal(NewSymbol("String")) {
			return &AtomKindString{}
		}
	}
	return nil
}
func AtomKindStringToPreserves(l AtomKindString) Value {
	return NewSymbol("String")
}

type AtomKindByteString struct {
}

func (*AtomKindByteString) IsAtomKind() {
}
func AtomKindByteStringFromPreserves(value Value) *AtomKindByteString {
	if v := SymbolFromPreserves(value); v != nil {
		if v.Equal(NewSymbol("ByteString")) {
			return &AtomKindByteString{}
		}
	}
	return nil
}
func AtomKindByteStringToPreserves(l AtomKindByteString) Value {
	return NewSymbol("ByteString")
}

type AtomKindSymbol struct {
}

func (*AtomKindSymbol) IsAtomKind() {
}
func AtomKindSymbolFromPreserves(value Value) *AtomKindSymbol {
	if v := SymbolFromPreserves(value); v != nil {
		if v.Equal(NewSymbol("Symbol")) {
			return &AtomKindSymbol{}
		}
	}
	return nil
}
func AtomKindSymbolToPreserves(l AtomKindSymbol) Value {
	return NewSymbol("Symbol")
}

type Binding struct {
	Name    Symbol
	Pattern SimplePattern
}

func NewBinding(name Symbol, pattern SimplePattern) *Binding {
	return &Binding{Name: name, Pattern: pattern}
}
func BindingFromPreserves(value Value) *Binding {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("named")) {
				if p0 := SymbolFromPreserves(rec.Fields[0]); p0 != nil {
					if p1 := SimplePatternFromPreserves(rec.Fields[1]); p1 != nil {
						return &Binding{Name: *p0, Pattern: p1}
					}
				}
			}
		}
	}
	return nil
}
func BindingToPreserves(s Binding) Value {
	return &Record{Key: NewSymbol("named"), Fields: []Value{SymbolToPreserves(s.Name), SimplePatternToPreserves(s.Pattern)}}
}

type Bundle struct {
	Modules Modules
}

func NewBundle(modules Modules) *Bundle {
	return &Bundle{Modules: modules}
}
func BundleFromPreserves(value Value) *Bundle {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("bundle")) {
				if p0 := ModulesFromPreserves(rec.Fields[0]); p0 != nil {
					return &Bundle{Modules: *p0}
				}
			}
		}
	}
	return nil
}
func BundleToPreserves(s Bundle) Value {
	return &Record{Key: NewSymbol("bundle"), Fields: []Value{ModulesToPreserves(s.Modules)}}
}

type CompoundPattern interface {
	IsCompoundPattern()
}

func CompoundPatternFromPreserves(value Value) CompoundPattern {
	if o := CompoundPatternRecFromPreserves(value); o != nil {
		return o
	}
	if o := CompoundPatternTupleFromPreserves(value); o != nil {
		return o
	}
	if o := CompoundPatternTuplePrefixFromPreserves(value); o != nil {
		return o
	}
	if o := CompoundPatternDictFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func CompoundPatternToPreserves(s CompoundPattern) Value {
	switch u := s.(type) {
	case *CompoundPatternRec:
		return CompoundPatternRecToPreserves(*u)
	case *CompoundPatternTuple:
		return CompoundPatternTupleToPreserves(*u)
	case *CompoundPatternTuplePrefix:
		return CompoundPatternTuplePrefixToPreserves(*u)
	case *CompoundPatternDict:
		return CompoundPatternDictToPreserves(*u)
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
func CompoundPatternRecFromPreserves(value Value) *CompoundPatternRec {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("rec")) {
				if p0 := NamedPatternFromPreserves(rec.Fields[0]); p0 != nil {
					if p1 := NamedPatternFromPreserves(rec.Fields[1]); p1 != nil {
						return &CompoundPatternRec{Label: p0, Fields: p1}
					}
				}
			}
		}
	}
	return nil
}
func CompoundPatternRecToPreserves(s CompoundPatternRec) Value {
	return &Record{Key: NewSymbol("rec"), Fields: []Value{NamedPatternToPreserves(s.Label), NamedPatternToPreserves(s.Fields)}}
}

type CompoundPatternTuple struct {
	Patterns []NamedPattern
}

func NewCompoundPatternTuple(patterns []NamedPattern) *CompoundPatternTuple {
	return &CompoundPatternTuple{Patterns: patterns}
}
func (*CompoundPatternTuple) IsCompoundPattern() {
}
func CompoundPatternTupleFromPreserves(value Value) *CompoundPatternTuple {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("tuple")) {
				if seq, ok := rec.Fields[0].(*Sequence); ok {
					var p0 []NamedPattern
					for _, item := range *seq {
						if itemParsed := NamedPatternFromPreserves(item); itemParsed != nil {
							p0 = append(p0, itemParsed)
						} else {
							return nil
						}
					}
					return &CompoundPatternTuple{Patterns: p0}
				}
			}
		}
	}
	return nil
}
func CompoundPatternTupleToPreserves(s CompoundPatternTuple) Value {
	var p0 = &Sequence{}
	for _, k := range s.Patterns {
		*p0 = append(*p0, NamedPatternToPreserves(k))
	}
	return &Record{Key: NewSymbol("tuple"), Fields: []Value{p0}}
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
func CompoundPatternTuplePrefixFromPreserves(value Value) *CompoundPatternTuplePrefix {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("tuplePrefix")) {
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
	}
	return nil
}
func CompoundPatternTuplePrefixToPreserves(s CompoundPatternTuplePrefix) Value {
	var p0 = &Sequence{}
	for _, k := range s.Fixed {
		*p0 = append(*p0, NamedPatternToPreserves(k))
	}
	return &Record{Key: NewSymbol("tuplePrefix"), Fields: []Value{p0, NamedSimplePatternToPreserves(s.Variable)}}
}

type CompoundPatternDict struct {
	Entries DictionaryEntries
}

func NewCompoundPatternDict(entries DictionaryEntries) *CompoundPatternDict {
	return &CompoundPatternDict{Entries: entries}
}
func (*CompoundPatternDict) IsCompoundPattern() {
}
func CompoundPatternDictFromPreserves(value Value) *CompoundPatternDict {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("dict")) {
				if p0 := DictionaryEntriesFromPreserves(rec.Fields[0]); p0 != nil {
					return &CompoundPatternDict{Entries: *p0}
				}
			}
		}
	}
	return nil
}
func CompoundPatternDictToPreserves(s CompoundPatternDict) Value {
	return &Record{Key: NewSymbol("dict"), Fields: []Value{DictionaryEntriesToPreserves(s.Entries)}}
}

type Definition interface {
	IsDefinition()
}

func DefinitionFromPreserves(value Value) Definition {
	if o := DefinitionOrFromPreserves(value); o != nil {
		return o
	}
	if o := DefinitionAndFromPreserves(value); o != nil {
		return o
	}
	if o := DefinitionPatternFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func DefinitionToPreserves(s Definition) Value {
	switch u := s.(type) {
	case *DefinitionOr:
		return DefinitionOrToPreserves(*u)
	case *DefinitionAnd:
		return DefinitionAndToPreserves(*u)
	case *DefinitionPattern:
		return DefinitionPatternToPreserves(*u)
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
func DefinitionOrFromPreserves(value Value) *DefinitionOr {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "or" {
			if seq, ok := rec.Fields[0].(*Sequence); ok {
				var patterns []NamedAlternative
				for _, item := range *seq {
					if itemParsed := NamedAlternativeFromPreserves(item); itemParsed != nil {
						patterns = append(patterns, *itemParsed)
					} else {
						return nil
					}
				}
				return &DefinitionOr{Pattern0: patterns[0], Pattern1: patterns[1], PatternN: patterns[2:]}
			}
		}
	}
	return nil
}
func DefinitionOrToPreserves(s DefinitionOr) Value {
	var values []Value
	values = append(values, NamedAlternativeToPreserves(s.Pattern0))
	values = append(values, NamedAlternativeToPreserves(s.Pattern1))
	for _, v := range s.PatternN {
		values = append(values, NamedAlternativeToPreserves(v))
	}
	return &Record{Key: NewSymbol("or"), Fields: []Value{extras.Reference(Sequence(values))}}
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
func DefinitionAndFromPreserves(value Value) *DefinitionAnd {
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
func DefinitionAndToPreserves(s DefinitionAnd) Value {
	var values []Value
	values = append(values, NamedPatternToPreserves(s.Pattern0))
	values = append(values, NamedPatternToPreserves(s.Pattern1))
	for _, v := range s.PatternN {
		values = append(values, NamedPatternToPreserves(v))
	}
	return &Record{Key: NewSymbol("and"), Fields: []Value{extras.Reference(Sequence(values))}}
}

type DefinitionPattern struct {
	Pattern
}

func NewDefinitionPattern(obj Pattern) *DefinitionPattern {
	return &DefinitionPattern{Pattern: obj}
}
func (*DefinitionPattern) IsDefinition() {
}
func DefinitionPatternFromPreserves(value Value) *DefinitionPattern {
	if o := PatternFromPreserves(value); o != nil {
		return &DefinitionPattern{Pattern: o}
	}
	return nil
}
func DefinitionPatternToPreserves(s DefinitionPattern) Value {
	return PatternToPreserves(s.Pattern)
}

type Definitions map[Symbol]Definition

func NewDefinitions() Definitions {
	return make(Definitions)
}
func DefinitionsFromPreserves(value Value) *Definitions {
	if dict, ok := value.(*Dictionary); ok {
		obj := NewDefinitions()
		for dictKey, dictValue := range *dict {
			if dKey := SymbolFromPreserves(dictKey); dKey != nil {
				if dValue := DefinitionFromPreserves(dictValue); dValue != nil {
					obj[*dKey] = dValue
					continue
				}
			}
			return nil
		}
		return &obj
	}
	return nil
}
func DefinitionsToPreserves(m Definitions) Value {
	var dictionary = make(Dictionary)
	for dKey, dValue := range m {
		dictionary[SymbolToPreserves(dKey)] = DefinitionToPreserves(dValue)
	}
	return &dictionary
}

type DictionaryEntries map[Value]NamedSimplePattern

func NewDictionaryEntries() DictionaryEntries {
	return make(DictionaryEntries)
}
func DictionaryEntriesFromPreserves(value Value) *DictionaryEntries {
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
func DictionaryEntriesToPreserves(m DictionaryEntries) Value {
	var dictionary = make(Dictionary)
	for dKey, dValue := range m {
		dictionary[ValueToPreserves(dKey)] = NamedSimplePatternToPreserves(dValue)
	}
	return &dictionary
}

type EmbeddedTypeName interface {
	IsEmbeddedTypeName()
}

func EmbeddedTypeNameFromPreserves(value Value) EmbeddedTypeName {
	if o := EmbeddedTypeNameFalseFromPreserves(value); o != nil {
		return o
	}
	if o := EmbeddedTypeNameRefFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func EmbeddedTypeNameToPreserves(s EmbeddedTypeName) Value {
	switch u := s.(type) {
	case *EmbeddedTypeNameFalse:
		return EmbeddedTypeNameFalseToPreserves(*u)
	case *EmbeddedTypeNameRef:
		return EmbeddedTypeNameRefToPreserves(*u)
	}
	return nil
}

type EmbeddedTypeNameFalse struct {
}

func (*EmbeddedTypeNameFalse) IsEmbeddedTypeName() {
}
func EmbeddedTypeNameFalseFromPreserves(value Value) *EmbeddedTypeNameFalse {
	if v := BooleanFromPreserves(value); v != nil {
		if v.Equal(NewBoolean(false)) {
			return &EmbeddedTypeNameFalse{}
		}
	}
	return nil
}
func EmbeddedTypeNameFalseToPreserves(l EmbeddedTypeNameFalse) Value {
	return NewBoolean(false)
}

type EmbeddedTypeNameRef struct {
	Ref
}

func NewEmbeddedTypeNameRef(obj Ref) *EmbeddedTypeNameRef {
	return &EmbeddedTypeNameRef{Ref: obj}
}
func (*EmbeddedTypeNameRef) IsEmbeddedTypeName() {
}
func EmbeddedTypeNameRefFromPreserves(value Value) *EmbeddedTypeNameRef {
	if o := RefFromPreserves(value); o != nil {
		return &EmbeddedTypeNameRef{Ref: *o}
	}
	return nil
}
func EmbeddedTypeNameRefToPreserves(s EmbeddedTypeNameRef) Value {
	return RefToPreserves(s.Ref)
}

type ModulePath []Symbol

func NewModulePath(items ...Symbol) *ModulePath {
	_items := ModulePath(items)
	return &_items
}
func ModulePathFromPreserves(value Value) *ModulePath {
	if seq, ok := value.(*Sequence); ok {
		var items []Symbol
		for _, item := range *seq {
			if itemParsed := SymbolFromPreserves(item); itemParsed != nil {
				items = append(items, *itemParsed)
			} else {
				return nil
			}
		}
		_items := ModulePath(items)
		return &_items
	}
	return nil
}
func ModulePathToPreserves(d ModulePath) Value {
	var values []Value
	for _, v := range d {
		values = append(values, SymbolToPreserves(v))
	}
	s := Sequence(values)
	return &s
}
func (m ModulePath) Hash() (s string) {
	for i, v := range m {
		s += v.String()
		if len(m) != i+1 {
			s += "."
		}
	}
	return
}
func (m *ModulePath) ToHash() extras.Hash[ModulePath] {
	return extras.NewHash(*m)
}

type Modules map[extras.Hash[ModulePath]]Schema

func NewModules() Modules {
	return make(Modules)
}
func ModulesFromPreserves(value Value) *Modules {
	if dict, ok := value.(*Dictionary); ok {
		obj := NewModules()
		for dictKey, dictValue := range *dict {
			if dKey := ModulePathFromPreserves(dictKey); dKey != nil {
				if dValue := SchemaFromPreserves(dictValue); dValue != nil {
					obj[dKey.ToHash()] = *dValue
					continue
				}
			}
			return nil
		}
		return &obj
	}
	return nil
}
func ModulesToPreserves(m Modules) Value {
	var dictionary = make(Dictionary)
	for dKey, dValue := range m {
		dictionary[ModulePathToPreserves(dKey.FromHash())] = SchemaToPreserves(dValue)
	}
	return &dictionary
}

type NamedAlternative struct {
	VariantLabel Pstring
	Pattern      Pattern
}

func NewNamedAlternative(variantLabel Pstring, pattern Pattern) *NamedAlternative {
	return &NamedAlternative{VariantLabel: variantLabel, Pattern: pattern}
}
func NamedAlternativeFromPreserves(value Value) *NamedAlternative {
	if seq, ok := value.(*Sequence); ok && len(*seq) == 2 {
		if p0 := PstringFromPreserves((*seq)[0]); p0 != nil {
			if p1 := PatternFromPreserves((*seq)[1]); p1 != nil {
				return &NamedAlternative{VariantLabel: *p0, Pattern: p1}
			}
		}
	}
	return nil
}
func NamedAlternativeToPreserves(d NamedAlternative) Value {
	return &Sequence{PstringToPreserves(d.VariantLabel), PatternToPreserves(d.Pattern)}
}

type NamedPattern interface {
	IsNamedPattern()
}

func NamedPatternFromPreserves(value Value) NamedPattern {
	if o := NamedPatternNamedFromPreserves(value); o != nil {
		return o
	}
	if o := NamedPatternAnonymousFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func NamedPatternToPreserves(s NamedPattern) Value {
	switch u := s.(type) {
	case *NamedPatternNamed:
		return NamedPatternNamedToPreserves(*u)
	case *NamedPatternAnonymous:
		return NamedPatternAnonymousToPreserves(*u)
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
func NamedPatternNamedFromPreserves(value Value) *NamedPatternNamed {
	if o := BindingFromPreserves(value); o != nil {
		return &NamedPatternNamed{Binding: *o}
	}
	return nil
}
func NamedPatternNamedToPreserves(s NamedPatternNamed) Value {
	return BindingToPreserves(s.Binding)
}

type NamedPatternAnonymous struct {
	Pattern
}

func NewNamedPatternAnonymous(obj Pattern) *NamedPatternAnonymous {
	return &NamedPatternAnonymous{Pattern: obj}
}
func (*NamedPatternAnonymous) IsNamedPattern() {
}
func NamedPatternAnonymousFromPreserves(value Value) *NamedPatternAnonymous {
	if o := PatternFromPreserves(value); o != nil {
		return &NamedPatternAnonymous{Pattern: o}
	}
	return nil
}
func NamedPatternAnonymousToPreserves(s NamedPatternAnonymous) Value {
	return PatternToPreserves(s.Pattern)
}

type NamedSimplePattern interface {
	IsNamedSimplePattern()
}

func NamedSimplePatternFromPreserves(value Value) NamedSimplePattern {
	if o := NamedSimplePatternNamedFromPreserves(value); o != nil {
		return o
	}
	if o := NamedSimplePatternAnonymousFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func NamedSimplePatternToPreserves(s NamedSimplePattern) Value {
	switch u := s.(type) {
	case *NamedSimplePatternNamed:
		return NamedSimplePatternNamedToPreserves(*u)
	case *NamedSimplePatternAnonymous:
		return NamedSimplePatternAnonymousToPreserves(*u)
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
func NamedSimplePatternNamedFromPreserves(value Value) *NamedSimplePatternNamed {
	if o := BindingFromPreserves(value); o != nil {
		return &NamedSimplePatternNamed{Binding: *o}
	}
	return nil
}
func NamedSimplePatternNamedToPreserves(s NamedSimplePatternNamed) Value {
	return BindingToPreserves(s.Binding)
}

type NamedSimplePatternAnonymous struct {
	SimplePattern
}

func NewNamedSimplePatternAnonymous(obj SimplePattern) *NamedSimplePatternAnonymous {
	return &NamedSimplePatternAnonymous{SimplePattern: obj}
}
func (*NamedSimplePatternAnonymous) IsNamedSimplePattern() {
}
func NamedSimplePatternAnonymousFromPreserves(value Value) *NamedSimplePatternAnonymous {
	if o := SimplePatternFromPreserves(value); o != nil {
		return &NamedSimplePatternAnonymous{SimplePattern: o}
	}
	return nil
}
func NamedSimplePatternAnonymousToPreserves(s NamedSimplePatternAnonymous) Value {
	return SimplePatternToPreserves(s.SimplePattern)
}

type Pattern interface {
	IsPattern()
}

func PatternFromPreserves(value Value) Pattern {
	if o := PatternSimplePatternFromPreserves(value); o != nil {
		return o
	}
	if o := PatternCompoundPatternFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func PatternToPreserves(s Pattern) Value {
	switch u := s.(type) {
	case *PatternSimplePattern:
		return PatternSimplePatternToPreserves(*u)
	case *PatternCompoundPattern:
		return PatternCompoundPatternToPreserves(*u)
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
func PatternSimplePatternFromPreserves(value Value) *PatternSimplePattern {
	if o := SimplePatternFromPreserves(value); o != nil {
		return &PatternSimplePattern{SimplePattern: o}
	}
	return nil
}
func PatternSimplePatternToPreserves(s PatternSimplePattern) Value {
	return SimplePatternToPreserves(s.SimplePattern)
}

type PatternCompoundPattern struct {
	CompoundPattern
}

func NewPatternCompoundPattern(obj CompoundPattern) *PatternCompoundPattern {
	return &PatternCompoundPattern{CompoundPattern: obj}
}
func (*PatternCompoundPattern) IsPattern() {
}
func PatternCompoundPatternFromPreserves(value Value) *PatternCompoundPattern {
	if o := CompoundPatternFromPreserves(value); o != nil {
		return &PatternCompoundPattern{CompoundPattern: o}
	}
	return nil
}
func PatternCompoundPatternToPreserves(s PatternCompoundPattern) Value {
	return CompoundPatternToPreserves(s.CompoundPattern)
}

type Ref struct {
	Module ModulePath
	Name   Symbol
}

func NewRef(module ModulePath, name Symbol) *Ref {
	return &Ref{Module: module, Name: name}
}
func RefFromPreserves(value Value) *Ref {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("ref")) {
				if p0 := ModulePathFromPreserves(rec.Fields[0]); p0 != nil {
					if p1 := SymbolFromPreserves(rec.Fields[1]); p1 != nil {
						return &Ref{Module: *p0, Name: *p1}
					}
				}
			}
		}
	}
	return nil
}
func RefToPreserves(s Ref) Value {
	return &Record{Key: NewSymbol("ref"), Fields: []Value{ModulePathToPreserves(s.Module), SymbolToPreserves(s.Name)}}
}

type Schema struct {
	Definitions  Definitions
	EmbeddedType EmbeddedTypeName
	Version      Version
}

func NewSchema(definitions Definitions, embeddedType EmbeddedTypeName, version Version) *Schema {
	return &Schema{Definitions: definitions, EmbeddedType: embeddedType, Version: version}
}
func SchemaFromPreserves(value Value) *Schema {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("schema")) {
				if dict, ok := rec.Fields[0].(*Dictionary); ok {
					var obj Schema
					for dictKey, dictValue := range *dict {
						if dictKey.Equal(NewSymbol("definitions")) {
							if p0 := DefinitionsFromPreserves(dictValue); p0 != nil {
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
							if p2 := VersionFromPreserves(dictValue); p2 != nil {
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
	}
	return nil
}
func SchemaToPreserves(d Schema) Value {
	return &Record{Key: NewSymbol("schema"), Fields: []Value{&Dictionary{NewSymbol("definitions"): DefinitionsToPreserves(d.Definitions), NewSymbol("embeddedType"): EmbeddedTypeNameToPreserves(d.EmbeddedType), NewSymbol("version"): VersionToPreserves(d.Version)}}}
}

type SimplePattern interface {
	IsSimplePattern()
}

func SimplePatternFromPreserves(value Value) SimplePattern {
	if o := SimplePatternAnyFromPreserves(value); o != nil {
		return o
	}
	if o := SimplePatternAtomFromPreserves(value); o != nil {
		return o
	}
	if o := SimplePatternEmbeddedFromPreserves(value); o != nil {
		return o
	}
	if o := SimplePatternLitFromPreserves(value); o != nil {
		return o
	}
	if o := SimplePatternSeqofFromPreserves(value); o != nil {
		return o
	}
	if o := SimplePatternSetofFromPreserves(value); o != nil {
		return o
	}
	if o := SimplePatternDictofFromPreserves(value); o != nil {
		return o
	}
	if o := SimplePatternRefFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func SimplePatternToPreserves(s SimplePattern) Value {
	switch u := s.(type) {
	case *SimplePatternAny:
		return SimplePatternAnyToPreserves(*u)
	case *SimplePatternAtom:
		return SimplePatternAtomToPreserves(*u)
	case *SimplePatternEmbedded:
		return SimplePatternEmbeddedToPreserves(*u)
	case *SimplePatternLit:
		return SimplePatternLitToPreserves(*u)
	case *SimplePatternSeqof:
		return SimplePatternSeqofToPreserves(*u)
	case *SimplePatternSetof:
		return SimplePatternSetofToPreserves(*u)
	case *SimplePatternDictof:
		return SimplePatternDictofToPreserves(*u)
	case *SimplePatternRef:
		return SimplePatternRefToPreserves(*u)
	}
	return nil
}

type SimplePatternAny struct {
}

func (*SimplePatternAny) IsSimplePattern() {
}
func SimplePatternAnyFromPreserves(value Value) *SimplePatternAny {
	if v := SymbolFromPreserves(value); v != nil {
		if v.Equal(NewSymbol("any")) {
			return &SimplePatternAny{}
		}
	}
	return nil
}
func SimplePatternAnyToPreserves(l SimplePatternAny) Value {
	return NewSymbol("any")
}

type SimplePatternAtom struct {
	AtomKind AtomKind
}

func NewSimplePatternAtom(atomKind AtomKind) *SimplePatternAtom {
	return &SimplePatternAtom{AtomKind: atomKind}
}
func (*SimplePatternAtom) IsSimplePattern() {
}
func SimplePatternAtomFromPreserves(value Value) *SimplePatternAtom {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("atom")) {
				if p0 := AtomKindFromPreserves(rec.Fields[0]); p0 != nil {
					return &SimplePatternAtom{AtomKind: p0}
				}
			}
		}
	}
	return nil
}
func SimplePatternAtomToPreserves(s SimplePatternAtom) Value {
	return &Record{Key: NewSymbol("atom"), Fields: []Value{AtomKindToPreserves(s.AtomKind)}}
}

type SimplePatternEmbedded struct {
	Interface SimplePattern
}

func NewSimplePatternEmbedded(_interface SimplePattern) *SimplePatternEmbedded {
	return &SimplePatternEmbedded{Interface: _interface}
}
func (*SimplePatternEmbedded) IsSimplePattern() {
}
func SimplePatternEmbeddedFromPreserves(value Value) *SimplePatternEmbedded {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("embedded")) {
				if p0 := SimplePatternFromPreserves(rec.Fields[0]); p0 != nil {
					return &SimplePatternEmbedded{Interface: p0}
				}
			}
		}
	}
	return nil
}
func SimplePatternEmbeddedToPreserves(s SimplePatternEmbedded) Value {
	return &Record{Key: NewSymbol("embedded"), Fields: []Value{SimplePatternToPreserves(s.Interface)}}
}

type SimplePatternLit struct {
	Value Value
}

func NewSimplePatternLit(value Value) *SimplePatternLit {
	return &SimplePatternLit{Value: value}
}
func (*SimplePatternLit) IsSimplePattern() {
}
func SimplePatternLitFromPreserves(value Value) *SimplePatternLit {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("lit")) {
				return &SimplePatternLit{Value: rec.Fields[0]}
			}
		}
	}
	return nil
}
func SimplePatternLitToPreserves(s SimplePatternLit) Value {
	return &Record{Key: NewSymbol("lit"), Fields: []Value{s.Value}}
}

type SimplePatternSeqof struct {
	Pattern SimplePattern
}

func NewSimplePatternSeqof(pattern SimplePattern) *SimplePatternSeqof {
	return &SimplePatternSeqof{Pattern: pattern}
}
func (*SimplePatternSeqof) IsSimplePattern() {
}
func SimplePatternSeqofFromPreserves(value Value) *SimplePatternSeqof {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("seqof")) {
				if p0 := SimplePatternFromPreserves(rec.Fields[0]); p0 != nil {
					return &SimplePatternSeqof{Pattern: p0}
				}
			}
		}
	}
	return nil
}
func SimplePatternSeqofToPreserves(s SimplePatternSeqof) Value {
	return &Record{Key: NewSymbol("seqof"), Fields: []Value{SimplePatternToPreserves(s.Pattern)}}
}

type SimplePatternSetof struct {
	Pattern SimplePattern
}

func NewSimplePatternSetof(pattern SimplePattern) *SimplePatternSetof {
	return &SimplePatternSetof{Pattern: pattern}
}
func (*SimplePatternSetof) IsSimplePattern() {
}
func SimplePatternSetofFromPreserves(value Value) *SimplePatternSetof {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("setof")) {
				if p0 := SimplePatternFromPreserves(rec.Fields[0]); p0 != nil {
					return &SimplePatternSetof{Pattern: p0}
				}
			}
		}
	}
	return nil
}
func SimplePatternSetofToPreserves(s SimplePatternSetof) Value {
	return &Record{Key: NewSymbol("setof"), Fields: []Value{SimplePatternToPreserves(s.Pattern)}}
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
func SimplePatternDictofFromPreserves(value Value) *SimplePatternDictof {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("dictof")) {
				if p0 := SimplePatternFromPreserves(rec.Fields[0]); p0 != nil {
					if p1 := SimplePatternFromPreserves(rec.Fields[1]); p1 != nil {
						return &SimplePatternDictof{Key: p0, Value: p1}
					}
				}
			}
		}
	}
	return nil
}
func SimplePatternDictofToPreserves(s SimplePatternDictof) Value {
	return &Record{Key: NewSymbol("dictof"), Fields: []Value{SimplePatternToPreserves(s.Key), SimplePatternToPreserves(s.Value)}}
}

type SimplePatternRef struct {
	Ref
}

func NewSimplePatternRef(obj Ref) *SimplePatternRef {
	return &SimplePatternRef{Ref: obj}
}
func (*SimplePatternRef) IsSimplePattern() {
}
func SimplePatternRefFromPreserves(value Value) *SimplePatternRef {
	if o := RefFromPreserves(value); o != nil {
		return &SimplePatternRef{Ref: *o}
	}
	return nil
}
func SimplePatternRefToPreserves(s SimplePatternRef) Value {
	return RefToPreserves(s.Ref)
}

type Version struct {
}

func VersionFromPreserves(value Value) *Version {
	if v := SignedIntegerFromPreserves(value); v != nil {
		if v.Equal(NewSignedInteger("1")) {
			return &Version{}
		}
	}
	return nil
}
func VersionToPreserves(l Version) Value {
	return NewSignedInteger("1")
}
