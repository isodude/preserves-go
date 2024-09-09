package text

import (
	"bufio"
	"io"
	"reflect"
	"regexp"
	"strings"

	"github.com/isodude/preserves-go/lib/preserves"
	log "github.com/sirupsen/logrus"
)

const (
	StringRunes           = "\"\""
	RecordRunes           = "<>"
	SequenceRunes         = "[]"
	SetRunes              = "#{}"
	DictionaryRunes       = "{}"
	AnnotationRunes       = "@"
	EmbeddedRunes         = "#:"
	QuotedSymbolRunes     = "||"
	BooleanTrueRunes      = "#t"
	BooleanFalseRunes     = "#f"
	BinByteStringRunes    = "#\"\""
	HexByteStringRunes    = "#x\"\""
	Base64ByteStringRunes = "#[]"
	HexDoubleRunes        = "#xd\"\""
	CommentRunes          = "#"
	WhiteSpaceRunes       = " \t\r\n"
	ShebangRunes          = "#!"
)

type Maze map[rune]*MazePath

func NewMaze(m map[rune]*MazePath) *Maze {
	n := Maze(m)
	return &n
}

func (rm *Maze) Get(r rune) (mp *MazePath, ok bool) {
	mp, ok = (*rm)[r]
	return
}

type MazePath struct {
	ReadFrom io.ReaderFrom
	Maze     *Maze
}

func (mp *MazePath) AddPath(path string, rf io.ReaderFrom) *MazePath {
	if len(path) == 0 {
		mp.ReadFrom = rf
		return mp
	}
	if mp.Maze == nil {
		mp.Maze = &Maze{}
	}
	m, ok := (*mp.Maze)[rune(path[0])]
	if ok {
		m.AddPath(path[1:], rf)
		return mp
	}
	newMp := (&MazePath{}).AddPath(path[1:], rf)
	(*mp.Maze)[rune(path[0])] = newMp
	return mp
}

var RuneMaze = (&MazePath{}).
	// Bare Atoms
	AddPath("", new(Bare)).

	// Atom
	AddPath(string(QuotedSymbolRunes[0]), new(QuotedSymbol)).
	AddPath(string(StringRunes[0]), new(String)).
	AddPath(BinByteStringRunes[0:2], new(BinByteString)).
	AddPath(Base64ByteStringRunes[0:2], new(Base64ByteString)).
	AddPath(BooleanTrueRunes, new(Boolean)).
	AddPath(BooleanFalseRunes, new(Boolean)).
	AddPath(HexDoubleRunes[0:4], new(HexDouble)).
	AddPath(HexByteStringRunes[0:3], new(HexByteString)).

	// Compound
	AddPath(SetRunes[0:2], new(Set)).
	AddPath(string(RecordRunes[0]), new(Record)).
	AddPath(string(DictionaryRunes[0]), new(Dictionary)).
	AddPath(string(SequenceRunes[0]), new(Sequence)).

	// Annotation
	AddPath(AnnotationRunes, new(Annotation)).
	AddPath(CommentRunes, new(Comment)).
	AddPath(ShebangRunes, new(Shebang)).

	// Embedded
	AddPath(EmbeddedRunes, new(Embedded))

