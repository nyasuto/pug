# 📊 Pugコンパイラ包括的性能ベンチマーク・比較システム

## 概要

このディレクトリには、Pugコンパイラの性能を定量的に測定し、産業標準コンパイラ（GCC、Rust）との比較を行う包括的ベンチマークシステムが含まれています。

## 🎯 目的

- **性能測定**: 各フェーズの性能を定量的に測定
- **進化追跡**: フェーズ間での性能向上を追跡
- **産業比較**: GCC、Rustとの性能比較
- **回帰検出**: 性能退化の早期発見
- **目標設定**: 次フェーズの具体的目標設定

## 📁 ファイル構成

```
benchmark/
├── README.md                 # このファイル
├── compiler_bench.go         # コンパイラベンチマーク（コア機能）
├── vs_gcc.go                # GCC比較ベンチマーク
├── vs_rust.go               # Rust比較ベンチマーク
├── report.go                # レポート生成・可視化
├── wiki_update.go           # GitHub Wiki自動更新
├── benchmark_test.go        # テストスイート
└── doc.go                   # パッケージドキュメント
```

## 🚀 クイックスタート

### 基本ベンチマーク実行

```bash
# 全ベンチマーク実行
make bench

# コンパイラベンチマークのみ
go test -bench=BenchmarkCompiler ./benchmark/...

# GCC比較ベンチマーク
go test -bench=BenchmarkVsGCC ./benchmark/...

# Rust比較ベンチマーク  
go test -bench=BenchmarkVsRust ./benchmark/...

# 包括的ベンチマークスイート
go test -bench=BenchmarkSuite ./benchmark/...
```

### レポート生成

```bash
# JSONレポート生成
go test -bench=. ./benchmark/... -run=TestGenerateReport

# HTMLレポート生成
go test -bench=. ./benchmark/... -run=TestHTMLReport
```

## 📊 ベンチマーク種類

### 1. コンパイラベンチマーク（compiler_bench.go）

**目的**: Pugコンパイラ自体の性能測定

**測定項目**:
- コンパイル時間
- 実行時間  
- メモリ使用量
- バイナリサイズ
- スループット（ops/sec）

**テストケース**:
- フィボナッチ計算（再帰）
- ソートアルゴリズム
- 数値計算（円周率計算）
- 複雑な制御構造

**実行方法**:
```bash
go test -bench=BenchmarkCompiler_Phase1 ./benchmark/...
go test -bench=BenchmarkCompiler_Phase2 ./benchmark/...
```

### 2. GCC比較ベンチマーク（vs_gcc.go）

**目的**: 産業標準Cコンパイラとの性能比較

**比較対象**:
- GCC -O0 (最適化なし)
- GCC -O1 (基本最適化)
- GCC -O2 (推奨最適化)
- GCC -O3 (積極的最適化)

**比較指標**:
- 実行時間比率（Pug/GCC）
- コンパイル時間比率
- バイナリサイズ比率
- メモリ使用量比率

**実行方法**:
```bash
go test -bench=BenchmarkVsGCC_O2 ./benchmark/...
go test -bench=BenchmarkPugCompilerEvolution ./benchmark/...
```

### 3. Rust比較ベンチマーク（vs_rust.go）

**目的**: 現代的システム言語との性能比較

**比較対象**:
- Rust debug ビルド
- Rust release ビルド（最適化済み）

**特徴**:
- ゼロコスト抽象化との比較
- メモリ安全性とパフォーマンスの関係
- 現代的コンパイラ技術との比較

**実行方法**:
```bash
go test -bench=BenchmarkVsRust_Release ./benchmark/...
go test -bench=BenchmarkPugVsRustEvolution ./benchmark/...
```

## 📈 レポートシステム

### レポート機能（report.go）

**生成されるレポート**:
- JSON形式の詳細データ
- HTML形式の可視化レポート
- 進化分析レポート
- 性能グレード評価

**主要コンポーネント**:
```go
type BenchmarkReport struct {
    Timestamp        time.Time
    Phase            string
    Environment      EnvironmentInfo
    CompilerResults  []*BenchmarkResult
    GCCComparisons   []*ComparisonResult
    RustComparisons  []*RustComparisonResult
    Summary          BenchmarkSummary
    Recommendations  []string
}
```

**使用例**:
```go
// レポート生成
report := GenerateComprehensiveReport("phase1", compilerResults, gccResults, rustResults)

// JSON保存
report.SaveReportJSON("benchmark-report.json")

// HTML生成
report.GenerateHTMLReport("benchmark-report.html")
```

### 性能グレード評価

**評価基準**:
- **S+**: 産業レベル（GCC同等以上）
- **S**: 優秀（GCCの2倍以内）
- **A**: 良好（GCCの5倍以内）
- **B**: 基本達成（GCCの10倍以内）
- **C**: 改善必要（GCCの50倍以内）
- **D**: 初期段階（それ以上）

## 🔄 CI/CD統合

### GitHub Actions統合

**自動実行タイミング**:
- メインブランチへのプッシュ時
- プルリクエスト作成時（オプション）

**実行内容**:
1. 基本ベンチマーク実行
2. 包括的ベンチマーク実行
3. 産業標準比較（GCC/Rust利用可能時）
4. レポート生成・保存
5. パフォーマンス回帰検出

**設定**:
```yaml
# .github/workflows/ci.yml の benchmark job で自動実行
# カスタマイズ可能な環境変数:
# - BENCHMARK_TIMEOUT: ベンチマークタイムアウト
# - ENABLE_WIKI_UPDATE: Wiki自動更新有効化
```

