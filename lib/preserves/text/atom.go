package text

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/isodude/preserves-go/lib/preserves"
	log "github.com/sirupsen/logrus"
)

type QuotedSymbol struct {
	preserves.Pstring
}

func (qs *QuotedSymbol) New() preserves.Value {
	return new(QuotedSymbol)
}

func (bs *QuotedSymbol) Equal(y preserves.Value) bool { return false }

func (qs *QuotedSymbol) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br      *bufio.Reader
		escaped bool
		ok      bool
		r       rune
		s       int
	)
	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}

	ok, s, err = EndRune(rune(QuotedSymbolRunes[0]), br)
	n += int64(s)
	if err != nil || !ok {
		return
	}

	parsed := []rune{}

	for {
		r, s, err = br.ReadRune()
		n += int64(s)
		log.Debugf("qs: %s, %v\n", string(r), err)
		if !escaped && r == rune(QuotedSymbolRunes[1]) {
			break
		}

		if escaped {
			escaped = false
			switch r {
			case '\\':
				parsed = append(parsed, '\\')
			case 'n':
				parsed = append(parsed, '\n')
			case 'r':
				parsed = append(parsed, '\r')
			case 't':
				parsed = append(parsed, '\t')
			default:
				parsed = append(parsed, r)
			}
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		log.Debugf("debug: %v\n", parsed)
		parsed = append(parsed, r)
	}
	qs.Pstring = preserves.Pstring(parsed)

	return
}

func (qs *QuotedSymbol) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	m, err = fmt.Fprintf(w, "|%s|", strings.Replace(string(qs.Pstring), "|", "\\|", -1))
	n += int64(m)
	return
}

func (qs *QuotedSymbol) MarshalText() (data []byte, err error) {
	buf := &bytes.Buffer{}
	_, err = qs.WriteTo(buf)
	data = buf.Bytes()
	return
}

func (qs *QuotedSymbol) String() string {
	b, _ := qs.MarshalText()
	return string(b)
}

type BareSymbol struct {
	preserves.Symbol
}

func NewBareSymbol(s string) preserves.Value {
	return &([]BareSymbol{{Symbol: *preserves.NewSymbol(s)}}[0])
}

func (bs *BareSymbol) New() preserves.Value {
	return new(BareSymbol)
}

func (bs *BareSymbol) Equal(y preserves.Value) bool {
	x, ok := y.(*BareSymbol)
	if !ok {
		return false
	}
	return *x == *bs
}

func (bs *BareSymbol) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br     *bufio.Reader
		ok     bool
		s      int
		r      rune
		parsed []rune
	)
	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}

	for {
		r, s, err = br.ReadRune()
		n += int64(s)
		if err != nil && err != io.EOF {
			return
		}
		if err == io.EOF || !SymbolOrNumberRegexp.MatchString(string(r)) {
			if err != io.EOF {
				err = br.UnreadRune()
				n -= int64(s)
				if err != nil {
					return
				}
			} else {
				err = nil
			}

			log.Debugf("baresymbol: read: %s", string(parsed))
			bs.Symbol = preserves.Symbol(parsed)
			return
		}
		parsed = append(parsed, r)
	}
}

func (bs *BareSymbol) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	m, err = w.Write([]byte(bs.Symbol))
	n += int64(m)
	return
}

func (bs *BareSymbol) MarshalText() (data []byte, err error) {
	buf := &bytes.Buffer{}
	_, err = bs.WriteTo(buf)
	data = buf.Bytes()
	return
}

func (bs *BareSymbol) String() string {
	b, _ := bs.MarshalText()
	return string(b)
}

type Boolean struct {
	preserves.Boolean
}

func (b *Boolean) New() preserves.Value {
	return new(Boolean)
}

func (b *Boolean) Equal(y preserves.Value) bool {
	x, ok := y.(*Boolean)
	if !ok {
		return false
	}
	return x.Boolean == b.Boolean
}