/*
	MazePath{
		Maze: NewMaze(map[rune]MazePath{
			rune(QuotedSymbolRunes[0]): {ReadFrom: new(QuotedSymbol)},
			rune(StringRunes[0]):       {ReadFrom: new(String)},
			rune(RecordRunes[0]):       {ReadFrom: new(Record)},
			rune(DictionaryRunes[0]):   {ReadFrom: new(Dictionary)},
			rune(SequenceRunes[0]):     {ReadFrom: new(Sequence)},
			rune(AnnotationRunes[0]):   {ReadFrom: new(Annotation)},
			rune(CommentRunes[0]): {
				ReadFrom: new(Comment),
				Maze: NewMaze(map[rune]MazePath{
					rune(ShebangRunes[1]):          {ReadFrom: new(Shebang)},
					rune(BinByteStringRunes[1]):    {ReadFrom: new(BinByteString)},
					rune(Base64ByteStringRunes[1]): {ReadFrom: new(Base64ByteString)},
					rune(BooleanTrueRunes[1]):      {ReadFrom: new(Boolean)},
					rune(BooleanFalseRunes[1]):     {ReadFrom: new(Boolean)},
					rune(EmbeddedRunes[1]):         {ReadFrom: new(Embedded)},
					rune(SetRunes[1]):              {ReadFrom: new(Set)},
					rune(HexDoubleRunes[1]): {
						Maze: NewMaze(map[rune]MazePath{
							rune(HexByteStringRunes[2]): {ReadFrom: new(HexByteString)},
							rune(HexDoubleRunes[2]): {
								Maze: NewMaze(map[rune]MazePath{
									rune(HexDoubleRunes[3]): {ReadFrom: new(HexDouble)},
								}),
							},
						}),
					},
				}),
			},
		}),
	}
*/
var (
	WhiteSpaceRegexp     = regexp.MustCompile("[ \t\r\n]")
	CommasRegexp         = regexp.MustCompile("([ \t\r\n]*,)*[ \t\r\n]*")
	SymbolOrNumberRegexp = regexp.MustCompile(`^[-a-zA-Z0-9~!$%^&*?_=+/.]+$`)
	NumberRegexp         = regexp.MustCompile(`^([-+]?\d+)((\.\d+([eE][-+]?\d+)?)|([eE][-+]?\d+))?$`)
	DoubleRegexp         = regexp.MustCompile(`^([-+]?\d+)((\.\d+([eE][-+]?\d+)?)|([eE][-+]?\d+))$`)
	SignedIntegerRegexp  = regexp.MustCompile(`^([-+]?\d+)$`)
	Base64Regexp         = regexp.MustCompile(`[a-zA-Z0-9+/_=']`)
)

func EndRune(er rune, ir io.Reader) (ok bool, n int, err error) {
	var (
		r rune
		s int
	)
	br := ir.(*bufio.Reader)

	r, s, err = br.ReadRune()
	n += s

	log.Debugf("endrune: %s, %v", string(r), err)

	if err == io.EOF && r == er {
		ok = true
		err = nil
		return
	}

	if err != nil {
		return
	}

	if r == er {
		ok = true
		return
	}

	err = br.UnreadRune()
	n -= s
	return
}

func SkipWhitespace(ir io.Reader) (n int64, err error) {
	var (
		r rune
		s int
	)
	br := ir.(*bufio.Reader)
	r, s, err = br.ReadRune()
	n += int64(s)
	log.Debugf("SkipWhitespace: %s, %v", string(r), err)
	if err != nil {
		return
	}
	for r == ' ' || r == '\t' || r == '\n' || r == '\r' {
		r, s, err = br.ReadRune()
		n += int64(s)
		log.Debugf("SkipWhitespace: loop: %s, %v", string(r), err)
		if err != nil {
			return
		}
	}
	err = br.UnreadRune()
	n -= int64(s)
	return
}

func SkipCommas(ir io.Reader) (n int64, err error) {
	var (
		r rune
		s int
	)
	br := ir.(*bufio.Reader)
	r, s, err = br.ReadRune()
	n += int64(s)
	log.Debugf("SkipCommas: %s, %v", string(r), err)
	if err != nil {
		return
	}
	for r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == ',' {
		r, s, err = br.ReadRune()
		n += int64(s)
		log.Debugf("SkipCommas: loop: %s, %v", string(r), err)
		if err != nil {
			return
		}
	}
	err = br.UnreadRune()
	n -= int64(s)
	return
}

type indentWriter struct {
	io.Writer
}

func (iw indentWriter) Write(p []byte) (int, error) {
	p = []byte(strings.ReplaceAll(string(p), "\n", "\n  "))
	return iw.Writer.Write(p)
}

type Bare struct {
	*bufio.Reader
	Error error
}

