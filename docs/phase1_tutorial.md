# Phase 1学習ガイド：基本言語処理システムの実装

**📚 コンパイラ技術の基礎から実践まで**

## 🎯 Phase 1で学べること

Phase 1では、プログラミング言語処理の基本的な仕組みを段階的に学びます：

1. **字句解析（Lexical Analysis）** - ソースコードをトークンに分割
2. **構文解析（Syntax Analysis）** - トークンから抽象構文木を構築  
3. **AST（Abstract Syntax Tree）** - プログラム構造の内部表現
4. **インタープリター（Interpreter）** - ASTを直接実行する仕組み

これらの技術により、シンプルながら動作するプログラミング言語処理系を完成させます。

---

## 🔤 第1章：字句解析器（Lexer）の理解

### 字句解析とは

字句解析は、ソースコードという文字列を**トークン**と呼ばれる意味のある単位に分割する処理です。

```dog
let x = 42 + 10;
```

上記のコードは以下のトークンに分割されます：

| トークン | 種類 | 説明 |
|---------|------|------|
| `let` | KEYWORD | 変数宣言キーワード |
| `x` | IDENTIFIER | 識別子（変数名） |
| `=` | ASSIGN | 代入演算子 |
| `42` | INTEGER | 整数リテラル |
| `+` | PLUS | 加算演算子 |
| `10` | INTEGER | 整数リテラル |
| `;` | SEMICOLON | セミコロン |

### Lexer実装の詳細解説

#### 基本構造

```go
type Lexer struct {
    input        string // 解析対象の入力文字列
    position     int    // 現在の文字位置
    readPosition int    // 次に読む文字位置
    ch           byte   // 現在検査中の文字
    line         int    // 現在の行番号
    column       int    // 現在の列番号
}
```

**設計のポイント：**
- `position`と`readPosition`の2つのポインタで先読み（look-ahead）を実現
- 行・列番号を保持することでエラー報告を詳細化
- バイト単位で処理（UTF-8対応は今後の拡張課題）

#### 文字読み込みメカニズム

```go
func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0 // ASCII NUL文字（EOF）
    } else {
        l.ch = l.input[l.readPosition]
    }

    // 行・列番号の更新
    if l.ch == '\n' {
        l.line++
        l.column = 0
    } else {
        l.column++
    }

    l.position = l.readPosition
    l.readPosition++
}
```

**実装の工夫：**
- EOF（End of File）を`0`で表現
- 改行文字を検出して行番号を自動更新
- エラー報告に必要な位置情報を正確に追跡

#### トークン識別の実装

```go
func (l *Lexer) NextToken() Token {
    var tok Token
    
    l.skipWhitespace() // 空白文字をスキップ
    
    switch l.ch {
    case '=':
        if l.peekChar() == '=' {
            // == 演算子の処理
            ch := l.ch
            l.readChar()
            tok = Token{Type: EQ, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(ASSIGN, l.ch)
        }
    case '+':
        tok = newToken(PLUS, l.ch)
    case '-':
        tok = newToken(MINUS, l.ch)
    // ... 他の演算子
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()
            tok.Type = LookupIdent(tok.Literal)
            return tok // readIdentifierで位置が進むため、ここでreturn
        } else if isDigit(l.ch) {
            tok.Type = INT
            tok.Literal = l.readNumber()
            return tok
        } else {
            tok = newToken(ILLEGAL, l.ch)
        }
    }
    
    l.readChar()
    return tok
}
```

**アルゴリズムの特徴：**
1. **先読み（Peek）による複数文字演算子の対応**：`==`, `!=`, `<=`, `>=`
2. **識別子とキーワードの区別**：`LookupIdent`で予約語を判定
3. **数値リテラルの読み取り**：連続する数字を一つのトークンとして処理
4. **エラーハンドリング**：不正な文字は`ILLEGAL`トークンとして処理

### 実践演習 1：Lexerの動作確認

以下のコードでLexerの動作を確認しましょう：

```go
func main() {
    input := `let five = 5;
