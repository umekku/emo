# Prefrontal Cortex モジュール使用例

## 概要

Prefrontal Cortex（前頭前皮質）モジュールは、**理性(Sanity)による感情抑制** を管理します。
高い理性値により、ネガティブな感情を減衰させ、ポジティブな感情を増幅することで、感情的な反応を制御します。

## 基本的な使用方法

```go
package main

import (
    "fmt"
    "github.com/your-username/mind-os/internal/pfc"
    "github.com/your-username/mind-os/internal/models"
)

func main() {
    // PrefrontalCortexインスタンスを作成
    pfc := pfc.New()
    
    fmt.Printf("初期理性値: %d\n", pfc.GetSanity())
    
    // 生の感情
    raw := []models.EmotionValue{
        {Code: models.EmotionAnger, Value: 80},
        {Code: models.EmotionFear, Value: 70},
        {Code: models.EmotionJoy, Value: 60},
    }
    
    // 理性による調停
    arbitrated := pfc.Arbitrate(raw)
    
    fmt.Println("調停後の感情:")
    for _, emotion := range arbitrated {
        fmt.Printf("  %s: %d\n", emotion.Code, emotion.Value)
    }
}
```

## 理性レベル

理性値は **0-100** の範囲で管理され、以下のレベルに分類されます:

| 理性値 | レベル | 説明 |
|--------|--------|------|
| 80-100 | `very_high` | 非常に高い理性 |
| 60-79  | `high` | 高い理性 |
| 40-59  | `normal` | 通常の理性 |
| 20-39  | `low` | 低い理性 |
| 0-19   | `very_low` | 非常に低い理性 |

**初期値**: 80 (very_high)

## 感情調停のメカニズム

### 高理性時 (Sanity ≥ 30)

#### ネガティブ感情の抑制
- **対象**: Anger (怒り), Fear (恐れ), Disgust (嫌悪)
- **抑制率**: `Sanity × 50%`
- **計算式**: `新しい値 = 元の値 - (元の値 × Sanity/100 × 0.5)`

**例:**
```
Sanity = 100, Anger = 80
抑制 = 80 × 1.0 × 0.5 = 40
結果 = 80 - 40 = 40
```

#### ポジティブ感情の増幅
- **対象**: Joy (喜び), Love (愛), Hope (希望)
- **増幅率**: `Sanity × 10%`
- **計算式**: `新しい値 = 元の値 + (元の値 × Sanity/100 × 0.1)`

**例:**
```
Sanity = 100, Joy = 60
増幅 = 60 × 1.0 × 0.1 = 6
結果 = 60 + 6 = 66
```

#### 中立感情
- **対象**: Surprise (驚き), Neutral (中立)
- **処理**: 変更なし

### 低理性時 (Sanity < 30)

- すべての感情がそのまま通過
- 感情制御不能状態

## 主要メソッド

### Arbitrate

生の感情を理性により調停します。

```go
pfc := pfc.New()

raw := []models.EmotionValue{
    {Code: models.EmotionAnger, Value: 80},
    {Code: models.EmotionJoy, Value: 60},
}

arbitrated := pfc.Arbitrate(raw)
// Anger: 80 → 40 (抑制)
// Joy: 60 → 66 (増幅)
```

### GetSanity

現在の理性値を取得します。

```go
sanity := pfc.GetSanity()
fmt.Printf("現在の理性: %d\n", sanity)
```

### SetSanity

理性値を設定します。

```go
pfc.SetSanity(60)
```

### UpdateSanity

理性値を増減します。

```go
pfc.UpdateSanity(10)   // +10
pfc.UpdateSanity(-20)  // -20
```

### ApplyStress

ストレスにより理性値を減少させます。

```go
// ストレスレベル 0-100
pfc.ApplyStress(80)  // 最大16減少
```

**計算式**: `減少量 = ストレスレベル / 5`

### Rest

休息により理性値を回復させます。

```go
// 休息の質 0-100
pfc.Rest(100)  // 最大25回復
```

**計算式**: `回復量 = 休息の質 / 4`

### CanControlEmotions

感情を制御できるか判定します。

```go
if pfc.CanControlEmotions() {
    fmt.Println("感情制御可能")
} else {
    fmt.Println("感情制御不能")
}
```

**条件**: Sanity ≥ 30

### GetSuppressionRate

抑制率を取得します。

```go
rate := pfc.GetSuppressionRate()
fmt.Printf("抑制率: %.2f\n", rate)
```

**戻り値**: 0.0-1.0 (Sanity < 30 の場合は 0.0)

### CalculateEmotionalImpact

感情の影響度を計算します。

```go
emotions := []models.EmotionValue{
    {Code: models.EmotionAnger, Value: 80},
    {Code: models.EmotionFear, Value: 60},
}

impact := pfc.CalculateEmotionalImpact(emotions)
fmt.Printf("感情影響度: %d\n", impact)
```

## 実践例

### 例1: ストレスと休息のサイクル

```go
pfc := pfc.New()

fmt.Printf("初期理性: %d\n", pfc.GetSanity())

// 仕事でストレス
pfc.ApplyStress(80)
fmt.Printf("ストレス後: %d (%s)\n", 
    pfc.GetSanity(), pfc.GetSanityLevel())

// 休息
pfc.Rest(100)
fmt.Printf("休息後: %d (%s)\n", 
    pfc.GetSanity(), pfc.GetSanityLevel())
```

