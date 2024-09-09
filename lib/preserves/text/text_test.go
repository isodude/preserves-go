package text

import (
	"bufio"
	"encoding"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestSamples(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	f, err := os.Open("../../../tests/samples.pr")
	if err != nil {
		t.Fatal(err)
	}
	buf := bufio.NewReader(f)
	buf = bufio.NewReaderSize(buf, 200000)
	tp := TextParser{}
	_, err = tp.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tp.Result.(encoding.TextMarshaler).MarshalText()
	if err != nil {
		t.Fatal(err)
	}
}
