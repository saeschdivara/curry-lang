package evaluator

import (
	"curryLang/lexer"
	"curryLang/object"
	"curryLang/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"5 > 3", true},
		{"9 < 7", false},
		{"9 == 7", false},
		{"9 == 9", true},
		{"9 != 7", true},
		{"9 != 9", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1 + 3", 4},
		{"1 - 3", -2},
		{"4 * 3", 12},
		{"12 / 3", 4},
		{"13 / 3", 4},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!false", true},
		{"!true", false},
		{"!!true", true},
		{"!!false", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-5", -5},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected object.Object
	}{
		{"if(true) { true }", &object.Boolean{Value: true}},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testObject(t, evaluated, tt.expected)
	}
}

func TestEvalLetStatement(t *testing.T) {
	l := lexer.New("let foo = 3;")
	p := parser.New(l)
	program := p.ParseProgram()
	engine := ExecutionEngine{}

	engine.Eval(program)

	if len(engine.Variables) != 1 {
		t.Errorf("Engine should contain 1 variable")
	}
}

func TestEvalIdentifierExpression(t *testing.T) {
	l := lexer.New("let foo = 3;foo;")
	p := parser.New(l)
	program := p.ParseProgram()
	engine := ExecutionEngine{}

	result := engine.Eval(program)

	if len(engine.Variables) != 1 {
		t.Errorf("Engine should contain 1 variable")
		return
	}

	intResult, ok := result.(*object.Integer)

	if !ok {
		t.Errorf("Result is not of type object.Integer but got %T", result)
		return
	}

	if intResult.Value != 3 {
		t.Errorf("Variable should contain 3")
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	engine := ExecutionEngine{}

	return engine.Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testObject(t *testing.T, obj object.Object, expected object.Object) bool {

	if obj == nil {
		t.Errorf("object is nil")
		return false
	}

	if obj.Type() != expected.Type() {
		t.Errorf("object is not %T. got=%T (%+v)", expected, obj, obj)
		return false
	}

	return false
}