### 例2: 理性レベルによる感情制御の違い

```go
raw := []models.EmotionValue{
    {Code: models.EmotionAnger, Value: 90},
}

// 高理性
pfc1 := pfc.New()
pfc1.SetSanity(100)
result1 := pfc1.Arbitrate(raw)
fmt.Printf("高理性(100): Anger %d → %d\n", 
    raw[0].Value, result1[0].Value)

// 中理性
pfc2 := pfc.New()
pfc2.SetSanity(50)
result2 := pfc2.Arbitrate(raw)
fmt.Printf("中理性(50): Anger %d → %d\n", 
    raw[0].Value, result2[0].Value)

// 低理性
pfc3 := pfc.New()
pfc3.SetSanity(20)
result3 := pfc3.Arbitrate(raw)
fmt.Printf("低理性(20): Anger %d → %d\n", 
    raw[0].Value, result3[0].Value)
```

### 例3: 複合感情の調停

```go
pfc := pfc.New()
pfc.SetSanity(80)

raw := []models.EmotionValue{
    {Code: models.EmotionAnger, Value: 80},
    {Code: models.EmotionFear, Value: 70},
    {Code: models.EmotionJoy, Value: 60},
    {Code: models.EmotionLove, Value: 50},
    {Code: models.EmotionNeutral, Value: 50},
}

arbitrated := pfc.Arbitrate(raw)

fmt.Println("感情調停結果:")
for i, emotion := range arbitrated {
    change := emotion.Value - raw[i].Value
    symbol := ""
    if change > 0 {
        symbol = "↑"
    } else if change < 0 {
        symbol = "↓"
    } else {
        symbol = "→"
    }
    fmt.Printf("  %s: %d %s %d (%+d)\n", 
        emotion.Code, raw[i].Value, symbol, emotion.Value, change)
}
```

### 例4: 感情影響度の計算

```go
pfc := pfc.New()

emotions := []models.EmotionValue{
    {Code: models.EmotionAnger, Value: 80},
    {Code: models.EmotionFear, Value: 60},
}

// 高理性時
pfc.SetSanity(100)
highImpact := pfc.CalculateEmotionalImpact(emotions)
fmt.Printf("高理性時の影響度: %d\n", highImpact)

// 低理性時
pfc.SetSanity(20)
lowImpact := pfc.CalculateEmotionalImpact(emotions)
fmt.Printf("低理性時の影響度: %d\n", lowImpact)

fmt.Printf("影響度の差: %d\n", lowImpact - highImpact)
```

## 統合例: 扁桃体 + 前頭前皮質

```go
import (
    "github.com/your-username/mind-os/internal/amygdala"
    "github.com/your-username/mind-os/internal/pfc"
)

func processEmotion(text string) {
    // 扁桃体で反射的な感情を生成
    amy := amygdala.New()
    rawEmotions := amy.Assess(text)
    
    fmt.Println("生の感情:")
    for _, e := range rawEmotions {
        fmt.Printf("  %s: %d\n", e.Code, e.Value)
    }
    
    // 前頭前皮質で理性的に調停
    pfc := pfc.New()
    controlledEmotions := pfc.Arbitrate(rawEmotions)
    
    fmt.Println("調停後の感情:")
    for _, e := range controlledEmotions {
        fmt.Printf("  %s: %d\n", e.Code, e.Value)
    }
}

// 使用例
processEmotion("バグが見つかった")
// 生の感情: Fear:80, Anger:40
// 調停後: Fear:48, Anger:24 (理性により抑制)
```

## スレッドセーフ性

PrefrontalCortexは `sync.RWMutex` を使用してスレッドセーフに実装されています。

```go
pfc := pfc.New()

// 複数のゴルーチンから安全にアクセス可能
go pfc.UpdateSanity(10)
go pfc.GetSanity()
go pfc.Arbitrate(emotions)
```

## 設計上の特徴

### 1. 非対称な処理
- ネガティブ感情: 50%抑制
- ポジティブ感情: 10%増幅

→ ネガティブ感情の抑制効果が大きい

### 2. 閾値ベースの制御
- Sanity ≥ 30: 感情制御可能
- Sanity < 30: 感情制御不能

### 3. 段階的な効果
- 理性値が高いほど抑制効果が大きい
- 線形的な関係

### 4. ストレスと休息
- ストレス: 理性を減少
- 休息: 理性を回復

## 理性値の管理戦略

### 理性を維持する方法
1. 定期的な休息 (`Rest()`)
2. ストレスの軽減
3. ポジティブな体験

### 理性が低下する原因
1. 高ストレス (`ApplyStress()`)
2. 長時間の活動
3. ネガティブな体験

## 今後の拡張予定

- [ ] 感情ごとの個別抑制率設定
- [ ] 学習による抑制パターンの最適化
- [ ] 長期的な理性トレンド分析
- [ ] 疲労度との連携
- [ ] 瞑想・マインドフルネスシミュレーション
- [ ] 理性履歴の記録と可視化
