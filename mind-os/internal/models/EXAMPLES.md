# データ構造使用例

## EmotionValue の作成

```go
import "github.com/your-username/mind-os/internal/models"

// 単一の感情値を作成
joy := models.EmotionValue{
    Code:  models.EmotionJoy,
    Value: 85,
}

// バリデーション
if !joy.Validate() {
    // エラー処理
}
```

## EmotionMap の使用

```go
// マップ形式で感情を管理
emotions := models.EmotionMap{
    models.EmotionJoy:  80,
    models.EmotionHope: 60,
    models.EmotionLove: 70,
}

// スライスに変換
values := emotions.ToEmotionValues()

// スライスからマップに変換
emotionMap := models.FromEmotionValues(values)
```

## RuneMemory の作成

```go
import (
    "time"
    "github.com/google/uuid"
    "github.com/your-username/mind-os/internal/models"
)

// 記憶ノードを作成
memory := models.RuneMemory{
    UUID: uuid.New().String(),
    Text: "ユーザーが「ありがとう」と言った",
    Emotions: []models.EmotionValue{
        {Code: models.EmotionJoy, Value: 70},
        {Code: models.EmotionLove, Value: 50},
    },
    Weight:     0.8,
    Type:       models.MemoryLTM,
    CreatedAt:  time.Now(),
    LastAccess: time.Now(),
    Tags:       []string{"感謝", "ポジティブ"},
}
```

## MindStateResponse の作成

```go
// クライアントへのレスポンスを構築
response := models.MindStateResponse{
    CurrentReaction: []models.EmotionValue{
        {Code: models.EmotionJoy, Value: 85},
        {Code: models.EmotionHope, Value: 60},
    },
    MoodStability: 0.75,
    PersonalityBias: []models.EmotionValue{
        {Code: models.EmotionHope, Value: 65},
        {Code: models.EmotionLove, Value: 55},
    },
    Motivation: 0.80,
    Sanity:     0.90,
}

// JSONに変換してレスポンス
c.JSON(http.StatusOK, response)
```

## 感情コードのバリデーション

```go
// 感情コードが有効かチェック
code := models.EmotionCode("J")
if models.IsValidEmotionCode(code) {
    // 有効な感情コード
}

// 無効な感情コード
invalidCode := models.EmotionCode("X")
if !models.IsValidEmotionCode(invalidCode) {
    // エラー処理
}
```

## 記憶タイプの判定

```go
memory := models.RuneMemory{
    Type: models.MemorySTM,
    // ...
}

switch memory.Type {
case models.MemorySTM:
    // 短期記憶の処理
    // 一定時間後に削除または長期記憶に昇格
case models.MemoryLTM:
    // 長期記憶の処理
    // 永続化または重み付けによる保持
}
```

## JSON シリアライゼーション例

### EmotionValue
```json
{
  "code": "J",
  "value": 85
}
```

### RuneMemory
```json
{
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "text": "ユーザーが「ありがとう」と言った",
  "emotions": [
    {"code": "J", "value": 70},
    {"code": "L", "value": 50}
  ],
  "weight": 0.8,
  "type": "LTM",
  "created_at": "2026-01-19T13:57:50+09:00",
  "last_access": "2026-01-19T13:57:50+09:00",
  "tags": ["感謝", "ポジティブ"]
}
```

### MindStateResponse
```json
{
  "current_reaction": [
    {"code": "J", "value": 85},
    {"code": "H", "value": 60}
  ],
  "mood_stability": 0.75,
  "personality_bias": [
    {"code": "H", "value": 65},
    {"code": "L", "value": 55}
  ],
  "motivation": 0.80,
  "sanity": 0.90
}
```
