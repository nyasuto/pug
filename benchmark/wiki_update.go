package benchmark

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// WikiUpdater はGitHub Wiki自動更新機能
type WikiUpdater struct {
	RepoURL     string
	WikiURL     string
	TempDir     string
	CommitUser  string
	CommitEmail string
}

// validateFilePath validates file paths to prevent directory traversal
func validateFilePath(path string) error {
	// Check for directory traversal attempts
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid file path: contains directory traversal")
	}

	// Check for absolute paths outside allowed directory
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid file path: %v", err)
	}

	// Only allow paths within temp directory or current working directory
	wd, _ := os.Getwd()
	if !strings.HasPrefix(abs, wd) && !strings.HasPrefix(abs, os.TempDir()) {
		return fmt.Errorf("invalid file path: outside allowed directory")
	}

	return nil
}

// validateGitInput validates git command inputs
func validateGitInput(input string) error {
	// Only allow alphanumeric, spaces, dots, hyphens, underscores, and @
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9\s\.\-_@]+$`)
	if !validPattern.MatchString(input) {
		return fmt.Errorf("invalid git input: contains unsafe characters")
	}
	return nil
}

// NewWikiUpdater は新しいWikiUpdaterを作成
func NewWikiUpdater(repoURL string) *WikiUpdater {
	// GitHubリポジトリURLからWiki URLを生成
	wikiURL := strings.ReplaceAll(repoURL, ".git", ".wiki.git")

	return &WikiUpdater{
		RepoURL:     repoURL,
		WikiURL:     wikiURL,
		CommitUser:  "pug-benchmark-bot",
		CommitEmail: "noreply@github.com",
	}
}

// UpdateBenchmarkWiki はベンチマーク結果でWikiを更新
func (wu *WikiUpdater) UpdateBenchmarkWiki(report *BenchmarkReport) error {
	// 一時ディレクトリ作成
	tempDir, err := os.MkdirTemp("", "wiki_update_*")
	if err != nil {
		return fmt.Errorf("一時ディレクトリ作成失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)

	wu.TempDir = tempDir

	// Wikiリポジトリクローン
	err = wu.cloneWikiRepo()
	if err != nil {
		return fmt.Errorf("wikiクローン失敗: %v", err)
	}

	// ベンチマークページ更新
	err = wu.updateBenchmarkPages(report)
	if err != nil {
		return fmt.Errorf("ベンチマークページ更新失敗: %v", err)
	}

	// 変更をコミット・プッシュ
	err = wu.commitAndPush(report)
	if err != nil {
		return fmt.Errorf("コミット・プッシュ失敗: %v", err)
	}

	return nil
}

// cloneWikiRepo はWikiリポジトリをクローン
func (wu *WikiUpdater) cloneWikiRepo() error {
	// Validate inputs
	if err := validateGitInput(wu.WikiURL); err != nil {
		return fmt.Errorf("invalid wiki URL: %v", err)
	}

	wikiPath := filepath.Join(wu.TempDir, "wiki")
	if err := validateFilePath(wikiPath); err != nil {
		return fmt.Errorf("invalid wiki path: %v", err)
	}

	cmd := exec.Command("git", "clone", wu.WikiURL, wikiPath) // #nosec G204 - controlled input for wiki management
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// updateBenchmarkPages はベンチマークページを更新
func (wu *WikiUpdater) updateBenchmarkPages(report *BenchmarkReport) error {
	wikiDir := filepath.Join(wu.TempDir, "wiki")

	// メインベンチマークページ更新
	err := wu.updateMainBenchmarkPage(wikiDir, report)
	if err != nil {
		return err
	}

	// フェーズ別詳細ページ更新
	err = wu.updatePhaseDetailPage(wikiDir, report)
	if err != nil {
		return err
	}

	// 比較結果ページ更新
	err = wu.updateComparisonPages(wikiDir, report)
	if err != nil {
		return err
	}

	// 進化履歴ページ更新
	err = wu.updateEvolutionHistoryPage(wikiDir, report)
	if err != nil {
		return err
	}

	return nil
}

// updateMainBenchmarkPage はメインベンチマークページを更新
func (wu *WikiUpdater) updateMainBenchmarkPage(wikiDir string, report *BenchmarkReport) error {
	content := wu.generateMainBenchmarkContent(report)

	filename := filepath.Join(wikiDir, "Performance-Benchmark.md")
	return os.WriteFile(filename, []byte(content), 0600)
}

// generateMainBenchmarkContent はメインベンチマークページの内容を生成
func (wu *WikiUpdater) generateMainBenchmarkContent(report *BenchmarkReport) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "# 📊 Pugコンパイラ性能ベンチマーク\n\n")
	fmt.Fprintf(&buf, "**最終更新**: %s\n", report.Timestamp.Format("2006-01-02 15:04:05 JST"))
	fmt.Fprintf(&buf, "**現在のフェーズ**: %s\n", strings.ToUpper(report.Phase))
	fmt.Fprintf(&buf, "**バージョン**: %s\n\n", report.Version)

	// 実行環境情報
	fmt.Fprintf(&buf, "## 🔧 実行環境\n\n")
	fmt.Fprintf(&buf, "| 項目 | 値 |\n")
	fmt.Fprintf(&buf, "|------|----|\n")
	fmt.Fprintf(&buf, "| OS | %s |\n", report.Environment.OS)
	fmt.Fprintf(&buf, "| アーキテクチャ | %s |\n", report.Environment.Arch)
	fmt.Fprintf(&buf, "| Go バージョン | %s |\n", report.Environment.GoVersion)
	fmt.Fprintf(&buf, "| CPU | %s |\n", report.Environment.CPUModel)
	fmt.Fprintf(&buf, "| CPU コア数 | %d |\n", report.Environment.CPUCores)
	fmt.Fprintf(&buf, "| メモリ | %d GB |\n\n", report.Environment.MemoryGB)

	// 総合性能グレード
	fmt.Fprintf(&buf, "## 🎯 総合性能グレード\n\n")
	fmt.Fprintf(&buf, "### %s\n\n", report.Summary.PerformanceGrade)

	// 基本統計
	fmt.Fprintf(&buf, "## 📈 基本統計\n\n")
	fmt.Fprintf(&buf, "| 指標 | 値 |\n")
	fmt.Fprintf(&buf, "|------|----|\n")
	fmt.Fprintf(&buf, "| 総テスト数 | %d |\n", report.Summary.TotalTests)
	fmt.Fprintf(&buf, "| 成功テスト数 | %d |\n", report.Summary.SuccessfulTests)
	fmt.Fprintf(&buf, "| 成功率 | %.1f%% |\n", report.Summary.SuccessRate)
	fmt.Fprintf(&buf, "| 平均コンパイル時間 | %v |\n", report.Summary.AvgCompileTime)
	fmt.Fprintf(&buf, "| 平均実行時間 | %v |\n", report.Summary.AvgExecuteTime)
	fmt.Fprintf(&buf, "| 平均メモリ使用量 | %d KB |\n", report.Summary.AvgMemoryUsage)
	fmt.Fprintf(&buf, "| 平均バイナリサイズ | %d bytes |\n\n", report.Summary.AvgBinarySize)

	// 産業標準比較
	fmt.Fprintf(&buf, "## 🏁 産業標準比較\n\n")
	fmt.Fprintf(&buf, "### vs GCC\n\n")
	fmt.Fprintf(&buf, "| 指標 | 比率 | グレード |\n")
	fmt.Fprintf(&buf, "|------|------|----------|\n")
	fmt.Fprintf(&buf, "| 実行時間 | %.2fx | %s |\n", report.Summary.GCCComparison.AvgRuntimeRatio, report.Summary.GCCComparison.Grade)
	fmt.Fprintf(&buf, "| コンパイル時間 | %.2fx | - |\n", report.Summary.GCCComparison.AvgCompileRatio)
	fmt.Fprintf(&buf, "| バイナリサイズ | %.2fx | - |\n", report.Summary.GCCComparison.AvgBinaryRatio)
	fmt.Fprintf(&buf, "| メモリ使用量 | %.2fx | - |\n\n", report.Summary.GCCComparison.AvgMemoryRatio)

	fmt.Fprintf(&buf, "### vs Rust\n\n")
	fmt.Fprintf(&buf, "| 指標 | 比率 | グレード |\n")
	fmt.Fprintf(&buf, "|------|------|----------|\n")
	fmt.Fprintf(&buf, "| 実行時間 | %.2fx | %s |\n", report.Summary.RustComparison.AvgRuntimeRatio, report.Summary.RustComparison.Grade)
	fmt.Fprintf(&buf, "| コンパイル時間 | %.2fx | - |\n", report.Summary.RustComparison.AvgCompileRatio)
	fmt.Fprintf(&buf, "| バイナリサイズ | %.2fx | - |\n", report.Summary.RustComparison.AvgBinaryRatio)
	fmt.Fprintf(&buf, "| メモリ使用量 | %.2fx | - |\n\n", report.Summary.RustComparison.AvgMemoryRatio)

	// 次フェーズ目標
	fmt.Fprintf(&buf, "## 🎯 次フェーズ目標\n\n")
	for _, goal := range report.Summary.NextPhaseGoals {
		fmt.Fprintf(&buf, "- %s\n", goal)
	}
	fmt.Fprintf(&buf, "\n")

	// 改善推奨事項
	fmt.Fprintf(&buf, "## 💡 改善推奨事項\n\n")
	for _, rec := range report.Recommendations {
		fmt.Fprintf(&buf, "- %s\n", rec)
	}
	fmt.Fprintf(&buf, "\n")

	// リンク
	fmt.Fprintf(&buf, "## 📚 詳細情報\n\n")
	fmt.Fprintf(&buf, "- [フェーズ別詳細](./%s-Detail)\n", strings.ToUpper(report.Phase))
	fmt.Fprintf(&buf, "- [GCC比較詳細](./GCC-Comparison)\n")
	fmt.Fprintf(&buf, "- [Rust比較詳細](./Rust-Comparison)\n")
	fmt.Fprintf(&buf, "- [進化履歴](./Performance-Evolution)\n\n")

	// フッター
	fmt.Fprintf(&buf, "---\n\n")
	fmt.Fprintf(&buf, "**📝 Creator**: Claude Code  \n")
	fmt.Fprintf(&buf, "**📅 Generated**: %s  \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&buf, "**🔍 Analysis Method**: 包括的ベンチマーク・比較システム  \n")
	fmt.Fprintf(&buf, "**📊 Data Reliability**: CI/CD自動生成（高信頼性）  \n\n")
	fmt.Fprintf(&buf, "🤖 Generated with [Claude Code](https://claude.ai/code)\n")

	return buf.String()
}

// updatePhaseDetailPage はフェーズ別詳細ページを更新
func (wu *WikiUpdater) updatePhaseDetailPage(wikiDir string, report *BenchmarkReport) error {
	content := wu.generatePhaseDetailContent(report)

	filename := filepath.Join(wikiDir, fmt.Sprintf("%s-Detail.md", strings.ToUpper(report.Phase)))
	return os.WriteFile(filename, []byte(content), 0600)
}

// generatePhaseDetailContent はフェーズ別詳細ページの内容を生成
func (wu *WikiUpdater) generatePhaseDetailContent(report *BenchmarkReport) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "# 📋 %s 詳細ベンチマーク結果\n\n", strings.ToUpper(report.Phase))
	fmt.Fprintf(&buf, "**更新日時**: %s\n\n", report.Timestamp.Format("2006-01-02 15:04:05 JST"))

	// コンパイラベンチマーク詳細
	fmt.Fprintf(&buf, "## 🔧 コンパイラベンチマーク詳細\n\n")

	if len(report.CompilerResults) > 0 {
		fmt.Fprintf(&buf, "| テスト | コンパイル時間 | 実行時間 | メモリ使用量 | バイナリサイズ | スループット | 状態 |\n")
		fmt.Fprintf(&buf, "|--------|----------------|----------|--------------|---------------|-------------|------|\n")

		for _, result := range report.CompilerResults {
			status := "❌"
			if result.Success {
				status = "✅"
			}

			fmt.Fprintf(&buf, "| %s | %v | %v | %d KB | %d bytes | %d ops/sec | %s |\n",
				result.Phase, result.CompileTime, result.ExecuteTime,
				result.MemoryUsage, result.BinarySize, result.ThroughputOps, status)
		}
		fmt.Fprintf(&buf, "\n")
	}

	// エラー詳細
	fmt.Fprintf(&buf, "## ❌ エラー詳細\n\n")
	hasErrors := false
	for _, result := range report.CompilerResults {
		if !result.Success && result.ErrorMessage != "" {
			fmt.Fprintf(&buf, "### %s\n\n", result.Phase)
			fmt.Fprintf(&buf, "```\n%s\n```\n\n", result.ErrorMessage)
			hasErrors = true
		}
	}

	if !hasErrors {
		fmt.Fprintf(&buf, "エラーはありません。\n\n")
	}

	// パフォーマンス分析
	fmt.Fprintf(&buf, "## 📊 パフォーマンス分析\n\n")
	fmt.Fprintf(&buf, "### 実行時間分布\n\n")

	if len(report.CompilerResults) > 0 {
		// 実行時間の分析
		var executeTimes []time.Duration
		for _, result := range report.CompilerResults {
			if result.Success {
				executeTimes = append(executeTimes, result.ExecuteTime)
			}
		}

		if len(executeTimes) > 0 {
			sort.Slice(executeTimes, func(i, j int) bool {
				return executeTimes[i] < executeTimes[j]
			})

			fmt.Fprintf(&buf, "- **最速**: %v\n", executeTimes[0])
			fmt.Fprintf(&buf, "- **最遅**: %v\n", executeTimes[len(executeTimes)-1])
			fmt.Fprintf(&buf, "- **中央値**: %v\n", executeTimes[len(executeTimes)/2])

			// 分散計算
			var sum time.Duration
			for _, t := range executeTimes {
				sum += t
			}
			avg := sum / time.Duration(len(executeTimes))
			fmt.Fprintf(&buf, "- **平均**: %v\n\n", avg)
		}
	}

	// 推奨事項
	fmt.Fprintf(&buf, "## 💡 %s固有の推奨事項\n\n", strings.ToUpper(report.Phase))

	switch report.Phase {
	case "phase1":
		fmt.Fprintf(&buf, "- 🔤 字句解析の最適化\n")
		fmt.Fprintf(&buf, "- 🌳 AST構築の効率化\n")
		fmt.Fprintf(&buf, "- 🧮 評価器の高速化\n")
		fmt.Fprintf(&buf, "- 📚 Phase2への準備\n")
	case "phase2":
		fmt.Fprintf(&buf, "- ⚙️ コード生成の最適化\n")
		fmt.Fprintf(&buf, "- 🎯 型システムの強化\n")
		fmt.Fprintf(&buf, "- 🔄 制御構造の効率化\n")
		fmt.Fprintf(&buf, "- 📈 Phase3への準備\n")
	case "phase3":
		fmt.Fprintf(&buf, "- 🎛️ IR最適化パスの実装\n")
		fmt.Fprintf(&buf, "- 🔍 SSA形式の効率化\n")
		fmt.Fprintf(&buf, "- 🚀 高度な最適化技法\n")
		fmt.Fprintf(&buf, "- 🔗 Phase4への準備\n")
	case "phase4":
		fmt.Fprintf(&buf, "- 🔗 LLVM統合の完成\n")
		fmt.Fprintf(&buf, "- 🌐 マルチターゲット対応\n")
		fmt.Fprintf(&buf, "- 🏭 産業レベル品質\n")
		fmt.Fprintf(&buf, "- 🚀 最終最適化\n")
	}
	fmt.Fprintf(&buf, "\n")

	// ナビゲーション
	fmt.Fprintf(&buf, "## 📚 ナビゲーション\n\n")
	fmt.Fprintf(&buf, "- [← メインページ](./Performance-Benchmark)\n")
	fmt.Fprintf(&buf, "- [GCC比較 →](./GCC-Comparison)\n")
	fmt.Fprintf(&buf, "- [Rust比較 →](./Rust-Comparison)\n")
	fmt.Fprintf(&buf, "- [進化履歴 →](./Performance-Evolution)\n\n")

	// フッター
	fmt.Fprintf(&buf, "---\n\n")
	fmt.Fprintf(&buf, "🤖 Generated with [Claude Code](https://claude.ai/code)\n")

	return buf.String()
}

// updateComparisonPages は比較結果ページを更新
func (wu *WikiUpdater) updateComparisonPages(wikiDir string, report *BenchmarkReport) error {
	// GCC比較ページ
	gccContent := wu.generateGCCComparisonContent(report)
	err := os.WriteFile(filepath.Join(wikiDir, "GCC-Comparison.md"), []byte(gccContent), 0600)
	if err != nil {
		return err
	}

	// Rust比較ページ
	rustContent := wu.generateRustComparisonContent(report)
	err = os.WriteFile(filepath.Join(wikiDir, "Rust-Comparison.md"), []byte(rustContent), 0600)
	if err != nil {
		return err
	}

	return nil
}

// generateGCCComparisonContent はGCC比較ページの内容を生成
func (wu *WikiUpdater) generateGCCComparisonContent(report *BenchmarkReport) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "# 🏁 Pug vs GCC 比較ベンチマーク\n\n")
	fmt.Fprintf(&buf, "**更新日時**: %s\n\n", report.Timestamp.Format("2006-01-02 15:04:05 JST"))

	// サマリー
	fmt.Fprintf(&buf, "## 📊 比較サマリー\n\n")
	fmt.Fprintf(&buf, "| 指標 | 平均比率 | 評価 |\n")
	fmt.Fprintf(&buf, "|------|----------|------|\n")
	fmt.Fprintf(&buf, "| 実行時間 | %.2fx | %s |\n", report.Summary.GCCComparison.AvgRuntimeRatio, report.Summary.GCCComparison.Grade)
	fmt.Fprintf(&buf, "| コンパイル時間 | %.2fx | - |\n", report.Summary.GCCComparison.AvgCompileRatio)
	fmt.Fprintf(&buf, "| バイナリサイズ | %.2fx | - |\n", report.Summary.GCCComparison.AvgBinaryRatio)
	fmt.Fprintf(&buf, "| メモリ使用量 | %.2fx | - |\n\n", report.Summary.GCCComparison.AvgMemoryRatio)

	// 詳細結果
	if len(report.GCCComparisons) > 0 {
		fmt.Fprintf(&buf, "## 📋 詳細比較結果\n\n")
		fmt.Fprintf(&buf, "| テスト | 最適化 | 実行時間比 | コンパイル比 | Pug状態 | GCC状態 |\n")
		fmt.Fprintf(&buf, "|--------|--------|------------|-------------|---------|----------|\n")

		for _, comp := range report.GCCComparisons {
			pugStatus := "❌"
			if comp.PugSuccess {
				pugStatus = "✅"
			}
			gccStatus := "❌"
			if comp.GCCSuccess {
				gccStatus = "✅"
			}

			fmt.Fprintf(&buf, "| %s | %s | %.2fx | %.2fx | %s | %s |\n",
				comp.TestName, comp.OptLevel, comp.RuntimeSpeedRatio, comp.CompileSpeedRatio,
				pugStatus, gccStatus)
		}
		fmt.Fprintf(&buf, "\n")
	}

	// 分析
	fmt.Fprintf(&buf, "## 🔍 分析\n\n")

	if report.Summary.GCCComparison.AvgRuntimeRatio <= 1.0 {
		fmt.Fprintf(&buf, "🎉 **優秀**: PugがGCCと同等以上の性能を発揮しています！\n\n")
	} else if report.Summary.GCCComparison.AvgRuntimeRatio <= 2.0 {
		fmt.Fprintf(&buf, "✅ **良好**: PugはGCCの2倍以内の実行時間です。\n\n")
	} else if report.Summary.GCCComparison.AvgRuntimeRatio <= 10.0 {
		fmt.Fprintf(&buf, "⚠️ **改善余地**: PugはGCCより遅いですが、許容範囲内です。\n\n")
	} else {
		fmt.Fprintf(&buf, "🔧 **要改善**: Pugの性能向上が必要です。\n\n")
	}

	// GCCについて
	fmt.Fprintf(&buf, "## 🏁 GCCについて\n\n")
	fmt.Fprintf(&buf, "GCC (GNU Compiler Collection) は業界標準のCコンパイラです。\n\n")
	fmt.Fprintf(&buf, "### 特徴\n")
	fmt.Fprintf(&buf, "- 🏭 産業レベルの成熟したコンパイラ\n")
	fmt.Fprintf(&buf, "- ⚡ 高度な最適化機能\n")
	fmt.Fprintf(&buf, "- 🌐 多プラットフォーム対応\n")
	fmt.Fprintf(&buf, "- 📈 長年の最適化ノウハウ蓄積\n\n")

	// 目標
	fmt.Fprintf(&buf, "## 🎯 フェーズ別目標\n\n")
	fmt.Fprintf(&buf, "| フェーズ | 目標実行時間比 | 現状 |\n")
	fmt.Fprintf(&buf, "|----------|----------------|------|\n")
	fmt.Fprintf(&buf, "| Phase 1 | 10-100x slower | - |\n")
	fmt.Fprintf(&buf, "| Phase 2 | 2-10x slower | - |\n")
	fmt.Fprintf(&buf, "| Phase 3 | 1-2x slower | - |\n")
	fmt.Fprintf(&buf, "| Phase 4 | GCC同等 | - |\n\n")

	// ナビゲーション
	fmt.Fprintf(&buf, "## 📚 ナビゲーション\n\n")
	fmt.Fprintf(&buf, "- [← メインページ](./Performance-Benchmark)\n")
	fmt.Fprintf(&buf, "- [Rust比較 →](./Rust-Comparison)\n")
	fmt.Fprintf(&buf, "- [進化履歴 →](./Performance-Evolution)\n\n")

	fmt.Fprintf(&buf, "---\n\n")
	fmt.Fprintf(&buf, "🤖 Generated with [Claude Code](https://claude.ai/code)\n")

	return buf.String()
}

// generateRustComparisonContent はRust比較ページの内容を生成
func (wu *WikiUpdater) generateRustComparisonContent(report *BenchmarkReport) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "# 🦀 Pug vs Rust 比較ベンチマーク\n\n")
	fmt.Fprintf(&buf, "**更新日時**: %s\n\n", report.Timestamp.Format("2006-01-02 15:04:05 JST"))

	// サマリー
	fmt.Fprintf(&buf, "## 📊 比較サマリー\n\n")
	fmt.Fprintf(&buf, "| 指標 | 平均比率 | 評価 |\n")
	fmt.Fprintf(&buf, "|------|----------|------|\n")
	fmt.Fprintf(&buf, "| 実行時間 | %.2fx | %s |\n", report.Summary.RustComparison.AvgRuntimeRatio, report.Summary.RustComparison.Grade)
	fmt.Fprintf(&buf, "| コンパイル時間 | %.2fx | - |\n", report.Summary.RustComparison.AvgCompileRatio)
	fmt.Fprintf(&buf, "| バイナリサイズ | %.2fx | - |\n", report.Summary.RustComparison.AvgBinaryRatio)
	fmt.Fprintf(&buf, "| メモリ使用量 | %.2fx | - |\n\n", report.Summary.RustComparison.AvgMemoryRatio)

	// 詳細結果
	if len(report.RustComparisons) > 0 {
		fmt.Fprintf(&buf, "## 📋 詳細比較結果\n\n")
		fmt.Fprintf(&buf, "| テスト | ビルド | 実行時間比 | コンパイル比 | Pug状態 | Rust状態 |\n")
		fmt.Fprintf(&buf, "|--------|--------|------------|-------------|---------|----------|\n")

		for _, comp := range report.RustComparisons {
			pugStatus := "❌"
			if comp.PugSuccess {
				pugStatus = "✅"
			}
			rustStatus := "❌"
			if comp.RustSuccess {
				rustStatus = "✅"
			}

			fmt.Fprintf(&buf, "| %s | %s | %.2fx | %.2fx | %s | %s |\n",
				comp.TestName, comp.OptLevel, comp.RuntimeSpeedRatio, comp.CompileSpeedRatio,
				pugStatus, rustStatus)
		}
		fmt.Fprintf(&buf, "\n")
	}

	// 分析
	fmt.Fprintf(&buf, "## 🔍 分析\n\n")

	if report.Summary.RustComparison.AvgRuntimeRatio <= 1.0 {
		fmt.Fprintf(&buf, "🎉 **驚異的**: PugがRustと同等以上の性能！ゼロコスト抽象化レベルです。\n\n")
	} else if report.Summary.RustComparison.AvgRuntimeRatio <= 2.0 {
		fmt.Fprintf(&buf, "🦀 **素晴らしい**: PugはRustの2倍以内の性能です。\n\n")
	} else if report.Summary.RustComparison.AvgRuntimeRatio <= 10.0 {
		fmt.Fprintf(&buf, "⚠️ **改善余地**: PugはRustより遅いですが、まだ実用的な範囲です。\n\n")
	} else if report.Summary.RustComparison.AvgRuntimeRatio <= 100.0 {
		fmt.Fprintf(&buf, "🔧 **要改善**: Pugの大幅な性能向上が必要です。\n\n")
	} else {
		fmt.Fprintf(&buf, "📚 **学習段階**: 基本機能の実装段階です。\n\n")
	}

	// Rustについて
	fmt.Fprintf(&buf, "## 🦀 Rustについて\n\n")
	fmt.Fprintf(&buf, "Rust は現代的なシステムプログラミング言語で、パフォーマンスと安全性を両立しています。\n\n")
	fmt.Fprintf(&buf, "### 特徴\n")
	fmt.Fprintf(&buf, "- 🚀 ゼロコスト抽象化による高速実行\n")
	fmt.Fprintf(&buf, "- 🛡️ メモリ安全性の保証\n")
	fmt.Fprintf(&buf, "- ⚡ 強力な最適化コンパイラ\n")
	fmt.Fprintf(&buf, "- 🦀 所有権システムによる効率的メモリ管理\n")
	fmt.Fprintf(&buf, "- ⏱️ 長いコンパイル時間（トレードオフ）\n\n")

	// 目標
	fmt.Fprintf(&buf, "## 🎯 フェーズ別目標（vs Rust）\n\n")
	fmt.Fprintf(&buf, "| フェーズ | 目標実行時間比 | 特徴 |\n")
	fmt.Fprintf(&buf, "|----------|----------------|------|\n")
	fmt.Fprintf(&buf, "| Phase 1 | 100-1000x slower | インタープリター段階 |\n")
	fmt.Fprintf(&buf, "| Phase 2 | 10-50x slower | 基本コンパイラ |\n")
	fmt.Fprintf(&buf, "| Phase 3 | 2-5x slower | 最適化コンパイラ |\n")
	fmt.Fprintf(&buf, "| Phase 4 | Rust同等 | ゼロコスト抽象化達成 |\n\n")

	// 学習ポイント
	fmt.Fprintf(&buf, "## 💡 学習ポイント\n\n")
	fmt.Fprintf(&buf, "- **コンパイル時間**: Rustは最適化に時間をかける。Pugは軽量・高速コンパイルを目指す\n")
	fmt.Fprintf(&buf, "- **実行時性能**: Rustのゼロコスト抽象化がゴール\n")
	fmt.Fprintf(&buf, "- **メモリ効率**: Rustの所有権システムから学ぶ\n")
	fmt.Fprintf(&buf, "- **最適化**: LLVMを活用した高度な最適化\n\n")

	// ナビゲーション
	fmt.Fprintf(&buf, "## 📚 ナビゲーション\n\n")
	fmt.Fprintf(&buf, "- [← メインページ](./Performance-Benchmark)\n")
	fmt.Fprintf(&buf, "- [← GCC比較](./GCC-Comparison)\n")
	fmt.Fprintf(&buf, "- [進化履歴 →](./Performance-Evolution)\n\n")

	fmt.Fprintf(&buf, "---\n\n")
	fmt.Fprintf(&buf, "🤖 Generated with [Claude Code](https://claude.ai/code)\n")

	return buf.String()
}

// updateEvolutionHistoryPage は進化履歴ページを更新
func (wu *WikiUpdater) updateEvolutionHistoryPage(wikiDir string, report *BenchmarkReport) error {
	// 既存の履歴を読み込み、新しい結果を追加
	filename := filepath.Join(wikiDir, "Performance-Evolution.md")

	var existingContent string
	// Validate file path before reading
	if err := validateFilePath(filename); err == nil {
		if data, err := os.ReadFile(filename); err == nil { // #nosec G304
			existingContent = string(data)
		}
	}

	// 新しい履歴エントリを生成
	newEntry := wu.generateEvolutionEntry(report)

	// 履歴ページの内容を生成
	content := wu.generateEvolutionContent(existingContent, newEntry, report)

	return os.WriteFile(filename, []byte(content), 0600)
}

// generateEvolutionEntry は新しい進化履歴エントリを生成
func (wu *WikiUpdater) generateEvolutionEntry(report *BenchmarkReport) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "### %s - %s\n\n", report.Timestamp.Format("2006-01-02"), strings.ToUpper(report.Phase))
	fmt.Fprintf(&buf, "- **性能グレード**: %s\n", report.Summary.PerformanceGrade)
	fmt.Fprintf(&buf, "- **成功率**: %.1f%% (%d/%d)\n", report.Summary.SuccessRate, report.Summary.SuccessfulTests, report.Summary.TotalTests)
	fmt.Fprintf(&buf, "- **平均実行時間**: %v\n", report.Summary.AvgExecuteTime)
	fmt.Fprintf(&buf, "- **vs GCC**: %.2fx (%s)\n", report.Summary.GCCComparison.AvgRuntimeRatio, report.Summary.GCCComparison.Grade)
	fmt.Fprintf(&buf, "- **vs Rust**: %.2fx (%s)\n", report.Summary.RustComparison.AvgRuntimeRatio, report.Summary.RustComparison.Grade)
	fmt.Fprintf(&buf, "\n")

	return buf.String()
}

// generateEvolutionContent は進化履歴ページの完全な内容を生成
func (wu *WikiUpdater) generateEvolutionContent(existingContent, newEntry string, report *BenchmarkReport) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "# 📈 Pugコンパイラ性能進化履歴\n\n")
	fmt.Fprintf(&buf, "**最終更新**: %s\n\n", report.Timestamp.Format("2006-01-02 15:04:05 JST"))

	// 進化グラフ（簡易版）
	fmt.Fprintf(&buf, "## 📊 性能進化グラフ\n\n")
	fmt.Fprintf(&buf, "```\n")
	fmt.Fprintf(&buf, "Phase 1 → Phase 2 → Phase 3 → Phase 4\n")
	fmt.Fprintf(&buf, "  📚      ⚙️       🎯       🚀\n")
	fmt.Fprintf(&buf, "学習     基盤     最適化   産業レベル\n")
	fmt.Fprintf(&buf, "```\n\n")

	// 目標と現状
	fmt.Fprintf(&buf, "## 🎯 フェーズ別目標と現状\n\n")
	fmt.Fprintf(&buf, "| フェーズ | 目標性能向上 | vs GCC目標 | vs Rust目標 | 現状 |\n")
	fmt.Fprintf(&buf, "|----------|-------------|-------------|------------|------|\n")
	fmt.Fprintf(&buf, "| Phase 1 | ベースライン | 10-100x slower | 100-1000x slower | 実装中 |\n")
	fmt.Fprintf(&buf, "| Phase 2 | 10x向上 | 2-10x slower | 10-50x slower | 準備中 |\n")
	fmt.Fprintf(&buf, "| Phase 3 | 50x向上 | 1-2x slower | 2-5x slower | 予定 |\n")
	fmt.Fprintf(&buf, "| Phase 4 | 100x向上 | GCC同等 | Rust同等 | 予定 |\n\n")

	// 履歴セクション
	fmt.Fprintf(&buf, "## 📅 実装履歴\n\n")

	// 新しいエントリを追加
	fmt.Fprintf(&buf, "%s", newEntry)

	// 既存の履歴エントリがあれば追加（簡易実装）
	if existingContent != "" {
		// 既存コンテンツから履歴部分を抽出（簡易版）
		if strings.Contains(existingContent, "## 📅 実装履歴") {
			parts := strings.Split(existingContent, "## 📅 実装履歴")
			if len(parts) > 1 {
				historyPart := parts[1]
				// 次のセクションまでを取得
				if idx := strings.Index(historyPart, "\n## "); idx != -1 {
					historyPart = historyPart[:idx]
				}
				// 新しいエントリ以外の部分を追加
				lines := strings.Split(historyPart, "\n")
				inEntry := false
				for _, line := range lines {
					if strings.HasPrefix(line, "### ") {
						inEntry = true
					}
					if inEntry && strings.TrimSpace(line) != "" {
						fmt.Fprintf(&buf, "%s\n", line)
					}
				}
			}
		}
	}

	// マイルストーン
	fmt.Fprintf(&buf, "\n## 🏆 マイルストーン\n\n")
	fmt.Fprintf(&buf, "- [ ] Phase 1完成: 基本インタープリター\n")
	fmt.Fprintf(&buf, "- [ ] Phase 2完成: アセンブリコンパイラ\n")
	fmt.Fprintf(&buf, "- [ ] Phase 3完成: 最適化コンパイラ\n")
	fmt.Fprintf(&buf, "- [ ] Phase 4完成: LLVM統合コンパイラ\n")
	fmt.Fprintf(&buf, "- [ ] GCC性能達成: 産業レベル到達\n")
	fmt.Fprintf(&buf, "- [ ] Rust性能達成: ゼロコスト抽象化\n\n")

	// 技術的進歩
	fmt.Fprintf(&buf, "## 🔧 技術的進歩\n\n")
	fmt.Fprintf(&buf, "### 実装済み\n")
	fmt.Fprintf(&buf, "- ✅ 字句解析器 (Lexer)\n")
	fmt.Fprintf(&buf, "- ✅ 構文解析器 (Parser)\n")
	fmt.Fprintf(&buf, "- ✅ 抽象構文木 (AST)\n")
	fmt.Fprintf(&buf, "- ✅ 評価器 (Evaluator)\n")
	fmt.Fprintf(&buf, "- ✅ 制御構造 (if/while/for)\n")
	fmt.Fprintf(&buf, "- ✅ 包括的ベンチマークシステム\n\n")

	fmt.Fprintf(&buf, "### 開発中\n")
	fmt.Fprintf(&buf, "- 🔧 x86_64アセンブリ生成\n")
	fmt.Fprintf(&buf, "- 🔧 型システム強化\n")
	fmt.Fprintf(&buf, "- 🔧 エラーハンドリング改善\n\n")

	fmt.Fprintf(&buf, "### 予定\n")
	fmt.Fprintf(&buf, "- 📋 IR (中間表現) 設計\n")
	fmt.Fprintf(&buf, "- 📋 SSA形式対応\n")
	fmt.Fprintf(&buf, "- 📋 最適化パス実装\n")
	fmt.Fprintf(&buf, "- 📋 LLVM統合\n")
	fmt.Fprintf(&buf, "- 📋 多言語バックエンド\n\n")

	// ナビゲーション
	fmt.Fprintf(&buf, "## 📚 ナビゲーション\n\n")
	fmt.Fprintf(&buf, "- [← メインページ](./Performance-Benchmark)\n")
	fmt.Fprintf(&buf, "- [← GCC比較](./GCC-Comparison)\n")
	fmt.Fprintf(&buf, "- [← Rust比較](./Rust-Comparison)\n\n")

	fmt.Fprintf(&buf, "---\n\n")
	fmt.Fprintf(&buf, "🤖 Generated with [Claude Code](https://claude.ai/code)\n")

	return buf.String()
}

