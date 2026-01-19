# Hippocampus モジュール使用例

## 概要

Hippocampus（海馬）モジュールは、**短期記憶(STM)** と **長期記憶(LTM)** を管理します。
人間の記憶システムを模倣し、重要な記憶は長期記憶へ固定化され、重要でない記憶は忘却されます。

## 基本的な使用方法

```go
package main

import (
    "fmt"
    "github.com/your-username/mind-os/internal/hippocampus"
    "github.com/your-username/mind-os/internal/models"
)

func main() {
    // Hippocampusインスタンスを作成
    h := hippocampus.New()
    
    // 記憶を追加
    emotions := []models.EmotionValue{
        {Code: models.EmotionJoy, Value: 80},
        {Code: models.EmotionHope, Value: 60},
    }
    h.AddEpisode("プロジェクトが成功した", emotions)
    
    // 記憶統計を確認
    fmt.Printf("STM: %d件, LTM: %d件\n", h.GetSTMCount(), h.GetLTMCount())
    
    // 睡眠処理（記憶の固定化）
    h.SleepAndConsolidate()
    
    fmt.Printf("睡眠後 - STM: %d件, LTM: %d件\n", h.GetSTMCount(), h.GetLTMCount())
}
```

## 記憶の種類

### 短期記憶 (STM - Short-Term Memory)

- 新しく追加された記憶は最初STMに保存される
- 最大100件まで保持（設定可能）
- 上限を超えると古い記憶から削除（FIFO）
- 睡眠処理で重要な記憶はLTMへ移行、それ以外は忘却

### 長期記憶 (LTM - Long-Term Memory)

- 重み（Weight）が閾値以上の記憶がSTMから移行
- デフォルト閾値: 0.6
- 最大1000件まで保持（設定可能）
- 重みとアクセス時刻でスコアリングし、低スコアの記憶から削除

## 記憶の重み計算

記憶の重要度は以下の要素で計算されます:

1. **感情の強度**: 感情値の平均（0-100を0.0-1.0に正規化）
2. **感情の多様性**: 感情の種類数 × 0.05

```go
// 例1: 高重み記憶（LTMへ移行される）
emotions := []models.EmotionValue{
    {Code: models.EmotionJoy, Value: 90},
    {Code: models.EmotionLove, Value: 85},
}
// 重み = (90+85)/200 + 2*0.05 = 0.875 + 0.1 = 0.975

// 例2: 低重み記憶（忘却される）
emotions := []models.EmotionValue{
    {Code: models.EmotionNeutral, Value: 30},
}
// 重み = 30/100 + 1*0.05 = 0.3 + 0.05 = 0.35
```

## 主要メソッド

### AddEpisode

新しいエピソード記憶をSTMに追加します。

```go
h := hippocampus.New()

emotions := []models.EmotionValue{
    {Code: models.EmotionJoy, Value: 75},
}

h.AddEpisode("新しい記憶", emotions)
```

**処理内容:**
1. UUIDを自動生成
2. 感情から重みを計算
3. タグを自動抽出
4. STMに追加
5. STMサイズ制限チェック

### GetRecentContext

直近の記憶を最大10件取得します。

```go
recent := h.GetRecentContext()

for _, memory := range recent {
    fmt.Printf("%s: %s (重み: %.2f)\n", 
        memory.Type, memory.Text, memory.Weight)
}
```

**特徴:**
- STMとLTMを統合
- 最終アクセス時刻でソート（新しい順）
- 最大10件を返す

### SleepAndConsolidate

睡眠処理 - 記憶の固定化と忘却を実行します。

```go
fmt.Printf("睡眠前 - STM: %d, LTM: %d\n", 
    h.GetSTMCount(), h.GetLTMCount())

h.SleepAndConsolidate()

fmt.Printf("睡眠後 - STM: %d, LTM: %d\n", 
    h.GetSTMCount(), h.GetLTMCount())
```

**処理フロー:**
1. STM内の各記憶の重みをチェック
2. 重み ≥ 0.6 → LTMへ移行
3. 重み < 0.6 → 忘却
4. STMをクリア
5. LTMサイズ制限チェック

## 実践例

### 例1: 日常会話の記憶

```go
h := hippocampus.New()

// 挨拶（低重み）
h.AddEpisode("おはよう", []models.EmotionValue{
    {Code: models.EmotionNeutral, Value: 50},
})

// 重要な出来事（高重み）
h.AddEpisode("プロジェクトが完成した！", []models.EmotionValue{
    {Code: models.EmotionJoy, Value: 95},
    {Code: models.EmotionHope, Value: 80},
})

// 睡眠処理
h.SleepAndConsolidate()

// 結果: 挨拶は忘却、プロジェクト完成はLTMへ
```

### 例2: 感情的な体験

```go
h := hippocampus.New()

// ネガティブな体験（高重み）
h.AddEpisode("バグで本番障害が発生", []models.EmotionValue{
    {Code: models.EmotionFear, Value: 90},
    {Code: models.EmotionAnger, Value: 75},
})

// ポジティブな解決（高重み）
h.AddEpisode("バグを修正できた", []models.EmotionValue{
    {Code: models.EmotionJoy, Value: 85},
    {Code: models.EmotionHope, Value: 70},
})

h.SleepAndConsolidate()

// 両方ともLTMへ移行（感情が強い）
```

### 例3: 記憶の検索

```go
h := hippocampus.New()

h.AddEpisode("重要な会議", []models.EmotionValue{
    {Code: models.EmotionJoy, Value: 70},
})

uuid := h.STM[0].UUID

// UUIDで記憶を検索
memory := h.GetMemoryByUUID(uuid)
if memory != nil {
    fmt.Printf("見つかった: %s\n", memory.Text)
}
```

## API経由での使用

### 記憶の追加

```bash
curl -X POST http://localhost:8081/api/memory/add \
  -H "Content-Type: application/json" \
  -d '{"text":"今日は良い一日だった"}'
```

### 直近の記憶取得

```bash
curl http://localhost:8081/api/memory/recent
```

### 記憶統計

```bash
curl http://localhost:8081/api/memory/stats
```

### 睡眠処理

```bash
curl -X POST http://localhost:8081/api/memory/sleep
```

## 設定のカスタマイズ

```go
h := hippocampus.New()

// 固定化閾値を変更（デフォルト: 0.6）
h.consolidationThreshold = 0.7 // より厳しく

// STM最大サイズを変更（デフォルト: 100）
h.maxSTMSize = 50

// LTM最大サイズを変更（デフォルト: 1000）
h.maxLTMSize = 500
```

## タグシステム

記憶には自動的にタグが付与されます:

1. **感情タグ**: 感情コード（J, S, A, F, L, D, H, N）
2. **長さタグ**:
   - `short`: 20文字未満
   - `long`: 100文字以上

```go
memory := h.STM[0]
fmt.Println("タグ:", memory.Tags)
// 出力例: [J H short]
```

## 記憶の削除アルゴリズム

LTMが最大サイズを超えた場合:

1. 各記憶のスコアを計算
   - スコア = 重み × 0.7 + 時間スコア × 0.3
   - 時間スコア = 1 / (1 + 経過日数)

2. スコアでソート（低い順）

3. 下位を削除（最大サイズの80%まで削減）

## 今後の拡張予定

- [ ] 記憶の関連性スコア計算
- [ ] セマンティック検索
- [ ] 記憶のクラスタリング
- [ ] 感情による記憶の想起
- [ ] データベースへの永続化
- [ ] 記憶の編集・削除API
