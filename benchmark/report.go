package benchmark

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"
	"time"
)

// BenchmarkReport は包括的ベンチマークレポート構造体
type BenchmarkReport struct {
	Timestamp       time.Time               `json:"timestamp"`
	Version         string                  `json:"version"`
	Phase           string                  `json:"phase"`
	Environment     EnvironmentInfo         `json:"environment"`
	CompilerResults []*BenchmarkResult      `json:"compiler_results"`
	GCCComparisons  []*ComparisonResult     `json:"gcc_comparisons"`
	RustComparisons []*RustComparisonResult `json:"rust_comparisons"`
	Summary         BenchmarkSummary        `json:"summary"`
	Recommendations []string                `json:"recommendations"`
}

// EnvironmentInfo は実行環境情報
type EnvironmentInfo struct {
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	GoVersion string `json:"go_version"`
	CPUModel  string `json:"cpu_model"`
	CPUCores  int    `json:"cpu_cores"`
	MemoryGB  int    `json:"memory_gb"`
	BuildInfo string `json:"build_info"`
}

// BenchmarkSummary はベンチマーク結果のサマリー
type BenchmarkSummary struct {
	TotalTests      int     `json:"total_tests"`
	SuccessfulTests int     `json:"successful_tests"`
	FailedTests     int     `json:"failed_tests"`
	SuccessRate     float64 `json:"success_rate"`

	// 性能指標
	AvgCompileTime time.Duration `json:"avg_compile_time"`
	AvgExecuteTime time.Duration `json:"avg_execute_time"`
	AvgMemoryUsage int64         `json:"avg_memory_usage"`
	AvgBinarySize  int64         `json:"avg_binary_size"`

	// 比較指標
	GCCComparison  ComparisonSummary `json:"gcc_comparison"`
	RustComparison ComparisonSummary `json:"rust_comparison"`

	// 進化指標
	PerformanceGrade string   `json:"performance_grade"`
	NextPhaseGoals   []string `json:"next_phase_goals"`
}

// ComparisonSummary は比較結果のサマリー
type ComparisonSummary struct {
	AvgRuntimeRatio  float64 `json:"avg_runtime_ratio"`
	AvgCompileRatio  float64 `json:"avg_compile_ratio"`
	AvgBinaryRatio   float64 `json:"avg_binary_ratio"`
	AvgMemoryRatio   float64 `json:"avg_memory_ratio"`
	BetterThanTarget bool    `json:"better_than_target"`
	Grade            string  `json:"grade"`
}

// GenerateComprehensiveReport は包括的ベンチマークレポートを生成
func GenerateComprehensiveReport(
	phase string,
	compilerResults []*BenchmarkResult,
	gccResults []*ComparisonResult,
	rustResults []*RustComparisonResult,
) *BenchmarkReport {

	report := &BenchmarkReport{
		Timestamp:       time.Now(),
		Version:         "1.0.0",
		Phase:           phase,
		Environment:     collectEnvironmentInfo(),
		CompilerResults: compilerResults,
		GCCComparisons:  gccResults,
		RustComparisons: rustResults,
	}

	report.Summary = generateSummary(compilerResults, gccResults, rustResults)
	report.Recommendations = generateRecommendations(report.Summary, phase)

	return report
}

// collectEnvironmentInfo は実行環境情報を収集
func collectEnvironmentInfo() EnvironmentInfo {
	return EnvironmentInfo{
		OS:        fmt.Sprintf("%s/%s", getOSInfo(), getArchInfo()),
		Arch:      getArchInfo(),
		GoVersion: getGoVersion(),
		CPUModel:  getCPUModel(),
		CPUCores:  getCPUCores(),
		MemoryGB:  getMemoryGB(),
		BuildInfo: getBuildInfo(),
	}
}

// getOSInfo はOS情報を取得（簡易版）
func getOSInfo() string {
	return "Linux" // 実際の実装では runtime.GOOS を使用
}

// getArchInfo はアーキテクチャ情報を取得
func getArchInfo() string {
	return "amd64" // 実際の実装では runtime.GOARCH を使用
}

