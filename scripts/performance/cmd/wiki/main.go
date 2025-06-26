// GitHub Wiki自動更新システム
// 性能測定結果を自動的にGitHub Wikiページに反映

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// WikiUpdater GitHub Wiki更新システム
type WikiUpdater struct {
	RepoURL  string
	WikiURL  string
	TempDir  string
	GitUser  string
	GitEmail string
}

// WikiPage Wiki更新ページ情報
type WikiPage struct {
	Filename string
	Title    string
	Content  string
}

// PerformanceWikiData Wiki更新用性能データ
type PerformanceWikiData struct {
	Timestamp       time.Time              `json:"timestamp"`
	CommitHash      string                 `json:"commit_hash"`
	Branch          string                 `json:"branch"`
	PerformanceData map[string]interface{} `json:"performance_data"`
	TrendAnalysis   map[string]interface{} `json:"trend_analysis"`
}

func main() {
	fmt.Println("📝 GitHub Wiki自動更新システム")
	fmt.Println("=================================")

	// 環境変数から設定取得
	repoURL := getEnvOrDefault("GITHUB_REPOSITORY", "nyasuto/pug")
	commitHash := getEnvOrDefault("GITHUB_SHA", "unknown")
	branch := getEnvOrDefault("GITHUB_REF_NAME", "main")

	updater := WikiUpdater{
		RepoURL:  fmt.Sprintf("https://github.com/%s", repoURL),
		WikiURL:  fmt.Sprintf("https://github.com/%s.wiki.git", repoURL),
		TempDir:  "/tmp/pug-wiki-update",
		GitUser:  "github-actions[bot]",
		GitEmail: "github-actions[bot]@users.noreply.github.com",
	}

	fmt.Printf("📊 Wiki更新対象: %s\n", updater.RepoURL)

	// 性能データ収集
	wikiData, err := collectWikiData(commitHash, branch)
	if err != nil {
		log.Printf("⚠️ Wiki更新データ収集失敗: %v", err)
		return
	}

	// Wiki更新実行
	if err := updater.UpdatePerformanceWiki(wikiData); err != nil {
		log.Printf("❌ Wiki更新失敗: %v", err)
		os.Exit(1)
	}

	fmt.Println("✅ GitHub Wiki更新完了")
}

// collectWikiData Wiki更新用データを収集
func collectWikiData(commitHash, branch string) (PerformanceWikiData, error) {
	data := PerformanceWikiData{
		Timestamp:       time.Now().UTC(),
		CommitHash:      commitHash,
		Branch:          branch,
		PerformanceData: make(map[string]interface{}),
		TrendAnalysis:   make(map[string]interface{}),
	}

	// 性能レポートJSONを読み込み
	if reportData, err := loadJSONFile("performance-report.json"); err == nil {
		data.PerformanceData = reportData
	} else {
		fmt.Printf("⚠️ performance-report.json 読み込み失敗: %v\n", err)
	}

	// トレンド分析JSONを読み込み
	if trendData, err := loadJSONFile("trend-analysis.json"); err == nil {
		data.TrendAnalysis = trendData
	} else {
		fmt.Printf("ℹ️ trend-analysis.json がありません（初回実行時は正常）\n")
	}

	return data, nil
}

// UpdatePerformanceWiki 性能Wiki更新メイン処理
func (w *WikiUpdater) UpdatePerformanceWiki(data PerformanceWikiData) error {
	fmt.Println("🔧 Wiki更新処理開始...")

	// Wikiリポジトリクローン
	if err := w.cloneWikiRepo(); err != nil {
		return fmt.Errorf("wikiクローン失敗: %v", err)
	}
	defer w.cleanup()

	// Git設定
	if err := w.setupGitConfig(); err != nil {
		return fmt.Errorf("git設定失敗: %v", err)
	}

	// Wiki ページ生成・更新
	pages := w.generateWikiPages(data)
	for _, page := range pages {
		if err := w.updateWikiPage(page); err != nil {
			log.Printf("⚠️ ページ更新失敗 %s: %v", page.Title, err)
		} else {
			fmt.Printf("  ✅ %s 更新完了\n", page.Title)
		}
	}

	// コミット・プッシュ
	if err := w.commitAndPush(data); err != nil {
		return fmt.Errorf("コミット・プッシュ失敗: %v", err)
	}

	return nil
}

// cloneWikiRepo Wikiリポジトリをクローン
func (w *WikiUpdater) cloneWikiRepo() error {
	fmt.Println("📥 Wikiリポジトリクローン中...")

	// 既存ディレクトリを削除
	_ = os.RemoveAll(w.TempDir)

	// Wikiクローン（失敗しても継続 - Wikiが初回作成の場合）
	cmd := exec.Command("git", "clone", w.WikiURL, w.TempDir) // #nosec G204 - controlled git operations for Wiki automation
	if err := cmd.Run(); err != nil {
		fmt.Printf("ℹ️ Wikiクローン失敗（初回作成時は正常）: %v\n", err)
		// 空のディレクトリを作成して初期化
		if err := os.MkdirAll(w.TempDir, 0750); err != nil {
			return err
		}

		// Git初期化
		cmd = exec.Command("git", "init")
		cmd.Dir = w.TempDir
		if err := cmd.Run(); err != nil {
			return err
		}

		// リモート追加
		cmd = exec.Command("git", "remote", "add", "origin", w.WikiURL) // #nosec G204 - controlled git operations
		cmd.Dir = w.TempDir
		_ = cmd.Run() // エラーは無視（すでに存在する場合）
	}

	return nil
}

