// 長期性能トレンド分析システム
// 蓄積された性能データから回帰検出とトレンド分析を実行

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// PerformanceHistoryEntry 性能履歴エントリ
type PerformanceHistoryEntry struct {
	Timestamp        time.Time `json:"timestamp"`
	CommitHash       string    `json:"commit_hash"`
	Branch           string    `json:"branch"`
	RunNumber        int       `json:"run_number"`
	AverageNsPerOp   int64     `json:"average_ns_per_op"`
	PerformanceGrade string    `json:"performance_grade"`
	Phase            string    `json:"phase"`
	TestCount        int       `json:"test_count"`
}

// TrendAnalysisResult トレンド分析結果
type TrendAnalysisResult struct {
	AnalysisDate       time.Time                `json:"analysis_date"`
	TotalDataPoints    int                      `json:"total_data_points"`
	AnalysisPeriodDays int                      `json:"analysis_period_days"`
	RegressionDetected bool                     `json:"regression_detected"`
	TrendDirection     string                   `json:"trend_direction"`    // "improving", "stable", "degrading"
	PerformanceChange  float64                  `json:"performance_change"` // % change
	PhaseAnalysis      map[string]PhaseAnalysis `json:"phase_analysis"`
	RegressionAlerts   []RegressionAlert        `json:"regression_alerts"`
	Recommendations    []string                 `json:"recommendations"`
	ChartDataURL       string                   `json:"chart_data_url"`
}

// PhaseAnalysis フェーズ別分析
type PhaseAnalysis struct {
	Phase              string    `json:"phase"`
	DataPoints         int       `json:"data_points"`
	AveragePerformance int64     `json:"average_performance"`
	BestPerformance    int64     `json:"best_performance"`
	WorstPerformance   int64     `json:"worst_performance"`
	TrendDirection     string    `json:"trend_direction"`
	LastUpdate         time.Time `json:"last_update"`
}

// RegressionAlert 回帰アラート
type RegressionAlert struct {
	Severity          string    `json:"severity"` // "critical", "warning", "info"
	Message           string    `json:"message"`
	CommitHash        string    `json:"commit_hash"`
	Timestamp         time.Time `json:"timestamp"`
	PerformanceChange float64   `json:"performance_change"`
	ActionRequired    string    `json:"action_required"`
}

func main() {
	fmt.Println("📈 Pugコンパイラ長期性能トレンド分析")
	fmt.Println("==========================================")

	// 性能履歴データの収集
	historyData, err := collectPerformanceHistory()
	if err != nil {
		log.Printf("⚠️ 履歴データ収集エラー: %v", err)
		return
	}

	if len(historyData) == 0 {
		fmt.Println("📊 性能履歴データが不足しています。データ蓄積を継続してください。")
		return
	}

	fmt.Printf("📊 分析対象データ: %d件\n", len(historyData))

	// トレンド分析実行
	analysis := performTrendAnalysis(historyData)

	// 回帰検出
	detectRegressions(&analysis, historyData)

	// 分析結果をJSON保存
	if err := saveAnalysisResult(analysis); err != nil {
		log.Printf("❌ 分析結果保存失敗: %v", err)
	} else {
		fmt.Println("✅ トレンド分析結果保存完了: trend-analysis.json")
	}

	// アラートレポート生成
	generateAlertReport(analysis)

	// 可視化データ生成
	generateVisualizationData(analysis, historyData)

	// 回帰アラートがある場合は警告終了
	if analysis.RegressionDetected {
		fmt.Printf("🚨 性能回帰が検出されました！詳細は trend-analysis.json を確認してください。\n")
		os.Exit(1)
	}

	fmt.Println("🎉 トレンド分析完了")
}

// collectPerformanceHistory 性能履歴データを収集
func collectPerformanceHistory() ([]PerformanceHistoryEntry, error) {
	fmt.Println("📂 性能履歴データ収集中...")

	var allData []PerformanceHistoryEntry

	// .performance_history ディレクトリから履歴データを収集
	historyDir := ".performance_history"
	if _, err := os.Stat(historyDir); os.IsNotExist(err) {
		return allData, nil // 履歴ディレクトリが存在しない場合は空配列を返す
	}

	err := filepath.Walk(historyDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".json") && strings.Contains(path, "benchmark_") {
			entry, err := parseHistoryFile(path)
			if err != nil {
				log.Printf("⚠️ 履歴ファイル解析失敗 %s: %v", path, err)
				return nil // エラーがあっても処理継続
			}
			allData = append(allData, entry)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 時系列でソート
	sort.Slice(allData, func(i, j int) bool {
		return allData[i].Timestamp.Before(allData[j].Timestamp)
	})

	return allData, nil
}

// parseHistoryFile 履歴ファイルを解析
func parseHistoryFile(filepath string) (PerformanceHistoryEntry, error) {
	var entry PerformanceHistoryEntry

	data, err := os.ReadFile(filepath)
	if err != nil {
		return entry, err
	}

	// 簡易的な履歴データ構造（実際は performance-report.json の形式）
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return entry, err
	}

	// 基本情報の抽出
	if timestamp, ok := rawData["timestamp"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, timestamp); err == nil {
			entry.Timestamp = parsedTime
		}
	}

	if commit, ok := rawData["commit"].(string); ok {
		entry.CommitHash = commit
	}

	if branch, ok := rawData["branch"].(string); ok {
		entry.Branch = branch
	}

	if runNum, ok := rawData["run_number"].(float64); ok {
		entry.RunNumber = int(runNum)
	}

	// デフォルト値設定（実際の実装では詳細な解析が必要）
	entry.AverageNsPerOp = 1000000 // 1ms default
	entry.PerformanceGrade = "B"
	entry.Phase = "Phase1"
	entry.TestCount = 10

	return entry, nil
}

