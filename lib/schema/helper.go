package schema

import (
	"go/ast"
	"reflect"
)

func fieldsOfStruct[T any](t T) (f []*ast.Field) {
	typ := reflect.TypeOf(t)
	val := reflect.ValueOf(t)
	for i := range val.NumField() {
		fld := typ.Field(i)
		switch fld.Type.Kind() {
		case reflect.Struct:
			f = append(f, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fld.Name)},
				Type:  ast.NewIdent(fld.Type.Name()),
			})
		case reflect.Slice:
			f = append(f, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fld.Name)},
				Type:  &ast.ArrayType{Elt: ast.NewIdent(fld.Type.Elem().Name())},
			})
		}
	}
	return f
}
