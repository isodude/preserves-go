package goast

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/isodude/preserves-go/lib/preserves"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Struct struct {
	Name            string
	Fields          []*ast.Field
	ASTFields       []*Field
	Identifier      []AST
	Value           preserves.Value
	Kind            ObjectType
	StructKind      ObjectType
	mapFieldsToType map[string]ObjectType
	MapKeyToField   []preserves.Value
}

func NewStruct(name string) *Struct {
	return &Struct{Name: name}
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
	return StructObjectType
}
func (s *Struct) Under(_ AST) {}
func (s *Struct) AddField(a *Field) {
	s.ASTFields = append(s.ASTFields, a)
}
func (s *Struct) GetName() string {
	return s.Name
}
func (s *Struct) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(s.Name)
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
			Name: s.Name,
			Type: s.Value,
		}
	case StructRecType:
		panic(fmt.Sprintf("%s: %d: %v", s.Name, s.StructKind, s))
		return &Rec{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			identifier:      s.Identifier,
		}
	case StructDictType:
		return &StructDict{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			mapKeyToField:   s.MapKeyToField,
			identifier:      s.Identifier,
		}
	case FirstArrayType:
		panic(fmt.Sprintf("%s: %d: %v", s.Name, s.StructKind, s))
	case LastArrayType:
		panic(fmt.Sprintf("%s: %d: %v", s.Name, s.StructKind, s))
	case AllSameTypeArrayType:
		panic(fmt.Sprintf("%s: %d: %v", s.Name, s.StructKind, s))
	case TupleType:
		return &Tuple{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			identifier:      s.Identifier,
		}
	case StructTupleType:
		return &StructTuple{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			identifier:      s.Identifier,
		}
	case StructTuplePrefixType:
		return &TuplePrefix{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			identifier:      s.Identifier,
		}
	case StructSeqofType:
		return &Seqof{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
		}
	default:
		panic(fmt.Sprintf("%s: %d", s.Name, s.StructKind))
	}
}

/*
func (s *Struct) FromPreserves(Value) *Struct {

}
*/
func (s *Struct) AST(above AST) (decl []ast.Decl) {
	name := s.GetTitle()
	if above != nil {
		name = fmt.Sprintf("%s%s", above.GetTitle(), name)
	}
	var fields []*ast.Field
	for _, field := range s.Fields {
		if len(field.Names) < 1 {
			continue
		}

		var nameType ast.Expr
		nameType = field.Type
		if ident, ok := nameType.(*ast.Ident); ok {
			if strings.ToLower(ident.String()) == "any" {
				nameType = ast.NewIdent("Value")
			}

			if ident.String() == "String" {
				nameType = ast.NewIdent("Pstring")
			}
		}

		fname := cases.Title(language.English, cases.NoLower).String(field.Names[0].String())
		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(fname)},
			Type:  nameType,
		})
	}
	decl = append(decl, &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// Generated via struct\n",
				},
			},
		},
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: fields,
					},
				},
			},
		},
	})
	var smallFields []*ast.Field
	for _, field := range fields {
		if len(field.Names) < 1 {
			continue
		}
		if len(field.Names[0].Name) < 2 {
			continue
		}
		fname := fmt.Sprintf("%s%s", strings.ToLower(field.Names[0].Name[0:1]), field.Names[0].Name[1:])
		switch fname {
		case "interface":
			fallthrough
		case "any":
			fname = fmt.Sprintf("_%s", fname)
		}
		smallFields = append(smallFields, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(fname)},
			Type:  field.Type,
		})
	}
	var smallValues []ast.Expr
	for _, field := range fields {
		if len(field.Names) < 1 {
			continue
		}
		if len(field.Names[0].Name) < 2 {
			continue
		}
		fname := fmt.Sprintf("%s%s", strings.ToLower(field.Names[0].Name[0:1]), field.Names[0].Name[1:])
		switch fname {
		case "interface":
			fallthrough
		case "any":
			fname = fmt.Sprintf("_%s", fname)
		}
		smallValues = append(smallValues, &ast.KeyValueExpr{
			Key:   field.Names[0],
			Value: ast.NewIdent(fname),
		})
	}

	decl = append(decl,
		&ast.FuncDecl{
			Doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{
						Text: "// Generated via struct\n",
					},
				},
			},
			Name: ast.NewIdent(fmt.Sprintf("New%s", name)),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: smallFields,
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
			Body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent(name),
								Elts: smallValues,
							},
						},
					},
				},
			},
			},
		},
	)

	sname := &ast.StarExpr{X: ast.NewIdent(name)}
	if above != nil {
		decl = append(decl,
			&ast.FuncDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// Generated via struct\n",
						},
					},
				},
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  sname,
						}},
				},
				Name: ast.NewIdent(fmt.Sprintf("Is%s", above.GetName())),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
				},
				Body: &ast.BlockStmt{},
			})
	}

	decl = append(decl, s.GetStructKind().AST(above)...)

	return
}
