package parser

import (
	"curryLang/ast"
	"curryLang/lexer"
	"fmt"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
   let x = 5;
   let y = 10;
   let foobar = 838383;
   `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedExpression string
	}{
		{"x", "5"},
		{"y", "10"},
		{"foobar", "838383"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedExpression) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
   return 5;
   return 10;
   return 993322;
   `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestWhileStatements(t *testing.T) {
	input := `
   while (foo) {
		let x = 10;
	}
   `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.WhileStatement. got=%T", program.Statements[0])
	}

	if stmt.Condition.String() != "foo" {
		t.Fatalf("Expected stmt.Condition.String() to be foo but was %s", stmt.Condition.String())
	}

	if len(stmt.Body) != 1 {
		t.Fatalf("Expected len(stmt.Body) to be 1 but was %v", len(stmt.Body))
	}
}

func TestVariableAssignment(t *testing.T) {
	input := "x = 10;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.String() != "x" {
		t.Fatalf("Expected stmt.Name.String() to be x but was %s", stmt.Name.String())
	}

	if stmt.Value.String() != "10" {
		t.Fatalf("Expected stmt.Value.String() to be 10 but was %s", stmt.Value.String())
	}
}

func TestStringLiteral(t *testing.T) {
	input := "x = \"f{}dfds)=\";"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.String() != "x" {
		t.Fatalf("Expected stmt.Name.String() to be x but was %s", stmt.Name.String())
	}

	if stmt.Value.String() != "f{}dfds)=" {
		t.Fatalf("Expected stmt.Value.String() to be \"f{}dfds)=\" but was %s", stmt.Value.String())
	}
}

func TestIntegerListExpression(t *testing.T) {
	input := "x = [10, 50];"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.String() != "x" {
		t.Fatalf("Expected stmt.Name.String() to be x but was %s", stmt.Name.String())
	}

	if stmt.Value.String() != "[10, 50]" {
		t.Fatalf("Expected stmt.Value.String() to be [10, 50] but was %s", stmt.Value.String())
	}
}

func TestStringListExpression(t *testing.T) {
	input := "x = [\"foo\", \"bar\"];"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.String() != "x" {
		t.Fatalf("Expected stmt.Name.String() to be x but was %s", stmt.Name.String())
	}

	if stmt.Value.String() != "[foo, bar]" {
		t.Fatalf("Expected stmt.Value.String() to be [foo, bar] but was %s", stmt.Value.String())
	}
}

func TestListListsExpression(t *testing.T) {
	input := "x = [[\"foo\", \"bar\"], [1, 20, 6]];"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.String() != "x" {
		t.Fatalf("Expected stmt.Name.String() to be x but was %s", stmt.Name.String())
	}

	if stmt.Value.String() != "[[foo, bar], [1, 20, 6]]" {
		t.Fatalf("Expected stmt.Value.String() to be [foo, bar] but was %s", stmt.Value.String())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestBooleanLiteralExpression(t *testing.T) {
	booleanTests := []struct {
		input        string
		booleanValue bool
	}{
		{"false;", false},
		{"true;", true},
	}
	for _, tt := range booleanTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("stmt is not ast.Boolean. got=%T", stmt.Expression)
		}

		if tt.booleanValue != exp.Value {
			return
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b);",
		}, {

			"!-a",
			"(!(-a));",
		},
		{
			"a + b + c",
			"((a + b) + c);",
		},
		{
			"a + b - c",
			"((a + b) - c);",
		},
		{
			"a * b * c",
			"((a * b) * c);",
		},
		{
			"a * b / c",
			"((a * b) / c);",
		},
		{
			"a + b / c",
			"(a + (b / c));",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f);",
		}, {
			"3 + 4; -5 * 5",
			"(3 + 4);((-5) * 5);",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4));",
		}, {
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4));",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		},
		{
			"true",
			"true;",
		},
		{
			"false",
			"false;",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false);",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true);",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4);",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2);",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5));",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5));",
		},
		{
			"!(true == true)",
			"(!(true == true));",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParsingIfElseExpressions(t *testing.T) {
	infixTests := []struct {
		input       string
		condition   string
		consequence string
		alternative string
	}{
		{"if (5 < 6) { true } else { false }", "(5 < 6)", "true;", "false;"},
		{"if (5 < 6) { true }", "(5 < 6)", "true;", ""},
		{"if 5 < 6 { true }", "(5 < 6)", "true;", ""},
		{`
	if (test == 5) {
		let x = 5;
	}
`, "(test == 5)", "let x = 5;", ""},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			return
		}

		expr, ok := stmt.Expression.(*ast.IfElseExpression)
		if !ok {
			t.Fatalf("Expression is not ast.IfElseExpression. got=%T", stmt.Expression)
			return
		}

		if expr.Condition.String() != tt.condition {
			t.Errorf("expr.Condition.String() not '%s'. got=%s", tt.condition, expr.Condition.String())
		}

		if expr.ConsequenceString() != tt.consequence {
			t.Errorf("expr.ConsequenceString() not '%s'. got=%s", tt.consequence, expr.ConsequenceString())
		}

		if expr.AlternativeString() != tt.alternative {
			t.Errorf("expr.AlternativeString() not '%s'. got=%s", tt.alternative, expr.AlternativeString())
		}
	}
}

func TestParsingFunctionExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		name       string
		parameters string
		body       string
	}{
		{"fn(x, y) { test; }", "", "(x, y)", "test;"},
		{"fn foo(x, y) { test; }", "foo", "(x, y)", "test;"},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			return
		}

		expr, ok := stmt.Expression.(*ast.FunctionExpression)
		if !ok {
			t.Fatalf("Expression is not ast.FunctionExpression. got=%T", stmt.Expression)
			return
		}

		if expr.Name != tt.name {
			t.Fatalf("expr.Name not '%s'. got=%s", tt.name, expr.Name)
			return
		}

		parametersOutput := expr.ParametersString()
		if parametersOutput != tt.parameters {
			t.Fatalf("expr.ParametersString() not '%s'. got=%s", tt.parameters, parametersOutput)
			return
		}

		functionBodyOutput := expr.BodyString()
		if functionBodyOutput != tt.body {
			t.Fatalf("expr.BodyString() not '%s'. got=%s", tt.body, functionBodyOutput)
			return
		}
	}
}

func TestFunctionCallExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		name       string
		parameters string
	}{
		{"add(x, y);", "add", "(x, y)"},
		{"add(x+1, y*2);", "add", "((x + 1), (y * 2))"},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			return
		}

		expr, ok := stmt.Expression.(*ast.FunctionCallExpression)
		if !ok {
			t.Fatalf("Expression is not ast.FunctionExpression. got=%T", stmt.Expression)
			return
		}

		if expr.FunctionExpr.String() != tt.name {
			t.Fatalf("expr.FunctionExpr.String() not '%s'. got=%s", tt.name, expr.FunctionExpr.String())
			return
		}

		if expr.ParametersString() != tt.parameters {
			t.Fatalf("expr.ParametersString() not '%s'. got=%s", tt.parameters, expr.ParametersString())
			return
		}
	}
}

