package phase2

import (
	"fmt"
	"strings"

	"github.com/nyasuto/pug/phase1"
)

// Type は型を表すインターフェース
type Type interface {
	String() string
	Equals(Type) bool
}

// 基本型の実装

// IntType は整数型
type IntType struct{}

func (t *IntType) String() string {
	return "int"
}

func (t *IntType) Equals(other Type) bool {
	_, ok := other.(*IntType)
	return ok
}

// FloatType は浮動小数点型
type FloatType struct{}

func (t *FloatType) String() string {
	return "float"
}

func (t *FloatType) Equals(other Type) bool {
	_, ok := other.(*FloatType)
	return ok
}

// StringType は文字列型
type StringType struct{}

func (t *StringType) String() string {
	return "string"
}

func (t *StringType) Equals(other Type) bool {
	_, ok := other.(*StringType)
	return ok
}

// BoolType はブール型
type BoolType struct{}

func (t *BoolType) String() string {
	return "bool"
}

func (t *BoolType) Equals(other Type) bool {
	_, ok := other.(*BoolType)
	return ok
}

// FunctionType は関数型
type FunctionType struct {
	Parameters []Type
	ReturnType Type
}

func (t *FunctionType) String() string {
	params := make([]string, len(t.Parameters))
	for i, param := range t.Parameters {
		params[i] = param.String()
	}
	return fmt.Sprintf("fn(%s) -> %s", strings.Join(params, ", "), t.ReturnType.String())
}

func (t *FunctionType) Equals(other Type) bool {
	otherFunc, ok := other.(*FunctionType)
	if !ok {
		return false
	}

	if len(t.Parameters) != len(otherFunc.Parameters) {
		return false
	}

	for i, param := range t.Parameters {
		if !param.Equals(otherFunc.Parameters[i]) {
			return false
		}
	}

	return t.ReturnType.Equals(otherFunc.ReturnType)
}

// ArrayType は配列型
type ArrayType struct {
	ElementType Type
}

func (t *ArrayType) String() string {
	return fmt.Sprintf("[%s]", t.ElementType.String())
}

func (t *ArrayType) Equals(other Type) bool {
	otherArray, ok := other.(*ArrayType)
	if !ok {
		return false
	}
	return t.ElementType.Equals(otherArray.ElementType)
}

// UnknownType は未知の型（型推論で使用）
type UnknownType struct {
	Name string
}

func (t *UnknownType) String() string {
	return fmt.Sprintf("?%s", t.Name)
}

func (t *UnknownType) Equals(other Type) bool {
	otherUnknown, ok := other.(*UnknownType)
	if !ok {
		return false
	}
	return t.Name == otherUnknown.Name
}

// 型の単体インスタンス（シングルトン）
var (
	INT_TYPE    = &IntType{}
	FLOAT_TYPE  = &FloatType{}
	STRING_TYPE = &StringType{}
	BOOL_TYPE   = &BoolType{}
)

// TypeEnvironment は型環境（変数名と型のマッピング）
type TypeEnvironment struct {
	store map[string]Type
	outer *TypeEnvironment
}

// NewTypeEnvironment は新しい型環境を作成する
func NewTypeEnvironment() *TypeEnvironment {
	return &TypeEnvironment{
		store: make(map[string]Type),
		outer: nil,
	}
}

// NewEnclosedTypeEnvironment は外側の環境を持つ新しい型環境を作成する
func NewEnclosedTypeEnvironment(outer *TypeEnvironment) *TypeEnvironment {
	env := NewTypeEnvironment()
	env.outer = outer
	return env
}

// Get は変数の型を取得する
func (e *TypeEnvironment) Get(name string) (Type, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(name)
	}
	return value, ok
}

// Set は変数に型を設定する
func (e *TypeEnvironment) Set(name string, val Type) Type {
	e.store[name] = val
	return val
}

// TypeChecker は型検査器
type TypeChecker struct {
	env    *TypeEnvironment
	errors []string
}

