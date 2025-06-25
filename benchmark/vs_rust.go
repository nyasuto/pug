package benchmark

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// RustBenchmark はRustとの比較ベンチマーク構造体
type RustBenchmark struct {
	Name           string
	PugSource      string
	RustSource     string
	ExpectedOutput string
	OptLevel       string // debug, release
}

// RustComparisonResult はRust比較ベンチマーク結果
type RustComparisonResult struct {
	TestName string
	OptLevel string

	// Pug結果
	PugCompileTime time.Duration
	PugExecuteTime time.Duration
	PugBinarySize  int64
	PugMemoryUsage int64
	PugSuccess     bool
	PugError       string

	// Rust結果
	RustCompileTime time.Duration
	RustExecuteTime time.Duration
	RustBinarySize  int64
	RustMemoryUsage int64
	RustSuccess     bool
	RustError       string

	// 比較指標
	CompileSpeedRatio float64 // Pug/Rust - 小さいほど高速
	RuntimeSpeedRatio float64 // Pug/Rust - 小さいほど高速
	BinarySizeRatio   float64 // Pug/Rust - 小さいほど効率的
	MemoryUsageRatio  float64 // Pug/Rust - 小さいほど効率的
}

// Rust比較テストケース
var rustBenchmarks = []RustBenchmark{
	{
		Name: "fibonacci_recursive",
		PugSource: `
let fib = fn(n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
};

let result = fib(20);
puts(result);
`,
		RustSource: `
fn fib(n: i32) -> i32 {
    if n <= 1 {
        return n;
    }
    fib(n - 1) + fib(n - 2)
}

fn main() {
    let result = fib(20);
    println!("{}", result);
}
`,
		ExpectedOutput: "6765",
	},

	{
		Name: "fibonacci_iterative",
		PugSource: `
let fib_iter = fn(n) {
    if (n <= 1) {
        return n;
    }
    let a = 0;
    let b = 1;
    let i = 2;
    while (i <= n) {
        let temp = a + b;
        a = b;
        b = temp;
        i = i + 1;
    }
    return b;
};

let result = fib_iter(30);
puts(result);
`,
		RustSource: `
fn fib_iter(n: i32) -> i32 {
    if n <= 1 {
        return n;
    }
    let mut a = 0;
    let mut b = 1;
    for _i in 2..=n {
        let temp = a + b;
        a = b;
        b = temp;
    }
    b
}

fn main() {
    let result = fib_iter(30);
    println!("{}", result);
}
`,
		ExpectedOutput: "832040",
	},

	{
		Name: "factorial",
		PugSource: `
let factorial = fn(n) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
};

let result = factorial(12);
puts(result);
`,
		RustSource: `
fn factorial(n: i32) -> i32 {
    if n <= 1 {
        1
    } else {
        n * factorial(n - 1)
    }
}

fn main() {
    let result = factorial(12);
    println!("{}", result);
}
`,
		ExpectedOutput: "479001600",
	},

	{
		Name: "nested_loops",
		PugSource: `
let nested_sum = fn(n) {
    let sum = 0;
    let i = 0;
    while (i < n) {
        let j = 0;
        while (j < n) {
            sum = sum + i * j;
            j = j + 1;
        }
        i = i + 1;
    }
    return sum;
};

let result = nested_sum(100);
puts(result);
`,
		RustSource: `
fn nested_sum(n: i32) -> i32 {
    let mut sum = 0;
    for i in 0..n {
        for j in 0..n {
            sum += i * j;
        }
    }
    sum
}

fn main() {
    let result = nested_sum(100);
    println!("{}", result);
}
`,
		ExpectedOutput: "24502500",
	},

	{
		Name: "vector_operations",
		PugSource: `
let vec_sum = fn(arr) {
    let sum = 0;
    let i = 0;
    while (i < len(arr)) {
        sum = sum + arr[i];
        i = i + 1;
    }
    return sum;
};

let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
let result = vec_sum(numbers);
puts(result);
`,
		RustSource: `
fn vec_sum(arr: &[i32]) -> i32 {
    arr.iter().sum()
}

fn main() {
    let numbers = vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    let result = vec_sum(&numbers);
    println!("{}", result);
}
`,
		ExpectedOutput: "55",
	},

	{
		Name: "string_processing",
		PugSource: `
let count_chars = fn(s) {
    return len(s);
};

let text = "Hello, Pug compiler world!";
let result = count_chars(text);
puts(result);
`,
		RustSource: `
fn count_chars(s: &str) -> usize {
    s.len()
}

fn main() {
    let text = "Hello, Pug compiler world!";
    let result = count_chars(text);
    println!("{}", result);
}
`,
		ExpectedOutput: "26",
	},

	{
		Name: "sorting_algorithm",
		PugSource: `
let quicksort = fn(arr) {
    if (len(arr) <= 1) {
        return arr;
    }
    
    let pivot = first(arr);
    let less = [];
    let greater = [];
    let rest_arr = rest(arr);
    
    let i = 0;
    while (i < len(rest_arr)) {
        let elem = rest_arr[i];
        if (elem < pivot) {
            less = push(less, elem);
        } else {
            greater = push(greater, elem);
        }
        i = i + 1;
    }
    
    return len(less) + 1 + len(greater);  // simplified
};

let numbers = [64, 34, 25, 12, 22, 11, 90, 5];
let result = quicksort(numbers);
puts(result);
`,
		RustSource: `
fn quicksort_len(arr: &[i32]) -> usize {
    if arr.len() <= 1 {
        return arr.len();
    }
    
    let pivot = arr[0];
    let mut less = Vec::new();
    let mut greater = Vec::new();
    
    for &elem in &arr[1..] {
        if elem < pivot {
            less.push(elem);
        } else {
            greater.push(elem);
        }
    }
    
    less.len() + 1 + greater.len()
}

fn main() {
    let numbers = vec![64, 34, 25, 12, 22, 11, 90, 5];
    let result = quicksort_len(&numbers);
    println!("{}", result);
}
`,
		ExpectedOutput: "8",
	},
}