// setupGitConfig Git設定をセットアップ
func (w *WikiUpdater) setupGitConfig() error {
	fmt.Println("⚙️ Git設定中...")

	commands := [][]string{
		{"git", "config", "user.name", w.GitUser},
		{"git", "config", "user.email", w.GitEmail},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) // #nosec G204 - controlled git config commands
		cmd.Dir = w.TempDir
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// generateWikiPages Wiki更新ページを生成
func (w *WikiUpdater) generateWikiPages(data PerformanceWikiData) []WikiPage {
	var pages []WikiPage

	// メイン性能ベンチマークページ
	pages = append(pages, WikiPage{
		Filename: "Performance-Benchmark.md",
		Title:    "性能ベンチマーク",
		Content:  w.generateMainBenchmarkPage(data),
	})

	// 進化履歴ページ
	pages = append(pages, WikiPage{
		Filename: "Performance-Evolution.md",
		Title:    "性能進化履歴",
		Content:  w.generateEvolutionPage(data),
	})

	// GCC比較ページ
	pages = append(pages, WikiPage{
		Filename: "GCC-Comparison.md",
		Title:    "GCC比較",
		Content:  w.generateGCCComparisonPage(data),
	})

	// Rust比較ページ
	pages = append(pages, WikiPage{
		Filename: "Rust-Comparison.md",
		Title:    "Rust比較",
		Content:  w.generateRustComparisonPage(data),
	})

	// トレンド分析ページ
	if len(data.TrendAnalysis) > 0 {
		pages = append(pages, WikiPage{
			Filename: "Performance-Trends.md",
			Title:    "性能トレンド分析",
			Content:  w.generateTrendAnalysisPage(data),
		})
	}

	return pages
}

// generateMainBenchmarkPage メインベンチマークページを生成
func (w *WikiUpdater) generateMainBenchmarkPage(data PerformanceWikiData) string {
	content := fmt.Sprintf(`# 📊 Pugコンパイラ性能ベンチマーク

**最終更新**: %s  
**コミット**: [%s](%s/commit/%s)  
**ブランチ**: %s

## 概要

このページには、Pugコンパイラプロジェクトの最新性能測定結果が自動更新されます。

## 🚀 最新ベンチマーク結果

`, data.Timestamp.Format("2006-01-02 15:04:05 UTC"),
		data.CommitHash[:8], w.RepoURL, data.CommitHash, data.Branch)

	// 性能データがある場合は詳細を追加
	if perfData, ok := data.PerformanceData["summary"].(map[string]interface{}); ok {
		content += fmt.Sprintf(`### 📈 性能サマリー

- **総ベンチマーク数**: %v
- **成功テスト数**: %v
- **平均実行時間**: %v ns/op
- **性能グレード**: %v

`,
			getValueOrDefault(perfData, "total_benchmarks", "N/A"),
			getValueOrDefault(perfData, "successful_tests", "N/A"),
			getValueOrDefault(perfData, "average_ns_per_op", "N/A"),
			getValueOrDefault(perfData, "performance_grade", "N/A"))
	}

	content += `## 📊 詳細レポート

- [性能進化履歴](Performance-Evolution) - フェーズ間の性能変遷
- [GCC比較](GCC-Comparison) - 産業標準Cコンパイラとの比較
- [Rust比較](Rust-Comparison) - 現代的システム言語との比較
- [性能トレンド分析](Performance-Trends) - 長期的な性能動向

## 🎯 フェーズ別目標

### Phase 1: 基本言語処理（インタープリター）
- ✅ 基本機能の安定動作
- ✅ 包括的テストカバレッジ（75%+）
- 🎯 GCCの10-100倍以内の実行時間

### Phase 2: コンパイラ基盤（アセンブリ生成）
- 🎯 Phase1から10倍性能向上
- 🎯 GCCの2-10倍以内の実行時間
- 🎯 基本的なコード最適化

### Phase 3: 最適化エンジン（IR最適化）
- 🎯 Phase2から5倍性能向上
- 🎯 GCCの1-2倍以内の実行時間

### Phase 4: 産業レベル（LLVM統合）
- 🎯 Phase3から2倍性能向上
- 🎯 GCC同等の実行時間

---

🤖 このページは自動生成されています。[Claude Code](https://claude.ai/code)により更新。
`

	return content
}

// generateEvolutionPage 進化履歴ページを生成
func (w *WikiUpdater) generateEvolutionPage(data PerformanceWikiData) string {
	return fmt.Sprintf(`# 📈 Pugコンパイラ性能進化履歴

**最終更新**: %s

## 📊 進化チャート

TODO: 過去の性能データから進化チャートを生成

## 🎯 マイルストーン達成状況

### Phase 1 → Phase 2
- **目標**: 10倍性能向上
- **現状**: 測定中...

### Phase 2 → Phase 3  
- **目標**: 5倍性能向上
- **現状**: 未実装

### Phase 3 → Phase 4
- **目標**: 2倍性能向上  
- **現状**: 未実装

---

🤖 自動更新: %s
`, data.Timestamp.Format("2006-01-02 15:04:05 UTC"), data.Timestamp.Format("2006-01-02 15:04:05 UTC"))
}

// generateGCCComparisonPage GCC比較ページを生成
func (w *WikiUpdater) generateGCCComparisonPage(data PerformanceWikiData) string {
	return fmt.Sprintf(`# 🏁 GCC比較ベンチマーク

**最終更新**: %s

## 概要

産業標準CコンパイラGCCとの性能比較結果です。

## 📊 比較結果

TODO: GCC比較データの詳細表示

---

🤖 自動更新: %s
`, data.Timestamp.Format("2006-01-02 15:04:05 UTC"), data.Timestamp.Format("2006-01-02 15:04:05 UTC"))
}

// generateRustComparisonPage Rust比較ページを生成
func (w *WikiUpdater) generateRustComparisonPage(data PerformanceWikiData) string {
	return fmt.Sprintf(`# 🦀 Rust比較ベンチマーク

**最終更新**: %s

## 概要

現代的システム言語Rustとの性能比較結果です。

## 📊 比較結果

TODO: Rust比較データの詳細表示

---

🤖 自動更新: %s
`, data.Timestamp.Format("2006-01-02 15:04:05 UTC"), data.Timestamp.Format("2006-01-02 15:04:05 UTC"))
}

// generateTrendAnalysisPage トレンド分析ページを生成
func (w *WikiUpdater) generateTrendAnalysisPage(data PerformanceWikiData) string {
	content := fmt.Sprintf(`# 📈 性能トレンド分析

**最終更新**: %s

## 概要

長期的な性能動向と回帰検出結果です。

`, data.Timestamp.Format("2006-01-02 15:04:05 UTC"))

	// トレンド分析データがある場合は詳細を追加
	if trendData, ok := data.TrendAnalysis["trend_direction"].(string); ok {
		var emoji string
		switch trendData {
		case "improving":
			emoji = "📈"
		case "degrading":
			emoji = "📉"
		default:
			emoji = "📊"
		}

		content += fmt.Sprintf(`## %s 現在のトレンド

**方向**: %s

`, emoji, trendData)
	}

	content += `## 🚨 回帰アラート

TODO: 回帰アラート情報の表示

---

🤖 自動更新: ` + data.Timestamp.Format("2006-01-02 15:04:05 UTC")

	return content
}

// updateWikiPage Wikiページを更新
func (w *WikiUpdater) updateWikiPage(page WikiPage) error {
	filePath := filepath.Join(w.TempDir, page.Filename)
	return os.WriteFile(filePath, []byte(page.Content), 0600)
}

// commitAndPush 変更をコミットしてプッシュ
func (w *WikiUpdater) commitAndPush(data PerformanceWikiData) error {
	fmt.Println("💾 Wikiコミット・プッシュ中...")

	// ステージング
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = w.TempDir
	if err := cmd.Run(); err != nil {
		return err
	}

	// 変更があるかチェック
	cmd = exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = w.TempDir
	if err := cmd.Run(); err == nil {
		fmt.Println("ℹ️ Wiki更新なし（変更なし）")
		return nil
	}

	// コミット
	commitMsg := fmt.Sprintf(`性能ベンチマーク自動更新

- コミット: %s
- ブランチ: %s
- 更新日時: %s

🤖 Generated with Claude Code
`, data.CommitHash[:8], data.Branch, data.Timestamp.Format("2006-01-02 15:04:05 UTC"))

	cmd = exec.Command("git", "commit", "-m", commitMsg) // #nosec G204 - controlled git commit with validated message
	cmd.Dir = w.TempDir
	if err := cmd.Run(); err != nil {
		return err
	}

	// プッシュ
	cmd = exec.Command("git", "push", "origin", "master")
	cmd.Dir = w.TempDir
	if err := cmd.Run(); err != nil {
		// masterがない場合はmainを試す
		cmd = exec.Command("git", "push", "origin", "main")
		cmd.Dir = w.TempDir
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	fmt.Println("✅ Wiki更新完了")
	return nil
}

// cleanup 一時ディレクトリをクリーンアップ
func (w *WikiUpdater) cleanup() {
	_ = os.RemoveAll(w.TempDir)
}

// ユーティリティ関数群

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadJSONFile(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filename) // #nosec G304 - controlled performance data file reading
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func getValueOrDefault(data map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, ok := data[key]; ok {
		return value
	}
	return defaultValue
}
