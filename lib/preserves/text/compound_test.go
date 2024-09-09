package text

import (
	"encoding"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestParseCompound(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	strMap := map[string]string{
		"    <symbol symbol>   ":         "<symbol symbol>",
		"    <symbol #f>   ":             "<symbol #f>",
		"    <symbol <#f #t>>   ":        "<symbol <#f #t>>",
		"    <>   ":                      "<>",
		"    <|symbol|>   ":              "<|symbol|>",
		"  {key: symbol}  ":              "{\n  key: symbol\n}",
		"  {key: symbol key1: symbol}  ": "{\n  key: symbol\n  key1: symbol\n}",
		"  {<key1>: <symbol symbols> <key>: <symbol symbols>}  ": "{\n  <key>: <symbol symbols>\n  <key1>: <symbol symbols>\n}",
		"  [symbol1 symbol2]  ":                                  "[\n  symbol1\n  symbol2\n]",
		"  #{symbol symbol}  ":                                   "#{\n  symbol\n  symbol\n}",
		"  #{1 1 1 1 1 1 1}    ":                                 "#{\n  1\n  1\n  1\n  1\n  1\n  1\n  1\n}",
	}

	buf := &strings.Reader{}
	for k, e := range strMap {
		buf.Reset(k)
		pt := &TextParser{}
		_, err := pt.ReadFrom(buf)
		if err != nil {
			t.Fatal(err)
		}
		data, err := (pt.Result.(encoding.TextMarshaler)).MarshalText()
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != e {
			t.Fatalf("not equal %s != %s", string(data), e)
		}
	}
}
