// original from https://interpreterbook.com/

package evaluator

import (
	"interpreter/ast"
	"interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {

	// 对语句求值
	case *ast.Program:
		return evalProgram(n)

	case *ast.BlockStatement:
		return evalBlockStatement(n)

	case *ast.ExpressionStatement:
		return Eval(n.Expression)

	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue)
		return &object.ReturnValue{Value: val} // 包裹待返回的 Object

	// 对表达式求值
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		return evalPrefixExpression(n.Operator, right)

	case *ast.InfixExpression:
		left := Eval(n.Left)
		right := Eval(n.Right)
		return evalInfixExpression(n.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(n)

	// 对标识符求值

	// 对字面量求值
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(n.Value)

	}

	return nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalProgram(program *ast.Program) object.Object {
	// evalProgram 跟 evalBlockStatement 很相似，但 BlockStatement 可以嵌套，当遇到
	// return 语句时，需要跳到最外一层 block，所以无法重用 evalBlockStatement
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		// 在一组语句中，存在 return 语句
		// if returnValue, ok := result.(*object.ReturnValue); ok {
		// 	return returnValue.Value
		// }

		// 因为 BlockStatement 可以嵌套，所以暂时不拆封 ReturnValue
		// 直到 Program 或者 函数字面量 才解封
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	// 返回最后一条语句的值
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "+":
		return evalPlusPrefixOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {

	/*
		// 对于数字，0 视为 false，非 0 视为 true
		if right.Type() == object.INTEGER_OBJ {
			value := right.(*object.Integer).Value
			if value == 0 {
				return TRUE
			} else {
				return FALSE
			}
		}
	*/

	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE // 所有非 false, null 的值都视为 true
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL // note:: 可返回错误
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalPlusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL // note:: 可返回错误
	}

	return right
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ &&
		right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	case operator == "&&":
		return nativeBoolToBooleanObject(
			left.(*object.Boolean).Value &&
				right.(*object.Boolean).Value)

	case operator == "||":
		return nativeBoolToBooleanObject(
			left.(*object.Boolean).Value ||
				right.(*object.Boolean).Value)
	default:
		return NULL // note:: 可返回错误
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}

	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)

	default:
		return NULL // note:: 可返回错误
	}
}

func evalIfExpression(expression *ast.IfExpression) object.Object {
	condition := Eval(expression.Condition)

	if isTruthy(condition) {
		return Eval(expression.Consequence)
	} else if expression.Alternative != nil {
		return Eval(expression.Alternative)
	} else {
		// Alternative 被选中但它不存在的情况，返回 NULL
		return NULL
	}
}

// NULL 和 FALSE 视为 false，其他视为 true
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}
