package phase1

// Environment は変数の束縛を管理する環境
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment は新しい環境を作成する
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment は外側の環境を持つ新しい環境を作成する
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get は変数の値を取得する
func (e *Environment) Get(name string) (Object, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(name)
	}
	return value, ok
}

// Set は変数に値を設定する
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
