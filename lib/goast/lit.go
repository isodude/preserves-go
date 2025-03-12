package goast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/isodude/preserves-go/lib/preserves"
	"github.com/isodude/preserves-go/lib/preserves/text"
)

type Lit struct {
	Type preserves.Value
	title
}

func NewLit(v preserves.Value) *Lit {
	return &Lit{Type: v}
}
func (l *Lit) GetObjectType() ObjectType {
	return StructObjectType
}
func (l *Lit) Under(_ AST)              {}
func (*Lit) SetKind(o ObjectType)       {}
func (*Lit) SetStructKind(o ObjectType) {}
func (l *Lit) Stmt(key ast.Expr, stmts []ast.Stmt) ast.Stmt {
	fieldType := reflect.TypeOf(l.Type).String()
	fieldType = strings.TrimPrefix(fieldType, "*preserves.")
	var fieldValue ast.Expr
	switch t := l.Type.(type) {
	case *preserves.Boolean:
		fieldValue = &ast.CallExpr{
			Fun:  ast.NewIdent("NewBoolean"),
			Args: []ast.Expr{ast.NewIdent(fmt.Sprintf("%t", reflect.ValueOf(t).Elem().Bool()))},
		}
	case *preserves.SignedInteger:
		a := big.Int(*t)
		fieldValue = &ast.CallExpr{
			Fun:  ast.NewIdent("NewSignedInteger"),
			Args: []ast.Expr{ast.NewIdent(strconv.Quote(a.String()))},
		}
	case *preserves.Pstring:
		fieldValue = &ast.CallExpr{
			Fun:  ast.NewIdent("NewPstring"),
			Args: []ast.Expr{ast.NewIdent(strconv.Quote(string(*t)))},
		}
	case *preserves.Symbol:
		fieldValue = &ast.CallExpr{
			Fun:  ast.NewIdent("NewSymbol"),
			Args: []ast.Expr{ast.NewIdent(strconv.Quote(t.String()))},
		}
	}
	return &ast.IfStmt{
		Init: &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				ast.NewIdent("v"),
			},
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent(fmt.Sprintf("%s%s", fieldType, "FromPreserves")),
					Args: []ast.Expr{
						key,
					},
				},
			},
		},
		Cond: &ast.BinaryExpr{
			Op: token.NEQ,
			X:  ast.NewIdent("v"),
			Y:  ast.NewIdent("nil"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.IfStmt{
					Cond: &ast.CallExpr{
						Fun:  ast.NewIdent("v.Equal"),
						Args: []ast.Expr{fieldValue},
					},
					Body: &ast.BlockStmt{List: stmts},
				},
			},
		},
	}
}
func (l *Lit) ToStmt(key ast.Expr) []ast.Expr {
	fieldType := reflect.TypeOf(l.Type).String()
	fieldType = strings.TrimPrefix(fieldType, "*preserves.")
	var fieldValue ast.Expr
	switch t := l.Type.(type) {
	case *preserves.Boolean:
		fieldValue = &ast.CallExpr{
			Fun:  ast.NewIdent("NewBoolean"),
			Args: []ast.Expr{ast.NewIdent(fmt.Sprintf("%t", reflect.ValueOf(t).Elem().Bool()))},
		}
	case *preserves.SignedInteger:
		a := big.Int(*t)
		fieldValue = &ast.CallExpr{
			Fun:  ast.NewIdent("NewSignedInteger"),
			Args: []ast.Expr{ast.NewIdent(strconv.Quote(a.String()))},
		}
	case *preserves.Pstring:
		fieldValue = &ast.CallExpr{
			Fun:  ast.NewIdent("NewPstring"),
			Args: []ast.Expr{ast.NewIdent(strconv.Quote(string(*t)))},
		}
	case *preserves.Symbol:
		fieldValue = &ast.CallExpr{
			Fun:  ast.NewIdent("NewSymbol"),
			Args: []ast.Expr{ast.NewIdent(strconv.Quote(t.String()))},
		}
	}
	return []ast.Expr{
		&ast.KeyValueExpr{
			Key:   key,
			Value: fieldValue,
		},
	}
}
func (l *Lit) GetType() string {
	var b bytes.Buffer
	_, err := l.Type.WriteTo(&b)
	if err != nil {
		return ""
	}
	return b.String()
}
func (l *Lit) AST(above AST) (decl []ast.Decl) {
	name := l.GetTitle()
	var aboveTitle *title
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
		aboveTitle = above.Title()
	}
	var b bytes.Buffer
	if l.Type == nil {
		panic(fmt.Sprintf("type in lit %s was nil", name))
	}
	_, err := text.FromPreserves(l.Type).WriteTo(&b)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	if above == nil {
		decl = append(decl, objectTypeSpec(nil, l.Title(), []*Field{}))
		decl = append(decl, objectFuncDeclNew(nil, l.Title(), []*Field{}))
	}

	decl = append(decl, objectFuncDeclFromPreserves(aboveTitle, l.Title(), []*Field{}, []ast.Stmt{
		l.Stmt(ast.NewIdent("value"), []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{
			&ast.UnaryExpr{Op: token.AND, X: &ast.CompositeLit{Type: ast.NewIdent(name)}}}}}),
		&ast.ReturnStmt{Results: []ast.Expr{ast.NewIdent("nil")}},
	}))

	decl = append(decl, objectFuncDeclToPreserves(aboveTitle, l.Title(), []*Field{}, []ast.Stmt{
		&ast.ReturnStmt{
			Results: []ast.Expr{
				l.ToStmt(nil)[0].(*ast.KeyValueExpr).Value,
			},
		},
	}))

	return
}