let ten = 10;
let add = fn(x, y) {
    x + y;
};
let result = add(five, ten);`

    l := lexer.New(input)
    
    for {
        tok := l.NextToken()
        fmt.Printf("Type: %s, Literal: %s\n", tok.Type, tok.Literal)
        if tok.Type == token.EOF {
            break
        }
    }
}
```

**期待される出力例：**
```
Type: LET, Literal: let
Type: IDENT, Literal: five
Type: ASSIGN, Literal: =
Type: INT, Literal: 5
Type: SEMICOLON, Literal: ;
...
```

---

## 🌳 第2章：構文解析器（Parser）の実装

### 構文解析とは

構文解析は、字句解析で生成されたトークン列を解析し、プログラムの構造を表現する**抽象構文木（AST）**を構築する処理です。

### 再帰下降構文解析法

pugは**再帰下降構文解析法（Recursive Descent Parsing）**を採用しています。

**特徴：**
- 文法規則に対応する関数を定義
- 各関数が対応する文法要素のASTノードを生成
- 左再帰を避けた文法設計が必要
- 理解しやすく、デバッグが容易

### 文法設計の基本方針

```ebnf
program = statement*

statement = letStatement
          | returnStatement  
          | expressionStatement

letStatement = "let" identifier "=" expression ";"
returnStatement = "return" expression ";"
expressionStatement = expression ";"

expression = infix expression
           | prefix expression
           | identifier
           | integer
           | "(" expression ")"
```

### Pratt Parser（演算子優先順位解析）

複雑な式の解析には**Pratt Parser**を使用します：

```go
type (
    prefixParseFn func() ast.Expression               // 前置式解析関数
    infixParseFn  func(ast.Expression) ast.Expression // 中置式解析関数
)

type Parser struct {
    l *lexer.Lexer
    
    curToken  token.Token
    peekToken token.Token
    
    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns  map[token.TokenType]infixParseFn
    
    errors []string
}
```

#### 演算子優先順位の定義

```go
const (
    _ int = iota
    LOWEST      // 最低優先度
    EQUALS      // ==, !=
    LESSGREATER // > または <
    SUM         // +, -
    PRODUCT     // *, /
    PREFIX      // -X または !X
    CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
    token.EQ:       EQUALS,
    token.NOT_EQ:   EQUALS,
    token.LT:       LESSGREATER,
    token.GT:       LESSGREATER,
    token.PLUS:     SUM,
    token.MINUS:    SUM,
    token.SLASH:    PRODUCT,
    token.ASTERISK: PRODUCT,
    token.LPAREN:   CALL,
}
```

#### 式解析のアルゴリズム

```go
func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefix := p.prefixParseFns[p.curToken.Type]
    if prefix == nil {
        p.noPrefixParseFnError(p.curToken.Type)
        return nil
    }
    
    leftExp := prefix()
    
    for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
        infix := p.infixParseFns[p.peekToken.Type]
        if infix == nil {
            return leftExp
        }
        
        p.nextToken()
        leftExp = infix(leftExp)
    }
    
    return leftExp
}
```

**アルゴリズムの動作例：**

入力：`2 + 3 * 4`

1. `2`を解析（prefix）
2. `+`を検出、右辺`3 * 4`を解析
3. `3`を解析後、`*`の優先度が`+`より高いため先に処理
4. 最終的に`+(2, *(3, 4))`のAST構造を生成

### 実践演習 2：Parser動作確認

```go
func main() {
    input := "let x = 2 + 3 * 4;"
    
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()
    
    // エラーチェック
    if errors := p.Errors(); len(errors) > 0 {
        for _, err := range errors {
            fmt.Println("Parse error:", err)
        }
        return
    }
    
    fmt.Println(program.String())
}
```

---

## 🌲 第3章：抽象構文木（AST）の設計

### AST設計の原則

AST（Abstract Syntax Tree）は、プログラムの構造を階層的に表現するデータ構造です。

**設計原則：**
1. **統一されたインターフェース**：全てのノードが共通の基底を持つ
2. **型安全性**：コンパイル時に不正な操作を検出
3. **拡張性**：新しい言語構文の追加が容易
4. **可読性**：デバッグ時にAST構造を視覚的に確認可能

