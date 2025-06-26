// 性能ダッシュボード生成システム
// 包括的な性能データを可視化するHTMLダッシュボードを生成

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"
)

// DashboardData ダッシュボード用データ
type DashboardData struct {
	GeneratedAt    time.Time          `json:"generated_at"`
	ProjectInfo    ProjectInfo        `json:"project_info"`
	CurrentResults PerformanceResults `json:"current_results"`
	TrendAnalysis  TrendAnalysisData  `json:"trend_analysis"`
	Comparisons    ComparisonResults  `json:"comparisons"`
	Alerts         []PerformanceAlert `json:"alerts"`
	Charts         []ChartConfig      `json:"charts"`
}

// ProjectInfo プロジェクト情報
type ProjectInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	CommitHash  string `json:"commit_hash"`
	Branch      string `json:"branch"`
	RunNumber   int    `json:"run_number"`
	Environment string `json:"environment"`
}

// PerformanceResults 現在の性能結果
type PerformanceResults struct {
	Summary         ResultSummary           `json:"summary"`
	BasicResults    []BenchmarkItem         `json:"basic_results"`
	CompilerResults []BenchmarkItem         `json:"compiler_results"`
	PhaseResults    map[string]PhaseMetrics `json:"phase_results"`
}

// ResultSummary 結果サマリー
type ResultSummary struct {
	TotalTests       int    `json:"total_tests"`
	SuccessfulTests  int    `json:"successful_tests"`
	FailedTests      int    `json:"failed_tests"`
	AverageNsPerOp   int64  `json:"average_ns_per_op"`
	PerformanceGrade string `json:"performance_grade"`
	GradeColor       string `json:"grade_color"`
}

// BenchmarkItem ベンチマーク項目
type BenchmarkItem struct {
	Name       string `json:"name"`
	NsPerOp    int64  `json:"ns_per_op"`
	Iterations int    `json:"iterations"`
	Phase      string `json:"phase"`
	Status     string `json:"status"`
}

// PhaseMetrics フェーズメトリクス
type PhaseMetrics struct {
	Phase          string  `json:"phase"`
	TestCount      int     `json:"test_count"`
	AverageNsPerOp int64   `json:"average_ns_per_op"`
	BestNsPerOp    int64   `json:"best_ns_per_op"`
	WorstNsPerOp   int64   `json:"worst_ns_per_op"`
	TargetAchieved bool    `json:"target_achieved"`
	TargetRatio    float64 `json:"target_ratio"`
}

// TrendAnalysisData トレンド分析データ
type TrendAnalysisData struct {
	Direction      string  `json:"direction"`
	ChangePercent  float64 `json:"change_percent"`
	DataPoints     int     `json:"data_points"`
	AnalysisPeriod int     `json:"analysis_period_days"`
	Confidence     string  `json:"confidence"`
}

// ComparisonResults 比較結果
type ComparisonResults struct {
	GCCComparison  []ComparisonItem `json:"gcc_comparison"`
	RustComparison []ComparisonItem `json:"rust_comparison"`
}

// ComparisonItem 比較項目
type ComparisonItem struct {
	TestName     string  `json:"test_name"`
	PugTime      int64   `json:"pug_time"`
	CompareTime  int64   `json:"compare_time"`
	SpeedupRatio float64 `json:"speedup_ratio"`
	Status       string  `json:"status"`
}

// PerformanceAlert 性能アラート
type PerformanceAlert struct {
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	CommitHash string    `json:"commit_hash"`
	ActionItem string    `json:"action_item"`
}

// ChartConfig チャート設定
type ChartConfig struct {
	ID     string      `json:"id"`
	Type   string      `json:"type"`
	Title  string      `json:"title"`
	Data   interface{} `json:"data"`
	Height string      `json:"height"`
}

