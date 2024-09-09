package lib

//Implements [[Preserves Schema](https://preserves.dev/preserves-schema.html) for Go-lang.

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	BooleanFalseRune = 0x80
	BooleanTrueRune  = 0x81
	Reserved82Rune   = 0x82
	Reserved83Rune   = 0x83
	EndRune          = 0x84
	AnnotationRune   = 0x85
	EmbeddedRune     = 0x86
	DoubleRune       = 0x87
	// up to 0xAF
	SignedIntegerRune = 0xB0
	StringRune        = 0xB1
	ByteStringRune    = 0xB2
	SymbolRune        = 0xB3
	RecordRune        = 0xB4
	SequenceRune      = 0xB5
	SetRune           = 0xB6
	DictionaryRune    = 0xB7
	ReservedB8Rune    = 0xB8
	ReservedB9Rune    = 0xB9
	ReservedBARune    = 0xBA
	ReservedBBRune    = 0xBB
	ReservedBCRune    = 0xBC
	ReservedBDRune    = 0xBD
	ReservedBERune    = 0xBE
	ReservedBFRune    = 0xBF
)

type MoreToRead struct {
	Read int
	Left int
}

func MaybeMoreToRead(l int, n int) (err error) {
	if l > 0 {
		err = &MoreToRead{
			Read: n,
			Left: l,
		}
	}
	return
}

func (m *MoreToRead) Error() string {
	return fmt.Sprintf("read %d but there's %d left", m.Read, m.Left)
}

func ReadRune(r Rune, rd io.Reader) (n int, err error) {
	b := make([]byte, 1)
	n, err = rd.Read(b)
	if err != nil {
		return
	}
	if b[0] != byte(r.BinaryRune()) {
		err = fmt.Errorf("parse error: missing rune: %x != %x", b[0], byte(r.BinaryRune()))
	}
	return
}
func ReadStringRune(r Rune, p Position, rd io.Reader) (n int, err error) {
	pr, ok := rd.(*PeekReader)
	if !ok {
		pr = NewPeekReader(rd)
	}
	s := r.TextRune(p)
	b := make([]byte, len(s))
	n, err = pr.Read(b)
	if err == io.EOF {
		err = nil
		return
	}
	if err != nil {
		err = fmt.Errorf("parse error: read: %v", err)
		return
	}
	if string(b) != s {
		err = fmt.Errorf("parse error: missing string rune %s != %s", s, string(b))
	}
	return
}

func ReadUntil(r *PeekReader, del []byte, endRune func(Position) string) (n int, err error) {
	endByte := byte(0x00)
	if endRune != nil {
		end := endRune(END)
		if len(end) > 0 {
			endByte = end[0]
		}
	}
	var bs []byte
	for {
	again:
		bs, err = r.Peek(1)
		if err != nil {
			return
		}

		if bs[0] == endByte {
			_, err = r.ReadByte()
			n += 1
			if err != nil {
				return
			}
			err = new(EndRuneError)
			return
		}

		for _, d := range del {
			if bs[0] == d {
				_, err = r.ReadByte()
				n += 1
				if err != nil {
					return
				}
				goto again
			}
		}

		return
	}

}

type EndRuneError string

func (e *EndRuneError) Error() string {
	return string("found endrune")
}

type Position int

const (
	START Position = iota
	END
	TRUE
	FALSE
)
const (
	COLON = ':'
	COMMA = ','
	WS    = ' '
)

var SYMBOL_OR_NUMBER_REGEXP = regexp.MustCompile(`^[-a-zA-Z0-9~!$%^&*?_=+/.]+$`)
var NUMBER_REGEXP = regexp.MustCompile(`^([-+]?\d+)((\.\d+([eE][-+]?\d+)?)|([eE][-+]?\d+))?$`)
var DOUBLE_REGEXP = regexp.MustCompile(`^([-+]?\d+)((\.\d+([eE][-+]?\d+)?)|([eE][-+]?\d+))$`)
var SIGNED_INTEGER_REGEXP = regexp.MustCompile(`^([-+]?\d+)$`)

func ReadValueFromBinary(r io.Reader) (v Value, n int, err error) {
	pr := NewPeekReader(r)
	var b []byte
	b, err = pr.Peek(1)
	n = len(b)
	if err != nil {
		return
	}
	if b[0] == EndRune {
		err = new(EndRuneError)
		return
	}
	if v, err = NewValueFromBinary(b[0]); err != nil {
		return
	}
	var m int
	m, err = v.UnmashalBinaryStream(pr)
	n += m
	return
}

func ReadValueFromText(r io.Reader, endRune func(Position) string) (v Value, n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}

	var bs []byte
	if bs, err = pr.Peek(1); err != nil {
		return
	}
	for bs[0] == ' ' || bs[0] == '\t' || bs[0] == '\n' || bs[0] == '\r' || bs[0] == ',' {
		_, err = pr.ReadByte()
		if err != nil {
			return
		}
		if bs, err = pr.Peek(1); err != nil {
			return
		}
	}

	if endRune != nil {
		if bs[0] == endRune(END)[0] {
			_, err = pr.ReadByte()
			if err != nil {
				return
			}
			err = new(EndRuneError)
			return
		}
	}

	v, err = NewValueFromTextPeek(pr)
	if err != nil {
		return
	}
	var m int
	m, err = v.UnmarshalTextStream(pr)
	n += m
	return
}

type PeekReader struct {
	io.Reader
	prepend  []byte
	Position int
	ReadData []byte
}

func NewPeekReader(r io.Reader) *PeekReader {
	return &PeekReader{Reader: r}
}

func (p *PeekReader) Peek(l int) (b []byte, err error) {
	if len(p.prepend) > 0 {
		for i, c := range p.prepend {
			if i >= l {
				return
			}
			b = append(b, c)
		}
		c := make([]byte, l-len(p.prepend))
		_, err = p.Reader.Read(c)
		if err != nil {
			return
		}
		b = append(b, c...)
		p.prepend = append(p.prepend, c...)
		p.ReadData = append(p.ReadData, c...)
		return
	}
	b = make([]byte, l)
	var n int
	n, err = p.Reader.Read(b)
	p.prepend = append(p.prepend, b...)
	p.Position += n
	p.ReadData = append(p.ReadData, b...)
	return
}
func (p *PeekReader) Read(b []byte) (n int, err error) {
	defer func() {
		p.Position += n
		p.ReadData = append(p.ReadData, b[:n]...)
	}()
	if len(b) == 0 {
		return
	}
	if len(p.prepend) > 0 {
		if len(b) == len(p.prepend) {
			copy(b, p.prepend)
			p.prepend = []byte{}
			return
		} else if len(b) < len(p.prepend) {
			copy(b, p.prepend[:len(b)])
			p.prepend = p.prepend[len(b):]
			return
		}
		copy(b, p.prepend)
		defer func() {
			p.prepend = []byte{}
		}()
		a := make([]byte, len(b)-len(p.prepend))
		n, err = p.Reader.Read(a)
		if err != nil {
			return
		}
		copy(b[len(p.prepend):], a)
		return
	}
	return p.Reader.Read(b)
}
func (p *PeekReader) ReadByte() (d byte, err error) {
	if len(p.prepend) > 0 {
		d = p.prepend[0]
		p.prepend = p.prepend[1:]
		p.Position += 1
		p.ReadData = append(p.ReadData, d)
		return
	}
	br, ok := p.Reader.(io.ByteReader)
	if !ok {
		b := make([]byte, 1)
		_, err = p.Reader.Read(b)
		d = b[0]
		p.Position += 1
		p.ReadData = append(p.ReadData, d)
		return
	}
	p.Position += 1
	p.ReadData = append(p.ReadData, d)
	return br.ReadByte()
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
		err = fmt.Errorf("reader not an bytereader")
		return
	}
	bl := &ByteLenReader{ByteReader: br}
	l, err = binary.ReadUvarint(bl)
	n = bl.readBytes
	return
}

