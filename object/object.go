package object

import (
	"curryLang/ast"
	"fmt"
)

type ObjectType string

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	FUNCITON_OBJ = "FUNCTION"
	NULL_OBJ     = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (i *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (i *Boolean) Inspect() string  { return fmt.Sprintf("%t", i.Value) }

type Function struct {
	Name       string
	Parameters []ast.Parameter
	Code       []ast.Statement
}

func (function *Function) Type() ObjectType { return FUNCITON_OBJ }
func (function *Function) Inspect() string  { return fmt.Sprintf("fn %s", function.Name) }

type Null struct{}

func (i *Null) Type() ObjectType { return NULL_OBJ }
func (i *Null) Inspect() string  { return "null" }
