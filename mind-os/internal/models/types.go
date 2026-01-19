package models

import (
	"log/slog"
	"time"
)

// EmotionCode は感情コードを表す文字列型
type EmotionCode string

// 感情コード定数 (EMO v1.1仕様)
// 『プルチックの感情の輪』および一般的な情動分類に基づく
const (
	EmotionJoy      EmotionCode = "J" // 喜び (Joy) - ドーパミン系、報酬予測
	EmotionSurprise EmotionCode = "S" // 驚き (Surprise) - 注意惹起、学習トリガー
	EmotionAnger    EmotionCode = "A" // 怒り (Anger) - ノルアドレナリン系、闘争反応
	EmotionFear     EmotionCode = "F" // 恐れ (Fear) - 扁桃体の活性化、逃走反応
	EmotionLove     EmotionCode = "L" // 愛 (Love) - オキシトシン系、社会的絆
	EmotionDisgust  EmotionCode = "D" // 嫌悪 (Disgust) - 島皮質の活性化、拒絶行動
	EmotionHope     EmotionCode = "H" // 希望 (Hope) - 期待値上昇、セロトニン系安定
	EmotionGrief    EmotionCode = "G" // 悲嘆 (Grief) - 喪失反応、前帯状皮質
	EmotionSadness  EmotionCode = "G" // 悲しみ (Sadness) - Griefのエイリアス
	EmotionNeutral  EmotionCode = "N" // 中立 (Neutral) - ベースライン状態
)

// EmotionValue は感情コードと強度値を持つ構造体
type EmotionValue struct {
	Code  EmotionCode `json:"code"`  // 感情コード
	Value int         `json:"value"` // 強度値 (0-100)
}

// MemoryType は記憶の種類を表す文字列型
type MemoryType string

// 記憶タイプ定数
const (
	MemorySTM MemoryType = "STM" // 短期記憶 (Short-Term Memory)
	MemoryLTM MemoryType = "LTM" // 長期記憶 (Long-Term Memory)
)

// RuneMemory は記憶ノードを表す構造体
type RuneMemory struct {
	UUID        string         `json:"uuid"`        // 一意識別子
	Text        string         `json:"text"`        // 記憶内容
	Emotions    []EmotionValue `json:"emotions"`    // 関連感情
	Weight      float64        `json:"weight"`      // 重み (重要度)
	Type        MemoryType     `json:"type"`        // 記憶タイプ
	CreatedAt   time.Time      `json:"createdAt"`   // 作成日時
	LastAccess  time.Time      `json:"lastAccess"`  // 最終アクセス日時
	RecallCount int            `json:"recallCount"` // 想起回数
	Tags        []string       `json:"tags"`        // タグ
}

// SignalType は刺激の種類を表す文字列型
type SignalType string

const (
	SignalChat     SignalType = "chat"     // 会話テキスト
	SignalPhysical SignalType = "physical" // 物理的刺激（食事、接触、痛み）
)

// SensoryInput は感覚入力を表す構造体
type SensoryInput struct {
	Type        SignalType `json:"type" validate:"required,oneof=chat physical"` // 刺激の種類
	InputText   string     `json:"text" validate:"required,max=500"`             // 記憶用のテキスト記述
	SignalValue int        `json:"signalValue" validate:"min=-100,max=100"`      // -100(不快/痛み) 〜 +100(快感/報酬)
}

// LogValue はslog.Valuerインターフェースの実装
// ログ出力時に InputText をマスクして機密情報を保護する
func (s SensoryInput) LogValue() slog.Value {
	// 部分的なマスキング (先頭3文字のみ表示)
	maskedText := "***"
	if len(s.InputText) > 3 {
		maskedText = s.InputText[:3] + "***"
	}

	return slog.GroupValue(
		slog.String("type", string(s.Type)),
		slog.String("text", maskedText), // マスク済みテキスト
		slog.Int("signalValue", s.SignalValue),
	)
}

// MindStateResponse はクライアントへのレスポンス用構造体
type MindStateResponse struct {
	CurrentReaction []EmotionValue `json:"currentReaction"` // 反応感情
	MoodStability   float64        `json:"moodStability"`   // 安定度 (0.0-1.0)
	PersonalityBias []EmotionValue `json:"personalityBias"` // 性格傾向
	Motivation      float64        `json:"motivation"`      // 意欲 (0.0-1.0)
	Sanity          float64        `json:"sanity"`          // 理性値 (0.0-1.0)
	// デバッグ用フィールド
	Cortisol        float64 `json:"cortisol"`
	Oxytocin        float64 `json:"oxytocin"`
	PredictedReward float64 `json:"predictedReward"`
	DaydreamLog     string  `json:"daydreamLog,omitempty"` // マインドワンダリングログ
	ReplyText       string  `json:"replyText,omitempty"`   // 生成された応答テキスト
}

// EmotionMap は感情コードから強度値へのマッピング
type EmotionMap map[EmotionCode]int

// ToEmotionValues は EmotionMap を EmotionValue のスライスに変換
func (em EmotionMap) ToEmotionValues() []EmotionValue {
	values := make([]EmotionValue, 0, len(em))
	for code, value := range em {
		values = append(values, EmotionValue{
			Code:  code,
			Value: value,
		})
	}
	return values
}

// FromEmotionValues は EmotionValue のスライスを EmotionMap に変換
func FromEmotionValues(values []EmotionValue) EmotionMap {
	em := make(EmotionMap)
	for _, v := range values {
		em[v.Code] = v.Value
	}
	return em
}

// Validate は EmotionValue の値が有効範囲内かチェック
func (ev *EmotionValue) Validate() bool {
	return ev.Value >= 0 && ev.Value <= 100
}

// IsValidEmotionCode は感情コードが有効かチェック
func IsValidEmotionCode(code EmotionCode) bool {
	switch code {
	case EmotionJoy, EmotionSurprise, EmotionAnger, EmotionFear,
		EmotionLove, EmotionDisgust, EmotionHope, EmotionGrief, EmotionNeutral:
		return true
	default:
		return false
	}
}
