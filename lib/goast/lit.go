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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Lit struct {
	Name string
	Type preserves.Value
}

func (l *Lit) GetObjectType() ObjectType {
	return StructObjectType
}
func (l *Lit) GetName() string {
	return l.Name
}
func (l *Lit) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(l.Name)
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
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
	}
	var b bytes.Buffer
	_, err := text.FromPreserves(l.Type).WriteTo(&b)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	//field := b.String()
	//var fieldType ast.Expr
	//fieldType = ast.NewIdent(fmt.Sprintf("\"%s\"", field))
	decl = append(decl, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.StructType{
					Fields: &ast.FieldList{},
				},
			},
		},
	})
	sname := &ast.StarExpr{X: ast.NewIdent(name)}
	if above != nil {
		decl = append(decl,
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  sname,
						}},
				},
				Name: ast.NewIdent(fmt.Sprintf("Is%s", above.GetName())),
				Type: &ast.FuncType{
					Func:   token.Pos(token.FUNC),
					Params: &ast.FieldList{},
				},
				Body: &ast.BlockStmt{},
			})
	}

	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "FromPreserves")),
			Type: &ast.FuncType{

				Func: token.Pos(token.FUNC),
				Params: &ast.FieldList{

					List: []*ast.Field{{
						Names: []*ast.Ident{ast.NewIdent("value")},
						Type:  ast.NewIdent("Value"),
					}},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  &ast.StarExpr{X: ast.NewIdent(name)},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					l.Stmt(ast.NewIdent("value"), []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{
						&ast.UnaryExpr{Op: token.AND, X: &ast.CompositeLit{Type: ast.NewIdent(name)}}}}}),
					&ast.ReturnStmt{Results: []ast.Expr{ast.NewIdent("nil")}},
				},
			},
		},
	)

	/*

	   func AtomKindStringToPreserves(l AtomKindString) Value {
	   	return SymbolToPreserves(NewSymbol("String"))
	   }
	*/
	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "ToPreserves")),
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("l")},
							Type:  ast.NewIdent(name),
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  ast.NewIdent("Value"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							l.ToStmt(nil)[0].(*ast.KeyValueExpr).Value,
						},
					},
				},
			},
		})
	return
}
