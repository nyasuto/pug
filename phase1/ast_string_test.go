package phase1

import (
	"testing"
)

// TestASTStringMethods は各ASTノードのString()メソッドをテストする
func TestASTStringMethods(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{
			name: "Program with statements",
			node: &Program{
				Statements: []Statement{
					&ExpressionStatement{
						Token: Token{Type: INT, Literal: "5"},
						Expression: &IntegerLiteral{
							Token: Token{Type: INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			expected: "5",
		},
		{
			name: "Empty Program",
			node: &Program{
				Statements: []Statement{},
			},
			expected: "",
		},
		{
			name: "LetStatement with value",
			node: &LetStatement{
				Token: Token{Type: LET, Literal: "let"},
				Name: &Identifier{
					Token: Token{Type: IDENT, Literal: "x"},
					Value: "x",
				},
				Value: &IntegerLiteral{
					Token: Token{Type: INT, Literal: "10"},
					Value: 10,
				},
			},
			expected: "let x = 10;",
		},
		{
			name: "LetStatement without value",
			node: &LetStatement{
				Token: Token{Type: LET, Literal: "let"},
				Name: &Identifier{
					Token: Token{Type: IDENT, Literal: "y"},
					Value: "y",
				},
				Value: nil,
			},
			expected: "let y = ;",
		},
		{
			name: "ReturnStatement with value",
			node: &ReturnStatement{
				Token: Token{Type: RETURN, Literal: "return"},
				ReturnValue: &IntegerLiteral{
					Token: Token{Type: INT, Literal: "42"},
					Value: 42,
				},
			},
			expected: "return 42;",
		},
		{
			name: "ReturnStatement without value",
			node: &ReturnStatement{
				Token:       Token{Type: RETURN, Literal: "return"},
				ReturnValue: nil,
			},
			expected: "return ;",
		},
		{
			name: "ExpressionStatement with expression",
			node: &ExpressionStatement{
				Token: Token{Type: INT, Literal: "100"},
				Expression: &IntegerLiteral{
					Token: Token{Type: INT, Literal: "100"},
					Value: 100,
				},
			},
			expected: "100",
		},
		{
			name: "ExpressionStatement without expression",
			node: &ExpressionStatement{
				Token:      Token{Type: INT, Literal: ""},
				Expression: nil,
			},
			expected: "",
		},
		{
			name: "StringLiteral",
			node: &StringLiteral{
				Token: Token{Type: STRING, Literal: "hello"},
				Value: "hello",
			},
			expected: "\"hello\"",
		},
		{
			name: "PrefixExpression",
			node: &PrefixExpression{
				Token:    Token{Type: MINUS, Literal: "-"},
				Operator: "-",
				Right: &IntegerLiteral{
					Token: Token{Type: INT, Literal: "5"},
					Value: 5,
				},
			},
			expected: "(-5)",
		},
		{
			name: "InfixExpression",
			node: &InfixExpression{
				Token: Token{Type: PLUS, Literal: "+"},
				Left: &IntegerLiteral{
					Token: Token{Type: INT, Literal: "1"},
					Value: 1,
				},
				Operator: "+",
				Right: &IntegerLiteral{
					Token: Token{Type: INT, Literal: "2"},
					Value: 2,
				},
			},
			expected: "(1 + 2)",
		},
		{
			name: "IfExpression with alternative",
			node: &IfExpression{
				Token: Token{Type: IF, Literal: "if"},
				Condition: &Boolean{
					Token: Token{Type: TRUE, Literal: "true"},
					Value: true,
				},
				Consequence: &BlockStatement{
					Token: Token{Type: LBRACE, Literal: "{"},
					Statements: []Statement{
						&ExpressionStatement{
							Token: Token{Type: INT, Literal: "1"},
							Expression: &IntegerLiteral{
								Token: Token{Type: INT, Literal: "1"},
								Value: 1,
							},
						},
					},
				},
				Alternative: &BlockStatement{
					Token: Token{Type: LBRACE, Literal: "{"},
					Statements: []Statement{
						&ExpressionStatement{
							Token: Token{Type: INT, Literal: "2"},
							Expression: &IntegerLiteral{
								Token: Token{Type: INT, Literal: "2"},
								Value: 2,
							},
						},
					},
				},
			},
			expected: "iftrue 1else 2",
		},
		{
			name: "IfExpression without alternative",
			node: &IfExpression{
				Token: Token{Type: IF, Literal: "if"},
				Condition: &Boolean{
					Token: Token{Type: TRUE, Literal: "true"},
					Value: true,
				},
				Consequence: &BlockStatement{
					Token: Token{Type: LBRACE, Literal: "{"},
					Statements: []Statement{
						&ExpressionStatement{
							Token: Token{Type: INT, Literal: "1"},
							Expression: &IntegerLiteral{
								Token: Token{Type: INT, Literal: "1"},
								Value: 1,
							},
						},
					},
				},
				Alternative: nil,
			},
			expected: "iftrue 1",
		},
		{
			name: "FunctionLiteral with parameters",
			node: &FunctionLiteral{
				Token: Token{Type: FN, Literal: "fn"},
				Parameters: []*Identifier{
					{
						Token: Token{Type: IDENT, Literal: "x"},
						Value: "x",
					},
					{
						Token: Token{Type: IDENT, Literal: "y"},
						Value: "y",
					},
				},
				Body: &BlockStatement{
					Token: Token{Type: LBRACE, Literal: "{"},
					Statements: []Statement{
						&ExpressionStatement{
							Token: Token{Type: INT, Literal: "42"},
							Expression: &IntegerLiteral{
								Token: Token{Type: INT, Literal: "42"},
								Value: 42,
							},
						},
					},
				},
			},
			expected: "fn(x, y) 42",
		},
		{
			name: "CallExpression with arguments",
			node: &CallExpression{
				Token: Token{Type: LPAREN, Literal: "("},
				Function: &Identifier{
					Token: Token{Type: IDENT, Literal: "add"},
					Value: "add",
				},
				Arguments: []Expression{
					&IntegerLiteral{
						Token: Token{Type: INT, Literal: "1"},
						Value: 1,
					},
					&IntegerLiteral{
						Token: Token{Type: INT, Literal: "2"},
						Value: 2,
					},
				},
			},
			expected: "add(1, 2)",
		},
		{
			name: "WhileStatement",
			node: &WhileStatement{
				Token: Token{Type: WHILE, Literal: "while"},
				Condition: &Boolean{
					Token: Token{Type: TRUE, Literal: "true"},
					Value: true,
				},
				Body: &BlockStatement{
					Token: Token{Type: LBRACE, Literal: "{"},
					Statements: []Statement{
						&ExpressionStatement{
							Token: Token{Type: INT, Literal: "1"},
							Expression: &IntegerLiteral{
								Token: Token{Type: INT, Literal: "1"},
								Value: 1,
							},
						},
					},
				},
			},
			expected: "while true 1",
		},
		{
			name: "ForStatement with all components",
			node: &ForStatement{
				Token: Token{Type: FOR, Literal: "for"},
				Initializer: &LetStatement{
					Token: Token{Type: LET, Literal: "let"},
					Name: &Identifier{
						Token: Token{Type: IDENT, Literal: "i"},
						Value: "i",
					},
					Value: &IntegerLiteral{
						Token: Token{Type: INT, Literal: "0"},
						Value: 0,
					},
				},
				Condition: &InfixExpression{
					Token: Token{Type: LT, Literal: "<"},
					Left: &Identifier{
						Token: Token{Type: IDENT, Literal: "i"},
						Value: "i",
					},
					Operator: "<",
					Right: &IntegerLiteral{
						Token: Token{Type: INT, Literal: "10"},
						Value: 10,
					},
				},
				Update: &InfixExpression{
					Token: Token{Type: PLUS, Literal: "+"},
					Left: &Identifier{
						Token: Token{Type: IDENT, Literal: "i"},
						Value: "i",
					},
					Operator: "+",
					Right: &IntegerLiteral{
						Token: Token{Type: INT, Literal: "1"},
						Value: 1,
					},
				},
				Body: &BlockStatement{
					Token: Token{Type: LBRACE, Literal: "{"},
					Statements: []Statement{
						&ExpressionStatement{
							Token: Token{Type: INT, Literal: "42"},
							Expression: &IntegerLiteral{
								Token: Token{Type: INT, Literal: "42"},
								Value: 42,
							},
						},
					},
				},
			},
			expected: "for let i = 0;; (i < 10); (i + 1) 42",
		},
		{
			name: "ForStatement with partial components",
			node: &ForStatement{
				Token:       Token{Type: FOR, Literal: "for"},
				Initializer: nil,
				Condition: &Boolean{
					Token: Token{Type: TRUE, Literal: "true"},
					Value: true,
				},
				Update: nil,
				Body: &BlockStatement{
					Token:      Token{Type: LBRACE, Literal: "{"},
					Statements: []Statement{},
				},
			},
			expected: "for ; true;  ",
		},
		{
			name: "BreakStatement",
			node: &BreakStatement{
				Token: Token{Type: BREAK, Literal: "break"},
			},
			expected: "break;",
		},
		{
			name: "ContinueStatement",
			node: &ContinueStatement{
				Token: Token{Type: CONTINUE, Literal: "continue"},
			},
			expected: "continue;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestTokenLiteralMethods は各ASTノードのTokenLiteral()メソッドをテストする
func TestTokenLiteralMethods(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{
			name: "Program with statements",
			node: &Program{
				Statements: []Statement{
					&ExpressionStatement{
						Token: Token{Type: INT, Literal: "5"},
						Expression: &IntegerLiteral{
							Token: Token{Type: INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			expected: "5",
		},
		{
			name: "Empty Program",
			node: &Program{
				Statements: []Statement{},
			},
			expected: "",
		},
		{
			name: "LetStatement",
			node: &LetStatement{
				Token: Token{Type: LET, Literal: "let"},
			},
			expected: "let",
		},
		{
			name: "ReturnStatement",
			node: &ReturnStatement{
				Token: Token{Type: RETURN, Literal: "return"},
			},
			expected: "return",
		},
		{
			name: "WhileStatement",
			node: &WhileStatement{
				Token: Token{Type: WHILE, Literal: "while"},
			},
			expected: "while",
		},
		{
			name: "ForStatement",
			node: &ForStatement{
				Token: Token{Type: FOR, Literal: "for"},
			},
			expected: "for",
		},
		{
			name: "BreakStatement",
			node: &BreakStatement{
				Token: Token{Type: BREAK, Literal: "break"},
			},
			expected: "break",
		},
		{
			name: "ContinueStatement",
			node: &ContinueStatement{
				Token: Token{Type: CONTINUE, Literal: "continue"},
			},
			expected: "continue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.TokenLiteral()
			if result != tt.expected {
				t.Errorf("TokenLiteral() = %q, want %q", result, tt.expected)
			}
		})
	}
}
