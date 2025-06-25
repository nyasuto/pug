package phase1

import (
	"testing"
)

// TestBuiltinLenFunction はlen関数のカバレッジを向上させるテスト
func TestBuiltinLenFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected interface{}
	}{
		{
			name:     "len of array",
			input:    &Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}}},
			expected: int64(2),
		},
		{
			name:     "len of string",
			input:    &String{Value: "hello"},
			expected: int64(5),
		},
		{
			name:     "len of empty array",
			input:    &Array{Elements: []Object{}},
			expected: int64(0),
		},
		{
			name:     "len of empty string",
			input:    &String{Value: ""},
			expected: int64(0),
		},
	}

	lenFn := builtins["len"]
	if lenFn == nil {
		t.Fatal("len builtin not found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lenFn.Fn(tt.input)

			integer, ok := result.(*Integer)
			if !ok {
				t.Errorf("result is not Integer. got=%T (%+v)", result, result)
				return
			}

			if integer.Value != tt.expected {
				t.Errorf("result has wrong value. got=%d, want=%d", integer.Value, tt.expected)
			}
		})
	}
}

// TestBuiltinLenErrors はlen関数のエラーケースをテストする
func TestBuiltinLenErrors(t *testing.T) {
	tests := []struct {
		name  string
		args  []Object
		error string
	}{
		{
			name:  "wrong number of arguments - none",
			args:  []Object{},
			error: "wrong number of arguments. got=0, want=1",
		},
		{
			name:  "wrong number of arguments - too many",
			args:  []Object{&String{Value: "test"}, &String{Value: "extra"}},
			error: "wrong number of arguments. got=2, want=1",
		},
		{
			name:  "unsupported argument type",
			args:  []Object{&Integer{Value: 42}},
			error: "argument to `len` not supported, got INTEGER",
		},
	}

	lenFn := builtins["len"]
	if lenFn == nil {
		t.Fatal("len builtin not found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lenFn.Fn(tt.args...)

			errObj, ok := result.(*Error)
			if !ok {
				t.Errorf("result is not Error. got=%T (%+v)", result, result)
				return
			}

			if errObj.Message != tt.error {
				t.Errorf("wrong error message. got=%q, want=%q", errObj.Message, tt.error)
			}
		})
	}
}

// TestBuiltinFirstFunction はfirst関数をテストする
func TestBuiltinFirstFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected Object
	}{
		{
			name:     "first of non-empty array",
			input:    &Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}}},
			expected: &Integer{Value: 1},
		},
		{
			name:     "first of empty array",
			input:    &Array{Elements: []Object{}},
			expected: NULL_OBJ_INSTANCE,
		},
	}

	firstFn := builtins["first"]
	if firstFn == nil {
		t.Fatal("first builtin not found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := firstFn.Fn(tt.input)

			if !compareObjects(result, tt.expected) {
				t.Errorf("result mismatch. got=%s, want=%s", result.Inspect(), tt.expected.Inspect())
			}
		})
	}
}

// TestBuiltinFirstErrors はfirst関数のエラーケースをテストする
func TestBuiltinFirstErrors(t *testing.T) {
	tests := []struct {
		name  string
		args  []Object
		error string
	}{
		{
			name:  "wrong number of arguments",
			args:  []Object{},
			error: "wrong number of arguments. got=0, want=1",
		},
		{
			name:  "wrong argument type",
			args:  []Object{&Integer{Value: 42}},
			error: "argument to `first` must be ARRAY, got *phase1.Integer",
		},
	}

	firstFn := builtins["first"]
	if firstFn == nil {
		t.Fatal("first builtin not found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := firstFn.Fn(tt.args...)

			errObj, ok := result.(*Error)
			if !ok {
				t.Errorf("result is not Error. got=%T (%+v)", result, result)
				return
			}

			if errObj.Message != tt.error {
				t.Errorf("wrong error message. got=%q, want=%q", errObj.Message, tt.error)
			}
		})
	}
}

