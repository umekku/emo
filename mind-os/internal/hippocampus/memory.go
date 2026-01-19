// memory.go: エピソード記憶の短期保持(STM)と長期固定化(LTM)、および忘却プロセスを管理する海馬モジュール
package hippocampus

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/umekku/mind-os/internal/models"
	"github.com/umekku/mind-os/internal/store"
)

// Hippocampus は海馬モジュール - 短期記憶と長期記憶を管理
type Hippocampus struct {
	STM []models.RuneMemory // 短期記憶 (Short-Term Memory) - メモリのみ
	// LTM []models.RuneMemory // 長期記憶 (LTM) - DBに移行するため削除

	store *store.DB // DB接続

	// 記憶の閾値設定
	consolidationThreshold float64 // LTMへの移行閾値
	maxSTMSize             int     // STMの最大サイズ
	maxLTMSize             int     // LTMの最大サイズ
}

// New は新しい Hippocampus インスタンスを作成
func New(db *store.DB) *Hippocampus {
	return &Hippocampus{
		STM:                    make([]models.RuneMemory, 0),
		store:                  db,
		consolidationThreshold: 0.6,  // 重み0.6以上でLTMへ移行
		maxSTMSize:             100,  // STM最大100件
		maxLTMSize:             1000, // LTM最大1000件
	}
}

// AddEpisode は新しいエピソード記憶をSTMに追加
func (h *Hippocampus) AddEpisode(text string, emotions []models.EmotionValue) {
	now := time.Now()

	// 感情の強度から重みを計算 (0.0-1.0)
	weight := h.calculateWeight(emotions)

	// 新しい記憶を作成
	memory := models.RuneMemory{
		UUID:       uuid.New().String(),
		Text:       text,
		Emotions:   emotions,
		Weight:     weight,
		Type:       models.MemorySTM,
		CreatedAt:  now,
		LastAccess: now,
		Tags:       h.extractTags(text, emotions),
	}

	// STMに追加
	h.STM = append(h.STM, memory)

	// STMサイズ制限チェック
	if len(h.STM) > h.maxSTMSize {
		// 古い記憶から削除（FIFO）
		h.STM = h.STM[1:]
	}
}

// AddMemory は外部ヘルパー用 (AddEpisodeのラッパー)
// MemoryHandlerとの互換性のため、RuneMemoryを返す
func (h *Hippocampus) AddMemory(text string, emotions []models.EmotionValue) models.RuneMemory {
	now := time.Now()
	weight := h.calculateWeight(emotions)

	memory := models.RuneMemory{
		UUID:       uuid.New().String(),
		Text:       text,
		Emotions:   emotions,
		Weight:     weight,
		Type:       models.MemorySTM,
		CreatedAt:  now,
		LastAccess: now,
		Tags:       h.extractTags(text, emotions),
	}

	h.STM = append(h.STM, memory)
	if len(h.STM) > h.maxSTMSize {
		h.STM = h.STM[1:]
	}

	return memory
}

// GetRecentContext は直近の記憶を返す
func (h *Hippocampus) GetRecentContext() []models.RuneMemory {
	// LTMから直近の記憶を取得
	var ltmMemories []models.RuneMemory
	if h.store != nil {
		var err error
		ltmMemories, err = h.store.GetRecentMemories(10)
		if err != nil {
			slog.Error("Failed to fetch LTM", "error", err)
		}
	}

	// STMとLTMを結合して最新順にソート
	allMemories := make([]models.RuneMemory, 0, len(h.STM)+len(ltmMemories))
	allMemories = append(allMemories, h.STM...)
	allMemories = append(allMemories, ltmMemories...)

	// 最終アクセス時間でソート（新しい順）
	for i := 0; i < len(allMemories); i++ {
		for j := i + 1; j < len(allMemories); j++ {
			if allMemories[i].LastAccess.Before(allMemories[j].LastAccess) {
				allMemories[i], allMemories[j] = allMemories[j], allMemories[i]
			}
		}
	}

	// 直近10件を返す
	limit := 10
	if len(allMemories) < limit {
		limit = len(allMemories)
	}

	return allMemories[:limit]
}

