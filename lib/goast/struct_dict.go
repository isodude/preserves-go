package goast

import (
	"fmt"
	"go/ast"
	"go/token"
	"math/big"
	"reflect"
	"strconv"

	"github.com/isodude/preserves-go/lib/preserves"
)

type StructDict struct {
	Fields          []*Field
	Kind            ObjectType
	StructKind      ObjectType
	mapFieldsToType map[string]ObjectType
	mapKeyToField   []preserves.Value
	identifier      Stmt
	title
}

func (d *StructDict) GetObjectType() ObjectType {
	return StructDictType
}
func (d *StructDict) Under(_ AST) {}
func (d *StructDict) SetKind(o ObjectType) {
	d.Kind = o
}
func (d *StructDict) SetStructKind(o ObjectType) {
	d.StructKind = o
}

func (d *StructDict) AST(above AST) (decl []ast.Decl) {
	name := d.GetName()
	sname := &ast.StarExpr{X: d.Ident(nil, false)}
	var firstStmt ast.Stmt
	/*
		if dict, ok := rec.Fields[0].(*Dictionary); ok {
			var obj {{name}}
			for dictKey, dictValue := range *dict {
			}
		}
	*/
	bodyStmt := &ast.BlockStmt{List: []ast.Stmt{}}
	firstStmt = &ast.IfStmt{
		Init: &ast.AssignStmt{Tok: token.DEFINE,
			Lhs: []ast.Expr{ast.NewIdent("dict"), ast.NewIdent("ok")},
			Rhs: []ast.Expr{&ast.TypeAssertExpr{X: &ast.IndexExpr{
				X:     &ast.SelectorExpr{X: ast.NewIdent("rec"), Sel: ast.NewIdent("Fields")},
				Index: ast.NewIdent("0"),
			}, Type: &ast.StarExpr{X: ast.NewIdent("Dictionary")},
			}}},
		Cond: ast.NewIdent("ok"),
		Body: &ast.BlockStmt{List: []ast.Stmt{
			&ast.DeclStmt{Decl: &ast.GenDecl{
				Tok:   token.VAR,
				Specs: []ast.Spec{&ast.TypeSpec{Name: ast.NewIdent("obj"), Type: ast.NewIdent(name)}},
			}},
			&ast.RangeStmt{
				Tok:   token.DEFINE,
				Key:   ast.NewIdent("dictKey"),
				Value: ast.NewIdent("dictValue"),
				X:     &ast.StarExpr{X: ast.NewIdent("dict")},
				Body:  bodyStmt,
			},
		}},
	}
	var toElements []ast.Expr
	for i, field := range d.Fields {
		varName := ast.NewIdent(fmt.Sprintf("p%d", i))

		fieldName := field.ToUpper()
		callExprFun := field.GetFromPreservesFunction(d.mapFieldsToType)
		if callExprFun == nil {
			continue
		}
		dVarName := field.AddMaybeRef(d.mapFieldsToType, varName)
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
		default:
			// TODO: support more preserves
			panic(fmt.Sprintf("%s not implemented", reflect.TypeOf(t)))
		}
		/*
			[NewSymbol("symbol")] = ModulesPathToPreserves(d.ModulesPath)
		*/
		fun := field.GetToPreservesFunction(d.mapFieldsToType)
		if fun == nil {
			continue
		}
		toElements = append(toElements, &ast.KeyValueExpr{Key: value, Value: &ast.CallExpr{
			Fun: fun,
			Args: []ast.Expr{&ast.SelectorExpr{
				X:   ast.NewIdent("d"),
				Sel: ast.NewIdent(fieldName),
			}},
		}})

		/*
			if dictKey.Equal({{value}}) {
				if {{varName}} := {{callExprFun}}(dictValue); {{varName}} != nil {
					obj[{{fieldName}}] = {{dVarName}}
					continue
				}
			}
		*/
		stmt := &ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun:  ast.NewIdent("dictKey.Equal"),
				Args: []ast.Expr{value},
			},
			Body: &ast.BlockStmt{List: []ast.Stmt{&ast.IfStmt{
				Init: &ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{varName},
					Rhs: []ast.Expr{&ast.CallExpr{Fun: callExprFun, Args: []ast.Expr{ast.NewIdent("dictValue")}}},
				},
				Cond: &ast.BinaryExpr{Op: token.NEQ, X: varName, Y: ast.NewIdent("nil")},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.AssignStmt{Tok: token.ASSIGN,
						Lhs: []ast.Expr{&ast.SelectorExpr{X: ast.NewIdent("obj"), Sel: ast.NewIdent(fieldName)}},
						Rhs: []ast.Expr{dVarName}},
					&ast.ExprStmt{X: ast.NewIdent("continue")},
				}},
			}}},
		}
		bodyStmt.List = append(bodyStmt.List, stmt)

		if len(d.Fields) == i+1 {
			returnStmt := ast.ReturnStmt{Results: []ast.Expr{&ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("obj")}}}
			firstStmt.(*ast.IfStmt).Body.List = append(firstStmt.(*ast.IfStmt).Body.List, &returnStmt)
		}
	}
	// return nil
	bodyStmt.List = append(bodyStmt.List, &ast.ReturnStmt{Results: []ast.Expr{ast.NewIdent("nil")}})
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
					Value: strconv.Quote(d.GetName()),
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{firstStmt},
		},
	}

	// &Record{Key: NewSymbol("schema"),
	var keyValues []ast.Expr
	if d.identifier != nil {
		ifStmt = d.identifier.Stmt(&ast.SelectorExpr{
			X:   ast.NewIdent("rec"),
			Sel: ast.NewIdent("Key"),
		}, []ast.Stmt{firstStmt})
		keyValues = d.identifier.ToStmt(ast.NewIdent("Key"))
	}
	// Fields: []Value{&Dictionary{NewSymbol("definitions"): DefinitionsToPreserves(d.Definitions), NewSymbol("embeddedType"): EmbeddedTypeNameToPreserves(d.EmbeddedType), NewSymbol("version"):   VersionToPreserves(d.Version)}}}
	keyValues = append(keyValues,
		&ast.KeyValueExpr{
			Key: ast.NewIdent("Fields"),
			Value: &ast.ArrayType{Elt: &ast.CompositeLit{
				Type: ast.NewIdent("Value"),
				Elts: []ast.Expr{&ast.UnaryExpr{
					Op: token.AND,
					X:  &ast.CompositeLit{Type: ast.NewIdent("Dictionary"), Elts: toElements},
				}},
			}},
		})

	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "FromPreserves")),
			Type: &ast.FuncType{
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
