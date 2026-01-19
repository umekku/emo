# Basal Ganglia モジュール使用例

## 概要

Basal Ganglia（大脳基底核）モジュールは、**ドーパミンによる意欲(Motivation)** と **行動強化** を管理します。
外部フィードバックや感情に基づいて意欲レベルを動的に調整し、行動の選択と実行を制御します。

## 基本的な使用方法

```go
package main

import (
    "fmt"
    "github.com/your-username/mind-os/internal/basal"
)

func main() {
    // BasalGangliaインスタンスを作成
    bg := basal.New()
    
    fmt.Printf("初期意欲: %d\n", bg.GetMotivation())
    
    // ポジティブフィードバック（褒められた）
    bg.UpdateMotivation(true)
    fmt.Printf("褒められた後: %d (%s)\n", 
        bg.GetMotivation(), bg.GetMotivationLevel())
    
    // ネガティブフィードバック（無視された）
    bg.UpdateMotivation(false)
    fmt.Printf("無視された後: %d (%s)\n", 
        bg.GetMotivation(), bg.GetMotivationLevel())
}
```

## 意欲レベル

意欲は **0-100** の範囲で管理され、以下のレベルに分類されます:

| 意欲値 | レベル | 説明 |
|--------|--------|------|
| 80-100 | `very_high` | 非常に高い意欲 |
| 60-79  | `high` | 高い意欲 |
| 40-59  | `normal` | 通常の意欲 |
| 20-39  | `low` | 低い意欲 |
| 0-19   | `very_low` | 非常に低い意欲 |

**初期値**: 50 (normal)

## 主要メソッド

### UpdateMotivation

外部フィードバックにより意欲を更新します。

```go
bg := basal.New()

// ポジティブフィードバック: +10
bg.UpdateMotivation(true)

// ネガティブフィードバック: -15
bg.UpdateMotivation(false)
```

**パラメータ:**
- `isPositiveFeedback bool`: true = ポジティブ、false = ネガティブ

**効果:**
- ポジティブ: 意欲 +10
- ネガティブ: 意欲 -15 (減少の方が大きい)

### GetMotivation

現在の意欲レベルを取得します。

```go
motivation := bg.GetMotivation()
fmt.Printf("現在の意欲: %d\n", motivation)
```

**戻り値:** int (0-100)

### GetMotivationLevel

意欲レベルを文字列で取得します。

```go
level := bg.GetMotivationLevel()
fmt.Printf("意欲レベル: %s\n", level)
// 出力例: "high"
```

**戻り値:** string (`very_high`, `high`, `normal`, `low`, `very_low`)

### RewardFromEmotion

感情値から報酬を計算して意欲を更新します。

```go
// 高い感情値（喜び）
bg.RewardFromEmotion(90)  // 意欲上昇

// 低い感情値（恐れ）
bg.RewardFromEmotion(20)  // 意欲低下

// 中立
bg.RewardFromEmotion(50)  // 変化なし
```

**計算式:**
```
報酬 = (感情値 - 50) × 3 / 10
```

**例:**
- 感情値 90 → 報酬 +12
- 感情値 50 → 報酬 0
- 感情値 20 → 報酬 -9

### ApplyDecay

時間経過による自然減衰を適用します。

```go
bg.SetMotivation(80)
bg.ApplyDecay()
// 意欲が50（中立）に向かって減衰
```

**特徴:**
- 50より高い → 50に向かって減衰
- 50より低い → 50に向かって回復
- 減衰率: 95%

### ShouldTakeAction

意欲レベルに基づいて行動を起こすべきか判定します。

```go
bg.SetMotivation(70)

if bg.ShouldTakeAction(60) {
    fmt.Println("行動を起こす")
} else {
    fmt.Println("行動しない")
}
```

**パラメータ:**
- `threshold int`: 行動閾値

**戻り値:** bool (意欲 ≥ 閾値)

## 実践例

### 例1: 学習サイクル

```go
bg := basal.New()

// 学習開始
fmt.Printf("学習開始: 意欲 %d\n", bg.GetMotivation())

// 問題を解く
bg.UpdateMotivation(true)  // 正解した
fmt.Printf("正解後: 意欲 %d\n", bg.GetMotivation())

// 難しい問題に挑戦
bg.UpdateMotivation(false)  // 間違えた
fmt.Printf("不正解後: 意欲 %d\n", bg.GetMotivation())

// 再度挑戦して成功
bg.UpdateMotivation(true)
bg.UpdateMotivation(true)
fmt.Printf("連続正解後: 意欲 %d\n", bg.GetMotivation())
```

