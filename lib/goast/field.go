package goast

import (
	"go/ast"

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

func (f *Field) Expr() ast.Expr {
	return f.fieldType.Expr()
}

func (f *Field) GetFieldTypeName() string {
	return f.fieldType.Name()
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
	/*allLit := true
	oneLit := false
	for _, a := range u.ASTs {
		if _, ok := a.(*Lit); !ok {
			allLit = false
		} else {
			oneLit = true
		}
	}
	if allLit {
		return UnionConstObjectType
	}
	if oneLit {
		return UnionVariantObjectType
	}*/
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