type BinaryMarshalStream interface {
	MarshalBinaryStream(io.Writer) (int, error)
}
type BinaryUnmarshalStream interface {
	UnmashalBinaryStream(io.Reader) (int, error)
}
type BinaryUnMarshalerStream interface {
	BinaryMarshalStream
	BinaryUnmarshalStream
}
type TextMarshalStream interface {
	MarshalTextStream(io.Writer) (n int, err error)
}
type TextUnmarshalerStream interface {
	UnmarshalTextStream(r io.Reader) (n int, err error)
}
type TextUnMarshalerStream interface {
	TextMarshalStream
	TextUnmarshalerStream
}

type Rune interface {
	BinaryRune() rune
	TextRune(Position) string
}

type Value interface {
	Equal(Value) bool
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	BinaryUnMarshalerStream
	TextUnMarshalerStream
	Rune
}

func NewValueFromBinary(b byte) (Value, error) {
	switch b {
	case BooleanFalseRune:
		return new(Boolean), nil
	case BooleanTrueRune:
		return new(Boolean), nil
	case DoubleRune:
		return new(Double), nil
	case SignedIntegerRune:
		return new(SignedInteger), nil
	case StringRune:
		return new(String), nil
	case ByteStringRune:
		return new(ByteString), nil
	case SymbolRune:
		return new(Symbol), nil
	case RecordRune:
		return new(Record), nil
	case SequenceRune:
		return new(Sequence), nil
	case SetRune:
		s := NewSet([]Value{})
		return Value(&s), nil
	case DictionaryRune:
		d := NewDictionary(map[Value]Value{})
		return Value(&d), nil
	case AnnotationRune:
		return new(Annotation), nil
	case EmbeddedRune:
		e := NewEmbedded(nil)
		return Value(&e), nil
	}
	return nil, fmt.Errorf("unknown byte: %x", b)
}

func NewValueFromTextPeek(r *PeekReader) (v Value, err error) {
	var bs []byte
	var str []byte
	bs, err = r.Peek(1)
	if err != nil {
		return
	}
	switch bs[0] {
	case '"':
		v = new(String)
		return
	case '<':
		v = new(Record)
		return
	case '[':
		v = new(Sequence)
		return
	case '{':
		d := NewDictionary(map[Value]Value{})
		v = Value(&d)
		return
	case '@':
		v = new(Annotation)
		return
	case '|':
		v = new(Symbol)
		return
	case '#':
		bs, err = r.Peek(2)
		if err != nil {
			return
		}
		switch bs[1] {
		case 't':
			b := Boolean(true)
			v = Value(&b)
			return
		case 'f':
			b := Boolean(false)
			v = Value(&b)
			return
		case '"':
			v = new(ByteString)
			return
		case 'x':
			bs, err = r.Peek(3)
			if err != nil {
				return
			}
			if bs[2] == '"' {
				v = new(ByteString)
				return
			}
			if bs[2] == 'd' {
				bs, err = r.Peek(4)
				if err != nil {
					return
				}
				if bs[3] == '"' {
					v = new(Double)
					return
				}
			}
		case '[':
			v = new(ByteString)
			return
		case '{':
			s := NewSet([]Value{})
			v = Value(&s)
			return
		case ':':
			e := NewEmbedded(nil)
			v = Value(&e)
			return
		}
		v = new(Comment)
		return
	}
	// it's a bare keyword
	i := 1
	for {
		bs, err = r.Peek(i)
		if err != nil {
			return
		}
		if !SYMBOL_OR_NUMBER_REGEXP.Match(bs[i-1:]) {
			if SIGNED_INTEGER_REGEXP.Match(str) {
				v = new(SignedInteger)
				return
			}
			if DOUBLE_REGEXP.Match(str) {
				v = new(Double)
				return
			}
			v = new(Symbol)
			return
		}
		str = append(str, bs[i-1])
		i += 1
	}
}

type BareError struct{}

func (b *BareError) Error() string {
	return "bareError"
}

//type Atom interface {
//	Boolean | Double | SignedInteger | String | ByteString | Symbol
//}

type Boolean bool

func NewBoolean(b bool) Boolean {
	return Boolean(b)
}
func (b *Boolean) Equal(y Value) bool {
	x, ok := y.(*Boolean)
	if !ok {
		return false
	}
	return bool(*b) == bool(*x)
}

func (b *Boolean) BinaryRune() rune {
	return BooleanFalseRune
}
func (b *Boolean) TextRune(p Position) string {
	if p == TRUE {
		return "#t"
	}
	return "#f"
}
func (b *Boolean) MarshalBinaryStream(w io.Writer) (n int, err error) {
	if *b {
		return w.Write([]byte{BooleanTrueRune})
	}
	return w.Write([]byte{byte(b.BinaryRune()), 0x08})
}
func (b *Boolean) UnmashalBinaryStream(r io.Reader) (n int, err error) {
	bs := make([]byte, 1)
	n, err = r.Read(bs)
	if err != nil {
		return
	}
	if bs[0] == BooleanFalseRune {
		*b = false
		return
	}
	if bs[0] == BooleanTrueRune {
		*b = true
	}

	return
}
func (b *Boolean) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	bs := make([]byte, 2)
	n, err = pr.Read(bs)
	if err != nil {
		return
	}
	if string(bs) == "#t" {
		*b = Boolean(true)
		return
	} else if string(bs) == "#f" {
		*b = Boolean(false)
		return
	}
	err = fmt.Errorf("parse error: invalid rune: %s", string(bs))
	return
}
func (b *Boolean) MarshalTextStream(w io.Writer) (n int, err error) {
	if *b {
		return w.Write([]byte("#t"))
	}
	return w.Write([]byte("#f"))
}
func (b *Boolean) MarshalBinary() (data []byte, err error) {
	if *b {
		data = []byte{BooleanTrueRune}
		return
	}
	data = []byte{BooleanFalseRune}
	return
}
func (b *Boolean) UnmarshalBinary(data []byte) (err error) {
	br := bytes.NewReader(data)
	var n int
	n, err = b.UnmashalBinaryStream(br)
	if err != nil {
		return
	}
	data = data[n:]
	err = MaybeMoreToRead(len(data), n)
	return
}
func (b *Boolean) MarshalText() (data []byte, err error) {
	buf := &bytes.Buffer{}
	_, err = b.MarshalTextStream(buf)
	return buf.Bytes(), err
}
func (b *Boolean) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = b.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type Double float64