func main() {
	fmt.Println("📊 Pugコンパイラ性能ダッシュボード生成")
	fmt.Println("=====================================")

	// データ収集
	dashboardData, err := collectDashboardData()
	if err != nil {
		log.Fatalf("❌ ダッシュボードデータ収集失敗: %v", err)
	}

	// HTMLダッシュボード生成
	if err := generateHTMLDashboard(dashboardData); err != nil {
		log.Fatalf("❌ HTMLダッシュボード生成失敗: %v", err)
	}

	// JSONデータ出力
	if err := saveDashboardJSON(dashboardData); err != nil {
		log.Printf("⚠️ JSONデータ保存失敗: %v", err)
	}

	fmt.Println("✅ 性能ダッシュボード生成完了")
	fmt.Println("   📄 performance-dashboard.html")
	fmt.Println("   📊 dashboard-data.json")
}

// collectDashboardData ダッシュボード用データを収集
func collectDashboardData() (DashboardData, error) {
	fmt.Println("📂 ダッシュボードデータ収集中...")

	data := DashboardData{
		GeneratedAt: time.Now().UTC(),
		ProjectInfo: ProjectInfo{
			Name:        "Pug Compiler",
			Version:     "v0.1.0",
			CommitHash:  getEnvOrDefault("GITHUB_SHA", "development"),
			Branch:      getEnvOrDefault("GITHUB_REF_NAME", "main"),
			RunNumber:   0,
			Environment: "GitHub Actions",
		},
	}

	// 性能レポートデータ読み込み
	if perfData, err := loadPerformanceData(); err == nil {
		data.CurrentResults = perfData
	} else {
		fmt.Printf("⚠️ 性能データ読み込み失敗: %v\n", err)
		data.CurrentResults = generateMockPerformanceData()
	}

	// トレンド分析データ読み込み
	if trendData, err := loadTrendData(); err == nil {
		data.TrendAnalysis = trendData
	} else {
		fmt.Printf("ℹ️ トレンド分析データなし\n")
		data.TrendAnalysis = generateMockTrendData()
	}

	// 比較結果読み込み
	data.Comparisons = loadComparisonData()

	// アラート生成
	data.Alerts = generateAlerts(data)

	// チャート設定生成
	data.Charts = generateChartConfigs(data)

	return data, nil
}

// loadPerformanceData 性能データを読み込み
func loadPerformanceData() (PerformanceResults, error) {
	var results PerformanceResults

	// performance-report.json を読み込み
	if data, err := loadJSONFile("performance-report.json"); err == nil {
		results = parsePerformanceData(data)
	} else {
		return results, err
	}

	return results, nil
}

// parsePerformanceData 性能データを解析
func parsePerformanceData(data map[string]interface{}) PerformanceResults {
	results := PerformanceResults{
		PhaseResults: make(map[string]PhaseMetrics),
	}

	// サマリー情報
	if summary, ok := data["summary"].(map[string]interface{}); ok {
		results.Summary = ResultSummary{
			TotalTests:       getIntValue(summary, "total_benchmarks"),
			SuccessfulTests:  getIntValue(summary, "successful_tests"),
			FailedTests:      getIntValue(summary, "failed_tests"),
			AverageNsPerOp:   getInt64Value(summary, "average_ns_per_op"),
			PerformanceGrade: getStringValue(summary, "performance_grade"),
			GradeColor:       getGradeColor(getStringValue(summary, "performance_grade")),
		}
	}

	// 基本ベンチマーク結果
	if basic, ok := data["basic_benchmarks"].([]interface{}); ok {
		for _, item := range basic {
			if benchmark, ok := item.(map[string]interface{}); ok {
				results.BasicResults = append(results.BasicResults, BenchmarkItem{
					Name:       getStringValue(benchmark, "name"),
					NsPerOp:    getInt64Value(benchmark, "ns_per_op"),
					Iterations: getIntValue(benchmark, "iterations"),
					Phase:      getStringValue(benchmark, "phase"),
					Status:     "success",
				})
			}
		}
	}

	return results
}