// SleepAndConsolidate は睡眠処理 - STMからLTMへの記憶の固定化
func (h *Hippocampus) SleepAndConsolidate() {
	if h.store == nil {
		return
	}

	// STMの記憶を重みでフィルタリング
	for _, memory := range h.STM {
		if memory.Weight >= h.consolidationThreshold {
			// 重みが閾値以上の記憶はLTMへ移行
			ltmMemory := memory
			ltmMemory.Type = models.MemoryLTM
			ltmMemory.LastAccess = time.Now()

			// DBに保存
			if err := h.store.SaveMemory(ltmMemory); err != nil {
				slog.Error("Failed to save LTM", "error", err)
			}
		}
		// 閾値未満の記憶は忘却（何もしない）
	}

	// STMをクリア
	h.STM = make([]models.RuneMemory, 0)

	// 古い記憶の風化処理
	h.FadeOldMemories()

	// LTMサイズ制限チェック（古い記憶の削除）
	if err := h.store.DeleteOldMemories(h.maxLTMSize); err != nil {
		slog.Warn("Failed to prune old LTM", "error", err)
	}
}

// FadeOldMemories は古い記憶の重みを減衰させる
func (h *Hippocampus) FadeOldMemories() {
	// DB側の処理が必要だが、ここでは簡易的実装
	// 実際にはDB内の全レコードの重みを下げるクエリなどを発行
}

// calculateWeight は感情の強度から記憶の重みを計算
func (h *Hippocampus) calculateWeight(emotions []models.EmotionValue) float64 {
	if len(emotions) == 0 {
		return 0.5 // デフォルト値
	}

	// 感情値の平均を計算
	total := 0
	for _, emotion := range emotions {
		total += emotion.Value
	}
	average := float64(total) / float64(len(emotions))

	// 0-100を0.0-1.0に正規化
	weight := average / 100.0

	// 感情の種類数も考慮（多様な感情ほど重要）
	diversityBonus := float64(len(emotions)) * 0.05
	weight += diversityBonus

	// 上限を1.0に制限
	if weight > 1.0 {
		weight = 1.0
	}

	return weight
}

// extractTags はテキストと感情からタグを抽出
func (h *Hippocampus) extractTags(text string, emotions []models.EmotionValue) []string {
	tags := make([]string, 0)

	// 感情コードをタグとして追加
	emotionMap := make(map[models.EmotionCode]bool)
	for _, emotion := range emotions {
		if !emotionMap[emotion.Code] {
			tags = append(tags, string(emotion.Code))
			emotionMap[emotion.Code] = true
		}
	}

	// テキストの長さに基づくタグ
	if len(text) > 100 {
		tags = append(tags, "long")
	} else if len(text) < 20 {
		tags = append(tags, "short")
	}

	return tags
}

// GetSTMCount はSTMの記憶数を返す
func (h *Hippocampus) GetSTMCount() int {
	return len(h.STM)
}

// GetLTMCount はLTMの記憶数を返す
func (h *Hippocampus) GetLTMCount() int {
	if h.store == nil {
		return 0
	}
	count, err := h.store.GetLTMCount()
	if err != nil {
		slog.Error("Failed to get LTM count", "error", err)
		return 0
	}
	return count
}

// GetMemoryByUUID はUUIDで記憶を検索
func (h *Hippocampus) GetMemoryByUUID(uuid string) *models.RuneMemory {
	// STMを検索
	for _, memory := range h.STM {
		if memory.UUID == uuid {
			return &memory
		}
	}

	// LTMを検索 (DB)
	if h.store != nil {
		memory, err := h.store.GetMemoryByUUID(uuid)
		if err != nil {
			slog.Warn("Failed to get memory by UUID", "error", err)
			return nil
		}
		return memory
	}

	return nil
}