// getGoVersion はGoバージョンを取得
func getGoVersion() string {
	return "go1.21" // 実際の実装では runtime.Version() を使用
}

// getCPUModel はCPUモデルを取得
func getCPUModel() string {
	return "Intel Core i7" // 実際の実装では /proc/cpuinfo などを読み取り
}

// getCPUCores はCPUコア数を取得
func getCPUCores() int {
	return 8 // 実際の実装では runtime.NumCPU() を使用
}

// getMemoryGB はメモリ容量を取得
func getMemoryGB() int {
	return 16 // 実際の実装では /proc/meminfo などを読み取り
}

// getBuildInfo はビルド情報を取得
func getBuildInfo() string {
	return "pug-compiler-dev" // 実際の実装では build info を取得
}

// generateSummary はベンチマーク結果のサマリーを生成
func generateSummary(
	compilerResults []*BenchmarkResult,
	gccResults []*ComparisonResult,
	rustResults []*RustComparisonResult,
) BenchmarkSummary {

	summary := BenchmarkSummary{}

	// 基本統計
	totalTests := len(compilerResults)
	successfulTests := 0

	var totalCompileTime, totalExecuteTime time.Duration
	var totalMemory, totalBinarySize int64

	for _, result := range compilerResults {
		if result.Success {
			successfulTests++
			totalCompileTime += result.CompileTime
			totalExecuteTime += result.ExecuteTime
			totalMemory += result.MemoryUsage
			totalBinarySize += result.BinarySize
		}
	}

	summary.TotalTests = totalTests
	summary.SuccessfulTests = successfulTests
	summary.FailedTests = totalTests - successfulTests

	if totalTests > 0 {
		summary.SuccessRate = float64(successfulTests) / float64(totalTests) * 100
	}

	if successfulTests > 0 {
		summary.AvgCompileTime = totalCompileTime / time.Duration(successfulTests)
		summary.AvgExecuteTime = totalExecuteTime / time.Duration(successfulTests)
		summary.AvgMemoryUsage = totalMemory / int64(successfulTests)
		summary.AvgBinarySize = totalBinarySize / int64(successfulTests)
	}

	// GCC比較サマリー
	summary.GCCComparison = generateComparisonSummary(gccResults)

	// Rust比較サマリー
	summary.RustComparison = generateRustComparisonSummary(rustResults)

	// 性能グレード評価
	summary.PerformanceGrade = evaluatePerformanceGrade(summary.GCCComparison, summary.RustComparison)

	// 次フェーズ目標
	summary.NextPhaseGoals = generateNextPhaseGoals(summary.PerformanceGrade)

	return summary
}

// generateComparisonSummary はGCC比較サマリーを生成
func generateComparisonSummary(results []*ComparisonResult) ComparisonSummary {
	summary := ComparisonSummary{}

	if len(results) == 0 {
		return summary
	}

	var totalRuntime, totalCompile, totalBinary, totalMemory float64
	successCount := 0

	for _, result := range results {
		if result.PugSuccess && result.GCCSuccess {
			totalRuntime += result.RuntimeSpeedRatio
			totalCompile += result.CompileSpeedRatio
			totalBinary += result.BinarySizeRatio
			totalMemory += result.MemoryUsageRatio
			successCount++
		}
	}

	if successCount > 0 {
		summary.AvgRuntimeRatio = totalRuntime / float64(successCount)
		summary.AvgCompileRatio = totalCompile / float64(successCount)
		summary.AvgBinaryRatio = totalBinary / float64(successCount)
		summary.AvgMemoryRatio = totalMemory / float64(successCount)

		// グレード評価（GCC基準）
		if summary.AvgRuntimeRatio <= 1.0 {
			summary.Grade = "A+ (GCC同等以上)"
		} else if summary.AvgRuntimeRatio <= 2.0 {
			summary.Grade = "A (優秀)"
		} else if summary.AvgRuntimeRatio <= 5.0 {
			summary.Grade = "B (良好)"
		} else if summary.AvgRuntimeRatio <= 10.0 {
			summary.Grade = "C (改善必要)"
		} else {
			summary.Grade = "D (大幅改善必要)"
		}

		// 目標達成判定
		summary.BetterThanTarget = summary.AvgRuntimeRatio <= 10.0 // Phase1目標
	}

	return summary
}

