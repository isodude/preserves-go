package text

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"slices"

	"github.com/isodude/preserves-go/lib/preserves"
	log "github.com/sirupsen/logrus"
)

type Record struct {
	preserves.Record
}

func (re *Record) New() preserves.Value {
	return &Record{Record: preserves.Record{}}
}

func (re *Record) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br *bufio.Reader
		ok bool
		s  int
		m  int64
	)

	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}

	ok, s, err = EndRune(rune(RecordRunes[0]), br)
	n += int64(s)
	if err != nil || !ok {
		return
	}

	tp := TextParser{}

	ok, s, err = EndRune(rune(RecordRunes[1]), br)
	n += int64(s)
	if err != nil || ok {
		n += 1
		return
	}

	m, err = tp.ReadFrom(br)
	n += m
	if err != nil {
		return
	}
	re.Key = tp.Result

	m, err = SkipWhitespace(br)
	n += m
	if err != nil {
		return
	}
	for {
		log.Debug("record: loop")
		ok, s, err = EndRune(rune(RecordRunes[1]), br)
		n += int64(s)
		if err != nil || ok {
			return
		}
		m, err = tp.ReadFrom(br)
		n += m
		if err != nil {
			return
		}
		re.Fields = append(re.Fields, tp.Result)
	}
}

func (re *Record) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	m, err = w.Write([]byte{RecordRunes[0]})
	n += int64(m)
	if err != nil {
		return
	}

	defer func() {
		m, err = w.Write([]byte{RecordRunes[1]})
		n += int64(m)
	}()
	k, ok := re.Key.(io.WriterTo)
	if !ok {
		return
	}
	var o int64
	o, err = k.WriteTo(w)
	n += o
	if err != nil {
		return
	}

	if len(re.Fields) == 1 {
		m, err = w.Write([]byte{' '})
		n += int64(m)
		if err != nil {
			return
		}

		f, ok := re.Fields[0].(io.WriterTo)
		if !ok {
			return
		}
		o, err = f.WriteTo(w)
		n += o
		if err != nil {
			return
		}
		return
	}
	var iw io.Writer
	iw = w
	if len(re.Fields) > 2 {
		iw = indentWriter{Writer: w}
	}

	for _, v := range re.Fields {
		if len(re.Fields) > 2 {

			m, err = w.Write([]byte{'\n', ' ', ' '})
			n += int64(m)
			if err != nil {
				return
			}
		} else {
			m, err = w.Write([]byte{' '})
			n += int64(m)
			if err != nil {
				return
			}
		}
		f, ok := v.(io.WriterTo)
		if !ok {
			return
		}
		o, err = f.WriteTo(iw)
		n += o
		if err != nil {
			return
		}
	}
	return
}

// encoding.TextMarshaler interface
func (re Record) MarshalText() (data []byte, err error) {
	buf := &bytes.Buffer{}
	_, err = re.WriteTo(buf)
	data = buf.Bytes()
	return
}

// fmt.Stringer interface
func (re Record) String() string {
	b, _ := re.MarshalText()
	return string(b)
}

type Dictionary struct {
	preserves.Dictionary
}

func (d *Dictionary) New() preserves.Value {
	return &Dictionary{Dictionary: make(preserves.Dictionary)}
}

// io.ReaderTo interface
func (d *Dictionary) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br *bufio.Reader
		ok bool
		s  int
		m  int64
		r  rune
	)

	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}

	ok, s, err = EndRune(rune(DictionaryRunes[0]), br)
	n += int64(s)
	if err != nil || !ok {
		return
	}

	tp := TextParser{SkipCommas: true}
	for {
		m, err = SkipCommas(br)
		n += m
		if err != nil {
			return
		}
		ok, s, err = EndRune(rune(DictionaryRunes[1]), br)
		n += int64(s)
		if err != nil || ok {
			return
		}
		var u, v preserves.Value
		m, err = tp.ReadFrom(br)
		n += m
		if err != nil {
			return
		}
		m, err = SkipCommas(br)
		n += m
		if err != nil {
			return
		}
		r, s, err = br.ReadRune()
		n += int64(s)
		if err != nil {
			return
		}
		if r != ':' {
			err = fmt.Errorf("parse error: missing ':' in dictionary, got %x", r)
			return
		}
		u = tp.Result

		m, err = tp.ReadFrom(br)
		n += m
		if err != nil {
			d.Dictionary[u] = nil
			return
		}
		v = tp.Result

		d.Dictionary[u] = v
	}
}

