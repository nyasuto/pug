package phase1

import "fmt"

// Eval はASTノードを評価してオブジェクトを返す
func Eval(node Node, env *Environment) Object {
	switch node := node.(type) {

	// Program
	case *Program:
		return evalProgram(node.Statements, env)

	// Statements
	case *ExpressionStatement:
		return Eval(node.Expression, env)

	case *LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val

	case *ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}

	case *BlockStatement:
		return evalBlockStatement(node, env)

	// Expressions
	case *IntegerLiteral:
		return &Integer{Value: node.Value}

	case *FloatLiteral:
		return &Float{Value: node.Value}

	case *StringLiteral:
		return &String{Value: node.Value}

	case *Boolean:
		return nativeBoolToPugBoolean(node.Value)

	case *PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *IfExpression:
		return evalIfExpression(node, env)

	case *Identifier:
		return evalIdentifier(node, env)

	case *FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &Function{Parameters: params, Env: env, Body: body}

	case *CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	default:
		return newError("unknown node type: %T", node)
	}
}

// evalProgram はプログラム全体を評価する
func evalProgram(stmts []Statement, env *Environment) Object {
	var result Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}

	return result
}

// evalBlockStatement はブロック文を評価する
func evalBlockStatement(block *BlockStatement, env *Environment) Object {
	var result Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// evalPrefixExpression は前置演算子式を評価する
func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "+":
		return evalPlusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// evalBangOperatorExpression は論理否定演算子を評価する
func evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE_OBJ_INSTANCE:
		return FALSE_OBJ_INSTANCE
	case FALSE_OBJ_INSTANCE:
		return TRUE_OBJ_INSTANCE
	case NULL_OBJ_INSTANCE:
		return TRUE_OBJ_INSTANCE
	default:
		return FALSE_OBJ_INSTANCE
	}
}

// evalMinusPrefixOperatorExpression はマイナス前置演算子を評価する
func evalMinusPrefixOperatorExpression(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: -right.Value}
	case *Float:
		return &Float{Value: -right.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

// evalPlusPrefixOperatorExpression はプラス前置演算子を評価する
func evalPlusPrefixOperatorExpression(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: +right.Value}
	case *Float:
		return &Float{Value: +right.Value}
	default:
		return newError("unknown operator: +%s", right.Type())
	}
}

// evalInfixExpression は中置演算子式を評価する
func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ:
		leftFloat := &Float{Value: float64(left.(*Integer).Value)}
		return evalFloatInfixExpression(operator, leftFloat, right)
	case left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ:
		rightFloat := &Float{Value: float64(right.(*Integer).Value)}
		return evalFloatInfixExpression(operator, left, rightFloat)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToPugBoolean(left == right)
	case operator == "!=":
		return nativeBoolToPugBoolean(left != right)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// evalIntegerInfixExpression は整数同士の中置演算子を評価する
func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError("modulo by zero")
		}
		return &Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToPugBoolean(leftVal < rightVal)
	case ">":
		return nativeBoolToPugBoolean(leftVal > rightVal)
	case "<=":
		return nativeBoolToPugBoolean(leftVal <= rightVal)
	case ">=":
		return nativeBoolToPugBoolean(leftVal >= rightVal)
	case "==":
		return nativeBoolToPugBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToPugBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

// evalFloatInfixExpression は浮動小数点同士の中置演算子を評価する
func evalFloatInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Float).Value
	rightVal := right.(*Float).Value

	switch operator {
	case "+":
		return &Float{Value: leftVal + rightVal}
	case "-":
		return &Float{Value: leftVal - rightVal}
	case "*":
		return &Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0.0 {
			return newError("division by zero")
		}
		return &Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToPugBoolean(leftVal < rightVal)
	case ">":
		return nativeBoolToPugBoolean(leftVal > rightVal)
	case "<=":
		return nativeBoolToPugBoolean(leftVal <= rightVal)
	case ">=":
		return nativeBoolToPugBoolean(leftVal >= rightVal)
	case "==":
		return nativeBoolToPugBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToPugBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

// evalStringInfixExpression は文字列同士の中置演算子を評価する
func evalStringInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value

	switch operator {
	case "+":
		return &String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToPugBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToPugBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// evalIfExpression はif式を評価する
func evalIfExpression(ie *IfExpression, env *Environment) Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL_OBJ_INSTANCE
	}
}

// evalIdentifier は識別子を評価する
func evalIdentifier(node *Identifier, env *Environment) Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

// evalExpressions は式のリストを評価する
func evalExpressions(exps []Expression, env *Environment) []Object {
	var result []Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// applyFunction は関数を適用する
func applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		// 引数の数をチェック
		if len(args) != len(fn.Parameters) {
			return newError("wrong number of arguments: want=%d, got=%d", len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %T", fn)
	}
}

// extendFunctionEnv は関数の環境を拡張する
func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

// unwrapReturnValue はReturnValueをアンラップする
func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// isTruthy はオブジェクトが真偽値として真かどうかを判定する
func isTruthy(obj Object) bool {
	switch obj {
	case NULL_OBJ_INSTANCE:
		return false
	case TRUE_OBJ_INSTANCE:
		return true
	case FALSE_OBJ_INSTANCE:
		return false
	default:
		return true
	}
}

// isError はオブジェクトがエラーかどうかを判定する
func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

// newError は新しいエラーオブジェクトを作成する
func newError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// nativeBoolToPugBoolean はGoのboolをPugのBooleanオブジェクトに変換する
func nativeBoolToPugBoolean(input bool) *BooleanObj {
	if input {
		return TRUE_OBJ_INSTANCE
	}
	return FALSE_OBJ_INSTANCE
}
