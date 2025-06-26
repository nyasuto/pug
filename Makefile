# Makefile for pug compiler project
# コンパイラ学習プロジェクト「pug」のビルドシステム

# Go settings
GO=go
GOCMD=$(GO)
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# プロジェクト設定
PROJECT_NAME=pug
BINARY_NAME=pug
BINARY_INTERP=interp
BINARY_TOOLS=tools

# ディレクトリ
BIN_DIR=bin
BUILD_DIR=build
DIST_DIR=dist

# Phase別ターゲット
PHASE1_SOURCES=phase1/*.go
PHASE2_SOURCES=phase2/*.go
PHASE3_SOURCES=phase3/**/*.go
PHASE4_SOURCES=phase4/**/*.go

.DEFAULT_GOAL := help

##@ ヘルプ
.PHONY: help
help: ## 利用可能コマンド一覧を表示
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1m%s\033[0m\n", "pug コンパイラ開発コマンド"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ 開発環境
.PHONY: dev
dev: deps dirs ## 開発環境セットアップ
	@echo "🚀 開発環境セットアップ完了"

.PHONY: deps
deps: ## 依存関係のインストール
	@echo "📦 依存関係をインストール中..."
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: dirs
dirs: ## プロジェクト構造の作成
	@echo "📁 ディレクトリ構造を作成中..."
	@mkdir -p $(BIN_DIR) $(BUILD_DIR) $(DIST_DIR)
	@mkdir -p phase1 phase2 phase3/ir phase3/optimizer phase3/backend
	@mkdir -p phase4/llvm phase4/runtime phase4/tools
	@mkdir -p cmd/μc cmd/μinterp cmd/μtools
	@mkdir -p examples benchmark docs

##@ Phase 1: 基本言語処理
.PHONY: phase1-build
phase1-build: ## Phase 1 インタープリターをビルド
	@echo "🔤 Phase 1 インタープリターをビルド中..."
	@if [ -f "cmd/interp/main.go" ]; then \
		$(GOBUILD) -o $(BIN_DIR)/$(BINARY_INTERP) cmd/interp/main.go; \
	else \
		echo "⚠️  cmd/interp/main.go が見つかりません"; \
	fi

.PHONY: phase1-test
phase1-test: ## Phase 1 テスト実行
	@echo "🧪 Phase 1 テスト実行中..."
	@if [ -d "phase1" ]; then \
		$(GOTEST) -v ./phase1/...; \
	else \
		echo "⚠️  phase1/ ディレクトリが見つかりません"; \
	fi

.PHONY: phase1-bench
phase1-bench: ## Phase 1 ベンチマーク実行
	@echo "⚡ Phase 1 ベンチマーク実行中..."
	@if [ -d "phase1" ]; then \
		$(GOTEST) -bench=. -benchmem ./phase1/...; \
	else \
		echo "⚠️  phase1/ ディレクトリが見つかりません"; \
	fi

##@ Phase 2: コンパイラ基盤
.PHONY: phase2-build
phase2-build: ## Phase 2 コンパイラをビルド
	@echo "⚙️ Phase 2 コンパイラをビルド中..."
	@if [ -f "cmd/compiler/main.go" ]; then \
		$(GOBUILD) -o $(BIN_DIR)/pugc cmd/compiler/main.go; \
		echo "✅ pugc コンパイラをビルドしました"; \
	elif [ -f "cmd/pug/main.go" ]; then \
		$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) cmd/pug/main.go; \
	else \
		echo "⚠️  コンパイラのmain.goが見つかりません"; \
	fi

.PHONY: phase2-test
phase2-test: ## Phase 2 テスト実行
	@echo "🧪 Phase 2 テスト実行中..."
	@if [ -d "phase2" ]; then \
		$(GOTEST) -v ./phase2/...; \
	else \
		echo "⚠️  phase2/ ディレクトリが見つかりません"; \
	fi

##@ テスト・品質チェック
.PHONY: test
test: ## 全てのテストを実行
	@echo "🧪 全テスト実行中..."
	$(GOTEST) -v ./...