// io.WriterTo interface
func (d *Dictionary) WriteTo(w io.Writer) (n int64, err error) {
	var (
		m     int64
		s     int
		space bool
	)
	s, err = w.Write([]byte{DictionaryRunes[0], '\n', ' ', ' '})
	n += int64(s)
	if err != nil {
		return
	}

	iw := indentWriter{Writer: w}
	var keys []preserves.Value
	for k := range d.Dictionary {
		keys = append(keys, k)
	}

	slices.SortFunc[[]preserves.Value](keys, func(a, b preserves.Value) int {
		return ToPreserves(a).Cmp(ToPreserves(b))
	})
	for _, u := range keys {
		v := d.Dictionary[u]

		if !space {
			space = true
		} else {
			s, err = w.Write([]byte{'\n', ' ', ' '})
			n += int64(s)
			if err != nil {
				return
			}
		}
		m, err = u.(io.WriterTo).WriteTo(iw)
		n += m
		if err != nil {
			return
		}
		s, err = w.Write([]byte{':', ' '})
		n += int64(s)
		if err != nil {
			return
		}
		m, err = v.(io.WriterTo).WriteTo(iw)
		n += m
		if err != nil {
			return
		}
	}
	s, err = w.Write([]byte{'\n', DictionaryRunes[1]})
	n += int64(s)
	return
}

// encoding.TextMarshaler interface
func (d *Dictionary) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = d.WriteTo(b)
	return b.Bytes(), err
}

// fmt.Stringer interface
func (d *Dictionary) String() string {
	b, _ := d.MarshalText()
	return string(b)
}

type Sequence struct {
	preserves.Sequence
}

func (s *Sequence) New() preserves.Value {
	return new(Sequence)
}

func (s *Sequence) WriteTo(w io.Writer) (n int64, err error) {
	var si int
	if len(s.Sequence) == 0 {
		si, err = fmt.Fprintf(w, "%s", SequenceRunes[0:2])
		n = int64(si)
		return
	}
	si, err = fmt.Fprintf(w, "%s", []byte{SequenceRunes[0], '\n', ' ', ' '})
	n = int64(si)
	if err != nil {
		return
	}
	var (
		notfirst bool
		m        int64
	)

	iw := indentWriter{Writer: w}
	for _, v := range s.Sequence {
		if notfirst {
			si, err = w.Write([]byte{'\n', ' ', ' '})
			n += int64(si)
			if err != nil {
				return
			}
		} else {
			notfirst = true
		}
		m, err = v.WriteTo(iw)
		n += m
		if err != nil {
			return
		}
	}
	si, err = w.Write([]byte{'\n', SequenceRunes[1]})
	n += int64(si)
	return
}
func (s *Sequence) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br *bufio.Reader
		ok bool
		si int
		m  int64
	)

	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}

	ok, si, err = EndRune(rune(SequenceRunes[0]), br)
	n += int64(si)
	if err != nil || !ok {
		return
	}

	tp := TextParser{SkipCommas: true}

	log.Debug("sequence: looping")
	for {
		m, err = SkipCommas(br)
		n += m
		if err != nil {
			return
		}
		ok, si, err = EndRune(rune(SequenceRunes[1]), br)
		n += int64(si)
		if ok {
			log.Debugf("sequence: endrune2: %v\n", err)
		}
		if err == io.EOF {
			log.Debugf("sequence: eof3: %v\n", err)
		}
		if err != nil || ok {
			return
		}
		m, err = tp.ReadFrom(br)
		n += m
		if err == io.EOF {
			log.Debugf("sequence: eof4: %v\n", err)
		}
		if err != nil {
			return
		}
		s.Sequence = append(s.Sequence, tp.Result)
	}
}

func (s *Sequence) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.WriteTo(b)
	return b.Bytes(), err
}
func (s *Sequence) String() string {
	b, _ := s.MarshalText()
	return string(b)
}

type Set struct {
	preserves.Set
}

func (s *Set) New() preserves.Value {
	return &([]Set{{Set: make(preserves.Set)}}[0])

}