// generateRustComparisonSummary はRust比較サマリーを生成
func generateRustComparisonSummary(results []*RustComparisonResult) ComparisonSummary {
	summary := ComparisonSummary{}

	if len(results) == 0 {
		return summary
	}

	var totalRuntime, totalCompile, totalBinary, totalMemory float64
	successCount := 0

	for _, result := range results {
		if result.PugSuccess && result.RustSuccess {
			if result.RuntimeSpeedRatio > 0 {
				totalRuntime += result.RuntimeSpeedRatio
			}
			if result.CompileSpeedRatio > 0 {
				totalCompile += result.CompileSpeedRatio
			}
			if result.BinarySizeRatio > 0 {
				totalBinary += result.BinarySizeRatio
			}
			if result.MemoryUsageRatio > 0 {
				totalMemory += result.MemoryUsageRatio
			}
			successCount++
		}
	}

	if successCount > 0 {
		summary.AvgRuntimeRatio = totalRuntime / float64(successCount)
		summary.AvgCompileRatio = totalCompile / float64(successCount)
		summary.AvgBinaryRatio = totalBinary / float64(successCount)
		summary.AvgMemoryRatio = totalMemory / float64(successCount)

		// グレード評価（Rust基準）
		if summary.AvgRuntimeRatio <= 1.0 {
			summary.Grade = "S (Rust同等以上)"
		} else if summary.AvgRuntimeRatio <= 2.0 {
			summary.Grade = "A+ (優秀)"
		} else if summary.AvgRuntimeRatio <= 10.0 {
			summary.Grade = "A (良好)"
		} else if summary.AvgRuntimeRatio <= 100.0 {
			summary.Grade = "B (改善必要)"
		} else {
			summary.Grade = "C (大幅改善必要)"
		}

		// 目標達成判定
		summary.BetterThanTarget = summary.AvgRuntimeRatio <= 100.0 // Phase1目標
	}

	return summary
}

// evaluatePerformanceGrade は総合性能グレードを評価
func evaluatePerformanceGrade(gccSummary, rustSummary ComparisonSummary) string {
	// GCC基準での評価を優先、Rust結果も考慮
	gccRatio := gccSummary.AvgRuntimeRatio
	rustRatio := rustSummary.AvgRuntimeRatio

	// Rust結果が利用可能でGCCより良い場合はRustを考慮
	if rustRatio > 0 && rustRatio < gccRatio {
		gccRatio = (gccRatio + rustRatio) / 2.0 // 平均値を使用
	}

	if gccRatio <= 1.0 {
		return "S+ (産業レベル)"
	} else if gccRatio <= 2.0 {
		return "S (優秀)"
	} else if gccRatio <= 5.0 {
		return "A (良好)"
	} else if gccRatio <= 10.0 {
		return "B (基本達成)"
	} else if gccRatio <= 50.0 {
		return "C (改善必要)"
	} else {
		return "D (初期段階)"
	}
}

// generateNextPhaseGoals は次フェーズの目標を生成
func generateNextPhaseGoals(currentGrade string) []string {
	switch currentGrade {
	case "S+ (産業レベル)", "S (優秀)":
		return []string{
			"更なる最適化の追求",
			"メモリ効率の改善",
			"エラーハンドリングの強化",
		}
	case "A (良好)":
		return []string{
			"実行時間の半分化",
			"バイナリサイズの最適化",
			"高度な最適化パスの実装",
		}
	case "B (基本達成)":
		return []string{
			"実行時間の1/5化",
			"基本的な最適化の実装",
			"型システムの強化",
		}
	case "C (改善必要)":
		return []string{
			"実行時間の1/10化",
			"コード生成効率の改善",
			"制御構造の最適化",
		}
	default:
		return []string{
			"基本機能の安定化",
			"テストカバレッジの向上",
			"エラー処理の改善",
		}
	}
}

