package benchmark

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// CompilerBenchmark はコンパイラ性能測定の構造体
type CompilerBenchmark struct {
	Phase          string
	SourceCode     string
	ExpectedOutput string
	CompileCommand []string
	ExecuteCommand []string
	TempDir        string
	SourceFile     string
	BinaryFile     string
}

// BenchmarkResult はベンチマーク結果の構造体
type BenchmarkResult struct {
	Phase         string
	CompileTime   time.Duration
	ExecuteTime   time.Duration
	MemoryUsage   int64 // KB
	BinarySize    int64 // bytes
	Success       bool
	ErrorMessage  string
	ThroughputOps int64 // operations per second
}

// 共通テストプログラム
var (
	// フィボナッチ計算プログラム（再帰）
	fibonacciProgram = `
let fib = fn(n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
};

let result = fib(20);
puts(result);
`

	// ソートアルゴリズムプログラム
	sortProgram = `
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
    
    return quicksort(less) + [pivot] + quicksort(greater);
};

let numbers = [64, 34, 25, 12, 22, 11, 90, 5, 77, 30];
let sorted = quicksort(numbers);
puts(sorted);
`

	// 数値計算プログラム
	numericalProgram = `
let calculate_pi = fn(iterations) {
    let pi = 0.0;
    let i = 0;
    while (i < iterations) {
        let term = 4.0 / (2.0 * i + 1.0);
        if (i % 2 == 0) {
            pi = pi + term;
        } else {
            pi = pi - term;
        }
        i = i + 1;
    }
    return pi;
};

let result = calculate_pi(1000);
puts(result);
`

	// 複雑な制御構造プログラム
	complexControlProgram = `
let factorial = fn(n) {
    let result = 1;
    let i = 1;
    while (i <= n) {
        result = result * i;
        i = i + 1;
    }
    return result;
};

let sum_factorials = fn(max) {
    let sum = 0;
    let i = 1;
    while (i <= max) {
        sum = sum + factorial(i);
        i = i + 1;
    }
    return sum;
};

let result = sum_factorials(10);
puts(result);
`
)

// setupBenchmark はベンチマーク環境をセットアップ
func setupBenchmark(phase string, sourceCode string) (*CompilerBenchmark, error) {
	tempDir, err := os.MkdirTemp("", "pug_benchmark_"+phase+"_*")
	if err != nil {
		return nil, fmt.Errorf("一時ディレクトリ作成失敗: %v", err)
	}

	sourceFile := filepath.Join(tempDir, "test.pug")
	binaryFile := filepath.Join(tempDir, "test")
	if runtime.GOOS == "windows" {
		binaryFile += ".exe"
	}

	err = os.WriteFile(sourceFile, []byte(sourceCode), 0600)
	if err != nil {
		_ = os.RemoveAll(tempDir) // #nosec G104 - ignore cleanup errors
		return nil, fmt.Errorf("ソースファイル作成失敗: %v", err)
	}

	var compileCommand, executeCommand []string

	switch phase {
	case "phase1":
		// Phase 1: インタープリター
		compileCommand = []string{} // インタープリターはコンパイル不要
		executeCommand = []string{"./bin/interp", sourceFile}
	case "phase2":
		// Phase 2: アセンブリコンパイラ
		compileCommand = []string{"./bin/pugc", sourceFile, "-o", binaryFile}
		executeCommand = []string{binaryFile}
	case "phase3":
		// Phase 3: 最適化コンパイラ（予定）
		compileCommand = []string{"./bin/pugc", "-O2", sourceFile, "-o", binaryFile}
		executeCommand = []string{binaryFile}
	case "phase4":
		// Phase 4: LLVM統合（予定）
		compileCommand = []string{"./bin/pugc", "--backend=llvm", "-O3", sourceFile, "-o", binaryFile}
		executeCommand = []string{binaryFile}
	default:
		_ = os.RemoveAll(tempDir) // #nosec G104 - ignore cleanup errors
		return nil, fmt.Errorf("未サポートフェーズ: %s", phase)
	}

	return &CompilerBenchmark{
		Phase:          phase,
		SourceCode:     sourceCode,
		CompileCommand: compileCommand,
		ExecuteCommand: executeCommand,
		TempDir:        tempDir,
		SourceFile:     sourceFile,
		BinaryFile:     binaryFile,
	}, nil
}

