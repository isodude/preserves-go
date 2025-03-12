package goast

import (
	"go/ast"
)

type Label struct {
	stmt Stmt
	// fullfill AST
	title
}

func NewLabel() *Label {
	return &Label{}
}

// fullfill AST
func (l *Label) AST(AST) (decl []ast.Decl)     { return }
func (l *Label) GetObjectType() (o ObjectType) { return }
func (l *Label) Under(AST)                     {}
func (l *Label) SetKind(ObjectType)            {}
func (l *Label) SetStructKind(ObjectType)      {}

// fullfill Stmt
func (l *Label) Stmt(key ast.Expr, stmts []ast.Stmt) ast.Stmt {
	if l.stmt != nil {
		return l.stmt.Stmt(key, stmts)
	}
	return nil
}

func (l *Label) ToStmt(key ast.Expr) []ast.Expr {
	if l.stmt != nil {
		return l.stmt.ToStmt(key)
	}
	return nil
}

func (l *Label) SetStmt(s Stmt) {
	l.stmt = s
}
