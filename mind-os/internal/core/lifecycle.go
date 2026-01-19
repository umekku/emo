package core

import (
	"github.com/umekku/mind-os/internal/models"
)

// Sleep は睡眠処理を実行（記憶の整理）
// 【神経科学的意味】睡眠中に短期記憶を長期記憶に固定化し、不要な記憶を忘却
// 【処理内容】Hippocampusの記憶固定化プロセスを実行
func (b *Brain) Sleep() SleepResult {
	b.mu.Lock()
	defer b.mu.Unlock()

	stmCountBefore := b.Hippocampus.GetSTMCount()
	ltmCountBefore := b.Hippocampus.GetLTMCount()

	// 海馬: 記憶の固定化
	b.Hippocampus.SleepAndConsolidate()

	ltmCountAfter := b.Hippocampus.GetLTMCount()

	consolidated := ltmCountAfter - ltmCountBefore
	forgotten := stmCountBefore - consolidated

	return SleepResult{
		ConsolidatedCount: consolidated,
		ForgottenCount:    forgotten,
		STMCount:          b.Hippocampus.GetSTMCount(),
		LTMCount:          b.Hippocampus.GetLTMCount(),
	}
}

// UpdateMotivation はフィードバックにより意欲を更新
// 【神経科学的意味】外部からの報酬/罰により大脳基底核の意欲を更新
// 【処理内容】報酬予測誤差(RPE)に基づいて意欲レベルを調整
func (b *Brain) UpdateMotivation(isPositive bool) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	var reward float64
	if isPositive {
		reward = 100.0
	} else {
		reward = 0.0
	}

	b.BasalGanglia.UpdateMotivation(reward)
	return b.BasalGanglia.GetMotivation()
}

// ApplyStress はストレスを適用
// 【神経科学的意味】外部ストレス要因により前頭前皮質の理性値を低下
// 【処理内容】PFCの理性値を減少させ、感情制御能力を低下させる
func (b *Brain) ApplyStress(stressLevel int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.PFC.ApplyStress(stressLevel)
}

// Rest は休息を適用
// 【神経科学的意味】休息により前頭前皮質の理性値を回復
// 【処理内容】PFCの理性値を増加させ、感情制御能力を向上させる
func (b *Brain) Rest(restQuality int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.PFC.Rest(restQuality)
}

// GetState は現在の脳の状態を取得
// 【役割】脳の主要パラメータ（意欲、理性、記憶数）を返す
// 【用途】デバッグやモニタリング用
func (b *Brain) GetState() BrainState {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return BrainState{
		Motivation:      b.BasalGanglia.GetMotivation(),
		MotivationLevel: b.BasalGanglia.GetMotivationLevel(),
		Sanity:          b.PFC.GetSanity(),
		SanityLevel:     b.PFC.GetSanityLevel(),
		STMCount:        b.Hippocampus.GetSTMCount(),
		LTMCount:        b.Hippocampus.GetLTMCount(),
	}
}

// GetRecentMemories は直近の記憶を取得
// 【役割】Hippocampusから最近の記憶を取得
// 【用途】文脈理解や性格傾向計算に使用
func (b *Brain) GetRecentMemories() []models.RuneMemory {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.Hippocampus.GetRecentContext()
}

// SleepResult は睡眠処理の結果
// 【用途】睡眠処理でどれだけの記憶が固定化/忘却されたかを報告
type SleepResult struct {
	ConsolidatedCount int // LTMに固定化された記憶数
	ForgottenCount    int // 忘却された記憶数
	STMCount          int // 残っている短期記憶数
	LTMCount          int // 総長期記憶数
}

// BrainState は脳の状態
// 【用途】現在の脳の主要パラメータを表現
type BrainState struct {
	Motivation      int    // 意欲値 (0-100)
	MotivationLevel string // 意欲レベル (文字列表現)
	Sanity          int    // 理性値 (0-100)
	SanityLevel     string // 理性レベル (文字列表現)
	STMCount        int    // 短期記憶数
	LTMCount        int    // 長期記憶数
}