.PHONY: test-cov
test-cov: ## テストカバレッジ付きでテスト実行
	@echo "📊 テストカバレッジ測定中..."
	@echo "mode: atomic" > coverage.out
	$(GOTEST) -race -coverprofile=phase1.cover -covermode=atomic ./phase1
	@tail -n +2 phase1.cover >> coverage.out || true
	@if [ -d "./phase2" ]; then \
		$(GOTEST) -race -coverprofile=phase2.cover -covermode=atomic ./phase2; \
		tail -n +2 phase2.cover >> coverage.out || true; \
	fi
	@if [ -d "./phase3" ]; then \
		$(GOTEST) -race -coverprofile=phase3.cover -covermode=atomic ./phase3; \
		tail -n +2 phase3.cover >> coverage.out || true; \
	fi
	@if [ -d "./phase4" ]; then \
		$(GOTEST) -race -coverprofile=phase4.cover -covermode=atomic ./phase4; \
		tail -n +2 phase4.cover >> coverage.out || true; \
	fi
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ カバレッジファイル生成完了:"
	@echo "   - coverage.out (統合)"
	@echo "   - coverage.html (HTMLレポート)"
	@ls -la phase*.cover 2>/dev/null | sed 's/^/   - /' || true

.PHONY: lint
lint: ## コードの静的解析
	@echo "🔍 静的解析実行中..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run ./...; \
	else \
		echo "⚠️  golangci-lint がインストールされていません"; \
		echo "   インストール: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: fmt
fmt: ## コードフォーマット
	@echo "✨ コードフォーマット中..."
	$(GOCMD) fmt ./...

.PHONY: sec
sec: ## セキュリティスキャン (gosec)
	@echo "🔒 セキュリティスキャン実行中..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec がインストールされていません"; \
		echo "   インストール: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
	fi

.PHONY: quality
quality: deps fmt lint build test sec ## 全品質チェック実行
	@echo "✅ 品質チェック完了"

.PHONY: quality-fix
quality-fix: ## 自動修正可能な品質問題を修正
	@echo "🔧 品質問題自動修正中..."
	$(GOCMD) fmt ./...
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run --fix ./...; \
	else \
		echo "⚠️  golangci-lint がインストールされていません"; \
	fi
	@echo "✅ 自動修正完了 - 手動でテストを実行してください: make test"

##@ ベンチマーク・性能測定
.PHONY: bench
bench: ## 全ベンチマーク実行
	@echo "⚡ ベンチマーク実行中..."
	$(GOTEST) -bench=. -benchmem ./...

.PHONY: bench-compiler
bench-compiler: ## コンパイラ性能ベンチマーク
	@echo "🔧 コンパイラベンチマーク実行中..."
	@if [ -d "benchmark" ]; then \
		$(GOTEST) -bench=BenchmarkCompiler -benchmem -v ./benchmark/...; \
	else \
		echo "⚠️  benchmark/ ディレクトリが見つかりません"; \
	fi

.PHONY: bench-gcc
bench-gcc: ## GCC比較ベンチマーク
	@echo "🏁 GCC比較ベンチマーク実行中..."
	@if command -v gcc >/dev/null 2>&1; then \
		if [ -d "benchmark" ]; then \
			$(GOTEST) -bench=BenchmarkVsGCC -benchmem -v -timeout=10m ./benchmark/...; \
		else \
			echo "⚠️  benchmark/ ディレクトリが見つかりません"; \
		fi \
	else \
		echo "⚠️  GCCが見つかりません。GCCをインストールしてください"; \
	fi

.PHONY: bench-rust
bench-rust: ## Rust比較ベンチマーク
	@echo "🦀 Rust比較ベンチマーク実行中..."
	@if command -v cargo >/dev/null 2>&1; then \
		if [ -d "benchmark" ]; then \
			$(GOTEST) -bench=BenchmarkVsRust -benchmem -v -timeout=15m ./benchmark/...; \
		else \
			echo "⚠️  benchmark/ ディレクトリが見つかりません"; \
		fi \
	else \
		echo "⚠️  Rust/Cargoが見つかりません。Rustをインストールしてください"; \
	fi