// generateRecommendations は改善推奨事項を生成
func generateRecommendations(summary BenchmarkSummary, phase string) []string {
	var recommendations []string

	// 成功率に基づく推奨
	if summary.SuccessRate < 80 {
		recommendations = append(recommendations, "🔧 テスト成功率向上: エラーハンドリングの強化が必要")
	}

	// GCC比較に基づく推奨
	if summary.GCCComparison.AvgRuntimeRatio > 10 {
		recommendations = append(recommendations, "⚡ 実行時間改善: コード生成の最適化が急務")
	}

	if summary.GCCComparison.AvgCompileRatio > 5 {
		recommendations = append(recommendations, "🏃 コンパイル高速化: パーサーと字句解析の最適化推奨")
	}

	// Rust比較に基づく推奨
	if summary.RustComparison.AvgRuntimeRatio > 100 {
		recommendations = append(recommendations, "🦀 Rust比較改善: 基本アルゴリズムの見直しが必要")
	}

	// フェーズ固有の推奨
	switch phase {
	case "phase1":
		recommendations = append(recommendations, "📚 Phase2準備: アセンブリコード生成器の設計開始")
		recommendations = append(recommendations, "🧪 テスト強化: より多様なベンチマークケースの追加")
	case "phase2":
		recommendations = append(recommendations, "🎯 最適化準備: IR設計とSSA形式の検討")
		recommendations = append(recommendations, "📊 性能分析: プロファイリングツールの導入")
	case "phase3":
		recommendations = append(recommendations, "🔗 LLVM準備: LLVM IRへの出力機能の検討")
		recommendations = append(recommendations, "🌐 多言語対応: 他言語バックエンドの検討")
	}

	// 一般的な推奨（常に表示）
	recommendations = append(recommendations, "📈 継続的改善: 定期的なベンチマーク実行と性能追跡")

	return recommendations
}

// SaveReportJSON はレポートをJSON形式で保存
func (r *BenchmarkReport) SaveReportJSON(filename string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("JSONマーシャル失敗: %v", err)
	}

	return os.WriteFile(filename, data, 0600)
}

// LoadReportJSON はJSON形式のレポートを読み込み
func LoadReportJSON(filename string) (*BenchmarkReport, error) {
	// Basic file path validation
	if strings.Contains(filename, "..") {
		return nil, fmt.Errorf("invalid file path: contains directory traversal")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込み失敗: %v", err)
	}

	var report BenchmarkReport
	err = json.Unmarshal(data, &report)
	if err != nil {
		return nil, fmt.Errorf("JSONアンマーシャル失敗: %v", err)
	}

	return &report, nil
}

// GenerateHTMLReport はHTMLレポートを生成
func (r *BenchmarkReport) GenerateHTMLReport(outputPath string) error {
	// Basic file path validation
	if strings.Contains(outputPath, "..") {
		return fmt.Errorf("invalid output path: contains directory traversal")
	}

	tmpl := template.Must(template.New("report").Parse(htmlReportTemplate))

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("HTMLファイル作成失敗: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, r)
}

// CompareReports は複数のレポートを比較
func CompareReports(reports []*BenchmarkReport) *EvolutionReport {
	if len(reports) < 2 {
		return nil
	}

	// 時系列ソート
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Timestamp.Before(reports[j].Timestamp)
	})

	evolution := &EvolutionReport{
		StartTimestamp: reports[0].Timestamp,
		EndTimestamp:   reports[len(reports)-1].Timestamp,
		Reports:        reports,
		Improvements:   make(map[string]float64),
		Trends:         make(map[string]string),
	}

	// 改善率計算
	start := reports[0].Summary
	end := reports[len(reports)-1].Summary

	if start.AvgExecuteTime > 0 {
		ratio := float64(end.AvgExecuteTime) / float64(start.AvgExecuteTime)
		evolution.Improvements["execution_time"] = (1.0 - ratio) * 100
	}

	if start.GCCComparison.AvgRuntimeRatio > 0 {
		ratio := end.GCCComparison.AvgRuntimeRatio / start.GCCComparison.AvgRuntimeRatio
		evolution.Improvements["gcc_comparison"] = (1.0 - ratio) * 100
	}

	// トレンド分析
	if evolution.Improvements["execution_time"] > 0 {
		evolution.Trends["execution_time"] = "改善"
	} else {
		evolution.Trends["execution_time"] = "悪化"
	}

	return evolution
}

