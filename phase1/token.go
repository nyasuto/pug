package phase1

// TokenType はトークンの種類を表す
type TokenType string

// Token はレクサーが生成するトークンを表す
type Token struct {
	Type     TokenType
	Literal  string
	Line     int
	Column   int
	Position int
}

// トークンタイプの定数定義
const (
	// 特殊トークン
	ILLEGAL TokenType = "ILLEGAL" // 不正な文字
	EOF     TokenType = "EOF"     // ファイル終端

	// 識別子とリテラル
	IDENT  TokenType = "IDENT"  // 識別子（変数名、関数名）
	INT    TokenType = "INT"    // 整数リテラル
	FLOAT  TokenType = "FLOAT"  // 浮動小数点リテラル
	STRING TokenType = "STRING" // 文字列リテラル

	// 演算子
	ASSIGN   TokenType = "=" // 代入
	PLUS     TokenType = "+" // 加算
	MINUS    TokenType = "-" // 減算
	MULTIPLY TokenType = "*" // 乗算
	DIVIDE   TokenType = "/" // 除算
	MODULO   TokenType = "%" // 剰余

	// 比較演算子
	EQ     TokenType = "==" // 等価
	NOT_EQ TokenType = "!=" // 非等価
	LT     TokenType = "<"  // 小なり
	GT     TokenType = ">"  // 大なり
	LTE    TokenType = "<=" // 以下
	GTE    TokenType = ">=" // 以上

	// 論理演算子
	AND TokenType = "&&" // 論理積
	OR  TokenType = "||" // 論理和
	NOT TokenType = "!"  // 論理否定

	// 区切り文字
	COMMA     TokenType = ","  // カンマ
	SEMICOLON TokenType = ";"  // セミコロン
	COLON     TokenType = ":"  // コロン
	ARROW     TokenType = "->" // 矢印（関数戻り値型）

	// 括弧類
	LPAREN   TokenType = "(" // 左丸括弧
	RPAREN   TokenType = ")" // 右丸括弧
	LBRACE   TokenType = "{" // 左波括弧
	RBRACE   TokenType = "}" // 右波括弧
	LBRACKET TokenType = "[" // 左角括弧
	RBRACKET TokenType = "]" // 右角括弧

	// キーワード
	LET      TokenType = "LET"      // let文
	FN       TokenType = "FN"       // function
	IF       TokenType = "IF"       // if文
	ELSE     TokenType = "ELSE"     // else文
	RETURN   TokenType = "RETURN"   // return文
	TRUE     TokenType = "TRUE"     // boolean true
	FALSE    TokenType = "FALSE"    // boolean false
	WHILE    TokenType = "WHILE"    // while文
	FOR      TokenType = "FOR"      // for文
	BREAK    TokenType = "BREAK"    // break文
	CONTINUE TokenType = "CONTINUE" // continue文

	// 型キーワード
	INT_TYPE    TokenType = "INT_TYPE"    // int型
	FLOAT_TYPE  TokenType = "FLOAT_TYPE"  // float型
	STRING_TYPE TokenType = "STRING_TYPE" // string型
	BOOL_TYPE   TokenType = "BOOL_TYPE"   // bool型
)

// keywords は予約語のマップ
var keywords = map[string]TokenType{
	"let":      LET,
	"fn":       FN,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"true":     TRUE,
	"false":    FALSE,
	"while":    WHILE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
	"int":      INT_TYPE,
	"float":    FLOAT_TYPE,
	"string":   STRING_TYPE,
	"bool":     BOOL_TYPE,
}

// LookupIdent は識別子がキーワードかどうかをチェックし、適切なTokenTypeを返す
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// String はTokenの文字列表現を返す
func (t Token) String() string {
	return string(t.Type) + ":" + t.Literal
}