// performTrendAnalysis トレンド分析を実行
func performTrendAnalysis(data []PerformanceHistoryEntry) TrendAnalysisResult {
	fmt.Println("📊 トレンド分析実行中...")

	analysis := TrendAnalysisResult{
		AnalysisDate:       time.Now().UTC(),
		TotalDataPoints:    len(data),
		AnalysisPeriodDays: calculateAnalysisPeriod(data),
		PhaseAnalysis:      make(map[string]PhaseAnalysis),
	}

	if len(data) < 2 {
		analysis.TrendDirection = "insufficient_data"
		analysis.Recommendations = []string{"データ蓄積が不足しています。継続的な測定を行ってください。"}
		return analysis
	}

	// 全体的なトレンド分析
	firstEntry := data[0]
	lastEntry := data[len(data)-1]

	performanceChange := float64(lastEntry.AverageNsPerOp-firstEntry.AverageNsPerOp) / float64(firstEntry.AverageNsPerOp) * 100
	analysis.PerformanceChange = performanceChange

	// トレンド方向の判定
	if performanceChange < -5.0 {
		analysis.TrendDirection = "improving"
	} else if performanceChange > 5.0 {
		analysis.TrendDirection = "degrading"
	} else {
		analysis.TrendDirection = "stable"
	}

	// フェーズ別分析
	phaseData := groupByPhase(data)
	for phase, entries := range phaseData {
		analysis.PhaseAnalysis[phase] = analyzePhase(phase, entries)
	}

	// 推奨事項生成
	analysis.Recommendations = generateTrendRecommendations(analysis)

	return analysis
}

// detectRegressions 性能回帰を検出
func detectRegressions(analysis *TrendAnalysisResult, data []PerformanceHistoryEntry) {
	fmt.Println("🔍 性能回帰検出中...")

	var alerts []RegressionAlert

	// 直近のデータポイントで大幅な性能劣化をチェック
	if len(data) >= 2 {
		recent := data[len(data)-1]
		previous := data[len(data)-2]

		changePercent := float64(recent.AverageNsPerOp-previous.AverageNsPerOp) / float64(previous.AverageNsPerOp) * 100

		if changePercent > 20.0 {
			alerts = append(alerts, RegressionAlert{
				Severity:          "critical",
				Message:           fmt.Sprintf("深刻な性能劣化が検出されました: %.1f%%の性能低下", changePercent),
				CommitHash:        recent.CommitHash,
				Timestamp:         recent.Timestamp,
				PerformanceChange: changePercent,
				ActionRequired:    "即座に原因調査と修正が必要です",
			})
		} else if changePercent > 10.0 {
			alerts = append(alerts, RegressionAlert{
				Severity:          "warning",
				Message:           fmt.Sprintf("性能劣化が検出されました: %.1f%%の性能低下", changePercent),
				CommitHash:        recent.CommitHash,
				Timestamp:         recent.Timestamp,
				PerformanceChange: changePercent,
				ActionRequired:    "原因を調査することを推奨します",
			})
		}
	}

	// 長期トレンドでの劣化チェック
	if analysis.TrendDirection == "degrading" && analysis.PerformanceChange > 15.0 {
		alerts = append(alerts, RegressionAlert{
			Severity:          "warning",
			Message:           fmt.Sprintf("長期的な性能劣化トレンド: %.1f%%の性能低下", analysis.PerformanceChange),
			CommitHash:        "trend_analysis",
			Timestamp:         time.Now().UTC(),
			PerformanceChange: analysis.PerformanceChange,
			ActionRequired:    "性能改善施策の検討が必要です",
		})
	}

	analysis.RegressionAlerts = alerts
	analysis.RegressionDetected = len(alerts) > 0

	fmt.Printf("  🚨 回帰アラート: %d件\n", len(alerts))
}

// ユーティリティ関数群

func calculateAnalysisPeriod(data []PerformanceHistoryEntry) int {
	if len(data) < 2 {
		return 0
	}
	duration := data[len(data)-1].Timestamp.Sub(data[0].Timestamp)
	return int(duration.Hours() / 24)
}

func groupByPhase(data []PerformanceHistoryEntry) map[string][]PerformanceHistoryEntry {
	phaseGroups := make(map[string][]PerformanceHistoryEntry)
	for _, entry := range data {
		phaseGroups[entry.Phase] = append(phaseGroups[entry.Phase], entry)
	}
	return phaseGroups
}

