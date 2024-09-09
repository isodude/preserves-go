package lib

import (
	"bytes"
	"io"
	"math/big"
	"os"
	"slices"
	"testing"
)

func TestValuesSymbol(t *testing.T) {
	a := NewSymbol("a")
	a1 := NewSymbol("a")
	b := NewSymbol("b")
	if a != a1 {
		t.Fatalf("%v != %v", a, a1)
	}
	if a == b {
		t.Fatalf("%v == %v", a, b)
	}
	if a > b {
		t.Fatalf("%v > %v", a, b)
	}
}

func TestValuesSymbolBinaryMarshal(t *testing.T) {
	hello := NewSymbol("titled")
	b, err := hello.MarshalBinary()
	if err != nil {
		t.Fatalf("unable to marshal symbol: %v", err)
	}
	c := []byte{0xB3, 0x06, 0x74, 0x69, 0x74, 0x6C, 0x65, 0x64}
	if slices.Compare[[]byte, byte](b, c) != 0 {
		t.Fatalf("marshal output was not correct (%v != %v)", b, c)
	}
	s := NewSymbol("")
	err = (&s).UnmarshalBinary(b)
	if err != nil {
		t.Fatalf("unable to unmarshal symbol: %v", err)
	}
	if !hello.Equal(&s) {
		t.Fatalf("symbols not equal: (%s != %s)", hello, s)
	}
}

func TestValuesBooleanTrueBinaryMarshal(t *testing.T) {
	True := NewBoolean(true)

	b, err := (&True).MarshalBinary()
	if err != nil {
		t.Fatalf("unable to marshal boolean: %v", err)
	}

	c := []byte{BooleanTrueRune}
	if slices.Compare[[]byte, byte](b, c) != 0 {
		t.Fatalf("marshal output was not correct (%v != %v)", b, c)
	}

	False := NewBoolean(false)
	err = (&False).UnmarshalBinary(b)
	if err != nil {
		t.Fatalf("unable to unmarshal boolean: %v", err)
	}
	if !(&True).Equal(&False) {
		t.Fatalf("booleans not equal: (%v != %v)", True, False)
	}
}
func TestValuesBooleanFalseBinaryMarshal(t *testing.T) {
	False := NewBoolean(false)
	b, err := (&False).MarshalBinary()
	if err != nil {
		t.Fatalf("unable to marshal boolean: %v", err)
	}

	c := []byte{BooleanFalseRune}
	if slices.Compare[[]byte, byte](b, c) != 0 {
		t.Fatalf("marshal output was not correct (%v != %v)", b, c)
	}

	True := NewBoolean(true)
	err = (&True).UnmarshalBinary(b)
	if err != nil {
		t.Fatalf("unable to unmarshal boolean: %v", err)
	}
	if !(&False).Equal(&True) {
		t.Fatalf("booleans not equal: (%v != %v)", True, False)
	}
}

func TestValuesSignedIntegerBinaryMarshal(t *testing.T) {
	for k, v := range map[int64][]byte{
		1:     {0x01, 0x01},
		0:     {0x00},
		-1:    {0x01, 0xFF},
		255:   {0x02, 0x00, 0xFF},
		-257:  {0x02, 0xFE, 0xFF},
		-2:    {0x01, 0xFE},
		-256:  {0x02, 0xFF, 0x00},
		256:   {0x02, 0x01, 0x00},
		-255:  {0x02, 0xFF, 0x01},
		32767: {0x02, 0x7F, 0xFF},
		-129:  {0x02, 0xFF, 0x7F},
		32768: {0x03, 0x00, 0x80, 0x00},
		-128:  {0x01, 0x80},
		127:   {0x01, 0x7F},
		65535: {0x03, 0x00, 0xFF, 0xFF},
		-127:  {0x01, 0x81},
		128:   {0x02, 0x00, 0x80},
		65536: {0x03, 0x01, 0x00, 0x00},
	} {
		integer := NewSignedInteger(k)
		b, err := (&integer).MarshalBinary()
		if err != nil {
			t.Fatalf("unable to marshal signedinteger: %v", err)
		}

		c := []byte{SignedIntegerRune}
		c = append(c, v...)
		if slices.Compare[[]byte, byte](b, c) != 0 {
			t.Fatalf("marshal output was not correct for %d: (%v != %v)", k, b, c)
		}

		zero := NewSignedInteger(0)
		err = (&zero).UnmarshalBinary(b)
		if err != nil {
			t.Fatalf("unable to unmarshal signedinteger: %v", err)
		}
		if !(&zero).Equal(&integer) {
			t.Fatalf("signedinteger not equal: %d: (%v(%v) != %v(%v))", k, zero, c, integer, b)
		}
	}
}

