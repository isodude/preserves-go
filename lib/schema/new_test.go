package schema

import (
	"testing"

	"github.com/isodude/preserves-go/lib/goast"
	. "github.com/isodude/preserves-go/lib/preserves"
	"github.com/kylelemons/godebug/diff"
)

func TestDefinition(t *testing.T) {
	definitions := map[string]struct {
		definition map[string]Definition
		result     string
	}{
		"Test": {map[string]Definition{"Test": NewDefinitionOr(
			*NewNamedAlternative(Pstring("first"), NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("test0")))),
			*NewNamedAlternative(Pstring("second"), NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("test1")))),
			[]NamedAlternative{
				*NewNamedAlternative(Pstring("third"), NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("testN")))),
			},
		)},
			`package beep

import (
	"github.com/isodude/preserves-go/lib/extras"
	. "github.com/isodude/preserves-go/lib/preserves"
)

type Test interface {
	IsTest()
}

func TestFromPreserves(value Value) (Test) {
	if o := TestFirstFromPreserves(value); o != nil {
		return o
	}
	if o := TestSecondFromPreserves(value); o != nil {
		return o
	}
	if o := TestThirdFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func TestToPreserves(s Test) (Value) {
	switch u := s.(type) {
	case *TestFirst:
		return TestFirstToPreserves(*u)
	case *TestSecond:
		return TestSecondToPreserves(*u)
	case *TestThird:
		return TestThirdToPreserves(*u)
	}
	return nil
}

type TestFirst struct {
}

func NewTestFirst() (*TestFirst) {
	return &TestFirst{}
}
func (*TestFirst) IsTest() {
}

type TestFirst struct {
}

func (*TestFirst) IsTest() {
}
func TestFirstFromPreserves(value Value) (*TestFirst) {
	if v := SymbolFromPreserves(value); v != nil {
		if v.Equal(NewSymbol("test0")) {
			return &TestFirst{}
		}
	}
	return nil
}
func TestFirstToPreserves(l TestFirst) (Value) {
	return NewSymbol("test0")
}

type TestSecond struct {
}

func NewTestSecond() (*TestSecond) {
	return &TestSecond{}
}
func (*TestSecond) IsTest() {
}

type TestSecond struct {
}

func (*TestSecond) IsTest() {
}
func TestSecondFromPreserves(value Value) (*TestSecond) {
	if v := SymbolFromPreserves(value); v != nil {
		if v.Equal(NewSymbol("test1")) {
			return &TestSecond{}
		}
	}
	return nil
}
func TestSecondToPreserves(l TestSecond) (Value) {
	return NewSymbol("test1")
}

type TestThird struct {
}

func NewTestThird() (*TestThird) {
	return &TestThird{}
}
func (*TestThird) IsTest() {
}

type TestThird struct {
}

func (*TestThird) IsTest() {
}
func TestThirdFromPreserves(value Value) (*TestThird) {
	if v := SymbolFromPreserves(value); v != nil {
		if v.Equal(NewSymbol("testN")) {
			return &TestThird{}
		}
	}
	return nil
}
func TestThirdToPreserves(l TestThird) (Value) {
	return NewSymbol("testN")
}
`,
		},
		"NamedAlternative": {
			map[string]Definition{"NamedAlternative": NewDefinitionPattern(NewPatternCompoundPattern(NewCompoundPatternTuple(
				[]NamedPattern{
					NewNamedPatternNamed(*NewBinding(*NewSymbol("variantLabel"), NewSimplePatternAtom(&AtomKindString{}))),
					NewNamedPatternNamed(*NewBinding(*NewSymbol("pattern"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Pattern"))))),
				},
			)))},
			`package beep

import (
	"github.com/isodude/preserves-go/lib/extras"
	. "github.com/isodude/preserves-go/lib/preserves"
)

type NamedAlternative struct {
	VariantLabel	AtomKindString
	Pattern		Pattern
}

func NewNamedAlternative(variantLabel AtomKindString, pattern Pattern) (*NamedAlternative) {
	return &NamedAlternative{VariantLabel: variantLabel, Pattern: pattern}
}
func NamedAlternativeFromPreserves(value Value) (*NamedAlternative) {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 2 {
		if sym, ok := rec.Key.(*Symbol); ok && sym.String() == "NamedAlternative" {
			// The variables of the struct do not have a corresponding ToPreserves/FromPreserves
			return nil
		}
	}
	return nil
}
func NamedAlternativeToPreserves(s NamedAlternative) (Value) {
	return &Record{Key: NewSymbol("NamedAlternative"), Fields: []Value{}}
}
`,
		},
	}

	for name, definition := range definitions {
		defs := []goast.AST{}
		for n, d := range definition.definition {
			defs = append(defs, DefinitionGenerator(n, d))
		}
		result := goast.Encode("beep", defs)
		if result != definition.result {
			t.Fatalf("%s not equal: %v", name, diff.Diff(result, definition.result))
		}
	}
}
