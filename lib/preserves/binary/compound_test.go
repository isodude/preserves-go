package binary

import (
	"bytes"
	"encoding"
	"slices"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestParseCompound(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	strMap := [][]byte{
		{AnnotationByte, SymbolByte, 0x01, 0x61, AnnotationByte, SymbolByte, 0x01, 0x62, SequenceByte, EndByte},
		{RecordByte, SymbolByte, 0x01, 0x61, SymbolByte, 0x01, 0x62, EndByte},
		{SequenceByte, SymbolByte, 0x01, 0x61, SymbolByte, 0x01, 0x62, EndByte},
		{SetByte, SymbolByte, 0x01, 0x61, SymbolByte, 0x01, 0x62, EndByte},
		{DictionaryByte, SymbolByte, 0x01, 0x61, SymbolByte, 0x01, 0x62, EndByte},
		{EmbeddedByte, SymbolByte, 0x01, 0x61},
		{RecordByte, EndByte},
		{SequenceByte, EndByte},
		{SetByte, EndByte},
		{DictionaryByte, EndByte},
	}
	for i, k := range strMap {
		buf := bytes.NewReader(k)
		pt := &BinaryParser{}
		_, err := pt.ReadFrom(buf)
		if err != nil {
			t.Fatalf("err (%d) %v", i, err)
		}
		data, err := (pt.Result.(encoding.BinaryMarshaler)).MarshalBinary()
		if err != nil {
			t.Fatalf("err (%d) %v", i, err)
		}
		if !slices.Equal(k, data) {
			t.Fatalf("not equal (%d) %x != %x", i, k, data)
		}
	}
}
