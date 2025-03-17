package goast

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/isodude/preserves-go/lib/preserves"
)

type fieldType interface {
	Expr() ast.Expr
	Name() string
}

type Field struct {
	ASTs            []AST
	Type            string
	fieldType       fieldType
	Array           bool
	Dict            bool
	DictKey         AST
	DictValue       AST
	Value           preserves.Value
	mapFieldsToType map[string]ObjectType
	title
}

func NewField(name string) *Field {
	f := &Field{}
	f.SetName(name)
	return f
}

func (f *Field) ToUpper() string {
	if f.name == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(string(f.name[0])), f.name[1:])
}
func (f *Field) Expr() ast.Expr {
	if f.Array {
		return &ast.ArrayType{Elt: ast.NewIdent(f.GetFieldTypeName())}
	}
	if f.Dict {
		return ast.NewIdent(f.convertType(f.DictValue.GetTitle()))
	}
	return ast.NewIdent(f.GetFieldTypeName())
	// return f.fieldType.Expr()
}
func (f *Field) convertType(s string) string {
	if strings.ToLower(s) == "any" {
		return "Value"
	}
	if s == "String" {
		return "Pstring"
	}
	return s
}
func (f *Field) GetFieldTypeName() string {
	return f.convertType(f.Type)
	// return f.fieldType.Name()
}

func (f *Field) GetToPreservesFunction(m map[string]ObjectType) *ast.Ident {
	if f.GetFieldTypeName() == "String" {
		return ast.NewIdent("ShimStringToPreserves")
	}
	if f.GetFieldTypeName() == "SignedInteger" {
		return ast.NewIdent("ShimIntToPreserves")
	}
	return ToMapToCallFunc(m, f.GetFieldTypeName())
}

func (f *Field) GetFromPreservesFunction(m map[string]ObjectType) *ast.Ident {
	if f.GetFieldTypeName() == "String" {
		return ast.NewIdent("ShimStringFromPreserves")
	}
	if f.GetFieldTypeName() == "SignedInteger" {
		return ast.NewIdent("ShimIntFromPreserves")
	}
	return MapToCallFunc(m, f.GetFieldTypeName())
}

func (f *Field) AddMaybeRef(m map[string]ObjectType, ident *ast.Ident) ast.Expr {
	o, ok := m[strings.ToLower(f.GetFieldTypeName())]
	if !ok {
		return nil
	}
	if !(o == UnionInterfaceObjectType || o == InterfaceObjectType) {
		return &ast.StarExpr{X: ident}
	}
	return ident
}

func (f *Field) SetArray() {
	f.Array = true
}
func (f *Field) SetDict(key AST, value AST) {
	f.Dict = true
	f.DictKey = key
	f.DictValue = value
}
func (f *Field) SetValue(v preserves.Value) {
	f.Value = v
}
func (f *Field) ConvertToMap(a AST) *Map {
	var key, value string
	if u, ok := f.DictKey.(*Field); ok {
		key = u.Type
	}
	if u, ok := f.DictValue.(*Field); ok {
		value = u.Type
	}
	return &Map{title: *a.Title(), Key: key, Value: value}
}
func (f *Field) GetObjectType() ObjectType {
	// TODO
	return UnionInterfaceObjectType
}

func (f *Field) Under(a AST) {
	f.ASTs = append(f.ASTs, a)
}
func (*Field) SetKind(o ObjectType) {}
func (f *Field) SetType(s string) {
	f.Type = s
}
func (*Field) SetStructKind(o ObjectType)    {}
func (f *Field) AST(_ AST) (decl []ast.Decl) { return }
func (f *Field) ASTField() *ast.Field {
	if f.Array {
		return &ast.Field{
			Names: []*ast.Ident{f.Ident(nil, false)},
			Type:  &ast.ArrayType{Elt: ast.NewIdent(f.Type)},
		}
	}
	if f.Dict {
		return &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(f.DictKey.GetName())},
			Type:  ast.NewIdent(f.DictValue.GetName()),
		}
	}
	return &ast.Field{
		Names: []*ast.Ident{f.Ident(nil, false)},
		Type:  ast.NewIdent(f.Type),
	}
}