### 成果物

**GitHub Actions Artifacts**:
- `benchmark-results-comprehensive`: 全ベンチマーク結果
- `benchmark-*.txt`: 個別ベンチマーク結果
- 保存期間: 30日

## 📝 GitHub Wiki自動更新

### Wiki自動更新機能（wiki_update.go）

**更新されるページ**:
- `Performance-Benchmark.md`: メインベンチマークページ
- `{Phase}-Detail.md`: フェーズ別詳細ページ
- `GCC-Comparison.md`: GCC比較詳細
- `Rust-Comparison.md`: Rust比較詳細
- `Performance-Evolution.md`: 進化履歴

**使用方法**:
```go
// Wiki更新
updater := NewWikiUpdater("https://github.com/owner/repo.git")
err := updater.UpdateBenchmarkWiki(report)
```

**セキュリティ**:
- 読み取り専用操作（既存Wiki構造を保持）
- 自動コミット・プッシュ
- エラー時の安全な停止

## 🎯 フェーズ別目標

### Phase 1: 基本言語処理（インタープリター）
**目標**:
- ✅ 基本機能の安定動作
- ✅ 包括的テストカバレッジ（75%+）
- 🎯 GCCの10-100倍以内の実行時間
- 🎯 Rustの100-1000倍以内の実行時間

### Phase 2: コンパイラ基盤（アセンブリ生成）
**目標**:
- 🎯 Phase1から10倍性能向上
- 🎯 GCCの2-10倍以内の実行時間
- 🎯 基本的なコード最適化

### Phase 3: 最適化エンジン（IR最適化）
**目標**:
- 🎯 Phase2から5倍性能向上
- 🎯 GCCの1-2倍以内の実行時間
- 🎯 高度な最適化パス実装

### Phase 4: 産業レベル（LLVM統合）
**目標**:
- 🎯 Phase3から2倍性能向上
- 🎯 GCC同等の実行時間
- 🎯 Rust同等の最適化レベル

## 🔧 カスタマイズ

### 新しいベンチマーク追加

1. **テストケース追加**:
```go
// compiler_bench.go に新しいプログラムを追加
const newTestProgram = `
let custom_function = fn(n) {
    // カスタムロジック
    return n * 2;
};
`

// 新しいベンチマーク関数を追加
func BenchmarkCompiler_Phase1_CustomTest(b *testing.B) {
    cb, err := setupBenchmark("phase1", newTestProgram)
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
```

2. **新しい比較対象追加**:
```go
// 新しい言語との比較（例：Go）
type GoComparisonResult struct {
    // 比較結果構造体
}

func BenchmarkVsGo(b *testing.B) {
    // Go比較ベンチマーク実装
}
```

### 環境固有設定

```go
// 環境変数での設定
const (
    DefaultTimeout = 30 * time.Second
    DefaultMemLimit = 1024 * 1024 * 1024 // 1GB
)

// カスタムベンチマーク設定
type BenchmarkConfig struct {
    Timeout     time.Duration
    MemoryLimit int64
    Iterations  int
    Parallel    bool
}
```

## 🚨 トラブルシューティング

### よくある問題

**1. GCC/Rustが見つからない**
```bash
# GCCインストール確認
gcc --version

# Rustインストール確認  
cargo --version

# 環境変数確認
echo $PATH
```

**2. ベンチマークタイムアウト**
```bash
# タイムアウト時間延長
go test -bench=. -timeout=10m ./benchmark/...

# 個別テスト実行
go test -bench=BenchmarkCompiler_Phase1_Fibonacci -v ./benchmark/...
```

**3. メモリ不足**
```bash
# メモリ使用量確認
go test -bench=. -benchmem ./benchmark/...

# 並列度調整
go test -bench=. -cpu=1 ./benchmark/...
```

### デバッグ情報

**詳細ログ有効化**:
```bash
# 詳細ベンチマーク情報
go test -bench=. -v ./benchmark/...

# レース検出
go test -bench=. -race ./benchmark/...

# プロファイリング
go test -bench=. -cpuprofile=cpu.prof ./benchmark/...
go test -bench=. -memprofile=mem.prof ./benchmark/...
```

## 📚 参考資料

### ベンチマーク設計
- [Go Testing Package](https://golang.org/pkg/testing/)
- [Benchmarking Go Programs](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
- [Performance Analysis](https://golang.org/doc/diagnostics.html)

### コンパイラ性能
- [GCC Optimization Options](https://gcc.gnu.org/onlinedocs/gcc/Optimize-Options.html)
- [Rust Performance Book](https://nnethercote.github.io/perf-book/)
- [LLVM Optimization Guide](https://llvm.org/docs/Passes.html)

### 統計・分析
- [Benchmarking Statistics](https://golang.org/x/perf/cmd/benchstat)
- [Performance Testing Best Practices](https://github.com/golang/go/wiki/Performance)

## 🤝 貢献

### ベンチマーク改善提案
1. 新しいテストケースの提案
2. 比較対象言語の追加
3. 可視化機能の改善
4. CI/CD統合の強化

### 報告・フィードバック
- 性能異常の報告
- ベンチマーク結果の解釈支援
- システム要求の提案

---

**📝 Creator**: Claude Code  
**📅 Created**: 2024-06-25  
**🔍 Analysis Method**: 包括的ベンチマーク・比較システム設計  
**📊 Data Reliability**: 産業標準比較による高信頼性  

🤖 Generated with [Claude Code](https://claude.ai/code)