package goast

import (
	"go/ast"
	"reflect"

	"github.com/isodude/preserves-go/lib/preserves"
)

type Match interface {
	Stmt() []ast.Stmt
	FromPreserves(preserves.Value) any
}

/*
	type LengthRangeMatch struct {
		From int
		To   int
		M    Match
	}

	func NewLengthRangeMatch(from int, to int, m Match) *LengthRangeMatch {
		return &LengthRangeMatch{From: from, To: to, M: m}
	}

	func (l *LengthRangeMatch) FromPreserves(values []preserves.Value) any {
		if len(values) > l.From && len(values) < l.To {
			return l.M.FromPreserves(values)
		}

		return nil
	}

	func (l *LengthRangeMatch) Stmt() (stmts []ast.Stmt) {
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X: &ast.CallExpr{
						Fun:  ast.NewIdent("len"),
						Args: []ast.Expr{ast.NewIdent("value")}},
					Op: token.GTR,
					Y:  ast.NewIdent(strconv.Itoa(l.From)),
				},
				Op: token.LAND,
				Y: &ast.BinaryExpr{
					X: &ast.CallExpr{
						Fun:  ast.NewIdent("len"),
						Args: []ast.Expr{ast.NewIdent("values")}},
					Op: token.LSS,
					Y:  ast.NewIdent(strconv.Itoa(l.To)),
				},
			},
			Body: &ast.BlockStmt{List: l.M.Stmt()},
		})
		return
	}

	type LengthExactMatch struct {
		L int
		M Match
	}

	func NewLengthExactMatch(l int, m Match) *LengthExactMatch {
		return &LengthExactMatch{L: l, M: m}
	}

	func (l *LengthExactMatch) FromPreserves(values preserves.Value) any {
		if len(values) == l.L {
			return l.M.FromPreserves(values)
		}

		return nil
	}

	func (l *LengthExactMatch) Stmt() (stmts []ast.Stmt) {
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun:  ast.NewIdent("len"),
					Args: []ast.Expr{ast.NewIdent("values")}},
				Op: token.EQL,
				Y:  ast.NewIdent(strconv.Itoa(l.From)),
			},
			Body: &ast.BlockStmt{List: l.M.Stmt()},
		})

		return
	}
*/
/*
func (r *Rec) Matcher() Match {
	return NewAssertTypeMatch(&preserves.Record{}, )
}
*/
type BoolMatch interface {
	Match(preserves.Value) bool
}
type ListMatch struct {
	B []BoolMatch
	M Match
}

func NewListMatch(m Match, b ...BoolMatch) *ListMatch {
	return &ListMatch{B: b, M: m}
}
func (l *ListMatch) FromPreserves(value preserves.Value) any {
	for _, k := range l.B {
		if !k.Match(value) {
			return nil
		}
	}
	return l.M.FromPreserves(value)
}
func (l *ListMatch) Stmt() (stmts []ast.Stmt) { return }

type AssertTypeMatch struct {
	Value preserves.Value
	M     Match
}

func NewAssertTypeMatch(value preserves.Value, m Match) *AssertTypeMatch {
	return &AssertTypeMatch{Value: value, M: m}
}
func (a *AssertTypeMatch) Match(value preserves.Value) bool {
	return reflect.TypeOf(value) == reflect.TypeOf(a.Value)
}
func (a *AssertTypeMatch) Stmt() (stmts []ast.Stmt) { return }

type RecordKeyBoolMatch string

func NewRecordKeyBoolMatch(s string) *RecordKeyBoolMatch {
	return &([]RecordKeyBoolMatch{RecordKeyBoolMatch(s)}[0])
}

func (r *RecordKeyBoolMatch) Match(value preserves.Value) bool {
	if rec, ok := value.(*preserves.Record); ok {
		if sym, ok := rec.Key.(*preserves.Symbol); ok {
			return sym.String() == string(*r)
		}
	}
	return false
}