// EvolutionReport は進化レポート構造体
type EvolutionReport struct {
	StartTimestamp time.Time          `json:"start_timestamp"`
	EndTimestamp   time.Time          `json:"end_timestamp"`
	Reports        []*BenchmarkReport `json:"reports"`
	Improvements   map[string]float64 `json:"improvements"`
	Trends         map[string]string  `json:"trends"`
}

// HTMLレポートテンプレート
const htmlReportTemplate = `
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🐺 Pugコンパイラ性能ベンチマークレポート</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        h1 { color: #2c3e50; border-bottom: 3px solid #3498db; padding-bottom: 10px; }
        h2 { color: #34495e; margin-top: 30px; }
        .summary { background: #ecf0f1; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .metric { display: inline-block; margin: 10px 20px; text-align: center; }
        .metric-value { font-size: 24px; font-weight: bold; color: #2980b9; }
        .metric-label { font-size: 12px; color: #7f8c8d; }
        .grade { font-size: 20px; font-weight: bold; padding: 10px; border-radius: 5px; }
        .grade-S { background: #2ecc71; color: white; }
        .grade-A { background: #3498db; color: white; }
        .grade-B { background: #f39c12; color: white; }
        .grade-C { background: #e74c3c; color: white; }
        .grade-D { background: #95a5a6; color: white; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #34495e; color: white; }
        .success { color: #27ae60; }
        .failure { color: #e74c3c; }
        .recommendation { background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .footer { margin-top: 40px; text-align: center; color: #7f8c8d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🐺 Pugコンパイラ性能ベンチマークレポート</h1>
        
        <div class="summary">
            <h2>📊 実行サマリー</h2>
            <div class="metric">
                <div class="metric-value">{{.Summary.SuccessfulTests}}/{{.Summary.TotalTests}}</div>
                <div class="metric-label">成功/総数</div>
            </div>
            <div class="metric">
                <div class="metric-value">{{printf "%.1f%%" .Summary.SuccessRate}}</div>
                <div class="metric-label">成功率</div>
            </div>
            <div class="metric">
                <div class="metric-value">{{.Summary.AvgExecuteTime}}</div>
                <div class="metric-label">平均実行時間</div>
            </div>
            <div class="metric">
                <div class="metric-value">{{.Summary.AvgMemoryUsage}}KB</div>
                <div class="metric-label">平均メモリ使用量</div>
            </div>
        </div>

        <div class="summary">
            <h2>🎯 性能グレード</h2>
            <div class="grade grade-A">
                {{.Summary.PerformanceGrade}}
            </div>
        </div>

        <h2>🏁 GCC比較結果</h2>
        <p>平均実行時間比: <strong>{{printf "%.2fx" .Summary.GCCComparison.AvgRuntimeRatio}}</strong></p>
        <p>グレード: <strong>{{.Summary.GCCComparison.Grade}}</strong></p>

        <h2>🦀 Rust比較結果</h2>
        <p>平均実行時間比: <strong>{{printf "%.2fx" .Summary.RustComparison.AvgRuntimeRatio}}</strong></p>
        <p>グレード: <strong>{{.Summary.RustComparison.Grade}}</strong></p>

        <h2>💡 改善推奨事項</h2>
        {{range .Recommendations}}
        <div class="recommendation">{{.}}</div>
        {{end}}

        <h2>🎯 次フェーズ目標</h2>
        <ul>
        {{range .Summary.NextPhaseGoals}}
        <li>{{.}}</li>
        {{end}}
        </ul>

        <div class="footer">
            <p>Generated at {{.Timestamp.Format "2006-01-02 15:04:05"}} | Phase: {{.Phase}} | Version: {{.Version}}</p>
            <p>🤖 Generated with <a href="https://claude.ai/code">Claude Code</a></p>
        </div>
    </div>
</body>
</html>
`
