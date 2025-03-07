package goast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"reflect"
	"slices"
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
	return Encode(name, asts)
}
func EncodeMapping(name string, asts []AST) []ast.Decl {
	m := make(map[string]AST)
	for _, t := range asts {
		m[strings.ToLower(t.GetName())] = t
	}

	var g func(AST)
	g = func(a AST) {
		switch b := a.(type) {
		case *Struct:
			for _, f := range b.ASTFields {
				if f.Dict {
					for i, c := range asts {
						if c == a {
							n := f.ConvertToMap(a)
							asts[i] = n
							m[strings.ToLower(n.GetName())] = n
						}
					}

					continue
				}
				b.Fields = append(b.Fields, f.ASTField())
			}
		case *Definition:
			for _, f := range b.ASTFields {
				if f.Dict {
					for i, c := range asts {
						if c == a {
							n := f.ConvertToMap(a)
							asts[i] = n
							m[strings.ToLower(n.GetName())] = n
						}
					}
					continue
				}
				b.Fields = append(b.Fields, f.ASTField())
			}
			for _, ast := range b.ASTs {
				g(ast)
			}
		case *Union:
			for _, ast := range b.ASTs {
				g(ast)
			}
		case *Passthrough:
			for _, ast := range b.ASTs {
				g(ast)
			}
		case *Boolean:
		case *SignedInteger:
		case *Pstring:
		case *Symbol:
		case *Value:
		case *Map:
		case *Field:
		default:
			panic(fmt.Sprintf("did not process %v", reflect.TypeOf(a)))
		}
	}
	for _, ast := range asts {
		g(ast)
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
			}
			if b.mapFieldsToType == nil {
				b.mapFieldsToType = make(map[string]ObjectType)
			}
			b.mapFieldsToType[strings.ToLower(name)] = o
		case *Definition:
			for _, c := range b.ASTs {
				f(c, name, o, again)
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
		case *Boolean:
		case *SignedInteger:
		case *Pstring:
		case *Symbol:
		case *Value:
		case *Field:
		default:
			panic(fmt.Sprintf("did not process %v", reflect.TypeOf(a)))
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

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		v := m[k]
		ts = append(ts, v.AST(nil)...)
	}

	return ts
}

func Encode(name string, asts []AST) string {
	asts = append(asts, []AST{&Boolean{}, &SignedInteger{}, &Pstring{}, &Symbol{}, &Value{}}...)

	astFile := &ast.File{
		Name:  ast.NewIdent(name),
		Decls: EncodeMapping(name, asts),
	}
	fset := token.NewFileSet()
	var bytes bytes.Buffer
	err := printer.Fprint(&bytes, fset, astFile)
	if err != nil {
		log.Errorf("%v", err)
	}
	return bytes.String()
}