### ノード階層の実装

```go
// 基底インターフェース
type Node interface {
    TokenLiteral() string
    String() string
}

// 文（Statement）のインターフェース
type Statement interface {
    Node
    statementNode()
}

// 式（Expression）のインターフェース  
type Expression interface {
    Node
    expressionNode()
}

// プログラム全体を表現するルートノード
type Program struct {
    Statements []Statement
}
```

### 具体的なノード実装例

#### Let文のAST表現

```go
type LetStatement struct {
    Token token.Token // LETトークン
    Name  *Identifier // 変数名
    Value Expression  // 代入する値
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
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
```

#### 中置式のAST表現

```go
type InfixExpression struct {
    Token    token.Token // 演算子トークン（+, -, *, / など）
    Left     Expression  // 左辺
    Operator string      // 演算子文字列
    Right    Expression  // 右辺
}
```

**AST例：**`2 + 3 * 4`
```
    InfixExpression(+)
    /               \
   2              InfixExpression(*)
                  /               \
                 3                 4
```

### AST可視化の実装

```go
func (p *Program) String() string {
    var out bytes.Buffer
    
    for _, s := range p.Statements {
        out.WriteString(s.String())
    }
    
    return out.String()
}
```

### 実践演習 3：AST構造の確認

以下のコードでAST構造を可視化：

```go
func printAST(node ast.Node, indent int) {
    spaces := strings.Repeat("  ", indent)
    
    switch n := node.(type) {
    case *ast.Program:
        fmt.Println(spaces + "Program")
        for _, stmt := range n.Statements {
            printAST(stmt, indent+1)
        }
    case *ast.LetStatement:
        fmt.Println(spaces + "LetStatement")
        fmt.Println(spaces + "  Name: " + n.Name.Value)
        if n.Value != nil {
            printAST(n.Value, indent+1)
        }
    case *ast.InfixExpression:
        fmt.Println(spaces + "InfixExpression (" + n.Operator + ")")
        printAST(n.Left, indent+1)
        printAST(n.Right, indent+1)
    // ... 他のノード型
    }
}
```

---

## ⚡ 第4章：インタープリター（Evaluator）の実装

### インタープリターとは

インタープリターは、ASTを直接実行してプログラムの結果を得る処理系です。コンパイルを行わず、ASTを辿りながらリアルタイムで実行します。

### オブジェクトシステムの設計

実行時の値を表現するオブジェクトシステム：

```go
type ObjectType string

const (
    INTEGER_OBJ  = "INTEGER"
    BOOLEAN_OBJ  = "BOOLEAN"
    NULL_OBJ     = "NULL"
    RETURN_OBJ   = "RETURN_VALUE"
    ERROR_OBJ    = "ERROR"
    FUNCTION_OBJ = "FUNCTION"
    STRING_OBJ   = "STRING"
    BUILTIN_OBJ  = "BUILTIN"
    ARRAY_OBJ    = "ARRAY"
    HASH_OBJ     = "HASH"
)

type Object interface {
    Type() ObjectType
    Inspect() string
}

type Integer struct {
    Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
```

### 環境（Environment）とスコープ管理

変数の管理には環境システムを使用：

```go
type Environment struct {
    store map[string]Object
    outer *Environment  // 外側のスコープへの参照
}

func NewEnvironment() *Environment {
    s := make(map[string]Object)
    return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
    value, ok := e.store[name]
    if !ok && e.outer != nil {
        value, ok = e.outer.Get(name)  // 外側のスコープを検索
    }
    return value, ok
}

func (e *Environment) Set(name string, val Object) Object {
    e.store[name] = val
    return val
}
```

### 評価器のコア実装

