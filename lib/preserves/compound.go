package preserves

import (
	"fmt"
	"io"
	"slices"
)

type Record struct {
	Key    Value
	Fields []Value
}

// io.WriterTo interface
func (re Record) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("record: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (re Record) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("record: readfrom: not implemented here")
	return
}

// New interface
func (re Record) New() Value {
	return new(Record)
}

// Comparable interface
func (re Record) Equal(y Value) (b bool) {
	x, ok := y.(*Record)
	if !ok {
		return false
	}
	c1, ok := (re).Key.(Comparable)
	if !ok {
		return false
	}
	if !c1.Equal(*x) {
		return false
	}
	return slices.EqualFunc[[]Value]((re).Fields, (*x).Fields, func(e1, e2 Value) bool {
		c1, ok := e1.(Comparable)
		if !ok {
			return false
		}
		return c1.Equal(e2)
	})
}
func (r Record) Cmp(y Value) int {
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
		return 1
	case *Record:
		if c := r.Key.Cmp((*u).Key); c != 0 {
			return c
		}
		return slices.CompareFunc(r.Fields, (*u).Fields, func(e1, e2 Value) int {
			return e1.Cmp(e2)
		})
	}
	return -1
}

type Sequence []Value

func (s Sequence) New() Value { return new(Sequence) }

func (s Sequence) Equal(y Value) bool {
	x, ok := y.(*Sequence)
	if !ok {
		return false
	}
	for _, u := range []Value(s) {
		for _, v := range []Value(*x) {
			if !u.Equal(v) {
				return false
			}
		}
	}
	return true
}

func (s Sequence) Cmp(y Value) int {
	switch u := y.(type) {
	case *Boolean, *SignedInteger, *Double, *Pstring, *ByteString, *Symbol, *Record:
		return 1
	case *Sequence:
		return slices.CompareFunc(s, *u, func(e1, e2 Value) int {
			return e1.Cmp(e2)
		})
	}
	return -1
}

// io.WriterTo interface
func (s Sequence) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("record: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (s Sequence) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("record: readfrom: not implemented here")
	return
}

type Set map[Value]struct{}

func (s Set) New() Value {
	return &([]Set{Set(map[Value]struct{}{})}[0])
}

func (s Set) Equal(y Value) bool {
	x, ok := y.(*Set)
	if !ok {
		return false
	}
	if len(s) != len(*x) {
		return false
	}
	// TODO: Very inefficient way of doing this
	// Value should be represented by a hash instead
	for u := range s {
		equal := false
		for v := range *x {
			if v.Equal(u) {
				equal = true
			}
		}
		if !equal {
			return false
		}
	}
	return true
}

func (s Set) Cmp(y Value) int {
	switch u := y.(type) {
	case *Boolean, *SignedInteger, *Double, *Pstring, *ByteString, *Symbol, *Record, *Sequence:
		return 1
	case *Set:
		var k1, k2 []Value
		for k := range s {
			k1 = append(k1, k)
		}
		for k := range *u {
			k2 = append(k1, k)
		}
		return slices.CompareFunc(k1, k2, func(e1, e2 Value) int {
			return e1.Cmp(e2)
		})
	}
	return -1
}

// io.WriterTo interface
func (s Set) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("record: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (s Set) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("record: readfrom: not implemented here")
	return
}

type Dictionary map[Value]Value

// io.WriterTo interface
func (d Dictionary) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("dictionary: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (d Dictionary) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("dictionary: readfrom: not implemented here")
	return
}

// New interface
func (d Dictionary) New() Value {
	m := make(map[Value]Value)
	n := Dictionary(m)
	return &n
}