// commitAndPush は変更をコミット・プッシュ
func (wu *WikiUpdater) commitAndPush(report *BenchmarkReport) error {
	wikiDir := filepath.Join(wu.TempDir, "wiki")

	// Git設定 - validate inputs
	if err := validateGitInput(wu.CommitUser); err != nil {
		return fmt.Errorf("invalid commit user: %v", err)
	}
	if err := validateGitInput(wu.CommitEmail); err != nil {
		return fmt.Errorf("invalid commit email: %v", err)
	}

	cmd := exec.Command("git", "config", "user.name", wu.CommitUser) // #nosec G204 - validated input for git config
	cmd.Dir = wikiDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git user設定失敗: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", wu.CommitEmail) // #nosec G204 - validated input for git config
	cmd.Dir = wikiDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git email設定失敗: %v", err)
	}

	// 変更をステージング
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = wikiDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git add失敗: %v", err)
	}

	// コミット
	commitMsg := fmt.Sprintf("🚀 %s性能ベンチマーク結果更新\n\n"+
		"- 実行日時: %s\n"+
		"- 性能グレード: %s\n"+
		"- 成功率: %.1f%%\n"+
		"- vs GCC: %.2fx\n"+
		"- vs Rust: %.2fx\n\n"+
		"🤖 Generated with Claude Code",
		strings.ToUpper(report.Phase),
		report.Timestamp.Format("2006-01-02 15:04:05"),
		report.Summary.PerformanceGrade,
		report.Summary.SuccessRate,
		report.Summary.GCCComparison.AvgRuntimeRatio,
		report.Summary.RustComparison.AvgRuntimeRatio)

	// Validate commit message
	if err := validateGitInput(commitMsg); err != nil {
		// If validation fails, use a safe default message
		commitMsg = "Update benchmark results"
	}

	cmd = exec.Command("git", "commit", "-m", commitMsg) // #nosec G204 - validated commit message
	cmd.Dir = wikiDir
	if err := cmd.Run(); err != nil {
		// コミットが失敗した場合（変更がない場合など）
		return nil // エラーとしない
	}

	// プッシュ
	cmd = exec.Command("git", "push", "origin", "master")
	cmd.Dir = wikiDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git push失敗: %v", err)
	}

	return nil
}
