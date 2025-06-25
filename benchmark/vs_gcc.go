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

// GCCBenchmark はGCCとの比較ベンチマーク構造体
type GCCBenchmark struct {
	Name           string
	PugSource      string
	CSource        string
	ExpectedOutput string
	OptLevel       string // -O0, -O1, -O2, -O3
}

// ComparisonResult は比較ベンチマーク結果
type ComparisonResult struct {
	TestName string
	OptLevel string

	// Pug結果
	PugCompileTime time.Duration
	PugExecuteTime time.Duration
	PugBinarySize  int64
	PugMemoryUsage int64
	PugSuccess     bool
	PugError       string

	// GCC結果
	GCCCompileTime time.Duration
	GCCExecuteTime time.Duration
	GCCBinarySize  int64
	GCCMemoryUsage int64
	GCCSuccess     bool
	GCCError       string

	// 比較指標
	CompileSpeedRatio float64 // Pug/GCC - 小さいほど高速
	RuntimeSpeedRatio float64 // Pug/GCC - 小さいほど高速
	BinarySizeRatio   float64 // Pug/GCC - 小さいほど効率的
	MemoryUsageRatio  float64 // Pug/GCC - 小さいほど効率的
}

// GCC比較テストケース
var gccBenchmarks = []GCCBenchmark{
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
		CSource: `
#include <stdio.h>

int fib(int n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
}

int main() {
    int result = fib(20);
    printf("%d\n", result);
    return 0;
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
		CSource: `
#include <stdio.h>

int fib_iter(int n) {
    if (n <= 1) {
        return n;
    }
    int a = 0, b = 1, temp;
    for (int i = 2; i <= n; i++) {
        temp = a + b;
        a = b;
        b = temp;
    }
    return b;
}

int main() {
    int result = fib_iter(30);
    printf("%d\n", result);
    return 0;
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
		CSource: `
#include <stdio.h>

int factorial(int n) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
}

int main() {
    int result = factorial(12);
    printf("%d\n", result);
    return 0;
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
		CSource: `
#include <stdio.h>

int nested_sum(int n) {
    int sum = 0;
    for (int i = 0; i < n; i++) {
        for (int j = 0; j < n; j++) {
            sum += i * j;
        }
    }
    return sum;
}

int main() {
    int result = nested_sum(100);
    printf("%d\n", result);
    return 0;
}
`,
		ExpectedOutput: "24502500",
	},

	{
		Name: "array_operations",
		PugSource: `
let array_sum = fn(arr) {
    let sum = 0;
    let i = 0;
    while (i < len(arr)) {
        sum = sum + arr[i];
        i = i + 1;
    }
    return sum;
};

let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
let result = array_sum(numbers);
puts(result);
`,
		CSource: `
#include <stdio.h>

int array_sum(int arr[], int size) {
    int sum = 0;
    for (int i = 0; i < size; i++) {
        sum += arr[i];
    }
    return sum;
}

int main() {
    int numbers[] = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
    int size = sizeof(numbers) / sizeof(numbers[0]);
    int result = array_sum(numbers, size);
    printf("%d\n", result);
    return 0;
}
`,
		ExpectedOutput: "55",
	},
}

// setupGCCComparison はGCC比較環境をセットアップ
func setupGCCComparison(benchmark GCCBenchmark, optLevel string) (*CompilerBenchmark, *CompilerBenchmark, string, error) {
	tempDir, err := os.MkdirTemp("", "gcc_comparison_*")
	if err != nil {
		return nil, nil, "", fmt.Errorf("一時ディレクトリ作成失敗: %v", err)
	}

	// Pugファイル設定
	pugSource := filepath.Join(tempDir, benchmark.Name+".pug")
	pugBinary := filepath.Join(tempDir, benchmark.Name+"_pug")
	if runtime.GOOS == "windows" {
		pugBinary += ".exe"
	}

	err = os.WriteFile(pugSource, []byte(benchmark.PugSource), 0644)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, nil, "", fmt.Errorf("pugソースファイル作成失敗: %v", err)
	}

	// Cファイル設定
	cSource := filepath.Join(tempDir, benchmark.Name+".c")
	cBinary := filepath.Join(tempDir, benchmark.Name+"_gcc")
	if runtime.GOOS == "windows" {
		cBinary += ".exe"
	}

	err = os.WriteFile(cSource, []byte(benchmark.CSource), 0644)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, nil, "", fmt.Errorf("cソースファイル作成失敗: %v", err)
	}

	// Pugベンチマーク設定
	pugBench := &CompilerBenchmark{
		Phase:          "pug_vs_gcc",
		SourceCode:     benchmark.PugSource,
		CompileCommand: []string{"./bin/pugc", pugSource, "-o", pugBinary},
		ExecuteCommand: []string{pugBinary},
		TempDir:        tempDir,
		SourceFile:     pugSource,
		BinaryFile:     pugBinary,
	}

	// GCCベンチマーク設定
	gccBench := &CompilerBenchmark{
		Phase:          "gcc",
		SourceCode:     benchmark.CSource,
		CompileCommand: []string{"gcc", optLevel, cSource, "-o", cBinary},
		ExecuteCommand: []string{cBinary},
		TempDir:        tempDir,
		SourceFile:     cSource,
		BinaryFile:     cBinary,
	}

	return pugBench, gccBench, tempDir, nil
}

// runGCCComparison はGCCとの比較ベンチマークを実行
func runGCCComparison(benchmark GCCBenchmark, optLevel string, timeout time.Duration) *ComparisonResult {
	result := &ComparisonResult{
		TestName: benchmark.Name,
		OptLevel: optLevel,
	}

	// GCCの存在確認
	if _, err := exec.LookPath("gcc"); err != nil {
		result.GCCError = "GCCが見つかりません"
		return result
	}

	// Pugコンパイラの存在確認
	if _, err := os.Stat("./bin/pugc"); os.IsNotExist(err) {
		result.PugError = "Pugコンパイラが見つかりません"
		return result
	}

	pugBench, gccBench, tempDir, err := setupGCCComparison(benchmark, optLevel)
	if err != nil {
		result.PugError = err.Error()
		result.GCCError = err.Error()
		return result
	}
	defer os.RemoveAll(tempDir)

	// Timeout control is handled in runBenchmark method

	// Pugベンチマーク実行
	pugResult := pugBench.runBenchmark(timeout)
	result.PugSuccess = pugResult.Success
	result.PugError = pugResult.ErrorMessage
	result.PugCompileTime = pugResult.CompileTime
	result.PugExecuteTime = pugResult.ExecuteTime
	result.PugBinarySize = pugResult.BinarySize
	result.PugMemoryUsage = pugResult.MemoryUsage

	// GCCベンチマーク実行
	gccResult := gccBench.runBenchmark(timeout)
	result.GCCSuccess = gccResult.Success
	result.GCCError = gccResult.ErrorMessage
	result.GCCCompileTime = gccResult.CompileTime
	result.GCCExecuteTime = gccResult.ExecuteTime
	result.GCCBinarySize = gccResult.BinarySize
	result.GCCMemoryUsage = gccResult.MemoryUsage

	// 比較指標計算
	if result.PugSuccess && result.GCCSuccess {
		result.CompileSpeedRatio = float64(result.PugCompileTime) / float64(result.GCCCompileTime)
		result.RuntimeSpeedRatio = float64(result.PugExecuteTime) / float64(result.GCCExecuteTime)

		if result.GCCBinarySize > 0 {
			result.BinarySizeRatio = float64(result.PugBinarySize) / float64(result.GCCBinarySize)
		}

		if result.GCCMemoryUsage > 0 {
			result.MemoryUsageRatio = float64(result.PugMemoryUsage) / float64(result.GCCMemoryUsage)
		}
	}

	return result
}

// BenchmarkVsGCC_O0 はGCC -O0との比較ベンチマーク
func BenchmarkVsGCC_O0(b *testing.B) {
	runGCCBenchmarkSuite(b, "-O0")
}

// BenchmarkVsGCC_O1 はGCC -O1との比較ベンチマーク
func BenchmarkVsGCC_O1(b *testing.B) {
	runGCCBenchmarkSuite(b, "-O1")
}

// BenchmarkVsGCC_O2 はGCC -O2との比較ベンチマーク
func BenchmarkVsGCC_O2(b *testing.B) {
	runGCCBenchmarkSuite(b, "-O2")
}

// BenchmarkVsGCC_O3 はGCC -O3との比較ベンチマーク
func BenchmarkVsGCC_O3(b *testing.B) {
	runGCCBenchmarkSuite(b, "-O3")
}

// runGCCBenchmarkSuite はGCC比較ベンチマークスイートを実行
func runGCCBenchmarkSuite(b *testing.B, optLevel string) {
	if _, err := exec.LookPath("gcc"); err != nil {
		b.Skip("GCCが見つかりません:", err)
		return
	}

	if _, err := os.Stat("./bin/pugc"); os.IsNotExist(err) {
		b.Skip("Pugコンパイラが見つかりません")
		return
	}

	var results []*ComparisonResult

	for _, benchmark := range gccBenchmarks {
		b.Run(benchmark.Name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := runGCCComparison(benchmark, optLevel, 60*time.Second)
				if i == 0 { // 最初の実行結果のみ保存
					results = append(results, result)
				}
			}
		})
	}

	// 結果レポート出力
	b.Cleanup(func() {
		printGCCComparisonReport(results, optLevel)
	})
}

// printGCCComparisonReport はGCC比較結果レポートを出力
func printGCCComparisonReport(results []*ComparisonResult, optLevel string) {
	fmt.Print("\n")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("🏁 Pug vs GCC %s 比較ベンチマーク結果\n", optLevel)
	fmt.Println(strings.Repeat("=", 100))

	// ヘッダー
	fmt.Printf("%-20s %-15s %-15s %-15s %-15s %-15s %-15s\n",
		"テスト", "コンパイル比", "実行時間比", "バイナリ比", "メモリ比", "Pug状態", "GCC状態")
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

		gccStatus := "❌"
		if result.GCCSuccess {
			gccStatus = "✅"
		}

		compileRatio := "N/A"
		runtimeRatio := "N/A"
		binaryRatio := "N/A"
		memoryRatio := "N/A"

		if result.PugSuccess && result.GCCSuccess {
			compileRatio = fmt.Sprintf("%.2fx", result.CompileSpeedRatio)
			runtimeRatio = fmt.Sprintf("%.2fx", result.RuntimeSpeedRatio)
			binaryRatio = fmt.Sprintf("%.2fx", result.BinarySizeRatio)
			memoryRatio = fmt.Sprintf("%.2fx", result.MemoryUsageRatio)

			avgCompileRatio += result.CompileSpeedRatio
			avgRuntimeRatio += result.RuntimeSpeedRatio
			avgBinaryRatio += result.BinarySizeRatio
			avgMemoryRatio += result.MemoryUsageRatio
			successCount++
		}

		fmt.Printf("%-20s %-15s %-15s %-15s %-15s %-15s %-15s\n",
			result.TestName, compileRatio, runtimeRatio, binaryRatio, memoryRatio, pugStatus, gccStatus)
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

		// 性能評価
		if avgRuntimeRatio <= 1.0 {
			fmt.Printf("🎉 実行速度: GCCと同等以上の性能!\n")
		} else if avgRuntimeRatio <= 2.0 {
			fmt.Printf("✅ 実行速度: 良好 (GCCの2倍以内)\n")
		} else if avgRuntimeRatio <= 10.0 {
			fmt.Printf("⚠️  実行速度: 改善の余地あり\n")
		} else {
			fmt.Printf("🔧 実行速度: 大幅改善が必要\n")
		}
	} else {
		fmt.Printf("❌ 比較可能な結果がありません\n")
	}

	fmt.Printf("\n🎯 目標:\n")
	fmt.Printf("  Phase 1 (Interpreter): 10-100x slower than GCC (学習段階)\n")
	fmt.Printf("  Phase 2 (Basic Compiler): 2-10x slower than GCC (基本実装)\n")
	fmt.Printf("  Phase 3 (Optimizer): 1-2x slower than GCC (最適化実装)\n")
	fmt.Printf("  Phase 4 (LLVM): GCCと同等性能 (産業レベル)\n")

	fmt.Println(strings.Repeat("=", 100))
}

// BenchmarkPugCompilerEvolution はPugコンパイラの進化をGCCと比較
func BenchmarkPugCompilerEvolution(b *testing.B) {
	optLevels := []string{"-O0", "-O1", "-O2", "-O3"}
	phases := []string{"phase1"}

	// Phase2が利用可能かチェック
	if _, err := os.Stat("./bin/pugc"); err == nil {
		phases = append(phases, "phase2")
	}

	evolutionResults := make(map[string]map[string]*ComparisonResult)

	for _, phase := range phases {
		evolutionResults[phase] = make(map[string]*ComparisonResult)

		for _, optLevel := range optLevels {
			b.Run(fmt.Sprintf("%s_vs_gcc_%s", phase, strings.TrimPrefix(optLevel, "-")), func(b *testing.B) {
				// 代表的なベンチマーク（フィボナッチ）で測定
				fibBench := gccBenchmarks[0] // fibonacci_recursive

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					result := runGCCComparison(fibBench, optLevel, 60*time.Second)
					if i == 0 {
						evolutionResults[phase][optLevel] = result
					}
				}
			})
		}
	}

	// 進化レポート出力
	b.Cleanup(func() {
		printEvolutionReport(evolutionResults)
	})
}

// printEvolutionReport はPugコンパイラ進化レポートを出力
func printEvolutionReport(results map[string]map[string]*ComparisonResult) {
	fmt.Print("\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("🚀 Pugコンパイラ進化分析 (vs GCC)\n")
	fmt.Println(strings.Repeat("=", 80))

	for phase, phaseResults := range results {
		fmt.Printf("\n📈 %s:\n", strings.ToUpper(phase))
		fmt.Printf("%-10s %-15s %-15s %-10s\n", "GCC Opt", "実行時間比", "コンパイル比", "状態")
		fmt.Println(strings.Repeat("-", 50))

		for _, optLevel := range []string{"-O0", "-O1", "-O2", "-O3"} {
			if result, exists := phaseResults[optLevel]; exists {
				status := "❌"
				runtimeRatio := "N/A"
				compileRatio := "N/A"

				if result.PugSuccess && result.GCCSuccess {
					status = "✅"
					runtimeRatio = fmt.Sprintf("%.2fx", result.RuntimeSpeedRatio)
					compileRatio = fmt.Sprintf("%.2fx", result.CompileSpeedRatio)
				}

				fmt.Printf("%-10s %-15s %-15s %-10s\n",
					optLevel, runtimeRatio, compileRatio, status)
			}
		}
	}

	fmt.Printf("\n🎯 Phase間性能目標:\n")
	fmt.Printf("  Phase 1 → Phase 2: 10x性能向上\n")
	fmt.Printf("  Phase 2 → Phase 3: 5x性能向上\n")
	fmt.Printf("  Phase 3 → Phase 4: 2x性能向上\n")
	fmt.Println(strings.Repeat("=", 80))
}