func NewDouble(f float64) Double {
	return Double(f)
}
func (d *Double) Equal(y Value) bool {
	x, ok := y.(*Double)
	if !ok {
		return false
	}
	return float64(*d) == float64(*x)
}

func (d *Double) BinaryRune() rune {
	return DoubleRune
}
func (d *Double) TextRune(_ Position) string {
	return ""
}
func (d *Double) MarshalBinaryStream(w io.Writer) (n int, err error) {
	// Length is always 64bit (0x08)
	n, err = w.Write([]byte{byte(d.BinaryRune()), 0x08})
	if err != nil {
		return
	}
	var m int

	n, err = w.Write(binary.BigEndian.AppendUint64([]byte{}, math.Float64bits(float64(*d))))
	n += m
	return
}

func (d *Double) UnmashalBinaryStream(r io.Reader) (n int, err error) {
	n, err = ReadRune(d, r)
	if err != nil {
		return
	}
	var m int
	b := make([]byte, 1)
	m, err = r.Read(b)
	n += m
	if err != nil {
		return
	}
	if b[0] != 0x08 {
		err = fmt.Errorf("could not read length")
		return
	}
	b = make([]byte, 8)
	m, err = r.Read(b)
	n += m
	if err != nil {
		return
	}
	*d = Double(math.Float64frombits(binary.BigEndian.Uint64(b)))
	return
}
func (d *Double) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	var bs []byte

	bs, err = pr.Peek(1)
	if err != nil {
		return
	}
	if bs[0] == '#' {
		bs, err = pr.Peek(4)
		if err != nil {
			return
		}
		if string(bs) == "#xd\"" {
			bs = make([]byte, 4)
			n, err = pr.Read(bs)
			if err != nil {
				return
			}
			var m int
			f := make([]byte, 8)
			for i := range []int{0, 1, 2, 3, 4, 5, 6, 7} {
				m, err = ReadUntil(pr, []byte{' ', '\t', '\n', '\r'}, nil)
				n += m
				if err != nil {
					return
				}
				bs = make([]byte, 2)
				m, err = pr.Read(bs)
				n += m
				if err != nil {
					return
				}
				m, err = fmt.Sscanf(string(bs), "%x", &f[i])
				if err != nil {
					return
				}
				if m != 1 {
					err = fmt.Errorf("parse error: no matches")
					return
				}
			}
			m, err = ReadUntil(pr, []byte{' ', '\t', '\n', '\r'}, nil)
			n += m
			if err != nil {
				return
			}
			bs = make([]byte, 1)
			m, err = pr.Read(bs)
			n += m
			if err != nil {
				return
			}
			if bs[0] != '"' {
				err = fmt.Errorf("parse error: missing end dquote %s != %s", string(bs), "\"")
				return
			}
			*d = Double(math.Float64frombits(binary.BigEndian.Uint64(f)))
			return
		}
	}
	var str []byte
	for {
		bs, err = pr.Peek(1)

		if err != nil {
			return
		}
		if !SYMBOL_OR_NUMBER_REGEXP.Match(bs) {
			var a float64
			a, err = strconv.ParseFloat(string(str), 64)
			if err != nil {
				return
			}
			*d = Double(a)
			return
		}
		_, err = pr.ReadByte()
		n += 1
		if err != nil {
			return
		}
		str = append(str, bs[0])
	}
}
func (d *Double) MarshalTextStream(w io.Writer) (n int, err error) {
	return fmt.Fprintf(w, "%f", *d)
}
func (d *Double) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = d.MarshalBinaryStream(b)
	return b.Bytes(), err

}
func (d *Double) UnmarshalBinary(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	var n int
	n, err = d.UnmashalBinaryStream(b)
	if err != nil {
		return
	}
	data = data[n:]
	return MaybeMoreToRead(len(data), n)
}
func (d *Double) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = d.MarshalTextStream(b)
	return b.Bytes(), err
}
func (d *Double) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = d.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type SignedInteger big.Int

func NewSignedInteger(i int64) SignedInteger {
	return SignedInteger(*big.NewInt(i))
}
func (s *SignedInteger) Equal(y Value) bool {
	x, ok := y.(*SignedInteger)
	if !ok {
		return false
	}
	a := big.Int(*s)
	b := big.Int(*x)
	return a.Cmp(&b) == 0
}

func (s *SignedInteger) BinaryRune() rune {
	return SignedIntegerRune
}
func (s *SignedInteger) TextRune(_ Position) string {
	return ""
}
func (s *SignedInteger) MarshalTextStream(w io.Writer) (n int, err error) {
	a := big.Int(*s)
	var str []byte
	str, err = a.MarshalText()
	n = len(str)
	if err != nil {
		return
	}
	var m int
	m, err = w.Write(str)
	n += m
	return
}
func (s *SignedInteger) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	var a big.Int
	var bs []byte
	var str []byte
	for {
		bs, err = pr.Peek(1)
		if err != nil {
			return
		}
		if !SYMBOL_OR_NUMBER_REGEXP.Match(bs) {
			a = big.Int(*s)
			err = a.UnmarshalText(str)
			if err != nil {
				return
			}
			*s = SignedInteger(a)
			return
		}
		_, err = pr.ReadByte()
		n += 1
		if err != nil {
			return
		}
		str = append(str, bs[0])
	}
}
func (s *SignedInteger) MarshalBinaryStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte{byte(s.BinaryRune())})
	if err != nil {
		return
	}
	var m int

	a := big.Int(*s)
	l := len(a.Bytes())
	b := a.Bytes()
	switch a.Cmp(big.NewInt(0)) {
	case 1:
		if b[0]&0x80 == 0x80 {
			l += 1
			b = append([]byte{0x00}, b...)
		}
	case -1:
		slices.Reverse(b)
		carry := 0x00
		for i, v := range b {
			b[i] = (v ^ 0xFF)
			if i == 0 || carry == 0x01 {
				if b[i] == 0xFF {
					if i == 0 {
						carry = 0x01
					}
					b[i] = 0x00
				} else {
					b[i] += 0x01
					carry = 0x00
				}
			}
		}
		slices.Reverse(b)
		if b[0]&0x80 == 0x00 {
			l += 1
			b = append([]byte{0xFF}, b...)
		}
	case 0:
		m, err = w.Write([]byte{0x00})
		n += m
		return
	}
	m, err = w.Write(binary.AppendUvarint([]byte{}, uint64(l)))
	n += m
	if err != nil {
		return
	}
	m, err = w.Write(b)
	n += m
	return
}
func (s *SignedInteger) UnmashalBinaryStream(r io.Reader) (n int, err error) {
	n, err = ReadRune(s, r)
	if err != nil {
		return
	}
	var (
		l uint64
		m int
	)
	l, m, err = ReadUvarint(r)
	n += m
	if err != nil {
		return
	}
	a := big.Int(*s)
	if l == 0 {
		a.SetUint64(0)
		*s = SignedInteger(a)
		return
	}
	b := make([]byte, l)
	err = binary.Read(r, binary.BigEndian, b)
	n += int(l)
	if err != nil {
		return
	}

	if b[0]&0x80 == 0x80 {
		slices.Reverse(b)
		carry := 0x00
		for i, v := range b {
			b[i] = (v ^ 0xFF)
			if i == 0 || carry == 0x01 {
				if b[i] == 0xFF {
					if i == 0 {
						carry = 0x01
					}
					b[i] = 0x00
				} else {
					b[i] += 0x01
					carry = 0x00
				}
			}
		}

		slices.Reverse(b)
		a.Neg(a.SetBytes(b))
	} else {
		a.SetBytes(b)
	}
	*s = SignedInteger(a)
	return
}
func (s *SignedInteger) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.MarshalBinaryStream(b)
	return b.Bytes(), err
}
func (s *SignedInteger) UnmarshalBinary(data []byte) (err error) {
	r := bytes.NewBuffer(data)
	var n int
	n, err = s.UnmashalBinaryStream(r)
	if err != nil {
		return
	}
	data = data[n:]

	return MaybeMoreToRead(len(data), n)
}
func (s *SignedInteger) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.MarshalTextStream(b)
	return b.Bytes(), err
}
func (s *SignedInteger) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = s.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type String string

