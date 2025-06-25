package benchmark

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestBenchmarkSetup はベンチマーク環境セットアップをテスト
func TestBenchmarkSetup(t *testing.T) {
	testProgram := `
let fib = fn(n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
};

let result = fib(5);
puts(result);
`

	cb, err := setupBenchmark("phase1", testProgram)
	if err != nil {
		t.Fatalf("ベンチマーク環境セットアップ失敗: %v", err)
	}
	defer cb.cleanup()

	// 一時ディレクトリとファイルの存在確認
	if _, err := os.Stat(cb.TempDir); os.IsNotExist(err) {
		t.Errorf("一時ディレクトリが作成されていません: %s", cb.TempDir)
	}

	if _, err := os.Stat(cb.SourceFile); os.IsNotExist(err) {
		t.Errorf("ソースファイルが作成されていません: %s", cb.SourceFile)
	}

	// ソースファイルの内容確認
	content, err := os.ReadFile(cb.SourceFile)
	if err != nil {
		t.Fatalf("ソースファイル読み込み失敗: %v", err)
	}

	if string(content) != testProgram {
		t.Errorf("ソースファイルの内容が異なります\n期待: %s\n実際: %s", testProgram, string(content))
	}
}

// TestBenchmarkCleanup はベンチマーク環境クリーンアップをテスト
func TestBenchmarkCleanup(t *testing.T) {
	cb, err := setupBenchmark("phase1", "let x = 5;")
	if err != nil {
		t.Fatalf("ベンチマーク環境セットアップ失敗: %v", err)
	}

	tempDir := cb.TempDir

	// クリーンアップ前の存在確認
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatalf("一時ディレクトリが存在しません: %s", tempDir)
	}

	// クリーンアップ実行
	cb.cleanup()

	// クリーンアップ後の確認
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Errorf("一時ディレクトリが削除されていません: %s", tempDir)
	}
}