.PHONY: bench-suite
bench-suite: ## 包括的ベンチマークスイート
	@echo "📊 包括的ベンチマークスイート実行中..."
	@if [ -d "benchmark" ]; then \
		$(GOTEST) -bench=BenchmarkSuite -benchmem -v -timeout=20m ./benchmark/...; \
	else \
		echo "⚠️  benchmark/ ディレクトリが見つかりません"; \
	fi

.PHONY: bench-evolution
bench-evolution: ## 進化分析ベンチマーク
	@echo "📈 進化分析ベンチマーク実行中..."
	@if [ -d "benchmark" ]; then \
		$(GOTEST) -bench=BenchmarkEvolution -benchmem -v -timeout=15m ./benchmark/...; \
	else \
		echo "⚠️  benchmark/ ディレクトリが見つかりません"; \
	fi

.PHONY: bench-report
bench-report: ## ベンチマークレポート生成
	@echo "📋 ベンチマークレポート生成中..."
	@if [ -d "benchmark" ]; then \
		$(GOTEST) -run=TestBenchmarkReport -v ./benchmark/...; \
		echo "✅ レポートが生成されました"; \
	else \
		echo "⚠️  benchmark/ ディレクトリが見つかりません"; \
	fi

.PHONY: bench-comprehensive
bench-comprehensive: build bench-compiler bench-gcc bench-rust bench-evolution ## 包括的ベンチマーク（全機能）
	@echo "🎯 包括的ベンチマーク完了"
	@echo ""
	@echo "📊 実行されたベンチマーク:"
	@echo "  ✅ コンパイラベンチマーク"
	@if command -v gcc >/dev/null 2>&1; then \
		echo "  ✅ GCC比較ベンチマーク"; \
	else \
		echo "  ⚠️  GCC比較ベンチマーク (スキップ)"; \
	fi
	@if command -v cargo >/dev/null 2>&1; then \
		echo "  ✅ Rust比較ベンチマーク"; \
	else \
		echo "  ⚠️  Rust比較ベンチマーク (スキップ)"; \
	fi
	@echo "  ✅ 進化分析ベンチマーク"
	@echo ""
	@echo "📈 詳細結果は個別ベンチマーク実行時に表示されます"

.PHONY: bench-analyze
bench-analyze: ## 性能データ分析・レポート生成
	@echo "📊 性能データ分析実行中..."
	@if [ -f "scripts/performance/cmd/analyzer/main.go" ]; then \
		go run scripts/performance/cmd/analyzer/main.go; \
	else \
		echo "⚠️  analyzer が見つかりません"; \
	fi

.PHONY: bench-trend
bench-trend: ## 長期性能トレンド分析
	@echo "📈 トレンド分析実行中..."
	@if [ -f "scripts/performance/cmd/trend/main.go" ]; then \
		go run scripts/performance/cmd/trend/main.go; \
	else \
		echo "⚠️  trend analyzer が見つかりません"; \
	fi

.PHONY: bench-dashboard
bench-dashboard: ## 性能ダッシュボード生成
	@echo "🎨 ダッシュボード生成中..."
	@if [ -f "scripts/performance/cmd/dashboard/main.go" ]; then \
		go run scripts/performance/cmd/dashboard/main.go; \
	else \
		echo "⚠️  dashboard generator が見つかりません"; \
	fi

.PHONY: bench-wiki
bench-wiki: ## GitHub Wiki自動更新
	@echo "📝 Wiki自動更新中..."
	@if [ -f "scripts/performance/cmd/wiki/main.go" ]; then \
		go run scripts/performance/cmd/wiki/main.go; \
	else \
		echo "⚠️  wiki updater が見つかりません"; \
	fi

