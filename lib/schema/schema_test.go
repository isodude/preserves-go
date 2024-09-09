package schema

import (
	"bytes"
	"encoding"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/isodude/preserves-go/lib/goast"
	. "github.com/isodude/preserves-go/lib/preserves"
	"github.com/isodude/preserves-go/lib/preserves/text"
	"github.com/kylelemons/godebug/diff"
)

func TestA(t *testing.T) {
	v, err := FromPreservesSchemaFile("../../tests/schema.prs")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatalf("v was nil")
	}

	_s := SchemaToPreserves(*v)

	__s, err := text.FromPreserves(_s).(encoding.TextMarshaler).MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(__s))

	s := SchemaToPreservesSchema(*v, "")

	f, err := os.Open("../../tests/schema.prs.generated")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(b)) != strings.TrimSpace(s) {
		fmt.Printf("%s\n", s)
		t.Fatalf("%v", diff.Diff(s, string(b)))
	}
}
func TestB(t *testing.T) {
	v, err := FromPreservesSchemaFile("../../tests/schema.prs.generated")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatalf("v was nil")
	}
	s := SchemaToPreserves(*v)

	m, err := text.FromPreserves(s).(encoding.TextMarshaler).MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("../../tests/schema.prs.ast")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(b)) != strings.TrimSpace(string(m)) {
		fmt.Printf("%s\n", s)
		t.Fatalf("%v", diff.Diff(string(m), string(b)))
	}
}
func TestC(t *testing.T) {
	v, err := FromPreserves("../../tests/schema.prs.ast")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatalf("v was nil")
	}
	s := SchemaToPreserves(*v)
	m, err := text.FromPreserves(s).(encoding.TextMarshaler).MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("../../tests/schema.prs.ast")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(b)) != strings.TrimSpace(string(m)) {
		t.Fatalf("%v", diff.Diff(string(m), string(b)))
	}
}
func TestD(t *testing.T) {
	v, err := FromPreservesSchemaFile("../../tests/schema.prs")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatalf("v was nil")
	}
	s := SchemaToPreserves(*v)

	m, err := text.FromPreserves(s).(encoding.TextMarshaler).MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("../../tests/schema.prs.ast")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(b)) != strings.TrimSpace(string(m)) {
		fmt.Printf("%s\n", s)
		t.Fatalf("%v", diff.Diff(string(m), string(b)))
	}
}
func TestE(t *testing.T) {
	v, err := FromPreserves("../../tests/schema.prs.ast")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatalf("v was nil")
	}
	fmt.Printf("%s\n", goast.EncodeToGoAST("schema", v))
	t.Fatal("beep")
}
func TestF(t *testing.T) {
	v, err := FromPreservesSchemaFile("../../tests/example.prs")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatalf("v was nil")
	}
	fmt.Printf("%s\n", goast.EncodeToGoAST("schema", v))
	fmt.Printf("%s\n", SchemaToPreservesSchema(*v, ""))
	t.Fatal("example")
}
func TestG(t *testing.T) {
	v, err := FromPreserves("../../tests/schema.prs.ast")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatalf("v was nil")
	}
	fmt.Printf("%s\n", goast.EncodeToGoAST("schema", v))
	t.Fatal("beep")
}
func TestDump(t *testing.T) {
	fset := token.NewFileSet()

	// Read in the original file
	snippet, err := parser.ParseFile(fset, "named_pattern.go", nil, parser.Trace|parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	var bytes bytes.Buffer
	cfg := &printer.Config{Mode: printer.UseSpaces}
	err = cfg.Fprint(&bytes, fset, snippet)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", bytes.String())
	t.Fatalf("beep")
}

type st[T any] struct {
	f func(Value) T
	a Value
	b Value
	c Value
}
