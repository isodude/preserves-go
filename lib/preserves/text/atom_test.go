package text

import (
	"encoding"
	"fmt"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestParseAtom(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	strMap := map[string]string{
		"    |str\\|in\\ng|   ":              "|str\\|in\ng|",
		"    string   ":                      "string",
		"    # string   \ncommented":         "# string   \ncommented",
		"#!/bin/bash\nstart":                 "#!/bin/bash\nstart",
		"    #f   ":                          "#f",
		"  \"string\"  ":                     "\"string\"",
		"   #xd\"00 00 00 00 00 00 00 00\" ": "#xd\"00 00 00 00 00 00 00 00\"",
		"  538 ":                             "538",
		"  0.5 ":                             "0.500000",
		"  #[ Y29 yeW 1i ]  ":                "#[Y29yeW1i]",
		"  #\"hello\"  ":                     "#\"hello\"",
		"  #x\" 00 00 00 \"    ":             "#x\"00 00 00\"",
		fmt.Sprintf("  #\"hello\\n%s\"   ", string([]byte{0xff})): "#\"hello\\n\\xff\"",
	}
	for k, e := range strMap {
		buf := strings.NewReader(k)
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
