package lib

import "fmt"

var (
	AND            = NewSymbol("and")
	ANY            = NewSymbol("any")
	ATOM           = NewSymbol("atom")
	BOOLEAN        = NewSymbol("Boolean")
	BUNDLE         = NewSymbol("bundle")
	BYTE_STRING    = NewSymbol("ByteString")
	DEFINITIONS    = NewSymbol("definitions")
	DICT           = NewSymbol("dict")
	DICTOF         = NewSymbol("dictof")
	DOUBLE         = NewSymbol("Double")
	EMBEDDED       = NewSymbol("embedded")
	LIT            = NewSymbol("lit")
	NAMED          = NewSymbol("named")
	OR             = NewSymbol("or")
	REC            = NewSymbol("rec")
	REF            = NewSymbol("ref")
	SCHEMA         = NewSymbol("schema")
	SEQOF          = NewSymbol("seqof")
	SETOF          = NewSymbol("setof")
	SIGNED_INTEGER = NewSymbol("SignedInteger")
	STRING         = NewSymbol("String")
	SYMBOL         = NewSymbol("Symbol")
	TUPLE          = NewSymbol("tuple")
	TUPLE_PREFIX   = NewSymbol("tuplePrefix")
	VERSION        = NewSymbol("version")
)

type SchemaDecodeFailed struct {
	Class    any
	Pattern  Value
	Value    Value
	Failures []error
}

func NewSchemaDecodeFailed(cls any, pattern, value Value, errors ...error) SchemaDecodeFailed {
	return SchemaDecodeFailed{
		Class:    cls,
		Pattern:  pattern,
		Value:    value,
		Failures: errors,
	}
}
func (s SchemaDecodeFailed) Error() (err string) {
	err = fmt.Sprintf("schemaDecodeFailed (%v, %v, %v): ", s.Class, s.Pattern, s.Value)
	for _, e := range s.Failures {
		err = fmt.Sprintf("%s: %v", err, e)
	}
	return
}

// Base class for classes representing grammatical productions in a schema: instances of
// [SchemaObject][preserves.schema.SchemaObject] represent schema *definitions*. This is an
// abstract class, as are its subclasses [Enumeration][preserves.schema.Enumeration] and
// [Definition][preserves.schema.Definition]. It is subclasses of *those* subclasses,
// automatically produced during schema loading, that are actually instantiated.
type SchemaObject struct {
	// A [Namespace][preserves.schema.Namespace] that is the top-level environment for all
	// bundles included in the [Compiler][preserves.schema.Compiler] run that produced this
	// [SchemaObject][preserves.schema.SchemaObject].
	RootNS string
	// A `Value` conforming to schema `meta.Definition` (and thus often to `meta.Pattern`
	// etc.), interpreted by the [SchemaObject][preserves.schema.SchemaObject] machinery to drive
	// parsing, unparsing and so forth.
	Schema string
	// A sequence (tuple) of [Symbol][preserves.values.Symbol]s naming the path from the root
	// to the schema module containing this definition.
	ModulePath string
	// A [Symbol][preserves.values.Symbol] naming this definition within its module.
	Name string
	/// `None` for [Definition][preserves.schema.Definition]s (such as
	// `bundle.stream.StreamListenerError` above) and for overall
	// [Enumeration][preserves.schema.Enumeration]s (such as `bundle.stream.Mode`), or a
	// [Symbol][preserves.values.Symbol] for variant definitions *contained within* an enumeration
	// (such as `bundle.stream.Mode.packet`).
	Variant string
}

func (s SchemaObject) Decode(v Value) (string, error) {
	panic(NotImplementedError(fmt.Errorf("implementers responsibility")))
}

func (s SchemaObject) TryDecode(v Value) (string, error) {
	r, err := s.Decode(v)
	if _, ok := err.(SchemaDecodeFailed); ok {
		return "", nil
	}
	return r, err
}

func (s SchemaObject) Parse(p, v Value, args ...[]string) Value {
	switch item := p.(type) {
	case *Symbol:
		if item.Equal(&ANY) {
			return v
		}
		return v
	case *Record:
		return v
	}
	return v
}
