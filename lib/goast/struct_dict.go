package goast

import (
	"fmt"
	"go/ast"
	"go/token"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/isodude/preserves-go/lib/preserves"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type StructDict struct {
	Name            string
	Fields          []*ast.Field
	Kind            ObjectType
	StructKind      ObjectType
	mapFieldsToType map[string]ObjectType
	mapKeyToField   []preserves.Value
	identifier      []AST
}

func (d *StructDict) GetObjectType() ObjectType {
	return StructDictType
}
func (d *StructDict) Under(_ AST) {}
func (d *StructDict) GetName() string {
	return d.Name
}
func (d *StructDict) GetTitle() string {
	return cases.Title(language.English, cases.NoLower).String(d.Name)
}
func (d *StructDict) SetKind(o ObjectType) {
	d.Kind = o
}
func (d *StructDict) SetStructKind(o ObjectType) {
	d.StructKind = o
}

/*
func (s *Struct) FromPreserves(Value) *Struct {

}
*/
func (d *StructDict) AST(above AST) (decl []ast.Decl) {
	name := d.Name
	sname := &ast.StarExpr{X: ast.NewIdent(d.Name)}
	/*
		func (*Lit) FromPreserves(value Value) *Lit {
			if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
				if s, ok := rec.Key.(*Symbol); ok && s.String() == "lit" {
					if a := AnyFromPreserves(rec.Fields[0]); a != nil {
						return &Lit{Value: a}
					}
				}
			}
			return nil
		}
	*/
	var firstStmt ast.Stmt
	/*
		if dict, ok := rec.Fields[0].(*Dictionary); ok && len(*dict) > 0 && len(*dict) < 4 {
			var schema Schema
			for k, v := range *dict {
				if sym := NewSymbol("").FromPreserves(k); sym != nil {
					switch string(*sym) {
	*/
	firstStmt = &ast.IfStmt{
		Init: &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				ast.NewIdent("dict"),
				ast.NewIdent("ok"),
			},
			Rhs: []ast.Expr{
				&ast.TypeAssertExpr{
					X: &ast.IndexExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("rec"),
							Sel: ast.NewIdent("Fields"),
						},
						Index: ast.NewIdent("0"),
					},
					Type: &ast.StarExpr{X: ast.NewIdent("Dictionary")},
				},
			},
		},
		Cond: ast.NewIdent("ok"),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.TypeSpec{
								Name: ast.NewIdent("obj"),
								Type: ast.NewIdent(name),
							},
						},
					},
				},
				&ast.RangeStmt{
					Tok:   token.DEFINE,
					Key:   ast.NewIdent("dictKey"),
					Value: ast.NewIdent("dictValue"),
					X:     &ast.StarExpr{X: ast.NewIdent("dict")},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{},
					},
				},
			},
		},
	}
	var toElements []ast.Expr
	for i, field := range d.Fields {
		if len(field.Names) != 1 {
			continue
		}

		varName := ast.NewIdent(fmt.Sprintf("p%d", i))
		var callExprFun *ast.Ident
		fieldName := field.Names[0].String()
		fieldName = fmt.Sprintf("%s%s", strings.ToUpper(string(fieldName[0])), fieldName[1:])
		var fieldType string
		switch t := field.Type.(type) {
		case *ast.ArrayType:
			switch u := t.Elt.(type) {
			case *ast.Ident:
				fieldType = u.String()
			default:
				panic(fmt.Sprintf("beep: %s", reflect.TypeOf(field.Type)))
			}
		case *ast.Ident:
			fieldType = t.String()
		default:
			panic(fmt.Sprintf("beep: %s", reflect.TypeOf(field.Type)))
		}

		if fieldType == "String" {
			fieldType = "Pstring"
		}
		if callExprFun = MapToCallFunc(d.mapFieldsToType, fieldType); callExprFun == nil {
			continue
		}
		o, ok := d.mapFieldsToType[strings.ToLower(fieldType)]
		if !ok {
			continue
		}
		var dVarName ast.Expr
		dVarName = varName
		if o != UnionInterfaceObjectType {
			dVarName = &ast.StarExpr{X: varName}
		}
		/*
			case "definitions":
							if d := (&Definitions{}).FromPreserves(v); d != nil {
								schema.Definitions = *d
							} else {
								return nil
							}
		*/
		var value ast.Expr
		switch t := d.mapKeyToField[i].(type) {
		case *preserves.Symbol:
			value = &ast.CallExpr{
				Fun:  ast.NewIdent("NewSymbol"),
				Args: []ast.Expr{ast.NewIdent(strconv.Quote(t.String()))},
			}
		case *preserves.SignedInteger:
			a := big.Int(*t)
			value = &ast.CallExpr{
				Fun:  ast.NewIdent("NewSignedInteger"),
				Args: []ast.Expr{ast.NewIdent(strconv.Quote(a.String()))},
			}
		case *preserves.Boolean:
			value = &ast.CallExpr{
				Fun:  ast.NewIdent("NewBoolean"),
				Args: []ast.Expr{ast.NewIdent(fmt.Sprintf("%t", bool(*t)))},
			}
		}
		toElements = append(toElements, &ast.KeyValueExpr{Key: value, Value: &ast.CallExpr{
			Fun: ast.NewIdent(fmt.Sprintf("%sToPreserves", fieldType)),
			Args: []ast.Expr{&ast.SelectorExpr{
				X:   ast.NewIdent("d"),
				Sel: ast.NewIdent(fieldName),
			}},
		}})
		stmt := &ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun:  ast.NewIdent("dictKey.Equal"),
				Args: []ast.Expr{value},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.IfStmt{
						Init: &ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{varName},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: callExprFun,
									Args: []ast.Expr{
										ast.NewIdent("dictValue"),
									},
								},
							},
						},
						Cond: &ast.BinaryExpr{
							Op: token.NEQ,
							X:  varName,
							Y:  ast.NewIdent("nil"),
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.AssignStmt{
									Tok: token.ASSIGN,
									Lhs: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("obj"),
											Sel: ast.NewIdent(fieldName),
										},
									},
									Rhs: []ast.Expr{
										dVarName,
									},
								},
								&ast.ExprStmt{X: ast.NewIdent("continue")},
							},
						},
					},
				},
			},
		}
		firstStmt.(*ast.IfStmt).Body.List[1].(*ast.RangeStmt).Body.List = append(firstStmt.(*ast.IfStmt).Body.List[1].(*ast.RangeStmt).Body.List, stmt)

		if len(d.Fields) == i+1 {
			returnStmt := ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X:  ast.NewIdent("obj"),
					},
				},
			}
			firstStmt.(*ast.IfStmt).Body.List = append(firstStmt.(*ast.IfStmt).Body.List, &returnStmt)
		}
	}
	// Add return nil inside the for loop
	firstStmt.(*ast.IfStmt).Body.List[1].(*ast.RangeStmt).Body.List = append(firstStmt.(*ast.IfStmt).Body.List[1].(*ast.RangeStmt).Body.List, &ast.ReturnStmt{Results: []ast.Expr{ast.NewIdent("nil")}})
	var ifStmt ast.Stmt

	ifStmt = &ast.IfStmt{
		Init: &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				ast.NewIdent("sym"),
				ast.NewIdent("ok"),
			},
			Rhs: []ast.Expr{
				&ast.TypeAssertExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("rec"),
						Sel: ast.NewIdent("Key"),
					},
					Type: &ast.StarExpr{X: ast.NewIdent("Symbol")},
				},
			},
		},
		Cond: &ast.BinaryExpr{
			Op: token.LAND,
			X:  ast.NewIdent("ok"),
			Y: &ast.BinaryExpr{
				Op: token.EQL,
				X: &ast.CallExpr{
					Fun: ast.NewIdent("sym.String"),
				},
				Y: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote(d.Name),
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{firstStmt},
		},
	}
	var keyValues []ast.Expr
	if len(d.identifier) > 0 {
		if stmter, ok := d.identifier[0].(Stmt); ok {
			ifStmt = stmter.Stmt(&ast.SelectorExpr{
				X:   ast.NewIdent("rec"),
				Sel: ast.NewIdent("Key"),
			}, []ast.Stmt{firstStmt})
		}
		if stmter, ok := d.identifier[0].(ToStmt); ok {
			keyValues = stmter.ToStmt(ast.NewIdent("Key"))
		}
	}
	keyValues = append(keyValues, &ast.KeyValueExpr{Key: ast.NewIdent("Fields"), Value: &ast.ArrayType{Elt: &ast.CompositeLit{Type: ast.NewIdent("Value"), Elts: []ast.Expr{&ast.UnaryExpr{Op: token.AND, X: &ast.CompositeLit{Type: ast.NewIdent("Dictionary"), Elts: toElements}}}}}})
	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "FromPreserves")),
			Type: &ast.FuncType{
				Func: token.Pos(token.FUNC),
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("value")},
							Type:  ast.NewIdent("Value"),
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{},
							Type:  sname,
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.IfStmt{
						Init: &ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{
								ast.NewIdent("rec"),
								ast.NewIdent("ok"),
							},
							Rhs: []ast.Expr{
								&ast.TypeAssertExpr{
									X:    ast.NewIdent("value"),
									Type: &ast.StarExpr{X: ast.NewIdent("Record")},
								},
							},
						},
						Cond: &ast.BinaryExpr{
							Op: token.LAND,
							X:  ast.NewIdent("ok"),
							Y: &ast.BinaryExpr{
								Op: token.EQL,
								X: &ast.CallExpr{
									Fun: ast.NewIdent("len"),
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("rec"),
											Sel: ast.NewIdent("Fields"),
										},
									},
								},
								Y: ast.NewIdent("1"),
							},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								ifStmt,
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{ast.NewIdent("nil")},
					},
				},
			},
		},
	)

	/*
			func SchemaToPreserves(d Schema) Value {
			return &Record{Key: NewSymbol("schema"), Fields: []Value{&Dictionary{
				NewSymbol("definitions"):  DefinitionsToPreserves(d.Definitions),
				NewSymbol("embeddedType"): EmbeddedTypeNameToPreserves(d.EmbeddedType),
				NewSymbol("version"):      VersionToPreserves(d.Version),
			}}}
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
							Names: []*ast.Ident{ast.NewIdent("d")},
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
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.CompositeLit{
									Type: ast.NewIdent("Record"),
									Elts: keyValues,
								},
							},
						},
					}},
			},
		},
	)
	return
}