// setupRustComparison はRust比較環境をセットアップ
func setupRustComparison(benchmark RustBenchmark, optLevel string) (*CompilerBenchmark, *CompilerBenchmark, string, error) {
	tempDir, err := os.MkdirTemp("", "rust_comparison_*")
	if err != nil {
		return nil, nil, "", fmt.Errorf("一時ディレクトリ作成失敗: %v", err)
	}

	// Pugファイル設定
	pugSource := filepath.Join(tempDir, benchmark.Name+".pug")
	pugBinary := filepath.Join(tempDir, benchmark.Name+"_pug")
	if runtime.GOOS == "windows" {
		pugBinary += ".exe"
	}

	err = os.WriteFile(pugSource, []byte(benchmark.PugSource), 0600)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, nil, "", fmt.Errorf("pugソースファイル作成失敗: %v", err)
	}

	// Rustプロジェクト設定
	rustProjectDir := filepath.Join(tempDir, "rust_"+benchmark.Name)
	err = os.MkdirAll(filepath.Join(rustProjectDir, "src"), 0750)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, nil, "", fmt.Errorf("rustプロジェクトディレクトリ作成失敗: %v", err)
	}

	// Cargo.toml作成
	cargoToml := fmt.Sprintf(`[package]
name = "%s"
version = "0.1.0"
edition = "2021"

[[bin]]
name = "%s"
path = "src/main.rs"
`, benchmark.Name, benchmark.Name)

	err = os.WriteFile(filepath.Join(rustProjectDir, "Cargo.toml"), []byte(cargoToml), 0600)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, nil, "", fmt.Errorf("cargo.toml作成失敗: %v", err)
	}

	// main.rs作成
	err = os.WriteFile(filepath.Join(rustProjectDir, "src", "main.rs"), []byte(benchmark.RustSource), 0600)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, nil, "", fmt.Errorf("rustソースファイル作成失敗: %v", err)
	}

	rustBinary := filepath.Join(rustProjectDir, "target", optLevel, benchmark.Name)
	if runtime.GOOS == "windows" {
		rustBinary += ".exe"
	}

	// Pugベンチマーク設定
	pugBench := &CompilerBenchmark{
		Phase:          "pug_vs_rust",
		SourceCode:     benchmark.PugSource,
		CompileCommand: []string{"./bin/pugc", pugSource, "-o", pugBinary},
		ExecuteCommand: []string{pugBinary},
		TempDir:        tempDir,
		SourceFile:     pugSource,
		BinaryFile:     pugBinary,
	}

	// Rustベンチマーク設定
	var cargoCmd []string
	if optLevel == "release" {
		cargoCmd = []string{"cargo", "build", "--release", "--manifest-path", filepath.Join(rustProjectDir, "Cargo.toml")}
	} else {
		cargoCmd = []string{"cargo", "build", "--manifest-path", filepath.Join(rustProjectDir, "Cargo.toml")}
	}

	rustBench := &CompilerBenchmark{
		Phase:          "rust",
		SourceCode:     benchmark.RustSource,
		CompileCommand: cargoCmd,
		ExecuteCommand: []string{rustBinary},
		TempDir:        tempDir,
		SourceFile:     filepath.Join(rustProjectDir, "src", "main.rs"),
		BinaryFile:     rustBinary,
	}

	return pugBench, rustBench, tempDir, nil
}

