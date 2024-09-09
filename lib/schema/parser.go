package schema

import (
	. "github.com/isodude/preserves-go/lib/preserves"
)

var _document = NewBundle(NewModules().Add(*NewModulePath(),
	*NewSchema(
		NewDefinitions().
			Add(*NewSymbol("AtomKind"), NewDefinitionOr(
				*NewNamedAlternative("Boolean", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("Boolean"))))),
				*NewNamedAlternative("Double", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("Double"))))),
				[]NamedAlternative{
					*NewNamedAlternative("SignedInteger", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("SignedInteger"))))),
					*NewNamedAlternative("String", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("String"))))),
					*NewNamedAlternative("ByteString", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("ByteString"))))),
					*NewNamedAlternative("Symbol", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("Symbol"))))),
				},
			)).
			Add(*NewSymbol("Binding"), NewDefinitionPattern(NewPatternCompoundPattern(NewCompoundPatternRec(
				NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("named")))), NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
					[]NamedPattern{
						NewNamedPatternNamed(*NewBinding(*NewSymbol("name"), NewSimplePatternAtom(&AtomKindSymbol{}))),
						NewNamedPatternNamed(*NewBinding(*NewSymbol("pattern"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("SimplePattern"))))),
					},
				))))))).
			Add(*NewSymbol("Bundle"), NewDefinitionPattern(NewPatternCompoundPattern(NewCompoundPatternRec(NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("bundle")))), NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
				[]NamedPattern{NewNamedPatternNamed(*NewBinding(*NewSymbol("modules"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Modules")))))},
			))))))).
			Add(*NewSymbol("CompoundPattern"), NewDefinitionOr(
				*NewNamedAlternative("rec", NewPatternCompoundPattern(NewCompoundPatternRec(
					NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("rec")))), NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
						[]NamedPattern{
							NewNamedPatternNamed(*NewBinding(*NewSymbol("label"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedPattern"))))),
							NewNamedPatternNamed(*NewBinding(*NewSymbol("fields"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedPattern"))))),
						},
					))),
				))),
				*NewNamedAlternative("tuple", NewPatternCompoundPattern(NewCompoundPatternRec(
					NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("tuple")))), NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
						[]NamedPattern{NewNamedPatternNamed(*NewBinding(*NewSymbol("patterns"), NewSimplePatternSeqof(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedPattern"))))))},
					))),
				))),
				[]NamedAlternative{
					*NewNamedAlternative("tuplePrefix", NewPatternCompoundPattern(NewCompoundPatternRec(
						NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("tuplePrefix")))), NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
							[]NamedPattern{
								NewNamedPatternNamed(*NewBinding(*NewSymbol("fixed"), NewSimplePatternSeqof(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedPattern")))))),
								NewNamedPatternNamed(*NewBinding(*NewSymbol("variable"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedPattern"))))),
							},
						))),
					))),

					*NewNamedAlternative("dict", NewPatternCompoundPattern(NewCompoundPatternRec(
						NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("dict")))), NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
							[]NamedPattern{NewNamedPatternNamed(*NewBinding(*NewSymbol("entries"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("DictionaryEntries")))))},
						))),
					)))},
			)).
			Add(*NewSymbol("Definition"), NewDefinitionOr(
				*NewNamedAlternative("or", NewPatternCompoundPattern(NewCompoundPatternRec(
					NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("or")))), NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
						[]NamedPattern{
							NewNamedPatternAnonymous(
								NewPatternCompoundPattern(
									NewCompoundPatternTuplePrefix(
										[]NamedPattern{
											NewNamedPatternNamed(*NewBinding(*NewSymbol("pattern0"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedAlternative"))))),
											NewNamedPatternNamed(*NewBinding(*NewSymbol("pattern1"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedAlternative"))))),
										},
										NewNamedSimplePatternNamed(*NewBinding(*NewSymbol("patternN"), NewSimplePatternSeqof(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedAlternative")))))),
									),
								),
							),
						},
					))),
				))),
				*NewNamedAlternative("and", NewPatternCompoundPattern(NewCompoundPatternRec(
					NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("and")))), NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
						[]NamedPattern{
							NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuplePrefix(
								[]NamedPattern{
									NewNamedPatternNamed(*NewBinding(*NewSymbol("pattern0"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedPattern"))))),
									NewNamedPatternNamed(*NewBinding(*NewSymbol("pattern1"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedPattern"))))),
								},
								NewNamedSimplePatternNamed(*NewBinding(*NewSymbol("patternN"), NewSimplePatternSeqof(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedPattern")))))),
							)))}))),
				))),
				[]NamedAlternative{
					*NewNamedAlternative("Pattern", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Pattern")))))),
				},
			)).
			Add(*NewSymbol("Definitions"), NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternDictof(NewSimplePatternAtom(&AtomKindSymbol{}), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Definition"))))))).
			Add(*NewSymbol("DictionaryEntries"), NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternDictof(&SimplePatternAny{}, NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("NamedSimplePattern"))))))).
			Add(*NewSymbol("EmbeddedTypeName"), NewDefinitionOr(
				*NewNamedAlternative("false", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternLit(NewBoolean(false))))),
				*NewNamedAlternative("Ref", NewPatternSimplePattern(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Ref"))))),
				[]NamedAlternative{},
			)).
			Add(*NewSymbol("ModulePath"), NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternSeqof(NewSimplePatternAtom(&AtomKindSymbol{}))))).
			Add(*NewSymbol("Modules"), NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternDictof(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("ModulePath"))), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Schema"))))))).
			Add(*NewSymbol("NamedAlternative"), NewDefinitionPattern(NewPatternCompoundPattern(NewCompoundPatternTuple(
				[]NamedPattern{
					NewNamedPatternNamed(*NewBinding(*NewSymbol("variantLabel"), NewSimplePatternAtom(&AtomKindString{}))),
					NewNamedPatternNamed(*NewBinding(*NewSymbol("pattern"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Pattern"))))),
				},
			)))).
			Add(*NewSymbol("NamedPattern"), NewDefinitionOr(
				*NewNamedAlternative("named", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Binding")))))),
				*NewNamedAlternative("anonymous", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Pattern")))))),
				[]NamedAlternative{},
			)).
			Add(*NewSymbol("NamedSimplePattern"), NewDefinitionOr(
				*NewNamedAlternative("named", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Binding")))))),
				*NewNamedAlternative("anonymous", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("SimplePattern")))))),
				[]NamedAlternative{},
			)).
			Add(*NewSymbol("Pattern"), NewDefinitionOr(
				*NewNamedAlternative("SimplePattern", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("SimplePattern")))))),
				*NewNamedAlternative("CompoundPattern", NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("CompoundPattern")))))),
				[]NamedAlternative{},
			)).
			Add(*NewSymbol("Schema"), NewDefinitionPattern(NewPatternCompoundPattern(NewCompoundPatternRec(
				NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("schema")))),
				NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
					[]NamedPattern{
						NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternDict(NewDictionaryEntries().
							Add(NewSymbol("definitions"), NewNamedSimplePatternNamed(*NewBinding(*NewSymbol("definitions"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Definitions")))))).
							Add(NewSymbol("embeddedType"), NewNamedSimplePatternNamed(*NewBinding(*NewSymbol("embeddedType"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("EmbeddedTypeName")))))).
							Add(NewSymbol("version"), NewNamedSimplePatternNamed(*NewBinding(*NewSymbol("version"), NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Version")))))),
						)))},
				))),
			)))).
			Add(*NewSymbol("SimplePattern"), NewDefinitionOr(
				*NewNamedAlternative(
					"any",
					NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("any")))),
				),
				*NewNamedAlternative(
					"atom",
					NewPatternCompoundPattern(NewCompoundPatternRec(
						NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("atom")))),
						NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
							[]NamedPattern{
								NewNamedPatternNamed(*NewBinding(
									*NewSymbol("atomKind"),
									NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("AtomKind"))),
								)),
							},
						))),
					)),
				),
				[]NamedAlternative{
					*NewNamedAlternative(
						"embedded",
						NewPatternCompoundPattern(NewCompoundPatternRec(
							NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("embedded")))),
							NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
								[]NamedPattern{
									NewNamedPatternNamed(*NewBinding(
										*NewSymbol("interface"),
										NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("SimplePattern"))),
									)),
								},
							))),
						)),
					),
					*NewNamedAlternative(
						"lit",
						NewPatternCompoundPattern(NewCompoundPatternRec(
							NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("lit")))),
							NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
								[]NamedPattern{
									NewNamedPatternNamed(*NewBinding(
										*NewSymbol("value"),
										NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("any"))),
									)),
								},
							))),
						)),
					),
					*NewNamedAlternative(
						"seqof",
						NewPatternCompoundPattern(NewCompoundPatternRec(
							NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("seqof")))),
							NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
								[]NamedPattern{
									NewNamedPatternNamed(*NewBinding(
										*NewSymbol("pattern"),
										NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("SimplePattern"))),
									)),
								},
							))),
						)),
					),
					*NewNamedAlternative(
						"setof",
						NewPatternCompoundPattern(NewCompoundPatternRec(
							NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("setof")))),
							NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
								[]NamedPattern{
									NewNamedPatternNamed(*NewBinding(
										*NewSymbol("pattern"),
										NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("SimplePattern"))),
									)),
								},
							))),
						)),
					),
					*NewNamedAlternative(
						"dictof",
						NewPatternCompoundPattern(NewCompoundPatternRec(
							NewNamedPatternAnonymous(NewPatternSimplePattern(NewSimplePatternLit(NewSymbol("dictof")))),
							NewNamedPatternAnonymous(NewPatternCompoundPattern(NewCompoundPatternTuple(
								[]NamedPattern{
									NewNamedPatternNamed(*NewBinding(
										*NewSymbol("key"),
										NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("SimplePattern"))),
									)),
									NewNamedPatternNamed(*NewBinding(
										*NewSymbol("value"),
										NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("SimplePattern"))),
									)),
								},
							))),
						)),
					),
					*NewNamedAlternative(
						"Ref",
						NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternRef(*NewRef(*NewModulePath(), *NewSymbol("Ref"))))),
					),
				},
			)).
			Add(*NewSymbol("Version"), NewDefinitionPattern(NewPatternSimplePattern(NewSimplePatternLit(NewSignedInteger("1"))))),
		&EmbeddedTypeNameFalse{},
		Version{},
	)))