func (s *Set) WriteTo(w io.Writer) (n int64, err error) {
	var si int
	si, err = fmt.Fprintf(w, "%s\n  ", SetRunes[0:2])
	n = int64(si)
	if err != nil {
		return
	}
	var (
		notfirst bool
		m        int64
	)

	iw := indentWriter{Writer: w}
	var keys []preserves.Value
	for k := range s.Set {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(e1, e2 preserves.Value) int {
		return e1.Cmp(e2)
	})
	for _, v := range keys {
		if notfirst {
			si, err = w.Write([]byte{'\n', ' ', ' '})
			n += int64(si)
			if err != nil {
				return
			}
		} else {
			notfirst = true
		}
		m, err = v.WriteTo(iw)
		n += m
		if err != nil {
			return
		}
	}
	si, err = w.Write([]byte{'\n', SetRunes[2]})
	n += int64(si)
	return
}
func (s *Set) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br *bufio.Reader
		ok bool
		si int
		m  int64
	)

	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}

	bs := make([]byte, 2)
	si, err = br.Read(bs)
	n = int64(si)
	if err != nil {
		return
	}
	if string(bs) != SetRunes[0:2] {
		err = fmt.Errorf("set: no rune, found: %s", bs)
		return
	}

	tp := TextParser{}

	for {
		m, err = SkipWhitespace(br)
		n += m
		if err != nil {
			return
		}
		ok, si, err = EndRune(rune(SetRunes[2]), br)
		n += int64(si)
		if err != nil || ok {
			return
		}
		m, err = tp.ReadFrom(br)
		n += m
		if err != nil {
			return
		}

		log.Debugf("set: read: %s, %v", reflect.TypeOf(tp.Result), tp.Result)
		s.Set[tp.Result] = struct{}{}
	}
}

func (s *Set) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = s.WriteTo(b)
	return b.Bytes(), err
}
func (s *Set) String() string {
	b, _ := s.MarshalText()
	return string(b)
}

type Annotation struct {
	preserves.Annotation
}

func (a *Annotation) New() preserves.Value {
	return new(Annotation)
}
func (a *Annotation) WriteTo(w io.Writer) (n int64, err error) {
	var s int
	s, err = fmt.Fprintf(w, "%s", AnnotationRunes)
	n = int64(s)
	if err != nil {
		return
	}
	var m int64
	if a.Value != nil {
		m, err = a.Value.WriteTo(w)
		n += m
		if err != nil {
			return
		}
	} else {
		return
	}

	s, err = fmt.Fprintf(w, " ")
	n += int64(s)
	if err != nil {
		return
	}

	if a.AnnotatedValue != nil {
		m, err = a.AnnotatedValue.WriteTo(w)
		n += m
	}
	return
}
func (a *Annotation) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br *bufio.Reader
		ok bool
		s  int
		m  int64
	)

	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}

	ok, s, err = EndRune(rune(AnnotationRunes[0]), br)
	n += int64(s)
	if err != nil || !ok {
		return
	}

	tp := TextParser{}

	m, err = tp.ReadFrom(br)
	n += m
	if err == io.EOF {
		log.Debugf("embedded: eof: %v\n", err)
	}
	if err != nil {
		return
	}

	a.Value = tp.Result

	m, err = tp.ReadFrom(br)
	n += m
	if err == io.EOF {
		log.Debugf("embedded: eof: %v\n", err)
	}
	if err != nil {
		return
	}

	a.AnnotatedValue = tp.Result
	return
}
func (a *Annotation) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = a.WriteTo(b)
	return b.Bytes(), err
}

func (a *Annotation) String() string {
	b, _ := a.MarshalText()
	return string(b)
}

type Embedded struct {
	preserves.Embedded
}

func (e *Embedded) New() preserves.Value {
	return new(Embedded)
}

func (e *Embedded) WriteTo(w io.Writer) (n int64, err error) {
	var s int
	s, err = fmt.Fprintf(w, "%s", EmbeddedRunes)
	n = int64(s)
	if err != nil {
		return
	}
	var m int64
	if e.Value != nil {
		m, err = e.Value.WriteTo(w)
		n += m
	}
	return
}
func (e *Embedded) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br *bufio.Reader
		ok bool
		s  int
		m  int64
	)

	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}

	ok, s, err = EndRune(rune(EmbeddedRunes[0]), br)
	n += int64(s)
	if err != nil || !ok {
		return
	}
	ok, s, err = EndRune(rune(EmbeddedRunes[1]), br)
	n += int64(s)
	if err != nil || !ok {
		return
	}

	tp := TextParser{}

	m, err = tp.ReadFrom(br)
	n += m
	if err == io.EOF {
		log.Debugf("embedded: eof: %v\n", err)
	}
	if err != nil {
		return
	}

	e.Value = tp.Result
	return
}
func (e *Embedded) MarshalText() (data []byte, err error) {
	b := &bytes.Buffer{}
	_, err = e.WriteTo(b)
	return b.Bytes(), err
}