// Comparable interface
func (d Dictionary) Equal(y Value) (b bool) {

	x, ok := y.(*Dictionary)
	if !ok {
		return false
	}
	if len(d) != len(*x) {
		return false
	}

	// TODO: Inefficient way of doing it
	for u, t := range d {
		equal := false
		for v, z := range *x {
			if u.(Comparable).Equal(v) {
				if t.(Comparable).Equal(z) {
					equal = true
				}
			}
		}
		if !equal {
			return false
		}
	}
	return true
}
func (d Dictionary) Cmp(y Value) int {
	switch u := y.(type) {
	case *Boolean, *SignedInteger, *Double, *Pstring, *ByteString, *Symbol, *Record, *Sequence, *Set:
		return 1
	case *Dictionary:
		if len(d) < len(*u) {
			return -1
		}
		if len(d) > len(*u) {
			return 1
		}
		var cmp int
		for a, b := range d {
			for k, v := range *u {
				c1 := a.Cmp(k)
				c2 := b.Cmp(v)
				cmp += c1 + c2
			}
		}
		return cmp
	}
	return -1
}
func (d Dictionary) Set(k, v Value) {
	d[k] = v
}

func (d Dictionary) Get(k Value) (Value, bool) {
	for key, value := range d {
		if key.Equal(k) {

			return value, true
		}
	}
	return nil, false
}

func (d Dictionary) Delete(k Value) bool {
	for key := range d {
		if key.Equal(k) {
			delete(d, key)
			return true
		}
	}
	return false
}

// sort.Interface interface is not really doable with maps.
// Dictionary however can be sorted within slices and compared to others.
// Thus Dictionary should not implement sort.Interface
//
// func (d *Dictionary) Len() int {
//	return len(*d)
// }
//
// func (d *Dictionary) Less(i, j int) bool {
// 	return false
// }

// func (d *Dictionary) Swap(i, j int) {
// 	return
// }

type Annotation struct {
	Value          Value
	AnnotatedValue Value
}

func (a Annotation) New() Value {
	return new(Annotation)
}

func (a Annotation) Equal(y Value) (b bool) {
	return a.AnnotatedValue.Equal(y)
}
func (Annotation) Cmp(y Value) int {
	switch y.(type) {
	case *Boolean, *SignedInteger, *Double, *Pstring, *ByteString, *Symbol, *Record, *Sequence, *Set, *Dictionary:
		return 1
	case *Annotation:
		return 0
	}
	return -1
}

// io.WriterTo interface
func (a Annotation) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("dictionary: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (a Annotation) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("dictionary: readfrom: not implemented here")
	return
}

type Embedded struct {
	Value Value
}

func (e Embedded) New() Value {
	return new(Embedded)
}

func (e Embedded) Equal(y Value) (b bool) {
	x, ok := y.(*Embedded)
	if !ok {
		return false
	}
	return e.Equal(x)
}
func (Embedded) Cmp(y Value) int {
	switch y.(type) {
	case *Boolean, *SignedInteger, *Double, *Pstring, *ByteString, *Symbol, *Record, *Sequence, *Set, *Dictionary, *Annotation, *Comment:
		return 1
	case *Embedded:
		return 0
	}
	return -1
}

// io.WriterTo interface
func (e Embedded) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("dictionary: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (e Embedded) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("dictionary: readfrom: not implemented here")
	return
}

type Comment struct {
	Value          Pstring
	AnnotatedValue Value
}

func (c Comment) New() Value {
	return new(Comment)
}

func (c Comment) Equal(y Value) (b bool) {
	x, ok := y.(*Comment)
	if !ok {
		return false
	}
	return c.Equal(x)
}

func (Comment) Cmp(y Value) int {
	switch y.(type) {
	case *Boolean, *SignedInteger, *Double, *Pstring, *ByteString, *Symbol, *Record, *Sequence, *Set, *Dictionary, *Annotation:
		return 1
	case *Comment:
		return 0
	}
	return -1
}

// io.WriterTo interface
func (c Comment) WriteTo(w io.Writer) (n int64, err error) {
	err = fmt.Errorf("dictionary: writeto: not implemented here")
	return
}

// io.ReaderFrom interface
func (c Comment) ReadFrom(ir io.Reader) (n int64, err error) {
	err = fmt.Errorf("dictionary: readfrom: not implemented here")
	return
}
