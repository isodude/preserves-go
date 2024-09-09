package binary

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/isodude/preserves-go/lib/preserves"
	log "github.com/sirupsen/logrus"
)

const (
	BooleanFalseByte = 0x80
	BooleanTrueByte  = 0x81
	Reserved82Byte   = 0x82
	Reserved83Byte   = 0x83
	EndByte          = 0x84
	AnnotationByte   = 0x85
	EmbeddedByte     = 0x86
	DoubleByte       = 0x87
	// up to 0xAF
	SignedIntegerByte = 0xB0
	StringByte        = 0xB1
	ByteStringByte    = 0xB2
	SymbolByte        = 0xB3
	RecordByte        = 0xB4
	SequenceByte      = 0xB5
	SetByte           = 0xB6
	DictionaryByte    = 0xB7
	ReservedB8Byte    = 0xB8
	ReservedB9Byte    = 0xB9
	ReservedBAByte    = 0xBA
	ReservedBBByte    = 0xBB
	ReservedBCByte    = 0xBC
	ReservedBDByte    = 0xBD
	ReservedBEByte    = 0xBE
	ReservedBFByte    = 0xBF
)

var (
	lookup = map[byte]preserves.New{
		BooleanTrueByte:   new(Boolean),
		BooleanFalseByte:  new(Boolean),
		DoubleByte:        new(Double),
		SignedIntegerByte: new(SignedInteger),
		ByteStringByte:    new(ByteString),
		StringByte:        new(String),
		SymbolByte:        new(Symbol),
		DictionaryByte:    new(Dictionary),
		SetByte:           new(Set),
		SequenceByte:      new(Sequence),
		RecordByte:        new(Record),
		EmbeddedByte:      new(Embedded),
		AnnotationByte:    new(Annotation),
	}
)

func ReadEndByte(ir io.Reader) (ok bool, n int, err error) {
	var (
		br *bufio.Reader
		d  byte
	)
	br, ok = ir.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(ir)
	}
	ok = false

	d, err = br.ReadByte()
	n = 1

	log.Debugf("endbyte: %x, %v", d, err)

	if err == io.EOF && d == EndByte {
		ok = true
		err = nil
		return
	}

	if err != nil {
		return
	}

	if d == EndByte {
		ok = true
		return
	}

	err = br.UnreadByte()
	n = 0
	return
}

type ByteLenReader struct {
	io.ByteReader
	readBytes int
}

func (b *ByteLenReader) ReadByte() (d byte, err error) {
	d, err = b.ByteReader.ReadByte()
	b.readBytes += 1
	return
}

func ReadUvarint(r io.Reader) (l uint64, n int, err error) {
	br, ok := r.(io.ByteReader)
	if !ok {
		br = bufio.NewReader(r)
	}
	bl := &ByteLenReader{ByteReader: br}
	l, err = binary.ReadUvarint(bl)
	n = bl.readBytes
	return
}

type BinaryParser struct {
	Result preserves.Value
}

func (bp *BinaryParser) ReadFrom(ir io.Reader) (n int64, err error) {
	br, ok := ir.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(ir)
	}

	var bs []byte
	bs, err = br.Peek(1)
	n += 1
	if err != nil {
		log.Debugf("binaryparser: %x, err: %v", bs, err)
		return
	}

	new, ok := lookup[bs[0]]
	if !ok {
		log.Debugf("binaryparser: %x: miss", bs)
		err = fmt.Errorf("binaryparser: invalid rune: %x", bs[0])
		return
	}
	log.Debugf("binaryparser: %x: hit: %v", bs, reflect.TypeOf(new))

	bp.Result = new.New()
	var m int64
	m, err = bp.Result.ReadFrom(br)
	n += m

	return
}

