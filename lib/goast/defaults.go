package goast

import (
	"go/ast"
)

type Boolean struct{ title }

func NewBoolean() *Boolean                       { return &Boolean{title: title{name: "boolean"}} }
func (*Boolean) GetObjectType() ObjectType       { return SimpleBoolType }
func (*Boolean) Under(_ AST)                     {}
func (*Boolean) SetKind(o ObjectType)            {}
func (*Boolean) SetStructKind(o ObjectType)      {}
func (*Boolean) AST(above AST) (decl []ast.Decl) { return }

type Symbol struct{ title }

func NewSymbol() *Symbol                        { return &Symbol{title: title{name: "symbol"}} }
func (*Symbol) GetObjectType() ObjectType       { return SimpleStringType }
func (*Symbol) Under(_ AST)                     {}
func (*Symbol) GetName() string                 { return "symbol" }
func (*Symbol) GetTitle() string                { return "Symbol" }
func (*Symbol) SetKind(o ObjectType)            {}
func (*Symbol) SetStructKind(o ObjectType)      {}
func (*Symbol) AST(above AST) (decl []ast.Decl) { return }

type SignedInteger struct{ title }

func NewSignedInteger() *SignedInteger                 { return &SignedInteger{title: title{name: "SignedInteger"}} }
func (*SignedInteger) GetObjectType() ObjectType       { return SimpleSignedIntegerType }
func (*SignedInteger) Under(_ AST)                     {}
func (*SignedInteger) GetName() string                 { return "SignedInteger" }
func (*SignedInteger) GetTitle() string                { return "SignedInteger" }
func (*SignedInteger) SetKind(o ObjectType)            {}
func (*SignedInteger) SetStructKind(o ObjectType)      {}
func (*SignedInteger) AST(above AST) (decl []ast.Decl) { return }

type Pstring struct{ title }

func NewPstring() *Pstring                       { return &Pstring{title: title{name: "Pstring"}} }
func (*Pstring) GetObjectType() ObjectType       { return SimplePstringType }
func (*Pstring) Under(_ AST)                     {}
func (*Pstring) GetName() string                 { return "Pstring" }
func (*Pstring) GetTitle() string                { return "Pstring" }
func (*Pstring) SetKind(o ObjectType)            {}
func (*Pstring) SetStructKind(o ObjectType)      {}
func (*Pstring) AST(above AST) (decl []ast.Decl) { return }

type Value struct{ title }

func NewValue() *Value                         { return &Value{title: title{name: "Value"}} }
func (*Value) GetObjectType() ObjectType       { return ValueType }
func (*Value) Under(_ AST)                     {}
func (*Value) SetKind(o ObjectType)            {}
func (*Value) SetStructKind(o ObjectType)      {}
func (*Value) AST(above AST) (decl []ast.Decl) { return }
