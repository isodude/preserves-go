package binary

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"slices"

	"github.com/isodude/preserves-go/lib/extras"
	"github.com/isodude/preserves-go/lib/preserves"
	log "github.com/sirupsen/logrus"
)

type Record struct {
	preserves.Record
}

func (r *Record) New() preserves.Value {
	return new(Record)
}

func (r *Record) Equal(y preserves.Value) bool {
	return r.Record.Equal(ToPreserves(y))
}
func (r *Record) Cmp(y preserves.Value) int {
	return r.Record.Cmp(ToPreserves(y))
}
func (r *Record) WriteTo(w io.Writer) (n int64, err error) {
	var (
		m int64
		s int
	)
	s, err = w.Write([]byte{RecordByte})
	n = int64(s)
	if err != nil {
		return
	}
	defer func() {
		s, err = w.Write([]byte{EndByte})
		n += int64(s)
	}()
	if r.Record.Key == nil {
		return
	}
	m, err = r.Record.Key.WriteTo(w)
	n += m
	if err != nil {
		return
	}

	for _, v := range r.Record.Fields {
		m, err = v.WriteTo(w)
		n += int64(m)
		if err != nil {
			return
		}
	}
	return
}
func (r *Record) ReadFrom(rd io.Reader) (n int64, err error) {
	br, ok := rd.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(rd)
	}
	var b byte
	b, err = br.ReadByte()
	if err != nil {
		return
	}
	if b != RecordByte {
		err = fmt.Errorf("record: no rune, found (%x)", b)
		return
	}
	var (
		m int64
		s int
	)

	ok, s, err = ReadEndByte(br)
	n += int64(s)
	if err != nil || ok {
		return
	}
	bp := BinaryParser{}
	m, err = bp.ReadFrom(br)
	n += m
	if err != nil {
		return
	}
	r.Record.Key = bp.Result
	r.Record.Fields = []preserves.Value{}
	for {
		ok, s, err = ReadEndByte(br)
		n += int64(s)
		if err != nil || ok {
			return
		}
		m, err = bp.ReadFrom(br)

		log.Debugf("record: loop: readfrom: %d, %v", m, err)
		n += m
		if err != nil {
			return
		}
		r.Record.Fields = append(r.Record.Fields, bp.Result)
	}
}
func (r *Record) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = r.WriteTo(b)
	return b.Bytes(), err
}

type Sequence struct {
	preserves.Sequence
}

func NewSequence(v ...preserves.Value) *Sequence {
	s := Sequence{}
	for _, u := range v {
		s.Sequence = append(s.Sequence, u)
	}
	return &s
}

func (s *Sequence) Equal(y preserves.Value) bool {
	return s.Sequence.Equal(ToPreserves(y))
}
func (s *Sequence) Cmp(y preserves.Value) int {
	return s.Sequence.Cmp(ToPreserves(y))
}
func (s *Sequence) New() preserves.Value {
	return new(Sequence)
}
func (s *Sequence) WriteTo(w io.Writer) (n int64, err error) {
	var (
		m  int64
		si int
	)
	si, err = w.Write([]byte{SequenceByte})
	n += int64(si)
	if err != nil {
		return
	}
	for _, v := range s.Sequence {
		m, err = v.WriteTo(w)
		n += m
		if err != nil {
			return
		}
	}
	si, err = w.Write([]byte{EndByte})
	n += int64(si)
	return
}
func (s *Sequence) ReadFrom(rd io.Reader) (n int64, err error) {
	br, ok := rd.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(rd)
	}
	var b byte

	log.Debugf("sequence: ")
	b, err = br.ReadByte()
	if err != nil {
		return
	}
	if b != SequenceByte {
		err = fmt.Errorf("sequence: no rune, found (%x)", b)
		return
	}

	var (
		m  int64
		si int
	)
	bp := BinaryParser{}
	s.Sequence = []preserves.Value{}
	for {
		ok, si, err = ReadEndByte(br)
		n += int64(si)
		if err != nil || ok {
			return
		}
		m, err = bp.ReadFrom(br)
		n += m
		if err != nil {
			return
		}
		s.Sequence = append(s.Sequence, bp.Result)
	}
}

func (s *Sequence) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.WriteTo(b)
	return b.Bytes(), err
}

