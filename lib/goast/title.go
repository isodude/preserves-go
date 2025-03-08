package goast

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Title interface {
	First(*title) string
	SetName(string)
	GetName() string
	Title() *title
	PrefixTitle(string) *title
	GetTitle() string
	GetPrefixTitle(*title) string
	GetFirstLower() string
	Ident(*title, bool) *ast.Ident
}
type title struct {
	name string
}

func (t *title) SetName(name string) {
	t.name = name
}

func (t *title) GetName() string {
	return t.name
}

func (t *title) Title() *title {
	return t
}
func (t *title) PrefixTitle(prefix string) *title {
	return &title{name: fmt.Sprintf("%s%s", prefix, cases.Title(language.English, cases.NoLower).String(t.name))}
}

func (t *title) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(t.name)
}

func (t *title) GetPrefixTitle(prefix *title) string {
	if prefix != nil {
		return fmt.Sprintf("%s%s", prefix.GetTitle(), t.GetTitle())
	}
	return t.GetTitle()
}

func (t *title) GetFirstLower() string {
	if len(t.name) > 0 {
		return fmt.Sprintf("%s%s", strings.ToLower(t.name[0:1]), t.name[1:])
	}
	return ""
}

func (t *title) Ident(prefix *title, lower bool) *ast.Ident {
	if lower {
		if prefix != nil {
			return ast.NewIdent(fmt.Sprintf("%s%s", prefix.GetTitle(), t.GetFirstLower()))
		}
		return ast.NewIdent(t.GetFirstLower())
	}
	return ast.NewIdent(t.GetPrefixTitle(prefix))
}

func (t *title) First(prefix *title) string {
	n := t.GetPrefixTitle(prefix)
	if len(n) > 0 {
		return strings.ToLower(string(n[0]))
	}
	return ""
}