```go
func Eval(node ast.Node, env *Environment) Object {
    switch node := node.(type) {
    
    case *ast.Program:
        return evalProgram(node.Statements, env)
        
    case *ast.ExpressionStatement:
        return Eval(node.Expression, env)
        
    case *ast.IntegerLiteral:
        return &Integer{Value: node.Value}
        
    case *ast.Boolean:
        return nativeBoolToPugBoolean(node.Value)
        
    case *ast.PrefixExpression:
        right := Eval(node.Right, env)
        if isError(right) {
            return right
        }
        return evalPrefixExpression(node.Operator, right)
        
    case *ast.InfixExpression:
        left := Eval(node.Left, env)
        if isError(left) {
            return left
        }
        right := Eval(node.Right, env)
        if isError(right) {
            return right
        }
        return evalInfixExpression(node.Operator, left, right)
        
    case *ast.IfExpression:
        return evalIfExpression(node, env)
        
    case *ast.Identifier:
        return evalIdentifier(node, env)
        
    case *ast.LetStatement:
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }
        env.Set(node.Name.Value, val)
        
    case *ast.ReturnStatement:
        val := Eval(node.ReturnValue, env)
        if isError(val) {
            return val
        }
        return &ReturnValue{Value: val}
    }
    
    return nil
}
```

### 演算処理の実装

```go
func evalInfixExpression(operator string, left, right Object) Object {
    switch {
    case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
        return evalIntegerInfixExpression(operator, left, right)
    case operator == "==":
        return nativeBoolToPugBoolean(left == right)
    case operator == "!=":
        return nativeBoolToPugBoolean(left != right)
    default:
        return newError("unknown operator: %s %s %s", 
            left.Type(), operator, right.Type())
    }
}

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
        return &Integer{Value: leftVal / rightVal}
    case "<":
        return nativeBoolToPugBoolean(leftVal < rightVal)
    case ">":
        return nativeBoolToPugBoolean(leftVal > rightVal)
    case "==":
        return nativeBoolToPugBoolean(leftVal == rightVal)
    case "!=":
        return nativeBoolToPugBoolean(leftVal != rightVal)
    default:
        return newError("unknown operator: %s", operator)
    }
}
```

### 関数定義と呼び出し

```go
type Function struct {
    Parameters []*ast.Identifier
    Body       *ast.BlockStatement
    Env        *Environment
}

func evalCallExpression(node *ast.CallExpression, env *Environment) Object {
    function := Eval(node.Function, env)
    if isError(function) {
        return function
    }
    
    args := evalExpressions(node.Arguments, env)
    if len(args) == 1 && isError(args[0]) {
        return args[0]
    }
    
    return applyFunction(function, args)
}

func applyFunction(fn Object, args []Object) Object {
    switch fn := fn.(type) {
    case *Function:
        extendedEnv := extendFunctionEnv(fn, args)
        evaluated := Eval(fn.Body, extendedEnv)
        return unwrapReturnValue(evaluated)
    default:
        return newError("not a function: %T", fn)
    }
}
```

### 実践演習 4：インタープリター動作確認

```go
func main() {
    input := `
let fibonacci = fn(x) {
    if (x < 2) {
        return x;
    } else {
        return fibonacci(x - 1) + fibonacci(x - 2);
    }
};

fibonacci(10);
`
    
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()
    
    env := object.NewEnvironment()
    evaluated := evaluator.Eval(program, env)
    
    fmt.Println(evaluated.Inspect())  // Output: 55
}
```

---

## 🔧 第5章：REPL（Read-Eval-Print Loop）の実装

### REPLとは

REPL（Read-Eval-Print Loop）は、対話式プログラミング環境です：

1. **Read**：ユーザー入力を読み取り
2. **Eval**：入力を評価・実行
3. **Print**：結果を表示
4. **Loop**：1-3を繰り返し

### REPL実装

```go
const PROMPT = "pug >> "

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)
    env := object.NewEnvironment()
    
    for {
        fmt.Printf(PROMPT)
        scanned := scanner.Scan()
        if !scanned {
            return
        }
        
        line := scanner.Text()
        l := lexer.New(line)
        p := parser.New(l)
        program := p.ParseProgram()
        
        if len(p.Errors()) != 0 {
            printParserErrors(out, p.Errors())
            continue
        }
        
        evaluated := evaluator.Eval(program, env)
        if evaluated != nil {
            io.WriteString(out, evaluated.Inspect())
            io.WriteString(out, "\n")
        }
    }
}
```

