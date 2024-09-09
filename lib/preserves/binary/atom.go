package binary

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"slices"

	"github.com/isodude/preserves-go/lib/preserves"
	log "github.com/sirupsen/logrus"
)

type Boolean struct {
	preserves.Boolean
}

func (b *Boolean) New() preserves.Value {
	return new(Boolean)
}
func (b *Boolean) Equal(y preserves.Value) bool {
	return b.Boolean.Equal(ToPreserves(y))
}
func (b *Boolean) Cmp(y preserves.Value) int {
	return b.Boolean.Cmp(ToPreserves(y))
}
func (b *Boolean) ReadFrom(ir io.Reader) (n int64, err error) {
	bs := make([]byte, 1)
	var m int
	m, err = ir.Read(bs)
	n = int64(m)
	if err != nil {
		return
	}
	if bs[0] == BooleanTrueByte {
		b.Boolean = preserves.Boolean(true)
		return
	}
	if bs[0] == BooleanFalseByte {
		b.Boolean = preserves.Boolean(false)
		return
	}

	err = fmt.Errorf("boolean: no byte, found (%x)", bs)
	return
}

func (b *Boolean) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	if bool(b.Boolean) {
		m, err = w.Write([]byte{BooleanTrueByte})
		n = int64(m)
		return
	}

	m, err = w.Write([]byte{BooleanFalseByte})
	n = int64(m)
	return
}

func (b *Boolean) MarshalBinary() (data []byte, err error) {
	bs := &bytes.Buffer{}
	_, err = b.WriteTo(bs)
	return bs.Bytes(), err

}

type Double struct {
	preserves.Double
}

func (d *Double) New() preserves.Value {
	return new(Double)
}

func (d *Double) Equal(y preserves.Value) bool {
	return d.Double.Equal(ToPreserves(y))
}

func (d *Double) Cmp(y preserves.Value) int {
	return d.Double.Cmp(ToPreserves(y))
}
func (d *Double) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	// Length is always 64bit (0x08)
	m, err = w.Write([]byte{DoubleByte, 0x08})
	n = int64(m)
	if err != nil {
		return
	}

	m, err = w.Write(binary.BigEndian.AppendUint64([]byte{}, math.Float64bits(float64(d.Double))))
	n += int64(m)
	return
}

func (d *Double) ReadFrom(r io.Reader) (n int64, err error) {
	var m int
	bs := make([]byte, 2)
	m, err = r.Read(bs)
	n = int64(m)
	if err != nil {
		return
	}
	if !slices.Equal(bs, []byte{DoubleByte, 0x08}) {
		err = fmt.Errorf("double: unable to read bytes")
		return
	}
	bs = make([]byte, 8)
	m, err = r.Read(bs)
	n += int64(m)
	if err != nil {
		return
	}
	d.Double = preserves.Double(math.Float64frombits(binary.BigEndian.Uint64(bs)))
	return
}
func (d *Double) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = d.WriteTo(b)
	return b.Bytes(), err

}

type SignedInteger struct {
	preserves.SignedInteger
}

func (s *SignedInteger) New() preserves.Value {
	return new(SignedInteger)
}

func (s *SignedInteger) Equal(y preserves.Value) bool {
	return s.SignedInteger.Equal(ToPreserves(y))
}
func (s *SignedInteger) Cmp(y preserves.Value) int {
	return s.SignedInteger.Cmp(ToPreserves(y))
}
func (s *SignedInteger) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	m, err = w.Write([]byte{SignedIntegerByte})
	n = int64(m)
	if err != nil {
		return
	}

	a := big.Int(s.SignedInteger)
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
		n += int64(m)
		return
	}
	m, err = w.Write(binary.AppendUvarint([]byte{}, uint64(l)))
	n += int64(m)
	if err != nil {
		return
	}
	m, err = w.Write(b)
	n += int64(m)
	return
}

func (s *SignedInteger) ReadFrom(r io.Reader) (n int64, err error) {
	var m int
	bs := make([]byte, 1)
	m, err = r.Read(bs)
	n = int64(m)
	if err != nil {
		return
	}

	if bs[0] != SignedIntegerByte {
		err = fmt.Errorf("signedinteger: no byte, found (%x)", bs)
		return
	}

	var l uint64
	l, m, err = ReadUvarint(r)
	n += int64(m)
	if err != nil {
		return
	}
	a := big.Int(s.SignedInteger)
	if l == 0 {
		a.SetUint64(0)
		s.SignedInteger = preserves.SignedInteger(a)
		return
	}
	b := make([]byte, l)
	err = binary.Read(r, binary.BigEndian, b)
	n += int64(l)
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
	s.SignedInteger = preserves.SignedInteger(a)
	return
}

func (s *SignedInteger) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.WriteTo(b)
	return b.Bytes(), err
}

