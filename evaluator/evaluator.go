package evaluator

import (
	"github.com/michaelzhan1/go-interpreter/ast"
	"github.com/michaelzhan1/go-interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
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
		if v.Value {
			return TRUE
		}
		return FALSE
	case *ast.PrefixExpression:
		right := Eval(v.Right)
		return evalPrefixExpression(v.Operator, right)
	}

	return nil
}

// evalStatements evaluates a slice of statements
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt) // TODO: Flesh this out
	}

	return result
}

// evalPrefixExpression evaluates a prefix expression
func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return NULL
	}
}

// evalBangOperatorExpression evaluates the ! prefix operator
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// evalMinusOperatorExpression evaluates the - prefix operator
func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