// NewTypeChecker は新しい型検査器を作成する
func NewTypeChecker() *TypeChecker {
	env := NewTypeEnvironment()

	// 組み込み関数の型を設定
	env.Set("len", &FunctionType{
		Parameters: []Type{&ArrayType{ElementType: &UnknownType{Name: "T"}}},
		ReturnType: INT_TYPE,
	})
	env.Set("first", &FunctionType{
		Parameters: []Type{&ArrayType{ElementType: &UnknownType{Name: "T"}}},
		ReturnType: &UnknownType{Name: "T"},
	})
	env.Set("rest", &FunctionType{
		Parameters: []Type{&ArrayType{ElementType: &UnknownType{Name: "T"}}},
		ReturnType: &ArrayType{ElementType: &UnknownType{Name: "T"}},
	})
	env.Set("push", &FunctionType{
		Parameters: []Type{
			&ArrayType{ElementType: &UnknownType{Name: "T"}},
			&UnknownType{Name: "T"},
		},
		ReturnType: &ArrayType{ElementType: &UnknownType{Name: "T"}},
	})
	env.Set("puts", &FunctionType{
		Parameters: []Type{}, // 可変長引数として扱う
		ReturnType: STRING_TYPE,
	})

	return &TypeChecker{
		env:    env,
		errors: []string{},
	}
}

// CheckProgram はプログラム全体の型検査を行う
func (tc *TypeChecker) CheckProgram(program *phase1.Program) (Type, []string) {
	var lastType Type = &UnknownType{Name: "void"}

	for _, stmt := range program.Statements {
		lastType = tc.CheckStatement(stmt)
	}

	return lastType, tc.errors
}

// CheckStatement は文の型検査を行う
func (tc *TypeChecker) CheckStatement(stmt phase1.Statement) Type {
	switch node := stmt.(type) {
	case *phase1.LetStatement:
		return tc.checkLetStatement(node)
	case *phase1.ReturnStatement:
		return tc.checkReturnStatement(node)
	case *phase1.ExpressionStatement:
		return tc.CheckExpression(node.Expression)
	default:
		tc.addError(fmt.Sprintf("unknown statement type: %T", stmt))
		return &UnknownType{Name: "error"}
	}
}

// CheckExpression は式の型検査を行う
func (tc *TypeChecker) CheckExpression(expr phase1.Expression) Type {
	switch node := expr.(type) {
	case *phase1.IntegerLiteral:
		return INT_TYPE
	case *phase1.FloatLiteral:
		return FLOAT_TYPE
	case *phase1.StringLiteral:
		return STRING_TYPE
	case *phase1.Boolean:
		return BOOL_TYPE
	case *phase1.Identifier:
		return tc.checkIdentifier(node)
	case *phase1.InfixExpression:
		return tc.checkInfixExpression(node)
	case *phase1.PrefixExpression:
		return tc.checkPrefixExpression(node)
	case *phase1.IfExpression:
		return tc.checkIfExpression(node)
	case *phase1.FunctionLiteral:
		return tc.checkFunctionLiteral(node)
	case *phase1.CallExpression:
		return tc.checkCallExpression(node)
	default:
		tc.addError(fmt.Sprintf("unknown expression type: %T", expr))
		return &UnknownType{Name: "error"}
	}
}

// checkLetStatement はlet文の型検査を行う
func (tc *TypeChecker) checkLetStatement(stmt *phase1.LetStatement) Type {
	valueType := tc.CheckExpression(stmt.Value)
	tc.env.Set(stmt.Name.Value, valueType)
	return valueType
}

// checkReturnStatement はreturn文の型検査を行う
func (tc *TypeChecker) checkReturnStatement(stmt *phase1.ReturnStatement) Type {
	if stmt.ReturnValue != nil {
		return tc.CheckExpression(stmt.ReturnValue)
	}
	return &UnknownType{Name: "void"}
}

