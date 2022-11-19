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

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.Program:
		return EvalStatements(node.Statements)
	}

	return nil
}

func EvalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}