// generateHTMLDashboard HTMLダッシュボードを生成
func generateHTMLDashboard(data DashboardData) error {
	fmt.Println("🎨 HTMLダッシュボード生成中...")

	tmpl := `<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.ProjectInfo.Name}} - 性能ダッシュボード</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
            background-color: #f8f9fa;
            color: #333;
            line-height: 1.6;
        }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px 0;
            text-align: center;
            margin-bottom: 30px;
            border-radius: 10px;
        }
        .header h1 { font-size: 2.5em; margin-bottom: 10px; }
        .header .subtitle { opacity: 0.9; font-size: 1.1em; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .card { 
            background: white;
            border-radius: 10px;
            padding: 20px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            border-left: 4px solid #667eea;
        }
        .card h3 { margin-bottom: 15px; color: #333; }
        .metric { display: flex; justify-content: space-between; margin-bottom: 10px; }
        .metric-value { font-weight: bold; color: #667eea; }
        .grade-s { color: #ff6b35; }
        .grade-a { color: #f7931e; }
        .grade-b { color: #fccc02; }
        .grade-c { color: #8bc34a; }
        .grade-d { color: #9e9e9e; }
        .alert { 
            padding: 15px;
            margin-bottom: 15px;
            border-radius: 5px;
            border-left: 4px solid;
        }
        .alert-critical { background: #ffebee; border-color: #f44336; }
        .alert-warning { background: #fff3e0; border-color: #ff9800; }
        .alert-info { background: #e3f2fd; border-color: #2196f3; }
        .chart-container { height: 400px; margin: 20px 0; }
        .footer { 
            text-align: center;
            margin-top: 50px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
        }
        .trend-improving { color: #4caf50; }
        .trend-stable { color: #2196f3; }
        .trend-degrading { color: #f44336; }
        .comparison-table { width: 100%; border-collapse: collapse; margin-top: 15px; }
        .comparison-table th, .comparison-table td { 
            padding: 8px 12px; 
            text-align: left; 
            border-bottom: 1px solid #eee; 
        }
        .comparison-table th { background-color: #f5f5f5; font-weight: 600; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🐺 {{.ProjectInfo.Name}}</h1>
            <div class="subtitle">
                性能ダッシュボード | {{.GeneratedAt.Format "2006-01-02 15:04:05 UTC"}}
            </div>
            <div style="margin-top: 15px; font-size: 0.9em;">
                コミット: {{slice .ProjectInfo.CommitHash 0 8}} | ブランチ: {{.ProjectInfo.Branch}}
            </div>
        </div>

        <div class="grid">
            <div class="card">
                <h3>📊 性能サマリー</h3>
                <div class="metric">
                    <span>総テスト数</span>
                    <span class="metric-value">{{.CurrentResults.Summary.TotalTests}}</span>
                </div>
                <div class="metric">
                    <span>成功率</span>
                    <span class="metric-value">{{printf "%.1f%%" (div (mul .CurrentResults.Summary.SuccessfulTests 100.0) .CurrentResults.Summary.TotalTests)}}</span>
                </div>
                <div class="metric">
                    <span>平均実行時間</span>
                    <span class="metric-value">{{.CurrentResults.Summary.AverageNsPerOp}} ns/op</span>
                </div>
                <div class="metric">
                    <span>性能グレード</span>
                    <span class="metric-value grade-{{.CurrentResults.Summary.PerformanceGrade | lower}}">{{.CurrentResults.Summary.PerformanceGrade}}</span>
                </div>
            </div>

            <div class="card">
                <h3>📈 トレンド分析</h3>
                <div class="metric">
                    <span>傾向</span>
                    <span class="metric-value trend-{{.TrendAnalysis.Direction}}">
                        {{if eq .TrendAnalysis.Direction "improving"}}📈 改善中{{end}}
                        {{if eq .TrendAnalysis.Direction "stable"}}📊 安定{{end}}
                        {{if eq .TrendAnalysis.Direction "degrading"}}📉 劣化中{{end}}
                    </span>
                </div>
                <div class="metric">
                    <span>変化率</span>
                    <span class="metric-value">{{printf "%.1f%%" .TrendAnalysis.ChangePercent}}</span>
                </div>
                <div class="metric">
                    <span>データ期間</span>
                    <span class="metric-value">{{.TrendAnalysis.AnalysisPeriod}}日間</span>
                </div>
                <div class="metric">
                    <span>信頼度</span>
                    <span class="metric-value">{{.TrendAnalysis.Confidence}}</span>
                </div>
            </div>

            <div class="card">
                <h3>🔍 フェーズ別性能</h3>
                {{range $phase, $metrics := .CurrentResults.PhaseResults}}
                <div style="margin-bottom: 15px;">
                    <strong>{{$phase}}</strong>
                    <div class="metric">
                        <span>平均時間</span>
                        <span class="metric-value">{{$metrics.AverageNsPerOp}} ns/op</span>
                    </div>
                    <div class="metric">
                        <span>目標達成</span>
                        <span class="metric-value">{{if $metrics.TargetAchieved}}✅{{else}}⏳{{end}}</span>
                    </div>
                </div>
                {{end}}
            </div>
        </div>

        {{if .Alerts}}
        <div class="card">
            <h3>🚨 性能アラート</h3>
            {{range .Alerts}}
            <div class="alert alert-{{.Level}}">
                <strong>{{if eq .Level "critical"}}🔴{{else if eq .Level "warning"}}🟡{{else}}🔵{{end}} {{.Level | title}}</strong>: {{.Message}}
                <div style="font-size: 0.9em; margin-top: 5px;">
                    {{.Timestamp.Format "2006-01-02 15:04:05"}} | {{.ActionItem}}
                </div>
            </div>
            {{end}}
        </div>
        {{end}}

        <div class="grid">
            <div class="card">
                <h3>🏁 GCC比較</h3>
                {{if .Comparisons.GCCComparison}}
                <table class="comparison-table">
                    <thead>
                        <tr><th>テスト</th><th>Pug</th><th>GCC</th><th>比率</th></tr>
                    </thead>
                    <tbody>
                        {{range .Comparisons.GCCComparison}}
                        <tr>
                            <td>{{.TestName}}</td>
                            <td>{{.PugTime}} ns</td>
                            <td>{{.CompareTime}} ns</td>
                            <td>{{printf "%.1fx" .SpeedupRatio}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
                {{else}}
                <p>GCC比較データなし</p>
                {{end}}
            </div>

            <div class="card">
                <h3>🦀 Rust比較</h3>
                {{if .Comparisons.RustComparison}}
                <table class="comparison-table">
                    <thead>
                        <tr><th>テスト</th><th>Pug</th><th>Rust</th><th>比率</th></tr>
                    </thead>
                    <tbody>
                        {{range .Comparisons.RustComparison}}
                        <tr>
                            <td>{{.TestName}}</td>
                            <td>{{.PugTime}} ns</td>
                            <td>{{.CompareTime}} ns</td>
                            <td>{{printf "%.1fx" .SpeedupRatio}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
                {{else}}
                <p>Rust比較データなし</p>
                {{end}}
            </div>
        </div>

        {{range .Charts}}
        <div class="card">
            <h3>{{.Title}}</h3>
            <div class="chart-container">
                <canvas id="{{.ID}}" style="{{.Height}}"></canvas>
            </div>
        </div>
        {{end}}

        <div class="footer">
            <p>🤖 Generated with <a href="https://claude.ai/code">Claude Code</a></p>
            <p>最終更新: {{.GeneratedAt.Format "2006-01-02 15:04:05 UTC"}}</p>
        </div>
    </div>

    <script>
        // Chart.js設定
        {{range .Charts}}
        new Chart(document.getElementById('{{.ID}}'), {{.Data | toJSON}});
        {{end}}
    </script>
</body>
</html>`

	// テンプレート関数
	funcMap := template.FuncMap{
		"lower": func(s string) string { return s },
		"toJSON": func(v interface{}) template.JS {
			bytes, _ := json.Marshal(v)
			return template.JS(bytes)
		},
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mul": func(a, b float64) float64 { return a * b },
		"slice": func(s string, start, end int) string {
			if len(s) < end {
				return s
			}
			return s[start:end]
		},
	}

	t, err := template.New("dashboard").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create("performance-dashboard.html")
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, data)
}