func (b *Bare) New() preserves.Value {
	//log.Debugf("bare-new")
	var (
		bs  []byte
		err error
		i   int
		str []byte
	)
	for {
		i++
		bs, err = b.Reader.Peek(i)
		b.Error = err

		if err != nil {
			if err != io.EOF {
				return nil
			}
		}
		if err == io.EOF || !SymbolOrNumberRegexp.Match(bs[i-1:]) {
			if SignedIntegerRegexp.Match(str) {
				return new(SignedInteger)
			}
			if DoubleRegexp.Match(str) {
				return new(BareDouble)
			}
			return new(BareSymbol)
		}
		str = append(str, bs[i-1])
	}
}

func (b *Bare) ReadFrom(_ io.Reader) (n int64, err error) { return }

type TextParser struct {
	Result     preserves.Value
	SkipCommas bool
	RuneMaze   *MazePath
}

func (tp *TextParser) GetValue() preserves.Value {
	return tp.Result
}

func (tp *TextParser) ReadFrom(ir io.Reader) (n int64, err error) {
	var (
		br *bufio.Reader
		ok bool
		i  int
		m  int64
		bs []byte
		mp *MazePath
	)

	if tp.RuneMaze != nil {
		mp = tp.RuneMaze
	} else {
		mp = RuneMaze
	}

	maze := mp.Maze
	br, ok = ir.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(ir)
	}
	if b, ok := mp.ReadFrom.(*Bare); ok {
		b.Reader = br
	}
	value := mp.ReadFrom

	if tp.SkipCommas {
		m, err = SkipCommas(br)
	} else {
		m, err = SkipWhitespace(br)
	}
	n += int64(m)
	if err != nil {
		return
	}

	for maze != nil {
		i++
		bs, err = br.Peek(i)
		if err != nil {
			return
		}

		if mp.ReadFrom != nil {
			value = mp.ReadFrom
		}

		mp, ok = maze.Get(rune(bs[len(bs)-1]))
		if ok {
			maze = mp.Maze
			if mp.ReadFrom != nil {
				value = mp.ReadFrom
			}
			continue
		}
		maze = nil
	}

	tp.Result = value.(preserves.New).New()
	if tp.Result == nil {
		//err = fmt.Errorf("textparser: %v", value.(*Bare).Error)
		return
	}
	log.Debugf("textparser: found: %s", reflect.TypeOf(tp.Result))
	m, err = tp.Result.(io.ReaderFrom).ReadFrom(br)
	n += m
	return
}

func ToPreserves(v preserves.Value) preserves.Value {
	switch c := v.(type) {
	case *BareSymbol:
		return &c.Symbol
	case *QuotedSymbol:
		return &c.Pstring
	case *preserves.Symbol:
		return c
	case *Boolean:
		return &c.Boolean
	case *SignedInteger:
		return &c.SignedInteger
	case *HexByteString:
		return &c.ByteString
	case *BinByteString:
		return &c.ByteString
	case *Base64ByteString:
		return &c.ByteString
	case *String:
		return &c.Pstring
	case *preserves.Pstring:
		return c
	case *Record:
		r := &preserves.Record{}
		r.Key = ToPreserves(c.Key)
		for _, v := range c.Fields {
			r.Fields = append(r.Fields, ToPreserves(v))
		}
		return r
	case *preserves.Record:
		r := &preserves.Record{}
		r.Key = ToPreserves(c.Key)
		for _, v := range c.Fields {
			r.Fields = append(r.Fields, ToPreserves(v))
		}
		return r
	case *Set:
		s := make(preserves.Set)
		for k := range c.Set {
			s[ToPreserves(k)] = struct{}{}
		}
		return &s
	case *Sequence:
		s := &preserves.Sequence{}
		for _, v := range c.Sequence {
			*s = append(*s, ToPreserves(v))
		}
		return s
	case *Dictionary:
		d := make(preserves.Dictionary)
		for k, v := range c.Dictionary {
			d[ToPreserves(k)] = ToPreserves(v)
		}
		return &d
	case *Annotation:
		a := &preserves.Annotation{}
		a.Value = ToPreserves(c.Value)
		a.AnnotatedValue = ToPreserves(c.AnnotatedValue)
		return a
	case *preserves.Annotation:
		a := &preserves.Annotation{}
		a.Value = ToPreserves(c.Value)
		a.AnnotatedValue = ToPreserves(c.AnnotatedValue)
		return a
	case *Shebang:
		a := &preserves.Comment{}
		a.Value = *ToPreserves(preserves.Value(&a.Value)).(*preserves.Pstring)
		a.AnnotatedValue = ToPreserves(c.AnnotatedValue)
		return a
	case *Comment:
		a := &preserves.Comment{}
		a.Value = *ToPreserves(preserves.Value(&a.Value)).(*preserves.Pstring)
		a.AnnotatedValue = ToPreserves(c.AnnotatedValue)
		return a
	case *Embedded:
		e := &preserves.Embedded{}
		e.Value = ToPreserves(c.Value)
		return e
	}
	log.Fatalf("unknown type(2): %v", reflect.TypeOf(v))
	return nil
}

