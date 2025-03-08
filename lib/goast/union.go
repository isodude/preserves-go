package goast

import (
	"fmt"
	"go/ast"
)

type Union struct {
	ASTs            []AST
	mapFieldsToType map[string]ObjectType
	title
}

func NewUnion(name string) *Union {
	return &Union{title: title{name: name}}
}
func (u *Union) GetObjectType() ObjectType {
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
func (u *Union) Under(a AST) {
	u.ASTs = append(u.ASTs, a)
}
func (*Union) SetKind(o ObjectType)       {}
func (*Union) SetStructKind(o ObjectType) {}
func (u *Union) AST(above AST) (decl []ast.Decl) {
	name := u.GetName()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetName(), name)
	}
	/*
		allLit := true
		oneLit := false
		for _, a := range u.ASTs {
			if _, ok := a.(*Lit); !ok {
				allLit = false
			} else {
				oneLit = true
			}
		}
		if allLit {
			unionConst := NewUnionConst(name)
			unionConst.ASTs = u.ASTs
			unionConst.mapFieldsToType = u.mapFieldsToType
			decl = append(decl, unionConst.AST(above)...)
			return
		}
		if oneLit {
			unionVariant := NewUnionVariant(name)
			unionVariant.ASTs = u.ASTs
			unionVariant.mapFieldsToType = u.mapFieldsToType
			decl = append(decl, unionVariant.AST(above)...)
			return
		}
	*/
	unionInterface := NewUnionInterface(name)
	unionInterface.ASTs = u.ASTs
	unionInterface.mapFieldsToType = u.mapFieldsToType
	decl = append(decl, unionInterface.AST(above)...)
	return

}
