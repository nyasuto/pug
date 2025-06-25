package phase2

import (
	"github.com/nyasuto/pug/phase1"
)

// Symbol は変数の情報を表す
type Symbol struct {
	Name  string
	Type  Type
	Index int    // スタックオフセットまたはインデックス
	Scope string // "local", "global", "function"
}

// SymbolTable はスコープ管理を行うシンボルテーブル
type SymbolTable struct {
	parent     *SymbolTable
	store      map[string]*Symbol
	numSymbols int
	scopeLevel int
}

// NewSymbolTable は新しいシンボルテーブルを作成する
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:      make(map[string]*Symbol),
		numSymbols: 0,
		scopeLevel: 0,
	}
}

// NewEnclosedSymbolTable は親スコープを持つシンボルテーブルを作成する
func NewEnclosedSymbolTable(parent *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.parent = parent
	s.scopeLevel = parent.scopeLevel + 1
	return s
}

// Define は新しいシンボルを定義する
func (s *SymbolTable) Define(name string, symbolType Type) *Symbol {
	symbol := &Symbol{
		Name:  name,
		Type:  symbolType,
		Index: s.numSymbols,
		Scope: s.getScopeString(),
	}
	s.store[name] = symbol
	s.numSymbols++
	return symbol
}

// Resolve はシンボルを解決する（現在のスコープから親スコープへと検索）
func (s *SymbolTable) Resolve(name string) (*Symbol, bool) {
	symbol, ok := s.store[name]
	if ok {
		return symbol, true
	}

	// 親スコープで検索
	if s.parent != nil {
		return s.parent.Resolve(name)
	}

	return nil, false
}

// getScopeString はスコープレベルに応じた文字列を返す
func (s *SymbolTable) getScopeString() string {
	switch s.scopeLevel {
	case 0:
		return "global"
	case 1:
		return "function"
	default:
		return "local"
	}
}

// GetScopeLevel はスコープレベルを返す
func (s *SymbolTable) GetScopeLevel() int {
	return s.scopeLevel
}

// GetNumSymbols はシンボル数を返す
func (s *SymbolTable) GetNumSymbols() int {
	return s.numSymbols
}

// GetParent は親シンボルテーブルを返す
func (s *SymbolTable) GetParent() *SymbolTable {
	return s.parent
}

// LoopContext はループコンテキストを管理する構造体
type LoopContext struct {
	BreakLabel    string
	ContinueLabel string
	Parent        *LoopContext
}

// NewLoopContext は新しいループコンテキストを作成する
func NewLoopContext(breakLabel, continueLabel string, parent *LoopContext) *LoopContext {
	return &LoopContext{
		BreakLabel:    breakLabel,
		ContinueLabel: continueLabel,
		Parent:        parent,
	}
}

// ControlFlowAnalyzer は制御フロー解析を行う
type ControlFlowAnalyzer struct {
	symbolTable *SymbolTable
	loopContext *LoopContext
	errors      []string
}

// NewControlFlowAnalyzer は新しい制御フロー解析器を作成する
func NewControlFlowAnalyzer() *ControlFlowAnalyzer {
	return &ControlFlowAnalyzer{
		symbolTable: NewSymbolTable(),
		errors:      []string{},
	}
}

// EnterScope は新しいスコープに入る
func (cfa *ControlFlowAnalyzer) EnterScope() {
	cfa.symbolTable = NewEnclosedSymbolTable(cfa.symbolTable)
}

// ExitScope は現在のスコープから出る
func (cfa *ControlFlowAnalyzer) ExitScope() {
	if cfa.symbolTable.parent != nil {
		cfa.symbolTable = cfa.symbolTable.parent
	}
}

// EnterLoop は新しいループコンテキストに入る
func (cfa *ControlFlowAnalyzer) EnterLoop(breakLabel, continueLabel string) {
	cfa.loopContext = NewLoopContext(breakLabel, continueLabel, cfa.loopContext)
}

// ExitLoop は現在のループコンテキストから出る
func (cfa *ControlFlowAnalyzer) ExitLoop() {
	if cfa.loopContext != nil {
		cfa.loopContext = cfa.loopContext.Parent
	}
}

// GetCurrentLoopContext は現在のループコンテキストを返す
func (cfa *ControlFlowAnalyzer) GetCurrentLoopContext() *LoopContext {
	return cfa.loopContext
}

// GetSymbolTable は現在のシンボルテーブルを返す
func (cfa *ControlFlowAnalyzer) GetSymbolTable() *SymbolTable {
	return cfa.symbolTable
}

// GetErrors はエラーリストを返す
func (cfa *ControlFlowAnalyzer) GetErrors() []string {
	return cfa.errors
}

// AddError はエラーを追加する
func (cfa *ControlFlowAnalyzer) AddError(message string) {
	cfa.errors = append(cfa.errors, message)
}

// ValidateControlFlow は制御フローの妥当性を検証する
func (cfa *ControlFlowAnalyzer) ValidateControlFlow(stmt phase1.Statement) {
	switch node := stmt.(type) {
	case *phase1.BreakStatement:
		if cfa.loopContext == nil {
			cfa.AddError("break statement outside of loop")
		}
	case *phase1.ContinueStatement:
		if cfa.loopContext == nil {
			cfa.AddError("continue statement outside of loop")
		}
	case *phase1.BlockStatement:
		cfa.EnterScope()
		for _, s := range node.Statements {
			cfa.ValidateControlFlow(s)
		}
		cfa.ExitScope()
	}
}
