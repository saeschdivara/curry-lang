package evaluator

import (
	"curryLang/ast"
	"curryLang/object"
	"curryLang/token"
)

var (
	NULL = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case *ast.IfElseExpression:
		return EvalIfElseExpression(node)
	case *ast.PrefixExpression:
		return EvalPrefixExpression(node)
	case *ast.InfixExpression:
		return EvalInfixExpression(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.Program:
		return EvalStatements(node.Statements)
	}

	return NULL
}

func EvalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}

func EvalIfElseExpression(ifElse *ast.IfElseExpression) object.Object {
	conditionResult := Eval(ifElse.Condition)
	if conditionResult.Type() != object.BOOLEAN_OBJ {
		return NULL
	}

	condition := conditionResult.(*object.Boolean)

	if condition.Value {
		return EvalStatements(ifElse.Consequence)
	} else {
		return EvalStatements(ifElse.Alternative)
	}
}

func EvalPrefixExpression(infix *ast.PrefixExpression) object.Object {
	value := Eval(infix.Right)

	if value.Type() == object.BOOLEAN_OBJ {
		return EvalBooleanPrefixOperations(value.(*object.Boolean), infix.Operator)
	}
	if value.Type() == object.INTEGER_OBJ {
		return EvalIntegerPrefixOperations(value.(*object.Integer), infix.Operator)
	}

	return NULL
}

func EvalBooleanPrefixOperations(val *object.Boolean, operator string) object.Object {

	if operator == token.BANG {
		return &object.Boolean{Value: !val.Value}
	}

	return NULL
}

func EvalIntegerPrefixOperations(val *object.Integer, operator string) object.Object {

	if operator == token.MINUS {
		return &object.Integer{Value: -1 * val.Value}
	}

	return NULL
}

func EvalInfixExpression(infix *ast.InfixExpression) object.Object {
	left := Eval(infix.Left)
	right := Eval(infix.Right)

	if left.Type() != right.Type() {
		return NULL
	}

	if left.Type() == object.INTEGER_OBJ {
		return EvalIntegerInfixOperations(left.(*object.Integer), right.(*object.Integer), infix.Operator)
	}

	return NULL
}

func EvalIntegerInfixOperations(left *object.Integer, right *object.Integer, operator string) object.Object {

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