func (b *Boolean) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br *bufio.Reader
		ok bool
		s  int
		bs []byte
	)

	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}

	bs = make([]byte, 2)
	s, err = br.Read(bs)
	n += int64(s)
	if err != nil {
		return
	}

	switch string(bs) {
	case BooleanTrueRunes:
		b.Boolean = preserves.Boolean(true)
	case BooleanFalseRunes:
		b.Boolean = preserves.Boolean(false)
	default:
		err = fmt.Errorf("parse error: unreachable code")
		return
	}
	return
}

func (b *Boolean) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	if b.Boolean {
		m, err = w.Write([]byte(BooleanTrueRunes))
	} else {
		m, err = w.Write([]byte(BooleanFalseRunes))
	}
	n += int64(m)
	return
}

func (b *Boolean) MarshalText() (data []byte, err error) {
	buf := &bytes.Buffer{}
	_, err = b.WriteTo(buf)
	data = buf.Bytes()
	return
}

func (b *Boolean) String() string {
	bs, _ := b.MarshalText()
	return string(bs)
}

type BareDouble struct {
	preserves.Double
}

func (b *BareDouble) New() preserves.Value {
	return new(BareDouble)
}
func (b *BareDouble) Equal(y preserves.Value) bool {
	x, ok := y.(*BareDouble)
	if !ok {
		return false
	}
	return x.Double == b.Double
}
func (b *BareDouble) WriteTo(w io.Writer) (n int64, err error) {
	var s int
	s, err = fmt.Fprintf(w, "%f", b.Double)
	n = int64(s)
	return
}
func (b *BareDouble) ReadFrom(r io.Reader) (n int64, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	var (
		bs  []byte
		str []byte
	)
	for {
		bs, err = br.Peek(1)
		n += 1
		if err != nil {
			return
		}
		if !SymbolOrNumberRegexp.Match(bs) {
			var a float64
			a, err = strconv.ParseFloat(string(str), 64)
			if err != nil {
				return
			}
			b.Double = preserves.Double(a)
			return
		}
		_, err = br.ReadByte()
		n += 1
		if err != nil {
			return
		}
		str = append(str, bs[0])
	}
}

func (b *BareDouble) MarshalText() (data []byte, err error) {
	buf := &bytes.Buffer{}
	_, err = b.WriteTo(buf)
	data = buf.Bytes()
	return
}

func (b *BareDouble) String() string {
	bs, _ := b.MarshalText()
	return string(bs)
}

type HexDouble struct {
	preserves.Double
}

func (h *HexDouble) New() preserves.Value {
	return new(HexDouble)
}
func (h *HexDouble) Equal(y preserves.Value) bool {
	x, ok := y.(*HexDouble)
	if !ok {
		return false
	}
	return h.Double == x.Double
}

func (h *HexDouble) WriteTo(w io.Writer) (n int64, err error) {
	var s int
	s, err = w.Write([]byte(HexDoubleRunes[0:4]))
	n = int64(s)
	if err != nil {
		return
	}
	bs := binary.BigEndian.AppendUint64([]byte{}, math.Float64bits(float64(h.Double)))
	var notfirst bool
	for _, b := range bs {
		if notfirst {
			s, err = w.Write([]byte{' '})
			n += int64(s)
			if err != nil {
				return
			}
		} else {
			notfirst = true
		}
		s, err = fmt.Fprintf(w, "%02x", b)
		n += int64(s)
		if err != nil {
			return
		}
	}

	s, err = w.Write([]byte{HexDoubleRunes[4]})
	n += int64(s)

	log.Debugf("hexdouble: end: %d, %v", s, err)
	return
}
func (h *HexDouble) ReadFrom(r io.Reader) (n int64, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	bs := make([]byte, 4)
	var s int
	s, err = br.Read(bs)
	n += int64(s)
	if err != nil {
		return
	}
	if string(bs) != HexDoubleRunes[0:4] {
		err = fmt.Errorf("hexdouble: no runes, found (%s)", bs)
		return
	}

	var m int64
	f := make([]byte, 8)
	for i := range []int{0, 1, 2, 3, 4, 5, 6, 7} {
		m, err = SkipWhitespace(br)
		n += m
		if err != nil {
			return
		}
		bs = make([]byte, 2)
		s, err = br.Read(bs)
		n += int64(s)
		if err != nil {
			return
		}
		s, err = fmt.Sscanf(string(bs), "%x", &f[i])
		if err != nil {
			return
		}
		if s != 1 {
			err = fmt.Errorf("hexdouble: no matches")
			return
		}
	}
	m, err = SkipWhitespace(br)
	n += m
	if err != nil {
		return
	}
	bs = make([]byte, 1)
	s, err = br.Read(bs)
	n += int64(s)
	if err != nil {
		return
	}
	if bs[0] != HexDoubleRunes[4] {
		err = fmt.Errorf("parse error: missing end dquote %s != %s", string(bs), "\"")
		return
	}
	h.Double = preserves.Double(math.Float64frombits(binary.BigEndian.Uint64(f)))
	log.Debugf("hexdouble: %f", *h)
	return
}
func (h *HexDouble) MarshalText() (data []byte, err error) {
	buf := &bytes.Buffer{}
	_, err = h.WriteTo(buf)
	data = buf.Bytes()
	return
}