// checkIdentifier は識別子の型検査を行う
func (tc *TypeChecker) checkIdentifier(node *phase1.Identifier) Type {
	typ, ok := tc.env.Get(node.Value)
	if !ok {
		tc.addError(fmt.Sprintf("identifier not found: %s", node.Value))
		return &UnknownType{Name: "error"}
	}
	return typ
}

// checkInfixExpression は中置式の型検査を行う
func (tc *TypeChecker) checkInfixExpression(node *phase1.InfixExpression) Type {
	leftType := tc.CheckExpression(node.Left)
	rightType := tc.CheckExpression(node.Right)

	switch node.Operator {
	case "+", "-", "*", "/", "%":
		// 算術演算子
		leftIsNumeric := tc.isNumericType(leftType) || tc.isUnknownType(leftType)
		rightIsNumeric := tc.isNumericType(rightType) || tc.isUnknownType(rightType)

		if !leftIsNumeric {
			tc.addError(fmt.Sprintf("left operand of %s must be numeric, got %s", node.Operator, leftType.String()))
		}
		if !rightIsNumeric {
			tc.addError(fmt.Sprintf("right operand of %s must be numeric, got %s", node.Operator, rightType.String()))
		}

		// 型昇格: int + float = float
		if leftType.Equals(FLOAT_TYPE) || rightType.Equals(FLOAT_TYPE) {
			return FLOAT_TYPE
		}
		return INT_TYPE

	case "==", "!=":
		// 等価演算子：同じ型同士で比較
		if !leftType.Equals(rightType) {
			tc.addError(fmt.Sprintf("cannot compare %s with %s", leftType.String(), rightType.String()))
		}
		return BOOL_TYPE

	case "<", ">", "<=", ">=":
		// 比較演算子：数値型のみ
		if !tc.isNumericType(leftType) {
			tc.addError(fmt.Sprintf("left operand of %s must be numeric, got %s", node.Operator, leftType.String()))
		}
		if !tc.isNumericType(rightType) {
			tc.addError(fmt.Sprintf("right operand of %s must be numeric, got %s", node.Operator, rightType.String()))
		}
		return BOOL_TYPE

	default:
		tc.addError(fmt.Sprintf("unknown infix operator: %s", node.Operator))
		return &UnknownType{Name: "error"}
	}
}

// checkPrefixExpression は前置式の型検査を行う
func (tc *TypeChecker) checkPrefixExpression(node *phase1.PrefixExpression) Type {
	rightType := tc.CheckExpression(node.Right)

	switch node.Operator {
	case "-":
		if !tc.isNumericType(rightType) {
			tc.addError(fmt.Sprintf("operand of unary - must be numeric, got %s", rightType.String()))
		}
		return rightType
	case "!":
		if !rightType.Equals(BOOL_TYPE) {
			tc.addError(fmt.Sprintf("operand of ! must be bool, got %s", rightType.String()))
		}
		return BOOL_TYPE
	default:
		tc.addError(fmt.Sprintf("unknown prefix operator: %s", node.Operator))
		return &UnknownType{Name: "error"}
	}
}

// checkIfExpression はif式の型検査を行う
func (tc *TypeChecker) checkIfExpression(node *phase1.IfExpression) Type {
	conditionType := tc.CheckExpression(node.Condition)
	if !conditionType.Equals(BOOL_TYPE) {
		tc.addError(fmt.Sprintf("if condition must be bool, got %s", conditionType.String()))
	}

	// if文のブロック内で新しいスコープを作成
	consequenceEnv := NewEnclosedTypeEnvironment(tc.env)
	oldEnv := tc.env
	tc.env = consequenceEnv

	var consequenceType Type
	for _, stmt := range node.Consequence.Statements {
		consequenceType = tc.CheckStatement(stmt)
	}

	tc.env = oldEnv

	if node.Alternative != nil {
		// else文のブロック内で新しいスコープを作成
		alternativeEnv := NewEnclosedTypeEnvironment(tc.env)
		tc.env = alternativeEnv

		var alternativeType Type
		for _, stmt := range node.Alternative.Statements {
			alternativeType = tc.CheckStatement(stmt)
		}

		tc.env = oldEnv

		// if-else式の型は両方の分岐の型が一致している必要がある
		if consequenceType != nil && alternativeType != nil && !consequenceType.Equals(alternativeType) {
			tc.addError(fmt.Sprintf("if-else branches have different types: %s vs %s",
				consequenceType.String(), alternativeType.String()))
		}

		if alternativeType != nil {
			return alternativeType
		}
	}

	if consequenceType != nil {
		return consequenceType
	}

	return &UnknownType{Name: "void"}
}