type String struct {
	preserves.Pstring
}

func (s *String) New() preserves.Value {
	return new(String)
}

func (s *String) Equal(y preserves.Value) bool {
	return s.Pstring.Equal(ToPreserves(y))
}

func (s *String) Cmp(y preserves.Value) int {
	return s.Pstring.Cmp(ToPreserves(y))
}
func (s *String) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	m, err = w.Write([]byte{StringByte})
	n = int64(m)
	if err != nil {
		return
	}
	b := binary.AppendUvarint([]byte{}, uint64(len(s.Pstring)))
	m, err = w.Write(b)
	n += int64(m)
	if err != nil {
		return
	}
	m, err = w.Write([]byte(s.Pstring))
	n += int64(m)
	return
}
func (s *String) ReadFrom(r io.Reader) (n int64, err error) {
	var m int
	bs := make([]byte, 1)
	m, err = r.Read(bs)
	n = int64(m)
	if err != nil {
		return
	}

	if bs[0] != StringByte {
		err = fmt.Errorf("string: no byte, found (%x)", bs)
		return
	}

	var l uint64
	l, m, err = ReadUvarint(r)
	n += int64(m)
	if err != nil {
		return
	}

	str := make([]byte, l)
	m, err = r.Read(str)
	n += int64(m)
	if err != nil {
		return
	}
	s.Pstring = preserves.Pstring(str)
	return
}
func (s *String) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.WriteTo(b)
	return b.Bytes(), err

}

type ByteString struct {
	preserves.ByteString
}

func (bs *ByteString) New() preserves.Value {
	return new(ByteString)
}
func (b *ByteString) Equal(y preserves.Value) bool {
	return b.ByteString.Equal(ToPreserves(y))
}

func (b *ByteString) Cmp(y preserves.Value) int {
	return b.ByteString.Cmp(ToPreserves(y))
}
func (bs *ByteString) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	m, err = w.Write([]byte{ByteStringByte})
	n = int64(m)
	if err != nil {
		return
	}
	b := binary.AppendUvarint([]byte{}, uint64(len(bs.ByteString)))
	m, err = w.Write(b)
	n += int64(m)
	if err != nil {
		return
	}
	m, err = w.Write([]byte(bs.ByteString))
	n += int64(m)
	return
}
func (bs *ByteString) ReadFrom(r io.Reader) (n int64, err error) {
	var m int
	b := make([]byte, 1)
	m, err = r.Read(b)
	n = int64(m)
	if err != nil {
		return
	}
	if b[0] != ByteStringByte {
		err = fmt.Errorf("bytestring: no byte, found (%x)", b)
		return
	}
	var l uint64
	l, m, err = ReadUvarint(r)
	n += int64(m)
	if err != nil {
		return
	}

	b = make([]byte, l)
	m, err = r.Read(b)
	n += int64(m)
	if err != nil {
		return
	}

	bs.ByteString = preserves.ByteString(b)
	return
}
func (bs *ByteString) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = bs.WriteTo(b)
	return b.Bytes(), err

}

type Symbol struct {
	preserves.Symbol
}

func (s *Symbol) New() preserves.Value {
	return new(Symbol)
}

func (s *Symbol) Equal(y preserves.Value) bool {
	return s.Symbol.Equal(ToPreserves(y))
}

func (s *Symbol) Cmp(y preserves.Value) int {
	return s.Symbol.Cmp(ToPreserves(y))
}
func (s *Symbol) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	m, err = w.Write([]byte{SymbolByte})
	n = int64(m)
	if err != nil {
		return
	}
	b := binary.AppendUvarint([]byte{}, uint64(len(string(s.Symbol))))
	m, err = w.Write(b)
	n += int64(m)
	if err != nil {
		return
	}
	m, err = w.Write([]byte(s.Symbol))
	n += int64(m)
	return
}
func (s *Symbol) ReadFrom(r io.Reader) (n int64, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	br.Size()
	var b byte
	var m int
	b, err = br.ReadByte()
	n = 1
	if err != nil {
		return
	}
	if b != SymbolByte {
		err = fmt.Errorf("symbol: no byte, found (%x)", b)
		return
	}
	var l uint64
	l, m, err = ReadUvarint(br)
	n += int64(m)
	if err != nil {
		return
	}

	log.Debugf("symbol: read length: %d", l)
	str := make([]byte, l)
	m, err = br.Read(str)
	n += int64(m)
	if uint64(m) != l {
		err = fmt.Errorf("symbol: did not read enough, %d != %d", l, m)
		return
	}
	if err != nil {
		return
	}
	log.Debugf("symbol: read: %s", str)
	s.Symbol = preserves.Symbol(str)
	return
}
func (s *Symbol) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.WriteTo(b)
	return b.Bytes(), err
}