func (h *HexDouble) String() string {
	b, _ := h.MarshalText()
	return string(b)
}

type String struct {
	preserves.Pstring
}

func (s *String) New() preserves.Value {
	return new(String)
}
func (s *String) Equal(y preserves.Value) bool {
	x, ok := y.(*String)
	if !ok {
		return false
	}
	return x.Pstring == s.Pstring
}

func (s *String) WriteTo(w io.Writer) (n int64, err error) {
	var si int
	si, err = fmt.Fprintf(w, "%s%s%s", []byte{StringRunes[0]}, s.Pstring, []byte{StringRunes[1]})
	n = int64(si)
	return
}
func (s *String) ReadFrom(r io.Reader) (n int64, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	bs := make([]byte, 1)
	var si int
	var b byte

	si, err = br.Read(bs)
	n += int64(si)
	if err != nil {
		return
	}
	if bs[0] != StringRunes[0] {
		err = fmt.Errorf("string: no runes, found (%s)", bs)
		return
	}
	bs = []byte{}

	var escaped bool
	for {
		b, err = br.ReadByte()
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

		if b == StringRunes[1] {
			s.Pstring = preserves.Pstring(bs)
			log.Debugf("string: '%s'", bs)
			return
		}
		bs = append(bs, b)
	}
}

func (s *String) MarshalText() (data []byte, err error) {
	buf := &bytes.Buffer{}
	_, err = s.WriteTo(buf)
	data = buf.Bytes()
	return
}

func (s *String) String() string {
	b, _ := s.MarshalText()
	return string(b)
}

type SignedInteger struct {
	preserves.SignedInteger
}

func NewSignedInteger(i int64) *SignedInteger {
	return &([]SignedInteger{{preserves.SignedInteger(*big.NewInt(i))}}[0])
}
func (s *SignedInteger) New() preserves.Value {
	return NewSignedInteger(0)
}
func (s *SignedInteger) Equal(y preserves.Value) bool {
	x, ok := y.(*SignedInteger)
	if !ok {
		return false
	}
	a := big.Int(s.SignedInteger)
	b := big.Int(x.SignedInteger)
	return a.Cmp(&b) == 0
}

func (s *SignedInteger) WriteTo(w io.Writer) (n int64, err error) {
	a := big.Int(s.SignedInteger)
	var str []byte
	str, err = a.MarshalText()
	n = int64(len(str))
	if err != nil {
		return
	}
	var m int
	m, err = w.Write(str)
	n += int64(m)
	return
}
func (s *SignedInteger) ReadFrom(r io.Reader) (n int64, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	var a big.Int
	var bs []byte
	var str []byte
	for {
		bs, err = br.Peek(1)
		if err != nil {
			return
		}
		if !SymbolOrNumberRegexp.Match(bs) {
			a = big.Int(s.SignedInteger)
			err = a.UnmarshalText(str)
			if err != nil {
				return
			}
			s.SignedInteger = preserves.SignedInteger(a)
			return
		}
		_, err = br.ReadByte()
		n += 1
		if err != nil {
			return
		}
		str = append(str, bs[0])
	}
}
func (s *SignedInteger) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.WriteTo(b)
	return b.Bytes(), err
}
func (s *SignedInteger) String() string {
	b, _ := s.MarshalText()
	return string(b)
}