---

## 📊 第6章：エラーハンドリングとデバッグ

### エラー処理の設計

効果的なエラーハンドリング：

```go
type Error struct {
    Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

func newError(format string, a ...interface{}) *Error {
    return &Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj Object) bool {
    if obj != nil {
        return obj.Type() == ERROR_OBJ
    }
    return false
}
```

---

## 🎯 第7章：実践課題と練習問題

### 課題1：基本機能拡張

以下の機能を追加してみましょう：

1. **文字列リテラル対応**
   ```dog
   let message = "Hello, pug!";
   ```

2. **配列リテラル**
   ```dog
   let numbers = [1, 2, 3, 4, 5];
   ```

3. **ハッシュ（連想配列）**
   ```dog
   let person = {"name": "Alice", "age": 30};
   ```

### 課題2：組み込み関数の実装

```go
var builtins = map[string]*Builtin{
    "len": {
        Fn: func(args ...Object) Object {
            if len(args) != 1 {
                return newError("wrong number of arguments. got=%d, want=1", len(args))
            }
            
            switch arg := args[0].(type) {
            case *String:
                return &Integer{Value: int64(len(arg.Value))}
            case *Array:
                return &Integer{Value: int64(len(arg.Elements))}
            default:
                return newError("argument to `len` not supported, got %T", arg)
            }
        },
    },
    "print": {
        Fn: func(args ...Object) Object {
            for _, arg := range args {
                fmt.Println(arg.Inspect())
            }
            return NULL
        },
    },
}
```

### 課題3：制御構造の追加

1. **while文**
   ```dog
   let i = 0;
   while (i < 10) {
       print(i);
       let i = i + 1;
   }
   ```

2. **for文**
   ```dog
   for (let i = 0; i < 10; i = i + 1) {
       print(i);
   }
   ```

---

## 📈 第8章：Phase 2への準備

### Phase 1の成果

Phase 1を完了すると、以下が実装されています：

✅ **字句解析器**：ソースコードのトークン化  
✅ **構文解析器**：ASTの構築  
✅ **インタープリター**：AST直接実行  
✅ **REPL**：対話式実行環境  
✅ **エラーハンドリング**：わかりやすいエラー報告  

### Phase 2で学ぶこと

次のPhase 2では、より高度な技術を学びます：

🎯 **コンパイラ技術**：機械語・アセンブリ生成  
🎯 **型システム**：静的型検査・型推論  
🎯 **最適化**：コード効率化技術  
🎯 **性能分析**：定量的な改善効果測定  

### 移行時のチェックポイント

Phase 2に進む前に以下を確認：

- [ ] 全てのテストが通過する
- [ ] REPLで基本的な計算ができる
- [ ] 関数定義・呼び出しが動作する
- [ ] エラーが適切に報告される
- [ ] コードの可読性が保たれている

---

## 🔗 参考資料とさらなる学習

### 推奨書籍

1. **「Go言語でつくるインタープリター」** - Thorsten Ball
2. **「コンパイラ 作りながら学ぶ」** - 中田育男
3. **「ドラゴンブック」** - Aho, Sethi, Ullman

### 関連技術

- **ANTLR**：パーサージェネレーター
- **LLVM**：コンパイラ基盤技術
- **Tree-sitter**：構文解析ライブラリ

### オンライン学習リソース

- [Stanford CS143 Compilers](https://web.stanford.edu/class/cs143/)
- [MIT 6.035 Computer Language Engineering](https://ocw.mit.edu/courses/electrical-engineering-and-computer-science/6-035-computer-language-engineering-sml-2005/)

---

**🎉 Phase 1完了おめでとうございます！**

基本的な言語処理システムの仕組みを理解し、動作するインタープリターを実装できました。この基礎知識を活かして、Phase 2ではより高度なコンパイラ技術に挑戦しましょう！

**次のステップ：[Phase 2学習ガイド](phase2_tutorial.md)**