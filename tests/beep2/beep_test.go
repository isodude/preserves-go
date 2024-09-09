package beep

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/isodude/preserves-go/lib/preserves/text"
	"github.com/kylelemons/godebug/diff"
)

func TestC(t *testing.T) {
	file := "../schema.prs.ast"
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
	s := SchemaFromPreserves(text.ToPreserves(tp.GetValue()))
	if s == nil {
		t.Fatal("s was nil")
	}
	r := SchemaToPreserves(*s)
	p := text.FromPreserves(r)
	var b bytes.Buffer
	p.WriteTo(&b)
	f, err = os.Open("../../tests/schema.prs.ast")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	bf, err := io.ReadAll(f)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(bf)) != strings.TrimSpace(b.String()) {
		//fmt.Printf("%s\n", s)
		t.Fatalf("%v", diff.Diff(b.String(), string(bf)))
	}
	//fmt.Printf("%s\n", b.String())
	//pp.Default.SetColoringEnabled(false)
	//pp.Print(r)
	//t.Fatal("beep")
}
