// original from https://interpreterbook.com/

package evaluator

import (
	"fmt"
	"interpreter/ast"
	"interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {

	// 对语句求值
	case *ast.Program:
		return evalProgram(n, env)

	case *ast.BlockStatement:
		return evalBlockStatement(n, env)

	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val} // 包裹待返回的 Object

	case *ast.LetStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)

	// 对表达式求值
	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right)

	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isError(function) {
			return function
		}

		// 先对每个实参求值
		args := evalExpressions(n.Arguments, env)

		// for _, arg := range args {
		// 	if isError(arg) {
		// 		return arg
		// 	}
		// }

		// 在上一步骤，如果有其中一个参数求值出错，则返回
		// 单一个元素的 []object.Object，所以不需要逐个参数值检查
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	// 对标识符求值
	case *ast.Identifier:
		val, ok := env.Get(n.Value)
		if !ok {
			return newError("identifier not found: " + n.Value)
		}
		return val

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

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	// evalProgram 跟 evalBlockStatement 很相似，但 BlockStatement 可以嵌套，当遇到
	// return 语句时，需要跳到最外一层 block，所以无法重用 evalBlockStatement
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)

		// if returnValue, ok := result.(*object.ReturnValue); ok {
		// 	return returnValue.Value
		// }

		switch r := result.(type) {
		case *object.ReturnValue: // 拆封 ReturnValue，并跳过剩余的语句
			return r.Value
		case *object.Error: // 返回 Error，并跳过剩余的语句
			return r
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		// 在一组语句中，存在 return 语句
		// if returnValue, ok := result.(*object.ReturnValue); ok {
		// 	return returnValue.Value
		// }

		// 因为 BlockStatement 可以嵌套，所以暂时不拆封 ReturnValue
		// 直到 Program 或者 函数字面量 才解封
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ {
				return result // 跳过剩余的语句
			}

			if result.Type() == object.ERROR_OBJ {
				return result // 跳过剩余的语句
			}
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
		// return NULL
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
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
		// return NULL
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalPlusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		// return NULL
		return newError("unknown operator: +%s", right.Type())
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

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	case operator == "&&":
		return nativeBoolToBooleanObject(
			left.(*object.Boolean).Value &&
				right.(*object.Boolean).Value)

	case operator == "||":
		return nativeBoolToBooleanObject(
			left.(*object.Boolean).Value ||
				right.(*object.Boolean).Value)

	default:
		// return NULL
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		// return NULL
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(expression *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(expression.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(expression.Consequence, env)
	} else if expression.Alternative != nil {
		return Eval(expression.Alternative, env)
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// 用在 "调用 Eval(...) 之后还需进一步执行其他运算" 的场合，用于提早返回
// 比如在调用 evalInfixExpression 之前需要先对 left node 和 right node
// 求值，如果任意一个返回 Error，都应该提前返回 Error，而不是继续执行 evalInfixExpression
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	} else {
		return false
	}
}

func evalExpressions(
	expressions []ast.Expression,
	env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range expressions {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	// 为函数的求值创造一个新的环境，该环境的上层环境为 "函数定义时" 的环境
	// 即静态范围(static scope)
	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	// 用实参填充每一个形参
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

// 拆封函数里 return 语句所包装的值（即 object.Return）给函数调用者
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
