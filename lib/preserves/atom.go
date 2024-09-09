package preserves

import (
	"fmt"
	"io"
	"math/big"
	"slices"
	"strings"
)

type Boolean bool

func NewBoolean(b bool) *Boolean {
	return &([]Boolean{Boolean(b)}[0])
}
func (b Boolean) New() Value {
	return new(Boolean)
}

func (b Boolean) Equal(y Value) bool {
	x, ok := y.(*Boolean)
	if !ok {
		return false
	}
	return b == *x
}
func (b Boolean) Cmp(y Value) int {
	// Boolean is the lowest Value.
	switch u := y.(type) {
	case *Boolean:
		if b == *u {
			return 0
		}
		if b && !(*u) {
			return 1
		}
	}

	return -1
}

// io.WriterTo interface
func (b Boolean) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("dictionary: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (b Boolean) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("dictionary: readfrom: not implemented here")
	return
}

type Double float64

func NewDouble(f float64) *Double {
	return &([]Double{Double(f)}[0])
}

func (d Double) New() Value {
	return new(Double)
}

func (d Double) Equal(y Value) bool {
	x, ok := y.(*Double)
	if !ok {
		return false
	}
	return d == *x
}
func (d Double) Cmp(y Value) int {
	switch u := y.(type) {
	case *Boolean:
		return 1
	case *Double:
		if d > *u {
			return 1
		}
		if d == *u {
			return 0
		}
	}
	return -1
}

// io.WriterTo interface
func (d Double) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("dictionary: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (d Double) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("dictionary: readfrom: not implemented here")
	return
}

type SignedInteger big.Int

func NewSignedInteger(s string) *SignedInteger {
	a := big.NewInt(0)
	a.SetString(s, 10)
	b := SignedInteger(*a)
	return &b
}
func (s SignedInteger) New() Value {
	return new(SignedInteger)
}

func (s SignedInteger) Equal(y Value) bool {
	x, ok := y.(*SignedInteger)
	if !ok {
		return false
	}
	a := big.Int(s)
	b := big.Int(*x)
	return a.Cmp(&b) == 0
}
func (s SignedInteger) Cmp(y Value) int {
	switch u := y.(type) {
	case *Boolean:
		return 1
	case *Double:
		return 1
	case *SignedInteger:
		a := big.Int(s)
		b := big.Int(*u)
		return a.Cmp(&b)
	}
	return -1
}

// io.WriterTo interface
func (s SignedInteger) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("dictionary: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (s SignedInteger) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("dictionary: readfrom: not implemented here")
	return
}

type Pstring string

func NewPstring(s string) *Pstring {
	return &([]Pstring{Pstring(s)}[0])
}
func (p Pstring) New() Value {
	return new(Pstring)
}

func (p Pstring) Equal(y Value) bool {
	x, ok := y.(Pstring)
	if !ok {
		return false
	}
	return p == x
}
func (p Pstring) Cmp(y Value) int {
	switch u := y.(type) {
	case *Boolean:
		return 1
	case *SignedInteger:
		return 1
	case *Double:
		return 1
	case *Pstring:
		return strings.Compare(string(p), string(*u))
	}
	return -1
}

// io.WriterTo interface
func (p Pstring) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("dictionary: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (p Pstring) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("dictionary: readfrom: not implemented here")
	return
}

type ByteString []byte

func NewByteString(s string) *ByteString {
	return &([]ByteString{ByteString([]byte(s))}[0])
}
func (b ByteString) New() Value {
	return new(ByteString)
}

func (b ByteString) Equal(y Value) bool {
	x, ok := y.(*ByteString)
	if !ok {
		return false
	}
	return slices.Equal(b, *x)
}

func (b ByteString) Cmp(y Value) int {
	switch u := y.(type) {
	case *Boolean:
		return 1
	case *SignedInteger:
		return 1
	case *Double:
		return 1
	case *Pstring:
		return 1
	case *ByteString:
		return slices.Compare(b, *u)
	}
	return -1
}

// io.WriterTo interface
func (b ByteString) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("dictionary: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (b ByteString) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("dictionary: readfrom: not implemented here")
	return
}

type Symbol string

func NewSymbol(s string) *Symbol {
	return &([]Symbol{Symbol(s)}[0])
}

func (s Symbol) New() Value {
	return new(Symbol)
}

func (s Symbol) Equal(y Value) bool {
	x, ok := y.(*Symbol)
	if !ok {
		return false
	}
	return s == *x
}
func (s Symbol) Cmp(y Value) int {
	switch u := y.(type) {
	case *Boolean:
		return 1
	case *SignedInteger:
		return 1
	case *Double:
		return 1
	case *Pstring:
		return 1
	case *ByteString:
		return 1
	case *Symbol:
		return strings.Compare(string(s), string(*u))
	}
	return -1
}

// io.WriterTo interface
func (s Symbol) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("dictionary: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (s Symbol) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("dictionary: readfrom: not implemented here")
	return
}

func (s Symbol) String() string {
	return string(s)
}

func BooleanFromPreserves(value Value) *Boolean {
	if obj, ok := value.(*Boolean); ok {
		return obj
	}
	return nil
}
func BooleanToPreserves(v Boolean) Value {
	return &v
}
func SignedIntegerFromPreserves(value Value) *SignedInteger {
	if obj, ok := value.(*SignedInteger); ok {
		return obj
	}
	return nil
}
func SignedIntegerToPreserves(v SignedInteger) Value {
	return &v
}
func PstringFromPreserves(value Value) *Pstring {
	if obj, ok := value.(*Pstring); ok {
		return obj
	}
	return nil
}
func PstringToPreserves(v Pstring) Value {
	return &v
}
func SymbolFromPreserves(value Value) *Symbol {
	if obj, ok := value.(*Symbol); ok {
		return obj
	}
	return nil
}
func SymbolToPreserves(v Symbol) Value {
	return &v
}
func ValueFromPreserves(value Value) Value {
	if a := BooleanFromPreserves(value); a != nil {
		return a
	}
	if a := SignedIntegerFromPreserves(value); a != nil {
		return a
	}
	if a := PstringFromPreserves(value); a != nil {
		return a
	}
	if a := SymbolFromPreserves(value); a != nil {
		return a
	}
	return nil
}
func ValueToPreserves(v Value) Value {
	return v
}
