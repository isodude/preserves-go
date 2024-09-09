package preserves

import (
	"io"
)

type Comparable interface {
	Equal(Value) bool
	Cmp(Value) int
}

type New interface {
	New() Value
}

type Value interface {
	Comparable
	New
	io.WriterTo
	io.ReaderFrom
}

type Atom interface {
	Boolean | Double | SignedInteger | Pstring | ByteString | Symbol | Comment
}

type Compound interface {
	Record | Sequence | Set | Dictionary | Annotation | Embedded
}

type Parser interface {
	io.ReaderFrom
	GetValue() Value
}
