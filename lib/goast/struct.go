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
	Identifier      []AST
	Kind            ObjectType
	StructKind      ObjectType
	mapFieldsToType map[string]ObjectType
	MapKeyToField   []preserves.Value
}

func (s *Struct) GetObjectType() ObjectType {
	return StructObjectType
}
func (s *Struct) Under(_ AST) {}
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
			Name: ast.NewIdent(fmt.Sprintf("New%s", name)),
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
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

	switch s.StructKind {
	case StructRecType:
		decl = append(decl, (&Rec{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			identifier:      s.Identifier,
		}).AST(above)...)
	case StructDictType:
		decl = append(decl, (&StructDict{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			mapKeyToField:   s.MapKeyToField,
			identifier:      s.Identifier,
		}).AST(above)...)
	case FirstArrayType:
	case LastArrayType:
	case AllSameTypeArrayType:
	case StructTupleType:
		decl = append(decl, (&StructTuple{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			identifier:      s.Identifier,
		}).AST(above)...)
	case StructTuplePrefixType:
		decl = append(decl, (&TuplePrefix{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
			identifier:      s.Identifier,
		}).AST(above)...)
	case StructSeqofType:
		decl = append(decl, (&Seqof{
			Name:            s.Name,
			Fields:          s.Fields,
			mapFieldsToType: s.mapFieldsToType,
		}).AST(above)...)
	default:
		panic("beep")
	}

	return
}