// checkFunctionLiteral は関数リテラルの型検査を行う
func (tc *TypeChecker) checkFunctionLiteral(node *phase1.FunctionLiteral) Type {
	// パラメータの型を推論（現在は全てUnknownType）
	params := make([]Type, len(node.Parameters))
	for i := range node.Parameters {
		params[i] = &UnknownType{Name: fmt.Sprintf("param%d", i)}
	}

	// 関数本体用の新しい環境を作成
	funcEnv := NewEnclosedTypeEnvironment(tc.env)
	oldEnv := tc.env
	tc.env = funcEnv

	// パラメータを環境に追加
	for i, param := range node.Parameters {
		tc.env.Set(param.Value, params[i])
	}

	// 関数本体の型検査
	var returnType Type = &UnknownType{Name: "void"}
	for _, stmt := range node.Body.Statements {
		stmtType := tc.CheckStatement(stmt)
		if _, ok := stmt.(*phase1.ReturnStatement); ok {
			returnType = stmtType
		}
	}

	tc.env = oldEnv

	return &FunctionType{
		Parameters: params,
		ReturnType: returnType,
	}
}

// checkCallExpression は関数呼び出し式の型検査を行う
func (tc *TypeChecker) checkCallExpression(node *phase1.CallExpression) Type {
	funcType := tc.CheckExpression(node.Function)

	// 組み込み関数の特別処理
	if ident, ok := node.Function.(*phase1.Identifier); ok {
		switch ident.Value {
		case "puts":
			// putsは任意の型の引数を受け取る
			for _, arg := range node.Arguments {
				tc.CheckExpression(arg)
			}
			return STRING_TYPE
		}
	}

	function, ok := funcType.(*FunctionType)
	if !ok {
		tc.addError(fmt.Sprintf("not a function: %s", funcType.String()))
		return &UnknownType{Name: "error"}
	}

	// 引数の数をチェック
	if len(node.Arguments) != len(function.Parameters) {
		tc.addError(fmt.Sprintf("wrong number of arguments: expected %d, got %d",
			len(function.Parameters), len(node.Arguments)))
		return function.ReturnType
	}

	// 各引数の型をチェック
	for i, arg := range node.Arguments {
		argType := tc.CheckExpression(arg)
		expectedType := function.Parameters[i]

		// UnknownTypeは任意の型にマッチする
		if _, ok := expectedType.(*UnknownType); ok {
			continue
		}

		if !argType.Equals(expectedType) {
			tc.addError(fmt.Sprintf("argument %d: expected %s, got %s",
				i, expectedType.String(), argType.String()))
		}
	}

	return function.ReturnType
}

// isNumericType は数値型かどうかをチェックする
func (tc *TypeChecker) isNumericType(t Type) bool {
	return t.Equals(INT_TYPE) || t.Equals(FLOAT_TYPE)
}

// isUnknownType は未知型かどうかをチェックする
func (tc *TypeChecker) isUnknownType(t Type) bool {
	_, ok := t.(*UnknownType)
	return ok
}

// addError はエラーメッセージを追加する
func (tc *TypeChecker) addError(message string) {
	tc.errors = append(tc.errors, message)
}

// GetErrors はエラーメッセージのリストを取得する
func (tc *TypeChecker) GetErrors() []string {
	return tc.errors
}

// HasErrors はエラーがあるかどうかをチェックする
func (tc *TypeChecker) HasErrors() bool {
	return len(tc.errors) > 0
}