// TestGCCComparisonSetup はGCC比較環境セットアップをテスト
func TestGCCComparisonSetup(t *testing.T) {
	benchmark := GCCBenchmark{
		Name: "test_fib",
		PugSource: `
let fib = fn(n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
};

let result = fib(3);
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
    int result = fib(3);
    printf("%d\n", result);
    return 0;
}
`,
		ExpectedOutput: "2",
	}

	pugBench, gccBench, tempDir, err := setupGCCComparison(benchmark, "-O2")
	if err != nil {
		t.Fatalf("GCC比較環境セットアップ失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Pugファイルの確認
	if _, err := os.Stat(pugBench.SourceFile); os.IsNotExist(err) {
		t.Errorf("Pugソースファイルが作成されていません: %s", pugBench.SourceFile)
	}

	// Cファイルの確認
	if _, err := os.Stat(gccBench.SourceFile); os.IsNotExist(err) {
		t.Errorf("Cソースファイルが作成されていません: %s", gccBench.SourceFile)
	}

	// ファイル内容確認
	pugContent, err := os.ReadFile(pugBench.SourceFile)
	if err != nil {
		t.Fatalf("Pugソースファイル読み込み失敗: %v", err)
	}

	if string(pugContent) != benchmark.PugSource {
		t.Errorf("Pugソースファイルの内容が異なります")
	}

	cContent, err := os.ReadFile(gccBench.SourceFile)
	if err != nil {
		t.Fatalf("Cソースファイル読み込み失敗: %v", err)
	}

	if string(cContent) != benchmark.CSource {
		t.Errorf("Cソースファイルの内容が異なります")
	}
}

// TestRustComparisonSetup はRust比較環境セットアップをテスト
func TestRustComparisonSetup(t *testing.T) {
	benchmark := RustBenchmark{
		Name: "test_fib",
		PugSource: `
let fib = fn(n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
};

let result = fib(3);
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
    let result = fib(3);
    println!("{}", result);
}
`,
		ExpectedOutput: "2",
	}

	pugBench, rustBench, tempDir, err := setupRustComparison(benchmark, "debug")
	if err != nil {
		t.Fatalf("Rust比較環境セットアップ失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Pugファイルの確認
	if _, err := os.Stat(pugBench.SourceFile); os.IsNotExist(err) {
		t.Errorf("Pugソースファイルが作成されていません: %s", pugBench.SourceFile)
	}

	// Rustプロジェクトファイルの確認
	rustProjectDir := filepath.Dir(filepath.Dir(rustBench.SourceFile))
	cargoToml := filepath.Join(rustProjectDir, "Cargo.toml")
	if _, err := os.Stat(cargoToml); os.IsNotExist(err) {
		t.Errorf("Cargo.tomlが作成されていません: %s", cargoToml)
	}

	if _, err := os.Stat(rustBench.SourceFile); os.IsNotExist(err) {
		t.Errorf("Rustソースファイルが作成されていません: %s", rustBench.SourceFile)
	}
}

// TestBenchmarkReport はベンチマークレポート生成をテスト
func TestBenchmarkReport(t *testing.T) {
	// テスト用のダミー結果を作成
	compilerResults := []*BenchmarkResult{
		{
			Phase:         "phase1",
			CompileTime:   100 * time.Millisecond,
			ExecuteTime:   50 * time.Millisecond,
			MemoryUsage:   1024,
			BinarySize:    8192,
			Success:       true,
			ThroughputOps: 20,
		},
		{
			Phase:         "phase1",
			CompileTime:   120 * time.Millisecond,
			ExecuteTime:   60 * time.Millisecond,
			MemoryUsage:   1200,
			BinarySize:    8500,
			Success:       true,
			ThroughputOps: 16,
		},
	}

	gccResults := []*ComparisonResult{
		{
			TestName:          "fibonacci",
			OptLevel:          "-O2",
			PugCompileTime:    100 * time.Millisecond,
			PugExecuteTime:    50 * time.Millisecond,
			PugSuccess:        true,
			GCCCompileTime:    50 * time.Millisecond,
			GCCExecuteTime:    10 * time.Millisecond,
			GCCSuccess:        true,
			CompileSpeedRatio: 2.0,
			RuntimeSpeedRatio: 5.0,
			BinarySizeRatio:   1.5,
			MemoryUsageRatio:  1.2,
		},
	}

	rustResults := []*RustComparisonResult{
		{
			TestName:          "fibonacci",
			OptLevel:          "release",
			PugCompileTime:    100 * time.Millisecond,
			PugExecuteTime:    50 * time.Millisecond,
			PugSuccess:        true,
			RustCompileTime:   200 * time.Millisecond,
			RustExecuteTime:   5 * time.Millisecond,
			RustSuccess:       true,
			CompileSpeedRatio: 0.5,
			RuntimeSpeedRatio: 10.0,
			BinarySizeRatio:   2.0,
			MemoryUsageRatio:  1.5,
		},
	}

	// レポート生成
	report := GenerateComprehensiveReport("phase1", compilerResults, gccResults, rustResults)

	// 基本検証
	if report.Phase != "phase1" {
		t.Errorf("フェーズが正しく設定されていません: %s", report.Phase)
	}

	if report.Summary.TotalTests != 2 {
		t.Errorf("総テスト数が正しくありません: %d", report.Summary.TotalTests)
	}

	if report.Summary.SuccessfulTests != 2 {
		t.Errorf("成功テスト数が正しくありません: %d", report.Summary.SuccessfulTests)
	}

	if report.Summary.SuccessRate != 100.0 {
		t.Errorf("成功率が正しくありません: %.1f", report.Summary.SuccessRate)
	}

	// GCC比較確認
	if report.Summary.GCCComparison.AvgRuntimeRatio != 5.0 {
		t.Errorf("GCC実行時間比が正しくありません: %.1f", report.Summary.GCCComparison.AvgRuntimeRatio)
	}

	// Rust比較確認
	if report.Summary.RustComparison.AvgRuntimeRatio != 10.0 {
		t.Errorf("Rust実行時間比が正しくありません: %.1f", report.Summary.RustComparison.AvgRuntimeRatio)
	}

	// 推奨事項確認
	if len(report.Recommendations) == 0 {
		t.Errorf("推奨事項が生成されていません")
	}
}

// TestReportJSONSerialization はレポートのJSON保存・読み込みをテスト
func TestReportJSONSerialization(t *testing.T) {
	// テスト用レポート作成
	report := GenerateComprehensiveReport("phase1",
		[]*BenchmarkResult{{Phase: "phase1", Success: true}},
		[]*ComparisonResult{},
		[]*RustComparisonResult{})

	// 一時ファイル作成
	tempFile, err := os.CreateTemp("", "test_report_*.json")
	if err != nil {
		t.Fatalf("一時ファイル作成失敗: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// JSON保存
	err = report.SaveReportJSON(tempFile.Name())
	if err != nil {
		t.Fatalf("JSONレポート保存失敗: %v", err)
	}

	// JSON読み込み
	loadedReport, err := LoadReportJSON(tempFile.Name())
	if err != nil {
		t.Fatalf("JSONレポート読み込み失敗: %v", err)
	}

	// 内容確認
	if loadedReport.Phase != report.Phase {
		t.Errorf("読み込んだレポートのフェーズが異なります: %s != %s", loadedReport.Phase, report.Phase)
	}

	if loadedReport.Version != report.Version {
		t.Errorf("読み込んだレポートのバージョンが異なります: %s != %s", loadedReport.Version, report.Version)
	}
}

// TestHTMLReportGeneration はHTMLレポート生成をテスト
func TestHTMLReportGeneration(t *testing.T) {
	// テスト用レポート作成
	report := GenerateComprehensiveReport("phase1",
		[]*BenchmarkResult{{Phase: "phase1", Success: true}},
		[]*ComparisonResult{},
		[]*RustComparisonResult{})

	// 一時ファイル作成
	tempFile, err := os.CreateTemp("", "test_report_*.html")
	if err != nil {
		t.Fatalf("一時ファイル作成失敗: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// HTMLレポート生成
	err = report.GenerateHTMLReport(tempFile.Name())
	if err != nil {
		t.Fatalf("HTMLレポート生成失敗: %v", err)
	}

	// ファイル存在確認
	if _, err := os.Stat(tempFile.Name()); os.IsNotExist(err) {
		t.Errorf("HTMLレポートファイルが作成されていません: %s", tempFile.Name())
	}

	// ファイル内容確認（基本的なHTMLタグの存在）
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("HTMLレポートファイル読み込み失敗: %v", err)
	}

	htmlContent := string(content)

	if !containsString(htmlContent, "<!DOCTYPE html>") {
		t.Errorf("HTMLファイルにDOCTYPE宣言がありません")
	}

	if !containsString(htmlContent, "Pugコンパイラ性能ベンチマークレポート") {
		t.Errorf("HTMLファイルにタイトルがありません")
	}

	if !containsString(htmlContent, "phase1") {
		t.Errorf("HTMLファイルにフェーズ情報がありません")
	}
}

// TestEnvironmentInfoCollection は環境情報収集をテスト
func TestEnvironmentInfoCollection(t *testing.T) {
	env := collectEnvironmentInfo()

	if env.OS == "" {
		t.Errorf("OS情報が取得されていません")
	}

	if env.Arch == "" {
		t.Errorf("アーキテクチャ情報が取得されていません")
	}

	if env.GoVersion == "" {
		t.Errorf("Goバージョン情報が取得されていません")
	}

	if env.CPUCores <= 0 {
		t.Errorf("CPUコア数が正しく取得されていません: %d", env.CPUCores)
	}

	if env.MemoryGB <= 0 {
		t.Errorf("メモリ容量が正しく取得されていません: %d", env.MemoryGB)
	}
}

// containsString は文字列に部分文字列が含まれているかをチェック
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BenchmarkReportGeneration はレポート生成の性能をベンチマーク
func BenchmarkReportGeneration(b *testing.B) {
	// 大量のテスト結果を生成
	compilerResults := make([]*BenchmarkResult, 100)
	for i := 0; i < 100; i++ {
		compilerResults[i] = &BenchmarkResult{
			Phase:         "phase1",
			CompileTime:   time.Duration(i) * time.Millisecond,
			ExecuteTime:   time.Duration(i*2) * time.Millisecond,
			MemoryUsage:   int64(1024 + i*10),
			BinarySize:    int64(8192 + i*100),
			Success:       true,
			ThroughputOps: int64(20 - i/10),
		}
	}

	gccResults := make([]*ComparisonResult, 50)
	for i := 0; i < 50; i++ {
		gccResults[i] = &ComparisonResult{
			TestName:          "test",
			OptLevel:          "-O2",
			PugSuccess:        true,
			GCCSuccess:        true,
			CompileSpeedRatio: float64(i) * 0.1,
			RuntimeSpeedRatio: float64(i) * 0.2,
			BinarySizeRatio:   1.5,
			MemoryUsageRatio:  1.2,
		}
	}

	rustResults := make([]*RustComparisonResult, 50)
	for i := 0; i < 50; i++ {
		rustResults[i] = &RustComparisonResult{
			TestName:          "test",
			OptLevel:          "release",
			PugSuccess:        true,
			RustSuccess:       true,
			CompileSpeedRatio: float64(i) * 0.05,
			RuntimeSpeedRatio: float64(i) * 0.1,
			BinarySizeRatio:   2.0,
			MemoryUsageRatio:  1.5,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GenerateComprehensiveReport("phase1", compilerResults, gccResults, rustResults)
	}
}
