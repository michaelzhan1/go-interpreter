package evaluator

import (
	"github.com/michaelzhan1/go-interpreter/ast"
	"github.com/michaelzhan1/go-interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval evaluates an ast.Node and returns the evaluated object
func Eval(node ast.Node) object.Object {
	switch v := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(v.Statements)
	case *ast.ExpressionStatement:
		return Eval(v.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: v.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(v.Value)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt) // TODO: Flesh this out
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
