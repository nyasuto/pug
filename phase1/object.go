package phase1

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

// ObjectType はオブジェクトの型を表す
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

// Object は全てのオブジェクトが実装するインターフェース
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Hashable はハッシュ可能なオブジェクトのインターフェース
type Hashable interface {
	HashKey() HashKey
}

// HashKey はハッシュのキーを表す
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Integer は整数オブジェクト
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	// Convert int64 to uint64 using Go's well-defined conversion
	// #nosec G115 -- Go's int64 to uint64 conversion is well-defined and safe
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Float は浮動小数点オブジェクト
type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }

// BooleanObj はブール値オブジェクト
type BooleanObj struct {
	Value bool
}

func (b *BooleanObj) Type() ObjectType { return BOOLEAN_OBJ }
func (b *BooleanObj) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *BooleanObj) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// String は文字列オブジェクト
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value)) // Hash.Write never returns an error
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Null はnullオブジェクト
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue はreturn文の値をラップするオブジェクト
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error はエラーオブジェクト
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Function は関数オブジェクト
type Function struct {
	Parameters []*Identifier
	Body       *BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// BuiltinFunction は組み込み関数の型
type BuiltinFunction func(args ...Object) Object

// Builtin は組み込み関数オブジェクト
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// Array は配列オブジェクト
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// HashPair はハッシュのキーと値のペア
type HashPair struct {
	Key   Object
	Value Object
}

// Hash はハッシュオブジェクト
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// よく使用されるオブジェクトのシングルトン
var (
	NULL_OBJ_INSTANCE  = &Null{}
	TRUE_OBJ_INSTANCE  = &BooleanObj{Value: true}
	FALSE_OBJ_INSTANCE = &BooleanObj{Value: false}
)
