package evaluator

import (
	"curryLang/ast"
	"curryLang/object"
	"curryLang/token"
)

var (
	NULL = &object.Null{}
)

type Variable struct {
	Name  string
	Value object.Object
}

type ExecutionEngine struct {
	Variables       []Variable
	CurrentStackPos uint32

	Functions map[string]*object.Function
}

func NewEngine() *ExecutionEngine {
	engine := ExecutionEngine{}
	engine.Functions = make(map[string]*object.Function)
	return &engine
}

func (engine *ExecutionEngine) PushStack() {
	engine.CurrentStackPos = uint32(len(engine.Variables))
}

func (engine *ExecutionEngine) PopStack() {
	engine.Variables = engine.Variables[engine.CurrentStackPos:]
}

func (engine *ExecutionEngine) Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case *ast.Identifier:
		return engine.EvalIdentifier(node)

	case *ast.IfElseExpression:
		return engine.EvalIfElseExpression(node)
	case *ast.PrefixExpression:
		return engine.EvalPrefixExpression(node)
	case *ast.InfixExpression:
		return engine.EvalInfixExpression(node)

	case *ast.LetStatement:
		return engine.EvalLetStatement(node)

	case *ast.ReturnStatement:
		return engine.EvalReturnStatement(node)

	case *ast.FunctionExpression:
		return engine.EvalFunctionExpression(node)

	case *ast.FunctionCallExpression:
		return engine.EvalFunctionCallExpression(node)

	case *ast.ExpressionStatement:
		return engine.Eval(node.Expression)

	case *ast.Program:
		return engine.EvalStatements(node.Statements)
	}

	return NULL
}

func (engine *ExecutionEngine) EvalLetStatement(statement *ast.LetStatement) object.Object {
	val := engine.Eval(statement.Value)
	variable := Variable{
		Name:  statement.Name.Value,
		Value: val,
	}

	engine.Variables = append(engine.Variables, variable)

	return NULL
}

func (engine *ExecutionEngine) EvalReturnStatement(statement *ast.ReturnStatement) object.Object {
	return engine.Eval(statement.ReturnValue)
}

func (engine *ExecutionEngine) EvalFunctionExpression(statement *ast.FunctionExpression) object.Object {
	function := &object.Function{
		Name:       statement.Name,
		Parameters: statement.Parameters,
		Code:       statement.Body,
	}

	if statement.Name != "" {
		engine.Functions[statement.Name] = function
	}

	return function
}

func (engine *ExecutionEngine) EvalFunctionCallExpression(statement *ast.FunctionCallExpression) object.Object {
	functionExpr := engine.Eval(statement.FunctionExpr)

	if functionExpr == NULL {
		funcIdentifier, ok := statement.FunctionExpr.(*ast.Identifier)
		if !ok {
			return NULL
		}

		functionExpr = engine.Functions[funcIdentifier.Value]
	}

	function, ok := functionExpr.(*object.Function)
	if !ok {
		return NULL
	}

	engine.PushStack()
	// add parameters as variables to current stack
	for i, parameter := range function.Parameters {
		paramValue := engine.Eval(statement.Parameters[i])
		engine.Variables = append(engine.Variables, Variable{
			Name:  parameter.Name,
			Value: paramValue,
		})
	}

	result := engine.EvalStatements(function.Code)

	engine.PopStack()

	return result
}

func (engine *ExecutionEngine) EvalIdentifier(identifier *ast.Identifier) object.Object {

	for _, variable := range engine.Variables {
		if variable.Name == identifier.Value {
			return variable.Value
		}
	}

	return NULL
}

func (engine *ExecutionEngine) EvalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = engine.Eval(stmt)
	}

	return result
}

func (engine *ExecutionEngine) EvalIfElseExpression(ifElse *ast.IfElseExpression) object.Object {
	conditionResult := engine.Eval(ifElse.Condition)
	if conditionResult.Type() != object.BOOLEAN_OBJ {
		return NULL
	}

	condition := conditionResult.(*object.Boolean)

	var statements []ast.Statement

	if condition.Value {
		statements = ifElse.Consequence
	} else {
		statements = ifElse.Alternative
	}

	engine.PushStack()
	result := engine.EvalStatements(statements)
	engine.PopStack()

	return result
}

func (engine *ExecutionEngine) EvalPrefixExpression(infix *ast.PrefixExpression) object.Object {
	value := engine.Eval(infix.Right)

	if value.Type() == object.BOOLEAN_OBJ {
		return engine.EvalBooleanPrefixOperations(value.(*object.Boolean), infix.Operator)
	}
	if value.Type() == object.INTEGER_OBJ {
		return engine.EvalIntegerPrefixOperations(value.(*object.Integer), infix.Operator)
	}

	return NULL
}

func (engine *ExecutionEngine) EvalBooleanPrefixOperations(val *object.Boolean, operator string) object.Object {

	if operator == token.BANG {
		return &object.Boolean{Value: !val.Value}
	}

	return NULL
}

func (engine *ExecutionEngine) EvalIntegerPrefixOperations(val *object.Integer, operator string) object.Object {

	if operator == token.MINUS {
		return &object.Integer{Value: -1 * val.Value}
	}

	return NULL
}

func (engine *ExecutionEngine) EvalInfixExpression(infix *ast.InfixExpression) object.Object {
	left := engine.Eval(infix.Left)
	right := engine.Eval(infix.Right)

	if left.Type() != right.Type() {
		return NULL
	}

	if left.Type() == object.INTEGER_OBJ {
		return engine.EvalIntegerInfixOperations(left.(*object.Integer), right.(*object.Integer), infix.Operator)
	}

	return NULL
}

func (engine *ExecutionEngine) EvalIntegerInfixOperations(left *object.Integer, right *object.Integer, operator string) object.Object {

	switch operator {
	// logical operations
	case token.LT:
		return &object.Boolean{Value: left.Value < right.Value}
	case token.GT:
		return &object.Boolean{Value: left.Value > right.Value}
	case token.EQ:
		return &object.Boolean{Value: left.Value == right.Value}
	case token.NOT_EQ:
		return &object.Boolean{Value: left.Value != right.Value}

	// arithmetic operations
	case token.PLUS:
		return &object.Integer{Value: left.Value + right.Value}
	case token.MINUS:
		return &object.Integer{Value: left.Value - right.Value}
	case token.ASTERISK:
		return &object.Integer{Value: left.Value * right.Value}
	case token.SLASH:
		return &object.Integer{Value: left.Value / right.Value}
	}

	return NULL
}