.PHONY: bench-cicd
bench-cicd: bench-comprehensive bench-analyze bench-trend bench-dashboard ## CI/CD統合ベンチマーク（完全自動化）
	@echo "🚀 CI/CD統合ベンチマーク完了"
	@echo ""
	@echo "📊 生成されたファイル:"
	@ls -la performance-*.* 2>/dev/null | sed 's/^/  /' || echo "  📊 レポートファイルなし"
	@ls -la trend-*.* 2>/dev/null | sed 's/^/  /' || echo "  📈 トレンドファイルなし"  
	@ls -la dashboard-*.* 2>/dev/null | sed 's/^/  /' || echo "  🎨 ダッシュボードファイルなし"
	@echo ""
	@echo "✅ 全ての性能分析が完了しました"

.PHONY: bench-compile
bench-compile: ## コンパイル時間測定（従来互換）
	@echo "⏱️  コンパイル時間測定中..."
	@make bench-compiler

.PHONY: bench-runtime
bench-runtime: ## 実行時間測定（従来互換）
	@echo "🏃 実行時間測定中..."
	@make bench-compiler

##@ ビルド・リリース
.PHONY: build
build: phase1-build phase2-build ## 全バイナリをビルド
	@echo "🔨 ビルド完了"

.PHONY: clean
clean: ## ビルド成果物を削除
	@echo "🧹 クリーンアップ中..."
	@rm -rf $(BIN_DIR) $(BUILD_DIR) $(DIST_DIR)
	@rm -f coverage.out coverage.html phase*.cover

