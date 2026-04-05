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
	case *ast.BlockStatement:
		return evalStatements(v.Statements)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: v.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(v.Value)
	case *ast.IfExpression:
		return evalIfExpression(v)
	case *ast.PrefixExpression:
		right := Eval(v.Right)
		return evalPrefixExpression(v.Operator, right)
	case *ast.InfixExpression:
		left := Eval(v.Left)
		right := Eval(v.Right)
		return evalInfixExpression(v.Operator, left, right)
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

// evalIfExpression evaluates an if expression node
func evalIfExpression(node *ast.IfExpression) object.Object {
	condition := Eval(node.Condition)
	if isTruthy(condition) {
		return Eval(node.Consequence)
	} else if node.Alternative != nil {
		return Eval(node.Alternative)
	} else {
		return NULL
	}
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

// evalInfixExpression evaluates a generic infix expression
func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)

	// catch-alls for other types. directly comparing left and right will be correct for bools since we use singleton objects
	case op == "==":
		return nativeBoolToBooleanObject(left == right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}

// evalIntegerInfixExpression evaluates an infix expression between two integers
func evalIntegerInfixExpression(op string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NULL
	}
}

// nativeBoolToBooleanObject converts a go bool into a singleton boolean object
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// isTruthy returns if an object is truthy or not
func isTruthy(obj object.Object) bool {
	switch obj {
	case TRUE:
		return true
	case FALSE, NULL:
		return false
	default:
		return true
	}
}
