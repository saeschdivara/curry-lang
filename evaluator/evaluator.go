package evaluator

import (
	"curryLang/ast"
	"curryLang/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case *ast.IfElseExpression:
		return EvalIfElseExpression(node)

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
