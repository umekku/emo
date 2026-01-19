package core

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/umekku/mind-os/internal/models"
)

// WanderMind はデフォルトモードネットワーク(DMN)による自発的思考
// 非アクティブ時間中に記憶を回想し、感情・気分に影響を与える
func (b *Brain) WanderMind(duration time.Duration) string {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 経過時間に基づいて回想する記憶の数を決定
	// 1時間ごとに1件、最大5件
	hours := int(duration.Hours())
	recallCount := hours
	if recallCount > 5 {
		recallCount = 5
	}
	if recallCount < 1 {
		recallCount = 1
	}

	// 記憶が存在しない場合は何もしない
	if b.Hippocampus == nil {
		return "（記憶がないため、何も思い出せなかった）"
	}

	// 現在の気分状態を取得
	currentMood := b.getCurrentMoodTendency()

	// 現在の感情状態を取得（再固定化用）
	// PFCから最近の感情を推定
	currentEmotions := b.getCurrentEmotions()

	// 回想ログを構築
	var thoughtLog strings.Builder
	thoughtLog.WriteString("【マインドワンダリング】\n")

	// 記憶を回想（再固定化を適用）
	memories := b.recallMemoriesWithReconsolidation(recallCount, currentMood, currentEmotions)

	if len(memories) == 0 {
		thoughtLog.WriteString("（まだ記憶が形成されていない...）")
		return thoughtLog.String()
	}

	// 各記憶について反芻
	for i, memory := range memories {
		// 記憶の内容を要約
		summary := summarizeMemory(memory)
		thoughtLog.WriteString(fmt.Sprintf("%d. %s\n", i+1, summary))

		// 記憶の感情価を現在の状態に微量加算（反芻効果）
		b.ruminateOnMemory(memory)
	}

	return thoughtLog.String()
}

// getCurrentMoodTendency は現在の気分傾向を取得
// ネガティブ・ポジティブ・ニュートラルを判定
func (b *Brain) getCurrentMoodTendency() string {
	// PFCの理性値とBasalの意欲値から判定
	sanity := b.PFC.GetSanity()
	motivation := b.BasalGanglia.GetMotivation()

	if sanity < 30 || motivation < 30 {
		return "negative" // 抑うつ的
	} else if sanity > 70 && motivation > 70 {
		return "positive" // ポジティブ
	}
	return "neutral" // 中立
}

// weightedMemory は重み付き記憶
type weightedMemory struct {
	memory models.RuneMemory
	weight float64
}

// calculateMemoryWeight は記憶の重みを計算
// 現在の気分に近い記憶ほど思い出しやすい（気分一致効果）
func calculateMemoryWeight(memory models.RuneMemory, moodTendency string) float64 {
	baseWeight := memory.Weight // 記憶の重要度

	// 感情価による重み付け
	emotionBonus := 0.0
	for _, emotion := range memory.Emotions {
		switch moodTendency {
		case "negative":
			// ネガティブな気分の時はネガティブな記憶を思い出しやすい
			if emotion.Code == models.EmotionAnger ||
				emotion.Code == models.EmotionFear ||
				emotion.Code == models.EmotionDisgust ||
				emotion.Code == models.EmotionGrief {
				emotionBonus += float64(emotion.Value) * 0.02
			}
		case "positive":
			// ポジティブな気分の時はポジティブな記憶を思い出しやすい
			if emotion.Code == models.EmotionJoy ||
				emotion.Code == models.EmotionLove ||
				emotion.Code == models.EmotionHope {
				emotionBonus += float64(emotion.Value) * 0.02
			}
		}
	}

	return baseWeight + emotionBonus
}

// weightedRandomSelection は重み付きランダム選択
func weightedRandomSelection(weighted []weightedMemory, count int) []models.RuneMemory {
	if len(weighted) == 0 {
		return nil
	}

	// 選択可能な数を制限
	if count > len(weighted) {
		count = len(weighted)
	}

	// 総重みを計算
	totalWeight := 0.0
	for _, w := range weighted {
		totalWeight += w.weight
	}

	if totalWeight <= 0 {
		// 重みがない場合はランダムに選択
		selected := make([]models.RuneMemory, 0, count)
		indices := rand.Perm(len(weighted))
		for i := 0; i < count; i++ {
			selected = append(selected, weighted[indices[i]].memory)
		}
		return selected
	}

	// 重み付き抽選
	selected := make([]models.RuneMemory, 0, count)
	remaining := make([]weightedMemory, len(weighted))
	copy(remaining, weighted)

	for i := 0; i < count && len(remaining) > 0; i++ {
		// ルーレット選択
		r := rand.Float64() * totalWeight
		cumulative := 0.0

		for j, w := range remaining {
			cumulative += w.weight
			if r <= cumulative {
				selected = append(selected, w.memory)
				// 選択した記憶を除外
				remaining = append(remaining[:j], remaining[j+1:]...)
				totalWeight -= w.weight
				break
			}
		}
	}

	return selected
}

// ruminateOnMemory は記憶を反芻し、現在の感情状態に影響を与える
func (b *Brain) ruminateOnMemory(memory models.RuneMemory) {
	// 記憶の感情を微量（10%）だけ現在の状態に加算
	for _, emotion := range memory.Emotions {
		ruminationValue := int(float64(emotion.Value) * 0.1)

		// 感情の種類に応じて脳の状態を更新
		switch emotion.Code {
		case models.EmotionAnger, models.EmotionFear, models.EmotionDisgust, models.EmotionGrief:
			// ネガティブ感情 → ストレス増加
			b.Hypothalamus.Update(float64(ruminationValue), 0)
		case models.EmotionJoy, models.EmotionLove:
			// ポジティブ感情 → 愛着増加
			b.Hypothalamus.Update(0, float64(ruminationValue))
		}
	}
}

// summarizeMemory は記憶を要約
func summarizeMemory(memory models.RuneMemory) string {
	// テキストの最初30文字を取得
	text := memory.Text
	if len(text) > 30 {
		text = text[:30] + "..."
	}

	// 主要な感情を取得
	dominantEmotion := "中立"
	maxValue := 0
	for _, emotion := range memory.Emotions {
		if emotion.Value > maxValue {
			maxValue = emotion.Value
			dominantEmotion = emotionCodeToJapanese(emotion.Code)
		}
	}

	return fmt.Sprintf("「%s」を思い出した（%s）", text, dominantEmotion)
}

// emotionCodeToJapanese は感情コードを日本語に変換
func emotionCodeToJapanese(code models.EmotionCode) string {
	switch code {
	case models.EmotionJoy:
		return "喜び"
	case models.EmotionAnger:
		return "怒り"
	case models.EmotionFear:
		return "恐れ"
	case models.EmotionLove:
		return "愛"
	case models.EmotionDisgust:
		return "嫌悪"
	case models.EmotionSurprise:
		return "驚き"
	case models.EmotionHope:
		return "希望"
	case models.EmotionGrief:
		return "悲嘆"
	default:
		return "中立"
	}
}