type Base64ByteString struct {
	preserves.ByteString
}

func (bs *Base64ByteString) New() preserves.Value {
	return &([]Base64ByteString{{preserves.ByteString{}}}[0])
}
func (bs *Base64ByteString) Equal(v preserves.Value) bool {
	x, ok := v.(*Base64ByteString)
	if !ok {
		return false
	}
	return slices.Equal([]byte(bs.ByteString), []byte(x.ByteString))
}
func (bs *Base64ByteString) WriteTo(w io.Writer) (n int64, err error) {
	var s int
	s, err = w.Write([]byte(Base64ByteStringRunes[0:2]))
	n += int64(s)
	if err != nil {
		return
	}

	s, err = w.Write(base64.RawURLEncoding.AppendEncode([]byte{}, bs.ByteString))
	n += int64(s)
	if err != nil {
		return
	}

	s, err = w.Write([]byte{Base64ByteStringRunes[2]})
	n += int64(s)
	return
}
func (bs *Base64ByteString) ReadFrom(r io.Reader) (n int64, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	var s int
	bsb := make([]byte, 2)
	s, err = br.Read(bsb)
	n += int64(s)
	if err != nil {
		return
	}
	if string(bsb) != Base64ByteStringRunes[0:2] {
		err = fmt.Errorf("base64bytestring: no rune, found %s", bsb)
		return
	}
	var b byte
	var m int64
	bsb = []byte{}
	for {
		m, err = SkipWhitespace(br)
		n += m
		if err != nil {
			return
		}
		b, err = br.ReadByte()
		n += 1
		if err != nil {
			return
		}
		if b == Base64ByteStringRunes[2] {
			bsb, err = base64.RawURLEncoding.AppendDecode([]byte{}, bsb)
			if err != nil {
				return
			}
			bs.ByteString = preserves.ByteString(bsb)
			return
		}

		if !Base64Regexp.Match([]byte{b}) {
			err = fmt.Errorf("base64bytestring: unknown byte: %x", b)
			return
		}
		if b == '=' {
			continue
		}
		bsb = append(bsb, b)
		log.Debugf("base64bytestring: loop: %s", string(bsb))
	}
}

func (bs *Base64ByteString) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = bs.WriteTo(b)
	return b.Bytes(), err
}
func (bs *Base64ByteString) String() string {
	b, _ := bs.MarshalText()
	return string(b)
}

type HexByteString struct {
	preserves.ByteString
}

func (hbs *HexByteString) New() preserves.Value {
	return new(HexByteString)
}

func (hbs *HexByteString) Equal(v preserves.Value) bool {
	x, ok := v.(*HexByteString)
	if !ok {
		return false
	}
	return slices.Equal([]byte(hbs.ByteString), []byte(x.ByteString))
}

