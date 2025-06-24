package phase1

import (
	"testing"
)

// ベンチマーク用のサンプルコード
const benchmarkCode = `
let fibonacci = fn(n: int) -> int {
    if n <= 1 {
        return n;
    }
    return fibonacci(n-1) + fibonacci(n-2);
};

let factorial = fn(n: int) -> int {
    if n <= 1 {
        return 1;
    }
    return n * factorial(n-1);
};

let quicksort = fn(arr: [int], low: int, high: int) -> [int] {
    if low < high {
        let pi = partition(arr, low, high);
        quicksort(arr, low, pi - 1);
        quicksort(arr, pi + 1, high);
    }
    return arr;
};

let partition = fn(arr: [int], low: int, high: int) -> int {
    let pivot = arr[high];
    let i = low - 1;
    
    for j in low..high {
        if arr[j] < pivot {
            i = i + 1;
            swap(arr, i, j);
        }
    }
    swap(arr, i + 1, high);
    return i + 1;
};

let swap = fn(arr: [int], i: int, j: int) {
    let temp = arr[i];
    arr[i] = arr[j];
    arr[j] = temp;
};

// メイン処理
let numbers = [64, 34, 25, 12, 22, 11, 90];
let sorted = quicksort(numbers, 0, 6);

let result1 = fibonacci(10);
let result2 = factorial(5);

// 浮動小数点数計算
let pi: float = 3.14159;
let radius: float = 5.0;
let area: float = pi * radius * radius;

// 文字列処理
let message: string = "Hello, pug compiler!";
let greeting: string = "Welcome to Phase 1.0";

// 論理演算
let condition1: bool = true;
let condition2: bool = false;
let result3: bool = condition1 && !condition2;

// 複雑な式
let complex = (result1 + result2) * 2 - factorial(3);
let comparison = complex > 100 && area < 100.0;

// ループとブレーク
while (condition1) {
    if (complex > 50) {
        break;
    }
    complex = complex + 1;
    if (complex % 2 == 0) {
        continue;
    }
}

// コメント付きコード
// これは複雑な計算です
let advanced_calc = fn(x: float, y: float) -> float {
    // 二次方程式の解
    let discriminant = y * y - 4.0 * x * 1.0;
    if discriminant >= 0.0 {
        return (-y + sqrt(discriminant)) / (2.0 * x);
    } else {
        return 0.0; // 複素数解は未対応
    }
};
`

// BenchmarkLexer は基本的なレクサー性能をベンチマーク
func BenchmarkLexer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := New(benchmarkCode)
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkLexerSmallInput は小さな入力に対するレクサー性能をベンチマーク
func BenchmarkLexerSmallInput(b *testing.B) {
	input := "let x = 5 + 10;"
	for i := 0; i < b.N; i++ {
		l := New(input)
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkLexerMediumInput は中程度の入力に対するレクサー性能をベンチマーク
func BenchmarkLexerMediumInput(b *testing.B) {
	input := `
	let add = fn(x: int, y: int) -> int {
		return x + y;
	};
	let result = add(5, 10);
	if result > 10 {
		return true;
	} else {
		return false;
	}
	`
	for i := 0; i < b.N; i++ {
		l := New(input)
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkLexerNumbers は数値トークン化の性能をベンチマーク
func BenchmarkLexerNumbers(b *testing.B) {
	input := "123 456 789 3.14 2.71828 0.5 1000.999"
	for i := 0; i < b.N; i++ {
		l := New(input)
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkLexerStrings は文字列トークン化の性能をベンチマーク
func BenchmarkLexerStrings(b *testing.B) {
	input := `"hello" "world" "this is a longer string" "special\ncharacters\t"`
	for i := 0; i < b.N; i++ {
		l := New(input)
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkLexerIdentifiers は識別子トークン化の性能をベンチマーク
func BenchmarkLexerIdentifiers(b *testing.B) {
	input := "variable_name camelCase snake_case _private identifier123 let fn if else return"
	for i := 0; i < b.N; i++ {
		l := New(input)
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkLexerOperators は演算子トークン化の性能をベンチマーク
func BenchmarkLexerOperators(b *testing.B) {
	input := "+ - * / % == != < > <= >= && || ! = -> , ; : ( ) { } [ ]"
	for i := 0; i < b.N; i++ {
		l := New(input)
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkLexerComments はコメント処理の性能をベンチマーク
func BenchmarkLexerComments(b *testing.B) {
	input := `
	let x = 5; // これはコメント
	// 完全にコメントの行
	let y = 10; // 別のコメント
	// もう一つのコメント行
	let z = x + y;
	`
	for i := 0; i < b.N; i++ {
		l := New(input)
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkLexerMixedContent は様々なトークンが混在した内容の性能をベンチマーク
func BenchmarkLexerMixedContent(b *testing.B) {
	input := `
	let result = calculate(123, 45.6, "text", true);
	if (result >= 100.0 && flag != false) {
		return "success";
	} else {
		return "failure";
	}
	`
	for i := 0; i < b.N; i++ {
		l := New(input)
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkNewLexer はレクサー作成自体の性能をベンチマーク
func BenchmarkNewLexer(b *testing.B) {
	input := benchmarkCode
	for i := 0; i < b.N; i++ {
		_ = New(input)
	}
}

// BenchmarkNextTokenOnly は単一のNextToken呼び出しの性能をベンチマーク
func BenchmarkNextTokenOnly(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New("let x = 5;") // 毎回リセット
		l.NextToken()          // 1回だけ呼び出し
	}
}

// BenchmarkTokenCreation はトークン作成の性能をベンチマーク
func BenchmarkTokenCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Token{
			Type:     LET,
			Literal:  "let",
			Line:     1,
			Column:   1,
			Position: 0,
		}
	}
}

// BenchmarkLookupIdent はキーワード検索の性能をベンチマーク
func BenchmarkLookupIdent(b *testing.B) {
	identifiers := []string{"let", "fn", "if", "else", "return", "notakeyword", "variable"}
	for i := 0; i < b.N; i++ {
		for _, ident := range identifiers {
			_ = LookupIdent(ident)
		}
	}
}
