# 🐶 pug: 段階的に学ぶコンパイラ実装プロジェクト

**シンプルなレクサーからLLVMレベルの最適化コンパイラまで、段階的に学ぶコンパイラ技術**

--- 

pugは、プログラミング言語処理系の内部構造を、最初は数十行のレクサーだけで理解し、その後段階的に本格的な最適化コンパイラへと進化させる学習プロジェクトです。


## 🌟 コンセプト

パグ犬のように、すこしぶさいけどあいきょうある見た目で、言語処理の中身は詰まっている。最初は簡易なトークン解析から始まり、最終的には高度な最適化を行うLLVM連携コンパイラへ成長します。

## **🎯 プロジェクトの目的**
 

- **コンパイラ技術の実践的学習**: 字句解析から最適化まで段階的に理解

- **言語処理系の比較学習**: 異なるアプローチでの実装比較・分析

- **最適化技術による学習**: 定量的な性能改善効果を可視化

- **実用的な開発ワークフロー習得**: テスト駆動開発、品質管理手法

  

## **📊 現在の実装状況**

  

### **✅ 完了予定機能**

  

| フェーズ | 機能 | 実装状況 | 説明 |

|---------|------|----------|------|

| **1.0** | **シンプルレクサー** | 🔄 計画中 | トークン化、識別子・数値・演算子 |

| **1.1** | **基本パーサー** | 🔄 計画中 | 再帰下降法、AST構築、四則演算 |

| **1.2** | **シンプルインタープリター** | 🔄 計画中 | AST直接実行、変数・関数対応 |

| **2.0** | **コード生成器** | 🔄 計画中 | アセンブリ出力、レジスタ割り当て |

| **2.1** | **型システム** | 🔄 計画中 | 静的型検査、型推論、エラー報告 |

| **2.2** | **制御構造** | 🔄 計画中 | if/while/for、スコープ管理 |

| **3.0** | **中間表現（IR）** | 🔄 計画中 | SSA形式、基本ブロック、CFG |

| **3.1** | **基本最適化** | 🔄 計画中 | 定数畳み込み、デッドコード除去 |

| **3.2** | **高度最適化** | 🔄 計画中 | ループ最適化、インライン展開 |

| **4.0** | **LLVM連携** | 🔄 計画中 | LLVM IR生成、プラットフォーム対応 |

  

## **🚀 劇的な性能向上を実現予定**

  

### **最新性能分析目標**

シンプルなインタープリター実装から最新LLVM最適化コンパイラへの移行により、**産業級コンパイラレベルの性能向上**を達成予定：

  

| 実装 | コンパイル時間 | 実行性能 | 生成コードサイズ | 特徴 |

|------|-------------|----------|----------------|------|

| **Phase 1 インタープリター** | - | 1.0x (ベースライン) | - | シンプル解釈実行 |

| **Phase 2 ナイーブコンパイラ** | 50ms | **10x faster** | 100% | 直接アセンブリ生成 |

| **Phase 3 最適化コンパイラ** | 100ms | **50x faster** | 70% | IR最適化 |

| **🚀 Phase 4 LLVM連携** | 200ms | **100x faster** | **50%** | **産業級最適化** |

  

### **🏆 期待される改善効果**

- **実行性能**: **100倍高速化** (インタープリター → LLVMコンパイラ)

- **コードサイズ**: **50%削減** (最適化による効率化)

- **開発効率**: **段階的学習** による理解度向上

- **移植性**: **LLVM対応** によるマルチプラットフォーム

  

