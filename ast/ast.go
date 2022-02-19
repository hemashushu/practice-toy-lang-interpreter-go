// original from https://interpreterbook.com/

package ast

import (
	"bytes"
	"interpreter/token"
)

type Node interface {
	// 节点/token 的原始值（字符串），仅用于调试、测试
	// 具体的节点，比如 literal 有对应的值，比如 IntegerLiteral 对应的是一个整数
	TokenLiteral() string

	// 用于打印 AST，类似 toString()
	String() string
}

type Expression interface {
	Node
	expressionNode() // 无实际用途的方法，仅用于将 struct 分类
}

type Statement interface {
	Node
	statementNode() // 无实际用途的方法，仅用于将 struct 分类
}

type Program struct {
	Statements []Statement // 注意这里 Statement 是接口，而不是结构体，所以可以装 *LetStatement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral() // 只返回第一个语句的 token 值
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct {
	Token token.Token // let 语句的开始 token，必定是 LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {
	// 实现接口 Statement 的方法 statementNode()
	// 用于表明 LetStatement 是一个 Statement
	// 相当于 Rust 的 `#[derive]`` 或者 Java 的 `class X implements ISomeInterface`
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type Identifier struct {
	Token token.Token // the IDENT token
	Value string
}

func (i *Identifier) expressionNode() {
	// 实现接口 Expression 的方法 expressionNode()
	// 用于表明 Identifier 是一个 Expression
	// 相当于 Rust 的 `#[derive]`` 或者 Java 的 `class X implements ISomeInterface`
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type ReturnStatement struct {
	Token       token.Token // return 语句的开始 token，必定是 RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// 一元运算符
type PrefixExpression struct {
	Token    token.Token // 运算符, e.g. !, -, +
	Operator string      // 运算符的符号
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // 运算符 token, e.g. + - * / > < == !=
	Operator string      // 运算符的符号
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// BooleanLiteral
type Boolean struct {
	Token token.Token
	Value bool
}

func (il *Boolean) expressionNode()      {}
func (il *Boolean) TokenLiteral() string { return il.Token.Literal }
func (il *Boolean) String() string       { return il.Token.Literal }
