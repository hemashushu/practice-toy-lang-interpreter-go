// original from https://interpreterbook.com/

package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for i, msg := range errors {
		t.Errorf("parser error #%d: %q",
			i, msg)
	}

	t.FailNow()
}

func TestLetStatements(t *testing.T) {

	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x= 1;", "x", 5},
		{"let y= 2;", "y", 2},
		{"let foobar = 1234;", "foobar", 1234},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, actual %d", len(program.Statements))
		}

		statement := program.Statements[0]
		if !testLetStatement(t, statement, test.expectedIdentifier) {
			return
		}

		// TODO::

		// letStatement, ok := statement.(*ast.LetStatement)
		// if !ok {
		// 	t.Fatalf("expected *ast.LetStatement, actual %T", statement)
		// }

		// valueExpression := letStatement.Value
		// if !testLiteralExpression(t, valueExpression, test.expectedValue) {
		// 	return
		// }
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, identifierName string) bool {
	if statement.TokenLiteral() != "let" { // LET token 本身
		t.Errorf("TokenLiteral expected 'let', actual '%q'", statement.TokenLiteral())
		return false
	}

	// x.(T) 为类型断言/类型转换，将接口类型的值转换为具体类型
	// 相当于 Java:
	// if (statement instanceof LetStatement) {
	//     LetStatement letStatement = (LetStatement)statement;
	//     ...
	// }
	//
	// 注意这里的 statement 是 *LetStatement
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected *ast.LetStatement, actual %T", statement)
		return false
	}

	if letStatement.Name.Value != identifierName { // .Name 是 *Identifier
		t.Errorf("letStatement.Name.Value expected %q, actual %q",
			identifierName, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != identifierName { // .Name 是 *Identifier
		t.Errorf("letStatement.Name expected %q, actual %q",
			identifierName, letStatement.Name)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	// input := `
	// 	return 1;
	// 	return 23;
	// 	return 456;
	// 	`

	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 1;", 1},
		{"return 23;", 23},
		{"return 456;", 456},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		// if len(program.Statements) != 3 {
		// 	t.Fatalf("expected 3 statements, actual %d",
		// 		len(program.Statements))
		// }

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, actual %d", len(program.Statements))
		}

		// for _, statement := range program.Statements {
		statement := program.Statements[0]

		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("expected *ast.returnStatement, actual %T", statement)
			// continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Fatalf("returnStatement.TokenLiteral expected 'return', actual %q",
				returnStatement.TokenLiteral())
		}

		// TODO::
		// returnValueExpression := returnStatement.ReturnValue
		// if !testLiteralExpression(t, returnValueExpression, test.expectedValue) {
		// 	return
		// }

		// }
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement. actual %d",
			len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] expected ast.ExpressionStatement, actual %T",
			program.Statements[0])
	}

	testIdentifier(t, statement.Expression, "foobar")
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expected *ast.Identifier, actual %T", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("ident.Value expected %q, actual %q", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral expected %q, actual %q", "foobar",
			identifier.TokenLiteral())
		return false
	}

	return true
}

func TestIntegerLiteralExpression(t *testing.T) {
	// input := "5;"

	tests := []struct {
		input         string
		expectedValue int64
	}{
		{"5;", 5},
		{"1234;", 1234},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, actual %d",
				len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] expected ast.ExpressionStatement, actual %T",
				program.Statements[0])
		}

		if !testLiteralExpression(t, statement.Expression, test.expectedValue) {
			return
		}

		// 	literal, ok := statement.Expression.(*ast.IntegerLiteral)
		// 	if !ok {
		// 		t.Fatalf("expected *ast.IntegerLiteral, actual %T", statement.Expression)
		// 	}
		//
		// 	if literal.Value != 5 {
		// 		t.Errorf("literal.Value expected %d, actual %d", 5, literal.Value)
		// 	}
		//
		// 	if literal.TokenLiteral() != "5" {
		// 		t.Errorf("literal.TokenLiteral expected %q, actual %q", "5",
		// 			literal.TokenLiteral())
		// 	}
	}

}

func TestBooleanLiteralExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, actual %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected ast.ExpressionStatement, actual %T",
				program.Statements[0])
		}

		//		expression := statement.Expression
		// 		boolLiteral, ok := expression.(*ast.Boolean)
		// 		if !ok {
		// 			t.Fatalf("expected *ast.Boolean, actual %T", expression)
		// 		}
		//
		// 		if boolLiteral.Value != test.value {
		// 			t.Fatalf("expected %t, actual %t", test.value, boolLiteral.Value)
		// 		}

		if !testLiteralExpression(t, statement.Expression, test.expectedValue) {
			return
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input         string
		operator      string
		expectedValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, test := range prefixTests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected %d statement, actualy %d\n",
				1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] expected ast.ExpressionStatement, actual %T",
				program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("statement expected ast.PrefixExpression, actual %T", statement.Expression)
		}

		if expression.Operator != test.operator {
			t.Fatalf("exp.Operator expected %q, actual %q",
				test.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, test.expectedValue) {
			return
		}
	}
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expectedValue interface{}) bool {

	switch v := expectedValue.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(v))
	case int64:
		return testIntegerLiteral(t, expression, v)
	case string:
		return testIdentifier(t, expression, v)
	case bool:
		return testBooleanLiteral(t, expression, v)
	}

	t.Errorf("unexpected expression type %T", expression)
	return false
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, expectedValue int64) bool {
	intLiteral, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expression expected *ast.IntegerLiteral, actual %T", expression)
		return false
	}

	if intLiteral.Value != expectedValue {
		t.Errorf("intLiteral.Value expected %d, actual %d", expectedValue, intLiteral.Value)
		return false
	}

	if intLiteral.TokenLiteral() != fmt.Sprintf("%d", expectedValue) {
		t.Errorf("intLiteral.TokenLiteral excepted %d, actual %q", expectedValue,
			intLiteral.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, expectedValue bool) bool {
	boolLiteral, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("expected *ast.Boolean, actual %T", expression)
		return false
	}

	if boolLiteral.Value != expectedValue {
		t.Errorf("expected %t, actual %t", expectedValue, boolLiteral.Value)
		return false
	}

	if boolLiteral.TokenLiteral() != fmt.Sprintf("%t", expectedValue) {
		t.Errorf("boolLiteral.TokenLiteral expected %t, actual %q",
			expectedValue, boolLiteral.TokenLiteral())
		return false
	}

	return true
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

	for _, test := range infixTests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected %d statement, actual %d\n",
				1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected ast.ExpressionStatement, actual %T",
				program.Statements[0])
		}

		if !testInfixExpression(t, statement.Expression,
			test.leftValue, test.operator, test.rightValue) {
			return
		}
	}
}

func testInfixExpression(t *testing.T, expression ast.Expression,
	leftValue interface{}, operator string, rightValue interface{}) bool {

	operatorExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expected ast.OperatorExpression, actual %T", expression)
		return false
	}

	if !testLiteralExpression(t, operatorExpression.Left, leftValue) {
		return false
	}

	if operatorExpression.Operator != operator {
		t.Errorf("exp.Operator expected %q, actual %q", operator, operatorExpression.Operator)
		return false
	}

	if !testLiteralExpression(t, operatorExpression.Right, rightValue) {
		return false
	}

	return false
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != test.expectedValue {
			t.Fatalf("expected %q, actual %q", test.expectedValue, actual)
		}
	}
}
