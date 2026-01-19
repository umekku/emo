# Store Module (Data Persistence)

Storeモジュールは、Mind OSのデータベース永続化層を提供します。
SQLiteを使用し、記憶(Memory)などのデータを永続化します。

## 機能概要

- **SQLite接続**: Pure Goドライバー (`modernc.org/sqlite`) を使用し、CGO不要で動作
- **スキーマ管理**: 起動時にテーブルとインデックスを自動生成
- **CRUD操作**: 記憶の保存、取得、検索、削除機能を提供

## データベーススキーマ

### memories テーブル

| カラム名 | 型 | 説明 |
|----------|----|------|
| `uuid` | TEXT (PK) | 記憶のユニークID |
| `text` | TEXT | 記憶の内容 |
| `emotions` | TEXT (JSON) | 感情値リスト |
| `weight` | REAL | 記憶の重要度 (0.0-1.0) |
| `type` | TEXT | 記憶の種類 (STM/LTM) |
| `created_at` | DATETIME | 作成日時 |
| `last_access` | DATETIME | 最終アクセス日時 |
| `tags` | TEXT (JSON) | タグリスト |

## 使用方法

### 初期化

`core` パッケージで初期化され、`Hippocampus` に注入されます。

```go
db, err := store.NewDB("mind.db")
if err != nil {
    // エラー処理
}
defer db.Close()
```

### 記憶の保存

```go
err := db.SaveMemory(memory)
```

### 記憶の取得

```go
// 直近の記憶を取得
memories, err := db.GetRecentMemories(10)

// UUIDで検索
memory, err := db.GetMemoryByUUID("uuid-string")
```

### 古い記憶の削除

LTMのサイズ制限のため、重要度が低く古い記憶を削除します。

```go
// LTMを1000件に制限（超過分を削除）
err := db.DeleteOldMemories(1000)
```
