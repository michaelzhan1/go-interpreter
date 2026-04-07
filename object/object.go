package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/michaelzhan1/go-interpreter/ast"
)

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
)

// ObjectType is the underlying type of an Object
type ObjectType string

// Object is the interface used for all values
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Hashable is the interface for all hashable objects
type Hashable interface {
	Object
	HashKey() HashKey
}

// HashKey is a struct that lets us hash objects
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Integer is an Object representing an integer
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey { return HashKey{Type: i.Type(), Value: uint64(i.Value)} }

var _ Object = &Integer{}
var _ Hashable = &Integer{}

// Boolean is an Object representing a boolean
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

var _ Object = &Boolean{}
var _ Hashable = &Boolean{}

// String is an Object representing a string
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

var _ Object = &String{}
var _ Hashable = &String{}

// Array is an Object representing an array
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	es := []string{}
	for _, e := range a.Elements {
		es = append(es, e.Inspect())
	}

	var out bytes.Buffer

	out.WriteString(("["))
	out.WriteString(strings.Join(es, ", "))
	out.WriteString("]")

	return out.String()
}

var _ Object = &Array{}

// HashPair defines a hashed pair. It is needed because Hash stores a map keyed by the key's hash => HashPair
type HashPair struct {
	Key   Object
	Value Object
}

// Hash is an Object representing a hash map
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	ss := make([]string, len(h.Pairs))
	for _, p := range h.Pairs {
		ss = append(ss, p.Key.Inspect()+":"+p.Value.Inspect())
	}

	var out bytes.Buffer

	out.WriteString("{")
	out.WriteString(strings.Join(ss, ", "))
	out.WriteString("}")

	return out.String()
}

var _ Object = &Hash{}

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

var _ Object = &Function{}

// BuiltInFunction is a built in function that operates on objects
type BuiltInFunction func(args ...Object) Object

// BuiltIn is an Object representing a built in function
type BuiltIn struct {
	Fn BuiltInFunction
}

func (bi *BuiltIn) Type() ObjectType { return BUILTIN_OBJ }
func (bi *BuiltIn) Inspect() string  { return "built-in function" }

var _ Object = &BuiltIn{}

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