func (e *Embedded) String() string {
	b, _ := e.MarshalText()
	return string(b)
}

type Comment struct {
	preserves.Annotation
}

func (c *Comment) New() preserves.Value {
	return new(Comment)
}

func (c *Comment) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br     *bufio.Reader
		ok     bool
		s      int
		parsed []rune
	)
	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}
	ok, s, err = EndRune(rune(CommentRunes[0]), br)
	n += int64(s)
	if err != nil || !ok {
		return
	}

	var r rune
	for {
		r, s, err = br.ReadRune()
		n += int64(s)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		if r == '\n' || r == '\r' {
			break
		}

		parsed = append(parsed, r)
	}
	log.Debugf("comment: %s\n", string(parsed))
	c.Value = &([]String{{*preserves.NewPstring(string(parsed))}}[0])

	var m int64
	tp := TextParser{}
	m, err = tp.ReadFrom(br)
	n += m
	if err != nil {
		return
	}

	c.AnnotatedValue = tp.Result
	return
}

func (c *Comment) WriteTo(w io.Writer) (n int64, err error) {
	var s int
	var m int64

	s, err = w.Write([]byte{'\n'})
	n += int64(s)
	if err != nil {
		return
	}
	s, err = w.Write([]byte(CommentRunes))
	n += int64(s)
	if err != nil {
		return
	}
	if c.Value != nil {
		m, err = c.Value.WriteTo(w)
		n += m
		if err != nil {
			return
		}
	} else {
		return
	}

	if _, ok := c.AnnotatedValue.(*Comment); !ok {
		s, err = w.Write([]byte{'\n'})
		n += int64(s)
		if err != nil {
			return
		}
	}

	if c.AnnotatedValue != nil {
		m, err = c.AnnotatedValue.WriteTo(w)
		n += m
	}
	return
}

func (c *Comment) MarshalText() (data []byte, err error) {
	buf := &bytes.Buffer{}
	_, err = c.WriteTo(buf)
	data = buf.Bytes()
	return
}

func (c *Comment) String() string {
	b, _ := c.MarshalText()
	return string(b)
}

type Shebang struct {
	preserves.Annotation
}

func (s *Shebang) New() preserves.Value {
	return new(Shebang)
}

func (s *Shebang) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br     *bufio.Reader
		ok     bool
		si     int
		parsed []rune
	)
	if br, ok = ir.(*bufio.Reader); !ok {
		br = bufio.NewReader(ir)
	}
	bs := make([]byte, 2)
	si, err = br.Read(bs)
	n = int64(si)
	if err != nil {
		return
	}
	if string(bs) != ShebangRunes {
		err = fmt.Errorf("shebang: no runes, found: %x", bs)
		return
	}

	var r rune
	for {
		r, si, err = br.ReadRune()
		n += int64(si)
		log.Debugf("shebang: %s, %v\n", string(r), err)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		if r == '\n' || r == '\r' {
			break
		}
		parsed = append(parsed, r)
	}
	s.Value = &([]preserves.Pstring{preserves.Pstring(parsed)}[0])

	var m int64
	tp := TextParser{}
	m, err = tp.ReadFrom(br)
	n += m
	if err != nil {
		return
	}

	s.AnnotatedValue = tp.Result
	return
}

func (s *Shebang) WriteTo(w io.Writer) (n int64, err error) {
	var si int
	var m int64
	si, err = w.Write([]byte(ShebangRunes))
	n += int64(si)
	if err != nil {
		return
	}
	if s.Value != nil {
		m, err = s.Value.WriteTo(w)
		n += m
		if err != nil {
			return
		}
	} else {
		return
	}

	si, err = w.Write([]byte{'\n'})
	n += int64(si)
	if err != nil {
		return
	}

	if s.AnnotatedValue != nil {
		m, err = s.AnnotatedValue.WriteTo(w)
		n += m
	}
	return
}

func (s *Shebang) MarshalText() (data []byte, err error) {
	buf := &bytes.Buffer{}
	_, err = s.WriteTo(buf)
	data = buf.Bytes()
	return
}

func (s *Shebang) String() string {
	b, _ := s.MarshalText()
	return string(b)
}