.PHONY: install
install: build ## システムにインストール
	@echo "📦 システムインストール中..."
	@cp $(BIN_DIR)/* $(GOPATH)/bin/ 2>/dev/null || echo "⚠️  GOPATH/bin にコピーできませんでした"

##@ ユーティリティ
.PHONY: env-info
env-info: ## 環境情報表示
	@echo "🔧 環境情報:"
	@echo "  Go version: $$($(GOCMD) version)"
	@echo "  Go path: $$($(GOCMD) env GOPATH)"
	@echo "  Go root: $$($(GOCMD) env GOROOT)"
	@echo "  Project: $(PROJECT_NAME)"
	@echo "  Module: $$(head -1 go.mod | cut -d' ' -f2)"

.PHONY: demo
demo: phase1-build ## デモンストレーション実行
	@echo "🎭 デモンストレーション:"
	@if [ -f "examples/hello.dog" ]; then \
		echo "📄 Hello World プログラム実行:"; \
		./$(BIN_DIR)/$(BINARY_INTERP) examples/hello.dog; \
	else \
		echo "⚠️  examples/hello.dog が見つかりません"; \
	fi

##@ Git・CI/CD
.PHONY: pr
pr: ## Pull Request作成（現在のブランチから）

.PHONY: pr-ready  
pr-ready: ## PR準備完了チェック

.PHONY: git-hooks
git-hooks: ## Git フック（ブランチルール強制）をセットアップ
	@echo "🪝 Git フックセットアップ中..."
	@if [ -d ".git" ]; then \
		echo "📋 フックファイルをコピー中..."; \
		cp .git-hooks/pre-commit .git/hooks/pre-commit; \
		cp .git-hooks/pre-push .git/hooks/pre-push; \
		cp .git-hooks/commit-msg .git/hooks/commit-msg; \
		chmod +x .git/hooks/pre-commit .git/hooks/pre-push .git/hooks/commit-msg; \
		echo "✅ Git フックが設定されました:"; \
		echo "  - pre-commit: 品質チェック + ブランチルール強制"; \
		echo "  - pre-push: CI事前チェック（品質・ベンチマーク）"; \
		echo "  - commit-msg: Conventional Commits形式強制"; \
	else \
		echo "⚠️  Git リポジトリが見つかりません"; \
	fi

.PHONY: git-hooks-disable
git-hooks-disable: ## Git フックを無効化
	@echo "🚫 Git フック無効化中..."
	@if [ -d ".git/hooks" ]; then \
		rm -f .git/hooks/pre-commit .git/hooks/pre-push .git/hooks/commit-msg; \
		echo "✅ Git フックが無効化されました"; \
	else \
		echo "⚠️  Git フックディレクトリが見つかりません"; \
	fi

.PHONY: git-hooks-test
git-hooks-test: ## Git フックのテスト実行
	@echo "🧪 Git フックテスト実行中..."
	@if [ -f ".git/hooks/pre-commit" ]; then \
		echo "📋 pre-commit フックテスト:"; \
		.git/hooks/pre-commit || echo "❌ pre-commit フック失敗"; \
	fi
	@if [ -f ".git/hooks/commit-msg" ]; then \
		echo "📝 commit-msg フックテスト:"; \
		echo "test: テストコミットメッセージ" > /tmp/test-commit-msg; \
		.git/hooks/commit-msg /tmp/test-commit-msg || echo "❌ commit-msg フック失敗"; \
		rm -f /tmp/test-commit-msg; \
	fi

pr: ## Pull Request作成（現在のブランチから）
	@echo "🔀 Pull Request作成中..."
	@current_branch=$$(git symbolic-ref --short HEAD 2>/dev/null || echo ""); \
	if [ "$$current_branch" = "main" ]; then \
		echo "❌ エラー: mainブランチからはPull Requestを作成できません"; \
		echo "   機能ブランチを作成してください: git checkout -b feat/your-feature"; \
		exit 1; \
	fi; \
	if [ -z "$$current_branch" ]; then \
		echo "❌ エラー: ブランチ名を取得できません"; \
		exit 1; \
	fi; \
	echo "📤 ブランチ: $$current_branch"; \
	if ! git diff --quiet; then \
		echo "⚠️  警告: コミットされていない変更があります"; \
		echo "   先にコミットしてください: git add . && git commit"; \
		exit 1; \
	fi; \
	echo "🔍 最終品質チェック実行中..."; \
	if ! make quality >/dev/null 2>&1; then \
		echo "❌ エラー: 品質チェックに失敗しました"; \
		echo "   修正してください: make quality"; \
		exit 1; \
	fi; \
	echo "📤 ブランチをプッシュ中..."; \
	git push -u origin "$$current_branch"; \
	echo "🔀 Pull Request作成中..."; \
	pr_title=$$(echo "$$current_branch" | sed 's/^[^/]*\///'); \
	issue_num=$$(echo "$$current_branch" | grep -oE 'issue-[0-9]+' | grep -oE '[0-9]+' || echo ""); \
	if [ -n "$$issue_num" ]; then \
		pr_body="Closes #$$issue_num"; \
	else \
		pr_body=""; \
	fi; \
	gh pr create --title "$$pr_title" --body "## Summary\n$$pr_title\n\n$$pr_body\n\n## Test plan\n- [ ] 品質チェック通過\n- [ ] テストケース実行・通過\n- [ ] 機能動作確認\n- [ ] ドキュメント更新（必要に応じて）\n\n🤖 Generated with [Claude Code](https://claude.ai/code)"; \
	echo "✅ Pull Request作成完了"

pr-ready: ## PR準備完了チェック
	@echo "🔍 PR準備完了チェック実行中..."
	@current_branch=$$(git symbolic-ref --short HEAD 2>/dev/null || echo ""); \
	if [ "$$current_branch" = "main" ]; then \
		echo "❌ エラー: mainブランチです"; \
		exit 1; \
	fi; \
	echo "📋 ブランチ: $$current_branch"; \
	if ! git diff --quiet; then \
		echo "❌ エラー: コミットされていない変更があります"; \
		exit 1; \
	fi; \
	if ! git diff --cached --quiet; then \
		echo "❌ エラー: ステージされた変更があります"; \
		exit 1; \
	fi; \
	echo "🔍 品質チェック..."; \
	make quality; \
	echo "🧪 テスト実行..."; \
	make test; \
	echo "⚡ ベンチマーク実行..."; \
	make bench >/dev/null 2>&1 || echo "⚠️  ベンチマーク失敗（継続）"; \
	echo "✅ PR準備完了 - make pr でPull Requestを作成できます"

# フォニーターゲットの設定
.PHONY: all
all: dev build test quality bench ## 全てのタスクを実行