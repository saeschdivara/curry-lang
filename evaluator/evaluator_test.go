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

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\"fllll{}\"", "fllll{}"},
		{"\"()=$$ääsas\"", "()=$$ääsas"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
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
	engine := NewEngine()

	engine.Eval(program)

	if len(engine.Variables) != 1 {
		t.Errorf("Engine should contain 1 variable")
	}
}

func TestEvalFunctionCall(t *testing.T) {
	result := testEval("fn foo() { 3; }; foo();")

	intResult, ok := result.(*object.Integer)

	if !ok {
		t.Errorf("Result is not of type object.Integer but got %T", result)
		return
	}

	if intResult.Value != 3 {
		t.Errorf("Variable should contain 3")
	}
}

func TestEvalReturnStatement(t *testing.T) {
	result := testEval("fn foo() { return 3; }; foo();")

	intResult, ok := result.(*object.Integer)

	if !ok {
		t.Errorf("Result is not of type object.Integer but got %T", result)
		return
	}

	if intResult.Value != 3 {
		t.Errorf("Variable should contain 3")
	}
}

func TestEvalNestedReturnStatements(t *testing.T) {
	result := testEval(`
	fn xxx() { 
	  return 3; 
	};
	fn foo() { 
	  return xxx(); 
	};
	fn bar(x) {

      if (x) {
	  	return foo();
	  }

      let f = 99;
      return f;
	};

	bar(true);
	`)

	intResult, ok := result.(*object.Integer)

	if !ok {
		t.Errorf("Result is not of type object.Integer but got %T", result)
		return
	}

	if intResult.Value != 3 {
		t.Errorf("Variable should contain 3")
	}
}

func TestEvalAnonymousFunctionCall(t *testing.T) {
	result := testEval("fn() { 3; }();")

	intResult, ok := result.(*object.Integer)

	if !ok {
		t.Errorf("Result is not of type object.Integer but got %T", result)
		return
	}

	if intResult.Value != 3 {
		t.Errorf("Variable should contain 3")
	}
}

func TestEvalFunctionCallWithParameters(t *testing.T) {
	result := testEval("fn foo(a, b) { a + b; }; foo(4, 5);")

	intResult, ok := result.(*object.Integer)

	if !ok {
		t.Errorf("Result is not of type object.Integer but got %T", result)
		return
	}

	if intResult.Value != 9 {
		t.Errorf("Variable should contain 3")
	}
}

func TestEvalIdentifierExpression(t *testing.T) {
	l := lexer.New("let foo = 3;foo;")
	p := parser.New(l)
	program := p.ParseProgram()
	engine := NewEngine()

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

func TestEvalWhileLoop(t *testing.T) {
	l := lexer.New("let foo = 3; while(foo < 10) { foo = foo + 1; }; foo;")
	p := parser.New(l)
	program := p.ParseProgram()
	engine := NewEngine()

	result := engine.Eval(program)

	intResult, ok := result.(*object.Integer)

	if !ok {
		t.Errorf("Result is not of type object.Integer but got %T", result)
		return
	}

	if intResult.Value != 10 {
		t.Errorf("Variable should contain 10")
	}
}

func TestEvalErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected object.Object
	}{
		{"foo;", &object.Error{Message: "Undeclared variable foo used"}},
		{"if(foo) {};", &object.Error{Message: "Undeclared variable foo used"}},
		{"if(true) {foo};", &object.Error{Message: "Undeclared variable foo used"}},
		{"fn() { foo; } ()", &object.Error{Message: "Undeclared variable foo used"}},
		{"fn() { foo; return 1; } ()", &object.Error{Message: "Undeclared variable foo used"}},
		{`
			fn test() {
				let x = fn() {
					foo;
				}

				if (true) {
					x();
				}

                return 1;
			}

			test();
		`, &object.Error{Message: "Undeclared variable foo used"}},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	engine := NewEngine()

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

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
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