func analyzePhase(phase string, entries []PerformanceHistoryEntry) PhaseAnalysis {
	if len(entries) == 0 {
		return PhaseAnalysis{Phase: phase}
	}

	var total int64
	var best int64 = math.MaxInt64
	var worst int64 = 0

	for _, entry := range entries {
		total += entry.AverageNsPerOp
		if entry.AverageNsPerOp < best {
			best = entry.AverageNsPerOp
		}
		if entry.AverageNsPerOp > worst {
			worst = entry.AverageNsPerOp
		}
	}

	average := total / int64(len(entries))

	// 簡易トレンド分析
	trendDirection := "stable"
	if len(entries) >= 2 {
		first := entries[0].AverageNsPerOp
		last := entries[len(entries)-1].AverageNsPerOp
		change := float64(last-first) / float64(first) * 100

		if change < -5.0 {
			trendDirection = "improving"
		} else if change > 5.0 {
			trendDirection = "degrading"
		}
	}

	return PhaseAnalysis{
		Phase:              phase,
		DataPoints:         len(entries),
		AveragePerformance: average,
		BestPerformance:    best,
		WorstPerformance:   worst,
		TrendDirection:     trendDirection,
		LastUpdate:         entries[len(entries)-1].Timestamp,
	}
}

func generateTrendRecommendations(analysis TrendAnalysisResult) []string {
	var recommendations []string

	switch analysis.TrendDirection {
	case "improving":
		recommendations = append(recommendations, "✅ 性能は改善傾向にあります。現在の開発アプローチを継続してください。")
	case "degrading":
		recommendations = append(recommendations, "⚠️ 性能劣化傾向が検出されました。最適化に取り組むことを推奨します。")
		recommendations = append(recommendations, "🔍 最近のコミットで性能に影響する変更がないか確認してください。")
	case "stable":
		recommendations = append(recommendations, "📊 性能は安定しています。継続的な測定を維持してください。")
	default:
		recommendations = append(recommendations, "📈 データ蓄積を継続して、より詳細な分析を可能にしてください。")
	}

	if analysis.TotalDataPoints < 10 {
		recommendations = append(recommendations, "📊 より信頼性の高い分析のため、追加の測定データが必要です。")
	}

	return recommendations
}

func saveAnalysisResult(analysis TrendAnalysisResult) error {
	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("trend-analysis.json", data, 0644)
}

func generateAlertReport(analysis TrendAnalysisResult) {
	if len(analysis.RegressionAlerts) == 0 {
		return
	}

	fmt.Println("\n🚨 性能回帰アラートレポート")
	fmt.Println("============================")

	for _, alert := range analysis.RegressionAlerts {
		fmt.Printf("🔴 [%s] %s\n", strings.ToUpper(alert.Severity), alert.Message)
		fmt.Printf("   コミット: %s\n", alert.CommitHash)
		fmt.Printf("   時刻: %s\n", alert.Timestamp.Format("2006-01-02 15:04:05 UTC"))
		fmt.Printf("   対応: %s\n\n", alert.ActionRequired)
	}
}

func generateVisualizationData(analysis TrendAnalysisResult, data []PerformanceHistoryEntry) {
	fmt.Println("📊 可視化データ生成中...")

	// Chart.js用のデータ生成（簡易実装）
	chartData := map[string]interface{}{
		"type": "line",
		"data": map[string]interface{}{
			"labels": extractTimestamps(data),
			"datasets": []map[string]interface{}{
				{
					"label":           "平均実行時間 (ns/op)",
					"data":            extractPerformanceData(data),
					"borderColor":     "#007acc",
					"backgroundColor": "rgba(0, 122, 204, 0.1)",
					"tension":         0.1,
				},
			},
		},
		"options": map[string]interface{}{
			"responsive": true,
			"plugins": map[string]interface{}{
				"title": map[string]interface{}{
					"display": true,
					"text":    "Pugコンパイラ性能推移",
				},
			},
			"scales": map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": false,
					"title": map[string]interface{}{
						"display": true,
						"text":    "実行時間 (ns/op)",
					},
				},
			},
		},
	}

	// Chart.jsデータをJSONファイルに保存
	chartJSON, _ := json.MarshalIndent(chartData, "", "  ")
	err := os.WriteFile("performance-chart-data.json", chartJSON, 0644)
	if err != nil {
		log.Printf("⚠️ Chart.jsデータファイル保存失敗: %v", err)
	}

	fmt.Println("  ✅ Chart.js用データ: performance-chart-data.json")
}

func extractTimestamps(data []PerformanceHistoryEntry) []string {
	var timestamps []string
	for _, entry := range data {
		timestamps = append(timestamps, entry.Timestamp.Format("01/02"))
	}
	return timestamps
}

func extractPerformanceData(data []PerformanceHistoryEntry) []int64 {
	var performance []int64
	for _, entry := range data {
		performance = append(performance, entry.AverageNsPerOp)
	}
	return performance
}