// runRustComparison はRustとの比較ベンチマークを実行
func runRustComparison(benchmark RustBenchmark, optLevel string, timeout time.Duration) *RustComparisonResult {
	result := &RustComparisonResult{
		TestName: benchmark.Name,
		OptLevel: optLevel,
	}

	// Rustの存在確認
	if _, err := exec.LookPath("cargo"); err != nil {
		result.RustError = "Rust/Cargoが見つかりません"
		return result
	}

	// Pugコンパイラの存在確認
	if _, err := os.Stat("./bin/pugc"); os.IsNotExist(err) {
		result.PugError = "Pugコンパイラが見つかりません"
		return result
	}

	pugBench, rustBench, tempDir, err := setupRustComparison(benchmark, optLevel)
	if err != nil {
		result.PugError = err.Error()
		result.RustError = err.Error()
		return result
	}
	defer os.RemoveAll(tempDir)

	// Pugベンチマーク実行
	pugBenchResult := pugBench.runBenchmark(timeout)
	result.PugSuccess = pugBenchResult.Success
	result.PugError = pugBenchResult.ErrorMessage
	result.PugCompileTime = pugBenchResult.CompileTime
	result.PugExecuteTime = pugBenchResult.ExecuteTime
	result.PugBinarySize = pugBenchResult.BinarySize
	result.PugMemoryUsage = pugBenchResult.MemoryUsage

	// Rustベンチマーク実行
	rustBenchResult := rustBench.runBenchmark(timeout)
	result.RustSuccess = rustBenchResult.Success
	result.RustError = rustBenchResult.ErrorMessage
	result.RustCompileTime = rustBenchResult.CompileTime
	result.RustExecuteTime = rustBenchResult.ExecuteTime
	result.RustBinarySize = rustBenchResult.BinarySize
	result.RustMemoryUsage = rustBenchResult.MemoryUsage

	// 比較指標計算
	if result.PugSuccess && result.RustSuccess {
		if result.RustCompileTime > 0 {
			result.CompileSpeedRatio = float64(result.PugCompileTime) / float64(result.RustCompileTime)
		}
		if result.RustExecuteTime > 0 {
			result.RuntimeSpeedRatio = float64(result.PugExecuteTime) / float64(result.RustExecuteTime)
		}
		if result.RustBinarySize > 0 {
			result.BinarySizeRatio = float64(result.PugBinarySize) / float64(result.RustBinarySize)
		}
		if result.RustMemoryUsage > 0 {
			result.MemoryUsageRatio = float64(result.PugMemoryUsage) / float64(result.RustMemoryUsage)
		}
	}

	return result
}

// BenchmarkVsRust_Debug はRust debugビルドとの比較ベンチマーク
func BenchmarkVsRust_Debug(b *testing.B) {
	runRustBenchmarkSuite(b, "debug")
}

// BenchmarkVsRust_Release はRust releaseビルドとの比較ベンチマーク
func BenchmarkVsRust_Release(b *testing.B) {
	runRustBenchmarkSuite(b, "release")
}

// runRustBenchmarkSuite はRust比較ベンチマークスイートを実行
func runRustBenchmarkSuite(b *testing.B, optLevel string) {
	if _, err := exec.LookPath("cargo"); err != nil {
		b.Skip("Rust/Cargoが見つかりません:", err)
		return
	}

	if _, err := os.Stat("./bin/pugc"); os.IsNotExist(err) {
		b.Skip("Pugコンパイラが見つかりません")
		return
	}

	var results []*RustComparisonResult

	for _, benchmark := range rustBenchmarks {
		b.Run(benchmark.Name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := runRustComparison(benchmark, optLevel, 120*time.Second)
				if i == 0 { // 最初の実行結果のみ保存
					results = append(results, result)
				}
			}
		})
	}

	// 結果レポート出力
	b.Cleanup(func() {
		printRustComparisonReport(results, optLevel)
	})
}

