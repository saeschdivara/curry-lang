package evaluator

import (
	"curryLang/ast"
	"curryLang/object"
	"curryLang/token"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case *ast.IfElseExpression:
		return EvalIfElseExpression(node)
	case *ast.InfixExpression:
		return EvalInfixExpression(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.Program:
		return EvalStatements(node.Statements)
	}

	return &object.Null{}
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
		return &object.Null{}
	}

	condition := conditionResult.(*object.Boolean)

	if condition.Value {
		return EvalStatements(ifElse.Consequence)
	} else {
		return EvalStatements(ifElse.Alternative)
	}
}

func EvalInfixExpression(infix *ast.InfixExpression) object.Object {
	left := Eval(infix.Left)
	right := Eval(infix.Right)

	if left.Type() != right.Type() {
		return &object.Null{}
	}

	if left.Type() == object.INTEGER_OBJ {
		return EvalIntegerOperations(left.(*object.Integer), right.(*object.Integer), infix.Operator)
	}

	return &object.Null{}
}

func EvalIntegerOperations(left *object.Integer, right *object.Integer, operator string) object.Object {

	switch operator {
	case token.LT:
		return &object.Boolean{Value: left.Value < right.Value}
	case token.GT:
		return &object.Boolean{Value: left.Value > right.Value}
	case token.EQ:
		return &object.Boolean{Value: left.Value == right.Value}
	case token.NOT_EQ:
		return &object.Boolean{Value: left.Value != right.Value}
	}

	return &object.Null{}
}