// measureCompileTime はコンパイル時間を測定
func (cb *CompilerBenchmark) measureCompileTime(ctx context.Context) (time.Duration, error) {
	if len(cb.CompileCommand) == 0 {
		return 0, nil // インタープリターはコンパイル時間0
	}

	start := time.Now()
	cmd := exec.CommandContext(ctx, cb.CompileCommand[0], cb.CompileCommand[1:]...) // #nosec G204 - safe for benchmark use
	cmd.Dir = "."

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("コンパイル失敗: %v", err)
	}

	return time.Since(start), nil
}

// measureExecuteTime は実行時間を測定
func (cb *CompilerBenchmark) measureExecuteTime(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	cmd := exec.CommandContext(ctx, cb.ExecuteCommand[0], cb.ExecuteCommand[1:]...) // #nosec G204 - safe for benchmark use
	cmd.Dir = "."

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("実行失敗: %v", err)
	}

	return time.Since(start), nil
}

// measureMemoryUsage はメモリ使用量を測定（概算）
func (cb *CompilerBenchmark) measureMemoryUsage() (int64, error) {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	// Safely convert uint64 to int64 to avoid overflow
	alloc := m.Alloc / 1024
	if alloc > 9223372036854775807 { // MaxInt64
		return 9223372036854775807, nil // Cap at MaxInt64
	}
	return int64(alloc), nil // KB
}

// measureBinarySize はバイナリサイズを測定
func (cb *CompilerBenchmark) measureBinarySize() (int64, error) {
	if len(cb.CompileCommand) == 0 {
		return 0, nil // インタープリターはバイナリ生成なし
	}

	stat, err := os.Stat(cb.BinaryFile)
	if err != nil {
		return 0, fmt.Errorf("バイナリファイル情報取得失敗: %v", err)
	}

	return stat.Size(), nil
}

// cleanup はベンチマーク環境をクリーンアップ
func (cb *CompilerBenchmark) cleanup() {
	if cb.TempDir != "" {
		_ = os.RemoveAll(cb.TempDir) // #nosec G104 - ignore cleanup errors
	}
}

// runBenchmark はベンチマークを実行
func (cb *CompilerBenchmark) runBenchmark(timeout time.Duration) *BenchmarkResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	defer cb.cleanup()

	result := &BenchmarkResult{
		Phase: cb.Phase,
	}

	// コンパイル時間測定
	compileTime, err := cb.measureCompileTime(ctx)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}
	result.CompileTime = compileTime

	// 実行時間測定
	executeTime, err := cb.measureExecuteTime(ctx)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}
	result.ExecuteTime = executeTime

	// メモリ使用量測定
	memUsage, err := cb.measureMemoryUsage()
	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}
	result.MemoryUsage = memUsage

	// バイナリサイズ測定
	binarySize, err := cb.measureBinarySize()
	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}
	result.BinarySize = binarySize

	// スループット計算（操作/秒）
	totalTime := compileTime + executeTime
	if totalTime > 0 {
		result.ThroughputOps = int64(time.Second / totalTime)
	}

	result.Success = true
	return result
}