func NewString(s string) String {
	return String(s)
}

func (s *String) Equal(y Value) bool {
	x, ok := y.(*String)
	if !ok {
		return false
	}
	return string(*s) == string(*x)
}
func (s *String) BinaryRune() rune {
	return StringRune
}
func (s *String) TextRune(_ Position) string {
	return "\""
}
func (s *String) MarshalBinaryStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte{byte(s.BinaryRune())})
	if err != nil {
		return
	}
	var m int
	b := binary.AppendUvarint([]byte{}, uint64(len(*s)))
	m, err = w.Write(b)
	n += m
	if err != nil {
		return
	}
	m, err = w.Write([]byte(*s))
	n += m
	return
}
func (s *String) UnmashalBinaryStream(r io.Reader) (n int, err error) {
	n, err = ReadRune(s, r)
	if err != nil {
		return
	}
	var (
		l uint64
		m int
	)
	l, m, err = ReadUvarint(r)
	n += m
	if err != nil {
		return
	}

	str := make([]byte, l)
	m, err = r.Read(str)
	n += m
	if err != nil {
		return
	}
	*s = String(str)
	return
}
func (s *String) MarshalTextStream(w io.Writer) (n int, err error) {
	return fmt.Fprintf(w, "%s%s%s", s.TextRune(START), *s, s.TextRune(END))
}
func (s *String) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	var b byte
	b, err = pr.ReadByte()
	n += 1
	if err != nil {
		return
	}
	if b != '"' {
		err = fmt.Errorf("parse error: no start dquote found")
		return
	}
	var bs []byte
	var escaped bool
	for {
		b, err = pr.ReadByte()
		n += 1
		if err != nil {
			return
		}
		if escaped {
			escaped = false
			bs = append(bs, b)
			continue
		}
		if b == '\\' {
			escaped = true
		}

		if b == '"' {
			*s = String(bs)
			return
		}
		bs = append(bs, b)
	}
}
func (s *String) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.MarshalBinaryStream(b)
	return b.Bytes(), err

}
func (s *String) UnmarshalBinary(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	var n int
	n, err = s.UnmashalBinaryStream(buf)
	if err != nil {
		return
	}
	data = data[n:]
	return MaybeMoreToRead(len(data), n)
}
func (s *String) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.MarshalTextStream(b)
	return b.Bytes(), err
}
func (s *String) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = s.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type ByteString []byte

func NewByteString(b []byte) ByteString {
	return ByteString(b)
}
func (bs *ByteString) Equal(v Value) bool {
	x, ok := v.(*ByteString)
	if !ok {
		return false
	}
	return slices.Equal([]byte(*bs), []byte(*x))
}
func (bs *ByteString) BinaryRune() rune {
	return ByteStringRune
}
func (bs *ByteString) TextRune(d Position) string {
	if d == START {
		return "#\""
	}
	return "\""
}
func (bs *ByteString) MarshalTextStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte(bs.TextRune(START)))
	if err != nil {
		return
	}
	var m int

	m, err = w.Write([]byte{' '})
	n += m
	if err != nil {
		return
	}
	bw := base64.NewEncoder(base64.RawStdEncoding, w)
	defer bw.Close()
	for _, v := range *bs {
		m, err = bw.Write([]byte{v})
		n += m
		if err != nil {
			return
		}

		m, err = w.Write([]byte{' '})
		n += m
		if err != nil {
			return
		}
	}
	m, err = w.Write([]byte(bs.TextRune(END)))
	n += m
	return
}
func (bs *ByteString) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	var bsb []byte
	bsb, err = pr.Peek(3)
	if err != nil {
		return
	}
	if bsb[0] != '#' {
		err = fmt.Errorf("parse error: missing rune %s != %s", "#", string(bsb[0]))
		return
	}
	_, err = pr.ReadByte()
	n += 1
	if err != nil {
		return
	}
	var b byte
	switch bsb[1] {
	case '[':
		_, err = pr.ReadByte()
		n += 1
		if err != nil {
			return
		}
		var bytes []byte
		for {
			b, err = pr.ReadByte()
			n += 1
			if err != nil {
				return
			}
			if b == ']' {
				bytes, err = base64.RawStdEncoding.DecodeString(string(bytes))
				if err != nil {
					return
				}
				*bs = ByteString(bytes)
				return
			}
			if b == ' ' || b == '\t' || b == '=' {
				continue
			}
			bytes = append(bytes, b)
		}
	case '"':
		_, err = pr.ReadByte()
		n += 1
		if err != nil {
			return
		}
		var bytes []byte
		var escaped bool
		var hex bool
		for {
			b, err = pr.ReadByte()
			n += 1
			if err != nil {
				return
			}
			if hex {
				hex = false
				b1 := b
				b, err = pr.ReadByte()
				n += 1
				if err != nil {
					return
				}
				var i uint64
				i, err = strconv.ParseUint(string([]byte{b, b1}), 16, 8)
				if err != nil {
					return
				}
				bytes = append(bytes, byte(i))
			}
			if escaped {
				escaped = false
				if b == 'n' {
					bytes = append(bytes, '\n')
					continue
				}
				if b == 'r' {
					bytes = append(bytes, '\r')
					continue
				}
				if b == 't' {
					bytes = append(bytes, '\t')
					continue
				}
				if b == 'b' {
					bytes = append(bytes, '\b')
					continue
				}
				if b == 'f' {
					bytes = append(bytes, '\f')
					continue
				}
				if b == 'x' {
					hex = true
					continue
				}
				if b == '"' {
					bytes = append(bytes, b)
					continue
				}
				if b == '\\' {
					bytes = append(bytes, b)
					continue
				}
			}
			if b == '\\' {
				escaped = true
			}
			if string(b) == bs.TextRune(END) {
				*bs = ByteString(bytes)
				return
			}
			bytes = append(bytes, b)
		}
	case 'x':
		_, err = pr.ReadByte()
		n += 1
		if err != nil {
			return
		}
		b, err = pr.ReadByte()
		n += 1
		if err != nil {
			return
		}
		if b != '"' {
			err = fmt.Errorf("parse error: missing rune %s != %s", "\"", string(b))
			return
		}
		var bytes []byte
		for {
			b, err = pr.ReadByte()
			n += 1
			if err != nil {
				return
			}
			if string(b) == bs.TextRune(END) {
				*bs = ByteString(bytes)
				return
			}
			bytes = append(bytes, b)
		}
	}
	err = fmt.Errorf("parse error: missing rune: %s", bsb)
	return
}

