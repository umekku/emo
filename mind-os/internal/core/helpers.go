package core

import (
	"github.com/umekku/mind-os/internal/models"
)

// addEmotion は感情リストに値を追加または加算する
// 【処理内容】既存の感情コードがあれば加算、なければ新規追加
// 【上限処理】感情値は100を超えないようにクリッピング
func addEmotion(emotions *[]models.EmotionValue, code models.EmotionCode, value int) {
	for i, e := range *emotions {
		if e.Code == code {
			(*emotions)[i].Value += value
			if (*emotions)[i].Value > 100 {
				(*emotions)[i].Value = 100
			}
			return
		}
	}
	// 新規追加
	if value > 100 {
		value = 100
	}
	*emotions = append(*emotions, models.EmotionValue{Code: code, Value: value})
}

// generateMindState はマインドステートレスポンスを生成
// 【役割】現在の脳の状態を統合してクライアント向けレスポンスを作成
// 【処理内容】性格傾向、気分安定度、ホルモン状態、概日リズム効果を統合
func (b *Brain) generateMindState(emotions []models.EmotionValue) models.MindStateResponse {
	// 性格傾向を計算（直近の記憶から）
	personalityBias := b.calculatePersonalityBias()

	// 気分の安定度を計算
	moodStability := b.calculateMoodStability(emotions)

	// ホルモン状態と予測報酬を取得
	cortisol, oxytocin := b.Hypothalamus.GetStatus()
	predictedReward := b.BasalGanglia.GetPredictedReward()

	// 概日リズムの効果を取得
	motivationCap, _, _ := b.Hypothalamus.GetCircadianEffects()

	// 意欲値を取得し、概日リズムによるキャップを適用
	rawMotivation := float64(b.BasalGanglia.GetMotivation()) / 100.0
	cappedMotivation := rawMotivation * motivationCap

	return models.MindStateResponse{
		CurrentReaction: emotions,
		MoodStability:   moodStability,
		PersonalityBias: personalityBias,
		Motivation:      cappedMotivation,
		Sanity:          float64(b.PFC.GetSanity()) / 100.0,
		Cortisol:        cortisol,
		Oxytocin:        oxytocin,
		PredictedReward: predictedReward,
	}
}

// calculatePersonalityBias は性格傾向を計算
// 【神経科学的意味】長期的かつ反復的な記憶パターンから形成される性格特性
// 【アルゴリズム】直近の記憶の感情を平均化し、閾値以上のものを性格傾向とする
func (b *Brain) calculatePersonalityBias() []models.EmotionValue {
	memories := b.Hippocampus.GetRecentContext()

	if len(memories) == 0 {
		return []models.EmotionValue{
			{Code: models.EmotionNeutral, Value: 50},
		}
	}

	// 感情の累積
	emotionSum := make(map[models.EmotionCode]int)
	emotionCount := make(map[models.EmotionCode]int)

	for _, memory := range memories {
		for _, emotion := range memory.Emotions {
			emotionSum[emotion.Code] += emotion.Value
			emotionCount[emotion.Code]++
		}
	}

	// 平均を計算
	bias := make([]models.EmotionValue, 0)
	for code, sum := range emotionSum {
		avg := sum / emotionCount[code]
		if avg > 30 { // 閾値以上のみ
			bias = append(bias, models.EmotionValue{
				Code:  code,
				Value: avg,
			})
		}
	}

	if len(bias) == 0 {
		return []models.EmotionValue{
			{Code: models.EmotionNeutral, Value: 50},
		}
	}

	return bias
}

// calculateMoodStability は気分の安定度を計算
// 【神経科学的意味】感情の変動の小ささ。安定しているほど予測可能な反応
// 【アルゴリズム】感情値の分散を計算し、分散が小さいほど安定度が高い
// 【数式】Stability = 1.0 - (Variance / 10000.0)
func (b *Brain) calculateMoodStability(currentEmotions []models.EmotionValue) float64 {
	if len(currentEmotions) == 0 {
		return 0.5
	}

	// 感情値の平均を計算
	total := 0
	for _, emotion := range currentEmotions {
		total += emotion.Value
	}
	avg := float64(total) / float64(len(currentEmotions))

	variance := 0.0
	for _, emotion := range currentEmotions {
		diff := float64(emotion.Value) - avg
		variance += diff * diff
	}
	variance /= float64(len(currentEmotions))

	// 分散が小さいほど安定
	// 0-100の二乗を0.0-1.0に正規化
	stability := 1.0 - (variance / 10000.0)
	if stability < 0.0 {
		stability = 0.0
	}
	if stability > 1.0 {
		stability = 1.0
	}

	return stability
}