// printRustComparisonReport はRust比較結果レポートを出力
func printRustComparisonReport(results []*RustComparisonResult, optLevel string) {
	fmt.Print("\n")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("🦀 Pug vs Rust %s 比較ベンチマーク結果\n", optLevel)
	fmt.Println(strings.Repeat("=", 100))

	// ヘッダー
	fmt.Printf("%-20s %-15s %-15s %-15s %-15s %-15s %-15s\n",
		"テスト", "コンパイル比", "実行時間比", "バイナリ比", "メモリ比", "Pug状態", "Rust状態")
	fmt.Println(strings.Repeat("-", 100))

	var (
		avgCompileRatio float64
		avgRuntimeRatio float64
		avgBinaryRatio  float64
		avgMemoryRatio  float64
		successCount    int
	)

	for _, result := range results {
		pugStatus := "❌"
		if result.PugSuccess {
			pugStatus = "✅"
		}

		rustStatus := "❌"
		if result.RustSuccess {
			rustStatus = "✅"
		}

		compileRatio := "N/A"
		runtimeRatio := "N/A"
		binaryRatio := "N/A"
		memoryRatio := "N/A"

		if result.PugSuccess && result.RustSuccess {
			if result.CompileSpeedRatio > 0 {
				compileRatio = fmt.Sprintf("%.2fx", result.CompileSpeedRatio)
			}
			if result.RuntimeSpeedRatio > 0 {
				runtimeRatio = fmt.Sprintf("%.2fx", result.RuntimeSpeedRatio)
			}
			if result.BinarySizeRatio > 0 {
				binaryRatio = fmt.Sprintf("%.2fx", result.BinarySizeRatio)
			}
			if result.MemoryUsageRatio > 0 {
				memoryRatio = fmt.Sprintf("%.2fx", result.MemoryUsageRatio)
			}

			if result.CompileSpeedRatio > 0 {
				avgCompileRatio += result.CompileSpeedRatio
			}
			if result.RuntimeSpeedRatio > 0 {
				avgRuntimeRatio += result.RuntimeSpeedRatio
			}
			if result.BinarySizeRatio > 0 {
				avgBinaryRatio += result.BinarySizeRatio
			}
			if result.MemoryUsageRatio > 0 {
				avgMemoryRatio += result.MemoryUsageRatio
			}
			successCount++
		}

		fmt.Printf("%-20s %-15s %-15s %-15s %-15s %-15s %-15s\n",
			result.TestName, compileRatio, runtimeRatio, binaryRatio, memoryRatio, pugStatus, rustStatus)
	}

	fmt.Println(strings.Repeat("-", 100))

	// 平均値計算
	if successCount > 0 {
		avgCompileRatio /= float64(successCount)
		avgRuntimeRatio /= float64(successCount)
		avgBinaryRatio /= float64(successCount)
		avgMemoryRatio /= float64(successCount)

		fmt.Printf("%-20s %-15s %-15s %-15s %-15s\n",
			"平均",
			fmt.Sprintf("%.2fx", avgCompileRatio),
			fmt.Sprintf("%.2fx", avgRuntimeRatio),
			fmt.Sprintf("%.2fx", avgBinaryRatio),
			fmt.Sprintf("%.2fx", avgMemoryRatio))
	}

	fmt.Printf("\n📊 解析:\n")
	if successCount > 0 {
		fmt.Printf("🚀 平均実行時間比: %.2fx (1.0以下なら高速)\n", avgRuntimeRatio)
		fmt.Printf("⚡ 平均コンパイル比: %.2fx (1.0以下なら高速)\n", avgCompileRatio)
		fmt.Printf("📦 平均バイナリ比: %.2fx (1.0以下なら効率的)\n", avgBinaryRatio)
		fmt.Printf("💾 平均メモリ比: %.2fx (1.0以下なら効率的)\n", avgMemoryRatio)

		// Rust比較特別評価
		if avgRuntimeRatio <= 1.0 {
			fmt.Printf("🎉 実行速度: Rustと同等以上の性能! (ゼロコスト抽象化レベル)\n")
		} else if avgRuntimeRatio <= 2.0 {
			fmt.Printf("🦀 実行速度: 素晴らしい (Rustの2倍以内)\n")
		} else if avgRuntimeRatio <= 10.0 {
			fmt.Printf("⚠️  実行速度: 改善の余地あり (Rustより遅い)\n")
		} else if avgRuntimeRatio <= 100.0 {
			fmt.Printf("🔧 実行速度: 大幅改善が必要 (インタープリターレベル)\n")
		} else {
			fmt.Printf("⚠️  実行速度: 基本実装段階\n")
		}

		// コンパイル時間評価
		if avgCompileRatio <= 0.5 {
			fmt.Printf("⚡ コンパイル時間: Rustより高速! (Rustは重いコンパイラ)\n")
		} else if avgCompileRatio <= 1.0 {
			fmt.Printf("✅ コンパイル時間: Rustと同等 (良好)\n")
		} else if avgCompileRatio <= 2.0 {
			fmt.Printf("⚠️  コンパイル時間: やや遅い\n")
		} else {
			fmt.Printf("🔧 コンパイル時間: 改善が必要\n")
		}
	} else {
		fmt.Printf("❌ 比較可能な結果がありません\n")
	}

	fmt.Printf("\n🎯 Rust比較目標:\n")
	fmt.Printf("  Phase 1 (Interpreter): 100-1000x slower than Rust (学習段階)\n")
	fmt.Printf("  Phase 2 (Basic Compiler): 10-50x slower than Rust (基本実装)\n")
	fmt.Printf("  Phase 3 (Optimizer): 2-5x slower than Rust (最適化実装)\n")
	fmt.Printf("  Phase 4 (LLVM): Rustと同等性能 (ゼロコスト抽象化)\n")

	fmt.Printf("\n🦀 Rustの特徴:\n")
	fmt.Printf("  - ゼロコスト抽象化による高速実行\n")
	fmt.Printf("  - 強力な最適化コンパイラ\n")
	fmt.Printf("  - メモリ安全性とパフォーマンスの両立\n")
	fmt.Printf("  - 長いコンパイル時間 (トレードオフ)\n")

	fmt.Println(strings.Repeat("=", 100))
}