func (bs *ByteString) MarshalBinaryStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte{byte(bs.BinaryRune())})
	if err != nil {
		return
	}
	var m int
	b := binary.AppendUvarint([]byte{}, uint64(len(*bs)))
	m, err = w.Write(b)
	n += m
	if err != nil {
		return
	}
	m, err = w.Write([]byte(*bs))
	n += m
	return
}
func (bs *ByteString) UnmashalBinaryStream(r io.Reader) (n int, err error) {
	n, err = ReadRune(bs, r)
	if err != nil {
		return
	}
	var (
		l uint64
		m int
	)
	l, m, err = ReadUvarint(r)
	n += m
	if err != nil {
		return
	}

	b := make([]byte, l)
	m, err = r.Read(b)
	n += m
	if err != nil {
		return
	}

	*bs = ByteString(b)
	return
}
func (bs *ByteString) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = bs.MarshalBinaryStream(b)
	return b.Bytes(), err

}
func (bs *ByteString) UnmarshalBinary(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	var n int
	n, err = bs.UnmashalBinaryStream(buf)
	if err != nil {
		return
	}
	data = data[n:]
	return MaybeMoreToRead(len(data), n)
}
func (bs *ByteString) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = bs.MarshalTextStream(b)
	return b.Bytes(), err
}
func (bs *ByteString) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = bs.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type ByteStringSecondForm ByteString

func (bs *ByteStringSecondForm) TextRune(d Position) string {
	if d == START {
		return "#x\""
	}
	return "\""
}

type ByteStringThirdForm ByteString

func (bs *ByteStringThirdForm) TextRune(d Position) string {
	if d == START {
		return "#["
	}
	return "]"
}
func (bs *ByteStringThirdForm) MarshalTextStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte(bs.TextRune(START)))
	if err != nil {
		return
	}
	var m int

	m, err = w.Write([]byte{' '})
	n += m
	if err != nil {
		return
	}
	bw := base64.NewEncoder(base64.RawStdEncoding, w)
	defer bw.Close()
	for _, v := range *bs {
		m, err = bw.Write([]byte{v})
		n += m
		if err != nil {
			return
		}

		m, err = w.Write([]byte{' '})
		n += m
		if err != nil {
			return
		}
	}
	m, err = w.Write([]byte(bs.TextRune(END)))
	n += m
	return
}

type Symbol string

var (
	RAW_SYMBOL_RE = regexp.MustCompile("^[-a-zA-Z0-9~!$%^&*?_=+/.]+$")
)

func NewSymbol(s any) Symbol {
	switch a := s.(type) {
	case string:
		return Symbol(a)
	case Symbol:
		return Symbol(a.String())
	default:
		panic(fmt.Errorf("newSymbol: only string | Symbol is supported"))
	}
}

func (s *Symbol) Equal(y Value) bool {
	x, ok := y.(*Symbol)
	if !ok {
		return false
	}
	return string(*s) == string(*x)
}

func (s *Symbol) String() string {
	return string(*s)
}

func (s *Symbol) Representation() string {
	return fmt.Sprintf("#%s", s)
}
func (s *Symbol) BinaryRune() rune {
	return SymbolRune
}
func (s *Symbol) TextRune(_ Position) string {
	return "|"
}
func (s *Symbol) MarshalTextStream(w io.Writer) (n int, err error) {
	if *s != "" {
		return fmt.Fprintf(w, "%s%s%s", s.TextRune(START), strings.ReplaceAll(string(*s), "|", "\\|"), s.TextRune(END))
	}
	return
}
func (s *Symbol) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	var bs []byte
	bs, err = pr.Peek(1)
	if err != nil {
		return
	}

	var b byte
	var str []byte
	switch bs[0] {
	case s.TextRune(START)[0]:
		// skip the quote
		_, err = pr.ReadByte()
		n += 1
		if err != nil {
			return
		}
		var escaped bool
		for {
			b, err = pr.ReadByte()
			n += 1
			if err != nil {
				return
			}

			if escaped {
				escaped = false
				if b == '|' {
					str = append(str, b)
					continue
				}

			}
			if b == '\\' {
				escaped = true
				continue
			}

			if b == '|' {
				*s = Symbol(string(str))
				return
			}
			str = append(str, b)
		}
	default:
		for {
			if !SYMBOL_OR_NUMBER_REGEXP.Match(bs) {
				n = len(str)
				*s = Symbol(string(str))
				return
			}

			_, err = pr.ReadByte()
			n += 1
			if err != nil && err != io.EOF {
				return
			}
			str = append(str, bs[0])
			bs, err = pr.Peek(1)
			if err != nil {
				return
			}
		}
	}
}
func (s *Symbol) MarshalBinaryStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte{byte(s.BinaryRune())})
	if err != nil {
		return
	}
	var m int
	b := binary.AppendUvarint([]byte{}, uint64(len(string(*s))))
	m, err = w.Write(b)
	n += m
	if err != nil {
		return
	}
	m, err = w.Write([]byte(*s))
	n += m
	return
}
func (s *Symbol) UnmashalBinaryStream(r io.Reader) (n int, err error) {
	n, err = ReadRune(s, r)
	if err != nil {
		return
	}
	var (
		l uint64
		m int
	)
	l, m, err = ReadUvarint(r)
	n += m
	if err != nil {
		return
	}
	str := make([]byte, l)
	m, err = r.Read(str)
	n += m
	if err != nil {
		return
	}
	*s = Symbol(str)
	return
}
func (s *Symbol) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.MarshalBinaryStream(b)
	return b.Bytes(), err
}
func (s *Symbol) UnmarshalBinary(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	var n int
	n, err = s.UnmashalBinaryStream(buf)
	if err != nil {
		return
	}
	data = data[n:]
	return MaybeMoreToRead(len(data), n)
}

