# Git Hooks - pug コンパイラプロジェクト

ブランチルールを強制し、コード品質を保証するGit hooksです。

## 🪝 導入されているフック


### pre-push  
- **mainブランチへの直接プッシュ禁止**
- ブランチ命名規則の強制
- mainブランチからの派生チェック
- CI事前チェック（品質、ベンチマーク）
- Issue番号参照の推奨


## 🚀 セットアップ

```bash
# Git フックを有効化
make git-hooks

# テスト実行
make git-hooks-test

# 無効化（必要な場合）
make git-hooks-disable
```

## 📋 ブランチ命名規則

### 必須パターン
- `feat/issue-X-feature-name` - 新機能
- `fix/issue-X-description` - バグ修正
- `hotfix/X-description` - 緊急修正
- `test/X-description` - テスト
- `docs/X-description` - ドキュメント
- `ci/X-description` - CI/CD
- `refactor/X-description` - リファクタリング
- `perf/X-description` - 性能改善
- `security/X-description` - セキュリティ
- `deps/X-description` - 依存関係

### 例
```bash
git checkout -b feat/issue-2-lexer-implementation
git checkout -b fix/issue-5-parser-error
git checkout -b docs/issue-10-readme-update
```

### 利用可能なtype
- `feat` - 新機能
- `fix` - バグ修正
- `docs` - ドキュメント
- `style` - コードスタイル
- `refactor` - リファクタリング
- `test` - テスト追加・修正
- `chore` - その他のタスク
- `ci` - CI/CD関連
- `perf` - 性能改善
- `security` - セキュリティ修正
- `deps` - 依存関係更新

### 例
```bash
git commit -m "feat: Phase 1.0 レクサー実装 (#2)"
git commit -m "fix: 型検査エラーの修正 (#7)"
git commit -m "docs: README の使用方法更新"
```

## 🔧 トラブルシューティング

### フックをバイパスしたい場合
```bash
# 一時的にスキップ（非推奨）
git commit --no-verify
git push --no-verify
```

### フックが失敗する場合
1. エラーメッセージを確認
2. 品質チェックを実行: `make quality`
3. 修正後に再度コミット

### mainブランチでの作業が必要な場合
```bash
# フックを一時的に無効化
make git-hooks-disable

# 作業実行

# フックを再有効化
make git-hooks
```

## ⚠️ 重要な注意事項

- **mainブランチへの直接変更は禁止**
- 必ずPull Requestワークフローを使用
- コミット前に `make quality` を実行推奨
- 機密情報（API キー等）は絶対にコミットしない

これらのルールにより、コードベースの品質とセキュリティが保証されます。