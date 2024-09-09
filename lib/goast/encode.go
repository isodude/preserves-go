package goast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func strToCamelCase(s string) string {
	if len(s) > 0 {
		return fmt.Sprintf("%s%s", strings.ToUpper(string(s[0])), s[1:])
	}
	return ""
}

type AST interface {
	AST(AST) []ast.Decl
	GetName() string
	GetTitle() string
	GetObjectType() ObjectType
	Under(AST)
	SetKind(ObjectType)
	SetStructKind(ObjectType)
}

type Stmt interface {
	Stmt(ast.Expr, []ast.Stmt) ast.Stmt
}
type ToStmt interface {
	ToStmt(ast.Expr) []ast.Expr
}
type Encoder interface {
	EncodeToGoAST(AST, string) []AST
}
type Fields interface {
	EncodeToGoASTFields(AST) []*ast.Field
}
type String interface {
	ASTString() string
}

func ToMapToCallFunc(m map[string]ObjectType, key string) *ast.Ident {
	return ast.NewIdent(fmt.Sprintf("%sToPreserves", key))
}
func MapToCallFunc(m map[string]ObjectType, key string) *ast.Ident {
	return ast.NewIdent(fmt.Sprintf("%sFromPreserves", key))
	/*
		o, ok := m[strings.ToLower(key)]
		if !ok {
			return nil
		}
		var callExprFun *ast.Ident
		if o == UnionInterfaceObjectType {
			callExprFun = ast.NewIdent(fmt.Sprintf("%sFromPreserves", key))
		} else if o == UnionConstObjectType {
			callExprFun = ast.NewIdent(fmt.Sprintf("(&%s{}).FromPreserves", key))
		} else if o == PassthroughObjectType {
			callExprFun = ast.NewIdent(fmt.Sprintf("(&%s{}).FromPreserves", key))
		} else if o == UnionVariantObjectType {
			callExprFun = ast.NewIdent(fmt.Sprintf("(&%s{}).FromPreserves", key))

		} else if o == StructObjectType {
			callExprFun = ast.NewIdent(fmt.Sprintf("(&%s{}).FromPreserves", key))
		} else if o == MapObjectType {
			callExprFun = ast.NewIdent(fmt.Sprintf("New%s().FromPreserves", key))
		} else if o == SimpleStringType {
			callExprFun =
		} else {
			callExprFun = ast.NewIdent(fmt.Sprintf("missing: %v (%d)", key, o))
		}
		return callExprFun */
}
func EncodeToGoAST(name string, e Encoder) string {
	asts := e.EncodeToGoAST(nil, name)
	asts = append(asts, []AST{&Boolean{}, &SignedInteger{}, &Pstring{}, &Symbol{}, &Value{}}...)
	m := make(map[string]AST)
	for _, t := range asts {
		m[strings.ToLower(t.GetName())] = t
	}

	var f func(AST, string, ObjectType, bool)
	f = func(a AST, name string, o ObjectType, again bool) {
		switch b := a.(type) {
		case *Union:
			for _, c := range b.ASTs {
				f(c, name, o, again)
			}
			if b.mapFieldsToType == nil {
				b.mapFieldsToType = make(map[string]ObjectType)
			}

			b.mapFieldsToType[strings.ToLower(name)] = o

		case *Passthrough:
			if strings.ToLower(b.Object) == name {
				b.ObjectType = o
				//if c, ok := obj.(*Union); ok {
				//	b.ASTs = append(b.ASTs, c.ASTs...)
				//}
			}
			if b.mapFieldsToType == nil {
				b.mapFieldsToType = make(map[string]ObjectType)
			}
			b.mapFieldsToType[strings.ToLower(name)] = o
		case *Struct:
			if b.mapFieldsToType == nil {
				b.mapFieldsToType = make(map[string]ObjectType)
			}
			b.mapFieldsToType[strings.ToLower(name)] = o
		case *Map:
			if b.mapFieldsToType == nil {
				b.mapFieldsToType = make(map[string]ObjectType)
			}
			b.mapFieldsToType[strings.ToLower(name)] = o

		case *Tuple:
			if b.mapFieldsToType == nil {
				b.mapFieldsToType = make(map[string]ObjectType)
			}
			b.mapFieldsToType[strings.ToLower(name)] = o
		case *Seqof:
			if b.mapFieldsToType == nil {
				b.mapFieldsToType = make(map[string]ObjectType)
			}
			b.mapFieldsToType[strings.ToLower(name)] = o
		}
	}
	for k, v := range m {
		for _, vv := range m {
			f(vv, k, v.GetObjectType(), false)
		}
		switch uv := v.(type) {
		case *Union:
			for _, uuv := range uv.ASTs {
				for _, vv := range m {
					name = fmt.Sprintf("%s%s", uv.GetName(), uuv.GetName())
					f(vv, name, uuv.GetObjectType(), true)
				}
			}
		}
	}
	ts := []ast.Decl{
		&ast.GenDecl{
			Tok:    token.IMPORT,
			Lparen: token.Pos(token.LPAREN),
			Rparen: token.Pos(token.RPAREN),
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote("github.com/isodude/preserves-go/lib/extras")},
				},
				&ast.ImportSpec{
					Name: ast.NewIdent("."),
					Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote("github.com/isodude/preserves-go/lib/preserves")},
				},
			},
		},
	}

	for _, v := range asts {
		ts = append(ts, v.AST(nil)...)
	}

	astFile := &ast.File{
		Name:  ast.NewIdent("beep"),
		Decls: ts,
	}
	fset := token.NewFileSet()
	var bytes bytes.Buffer
	err := printer.Fprint(&bytes, fset, astFile)
	if err != nil {
		log.Errorf("%v", err)
	}
	return bytes.String()
}