func (s *Symbol) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = s.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

func (s *Symbol) MarshalText() (text []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.MarshalTextStream(b)
	return b.Bytes(), err
}

//type Compound interface {
//	Record | Sequence | Set | Dictionary
//}

// Representation of Preserves `Record`s, which are a pair of a *label* `Value` and a sequence of *field* `Value`s.
type Record struct {
	Key    *Value
	Fields []*Value
}

func NewRecord(key Value, values []Value) Record {
	r := Record{Key: &key}
	for _, v := range values {
		r.Fields = append(r.Fields, &v)
	}
	return r
}
func (r *Record) Equal(v Value) bool {
	x, ok := v.(*Record)
	if !ok {
		return false
	}
	if !(*(r.Key)).Equal(*(x.Key)) {
		return false
	}
	ret := false
	for _, u := range r.Fields {
		for _, v := range x.Fields {
			if (*u).Equal(*v) {
				return true
			}
		}
	}
	return ret
}
func (r *Record) BinaryRune() rune {
	return RecordRune
}
func (r *Record) TextRune(d Position) string {
	if d == START {
		return "<"
	}
	return ">"
}
func (r *Record) MarshalTextStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte(r.TextRune(START)))
	if err != nil {
		return
	}
	var m int

	m, err = (*r.Key).MarshalTextStream(w)
	n += m
	if err != nil {
		return
	}
	m, err = w.Write([]byte{WS})
	n += m
	if err != nil {
		return
	}

	for _, v := range r.Fields {
		m, err = (*v).MarshalTextStream(w)
		n += m
		if err != nil {
			return
		}
		m, err = w.Write([]byte{WS})
		n += m
		if err != nil {
			return
		}
	}
	m, err = w.Write([]byte(r.TextRune(END)))
	n += m
	return
}
func (r *Record) UnmarshalTextStream(rd io.Reader) (n int, err error) {
	pr, ok := rd.(*PeekReader)
	if !ok {
		pr = NewPeekReader(rd)
	}
	n, err = ReadStringRune(r, START, pr)
	if err != nil {
		return
	}
	var m int

	var u Value
	u, m, err = ReadValueFromText(pr, r.TextRune)
	n += m
	if _, ok := err.(*EndRuneError); ok {
		err = nil
		return
	}
	if err != nil {
		return
	}
	if r.Key == nil {
		r.Key = &u
	} else {
		*r.Key = u
	}
	for {
		var v Value
		v, m, err = ReadValueFromText(pr, r.TextRune)
		n += m
		if _, ok := err.(*EndRuneError); ok {
			err = nil
			return
		}
		if err != nil {
			return
		}
		r.Fields = append(r.Fields, &v)
	}
}
func (r *Record) MarshalBinaryStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte{byte(r.BinaryRune())})
	if err != nil {
		return
	}
	var m int
	m, err = (*r.Key).MarshalBinaryStream(w)
	n += m
	if err != nil {
		return
	}

	for _, v := range r.Fields {
		m, err = (*v).MarshalBinaryStream(w)
		n += m
		if err != nil {
			return
		}
	}
	m, err = w.Write([]byte{EndRune})
	n += m
	return
}
func (r *Record) UnmashalBinaryStream(rd io.Reader) (n int, err error) {
	n, err = ReadRune(r, rd)
	if err != nil {
		return
	}
	var (
		m  int
		ok bool
		v  Value
	)
	v, m, err = ReadValueFromBinary(rd)
	n += m
	if _, ok = err.(*EndRuneError); ok {
		err = nil
		return
	}
	r.Key = &v
	r.Fields = []*Value{}
	for {
		var v Value
		v, m, err = ReadValueFromBinary(rd)
		n += m
		if _, ok = err.(*EndRuneError); ok {
			err = nil
			return
		}
		if err != nil {
			return
		}
		r.Fields = append(r.Fields, &v)
	}
}
func (r *Record) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = r.MarshalBinaryStream(b)
	return b.Bytes(), err
}

func (r *Record) UnmarshalBinary(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	var n int
	n, err = r.UnmashalBinaryStream(b)
	if err != nil {
		return
	}
	data = data[n:]
	return MaybeMoreToRead(len(data), n)
}
func (r *Record) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = r.MarshalTextStream(b)
	return b.Bytes(), err
}

func (r *Record) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = r.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type Sequence []Value

func NewSequence(v []Value) Sequence {
	return Sequence(v)
}
func (s *Sequence) Equal(y Value) bool {
	x, ok := y.(*Sequence)
	if !ok {
		return false
	}
	for _, u := range []Value(*s) {
		for _, v := range []Value(*x) {
			if !u.Equal(v) {
				return false
			}
		}
	}
	return true
}
func (s *Sequence) BinaryRune() rune {
	return SequenceRune
}
func (s *Sequence) TextRune(d Position) string {
	if d == START {
		return "["
	}
	return "]"
}
func (s *Sequence) MarshalTextStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte(s.TextRune(START)))
	if err != nil {
		return
	}
	var m int

	var firstWritten bool
	for _, v := range *s {
		if firstWritten {
			m, err = w.Write([]byte{COMMA})
			n += m
			if err != nil {
				return
			}
		}
		firstWritten = true
		m, err = v.MarshalTextStream(w)
		n += m
		if err != nil {
			return
		}
	}
	m, err = w.Write([]byte(s.TextRune(END)))
	n += m
	return
}
func (s *Sequence) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	n, err = ReadStringRune(s, START, pr)
	if err != nil {
		return
	}
	var m int
	for {
		var v Value
		v, m, err = ReadValueFromText(pr, s.TextRune)
		n += m
		if _, ok := err.(*EndRuneError); ok {
			err = nil
			return
		}
		if err != nil {
			return
		}
		*s = append(*s, v)
	}
}
func (s *Sequence) MarshalBinaryStream(w io.Writer) (n int, err error) {
	w.Write([]byte{byte(s.BinaryRune())})
	for _, v := range *s {
		if _, err = v.MarshalBinaryStream(w); err != nil {
			return
		}
	}
	return w.Write([]byte{EndRune})
}
func (s *Sequence) UnmashalBinaryStream(rd io.Reader) (n int, err error) {
	n, err = ReadRune(s, rd)
	if err != nil {
		return
	}
	var (
		m  int
		ok bool
		v  Value
	)
	for {
		v, m, err = ReadValueFromBinary(rd)
		n += m
		if _, ok = err.(*EndRuneError); ok {
			err = nil
			return
		}
		if err != nil {
			return
		}
		*s = append(*s, v)
	}
}
func (s *Sequence) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.MarshalBinaryStream(b)
	return b.Bytes(), err
}
func (s *Sequence) UnmarshalBinary(data []byte) (err error) {
	var n int
	n, err = s.UnmashalBinaryStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	data = data[n:]
	return MaybeMoreToRead(len(data), n)
}
func (s *Sequence) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.MarshalTextStream(b)
	return b.Bytes(), err
}
func (s *Sequence) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = s.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type Set map[Value]struct{}

