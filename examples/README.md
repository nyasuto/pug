# Dog言語サンプルプログラム集

**🐶 pugプロジェクトの包括的なサンプルコード集**

このディレクトリには、Dog言語の機能を実践的に学ぶための豊富なサンプルプログラムが含まれています。初心者から上級者まで、段階的に言語機能を理解し、実用的なプログラミングスキルを身につけることができます。

## 📚 サンプルプログラム一覧

### 🎯 基本機能デモ

| ファイル名 | 内容 | 対象レベル | 学習時間 |
|-----------|------|----------|----------|
| **[hello.dog](hello.dog)** | Hello World・基本的な変数・関数 | 初級 | 15分 |
| **[data_types_demo.dog](data_types_demo.dog)** | 全データ型の詳細解説・操作例 | 初級 | 30分 |
| **[control_flow_demo.dog](control_flow_demo.dog)** | if文・while文・複雑な制御フロー | 初級-中級 | 45分 |
| **[functions_demo.dog](functions_demo.dog)** | 関数定義・クロージャ・高階関数 | 中級 | 60分 |

### 🧮 アルゴリズム実装集

| ファイル名 | 内容 | 対象レベル | 学習時間 |
|-----------|------|----------|----------|
| **[sorting_algorithms.dog](sorting_algorithms.dog)** | 7種類のソートアルゴリズム実装・比較 | 中級 | 90分 |
| **[search_algorithms.dog](search_algorithms.dog)** | 8種類の探索アルゴリズム・データ構造 | 中級-上級 | 75分 |
| **[mathematical_algorithms.dog](mathematical_algorithms.dog)** | 数論・組み合わせ論・数値計算 | 上級 | 120分 |

### 🛠️ 実用的アプリケーション

| ファイル名 | 内容 | 対象レベル | 学習時間 |
|-----------|------|----------|----------|
| **[calculator.dog](calculator.dog)** | 高機能電卓・数学関数実装 | 初級-中級 | 45分 |
| **[text_processor.dog](text_processor.dog)** | テキスト解析・パターンマッチング | 中級 | 60分 |
| **[simple_game.dog](simple_game.dog)** | 6種類のゲーム実装 | 中級 | 90分 |

### ⚡ 性能測定・ベンチマーク

| ファイル名 | 内容 | 対象レベル | 学習時間 |
|-----------|------|----------|----------|
| **[fibonacci.dog](fibonacci.dog)** | フィボナッチ・クロージャデモ | 初級-中級 | 30分 |
| **[performance_benchmarks.dog](performance_benchmarks.dog)** | 包括的性能測定スイート | 上級 | 45分 |

### 🔬 既存のサンプル

| ファイル名 | 内容 | 対象レベル |
|-----------|------|----------|
| **[algorithms.dog](algorithms.dog)** | 基本アルゴリズム集 | 中級 |
| **[closures.dog](closures.dog)** | クロージャの実践例 | 中級 |
| **[sorting.dog](sorting.dog)** | ソートアルゴリズム | 中級 |

## 🗺️ 推奨学習パス

### 🚀 初心者コース（8-10時間）

```
1. hello.dog                    (15分) - Dog言語の基本を理解
2. data_types_demo.dog          (30分) - データ型をマスター
3. control_flow_demo.dog        (45分) - 制御構造を学習
4. calculator.dog               (45分) - 実用的な計算プログラム
5. fibonacci.dog                (30分) - 再帰とクロージャ
6. simple_game.dog              (90分) - ゲームプログラミング
```

### 🎯 中級者コース（12-15時間）

```
1. functions_demo.dog           (60分) - 高度な関数プログラミング
2. sorting_algorithms.dog       (90分) - アルゴリズムの実装と分析
3. search_algorithms.dog        (75分) - データ構造と探索
4. text_processor.dog           (60分) - 文字列処理とパターンマッチング
5. performance_benchmarks.dog   (45分) - 性能測定と最適化
```

### 🧠 上級者コース（15-20時間）

```
1. mathematical_algorithms.dog  (120分) - 高度な数学アルゴリズム
2. 全サンプルの拡張・カスタマイズ      - 独自機能の追加
3. 複数サンプルの組み合わせ           - 大規模アプリケーション開発
4. 性能最適化チャレンジ              - アルゴリズム改善
```

## 📈 Phase別対応状況

### Phase 1（インタープリター）対応

✅ **完全対応サンプル**
- hello.dog
- data_types_demo.dog  
- control_flow_demo.dog
- functions_demo.dog
- fibonacci.dog
- calculator.dog
- simple_game.dog

### Phase 2（コンパイラ）対応

🔄 **テスト済みサンプル**
- sorting_algorithms.dog
- search_algorithms.dog
- mathematical_algorithms.dog
- text_processor.dog
- performance_benchmarks.dog

### Phase 3-4（最適化・LLVM）

🚀 **性能測定用サンプル**
- performance_benchmarks.dog（各Phaseでの性能比較）
- 大規模データセット処理サンプル

## 🎮 実行方法

### Phase 1（インタープリター）で実行

