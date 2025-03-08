package goast

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Passthrough struct {
	Object          string
	ObjectType      ObjectType
	ASTs            []AST
	mapFieldsToType map[string]ObjectType
	title
}

func NewPassthrough(name string) *Passthrough {
	return &Passthrough{title: title{name: name}}
}

type ObjectType int

const (
	InvalidObjectType ObjectType = iota

	DefinitionObjectType
	InterfaceObjectType
	PassthroughObjectType
	UnionInterfaceObjectType
	UnionConstObjectType
	UnionVariantObjectType
	StructObjectType
	MapObjectType
	SeqofSymbolObjectType
	StructDictType
	StructRecType
	FirstArrayType
	LastArrayType
	LitType
	AllSameTypeArrayType
	StructSeqofType
	StructTupleType
	StructTuplePrefixType
	SimpleStringType
	SimpleBoolType
	SimpleSignedIntegerType
	SimplePstringType
	ValueType
	TupleType
	ArrayField
	SetField
	DictField
)

func (p *Passthrough) GetObjectType() ObjectType {
	return p.ObjectType
}
func (p *Passthrough) Under(a AST) {
	p.ASTs = append(p.ASTs, a)
}
func (*Passthrough) SetKind(o ObjectType)       {}
func (*Passthrough) SetStructKind(o ObjectType) {}
func (p *Passthrough) AST(above AST) (decl []ast.Decl) {
	name := p.GetName()
	decl = append(decl, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{{Type: ast.NewIdent(p.Object)}},
					},
				},
			},
		},
	})

	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("New%s", name)),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{{
						Names: []*ast.Ident{ast.NewIdent("obj")},
						Type:  ast.NewIdent(p.Object),
					}},
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
								Elts: []ast.Expr{&ast.KeyValueExpr{
									Key:   ast.NewIdent(p.Object),
									Value: ast.NewIdent("obj"),
								}},
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
					Params: &ast.FieldList{},
				},
				Body: &ast.BlockStmt{},
			})
	}
	/*
		func (p *{{name}}) FromPreserves(value Value) *{{name}} {
			var o {{p.Object}}
			for _, v := range _pattern {
				switch u := v.(type) {
				case *{{p.Structs[0].Name}}:
					o = u.FromPerservesSchemaAST(value)
				}
			}
			if o != nil {
				if p, ok := o.({{p.Object}}); ok {
					return &{{name}}}}{p}
				}
			}
		}
	*/

	objName := fmt.Sprintf("%s%s", p.GetName(), cases.Title(language.English, cases.NoLower).String(p.Object))
	if _, ok := p.mapFieldsToType[strings.ToLower(p.Object)]; ok {
		objName = cases.Title(language.English, cases.NoLower).String(p.Object)
	}

	callExprFun := ast.NewIdent(fmt.Sprintf("%sFromPreserves", objName))
	varName := ast.NewIdent("o")
	var dVarName ast.Expr

	switch p.ObjectType {
	case StructObjectType:
		fallthrough
	case StructTupleType:
		fallthrough
	case SimpleStringType:
		fallthrough
	case SimpleBoolType:
		fallthrough
	case SimpleSignedIntegerType:
		fallthrough
	case SimplePstringType:
		dVarName = &ast.StarExpr{X: varName}
	case InvalidObjectType:
		panic(fmt.Sprintf("could not find %s, %v", p.Object, reflect.TypeOf(p.ObjectType)))
	default:
		dVarName = varName
	}

	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "FromPreserves")),
			Type: &ast.FuncType{

				Params: &ast.FieldList{

					List: []*ast.Field{{
						Names: []*ast.Ident{ast.NewIdent("value")},
						Type:  ast.NewIdent("Value"),
					}},
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
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.IfStmt{
						Init: &ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{
								varName,
							},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: callExprFun,
									Args: []ast.Expr{
										ast.NewIdent("value"),
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
								&ast.ReturnStmt{
									Results: []ast.Expr{
										&ast.UnaryExpr{
											Op: token.AND,
											X: &ast.CompositeLit{
												Type: ast.NewIdent(name),
												Elts: []ast.Expr{
													&ast.KeyValueExpr{
														Key:   ast.NewIdent(objName),
														Value: dVarName,
													},
												},
											},
										},
									},
								},
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("nil"),
						},
					},
				},
			},
		},
	)
	decl = append(decl,
		&ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%s%s", name, "ToPreserves")),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("s")},
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
				List: []ast.Stmt{&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(
								fmt.Sprintf("%sToPreserves", objName)),
							Args: []ast.Expr{
								&ast.SelectorExpr{X: ast.NewIdent(`s`), Sel: ast.NewIdent(objName)},
							},
						},
					},
				},
				},
			},
		})
	for _, a := range p.ASTs {
		decl = append(decl, a.AST(p)...)
	}
	return
}
