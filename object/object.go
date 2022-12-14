package object

import (
	"curryLang/ast"
	"fmt"
)

type ObjectType string

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	STRING_OBJ   = "STRING"
	FUNCITON_OBJ = "FUNCTION"
	LIST_OBJ     = "LIST"
	PACKAGE_OBJ  = "PACKAGE"
	ERROR_OBJ    = "ERROR"
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

type String struct {
	Value string
}

func (str *String) Type() ObjectType { return STRING_OBJ }
func (str *String) Inspect() string  { return str.Value }

type Function struct {
	Name       string
	Parameters []ast.Parameter
	Code       []ast.Statement
}

func (function *Function) Type() ObjectType { return FUNCITON_OBJ }
func (function *Function) Inspect() string  { return fmt.Sprintf("fn %s", function.Name) }

type List struct {
	ValueType ObjectType
	Value     []Object
}

func (list *List) Type() ObjectType { return LIST_OBJ }
func (list *List) Inspect() string  { return "list<" + string(list.ValueType) + ">" }

type Package struct {
	ValueType ObjectType
	Name      string
	Globals   map[string]Object
	Functions map[string]*Function
}

func (pkg *Package) Type() ObjectType { return PACKAGE_OBJ }
func (pkg *Package) Inspect() string  { return "Package " + pkg.Name }

type Null struct{}

func (i *Null) Type() ObjectType { return NULL_OBJ }
func (i *Null) Inspect() string  { return "null" }

type Error struct {
	Message string
}

func (err *Error) Type() ObjectType { return ERROR_OBJ }
func (err *Error) Inspect() string  { return fmt.Sprintf("error#%s", err.Message) }