type Set struct {
	preserves.Set
}

func NewSet(v ...preserves.Value) *Set {
	s := Set{}
	for _, u := range v {
		s.Set[u] = struct{}{}
	}
	return &s
}

func (s *Set) New() preserves.Value {
	return extras.Reference(Set{preserves.Set(map[preserves.Value]struct{}{})})
}

func (s *Set) Equal(y preserves.Value) bool {
	return s.Set.Equal(ToPreserves(y))
}
func (s *Set) Cmp(y preserves.Value) int {
	return s.Set.Cmp(ToPreserves(y))
}
func (s *Set) WriteTo(w io.Writer) (n int64, err error) {
	var (
		m  int64
		si int
	)
	si, err = w.Write([]byte{SetByte})
	n = int64(si)
	if err != nil {
		return
	}
	var keys []preserves.Value
	for k := range s.Set {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(e1, e2 preserves.Value) int {
		return e1.Cmp(e2)
	})
	for _, v := range keys {
		m, err = v.WriteTo(w)
		n += m
		if err != nil {

			return
		}
	}

	si, err = w.Write([]byte{EndByte})
	n = int64(si)
	return
}
func (s *Set) ReadFrom(rd io.Reader) (n int64, err error) {
	br, ok := rd.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(rd)
	}
	var b byte
	b, err = br.ReadByte()
	if err != nil {
		return
	}
	if b != SetByte {
		err = fmt.Errorf("set: no rune, found (%x)", b)
		return
	}

	var (
		m  int64
		si int
	)

	ok, si, err = ReadEndByte(br)
	n += int64(si)
	if err != nil || ok {
		return
	}
	bp := BinaryParser{}
	s.Set = make(map[preserves.Value]struct{})
	for {
		ok, si, err = ReadEndByte(br)
		n += int64(si)
		if err != nil || ok {
			return
		}
		m, err = bp.ReadFrom(br)
		n += m
		if err != nil {
			return
		}
		s.Set[bp.Result] = struct{}{}
	}
}
func (s *Set) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.WriteTo(b)
	return b.Bytes(), err
}

type Dictionary struct {
	preserves.Dictionary
}

func (d *Dictionary) New() preserves.Value {
	return extras.Reference(Dictionary{preserves.Dictionary(map[preserves.Value]preserves.Value{})})
}

func NewDictionary(pairs ...struct{ key, field preserves.Value }) *Dictionary {
	v := Dictionary{preserves.Dictionary(map[preserves.Value]preserves.Value{})}
	for _, k := range pairs {
		v.Dictionary[k.key] = k.field
	}
	return &v
}

func (d *Dictionary) Equal(y preserves.Value) bool {
	return d.Dictionary.Equal(ToPreserves(y))
}
func (d *Dictionary) Cmp(y preserves.Value) int {
	return d.Dictionary.Cmp(ToPreserves(y))
}
func (d *Dictionary) WriteTo(w io.Writer) (n int64, err error) {
	var (
		m int64
		s int
	)
	s, err = w.Write([]byte{DictionaryByte})
	n = int64(s)
	if err != nil {
		return
	}
	var keys []preserves.Value
	for k := range d.Dictionary {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(e1, e2 preserves.Value) int {
		return e1.Cmp(e2)
	})
	for _, u := range keys {
		v := d.Dictionary[u]
		m, err = u.WriteTo(w)
		n += m
		if err != nil {
			return
		}
		m, err = v.WriteTo(w)
		n += m
		if err != nil {
			return
		}
	}

	s, err = w.Write([]byte{EndByte})
	n = int64(s)
	return
}
func (d *Dictionary) ReadFrom(rd io.Reader) (n int64, err error) {
	br, ok := rd.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(rd)
	}
	var b byte
	b, err = br.ReadByte()
	if err != nil {
		return
	}
	if b != DictionaryByte {
		err = fmt.Errorf("dictionary: no rune, found (%x)", b)
		return
	}

	var (
		m int64
		s int
		u preserves.Value
	)

	ok, s, err = ReadEndByte(br)
	n += int64(s)
	if err != nil || ok {
		return
	}
	bp := BinaryParser{}
	d.Dictionary = make(map[preserves.Value]preserves.Value)
	for {
		ok, s, err = ReadEndByte(br)
		n += int64(s)
		if err != nil || ok {
			return
		}
		m, err = bp.ReadFrom(br)
		n += m
		if err != nil {
			return
		}
		u = bp.Result
		ok, s, err = ReadEndByte(br)
		n += int64(s)
		if err != nil || ok {
			return
		}
		m, err = bp.ReadFrom(br)
		n += m
		if err != nil {
			return
		}
		d.Dictionary[u] = bp.Result
	}
}
func (d *Dictionary) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = d.WriteTo(b)
	return b.Bytes(), err
}