func NewSet(v []Value) Set {
	s := Set{}
	for _, u := range v {
		s[u] = struct{}{}
	}
	return s
}
func (s *Set) Equal(y Value) bool {
	x, ok := y.(*Set)
	if !ok {
		return false
	}
	if len(*s) != len(*x) {
		return false
	}
	// TODO: Very inefficient way of doing this
	// Value should be represented by a hash instead
	for u := range *s {
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
func (s *Set) BinaryRune() rune {
	return SetRune
}
func (s *Set) TextRune(d Position) string {
	if d == START {
		return "#{"
	}
	return "}"
}
func (s *Set) MarshalTextStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte(s.TextRune(START)))
	if err != nil {
		return
	}
	var m int

	for v := range *s {
		m, err = v.MarshalTextStream(w)
		n += m
		if err != nil {
			return
		}
		m, err = w.Write([]byte{COMMA})
		n += m
		if err != nil {
			return
		}
	}
	m, err = w.Write([]byte(s.TextRune(END)))
	n += m
	return
}
func (s *Set) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	n, err = ReadStringRune(s, START, pr)
	if err != nil {
		return
	}
	var m int
	for {
		var v Value
		v, m, err = ReadValueFromText(pr, s.TextRune)
		n += m
		if _, ok := err.(*EndRuneError); ok {
			err = nil
			return
		}
		if err != nil {
			return
		}
		(*s)[v] = struct{}{}
	}
}
func (s *Set) MarshalBinaryStream(w io.Writer) (n int, err error) {
	w.Write([]byte{byte(s.BinaryRune())})
	for v := range *s {
		if _, err = v.MarshalBinaryStream(w); err != nil {
			return
		}
	}

	return w.Write([]byte{EndRune})
}
func (s *Set) UnmashalBinaryStream(rd io.Reader) (n int, err error) {
	n, err = ReadRune(s, rd)
	if err != nil {
		return
	}
	var (
		m  int
		ok bool
		v  Value
	)
	for {
		v, m, err = ReadValueFromBinary(rd)
		n += m
		if _, ok = err.(*EndRuneError); ok {
			err = nil
			return
		}
		if err != nil {
			return
		}
		(*s)[v] = struct{}{}
	}
}
func (s *Set) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.MarshalBinaryStream(b)
	return b.Bytes(), err
}
func (s *Set) UnmarshalBinary(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	var n int
	n, err = s.UnmashalBinaryStream(b)
	if err != nil {
		return
	}
	data = data[n:]
	return MaybeMoreToRead(len(data), n)
}
func (s *Set) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.MarshalTextStream(b)
	return b.Bytes(), err
}
func (s *Set) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = s.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type Dictionary map[Value]Value

func NewDictionary(v map[Value]Value) Dictionary {
	return Dictionary(v)
}

func (d *Dictionary) Equal(y Value) (b bool) {

	x, ok := y.(*Dictionary)
	if !ok {
		return false
	}
	if len(*d) != len(*x) {
		return false
	}

	// TODO: Inefficient way of doing it
	for u, t := range *d {
		equal := false
		for v, z := range *x {
			if u.Equal(v) {
				if t.Equal(z) {
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
func (d *Dictionary) BinaryRune() rune {
	return DictionaryRune
}
func (d *Dictionary) TextRune(p Position) string {
	if p == START {
		return "{"
	}
	return "}"
}
func (d *Dictionary) MarshalTextStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte(d.TextRune(START)))
	if err != nil {
		return
	}
	var m int

	for u, v := range *d {
		m, err = u.MarshalTextStream(w)
		n += m
		if err != nil {
			return
		}
		m, err = w.Write([]byte{COLON})
		n += m
		if err != nil {
			return
		}
		m, err = v.MarshalTextStream(w)
		n += m
		if err != nil {
			return
		}
		m, err = w.Write([]byte{COMMA})
		n += m
		if err != nil {
			return
		}
	}
	m, err = w.Write([]byte(d.TextRune(END)))
	n += m
	return
}
func (d *Dictionary) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	n, err = ReadStringRune(d, START, pr)
	if err != nil {
		return
	}
	var m int
	for {
		var u Value
		u, m, err = ReadValueFromText(pr, d.TextRune)
		n += m
		if _, ok := err.(*EndRuneError); ok {
			err = nil
			return
		}
		if err != nil {
			return
		}
		m, err = ReadUntil(pr, []byte{':'}, d.TextRune)
		n += m
		if err != nil {
			return
		}

		var v Value
		v, m, err = ReadValueFromText(pr, d.TextRune)
		n += m
		if _, ok := err.(*EndRuneError); ok {
			err = nil
			(*d)[u] = nil
			return
		}
		if err != nil {
			(*d)[u] = nil
			return
		}

		(*d)[u] = v
	}
}
func (d *Dictionary) MarshalBinaryStream(w io.Writer) (n int, err error) {
	w.Write([]byte{byte(d.BinaryRune())})
	for u, v := range *d {
		if _, err = u.MarshalBinaryStream(w); err != nil {
			return
		}
		if v != nil {
			if _, err = v.MarshalBinaryStream(w); err != nil {
				return
			}
		}
	}

	return w.Write([]byte{EndRune})
}
func (d *Dictionary) UnmashalBinaryStream(rd io.Reader) (n int, err error) {
	n, err = ReadRune(d, rd)
	if err != nil {
		return
	}
	var (
		m  int
		ok bool
		u  Value
		v  Value
	)
	for {
		u, m, err = ReadValueFromBinary(rd)
		n += m
		if _, ok = err.(*EndRuneError); ok {
			err = nil
			return
		}
		if err != nil {
			return
		}
		v, m, err = ReadValueFromBinary(rd)
		n += m
		if _, ok = err.(*EndRuneError); ok {
			err = nil
			return
		}
		if err != nil {
			return
		}
		(*d)[u] = v
	}
}
func (d *Dictionary) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = d.MarshalBinaryStream(b)
	return b.Bytes(), err
}
func (d *Dictionary) UnmarshalBinary(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	var n int
	n, err = d.UnmashalBinaryStream(b)
	if err != nil {
		return
	}
	data = data[n:]
	return MaybeMoreToRead(len(data), n)
}
func (d *Dictionary) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = d.MarshalTextStream(b)
	return b.Bytes(), err
}
func (d *Dictionary) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = d.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type Annotation struct {
	Value          *Value
	AnnotatedValue *Value
}