func ToPreserves(v preserves.Value) preserves.Value {
	switch c := v.(type) {
	case *preserves.Symbol:
		return c
	case *Boolean:
		return &c.Boolean
	case *SignedInteger:
		return &c.SignedInteger
	case *String:
		return &c.Pstring
	case *preserves.Pstring:
		return c
	case *ByteString:
		return &c.ByteString
	case *Double:
		return &c.Double
	case *Symbol:
		return &c.Symbol
	case *Record:
		r := &preserves.Record{}
		r.Key = ToPreserves(c.Key)
		for _, v := range c.Fields {
			r.Fields = append(r.Fields, ToPreserves(v))
		}
		return r
	case *preserves.Record:
		r := &preserves.Record{}
		r.Key = ToPreserves(c.Key)
		for _, v := range c.Fields {
			r.Fields = append(r.Fields, ToPreserves(v))
		}
		return r
	case *Set:
		s := make(preserves.Set)
		for k := range c.Set {
			s[ToPreserves(k)] = struct{}{}
		}
		return &s
	case *Sequence:
		s := &preserves.Sequence{}
		for _, v := range c.Sequence {
			*s = append(*s, ToPreserves(v))
		}
		return s
	case *Dictionary:
		d := make(preserves.Dictionary)
		for k, v := range c.Dictionary {
			d[ToPreserves(k)] = ToPreserves(v)
		}
		return &d
	case *Annotation:
		a := &preserves.Annotation{}
		a.Value = ToPreserves(c.Value)
		a.AnnotatedValue = ToPreserves(c.AnnotatedValue)
		return a
	case *preserves.Annotation:
		a := &preserves.Annotation{}
		a.Value = ToPreserves(c.Value)
		a.AnnotatedValue = ToPreserves(c.AnnotatedValue)
		return a
	case *Embedded:
		e := &preserves.Embedded{}
		e.Value = ToPreserves(c.Value)
		return e
	}
	log.Fatalf("unknown type(2): %v", reflect.TypeOf(v))
	return nil
}

func FromPreserves(v preserves.Value) preserves.Value {
	switch c := v.(type) {
	case *preserves.Symbol:
		return &Symbol{Symbol: *c}
	case *preserves.Pstring:
		return &String{Pstring: *c}
	case *preserves.SignedInteger:
		return &SignedInteger{SignedInteger: *c}
	case *preserves.Boolean:
		return &Boolean{Boolean: *c}
	case *preserves.Dictionary:
		d := &Dictionary{Dictionary: make(preserves.Dictionary)}
		for k, v := range *c {
			d.Dictionary[FromPreserves(k)] = FromPreserves(v)
		}
		return d
	case *preserves.Sequence:
		s := &Sequence{}
		for _, v := range *c {
			s.Sequence = append(s.Sequence, FromPreserves(v))
		}
		return s
	case *Sequence:
		s := &Sequence{}
		for _, v := range c.Sequence {
			s.Sequence = append(s.Sequence, FromPreserves(v))
		}
		return s
	case *preserves.Set:
		s := &Set{Set: make(preserves.Set)}
		for k := range *c {
			s.Set[FromPreserves(k)] = struct{}{}
		}
		return s
	case *preserves.Record:
		r := &Record{}
		r.Key = FromPreserves(c.Key)
		for _, v := range c.Fields {
			r.Fields = append(r.Fields, FromPreserves(v))
		}
		return r
	case *preserves.Annotation:
		return &Annotation{
			Annotation: preserves.Annotation{
				Value:          FromPreserves(c.Value),
				AnnotatedValue: FromPreserves(c.AnnotatedValue),
			},
		}
	case *preserves.Embedded:
		e := &Embedded{}
		e.Value = FromPreserves(c.Value)
		return e
	case nil:
		return c
	case *Record:
		r := &Record{}
		r.Key = FromPreserves(c.Key)
		for _, v := range c.Fields {
			r.Fields = append(r.Fields, FromPreserves(v))
		}
		return r
	case *Annotation:
		return &Annotation{
			Annotation: preserves.Annotation{
				Value:          FromPreserves(c.Value),
				AnnotatedValue: FromPreserves(c.AnnotatedValue),
			},
		}
	case *Dictionary:
		d := &Dictionary{Dictionary: make(preserves.Dictionary)}
		for k, v := range c.Dictionary {
			d.Dictionary[FromPreserves(k)] = FromPreserves(v)
		}
		return d
	}
	log.Fatalf("unknown type(3): %v", reflect.TypeOf(v))
	return nil
}
