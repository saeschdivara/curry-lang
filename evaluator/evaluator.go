package evaluator

import (
	"curryLang/ast"
	"curryLang/object"
	"curryLang/token"
	"fmt"
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
	CurrentStackPos []uint32

	Functions map[string]*object.Function

	// engine state flags
	IsReturnTriggered bool
	HasError          bool
}

func NewEngine() *ExecutionEngine {
	engine := ExecutionEngine{IsReturnTriggered: false}
	engine.Variables = make([]Variable, 0)
	engine.CurrentStackPos = make([]uint32, 0)
	engine.Functions = make(map[string]*object.Function)
	return &engine
}

func (engine *ExecutionEngine) PushStack() {
	variablesSize := uint32(len(engine.Variables))
	stackPosLength := len(engine.CurrentStackPos)

	if stackPosLength == 0 {
		engine.CurrentStackPos = append(engine.CurrentStackPos, variablesSize)
	} else {
		currentPos := engine.CurrentStackPos[stackPosLength-1]

		if currentPos != variablesSize {
			engine.CurrentStackPos = append(engine.CurrentStackPos, variablesSize)
		}
	}
}

func (engine *ExecutionEngine) PopStack() {
	pos := len(engine.CurrentStackPos) - 1

	if pos == -1 {
		return
	}

	currentStackPos := engine.CurrentStackPos[pos]
	engine.Variables = engine.Variables[:currentStackPos]
	engine.CurrentStackPos = engine.CurrentStackPos[:pos]
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

	case *ast.AssignmentStatement:
		return engine.EvalAssignmentStatement(node)

	case *ast.WhileStatement:
		return engine.EvalWhileStatement(node)

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

func (engine *ExecutionEngine) EvalAssignmentStatement(statement *ast.AssignmentStatement) object.Object {

	identifierName := statement.Name.Value

	for i, variable := range engine.Variables {
		if variable.Name == identifierName {
			engine.Variables[i].Value = engine.Eval(statement.Value)
			return NULL
		}
	}

	return engine.createError(fmt.Sprintf("Tried to assign value to not existing variable %s", identifierName))
}

func (engine *ExecutionEngine) EvalWhileStatement(statement *ast.WhileStatement) object.Object {
	conditionResult := engine.Eval(statement.Condition)
	condition, ok := conditionResult.(*object.Boolean)

	if !ok {
		return engine.createError("Condition resulted with no boolean result")
	}

	for condition.Value {
		engine.PushStack()
		result := engine.EvalStatements(statement.Body)
		engine.PopStack()

		if result.Type() == object.ERROR_OBJ {
			return result
		}

		conditionResult = engine.Eval(statement.Condition)
		condition, ok = conditionResult.(*object.Boolean)

		if !ok {
			return engine.createError("Condition resulted with no boolean result")
		}
	}

	return NULL
}

func (engine *ExecutionEngine) EvalReturnStatement(statement *ast.ReturnStatement) object.Object {
	result := engine.Eval(statement.ReturnValue)
	engine.IsReturnTriggered = true
	return result
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
	engine.IsReturnTriggered = false

	engine.PopStack()

	return result
}

func (engine *ExecutionEngine) EvalIdentifier(identifier *ast.Identifier) object.Object {

	identifierName := identifier.Value

	for _, variable := range engine.Variables {
		if variable.Name == identifierName {
			return variable.Value
		}
	}

	if val, ok := engine.Functions[identifierName]; ok {
		return val
	}

	return engine.createError(fmt.Sprintf("Undeclared variable %s used", identifier.Value))
}

func (engine *ExecutionEngine) EvalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = engine.Eval(stmt)

		if engine.IsReturnTriggered || engine.HasError {
			break
		}
	}

	return result
}

func (engine *ExecutionEngine) EvalIfElseExpression(ifElse *ast.IfElseExpression) object.Object {
	conditionResult := engine.Eval(ifElse.Condition)

	conditionResType := conditionResult.Type()
	if conditionResType != object.BOOLEAN_OBJ {
		if conditionResType == object.ERROR_OBJ {
			return conditionResult
		} else {
			return engine.createError(fmt.Sprintf("Non boolean type (%s) was returned for condition", conditionResType))
		}
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

func (engine *ExecutionEngine) EvalPrefixExpression(prefix *ast.PrefixExpression) object.Object {
	value := engine.Eval(prefix.Right)
	valueType := value.Type()

	if valueType == object.BOOLEAN_OBJ {
		return engine.EvalBooleanPrefixOperations(value.(*object.Boolean), prefix.Operator)
	}
	if valueType == object.INTEGER_OBJ {
		return engine.EvalIntegerPrefixOperations(value.(*object.Integer), prefix.Operator)
	}

	return engine.createError(
		fmt.Sprintf("Not supported prefix operator (%s) was used for type %s", prefix.Operator, valueType),
	)
}

func (engine *ExecutionEngine) EvalBooleanPrefixOperations(val *object.Boolean, operator string) object.Object {

	if operator == token.BANG {
		return &object.Boolean{Value: !val.Value}
	}

	return engine.createError(fmt.Sprintf("Not supported prefix operator (%s) was used for boolean", operator))
}

func (engine *ExecutionEngine) EvalIntegerPrefixOperations(val *object.Integer, operator string) object.Object {

	if operator == token.MINUS {
		return &object.Integer{Value: -1 * val.Value}
	}

	return engine.createError(fmt.Sprintf("Not supported prefix operator (%s) was used for integer", operator))
}

func (engine *ExecutionEngine) EvalInfixExpression(infix *ast.InfixExpression) object.Object {
	left := engine.Eval(infix.Left)
	if engine.HasError {
		return left
	}

	right := engine.Eval(infix.Right)
	if engine.HasError {
		return right
	}

	if left.Type() != right.Type() {
		return engine.createError(
			fmt.Sprintf("Left and right variable share not the same type(%s and %s)", left.Type(), right.Type()),
		)
	}

	operator := infix.Operator
	if left.Type() == object.INTEGER_OBJ {
		return engine.EvalIntegerInfixOperations(left.(*object.Integer), right.(*object.Integer), operator)
	}

	return engine.createError(fmt.Sprintf("Not supported infix operator (%s) was used for type %s", operator, left.Type()))
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

	return engine.createError(fmt.Sprintf("Not supported infix operator (%s) was used for integers", operator))
}

func (engine *ExecutionEngine) createError(message string) *object.Error {
	engine.HasError = true
	return &object.Error{Message: message}
}