func TestValuesBigSignedIntegerBinaryMarshal(t *testing.T) {
	a := big.Int{}
	a.SetString("87112285931760246646623899502532662132736", 10)
	c := []byte{SignedIntegerRune,
		0x12, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
	}
	integer := SignedInteger(a)
	b, err := (&integer).MarshalBinary()
	if err != nil {
		t.Fatalf("unable to marshal signedinteger: %v", err)
	}

	if slices.Compare[[]byte, byte](b, c) != 0 {
		t.Fatalf("marshal output was not correct for (%v != %v)", b, c)
	}
	zero := NewSignedInteger(0)
	err = (&zero).UnmarshalBinary(b)
	if err != nil {
		t.Fatalf("unable to unmarshal signedinteger: %v", err)
	}
	if !(&zero).Equal(&integer) {
		t.Fatalf("signedinteger not equal: (%v(%v) != %v(%v))", zero, c, integer, b)
	}
}
func TestValuesDoubleBinaryMarshal(t *testing.T) {
	for k, v := range map[float64][]byte{
		1.0:        {0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		-1.202e300: {0xFE, 0x3C, 0xB7, 0xB7, 0x59, 0xBF, 0x04, 0x26},
	} {
		double := NewDouble(k)
		b, err := (&double).MarshalBinary()
		if err != nil {
			t.Fatalf("unable to marshal double: %v", err)
		}

		c := []byte{DoubleRune, 0x08}
		c = append(c, v...)
		if slices.Compare[[]byte, byte](b, c) != 0 {
			t.Fatalf("marshal output was not correct for %f: (%v != %v)", k, b, c)
		}

		zero := NewDouble(0)
		err = (&zero).UnmarshalBinary(b)
		if err != nil {
			t.Fatalf("unable to unmarshal double: %v", err)
		}
		if !(&zero).Equal(&double) {
			t.Fatalf("double not equal: %f: (%v(%v) != %v(%v))", k, zero, c, double, b)
		}
	}
}

func TestValuesSymbolTextMarshal(t *testing.T) {
	payload := "titled"
	hello := NewSymbol(payload)
	b, err := hello.MarshalText()
	if err != nil {
		t.Fatalf("unable to marshal symbol: %v", err)
	}
	if string(b) != payload {
		t.Fatalf("marshal output was not correct (%s != %s)", string(b), payload)
	}
	s := NewSymbol("")
	err = (&s).UnmarshalText(b)
	if err != nil {
		t.Fatalf("unable to unmarshal symbol: %v", err)
	}
	if !hello.Equal(&s) {
		t.Fatalf("symbols not equal: (%s != %s)", hello, s)
	}
}

func TestValuesStringBinaryMarshal(t *testing.T) {
	payload := "titled"
	text := NewString(payload)
	b, err := text.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	c := []byte{StringRune, byte(len(payload))}
	c = append(c, []byte(payload)...)
	if slices.Compare[[]byte, byte](b, c) != 0 {
		t.Fatalf("marshal output was not correct for (%v != %v)", b, c)
	}
	empty := NewString("")
	err = empty.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if !text.Equal(&empty) {
		t.Fatalf("text not equal (%s != %s)", text, empty)
	}
}

func TestValuesByteStringBinaryMarshal(t *testing.T) {
	payload := "titled"
	text := NewByteString([]byte(payload))
	b, err := text.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	c := []byte{ByteStringRune, byte(len(payload))}
	c = append(c, []byte(payload)...)
	if slices.Compare[[]byte, byte](b, c) != 0 {
		t.Fatalf("marshal output was not correct for (%v != %v)", b, c)
	}
	empty := NewByteString([]byte{})
	err = empty.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if !text.Equal(&empty) {
		t.Fatalf("text not equal (%s != %s)", text, empty)
	}
}
func TestValuesEmptySequenceBinaryMarshal(t *testing.T) {
	empty := NewSequence([]Value{})
	b, err := empty.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	c := []byte{SequenceRune, EndRune}
	if !slices.Equal(b, c) {
		t.Fatalf("slices not equal (%v != %v)", b, c)
	}
	str := NewString("text")
	text := NewSequence([]Value{&str})
	if err = text.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if !text.Equal(&empty) {
		t.Fatalf("not equal: (%v != %v)", empty, text)
	}

}

func TestValuesStringSequenceBinaryMarshal(t *testing.T) {
	str := NewString("text")
	text := NewSequence([]Value{&str})
	b, err := text.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	c := []byte{SequenceRune, StringRune, byte(len(string(str)))}
	c = append(c, []byte(string(str))...)
	c = append(c, EndRune)
	if !slices.Equal(b, c) {
		t.Fatalf("slices not equal (%v != %v)", b, c)
	}
	empty := NewSequence([]Value{})
	if err = empty.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if !text.Equal(&empty) {
		t.Fatalf("not equal: (%v != %v)", empty, text)
	}
}
func TestValuesStringSequenceTextMarshal(t *testing.T) {
	str := NewString("text")
	text := NewSequence([]Value{&str})
	b, err := text.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	c := "[\"text\"]"
	if string(b) != c {
		t.Fatalf("slices not equal (%v != %v)", string(b), c)
	}
	empty := NewSequence([]Value{})
	if err = empty.UnmarshalText(b); err != nil {
		t.Fatal(err)
	}
	if !text.Equal(&empty) {
		t.Fatalf("not equal: (%v != %v)", empty, text)
	}
}

func TestValuesStringRecordBinaryMarshal(t *testing.T) {
	str := NewSymbol("date")
	s1 := NewString("year")
	s2 := NewString("month")
	s3 := NewString("day")
	text := NewRecord(&str, []Value{&s1, &s2, &s3})
	b, err := text.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	c := []byte{RecordRune, SymbolRune, byte(len(string(str)))}
	c = append(c, []byte(string(str))...)
	c = append(c, StringRune, byte(len(string(s1))))
	c = append(c, []byte(s1)...)
	c = append(c, StringRune, byte(len(string(s2))))
	c = append(c, []byte(s2)...)
	c = append(c, StringRune, byte(len(string(s3))))
	c = append(c, []byte(s3)...)
	c = append(c, EndRune)
	if !slices.Equal(b, c) {
		t.Fatalf("slices not equal (%v != %v)", b, c)
	}
	empty := NewRecord(nil, []Value{})
	if err = empty.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if !text.Equal(&empty) {
		t.Fatalf("not equal: (%v != %v)", empty, text)
	}
}

func TestValuesStringSetBinaryMarshal(t *testing.T) {
	str := NewSymbol("date")
	s1 := NewString("year")
	s2 := NewString("month")
	s3 := NewString("day")
	text := NewSet([]Value{&str, &s1, &s2, &s3})
	b, err := text.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	c := []byte{SetRune, SymbolRune, byte(len(string(str)))}
	c = append(c, []byte(string(str))...)
	c = append(c, StringRune, byte(len(string(s1))))
	c = append(c, []byte(s1)...)
	c = append(c, StringRune, byte(len(string(s2))))
	c = append(c, []byte(s2)...)
	c = append(c, StringRune, byte(len(string(s3))))
	c = append(c, []byte(s3)...)
	c = append(c, EndRune)
	if len(b) != len(c) {
		t.Fatalf("slices not equal length (%d != %d) (%v != %v)", len(b), len(c), b, c)
	}
	empty := NewSet([]Value{})
	if err = empty.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if !text.Equal(&empty) {
		t.Fatalf("not equal: (%v != %v)", empty, text)
	}
}

func TestValuesStringDictionaryBinaryMarshal(t *testing.T) {
	str := NewSymbol("date")
	s1 := NewString("year")
	s2 := NewString("month")
	s3 := NewString("day")
	text := NewDictionary(map[Value]Value{&str: &s1, &s2: &s3})
	b, err := text.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	c := []byte{DictionaryRune, SymbolRune, byte(len(string(str)))}
	c = append(c, []byte(string(str))...)
	c = append(c, StringRune, byte(len(string(s1))))
	c = append(c, []byte(s1)...)
	c = append(c, StringRune, byte(len(string(s2))))
	c = append(c, []byte(s2)...)
	c = append(c, StringRune, byte(len(string(s3))))
	c = append(c, []byte(s3)...)
	c = append(c, EndRune)
	if len(b) != len(c) {
		t.Fatalf("slices not equal length (%d != %d) (%v != %v)", len(b), len(c), b, c)
	}
	empty := NewDictionary(map[Value]Value{})
	if err = empty.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if !text.Equal(&empty) {
		t.Fatalf("not equal: (%v != %v)", empty, text)
	}
}

func TestSamples(t *testing.T) {
	f, err := os.Open("../tests/samples.bin")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	for len(data) > 0 {
		v, err := NewValueFromBinary(data[0])
		if err != nil {
			t.Fatalf("loop: %v", err)
		}

		err = v.UnmarshalBinary(data)
		if err != nil {
			m, ok := err.(*MoreToRead)
			if !ok {
				t.Fatalf("err: %v", err)
			}
			data = data[m.Read:]
		} else {
			break
		}
	}
}

func TestPeekReader(t *testing.T) {
	b := []byte{'a', 'b', 'c', 'd'}
	buf := bytes.NewBuffer(b)
	pr := NewPeekReader(buf)
	bs, err := pr.Peek(1)
	if err != nil {
		t.Fatal(err)
	}
	if bs[0] != b[0] {
		t.Fatalf("%s != %s", string(bs[0]), string(b[0]))
	}
	bs, err = pr.Peek(1)
	if err != nil {
		t.Fatal(err)
	}
	if bs[0] != b[0] {
		t.Fatalf("%s != %s", string(bs[0]), string(b[0]))
	}

	bs, err = pr.Peek(2)
	if err != nil {
		t.Fatal(err)
	}
	if bs[0] != b[0] {
		t.Fatalf("%s != %s", string(bs[0]), string(b[0]))
	}

	if bs[1] != b[1] {
		t.Fatalf("%s != %s", string(bs[1]), string(b[1]))
	}

	bs, err = pr.Peek(2)
	if err != nil {
		t.Fatal(err)
	}
	if bs[0] != b[0] {
		t.Fatalf("%s != %s", string(bs[0]), string(b[0]))
	}

	if bs[1] != b[1] {
		t.Fatalf("%s != %s", string(bs[1]), string(b[1]))
	}

	by, err := pr.ReadByte()
	if err != nil {
		t.Fatal(err)
	}

	if by != b[0] {
		t.Fatalf("%s != %s", string(by), string(b[0]))
	}

	bs, err = pr.Peek(1)
	if err != nil {
		t.Fatal(err)
	}

	if bs[0] != b[1] {
		t.Fatalf("%s != %s", string(bs[0]), string(b[1]))
	}

	bsb := make([]byte, 1)

	_, err = pr.Read(bsb)
	if err != nil {
		t.Fatal(err)
	}

	if bsb[0] != b[1] {
		t.Fatalf("%s != %s", string(bsb[0]), string(b[1]))
	}

}

func TestSchema(t *testing.T) {
	f, err := os.Open("../tests/samples.pr")
	if err != nil {
		t.Fatal(err)
	}
	pr := PeekReader{Reader: f}
	if err != nil {
		t.Fatal(err)
	}

	var v Value
	for {
		v, _, err = ReadValueFromText(&pr, nil)
		if err == io.EOF {
			return
		}
		if err != nil {
			t.Fatalf("loop: %v, %d, %s", err, pr.Position, pr.ReadData)
		}
		_, err = v.MarshalText()
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
	}
}
