package binary

import (
	"bytes"
	"math/big"
	"slices"
	"testing"

	"github.com/isodude/preserves-go/lib/extras"
	"github.com/isodude/preserves-go/lib/preserves"
)

func TestParseAtomMarshal(t *testing.T) {
	bMap := [][]byte{
		{BooleanTrueByte},
		{BooleanFalseByte},
		{DoubleByte, 0x08, 0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		{DoubleByte, 0x08, 0xFE, 0x3C, 0xB7, 0xB7, 0x59, 0xBF, 0x04, 0x26},
		{SignedIntegerByte, 0x01, 0x01},
		{SignedIntegerByte, 0x00},
		{SignedIntegerByte, 0x01, 0xFF},
		{SignedIntegerByte, 0x02, 0x00, 0xFF},
		{SignedIntegerByte, 0x02, 0xFE, 0xFF},
		{SignedIntegerByte, 0x01, 0xFE},
		{SignedIntegerByte, 0x02, 0xFF, 0x00},
		{SignedIntegerByte, 0x02, 0x01, 0x00},
		{SignedIntegerByte, 0x02, 0xFF, 0x01},
		{SignedIntegerByte, 0x02, 0x7F, 0xFF},
		{SignedIntegerByte, 0x02, 0xFF, 0x7F},
		{SignedIntegerByte, 0x03, 0x00, 0x80, 0x00},
		{SignedIntegerByte, 0x01, 0x80},
		{SignedIntegerByte, 0x01, 0x7F},
		{SignedIntegerByte, 0x03, 0x00, 0xFF, 0xFF},
		{SignedIntegerByte, 0x01, 0x81},
		{SignedIntegerByte, 0x02, 0x00, 0x80},
		{SignedIntegerByte, 0x03, 0x01, 0x00, 0x00},
		{SignedIntegerByte,
			0x12, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00,
		},
		{SymbolByte, 0x06, 0x74, 0x69, 0x74, 0x6C, 0x65, 0x64},
		{ByteStringByte, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67},
	}

	bp := BinaryParser{}
	for _, k := range bMap {
		buf := bytes.NewBuffer(k)
		_, err := bp.ReadFrom(buf)
		if err != nil {
			t.Fatal(err)
		}

		wbuf := new(bytes.Buffer)

		_, err = bp.Result.WriteTo(wbuf)
		if err != nil {
			t.Fatal(err)
		}

		if !slices.Equal[[]byte](wbuf.Bytes(), k) {
			t.Fatalf("not equal: %x != %x", wbuf.Bytes(), k)
		}

	}
}

type equal struct {
	bs    []byte
	value preserves.Value
}

func TestParseAtomEqual(t *testing.T) {
	bigInt := big.NewInt(0)
	bigInt.SetString("87112285931760246646623899502532662132736", 10)
	bMap := []equal{
		{[]byte{BooleanTrueByte}, extras.Reference(Boolean{preserves.Boolean(true)})},
		{[]byte{BooleanFalseByte}, extras.Reference(Boolean{preserves.Boolean(false)})},
		{[]byte{DoubleByte, 0x08, 0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, extras.Reference(Double{preserves.Double(1.0)})},
		{[]byte{DoubleByte, 0x08, 0xFE, 0x3C, 0xB7, 0xB7, 0x59, 0xBF, 0x04, 0x26}, extras.Reference(Double{preserves.Double(-1.202e300)})},
		{[]byte{SignedIntegerByte, 0x01, 0x01}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(1))})},
		{[]byte{SignedIntegerByte, 0x00}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(0))})},
		{[]byte{SignedIntegerByte, 0x01, 0xFF}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(-1))})},
		{[]byte{SignedIntegerByte, 0x02, 0x00, 0xFF}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(255))})},
		{[]byte{SignedIntegerByte, 0x02, 0xFE, 0xFF}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(-257))})},
		{[]byte{SignedIntegerByte, 0x01, 0xFE}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(-2))})},
		{[]byte{SignedIntegerByte, 0x02, 0xFF, 0x00}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(-256))})},
		{[]byte{SignedIntegerByte, 0x02, 0x01, 0x00}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(256))})},
		{[]byte{SignedIntegerByte, 0x02, 0xFF, 0x01}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(-255))})},
		{[]byte{SignedIntegerByte, 0x02, 0x7F, 0xFF}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(32767))})},
		{[]byte{SignedIntegerByte, 0x02, 0xFF, 0x7F}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(-129))})},
		{[]byte{SignedIntegerByte, 0x03, 0x00, 0x80, 0x00}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(32768))})},
		{[]byte{SignedIntegerByte, 0x01, 0x80}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(-128))})},
		{[]byte{SignedIntegerByte, 0x01, 0x7F}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(127))})},
		{[]byte{SignedIntegerByte, 0x03, 0x00, 0xFF, 0xFF}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(65535))})},
		{[]byte{SignedIntegerByte, 0x01, 0x81}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(-127))})},
		{[]byte{SignedIntegerByte, 0x02, 0x00, 0x80}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(128))})},
		{[]byte{SignedIntegerByte, 0x03, 0x01, 0x00, 0x00}, extras.Reference(SignedInteger{preserves.SignedInteger(*big.NewInt(65536))})},
		{[]byte{SignedIntegerByte, 0x12, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, extras.Reference(SignedInteger{preserves.SignedInteger(*bigInt)})},
		{[]byte{SymbolByte, 0x06, 0x74, 0x69, 0x74, 0x6C, 0x65, 0x64}, extras.Reference(Symbol{preserves.Symbol("titled")})},
		{[]byte{ByteStringByte, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67}, extras.Reference(ByteString{preserves.ByteString([]byte("teststring"))})},
	}

	for i, k := range bMap {
		buf := bytes.NewBuffer(k.bs)
		v := k.value.New()
		_, err := v.ReadFrom(buf)
		if err != nil {
			t.Fatal(err)
		}

		if !v.Equal(k.value) {
			t.Fatalf("not equal: (%d) %v != %v", i, v, k.value)
		}
	}
}
