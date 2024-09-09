package beep

import (
	"bufio"
	"io"
	"os"
	"testing"

	"github.com/isodude/preserves-go/lib/preserves/text"
	"github.com/k0kubun/pp/v3"
)

func TestC(t *testing.T) {
	file := "schema.prs.ast"
	f, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	buf := bufio.NewReader(f)
	buf = bufio.NewReaderSize(buf, 200000)
	tp := text.TextParser{}
	_, err = tp.ReadFrom(buf)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	s := (&Schema{}).FromPreserves(text.ToPreserves(tp.GetValue()))
	if s == nil {
		t.Fatal("s was nil")
	}

	pp.Default.SetColoringEnabled(false)
	pp.Print(s)
	t.Fatal("beep")
}