### 例2: 感情と意欲の連携

```go
bg := basal.New()

// ポジティブな体験
joyEmotion := 85
bg.RewardFromEmotion(joyEmotion)
fmt.Printf("喜びの体験後: 意欲 %d (%s)\n", 
    bg.GetMotivation(), bg.GetMotivationLevel())

// ネガティブな体験
fearEmotion := 25
bg.RewardFromEmotion(fearEmotion)
fmt.Printf("恐怖の体験後: 意欲 %d (%s)\n", 
    bg.GetMotivation(), bg.GetMotivationLevel())
```

### 例3: 時間経過シミュレーション

```go
bg := basal.New()
bg.SetMotivation(90)

fmt.Println("時間経過による減衰:")
for i := 0; i < 10; i++ {
    fmt.Printf("時刻 %d: 意欲 %d\n", i, bg.GetMotivation())
    bg.ApplyDecay()
}
```

### 例4: 行動閾値による制御

```go
bg := basal.New()

tasks := []struct {
    name      string
    threshold int
}{
    {"簡単なタスク", 30},
    {"普通のタスク", 50},
    {"難しいタスク", 70},
}

for _, task := range tasks {
    if bg.ShouldTakeAction(task.threshold) {
        fmt.Printf("✓ %s を実行\n", task.name)
    } else {
        fmt.Printf("✗ %s をスキップ（意欲不足）\n", task.name)
    }
}
```

## API経由での使用

### 意欲の取得

```bash
curl http://localhost:8081/api/motivation
```

レスポンス:
```json
{
  "motivation": 50,
  "level": "normal"
}
```

### ポジティブフィードバック

```bash
curl -X POST http://localhost:8081/api/motivation/feedback \
  -H "Content-Type: application/json" \
  -d '{"is_positive": true}'
```

### ネガティブフィードバック

```bash
curl -X POST http://localhost:8081/api/motivation/feedback \
  -H "Content-Type: application/json" \
  -d '{"is_positive": false}'
```

### 感情からの報酬

```bash
curl -X POST http://localhost:8081/api/motivation/emotion-reward \
  -H "Content-Type: application/json" \
  -d '{"emotion_value": 85}'
```

### 自然減衰の適用

```bash
curl -X POST http://localhost:8081/api/motivation/decay
```

### リセット

```bash
curl -X POST http://localhost:8081/api/motivation/reset
```

## パラメータのカスタマイズ

```go
bg := basal.New()

// 報酬パラメータを変更
bg.SetRewardParameters(15, 20)  // ポジティブ+15, ネガティブ-20

// 現在のパラメータを取得
positive, negative := bg.GetRewardParameters()
fmt.Printf("報酬: +%d, -%d\n", positive, negative)
```

## スレッドセーフ性

BasalGangliaは `sync.RWMutex` を使用してスレッドセーフに実装されています。

```go
bg := basal.New()

// 複数のゴルーチンから安全にアクセス可能
go bg.UpdateMotivation(true)
go bg.GetMotivation()
go bg.RewardFromEmotion(70)
```

## ドーパミンモデル

このモジュールは、脳内のドーパミンシステムをシミュレートしています:

1. **報酬予測誤差**: ポジティブフィードバック = ドーパミン放出
2. **報酬の欠如**: ネガティブフィードバック = ドーパミン減少
3. **自然減衰**: 時間経過で中立状態に戻る
4. **感情との連動**: 感情の強度が意欲に影響

## 設計上の特徴

### 1. 非対称な報酬
- ポジティブ: +10
- ネガティブ: -15

→ ネガティブフィードバックの影響が大きい（現実の心理に近い）

### 2. 自然減衰
- 極端な意欲は時間とともに中立(50)に戻る
- 減衰率: 95%

### 3. 範囲制限
- 最小値: 0
- 最大値: 100
- 自動的にクランプ

### 4. 感情連動
- 感情値から自動的に報酬を計算
- 高い感情 → 意欲上昇
- 低い感情 → 意欲低下

## 今後の拡張予定

- [ ] 学習率の動的調整
- [ ] 複数の報酬チャネル
- [ ] 長期的な意欲トレンド分析
- [ ] 目標設定と達成度追跡
- [ ] 習慣形成のシミュレーション
- [ ] 意欲履歴の記録と可視化
