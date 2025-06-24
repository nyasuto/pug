package phase1

// Lexer は字句解析器の構造体
type Lexer struct {
	input        string // 解析対象の入力文字列
	position     int    // 現在の文字位置（current charを指す）
	readPosition int    // 次に読む文字位置（current charの次）
	ch           byte   // 現在検査中の文字
	line         int    // 現在の行番号（エラー報告用）
	column       int    // 現在の列番号（エラー報告用）
}

// New は新しいLexerを作成する
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar() // 最初の文字を読み込む
	return l
}

// readChar は次の文字を読み込み、現在位置を進める
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

// peekChar は次の文字を先読みする（位置は進めない）
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace は空白文字をスキップする
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment は行コメント（//）をスキップする
func (l *Lexer) skipComment() {
	if l.ch == '/' && l.peekChar() == '/' {
		// 行末まで読み飛ばす
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	}
}

// readIdentifier は識別子を読み取る
func (l *Lexer) readIdentifier() string {
	position := l.position

	// 最初の文字はアルファベットまたはアンダースコア
	if isLetter(l.ch) {
		l.readChar()

		// 2文字目以降はアルファベット、数字、アンダースコア
		for isLetter(l.ch) || isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

// readNumber は数値を読み取る（整数・浮動小数点数対応）
func (l *Lexer) readNumber() (string, TokenType) {
	position := l.position
	tokenType := INT

	// 整数部分を読む
	for isDigit(l.ch) {
		l.readChar()
	}

	// 小数点がある場合は浮動小数点数
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = FLOAT
		l.readChar() // '.'をスキップ

		// 小数部分を読む
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position], tokenType
}

// readString は文字列リテラルを読み取る
func (l *Lexer) readString() string {
	result := ""

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		// エスケープシーケンスの処理
		if l.ch == '\\' {
			l.readChar() // バックスラッシュをスキップ
			if l.ch != 0 {
				switch l.ch {
				case 'n':
					result += "\n"
				case 't':
					result += "\t"
				case 'r':
					result += "\r"
				case '\\':
					result += "\\"
				case '"':
					result += "\\\""
				default:
					// 未対応のエスケープシーケンスはそのまま
					result += "\\" + string(l.ch)
				}
			}
		} else {
			result += string(l.ch)
		}
	}

	return result
}

// isLetter は文字がアルファベットまたはアンダースコアかチェックする
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit は文字が数字かチェックする
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// NextToken は次のトークンを解析して返す
func (l *Lexer) NextToken() Token {
	var tok Token

	// 空白とコメントをスキップ
	for {
		l.skipWhitespace()

		// コメントチェック
		if l.ch == '/' && l.peekChar() == '/' {
			l.skipComment()
			continue
		}
		break
	}

	// 現在位置を記録
	tok.Line = l.line
	tok.Column = l.column
	tok.Position = l.position

	switch l.ch {
	case '.':
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			tok.Type = IDENT
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = EQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = ASSIGN
			tok.Literal = string(l.ch)
		}
	case '+':
		tok.Type = PLUS
		tok.Literal = string(l.ch)
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok.Type = ARROW
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = MINUS
			tok.Literal = string(l.ch)
		}
	case '*':
		tok.Type = MULTIPLY
		tok.Literal = string(l.ch)
	case '/':
		tok.Type = DIVIDE
		tok.Literal = string(l.ch)
	case '%':
		tok.Type = MODULO
		tok.Literal = string(l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = NOT_EQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = NOT
			tok.Literal = string(l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = LTE
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = LT
			tok.Literal = string(l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = GTE
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = GT
			tok.Literal = string(l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok.Type = AND
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok.Type = OR
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch)
		}
	case ',':
		tok.Type = COMMA
		tok.Literal = string(l.ch)
	case ';':
		tok.Type = SEMICOLON
		tok.Literal = string(l.ch)
	case ':':
		tok.Type = COLON
		tok.Literal = string(l.ch)
	case '(':
		tok.Type = LPAREN
		tok.Literal = string(l.ch)
	case ')':
		tok.Type = RPAREN
		tok.Literal = string(l.ch)
	case '{':
		tok.Type = LBRACE
		tok.Literal = string(l.ch)
	case '}':
		tok.Type = RBRACE
		tok.Literal = string(l.ch)
	case '[':
		tok.Type = LBRACKET
		tok.Literal = string(l.ch)
	case ']':
		tok.Type = RBRACKET
		tok.Literal = string(l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok // readIdentifierで既に位置が進んでいるため
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.readNumber()
			return tok // readNumberで既に位置が進んでいるため
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch)
		}
	}

	l.readChar()
	return tok
}
