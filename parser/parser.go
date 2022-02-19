// original from https://interpreterbook.com/

package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
	"strconv"
)

// 表达式运算符的优先级别列表
const (
	_           int = iota
	LOWEST          // 最低优先级，比如从 “语句” 进来的 "表达式" 解析阶段。
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	PREFIX          // -X or !X
	CALL            // myFunction(X)
)

// 各个运算符 token 对应的优先级
var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS, // ==
	token.NOT_EQ: EQUALS, // "!="

	token.LT: LESSGREATER, // <
	token.GT: LESSGREATER, // >

	token.PLUS:     SUM,     // +
	token.MINUS:    SUM,     // -
	token.SLASH:    PRODUCT, // /
	token.ASTERISK: PRODUCT, // *
}

// 解析的顺序按照各种 "语言元素（即语句、表达式等）" 的优先级来进行。
//
// 1. 先解析各种 "语句"，比如变量声明语句，函数定义语句等，
//    语句之间是并排关系，没有优先级之分。
// 2. "语句" 当中有 "表达式语句" 一类
// 3. 解析 "表达式语句" 当中的 "表达式"
//    a. 赋值表达式（假如语言允许连续赋值）
//    b. 逻辑 or
//    c. 逻辑 and
//    d. 相等比较（==, !=）
//    e. 大小比较（>, <, >=, <=）
//	  f. 加减（+, -）
//    g. 乘除（*, /）
//
//    (以上是二元运算/表达式，以下可以视为是一元运算/表达式)
//
//    h. 一元运算（正负数，逻辑非）
//	  i. 对象成员或者函数调用（obj.prop, obj[index], func(...)）
//    j. 基础表达式
//
//    基础表达式包括字面量（包括元组、列表、字典等字面量）、括号、标识符、new 表达式
//    基础表达式单独出现，所以没有先后顺序。
//    括号、元组、列表、字典当中允许任何 "表达式"，所以又会回到第 3 步。
//
//    注意如果语言支持 new 表达式的话：
//    new (...) 表达式的优先级要比成员表达式的高，
//    即 `new a(...).c(...)` 是 `(new a(...)).c(...)`
//    new ... 表达式跟成员表达式一样，
//    即 `new a.b.c(...)` 是 `new (a.b.c)()`

