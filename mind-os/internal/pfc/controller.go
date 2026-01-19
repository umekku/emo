package pfc

import (
	"math/rand"
	"sync"
	"time" // for random seed

	"github.com/umekku/mind-os/internal/models"
)

// PrefrontalCortex は前頭前皮質モジュール - 理性による感情抑制を管理
type PrefrontalCortex struct {
	mu     sync.RWMutex
	Sanity int // 理性値 (0-100)

	// 内部パラメータ
	minSanity int     // 最小理性値
	maxSanity int     // 最大理性値
	decayRate float64 // ストレスによる減衰率
}

// New は新しい PrefrontalCortex インスタンスを作成
func New() *PrefrontalCortex {
	// 乱数初期化（簡易的）
	rand.Seed(time.Now().UnixNano())

	return &PrefrontalCortex{
		Sanity:    80,   // 初期値: 高め
		minSanity: 0,    // 最小値
		maxSanity: 100,  // 最大値
		decayRate: 0.98, // ストレス減衰率
	}
}

// Arbitrate は理性とホルモンバランスによる感情の調整を行う
// cortisol: ストレスホルモン (0-100)
// oxytocin: 愛着ホルモン (0-100)
func (pfc *PrefrontalCortex) Arbitrate(raw []models.EmotionValue, cortisol float64, oxytocin float64) []models.EmotionValue {
	pfc.mu.RLock()
	defer pfc.mu.RUnlock()

	// Sanityが極端に低い場合はそのまま通す（暴走）
	if pfc.Sanity < 30 {
		return raw
	}

	// 調整後の感情値
	arbitrated := make([]models.EmotionValue, len(raw))
	copy(arbitrated, raw)

	// Sanityに基づく基本抑制率を計算 (0.0-1.0)
	// Sanity 100 -> 1.0 (完全抑制可能), Sanity 0 -> 0.0
	suppressionRate := float64(pfc.Sanity) / 100.0

	// ストレス過多 (Cortisol > 50) による理性の弱体化
	// ストレスが高いと抑制が効きにくくなる
	if cortisol > 50 {
		// 50を超えた分だけ抑制率を下げる
		// 例: Cortisol 100 -> suppressionRate *= 0.5 (半減)
		weakeningFactor := 1.0 - ((cortisol - 50.0) / 100.0)
		suppressionRate *= weakeningFactor
	}

	// 各感情を調整
	for i, emotion := range arbitrated {
		switch emotion.Code {
		case models.EmotionAnger, models.EmotionFear, models.EmotionDisgust, models.EmotionGrief:
			// 値の補正用変数
			currentValue := float64(emotion.Value)

			// 1. ストレス過多による増幅(イライラ)
			if cortisol > 50 {
				// 最大1.5倍に増幅
				// Cortisol 50 -> 1.0x, Cortisol 100 -> 1.5x
				boost := 1.0 + ((cortisol - 50.0) / 100.0)
				currentValue *= boost
			}

			// 2. 愛着過多 (Oxytocin > 50) による感情変換
			// 怒り(Anger) を悲しみ(Grief) に変換（攻撃を抑える）
			if emotion.Code == models.EmotionAnger && oxytocin > 50 {
				// 変換確率: Oxytocin 50 -> 0%, 100 -> 50%
				prob := (oxytocin - 50.0) / 100.0
				if rand.Float64() < prob {
					arbitrated[i].Code = models.EmotionGrief
					// 変換時の少し強度を抑える（怒りよりはマイルドに）
					currentValue *= 0.9
				}
			}

			// 3. 理性による抑制
			// ネガティブ感情を削減
			reduction := currentValue * suppressionRate * 0.5
			newValue := int(currentValue - reduction)

			// 閾値制限
			if newValue < 0 {
				newValue = 0
			}
			if newValue > 100 {
				newValue = 100
			}
			arbitrated[i].Value = newValue

		case models.EmotionJoy, models.EmotionLove, models.EmotionHope:
			// ポジティブ感情も軽く増幅
			boost := float64(emotion.Value) * suppressionRate * 0.1
			newValue := emotion.Value + int(boost)
			if newValue > 100 {
				newValue = 100
			}
			arbitrated[i].Value = newValue

		case models.EmotionSurprise, models.EmotionNeutral:
			// 中立的な感情はそのまま
		}
	}

	return arbitrated
}

