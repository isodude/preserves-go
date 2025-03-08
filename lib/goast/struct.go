package goast

import (
	"fmt"
	"go/ast"

	"github.com/isodude/preserves-go/lib/preserves"
)

type Struct struct {
	Fields          []*ast.Field
	ASTFields       []*Field
	Identifier      []AST
	Value           preserves.Value
	Kind            ObjectType
	StructKind      ObjectType
	mapFieldsToType map[string]ObjectType
	MapKeyToField   []preserves.Value
	title
}

func NewStruct(name string) *Struct {
	s := &Struct{}
	s.SetName(name)
	return s
}
func (s *Struct) Stmt(key ast.Expr, stmts []ast.Stmt) ast.Stmt {
	if u, ok := s.GetStructKind().(Stmt); ok {
		return u.Stmt(key, stmts)
	}
	panic("should not be reached")
}
func (s *Struct) ToStmt(key ast.Expr) []ast.Expr {
	if u, ok := s.GetStructKind().(ToStmt); ok {
		return u.ToStmt(key)
	}
	panic("should not be reached")
}

func (s *Struct) SetValue(v preserves.Value) {
	s.StructKind = LitType
	s.Value = v
}
func (s *Struct) GetObjectType() ObjectType {
	if u, ok := s.GetStructKind().(AST); ok {
		return u.GetObjectType()
	}
	return StructObjectType
}
func (s *Struct) Under(_ AST) {}
func (s *Struct) AddField(a *Field) {
	s.ASTFields = append(s.ASTFields, a)
}
func (s *Struct) SetKind(o ObjectType) {
	s.Kind = o
}
func (s *Struct) SetStructKind(o ObjectType) {
	s.StructKind = o
}
func (s *Struct) GetStructKind() AST {
	switch s.StructKind {
	case LitType:
		return &Lit{
			title: s.title,
			Type:  s.Value,
		}
	case StructRecType:
		panic(fmt.Sprintf("%s: %d: %v", s.GetName(), s.StructKind, s))
	case StructDictType:
		return &StructDict{
			title:           s.title,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			mapKeyToField:   s.MapKeyToField,
			identifier:      s.Identifier,
		}
	case FirstArrayType:
		panic(fmt.Sprintf("%s: %d: %v", s.GetName(), s.StructKind, s))
	case LastArrayType:
		panic(fmt.Sprintf("%s: %d: %v", s.GetName(), s.StructKind, s))
	case AllSameTypeArrayType:
		panic(fmt.Sprintf("%s: %d: %v", s.GetName(), s.StructKind, s))
	case TupleType:
		return &Tuple{
			title:           s.title,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			identifier:      s.Identifier,
		}
	case StructTupleType:
		return &StructTuple{
			title:           s.title,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			identifier:      s.Identifier,
		}
	case StructTuplePrefixType:
		return &TuplePrefix{
			title:           s.title,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			identifier:      s.Identifier,
		}
	case StructSeqofType:
		return &Seqof{
			title:           s.title,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
		}
	default:
		panic(fmt.Sprintf("%s: %d", s.GetName(), s.StructKind))
	}
}

func (s *Struct) AST(above AST) (decl []ast.Decl) {
	name := s.GetTitle()
	var aboveTitle *title
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
		aboveTitle = above.Title()
	}
	if !(s.StructKind == StructSeqofType || (s.StructKind == LitType && above == nil) || s.StructKind == TupleType) {
		decl = append(decl, objectTypeSpec(aboveTitle, s.Title(), s.ASTFields))
		decl = append(decl, objectFuncDeclNew(aboveTitle, s.Title(), s.ASTFields))
		if above != nil {
			decl = append(decl, objectFuncDeclIs(aboveTitle, s.Title(), aboveTitle))
		}
	}

	decl = append(decl, s.GetStructKind().AST(above)...)

	return
}
