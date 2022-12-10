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

type Package struct {
	Name      string
	Globals   []Variable
	Functions map[string]*object.Function
}

type Module struct {
	Name     string
	Packages map[string]*Package
}

type ExecutionEngine struct {
	Variables       []Variable
	CurrentStackPos []uint32

	Functions map[string]*object.Function
	Modules   map[string]*Module

	// engine state flags
	IsReturnTriggered bool
	HasError          bool
}

func NewPackage(name string) *Package {
	return &Package{
		Name:      name,
		Globals:   make([]Variable, 0),
		Functions: map[string]*object.Function{},
	}
}

func NewModule(name string) *Module {
	return &Module{
		Name:     name,
		Packages: make(map[string]*Package),
	}
}

func NewEngine() *ExecutionEngine {
	engine := ExecutionEngine{IsReturnTriggered: false}
	engine.Variables = make([]Variable, 0)
	engine.CurrentStackPos = make([]uint32, 0)
	engine.Functions = make(map[string]*object.Function)
	engine.Modules = make(map[string]*Module)
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
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ListExpression:
		return engine.EvalListExpression(node)
	case *ast.Identifier:
		return engine.EvalIdentifier(node)

	case *ast.IfElseExpression:
		return engine.EvalIfElseExpression(node)
	case *ast.PrefixExpression:
		return engine.EvalPrefixExpression(node)
	case *ast.InfixExpression:
		return engine.EvalInfixExpression(node)
	case *ast.IndexAccessExpression:
		return engine.EvalIndexAccessExpression(node)

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

func (engine *ExecutionEngine) EvalListExpression(identifier *ast.ListExpression) object.Object {
	obj := &object.List{}
	obj.Value = make([]object.Object, 0)

	for i, valExpr := range identifier.Value {
		val := engine.Eval(valExpr)
		if i == 0 {
			obj.ValueType = val.Type()
		} else {
			if val.Type() != obj.ValueType {
				return engine.createError(
					fmt.Sprintf(
						"List members have to be all of the same type, value #%v has type %s instead of %s",
						i,
						val.Type(),
						obj.Type(),
					),
				)
			}
		}

		obj.Value = append(obj.Value, val)
	}

	return obj
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

	leftType := left.Type()
	rightType := right.Type()

	if leftType != right.Type() {

		if leftType == object.STRING_OBJ && rightType == object.INTEGER_OBJ {
			strVal := left.(*object.String)
			intVal := right.(*object.Integer)

			return &object.String{Value: fmt.Sprintf("%s%v", strVal.Value, intVal.Value)}
		}

		if leftType == object.INTEGER_OBJ && rightType == object.STRING_OBJ {
			intVal := left.(*object.Integer)
			strVal := right.(*object.String)

			return &object.String{Value: fmt.Sprintf("%v%s", intVal.Value, strVal.Value)}
		}

		return engine.createError(
			fmt.Sprintf("Left and right variable share not the same type(%s and %s)", leftType, rightType),
		)
	}

	operator := infix.Operator

	if leftType == object.INTEGER_OBJ {
		return engine.EvalIntegerInfixOperations(left.(*object.Integer), right.(*object.Integer), operator)
	}
	if leftType == object.STRING_OBJ {
		return engine.EvalStringInfixOperations(left.(*object.String), right.(*object.String), operator)
	}

	return engine.createError(fmt.Sprintf("Not supported infix operator (%s) was used for type %s", operator, leftType))
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

func (engine *ExecutionEngine) EvalStringInfixOperations(left *object.String, right *object.String, operator string) object.Object {

	switch operator {
	// logical operations
	//case token.LT:
	//	return &object.Boolean{Value: left.Value < right.Value}
	//case token.GT:
	//	return &object.Boolean{Value: left.Value > right.Value}
	//case token.EQ:
	//	return &object.Boolean{Value: left.Value == right.Value}
	//case token.NOT_EQ:
	//	return &object.Boolean{Value: left.Value != right.Value}

	case token.PLUS:
		return &object.String{Value: left.Value + right.Value}
	}

	return engine.createError(fmt.Sprintf("Not supported infix operator (%s) was used for integers", operator))
}

func (engine *ExecutionEngine) EvalIndexAccessExpression(indexAccess *ast.IndexAccessExpression) object.Object {

	indexExpr := engine.Eval(indexAccess.Value)

	if indexExpr.Type() == object.ERROR_OBJ {
		return indexExpr
	}

	if indexExpr.Type() != object.INTEGER_OBJ {
		return engine.createError(fmt.Sprintf("Index type has to be integer but is %s", indexExpr.Type()))
	}

	sourceExpr := engine.Eval(indexAccess.Source)

	if sourceExpr.Type() == object.ERROR_OBJ {
		return sourceExpr
	}

	if sourceExpr.Type() != object.LIST_OBJ {
		return engine.createError(fmt.Sprintf("Source type has to be list but is %s", sourceExpr.Type()))
	}

	indexObj := indexExpr.(*object.Integer)
	sourceList := sourceExpr.(*object.List)

	if int(indexObj.Value) >= len(sourceList.Value) {
		return engine.createError(
			fmt.Sprintf("List is too small (%v) for index %v", len(sourceList.Value), indexObj.Value),
		)
	}

	return sourceList.Value[indexObj.Value]
}

func (engine *ExecutionEngine) createError(message string) *object.Error {
	engine.HasError = true
	return &object.Error{Message: message}
}