type prefixParseFn func() ast.Expression              // Unary operator
type infixParseFn func(ast.Expression) ast.Expression // Binary operator

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token // current token
	peekToken token.Token // next token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 读两次，让 current token 和 peek token 都赋予值
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	// 注册 primary 表达式（字面量、标识符等）解析过程
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)

	// 注册一元操作符解析过程
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// 注册二元操作符解析过程
	p.registerInfix(token.PLUS, p.parseInfixExpression)     // +
	p.registerInfix(token.MINUS, p.parseInfixExpression)    // -
	p.registerInfix(token.SLASH, p.parseInfixExpression)    // /
	p.registerInfix(token.ASTERISK, p.parseInfixExpression) // *
	p.registerInfix(token.EQ, p.parseInfixExpression)       // ==
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)   // "!="
	p.registerInfix(token.LT, p.parseInfixExpression)       // <
	p.registerInfix(token.GT, p.parseInfixExpression)       // >

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token type %q, actual %q",
		t,
		p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 断言并消耗指定的 type 的 token
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		statement := p.parseStatement()
		// if statement != nil {
		program.Statements = append(program.Statements, statement)
		// }

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	statement.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// 检测到 ";" 就退出，并不消耗 ";" 符号
	// i.e. 当前 token 停留在 ';' 位置
	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.nextToken()

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// 当前 token 停留在 ';' 位置

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{
		Token: p.curToken,
	}

	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) { // 让当前 token 移动到 ';' 位置
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {

	// 先解析 primary expression （字面量、标识符等）和一元运算符
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// 比较式 precedence < p.peekPrecedence() 表示：
		// "下一个运算符" 的优先级比 "预想的" 要高，
		// "预想的" 是指调用 "parseExpression" 时，当前所处的优先级，一旦进入
		// "parseExpression" 阶段，所有比 "预想的" 优先级要高的连续 "运算符" 都会解析，
		// 直到碰到跟预想的优先级一样的，或者更低的，才会停止。
		//
		// 比如解析 "1+2+3"
		// 1. 一开始从 LOWEST 开始，解析了字面量 "1"，置为 left，
		// 注：所有 "语句表达式" 开始之前 "预想" 的都是最低优先级 LOWEST
		// 2. 然后在这里遇到了 "+" 运算符，"+" 的优先级比 LOWEST 高，
		// 3. 将 left (literal:1) 带入 infix，infix 构建 InfixExpression，消耗了 "+" 运算符
		// 4. 然后 infix 调用 parseExpression("+运算符的优先级") ，并准备将返回值作为 right
		//
		// 5. 解析了字面量 "2",置为 left，
		// 6. 然后在这里遇到了 "+" 运算符，"+" 的优先级跟 "+" 一致，
		// 7. parseExpression 返回 (literal:2)
		// 8. infix 返回 InfixExpression(1 "+" 2)
		//
		// 9. parseExpression 将 (1 "+" 2) 置为 left，然后再次查找下一个 token 的优先级
		// 10. 然后在这里遇到了 "+" 运算符，"+" 的优先级比 LOWEST 高，
		// 11. 将 (1 "+" 2) 带入 infox, infix 构建 InfixExpression，消耗了 "+" 运算符
		// 12. 然后 infix 调用 parseExpression("+运算符的优先级") ，并准备将返回值作为 right
		//
		// 13. parseExpression 返回 left (literal:3)
		// 14. infix 返回 InfixExpression ((1 "+" 2) "+" 3)
		// 15. parseExpression 返回 ((1 "+" 2) "+" 3)
		//
		// 比如解析 "1+2*3"，
		// 1. 一开始从 LOWEST 开始，解析了字面量 "1"，置为 left，
		// 2. 然后在这里遇到了 "+" 运算符，"+" 的优先级比 LOWEST 高，
		// 3. 将 left (literal:1) 带入 infix，infix 构建 InfixExpression，消耗了 "+" 运算符
		// 4. 然后 infix 调用 parseExpression("+运算符的优先级") ，并准备将返回值作为 right
		//
		// 5. 解析了字面量 "2",置为 left，
		// 6. 然后在这里遇到了 "*" 运算符，"*" 的优先级比 "+" 高，
		// 7. 将 left (literal:2) 带入 infox, infix 构建 InfixExpression，消耗了 "*" 运算符
		// 8. 然后 infix 调用 parseExpression("*运算符的优先级") ，并准备将返回值作为 right
		//
		// 9. parseExpression 返回 left (literal:3)
		// 10. infix 返回 InfixExpression(2 "*" 3)
		// 11. infix 返回 InfixExpression(1 "+" ...)
		// 12. parseExpression 返回 (1 "+" (2 "*" 3))

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken() // 消耗掉当前的 token

		leftExp = infix(leftExp)
	}

	return leftExp
}

// 查找当前 token 的运算符优先级别（假如存在的话，否则返回 LOWEST）
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

// 查找下一个 token 的运算符优先级别（假如存在的话，否则返回 LOWEST）
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{
		Token: p.curToken,
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = value
	return literal
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	literal := &ast.Boolean{
		Token: p.curToken,
	}

	// value, err := strconv.ParseBool(p.curToken.Literal)
	// if err != nil {
	// 	msg := fmt.Sprintf("could not parse %q as bool", p.curToken.Literal)
	// 	p.errors = append(p.errors, msg)
	// 	return nil
	// }
	//literal.Value = value

	literal.Value = p.curTokenIs(token.TRUE) // 因为只有 token.TRUE 和 token.FALSE 两种情况
	return literal
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %q found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()

	// 如果这里让 parseExpression(precedence -1) 可以实现
	// 同一个运算符实现右->左结合
	expression.Right = p.parseExpression(precedence)
	return expression
}
