package binary

import (
	"bufio"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestSamples(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	f, err := os.Open("../../../tests/samples.bin")
	if err != nil {
		t.Fatal(err)
	}
	buf := bufio.NewReader(f)
	buf = bufio.NewReaderSize(buf, 20000)
	bp := BinaryParser{}
	_, err = bp.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
}