// BenchmarkCompiler_Phase1_Fibonacci はPhase1フィボナッチベンチマーク
func BenchmarkCompiler_Phase1_Fibonacci(b *testing.B) {
	cb, err := setupBenchmark("phase1", fibonacciProgram)
	if err != nil {
		b.Skip("Phase1環境セットアップ失敗:", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := cb.runBenchmark(30 * time.Second)
		if !result.Success {
			b.Fatalf("ベンチマーク失敗: %s", result.ErrorMessage)
		}
	}
}

// BenchmarkCompiler_Phase1_Sort はPhase1ソートベンチマーク
func BenchmarkCompiler_Phase1_Sort(b *testing.B) {
	cb, err := setupBenchmark("phase1", sortProgram)
	if err != nil {
		b.Skip("Phase1環境セットアップ失敗:", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := cb.runBenchmark(30 * time.Second)
		if !result.Success {
			b.Fatalf("ベンチマーク失敗: %s", result.ErrorMessage)
		}
	}
}

// BenchmarkCompiler_Phase1_Numerical はPhase1数値計算ベンチマーク
func BenchmarkCompiler_Phase1_Numerical(b *testing.B) {
	cb, err := setupBenchmark("phase1", numericalProgram)
	if err != nil {
		b.Skip("Phase1環境セットアップ失敗:", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := cb.runBenchmark(30 * time.Second)
		if !result.Success {
			b.Fatalf("ベンチマーク失敗: %s", result.ErrorMessage)
		}
	}
}

// BenchmarkCompiler_Phase1_ComplexControl はPhase1複雑制御構造ベンチマーク
func BenchmarkCompiler_Phase1_ComplexControl(b *testing.B) {
	cb, err := setupBenchmark("phase1", complexControlProgram)
	if err != nil {
		b.Skip("Phase1環境セットアップ失敗:", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := cb.runBenchmark(30 * time.Second)
		if !result.Success {
			b.Fatalf("ベンチマーク失敗: %s", result.ErrorMessage)
		}
	}
}

// 将来のPhase2-4のベンチマーク関数（現在はスキップ）

// BenchmarkCompiler_Phase2_Fibonacci はPhase2フィボナッチベンチマーク
func BenchmarkCompiler_Phase2_Fibonacci(b *testing.B) {
	if _, err := os.Stat("./bin/pugc"); os.IsNotExist(err) {
		b.Skip("Phase2コンパイラが見つかりません")
		return
	}

	cb, err := setupBenchmark("phase2", fibonacciProgram)
	if err != nil {
		b.Skip("Phase2環境セットアップ失敗:", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := cb.runBenchmark(30 * time.Second)
		if !result.Success {
			b.Fatalf("ベンチマーク失敗: %s", result.ErrorMessage)
		}
	}
}

// Phase3とPhase4は現在未実装のためスキップ
func BenchmarkCompiler_Phase3_Fibonacci(b *testing.B) {
	b.Skip("Phase3は未実装")
}

func BenchmarkCompiler_Phase4_Fibonacci(b *testing.B) {
	b.Skip("Phase4は未実装")
}

// BenchmarkSuite は全フェーズの包括的ベンチマークを実行
func BenchmarkSuite(b *testing.B) {
	phases := []string{"phase1"}
	programs := map[string]string{
		"fibonacci":    fibonacciProgram,
		"sort":         sortProgram,
		"numerical":    numericalProgram,
		"complex_ctrl": complexControlProgram,
	}

	// 利用可能なPhaseを検出
	if _, err := os.Stat("./bin/pugc"); err == nil {
		phases = append(phases, "phase2")
	}

	var results []*BenchmarkResult

	for _, phase := range phases {
		for name, program := range programs {
			b.Run(fmt.Sprintf("%s_%s", phase, name), func(b *testing.B) {
				cb, err := setupBenchmark(phase, program)
				if err != nil {
					b.Skip("環境セットアップ失敗:", err)
					return
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					result := cb.runBenchmark(30 * time.Second)
					if !result.Success {
						b.Fatalf("ベンチマーク失敗: %s", result.ErrorMessage)
					}
					if i == 0 { // 最初の実行結果のみ保存
						results = append(results, result)
					}
				}
			})
		}
	}

	// 結果をレポート出力（ベンチマーク終了後）
	b.Cleanup(func() {
		printBenchmarkReport(results)
	})
}

// printBenchmarkReport はベンチマーク結果のレポートを出力
func printBenchmarkReport(results []*BenchmarkResult) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 コンパイラ性能ベンチマーク結果")
	fmt.Println(strings.Repeat("=", 80))

	for _, result := range results {
		fmt.Printf("\n🔧 フェーズ: %s\n", result.Phase)
		fmt.Printf("⏱️  コンパイル時間: %v\n", result.CompileTime)
		fmt.Printf("🏃 実行時間: %v\n", result.ExecuteTime)
		fmt.Printf("💾 メモリ使用量: %d KB\n", result.MemoryUsage)
		fmt.Printf("📦 バイナリサイズ: %d bytes\n", result.BinarySize)
		fmt.Printf("🚀 スループット: %d ops/sec\n", result.ThroughputOps)
		fmt.Printf("✅ 成功: %t\n", result.Success)
		if result.ErrorMessage != "" {
			fmt.Printf("❌ エラー: %s\n", result.ErrorMessage)
		}
		fmt.Println(strings.Repeat("-", 40))
	}

	fmt.Println("\n📈 Performance Trend:")
	fmt.Println("Phase 1 (Interpreter) → Phase 2 (Basic Compiler) → Phase 3 (Optimizer) → Phase 4 (LLVM)")
	fmt.Println("目標: 10x → 50x → 100x 性能向上")
	fmt.Println(strings.Repeat("=", 80))
}
