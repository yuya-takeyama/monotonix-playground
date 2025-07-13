# Monotonix依存関係の粒度に関する問題と提案

## 問題の概要

現在のMonotonixの依存関係指定では、ディレクトリ単位での指定しかできないため、以下のような課題があります：

### 具体的なケース

```yaml
# web-app/cmd/api-server/monotonix.yaml
app:
  name: web-app/cmd/api-server
  depends_on:
    - web-app/pkg  # 共有ライブラリへの依存
```

この設定では：
- ✅ `web-app/pkg/` 内のファイル変更時に api-server がビルドされる
- ❌ `web-app/go.mod` や `web-app/go.sum` 変更時に api-server がビルドされない

### 現状の回避策とトレードオフ

1. **`depends_on: [web-app]` にする**
   - ✅ go.mod/go.sum変更時にビルドされる
   - ❌ api-server変更時に worker もビルドされる（不要なビルド）

2. **`depends_on: [web-app/pkg]` のまま**
   - ✅ 細かい依存関係管理
   - ❌ go.mod/go.sum変更時にビルドされない

## 提案：ファイル単位の依存関係指定

### 理想的な設定例

```yaml
app:
  name: web-app/cmd/api-server
  depends_on:
    - web-app/pkg/          # ディレクトリ
    - web-app/go.mod        # ファイル
    - web-app/go.sum        # ファイル
```

または

```yaml
app:
  name: web-app/cmd/api-server
  depends_on:
    directories:
      - web-app/pkg/
    files:
      - web-app/go.mod
      - web-app/go.sum
```

### 想定される技術的課題と副作用

1. **パフォーマンスへの影響**
   - ファイル単位での変更検知により、チェック対象が増加
   - 大量のファイル指定時のスキャン時間増加

2. **設定の複雑化**
   - 依存関係の記述が冗長になる可能性
   - 保守性の低下リスク

3. **ファイルパスの管理**
   - 相対パス vs 絶対パス
   - ファイル移動時の設定更新
   - glob パターンサポートの必要性

4. **変更検知の実装複雑化**
   - 現在のディレクトリベースの検知ロジックの拡張
   - ファイル削除/追加の処理

### 代替案

1. **パターンマッチング対応**
   ```yaml
   depends_on:
     - web-app/pkg/
     - web-app/*.{mod,sum}
   ```

2. **プリセット機能**
   ```yaml
   depends_on:
     - web-app/pkg/
     - preset:go-dependencies  # go.mod, go.sum を自動含有
   ```

3. **除外指定との組み合わせ**
   ```yaml
   depends_on:
     - web-app/
   excludes:
     - web-app/cmd/worker/  # worker変更時は除外
   ```

## モノレポでの実際のニーズ

- Go: go.mod, go.sum
- Node.js: package.json, package-lock.json, yarn.lock
- Python: requirements.txt, pyproject.toml, poetry.lock
- Rust: Cargo.toml, Cargo.lock

これらの依存関係ファイルの変更は、関連するすべてのアプリケーションに影響するため、ファイル単位での指定は実用的なニーズがあります。

## 結論

ファイル単位での依存関係指定は、モノレポにおける細かい依存関係管理において重要な機能だと考えられます。実装時はパフォーマンスと設定の複雑さのバランスを考慮する必要があります。