// TestBuiltinLastFunction はlast関数をテストする
func TestBuiltinLastFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected Object
	}{
		{
			name:     "last of non-empty array",
			input:    &Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}}},
			expected: &Integer{Value: 2},
		},
		{
			name:     "last of empty array",
			input:    &Array{Elements: []Object{}},
			expected: NULL_OBJ_INSTANCE,
		},
	}

	lastFn := builtins["last"]
	if lastFn == nil {
		t.Fatal("last builtin not found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lastFn.Fn(tt.input)

			if !compareObjects(result, tt.expected) {
				t.Errorf("result mismatch. got=%s, want=%s", result.Inspect(), tt.expected.Inspect())
			}
		})
	}
}

// TestBuiltinRestFunction はrest関数をテストする
func TestBuiltinRestFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected Object
	}{
		{
			name:     "rest of non-empty array",
			input:    &Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}}},
			expected: &Array{Elements: []Object{&Integer{Value: 2}, &Integer{Value: 3}}},
		},
		{
			name:     "rest of single element array",
			input:    &Array{Elements: []Object{&Integer{Value: 1}}},
			expected: &Array{Elements: []Object{}},
		},
		{
			name:     "rest of empty array",
			input:    &Array{Elements: []Object{}},
			expected: NULL_OBJ_INSTANCE,
		},
	}

	restFn := builtins["rest"]
	if restFn == nil {
		t.Fatal("rest builtin not found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := restFn.Fn(tt.input)

			if !compareObjects(result, tt.expected) {
				t.Errorf("result mismatch. got=%s, want=%s", result.Inspect(), tt.expected.Inspect())
			}
		})
	}
}

// TestBuiltinPushFunction はpush関数をテストする
func TestBuiltinPushFunction(t *testing.T) {
	tests := []struct {
		name     string
		array    Object
		element  Object
		expected Object
	}{
		{
			name:     "push to non-empty array",
			array:    &Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}}},
			element:  &Integer{Value: 3},
			expected: &Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}}},
		},
		{
			name:     "push to empty array",
			array:    &Array{Elements: []Object{}},
			element:  &Integer{Value: 1},
			expected: &Array{Elements: []Object{&Integer{Value: 1}}},
		},
	}

	pushFn := builtins["push"]
	if pushFn == nil {
		t.Fatal("push builtin not found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pushFn.Fn(tt.array, tt.element)

			if !compareObjects(result, tt.expected) {
				t.Errorf("result mismatch. got=%s, want=%s", result.Inspect(), tt.expected.Inspect())
			}
		})
	}
}

// TestBuiltinTypeFunction はtype関数をテストする
func TestBuiltinTypeFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected string
	}{
		{
			name:     "type of integer",
			input:    &Integer{Value: 42},
			expected: "INTEGER",
		},
		{
			name:     "type of string",
			input:    &String{Value: "hello"},
			expected: "STRING",
		},
		{
			name:     "type of boolean",
			input:    &BooleanObj{Value: true},
			expected: "BOOLEAN",
		},
		{
			name:     "type of array",
			input:    &Array{Elements: []Object{}},
			expected: "ARRAY",
		},
		{
			name:     "type of null",
			input:    NULL_OBJ_INSTANCE,
			expected: "NULL",
		},
	}

	typeFn := builtins["type"]
	if typeFn == nil {
		t.Fatal("type builtin not found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := typeFn.Fn(tt.input)

			strObj, ok := result.(*String)
			if !ok {
				t.Errorf("result is not String. got=%T (%+v)", result, result)
				return
			}

			if strObj.Value != tt.expected {
				t.Errorf("wrong type string. got=%q, want=%q", strObj.Value, tt.expected)
			}
		})
	}
}

// compareObjects はオブジェクトを比較するヘルパー関数
func compareObjects(a, b Object) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch obj := a.(type) {
	case *Integer:
		return obj.Value == b.(*Integer).Value
	case *String:
		return obj.Value == b.(*String).Value
	case *BooleanObj:
		return obj.Value == b.(*BooleanObj).Value
	case *Array:
		other := b.(*Array)
		if len(obj.Elements) != len(other.Elements) {
			return false
		}
		for i, elem := range obj.Elements {
			if !compareObjects(elem, other.Elements[i]) {
				return false
			}
		}
		return true
	case *Null:
		return true
	default:
		return false
	}
}
