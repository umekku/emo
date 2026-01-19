# Amygdala モジュール使用例

## 概要

Amygdala（扁桃体）モジュールは、入力テキストから**反射的な感情**を即座に生成します。
キーワードマッチングベースのシンプルな実装で、高速な感情評価を実現します。

## 基本的な使用方法

```go
package main

import (
    "fmt"
    "github.com/your-username/mind-os/internal/amygdala"
)

func main() {
    // Amygdalaインスタンスを作成
    amy := amygdala.New()
    
    // テキストから感情を評価
    emotions := amy.Assess("できた！")
    
    // 結果を表示
    for _, emotion := range emotions {
        fmt.Printf("%s: %d\n", emotion.Code, emotion.Value)
    }
    // 出力例:
    // J: 90  (Joy - 喜び)
    // H: 60  (Hope - 希望)
}
```

## キーワードと感情のマッピング

### ポジティブ感情

| キーワード | 感情コード | 強度 |
|-----------|-----------|------|
| できた | J (Joy) | 90 |
| | H (Hope) | 60 |
| 成功 | J (Joy) | 85 |
| | H (Hope) | 70 |
| 嬉しい | J (Joy) | 95 |
| ありがとう | J (Joy) | 70 |
| | L (Love) | 60 |
| 好き | L (Love) | 80 |
| | J (Joy) | 60 |

### ネガティブ感情

| キーワード | 感情コード | 強度 |
|-----------|-----------|------|
| バグ | F (Fear) | 80 |
| | A (Anger) | 40 |
| エラー | F (Fear) | 70 |
| | A (Anger) | 35 |
| むかつく | A (Anger) | 85 |
| | D (Disgust) | 60 |
| 嫌い | D (Disgust) | 80 |

### 挨拶

| キーワード | 感情コード | 強度 |
|-----------|-----------|------|
| おはよう | N (Neutral) | 50 |
| | J (Joy) | 30 |
| こんにちは | N (Neutral) | 50 |
| | J (Joy) | 30 |

## 実践例

### 例1: ポジティブなメッセージ

```go
amy := amygdala.New()
emotions := amy.Assess("テストが成功した！嬉しい！")

// 結果:
// J (Joy): 100  (95 + 85 = 180 → 上限100に制限)
// H (Hope): 70
```

### 例2: ネガティブなメッセージ

```go
amy := amygdala.New()
emotions := amy.Assess("バグが見つかった。最悪だ。")

// 結果:
// F (Fear): 80
// A (Anger): 100  (40 + 70 = 110 → 上限100に制限)
// D (Disgust): 85
```

### 例3: 混合感情

```go
amy := amygdala.New()
emotions := amy.Assess("バグが見つかったけど、修正できた！")

// 結果:
// F (Fear): 80   (バグから)
// A (Anger): 40  (バグから)
// J (Joy): 90    (できたから)
// H (Hope): 60   (できたから)
```

### 例4: 未知の入力

```go
amy := amygdala.New()
emotions := amy.Assess("xyz123")

// 結果:
// N (Neutral): 50  (マッチするキーワードがない場合)
```

## API経由での使用

### リクエスト例

```bash
curl -X POST http://localhost:8081/api/emotion/assess \
  -H "Content-Type: application/json" \
  -d '{"text":"ありがとう！助かった！"}'
```

### レスポンス例

```json
{
  "text": "ありがとう！助かった！",
  "emotions": [
    {"code": "J", "value": 70},
    {"code": "L", "value": 60}
  ]
}
```

## 設計上の特徴

### 1. 即座の反応
- キーワードマッチングによる高速処理
- 複雑な自然言語処理を行わず、反射的な感情を生成

### 2. 感情の重ね合わせ
- 複数のキーワードがマッチした場合、感情値を加算
- 上限は100に制限

### 3. デフォルト動作
- マッチするキーワードがない場合は Neutral:50 を返す

### 4. 大文字小文字の区別なし
- 入力テキストは自動的に小文字に変換されて評価

## カスタマイズ

現在のキーワードマッピングは `amygdala.go` の `initializeEmotionKeywords()` 関数で定義されています。
新しいキーワードを追加する場合は、この関数を編集してください。

```go
func initializeEmotionKeywords() map[string][]models.EmotionValue {
    return map[string][]models.EmotionValue{
        // 新しいキーワードを追加
        "最高": {
            {Code: models.EmotionJoy, Value: 95},
        },
        // ...
    }
}
```

## 今後の拡張予定

- [ ] 感情の減衰機能（時間経過で感情値が減少）
- [ ] 文脈を考慮した感情評価
- [ ] 機械学習モデルとの統合
- [ ] 多言語対応
- [ ] カスタム辞書のロード機能