詳細は [GitHub Wiki - Compiler Performance Analysis Report](https://github.com/nyasuto/pug/wiki/Performance-Analysis-Report) を参照予定

  

## **📁 プロジェクト構成予定**

  

```

pug/

├── phase1/ # Phase 1: 基本言語処理

│ ├── lexer.go # 字句解析器（トークナイザー）

│ ├── parser.go # 構文解析器（再帰下降法）

│ ├── ast.go # 抽象構文木定義

│ ├── interpreter.go # シンプルインタープリター

│ └── *_test.go # 各フェーズのテストスイート

│

├── phase2/ # Phase 2: コンパイラ基盤

│ ├── codegen.go # アセンブリコード生成

│ ├── types.go # 型システム・型検査

│ ├── symbols.go # シンボルテーブル・スコープ

│ ├── control.go # 制御構造（if/while/for）

│ └── *_test.go # コンパイラ機能テスト

│

├── phase3/ # Phase 3: 中間表現・最適化

│ ├── ir/ # 中間表現（SSA形式）

│ │ ├── builder.go # IR構築器

│ │ ├── ssa.go # SSA変換・PHI関数

│ │ ├── cfg.go # 制御フローグラフ

│ │ └── analysis.go # データフロー解析

│ │

│ ├── optimizer/ # 最適化パス

│ │ ├── constant_fold.go # 定数畳み込み

│ │ ├── dead_code.go # デッドコード除去

│ │ ├── loop_opt.go # ループ最適化

│ │ └── inline.go # インライン展開

│ │

│ └── backend/ # コード生成バックエンド

│ ├── x86_64.go # x86-64ターゲット

│ ├── arm64.go # ARM64ターゲット

│ └── regalloc.go # レジスタ割り当て

│

├── phase4/ # Phase 4: LLVM統合・産業級

│ ├── llvm/ # LLVM連携

│ │ ├── ir_gen.go # LLVM IR生成

│ │ ├── optimization.go # LLVM最適化パス

│ │ └── targets.go # マルチターゲット対応

│ │

│ ├── runtime/ # ランタイムシステム

│ │ ├── gc.go # ガベージコレクション

│ │ ├── memory.go # メモリ管理

│ │ └── exceptions.go # 例外処理

│ │

│ └── tools/ # 開発ツール

│ ├── debugger.go # デバッガー

│ ├── profiler.go # プロファイラー

│ └── lsp.go # Language Server Protocol

│

├── cmd/

│ ├── pug/main.go # コンパイラ CLI エントリーポイント

│ ├── interp/main.go # インタープリター CLI

│ └── tools/main.go # 開発ツール統合CLI

│

├── examples/ # サンプルプログラム

│ ├── hello.dog # Hello World

│ ├── fibonacci.dog # フィボナッチ数列

│ ├── sorting.dog # ソートアルゴリズム

│ ├── web_server.dog # 簡易Webサーバー

│ └── ray_tracer.dog # レイトレーサー

│

├── benchmark/ # 性能測定・比較ツール

│ ├── compiler_bench.go # コンパイル時間測定

│ ├── runtime_bench.go # 実行時間測定

│ ├── vs_gcc.go # GCC比較ベンチマーク

│ └── vs_rust.go # Rust比較ベンチマーク

│

├── docs/ # 段階的学習ドキュメント

│ ├── phase1_tutorial.md # Phase 1 学習ガイド

│ ├── phase2_tutorial.md # Phase 2 学習ガイド

│ ├── compiler_theory.md # コンパイラ理論解説

│ └── optimization_guide.md # 最適化技術ガイド

│

├── .github/workflows/ci.yml # CI/CD パイプライン

├── Makefile # 統合開発コマンド

├── go.mod # Go module 定義

└── language_spec.md # Dog 言語仕様書

```

  

## **🚀 使用方法予定**

  

### **開発環境セットアップ**

```bash

make help # 利用可能コマンド一覧

make dev # 開発環境セットアップ

```

  

### **Phase 1: インタープリター（学習開始）**

```bash

# ビルド

make phase1-build

  

# Hello World実行

echo 'print("Hello, pug!")' > hello.dog

./bin/interp hello.dog # Output: Hello, pug!

  

# インタラクティブREPL

./bin/interp --repl

pug > 2 + 3 * 4 # Output: 14

pug > let x = 10; x * x # Output: 100

```

  

### **Phase 2: 基本コンパイラ**

```bash

# コンパイル & 実行

./bin/pug hello.dog -o hello # アセンブリ生成 & リンク

./hello # Output: Hello, pug compiler!

  

# アセンブリ出力確認

./bin/pug hello.dog --emit-asm # アセンブリコード表示

./bin/pug hello.dog --emit-ast # AST構造表示

```

  

### **Phase 3: 最適化コンパイラ**

```bash

# 最適化レベル指定

./bin/pug program.dog -O0 # 最適化なし

./bin/pug program.dog -O1 # 基本最適化

./bin/pug program.dog -O2 # 高度最適化

  

# 最適化過程確認

./bin/pug program.dog --emit-ir # 中間表現出力

./bin/pug program.dog --emit-ir-opt # 最適化後IR出力

./bin/pug program.dog --optimization-report # 最適化レポート

```

  

### **Phase 4: LLVM連携（産業級）**

```bash

# LLVM バックエンド使用

./bin/pug program.dog --backend=llvm --target=x86_64

./bin/pug program.dog --backend=llvm --target=arm64

./bin/pug program.dog --backend=llvm --target=wasm

  

# 高度なツール使用

./bin/tools debug program.dog # デバッガー起動

./bin/tools profile program.dog # プロファイリング

./bin/tools lsp # Language Server起動

```

  

### **性能分析・比較**

```bash

# 自動化された包括的性能分析

./benchmark/compiler_comparison.sh # デフォルト測定

./benchmark/vs_industry.sh # 業界標準比較

  

# 個別ベンチマーク

make bench-compile # コンパイル時間測定

make bench-runtime # 実行時間測定

make bench-vs-gcc # GCC比較

make bench-vs-rust # Rust比較

```

  

## **🔧 主要機能詳細予定**

  

### **🎯 Dog 言語仕様（段階的拡張）**

  

#### **Phase 1: 基本機能**

```rust

// 基本データ型

let x: int = 42;

let y: float = 3.14;

let name: string = "Compiler";

let flag: bool = true;

  

// 基本演算

let result = (x + 10) * 2;

print(result); // Output: 104

  

// 関数定義

fn fibonacci(n: int) -> int {

if n <= 1 {

return n;

}

return fibonacci(n-1) + fibonacci(n-2);

}

```

  

#### **Phase 2: 拡張機能**

```rust

// 構造体

struct Point {

x: float,

y: float,

}

  

// 配列・スライス

let numbers: [int] = [1, 2, 3, 4, 5];

let slice = numbers[1..4]; // [2, 3, 4]

  

// ループ・制御構造

for i in 0..10 {

if i % 2 == 0 {

print(i);

}

}

  

// エラーハンドリング

fn divide(a: int, b: int) -> Result<int, string> {

if b == 0 {

return Err("Division by zero");

}

return Ok(a / b);

}

```

  

#### **Phase 3-4: 高度機能**

```rust

// ジェネリクス

fn max<T: Comparable>(a: T, b: T) -> T {

return a > b ? a : b;

}

  

// トレイト・インターフェース

trait Drawable {

fn draw(self);

}

  

// 非同期・並行処理

async fn fetch_data(url: string) -> Result<Data, Error> {

let response = await http::get(url);

return parse_json(response.body);

}

  

// メモリ管理（所有権システム）

fn process_data(data: owned Data) -> ProcessedData {

// dataの所有権を受け取り、変換後に返す

return transform(data);

}

```

  

### **⚡ 最適化技術実装予定**

  

#### **Phase 3: 基本最適化**

- **定数畳み込み**: コンパイル時計算により実行時負荷削減

- **デッドコード除去**: 到達不可能コードの自動削除

- **共通部分式除去**: 重複計算の最適化

- **ループ不変式移動**: ループ内不変計算の外側移動

  

#### **Phase 4: 高度最適化**

- **インライン展開**: 関数呼び出しオーバーヘッド除去

- **ループ最適化**: アンローリング、ベクトル化

- **レジスタ割り当て**: グラフ彩色法による効率的割り当て

- **命令スケジューリング**: CPU最適化・パイプライン活用

  

### **🔍 高度解析機能**

  

#### **静的解析**

- **制御フロー解析**: デッドコード・到達可能性検出

- **データフロー解析**: 変数生存期間・使用-定義チェーン

- **エイリアス解析**: ポインタ・参照の別名解析

- **エスケープ解析**: ヒープ割り当て最適化

  

#### **動的解析・プロファイリング**

- **実行時プロファイリング**: ホットスポット特定

- **メモリ使用量追跡**: リーク検出・最適化指針

- **分岐予測情報**: 条件分岐最適化

- **キャッシュ効率分析**: メモリアクセスパターン最優化

  

## **📈 実用性とスケーラビリティ予定**

  

### **リアルタイム性能目標**

```bash

# コンパイル性能目標（1万行プログラム）

Phase 1 インタープリター: 即座実行 (0.1秒)

Phase 2 基本コンパイラ: 高速コンパイル (1秒)

Phase 3 最適化コンパイラ: バランス (5秒)

Phase 4 LLVM連携: 最高品質 (10秒)

  

# 実行性能目標（フィボナッチ40回）

Phase 1 インタープリター: 10秒

Phase 2 ナイーブコンパイラ: 1秒 (10x faster)

Phase 3 最適化コンパイラ: 0.2秒 (50x faster)

Phase 4 LLVM連携: 0.1秒 (100x faster)

```

  

### **コードサイズ効率性**

- **動的最適化**: 実行時情報による最適化

- **効率的エンコーディング**: 命令選択・アドレッシング最適化

- **リンク時最適化**: 全体最適化・デッドコード除去

  

### **エンタープライズ対応予定**

- **マルチターゲット**: x86_64, ARM64, WASM, RISC-V対応

- **デバッグ情報**: DWARF形式デバッグ情報生成

- **プロファイリング**: gprof, perf対応

- **IDE統合**: Language Server Protocol実装

  

## **🔄 開発ワークフロー予定**

  

### **CI/CDパイプライン**

```bash

# GitHub Actions で自動実行予定

- 品質チェック（lint, format, type-check）

- セキュリティスキャン（gosec, govulncheck）

- 統合テスト（各Phase機能、性能テスト）

- 業界標準比較（GCC, Rust, Clangとの性能比較）

- ブランチ保護・命名規則チェック

```

  

### **品質管理予定**

- **テスト駆動開発**: 各機能の包括的テストスイート

- **性能回帰防止**: 自動化された性能ベンチマーク

- **コードカバレッジ**: 95%以上のテストカバレッジ維持

- **静的解析**: 品質・セキュリティチェック

  

## **🎓 学習効果予定**

  

このプロジェクトを通じて習得できる技術:

  

### **コンパイラ技術**

- **言語処理系**: 字句解析・構文解析・意味解析

- **コード生成**: アセンブリ生成・機械語変換

- **最適化理論**: データフロー解析・制御フロー解析

- **型システム**: 静的型検査・型推論・多態型

  

### **システムプログラミング**

- **低レベルプログラミング**: アセンブリ・機械語理解

- **メモリ管理**: ヒープ・スタック・ガベージコレクション

- **並行プログラミング**: スレッド・非同期処理

- **プラットフォーム依存**: CPU アーキテクチャ・OS API

  

### **ソフトウェア設計**

- **大規模システム設計**: モジュール化・インターフェース設計

- **性能最適化**: プロファイリング・ボトルネック解析

- **テスト技法**: 単体テスト・統合テスト・性能テスト

- **DevOps**: CI/CD・自動化・品質管理

  

## **🚀 今後の開発ロードマップ予定**

  

### **Phase 1: 基本言語処理（3-6ヶ月）**

- ✅ 字句解析器実装（トークナイザー）

- ✅ 再帰下降構文解析器（AST生成）

- ✅ シンプルインタープリター（直接実行）

- ✅ 基本的な型システム（int, float, string, bool）

  

### **Phase 2: コンパイラ基盤（6-12ヶ月）**

- ✅ アセンブリコード生成（x86_64）

- ✅ シンボルテーブル・スコープ管理

- ✅ 制御構造（if/while/for）

- ✅ 関数定義・呼び出し

  

### **Phase 3: 最適化エンジン（12-18ヶ月）**

- ✅ SSA中間表現（Static Single Assignment）

- ✅ 基本最適化パス（定数畳み込み、デッドコード除去）

- ✅ 高度最適化（ループ最適化、インライン展開）

- ✅ レジスタ割り当て（グラフ彩色法）

  

### **Phase 4: 産業級コンパイラ（18-24ヶ月）**

- ✅ LLVM IR生成・連携

- ✅ マルチターゲット対応（ARM64, WASM）

- ✅ 高度な型システム（ジェネリクス、トレイト）

- ✅ デバッグ・プロファイリングツール

  

## **🤖 AI開発支援向け設計予定**

  

このプロジェクトはAI開発支援ツールとの協調を前提に設計予定:

  

- **段階的進化**: 複雑さを段階的に導入、理解しやすい構造

- **包括的ドキュメント**: 各Phaseの詳細学習ガイド

- **自動化**: Makefile、CI/CDによる一貫した開発体験

- **テスタビリティ**: 包括的テストによる安全な変更

- **可視化**: AST・IR・最適化過程の可視化ツール

  

## **📝 貢献・フィードバック予定**

  

### **開発状況**

- **GitHub Issues**: [プロジェクトボード](https://github.com/nyasuto/pug/issues)

- **Performance Wiki**: [コンパイラ性能分析レポート](https://github.com/nyasuto/pug/wiki/Performance-Analysis-Report)

- **Learning Guide**: [段階的学習ガイド](https://github.com/nyasuto/pug/wiki/Learning-Guide)

- **Pull Requests**: コードレビュー歓迎

- **Discussions**: アイデア・提案の議論

  

### **ライセンス**

MITライセンス。学習目的のため、フォーク・改変・提案すべて歓迎です。

  

---

  

**🎉 学習からプロダクションまでのコンパイラ技術完全習得！**

シンプルなレクサーから**LLVM級の最適化コンパイラ**まで進化を遂げる学習プロジェクトです。段階的実装により**コンパイラ技術の全体像**を体系的に理解し、最終的には**産業級の性能と機能**を持つ本格的なコンパイラの開発スキルを習得できます！