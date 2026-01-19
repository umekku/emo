package hippocampus

import (
	"log/slog"
	"time"

	"github.com/umekku/mind-os/internal/models"
)

// ReconsolidateMemory は記憶の再固定化を行う
// 【脳科学的意味】
// 記憶は想起されると不安定な状態に戻り、再固定化(Reconsolidation)される必要があります。
// この過程で、現在の感情や新しい情報が記憶に統合され、記憶の内容が書き換わります。
// これにより、トラウマの克服や記憶のアップデートが可能になります。
//
// 記憶を想起した際に、現在の感情状態で記憶を書き換える
func (h *Hippocampus) ReconsolidateMemory(memory *models.RuneMemory, currentEmotions []models.EmotionValue) {
	if memory == nil {
		return
	}

	// 1. 感情の上書き(Affective Coloring)
	// 現在の感情を0.1の係数でブレンド
	blendedEmotions := blendEmotions(memory.Emotions, currentEmotions, 0.1)
	memory.Emotions = blendedEmotions

	// 2. アクセスカウント更新
	memory.RecallCount++
	memory.LastAccess = time.Now()

	// 3. 重みの再計算（想起回数が多いほど重要）
	recallBonus := float64(memory.RecallCount) * 0.05
	memory.Weight += recallBonus
	if memory.Weight > 1.0 {
		memory.Weight = 1.0
	}

	// 4. DB更新（LTMの場合のみ）
	// TODO: Implement UpdateMemory in store package
	if memory.Type == models.MemoryLTM && h.store != nil {
		// if err := h.store.UpdateMemory(*memory); err != nil {
		// 	slog.Warn("Failed to update memory during reconsolidation", "error", err)
		// }
		slog.Info("Memory reconsolidated", "uuid", memory.UUID, "recallCount", memory.RecallCount)
	}
}

// blendEmotions は2つの感情リストをブレンド
func blendEmotions(original []models.EmotionValue, current []models.EmotionValue, blendRatio float64) []models.EmotionValue {
	// 現在の感情をマップに変換
	currentMap := make(map[models.EmotionCode]int)
	for _, e := range current {
		currentMap[e.Code] = e.Value
	}

	// 元の感情に現在の感情をブレンド
	blended := make([]models.EmotionValue, len(original))
	for i, origEmotion := range original {
		blended[i] = origEmotion

		// 同じ感情コードが現在の感情にあればブレンド
		if currentValue, exists := currentMap[origEmotion.Code]; exists {
			// 元の感情値 × (1 - ratio) + 現在の感情値 × ratio
			newValue := float64(origEmotion.Value)*(1.0-blendRatio) + float64(currentValue)*blendRatio
			blended[i].Value = int(newValue)
			if blended[i].Value > 100 {
				blended[i].Value = 100
			}
		}
	}

	return blended
}

// GetRecentContextWithReconsolidation は記憶を取得し、再固定化を適用
func (h *Hippocampus) GetRecentContextWithReconsolidation(currentEmotions []models.EmotionValue) []models.RuneMemory {
	// 元のGetRecentContextで記憶を取得
	memories := h.GetRecentContext()

	// 各記憶に再固定化を適用
	for i := range memories {
		h.ReconsolidateMemory(&memories[i], currentEmotions)
	}

	return memories
}
