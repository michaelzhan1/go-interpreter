package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/michaelzhan1/go-interpreter/ast"
)

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
)

// ObjectType is the underlying type of an Object
type ObjectType string

// Object is a the interface used for all values
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer is an Object representing an integer
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

var _ Object = &Integer{}

// Boolean is an Object representing a boolean
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

var _ Object = &Boolean{}

// String is an Object representing a string
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

var _ Object = &String{}

// Function is an Object representing a function
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	var out bytes.Buffer

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// BuiltInFunction is a built in function that operates on objects
type BuiltInFunction func(args ...Object) Object

// BuiltIn is an Object representing a built in function
type BuiltIn struct {
	Fn BuiltInFunction
}

func (bi *BuiltIn) Type() ObjectType { return BUILTIN_OBJ }
func (bi *BuiltIn) Inspect() string  { return "built-in function" }

// Null is an Object representing null
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

var _ Object = &Null{}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

var _ Object = &ReturnValue{}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

var _ Object = &Error{}
