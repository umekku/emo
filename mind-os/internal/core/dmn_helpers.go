package core

import (
	"github.com/umekku/mind-os/internal/models"
)

// getCurrentEmotions は現在の感情状態を取得
// 直近の記憶から感情を推定
func (b *Brain) getCurrentEmotions() []models.EmotionValue {
	// 直近の記憶から感情を取得
	recentMemories := b.Hippocampus.GetRecentContext()

	if len(recentMemories) == 0 {
		// 記憶がない場合はニュートラル
		return []models.EmotionValue{
			{Code: models.EmotionNeutral, Value: 50},
		}
	}

	// 最新の記憶の感情を返す
	return recentMemories[0].Emotions
}

// recallMemoriesWithReconsolidation は記憶を回想し、再固定化を適用
func (b *Brain) recallMemoriesWithReconsolidation(count int, moodTendency string, currentEmotions []models.EmotionValue) []models.RuneMemory {
	// Hippocampusから全記憶を取得（再固定化適用）
	allMemories := b.Hippocampus.GetRecentContextWithReconsolidation(currentEmotions)

	if len(allMemories) == 0 {
		return nil
	}

	// 気分に応じた重み付け
	weightedMemories := make([]weightedMemory, 0, len(allMemories))
	for _, mem := range allMemories {
		weight := calculateMemoryWeight(mem, moodTendency)
		weightedMemories = append(weightedMemories, weightedMemory{
			memory: mem,
			weight: weight,
		})
	}

	// 重み付き抽選で記憶を選抜
	selected := weightedRandomSelection(weightedMemories, count)

	return selected
}