// BenchmarkPugVsRustEvolution はPugコンパイラのRust比較進化分析
func BenchmarkPugVsRustEvolution(b *testing.B) {
	optLevels := []string{"debug", "release"}
	phases := []string{"phase1"}

	// Phase2が利用可能かチェック
	if _, err := os.Stat("./bin/pugc"); err == nil {
		phases = append(phases, "phase2")
	}

	evolutionResults := make(map[string]map[string]*RustComparisonResult)

	for _, phase := range phases {
		evolutionResults[phase] = make(map[string]*RustComparisonResult)

		for _, optLevel := range optLevels {
			b.Run(fmt.Sprintf("%s_vs_rust_%s", phase, optLevel), func(b *testing.B) {
				// 代表的なベンチマーク（フィボナッチ反復）で測定
				fibBench := rustBenchmarks[1] // fibonacci_iterative

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					result := runRustComparison(fibBench, optLevel, 120*time.Second)
					if i == 0 {
						evolutionResults[phase][optLevel] = result
					}
				}
			})
		}
	}

	// 進化レポート出力
	b.Cleanup(func() {
		printRustEvolutionReport(evolutionResults)
	})
}

// printRustEvolutionReport はPug vs Rust進化レポートを出力
func printRustEvolutionReport(results map[string]map[string]*RustComparisonResult) {
	fmt.Print("\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("🦀 Pug vs Rust 進化分析\n")
	fmt.Println(strings.Repeat("=", 80))

	for phase, phaseResults := range results {
		fmt.Printf("\n📈 %s vs Rust:\n", strings.ToUpper(phase))
		fmt.Printf("%-12s %-15s %-15s %-10s\n", "Rust Build", "実行時間比", "コンパイル比", "状態")
		fmt.Println(strings.Repeat("-", 55))

		for _, optLevel := range []string{"debug", "release"} {
			if result, exists := phaseResults[optLevel]; exists {
				status := "❌"
				runtimeRatio := "N/A"
				compileRatio := "N/A"

				if result.PugSuccess && result.RustSuccess {
					status = "✅"
					if result.RuntimeSpeedRatio > 0 {
						runtimeRatio = fmt.Sprintf("%.2fx", result.RuntimeSpeedRatio)
					}
					if result.CompileSpeedRatio > 0 {
						compileRatio = fmt.Sprintf("%.2fx", result.CompileSpeedRatio)
					}
				}

				fmt.Printf("%-12s %-15s %-15s %-10s\n",
					optLevel, runtimeRatio, compileRatio, status)
			}
		}
	}

	fmt.Printf("\n🎯 進化目標 (vs Rust release):\n")
	fmt.Printf("  Phase 1 → Phase 2: 実行時間 1/10\n")
	fmt.Printf("  Phase 2 → Phase 3: 実行時間 1/5\n")
	fmt.Printf("  Phase 3 → Phase 4: Rustレベル到達\n")

	fmt.Printf("\n💡 学習ポイント:\n")
	fmt.Printf("  - Rustのコンパイル時間は長い（最適化のため）\n")
	fmt.Printf("  - Pugは軽量・高速コンパイルを目指す\n")
	fmt.Printf("  - 実行時性能でRustに追いつくのが目標\n")
	fmt.Println(strings.Repeat("=", 80))
}