func TestParsingComplexStatements(t *testing.T) {
	input := `
	if (test == 5) {
		let x = if (f < 5) {
			let t = 33;
			40;
		};
	} else {
		let x = if (55 != 5) {
			40;
		} else {
			aaaa;
		};
	}
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		return
	}

	expr, ok := stmt.Expression.(*ast.IfElseExpression)
	if !ok {
		t.Fatalf("Expression is not ast.IfElseExpression. got=%T", stmt.Expression)
		return
	}

	if expr.Condition.String() == "test == 5" {
		t.Errorf("expr.Condition.String() not '%s'. got=%s", "test == 5", expr.Condition.String())
		return
	}
}

func TestParsingAnonymousFunctionCall(t *testing.T) {
	input := `
	fn() {} ();
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		return
	}

	expr, ok := stmt.Expression.(*ast.FunctionCallExpression)
	if !ok {
		t.Fatalf("Expression is not ast.FunctionCallExpression. got=%T", stmt.Expression)
		return
	}

	funcExpr, ok := expr.FunctionExpr.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("Expression is not ast.FunctionExpression. got=%T", expr.FunctionExpr)
		return
	}

	if funcExpr.Name != "" {
		t.Fatalf("funcExpr.Name is not anonymous function. got=%s\n", funcExpr.Name)
	}
}

func TestParsingComplexFunctionStatements(t *testing.T) {
	input := `
	if (test == 5) {
		let x = fn() {
			let foo = fn test(f, t, k) {};
		};
	}
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		return
	}

	expr, ok := stmt.Expression.(*ast.IfElseExpression)
	if !ok {
		t.Fatalf("Expression is not ast.IfElseExpression. got=%T", stmt.Expression)
		return
	}

	if expr.Condition.String() == "test == 5" {
		t.Errorf("expr.Condition.String() not '%s'. got=%s", "test == 5", expr.Condition.String())
		return
	}

	ifBodyStatements := expr.Consequence
	if len(ifBodyStatements) != 1 {
		t.Errorf("len(ifBodyStatements) not '%s'. got=%v", "1", len(ifBodyStatements))
		return
	}

	letStatement, ok := ifBodyStatements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("Expression is not ast.letStatement. got=%T", ifBodyStatements[0])
		return
	}

	if letStatement.Name.String() != "x" {
		t.Errorf("letStatement.Name.String() not '%s'. got=%s", "x", letStatement.Name.String())
		return
	}

	outerFuncExpr, ok := letStatement.Value.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("Expression is not ast.FunctionExpression. got=%T", outerFuncExpr)
		return
	}

	if outerFuncExpr.Name != "" {
		t.Errorf("outerFuncExpr.Name not '%s'. got=%s", "x", outerFuncExpr.Name)
		return
	}

	if len(outerFuncExpr.Parameters) != 0 {
		t.Errorf("len(outerFuncExpr.Parameters) not '%v'. got=%v", 0, len(outerFuncExpr.Parameters))
		return
	}

	if len(outerFuncExpr.Body) != 1 {
		t.Errorf("len(outerFuncExpr.Body) not '%v'. got=%v", 1, len(outerFuncExpr.Body))
		return
	}

	innerLetStatement, ok := outerFuncExpr.Body[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("Expression is not ast.LetStatement. got=%T", outerFuncExpr.Body[0])
		return
	}

	if innerLetStatement.Name.String() != "foo" {
		t.Errorf("innerLetStatement.Name.String() not '%s'. got=%s", "x", innerLetStatement.Name.String())
		return
	}

	innerFunctionExpr, ok := innerLetStatement.Value.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("Expression is not ast.FunctionExpression. got=%T", innerFunctionExpr)
		return
	}

	if innerFunctionExpr.Name != "test" {
		t.Errorf("innerFunctionExpr.Name not '%s'. got=%s", "(f, t, k)", innerFunctionExpr.ParametersString())
		return
	}

	if innerFunctionExpr.ParametersString() != "(f, t, k)" {
		t.Errorf("innerFunctionExpr.ParametersString() not '%s'. got=%s", "(f, t, k)", innerFunctionExpr.ParametersString())
		return
	}

	if len(innerFunctionExpr.Body) != 0 {
		t.Errorf("len(innerFunctionExpr.Body) not '%v'. got=%v", 0, len(innerFunctionExpr.Body))
		return
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string, value string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}
	if letStmt.Value.String() != value {
		t.Errorf("letStmt.Value.String() not '%s'. got=%s", value, letStmt.Value.String())
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