func (hbs *HexByteString) ReadFrom(r io.Reader) (n int64, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	var s int
	bsb := make([]byte, 3)
	s, err = br.Read(bsb)
	n += int64(s)
	if err != nil {
		return
	}
	if string(bsb) != HexByteStringRunes[0:3] {
		err = fmt.Errorf("base64bytestring: no rune, found %s", bsb)
		return
	}

	var (
		m int64
		u uint64
	)
	bs := make([]byte, 2)
	bsb = []byte{}
	for {
		m, err = SkipWhitespace(br)
		n += m
		if err != nil {
			return
		}

		ok, s, err = EndRune(rune(HexByteStringRunes[3]), br)
		n += int64(s)
		if ok {
			hbs.ByteString = preserves.ByteString(bsb)
			return
		}
		if err != nil {
			return
		}
		s, err = br.Read(bs)
		n += int64(s)
		if err != nil {
			return
		}

		u, err = strconv.ParseUint(string(bs), 16, 8)
		if err != nil {
			return
		}
		bsb = append(bsb, byte(u))
	}
}
func (hbs *HexByteString) WriteTo(w io.Writer) (n int64, err error) {
	var s int
	s, err = w.Write([]byte(HexByteStringRunes[0:3]))
	n = int64(s)
	if err != nil {
		return
	}
	var notfirst bool
	for _, b := range hbs.ByteString {
		if notfirst {
			s, err = w.Write([]byte{' '})
			n += int64(s)
			if err != nil {
				return
			}
		} else {
			notfirst = true
		}
		s, err = fmt.Fprintf(w, "%02x", b)
		n += int64(s)
		if err != nil {
			return
		}
	}

	s, err = w.Write([]byte{HexByteStringRunes[3]})
	n += int64(s)
	return
}

func (hbs *HexByteString) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = hbs.WriteTo(b)
	return b.Bytes(), err
}

func (hbs *HexByteString) String() string {
	b, _ := hbs.MarshalText()
	return string(b)
}

type BinByteString struct {
	preserves.ByteString
}

func (bbs *BinByteString) New() preserves.Value {
	return new(BinByteString)
}

func (bbs *BinByteString) Equal(v preserves.Value) bool {
	x, ok := v.(*BinByteString)
	if !ok {
		return false
	}
	return slices.Equal([]byte(bbs.ByteString), []byte(x.ByteString))
}
func (bbs *BinByteString) WriteTo(w io.Writer) (n int64, err error) {
	var s int
	s, err = w.Write([]byte(BinByteStringRunes[0:2]))
	n = int64(s)
	if err != nil {
		return
	}
	for _, b := range bbs.ByteString {
		escape := []byte{'"', '/', '\b', '\f', '\n', '\r', '\t'}
		escapeMap := map[byte]string{
			'"':  "\"",
			'/':  "/",
			'\b': "b",
			'\f': "f",
			'\n': "n",
			'\r': "r",
			'\t': "t",
		}
		if b >= 32 && b <= 126 {
			s, err = fmt.Fprintf(w, "%s", string(b))
			n += int64(s)
			if err != nil {
				return
			}
		} else if slices.Contains(escape, b) {
			s, err = fmt.Fprintf(w, "\\%s", escapeMap[b])
			n += int64(s)
			if err != nil {
				return
			}
		} else {
			s, err = fmt.Fprintf(w, "\\x%02x", b)
			n += int64(s)
			if err != nil {
				return
			}
		}

	}

	s, err = w.Write([]byte{BinByteStringRunes[2]})
	n += int64(s)
	return
}
func (bbs *BinByteString) ReadFrom(r io.Reader) (n int64, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	var s int
	bsb := make([]byte, 2)
	s, err = br.Read(bsb)
	n += int64(s)
	if err != nil {
		return
	}
	if string(bsb) != BinByteStringRunes[0:2] {
		err = fmt.Errorf("binbytestring: no rune, found %s", bsb)
		return
	}

	var (
		b       byte
		bytes   []byte
		escaped bool
		hex     bool
	)
	for {
		if !escaped {
			ok, s, err = EndRune(rune(BinByteStringRunes[2]), br)
			n += int64(s)
			if ok {
				bbs.ByteString = preserves.ByteString(bytes)
				return
			}
			if err != nil {
				return
			}
		}
		b, err = br.ReadByte()
		n += 1
		if err != nil {
			return
		}

		if hex {
			hex = false
			b1 := b
			b, err = br.ReadByte()
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
			continue
		}
		bytes = append(bytes, b)
	}
}

func (bbs *BinByteString) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = bbs.WriteTo(b)
	return b.Bytes(), err
}

func (bbs *BinByteString) String() string {
	b, _ := bbs.MarshalText()
	return string(b)
}
