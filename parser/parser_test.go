/**
 * Copyright (c) 2022 Hemashushu <hippospark@gmail.com>, All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x= 1;
	let y= 2;
	let foobar = 1234;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	// if program == nil {
	// 	t.Fatalf("ParseProgram() return null")
	// }

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, actual %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, test.expectedIdentifier) {
			return
		}
	}
}

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

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
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

	if letStatement.Name.Value != name { // .Name 是 *Identifier
		t.Errorf("letStatement.Name.Value expected %q, actual %q",
			name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name { // .Name 是 *Identifier
		t.Errorf("letStatement.Name expected %q, actual %q",
			name, letStatement.Name)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 1;
		return 23;
		return 456;
		`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, actual %d",
			len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("expected *ast.returnStatement, actual %T", statement)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral expected 'return', actual %q",
				returnStatement.TokenLiteral())
		}
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

	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier, actual %T", statement.Expression)
	}

	if identifier.Value != "foobar" {
		t.Errorf("ident.Value expected %q, actual %q", "foobar", identifier.Value)
	}

	if identifier.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral expected %q, actual %q", "foobar",
			identifier.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
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

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, actual %T", statement.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value expected %d, actual %d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral expected %q, actual %q", "5",
			literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
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

		if !testIntegerLiteral(t, expression.Right, test.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) bool {
	intLiteral, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expression expected *ast.IntegerLiteral, actual %T", expression)
		return false
	}

	if intLiteral.Value != value {
		t.Errorf("intLiteral.Value expected %d, actual %d", value, intLiteral.Value)
		return false
	}

	if intLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("intLiteral.TokenLiteral excepted %d, actual %q", value,
			intLiteral.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
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

		expression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expected ast.InfixExpression, actual %T", statement.Expression)
		}

		if !testIntegerLiteral(t, expression.Left, test.leftValue) {
			return
		}

		if expression.Operator != test.operator {
			t.Fatalf("exp.Operator expected %q, actual %q",
				test.operator, expression.Operator)
		}

		if !testIntegerLiteral(t, expression.Right, test.rightValue) {
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
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != test.expected {
			t.Errorf("expected %q, actual %q", test.expected, actual)
		}
	}
}