type Annotation struct {
	preserves.Annotation
}

func NewAnnotation(value, annotatedValue preserves.Value) *Annotation {
	return &([]Annotation{{Annotation: preserves.Annotation{Value: value, AnnotatedValue: annotatedValue}}}[0])
}

func (a *Annotation) New() preserves.Value {
	return new(Annotation)
}

func (a *Annotation) Equal(y preserves.Value) bool {
	return a.Annotation.Equal(ToPreserves(y))
}
func (a *Annotation) Cmp(y preserves.Value) int {
	return a.Annotation.Cmp(ToPreserves(y))
}
func (a *Annotation) WriteTo(w io.Writer) (n int64, err error) {
	var (
		m int64
		s int
	)
	s, err = w.Write([]byte{AnnotationByte})
	n = int64(s)
	if err != nil {
		return
	}
	m, err = a.Annotation.Value.WriteTo(w)
	n += m
	if err != nil {
		return
	}
	if a.Annotation.AnnotatedValue != nil {
		m, err = a.Annotation.AnnotatedValue.WriteTo(w)
		n += m
	}
	return
}

func (a *Annotation) ReadFrom(rd io.Reader) (n int64, err error) {
	br, ok := rd.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(rd)
	}
	var b byte
	b, err = br.ReadByte()
	if err != nil {
		return
	}
	if b != AnnotationByte {
		err = fmt.Errorf("annotation: no rune, found (%x)", b)
		return
	}

	var (
		bs []byte
		m  int64
	)

	bp := BinaryParser{}
	log.Debugf("annotation: read value")
	m, err = bp.ReadFrom(br)
	n += m
	if err != nil {
		return
	}
	a.Annotation.Value = bp.Result

	bs, err = br.Peek(1)
	if err != nil {
		return
	}
	log.Debugf("annotation: read annotatedvalue: %x", bs)
	m, err = bp.ReadFrom(br)
	n += m
	if err != nil {
		return
	}
	a.Annotation.AnnotatedValue = bp.Result
	return
}
func (a *Annotation) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = a.WriteTo(b)
	return b.Bytes(), err
}

type Embedded struct {
	preserves.Embedded
}

func NewEmbedded(v preserves.Value) *Embedded {
	return &([]Embedded{{preserves.Embedded{Value: v}}}[0])
}

func (e *Embedded) New() preserves.Value {
	return new(Embedded)
}

func (e *Embedded) Equal(y preserves.Value) bool {
	return e.Embedded.Equal(ToPreserves(y))
}
func (e *Embedded) Cmp(y preserves.Value) int {
	return e.Embedded.Cmp(ToPreserves(y))
}
func (e *Embedded) WriteTo(w io.Writer) (n int64, err error) {
	var (
		m int64
		s int
	)
	s, err = w.Write([]byte{EmbeddedByte})
	n = int64(s)
	if err != nil {
		return
	}
	m, err = e.Embedded.Value.WriteTo(w)
	n += m
	return
}
func (e *Embedded) ReadFrom(rd io.Reader) (n int64, err error) {
	br, ok := rd.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(rd)
	}
	var b byte
	b, err = br.ReadByte()
	if err != nil {
		return
	}
	if b != EmbeddedByte {
		err = fmt.Errorf("embedded: no rune, found (%x)", b)
		return
	}

	var (
		m int64
	)

	bp := BinaryParser{}
	m, err = bp.ReadFrom(br)
	n += m
	if err != nil {
		return
	}
	e.Embedded.Value = bp.Result
	return
}
func (e *Embedded) MarshalBinary() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = e.WriteTo(b)
	return b.Bytes(), err
}