func FromPreserves(v preserves.Value) preserves.Value {
	switch c := v.(type) {
	case *preserves.Symbol:
		if strings.Contains(c.String(), " ") {
			return &QuotedSymbol{Pstring: preserves.Pstring(c.String())}
		}
		return &BareSymbol{Symbol: *c}
	case *preserves.Pstring:
		return &String{Pstring: *c}
	case *preserves.SignedInteger:
		return &SignedInteger{SignedInteger: *c}
	case *preserves.Boolean:
		return &Boolean{Boolean: *c}
	case *preserves.ByteString:
		return &Base64ByteString{ByteString: *c}
	case *preserves.Dictionary:
		d := &Dictionary{Dictionary: make(preserves.Dictionary)}
		for k, v := range *c {
			d.Dictionary[FromPreserves(k)] = FromPreserves(v)
		}
		return d
	case *preserves.Sequence:
		s := &Sequence{}
		for _, v := range *c {
			s.Sequence = append(s.Sequence, FromPreserves(v))
		}
		return s
	case *Sequence:
		s := &Sequence{}
		for _, v := range c.Sequence {
			s.Sequence = append(s.Sequence, FromPreserves(v))
		}
		return s
	case *preserves.Set:
		s := &Set{Set: make(preserves.Set)}
		for k := range *c {
			s.Set[FromPreserves(k)] = struct{}{}
		}
		return s
	case *preserves.Record:
		r := &Record{}
		r.Key = FromPreserves(c.Key)
		for _, v := range c.Fields {
			r.Fields = append(r.Fields, FromPreserves(v))
		}
		return r
	case *preserves.Annotation:
		return &Annotation{
			Annotation: preserves.Annotation{
				Value:          FromPreserves(c.Value),
				AnnotatedValue: FromPreserves(c.AnnotatedValue),
			},
		}
	case *preserves.Embedded:
		e := &Embedded{}
		e.Value = FromPreserves(c.Value)
		return e
	case nil:
		return c
	case *Record:
		r := &Record{}
		r.Key = FromPreserves(c.Key)
		for _, v := range c.Fields {
			r.Fields = append(r.Fields, FromPreserves(v))
		}
		return r
	case *BareDouble:
		return c
	case *BareSymbol:
		return c
	case *Annotation:
		return &Annotation{
			Annotation: preserves.Annotation{
				Value:          FromPreserves(c.Value),
				AnnotatedValue: FromPreserves(c.AnnotatedValue),
			},
		}
	case *Dictionary:
		d := &Dictionary{Dictionary: make(preserves.Dictionary)}
		for k, v := range c.Dictionary {
			d.Dictionary[FromPreserves(k)] = FromPreserves(v)
		}
		return d
	}
	log.Fatalf("unknown type(3): %v", reflect.TypeOf(v))
	return nil
}