// GetSanity は現在の理性値を返す
func (pfc *PrefrontalCortex) GetSanity() int {
	pfc.mu.RLock()
	defer pfc.mu.RUnlock()
	return pfc.Sanity
}

// SetSanity は理性値を設定
func (pfc *PrefrontalCortex) SetSanity(value int) {
	pfc.mu.Lock()
	defer pfc.mu.Unlock()
	pfc.Sanity = value
	pfc.clampSanity()
}

// UpdateSanity は理性値を増減
func (pfc *PrefrontalCortex) UpdateSanity(delta int) {
	pfc.mu.Lock()
	defer pfc.mu.Unlock()
	pfc.Sanity += delta
	pfc.clampSanity()
}

// ApplyStress はストレスにより理性値を減少
func (pfc *PrefrontalCortex) ApplyStress(stressLevel int) {
	pfc.mu.Lock()
	defer pfc.mu.Unlock()

	// ストレスレベル(0-100)に応じて理性を減少
	reduction := stressLevel / 5 // 最大20減少
	pfc.Sanity -= reduction
	pfc.clampSanity()
}

// Rest は休息により理性値を回復
func (pfc *PrefrontalCortex) Rest(restQuality int) {
	pfc.mu.Lock()
	defer pfc.mu.Unlock()

	// 休息の質(0-100)に応じて理性を回復
	recovery := restQuality / 4 // 最大25回復
	pfc.Sanity += recovery
	pfc.clampSanity()
}

// GetSanityLevel は理性レベルを文字列で返す
func (pfc *PrefrontalCortex) GetSanityLevel() string {
	pfc.mu.RLock()
	defer pfc.mu.RUnlock()

	switch {
	case pfc.Sanity >= 80:
		return "very_high"
	case pfc.Sanity >= 60:
		return "high"
	case pfc.Sanity >= 40:
		return "normal"
	case pfc.Sanity >= 20:
		return "low"
	default:
		return "very_low"
	}
}

// CanControlEmotions は感情を制御できるか判定
func (pfc *PrefrontalCortex) CanControlEmotions() bool {
	pfc.mu.RLock()
	defer pfc.mu.RUnlock()
	return pfc.Sanity >= 30
}

// GetSuppressionRate は抑制率を返す (0.0-1.0)
func (pfc *PrefrontalCortex) GetSuppressionRate() float64 {
	pfc.mu.RLock()
	defer pfc.mu.RUnlock()

	if pfc.Sanity < 30 {
		return 0.0
	}
	return float64(pfc.Sanity) / 100.0
}

// CalculateEmotionalImpact は感情の影響度を計算
func (pfc *PrefrontalCortex) CalculateEmotionalImpact(emotions []models.EmotionValue) int {
	pfc.mu.RLock()
	defer pfc.mu.RUnlock()

	if len(emotions) == 0 {
		return 0
	}

	// ネガティブ感情の合計
	negativeTotal := 0
	for _, emotion := range emotions {
		switch emotion.Code {
		case models.EmotionAnger, models.EmotionFear, models.EmotionDisgust:
			negativeTotal += emotion.Value
		}
	}

	// 理性値が高いほど影響が小さい
	suppressionRate := pfc.GetSuppressionRate()
	impact := int(float64(negativeTotal) * (1.0 - suppressionRate*0.5))

	return impact
}

// Reset は理性値を初期状態にリセット
func (pfc *PrefrontalCortex) Reset() {
	pfc.mu.Lock()
	defer pfc.mu.Unlock()
	pfc.Sanity = 80
}

// clampSanity は理性値を範囲内に制限（内部用）
func (pfc *PrefrontalCortex) clampSanity() {
	if pfc.Sanity < pfc.minSanity {
		pfc.Sanity = pfc.minSanity
	}
	if pfc.Sanity > pfc.maxSanity {
		pfc.Sanity = pfc.maxSanity
	}
}