func NewAnnotation(v Value) Annotation {
	return Annotation{Value: &v}
}

func (a *Annotation) Equal(y Value) (b bool) {
	x, ok := y.(*Annotation)
	if !ok {
		return false
	}
	return a.Equal(x)
}
func (a *Annotation) BinaryRune() rune {
	return AnnotationRune
}
func (a *Annotation) TextRune(p Position) string {
	if p == START {
		return "@"
	}
	return ""
}
func (a *Annotation) MarshalTextStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte(a.TextRune(START)))
	if err != nil {
		return
	}
	var m int

	if a.Value == nil {
		return
	}
	m, err = (*a.Value).MarshalTextStream(w)
	n += m

	if a.AnnotatedValue == nil {
		return
	}
	m, err = (*a.AnnotatedValue).MarshalTextStream(w)
	n += m
	return
}
func (a *Annotation) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	n, err = ReadStringRune(a, START, pr)
	if err != nil {
		return
	}
	var (
		m int
		v Value
	)

	v, m, err = ReadValueFromText(pr, nil)
	n += m
	if err == nil || err == io.EOF {
		if a.Value == nil {
			a.Value = &v
		} else {
			(*a.Value) = v
		}
	}

	v, m, err = ReadValueFromText(pr, nil)
	n += m
	if err == nil || err == io.EOF {
		if a.AnnotatedValue == nil {
			a.AnnotatedValue = &v
		} else {
			(*a.AnnotatedValue) = v
		}
	}
	return
}
func (a *Annotation) MarshalBinaryStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte{byte(a.BinaryRune())})
	if err != nil {
		return
	}
	var m int
	m, err = Value(*a.Value).MarshalBinaryStream(w)
	n += m
	if err != nil {
		return
	}
	m, err = Value(*a.AnnotatedValue).MarshalBinaryStream(w)
	n += m
	return
}
func (a *Annotation) UnmashalBinaryStream(r io.Reader) (n int, err error) {
	if n, err = ReadRune(a, r); err != nil {
		return
	}

	var m int
	var u, v Value
	u, m, err = ReadValueFromBinary(r)
	n += m
	if err != nil {
		return
	}
	v, m, err = ReadValueFromBinary(r)
	n += m

	if err != nil {
		return
	}
	*a = Annotation{Value: &u, AnnotatedValue: &v}
	return
}
func (a *Annotation) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = a.MarshalBinaryStream(b)
	return b.Bytes(), err
}
func (a *Annotation) UnmarshalBinary(data []byte) (err error) {
	var n int
	n, err = a.UnmashalBinaryStream(bytes.NewBuffer(data))
	data = data[n:]
	if len(data) > 0 {
		err = &MoreToRead{
			Read: n,
			Left: len(data),
		}
	}
	return
}
func (a *Annotation) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = a.MarshalTextStream(b)
	return b.Bytes(), err
}
func (a *Annotation) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = a.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type Comment string

func (c *Comment) Equal(y Value) (b bool) {
	x, ok := y.(*Comment)
	if !ok {
		return false
	}
	return c.Equal(x)
}

// no binary representation
func (c *Comment) BinaryRune() rune {
	return '#'
}
func (c *Comment) TextRune(p Position) string {
	if p == START {
		return "#"
	}
	return "\n"
}
func (c *Comment) MarshalTextStream(w io.Writer) (n int, err error) {
	return fmt.Fprintf(w, "%s%s%s", c.TextRune(START), *c, c.TextRune(END))
}
func (c *Comment) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	var b byte
	var bs []byte
	for {
		b, err = pr.ReadByte()
		n += 1
		if err != nil {
			return
		}
		if b == '\n' {
			*c = Comment(string(bs))
			return
		}
		bs = append(bs, b)
	}
}
func (c *Comment) MarshalBinaryStream(w io.Writer) (n int, err error) {
	return
}
func (c *Comment) UnmashalBinaryStream(r io.Reader) (n int, err error) {
	return
}

// no binary representation
func (c *Comment) MarshalBinary() (data []byte, err error) {
	return
}

// no binary representation
func (c *Comment) UnmarshalBinary(data []byte) (err error) {
	return
}
func (c *Comment) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = c.MarshalTextStream(b)
	return b.Bytes(), err
}
func (c *Comment) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = c.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}

type Embedded struct {
	Value *Value
}

func NewEmbedded(v Value) Embedded {
	return Embedded{Value: &v}
}

func (e *Embedded) Equal(y Value) (b bool) {
	x, ok := y.(*Embedded)
	if !ok {
		return false
	}
	return e.Equal(x)
}

func (e *Embedded) BinaryRune() rune {
	return EmbeddedRune
}
func (e *Embedded) TextRune(p Position) string {
	if p == START {
		return "#:"
	}
	return ""
}
func (e *Embedded) MarshalBinaryStream(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte{byte(e.BinaryRune())})
	if err != nil {
		return
	}
	var m int

	m, err = Value(*e.Value).MarshalBinaryStream(w)
	n += m
	return
}
func (e *Embedded) UnmashalBinaryStream(r io.Reader) (n int, err error) {
	if n, err = ReadRune(e, r); err != nil {
		return
	}

	var m int
	var v Value
	v, m, err = ReadValueFromBinary(r)
	n += m
	*e = Embedded{Value: &v}
	return
}
func (e *Embedded) MarshalTextStream(w io.Writer) (n int, err error) {

	n, err = fmt.Fprintf(w, "%s", e.TextRune(START))
	if err != nil {
		return
	}
	var m int
	m, err = (*e.Value).MarshalTextStream(w)
	n += m
	return
}
func (e *Embedded) UnmarshalTextStream(r io.Reader) (n int, err error) {
	pr, ok := r.(*PeekReader)
	if !ok {
		pr = NewPeekReader(r)
	}
	n, err = ReadStringRune(e, START, pr)
	if err != nil {
		return
	}
	var (
		m int
		v Value
	)
	v, m, err = ReadValueFromText(pr, nil)
	n += m
	if err != nil {
		return
	}
	(*e.Value) = v
	return
}
func (e *Embedded) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = e.MarshalBinaryStream(b)
	return b.Bytes(), err
}
func (e *Embedded) UnmarshalBinary(data []byte) (err error) {
	var n int
	n, err = e.UnmashalBinaryStream(bytes.NewBuffer(data))
	data = data[n:]
	if len(data) > 0 {
		err = &MoreToRead{
			Read: n,
			Left: len(data),
		}
	}
	return
}
func (e *Embedded) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = e.MarshalTextStream(b)
	return b.Bytes(), err
}
func (e *Embedded) UnmarshalText(data []byte) (err error) {
	var n int
	n, err = e.UnmarshalTextStream(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	return MaybeMoreToRead(len(data[n:]), n)
}
