package evaluator

import (
	"fmt"

	"github.com/michaelzhan1/go-interpreter/ast"
	"github.com/michaelzhan1/go-interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

// Eval evaluates an ast.Node and returns the evaluated object
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch v := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(v, env)
	case *ast.LetStatement:
		val := Eval(v.Value, env)
		if isError(val) {
			return val
		}
		env.Set(v.Name.TokenLiteral(), val)
	case *ast.ReturnStatement:
		val := Eval(v.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.ExpressionStatement:
		return Eval(v.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(v, env)

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(v, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: v.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(v.Value)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: v.Parameters,
			Body:       v.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(v.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(v.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.IfExpression:
		return evalIfExpression(v, env)
	case *ast.PrefixExpression:
		right := Eval(v.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(v.Operator, right)
	case *ast.InfixExpression:
		left := Eval(v.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(v.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(v.Operator, left, right)
	}

	return nil
}

// evalProgram evaluates a program
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch v := result.(type) {
		case *object.ReturnValue:
			return v.Value
		case *object.Error:
			return v
		}
	}

	return result
}

// evalBlockStatement evaluates a block statement
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// evalExpressions evaluates a slice of ast.Expression
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

// evalIdentifier evaluates an identifier
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.TokenLiteral())
	if !ok {
		return newError("identifier not found: %s", node.TokenLiteral())
	}
	return val
}

// evalIfExpression evaluates an if expression node
func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

// applyFunction applies a function with args
func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

// extendFunctionEnv takes a function's env and creates an inner env with the outer's scope
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		env.Set(param.TokenLiteral(), args[i]) // set each param to its evaluated arg value
	}
	return env
}

// unwrapReturnValue bubbles up the return value from the function
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value // stops a return value from returning all function call stack early
	}
	return obj
}

// evalPrefixExpression evaluates a prefix expression
func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", op, right.Type())
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
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// evalInfixExpression evaluates a generic infix expression
func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())

	// catch-alls for other types. directly comparing left and right will be correct for bools since we use singleton objects
	case op == "==":
		return nativeBoolToBooleanObject(left == right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
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

// newError returns a new error
func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// isError returns if an object is an error
func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}