```bash
# 基本的な実行
./bin/interp examples/hello.dog

# REPLモード
./bin/interp --repl

# 特定のサンプル実行
./bin/interp examples/data_types_demo.dog
./bin/interp examples/fibonacci.dog
```

### Phase 2（コンパイラ）で実行

```bash
# コンパイルして実行
./bin/pug examples/hello.dog -o hello
./hello

# アセンブリ出力確認
./bin/pug examples/fibonacci.dog --emit-asm

# デバッグ情報付きコンパイル
./bin/pug examples/calculator.dog --debug -o calculator
```

### 性能測定での活用

```bash
# ベンチマーク実行
make bench-interpreter    # Phase 1での測定
make bench-compiler      # Phase 2での測定

# 特定サンプルでの性能比較
time ./bin/interp examples/performance_benchmarks.dog
time ./bin/pug examples/performance_benchmarks.dog -o bench && ./bench
```

## 📊 サンプルの特徴と学習効果

### 🎯 教育的価値

| 特徴 | 効果 |
|------|------|
| **段階的複雑さ** | 基礎から応用まで無理なく学習 |
| **実用的内容** | 現実的な問題解決能力を育成 |
| **豊富なコメント** | 実装の理解を深める詳細解説 |
| **性能分析** | アルゴリズムの効率性を実感 |

### 🔬 技術的網羅性

| 分野 | カバー内容 |
|------|-----------|
| **言語機能** | 全データ型・制御構造・関数・クロージャ |
| **アルゴリズム** | ソート・探索・数学・グラフ理論 |
| **アプリケーション** | ゲーム・計算機・テキスト処理 |
| **性能最適化** | ベンチマーク・比較分析・改善手法 |

## 🛠️ カスタマイズとExtending

### サンプルの改造例

1. **calculator.dogの拡張**
   ```dog
   // 新しい数学関数を追加
   let sin_approximation = fn(x) { /* テイラー展開実装 */ };
   let cos_approximation = fn(x) { /* テイラー展開実装 */ };
   ```

2. **simple_game.dogのエンハンス**
   ```dog
   // セーブ・ロード機能追加
   let save_game_state = fn(state) { /* 状態保存 */ };
   let load_game_state = fn() { /* 状態復元 */ };
   ```

3. **performance_benchmarks.dogの拡張**
   ```dog
   // 新しいベンチマーク追加
   let matrix_multiplication_benchmark = fn(size) { /* 行列演算 */ };
   ```

### 独自サンプルの作成指針

1. **明確な学習目標**：何を学ぶかを明確に定義
2. **段階的な実装**：簡単な部分から複雑な部分へ
3. **豊富なコメント**：理解を助ける詳細な説明
4. **実行可能性**：実際に動作することを確認
5. **拡張可能性**：さらなる改良の余地を残す

## 🎓 学習成果の確認

### チェックリスト

#### 基本レベル
- [ ] Dog言語の基本文法を理解している
- [ ] 全データ型を使い分けできる
- [ ] if文・while文を適切に使用できる
- [ ] 関数定義と呼び出しができる

#### 中級レベル
- [ ] クロージャと高階関数を理解している
- [ ] 基本的なアルゴリズムを実装できる
- [ ] 配列・ハッシュを効果的に使用できる
- [ ] エラーハンドリングを適切に行える

#### 上級レベル
- [ ] 高度なアルゴリズムを設計・実装できる
- [ ] 性能を意識したコード最適化ができる
- [ ] 大規模なプログラムを構造化できる
- [ ] 他の言語との性能比較ができる

### 実践プロジェクト提案

1. **個人プロジェクト**
   - オリジナルゲームの開発
   - 数値計算ライブラリの作成
   - テキスト解析ツールの実装

2. **性能改善チャレンジ**
   - 既存アルゴリズムの最適化
   - 新しいデータ構造の実装
   - 並行処理の導入

3. **言語拡張プロジェクト**
   - 新しい組み込み関数の追加
   - ライブラリシステムの設計
   - デバッグ支援ツールの開発

## 🔗 関連リソース

### プロジェクト内ドキュメント
- **[メインREADME](../README.md)** - プロジェクト全体概要
- **[Phase 1学習ガイド](../docs/phase1_tutorial.md)** - インタープリター実装
- **[Phase 2学習ガイド](../docs/phase2_tutorial.md)** - コンパイラ実装
- **[コンパイラ理論解説](../docs/compiler_theory.md)** - 理論的背景

### 性能測定システム
- **[benchmark/](../benchmark/)** - 自動性能測定システム
- **[性能レポート](../performance_report.html)** - 最新の測定結果

### 外部学習リソース
- **アルゴリズム学習**: LeetCode, AtCoder, Codeforces
- **言語設計**: SICP, PLAI, Types and Programming Languages
- **コンパイラ技術**: Dragon Book, Modern Compiler Implementation

---

## 🎉 まとめ

この豊富なサンプルプログラム集を通じて、Dog言語の実用性と教育価値を最大限に活用してください。各サンプルは独立して学習できる設計になっていますが、段階的に進めることでより深い理解が得られます。

**Happy Coding with Dog Language! 🐶✨**

**質問や改善提案があれば、GitHubのIssuesでお気軽にお聞かせください。**