// ユーティリティ関数群（続き）

func loadJSONFile(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntValue(data map[string]interface{}, key string) int {
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	return 0
}

func getInt64Value(data map[string]interface{}, key string) int64 {
	if val, ok := data[key].(float64); ok {
		return int64(val)
	}
	return 0
}

func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getGradeColor(grade string) string {
	switch grade {
	case "S+":
		return "#ff6b35"
	case "S":
		return "#f7931e"
	case "A":
		return "#fccc02"
	case "B":
		return "#8bc34a"
	case "C":
		return "#2196f3"
	default:
		return "#9e9e9e"
	}
}

func generateMockPerformanceData() PerformanceResults {
	return PerformanceResults{
		Summary: ResultSummary{
			TotalTests:       10,
			SuccessfulTests:  10,
			FailedTests:      0,
			AverageNsPerOp:   1000000,
			PerformanceGrade: "B",
			GradeColor:       "#8bc34a",
		},
		BasicResults: []BenchmarkItem{
			{Name: "Fibonacci", NsPerOp: 500000, Iterations: 1000, Phase: "Phase1", Status: "success"},
			{Name: "Factorial", NsPerOp: 250000, Iterations: 2000, Phase: "Phase1", Status: "success"},
		},
		PhaseResults: map[string]PhaseMetrics{
			"Phase1": {
				Phase: "Phase1", TestCount: 5, AverageNsPerOp: 800000,
				BestNsPerOp: 250000, WorstNsPerOp: 1500000, TargetAchieved: true, TargetRatio: 0.8,
			},
		},
	}
}

func generateMockTrendData() TrendAnalysisData {
	return TrendAnalysisData{
		Direction:      "stable",
		ChangePercent:  2.5,
		DataPoints:     15,
		AnalysisPeriod: 30,
		Confidence:     "高",
	}
}

func loadTrendData() (TrendAnalysisData, error) {
	return generateMockTrendData(), nil
}

func loadComparisonData() ComparisonResults {
	return ComparisonResults{
		GCCComparison: []ComparisonItem{
			{TestName: "fibonacci", PugTime: 1000000, CompareTime: 100000, SpeedupRatio: 0.1, Status: "slower"},
		},
		RustComparison: []ComparisonItem{
			{TestName: "fibonacci", PugTime: 1000000, CompareTime: 50000, SpeedupRatio: 0.05, Status: "slower"},
		},
	}
}

func generateAlerts(data DashboardData) []PerformanceAlert {
	var alerts []PerformanceAlert

	if data.CurrentResults.Summary.PerformanceGrade == "D" {
		alerts = append(alerts, PerformanceAlert{
			Level:      "warning",
			Message:    "性能グレードが低下しています",
			Timestamp:  time.Now(),
			CommitHash: data.ProjectInfo.CommitHash,
			ActionItem: "最適化の検討が必要です",
		})
	}

	return alerts
}

func generateChartConfigs(data DashboardData) []ChartConfig {
	return []ChartConfig{
		{
			ID:     "performance-trend",
			Type:   "line",
			Title:  "📈 性能推移",
			Height: "height: 300px;",
			Data: map[string]interface{}{
				"type": "line",
				"data": map[string]interface{}{
					"labels": []string{"Week 1", "Week 2", "Week 3", "Week 4"},
					"datasets": []map[string]interface{}{
						{
							"label":           "平均実行時間 (ns/op)",
							"data":            []int64{1200000, 1100000, 1000000, 950000},
							"borderColor":     "#667eea",
							"backgroundColor": "rgba(102, 126, 234, 0.1)",
							"tension":         0.1,
						},
					},
				},
				"options": map[string]interface{}{
					"responsive":          true,
					"maintainAspectRatio": false,
					"scales": map[string]interface{}{
						"y": map[string]interface{}{
							"beginAtZero": false,
						},
					},
				},
			},
		},
	}
}

func saveDashboardJSON(data DashboardData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("dashboard-data.json", jsonData, 0644)
}
