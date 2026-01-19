// ganglia.go: 報酬予測誤差(RPE)に基づいてドーパミン量を調整し、行動意欲(Motivation)を生成する大脳基底核モジュール
package basal

import (
	"sync"
)

// BasalGanglia は大脳基底核モジュール - ドーパミンによる意欲と行動強化を管理
// 報酬予測誤差 (RPE) モデルを採用し、期待値との差分で学習する
type BasalGanglia struct {
	mu              sync.RWMutex
	Motivation      float64 // 意欲レベル (0-100)
	PredictedReward float64 // 期待報酬値 (0-100)

	// 内部パラメータ
	decayRate     float64 // 自然減衰率
	minMotivation float64 // 最小意欲値
	maxMotivation float64 // 最大意欲値
}

// New は新しい BasalGanglia インスタンスを作成
func New() *BasalGanglia {
	return &BasalGanglia{
		Motivation:      50.0, // 初期値: 中立
		PredictedReward: 50.0, // 初期期待値: 中立
		decayRate:       0.95, // 自然減衰率 (5%)
		minMotivation:   0.0,
		maxMotivation:   100.0,
	}
}

// UpdateMotivation は報酬予測誤差に基づいて意欲と期待値を更新
// 【神経科学的意味】
// ドーパミンニューロンの発火パターン(RPE)をシミュレートする。
// 予期せぬ報酬（ポジティブRPE）はドーパミンを放出し、意欲を高め、行動を強化する。
// 期待外れ（ネガティブRPE）はドーパミンを抑制し、意欲を低下させる。
//
// actualReward: 実際に得られた快感・報酬値 (0-100)
func (bg *BasalGanglia) UpdateMotivation(actualReward float64) {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	// 報酬予測誤差 (Reward Prediction Error, RPE) の計算
	// 【数式】δ = R_actual - V_predicted
	// δ > 0: 期待以上の結果(Supprise!) -> ドーパミン放出
	// δ < 0: 期待外れ(Disappointment) -> ドーパミン抑制
	predictionError := actualReward - bg.PredictedReward

	// 1. 意欲(Motivation/Dopamine)の更新
	// 【数式】M_{t+1} = M_t + α × δ
	// α = 0.5 (感度係数)
	bg.Motivation += predictionError * 0.5

	// 2. 期待値(Value Function)の更新
	// 【数式】V_{t+1} = V_t + β × δ
	// β = 0.3 (学習率)
	// 将来の予測を現実に近づける（TD学習的な振る舞い）
	bg.PredictedReward += predictionError * 0.3

	bg.clampValues()
}

// GetMotivation は現在の意欲レベルを返す
// 外部I/F互換のため int で返す
func (bg *BasalGanglia) GetMotivation() int {
	bg.mu.RLock()
	defer bg.mu.RUnlock()
	return int(bg.Motivation)
}

// GetPredictedReward は現在の予測報酬値を返す
func (bg *BasalGanglia) GetPredictedReward() float64 {
	bg.mu.RLock()
	defer bg.mu.RUnlock()
	return bg.PredictedReward
}

// SetMotivation は意欲レベルを直接設定（テスト用）
func (bg *BasalGanglia) SetMotivation(value int) {
	bg.mu.Lock()
	defer bg.mu.Unlock()
	bg.Motivation = float64(value)
	bg.clampValues()
}

// ApplyDecay は時間経過による自然減衰を適用
func (bg *BasalGanglia) ApplyDecay() {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	// 50（中立）に向かって減衰
	if bg.Motivation > 50 {
		diff := bg.Motivation - 50.0
		bg.Motivation = 50.0 + (diff * bg.decayRate)
	} else if bg.Motivation < 50 {
		diff := 50.0 - bg.Motivation
		bg.Motivation = 50.0 - (diff * bg.decayRate)
	}

	bg.clampValues()
}

// GetMotivationLevel は意欲レベルを文字列で返す
func (bg *BasalGanglia) GetMotivationLevel() string {
	bg.mu.RLock()
	defer bg.mu.RUnlock()

	m := bg.Motivation
	switch {
	case m >= 80:
		return "very_high"
	case m >= 60:
		return "high"
	case m >= 40:
		return "normal"
	case m >= 20:
		return "low"
	default:
		return "very_low"
	}
}

// ShouldTakeAction は意欲レベルに基づいて行動を起こすべきか判定
func (bg *BasalGanglia) ShouldTakeAction(threshold int) bool {
	bg.mu.RLock()
	defer bg.mu.RUnlock()
	return int(bg.Motivation) >= threshold
}

// RewardFromEmotion は感情値から意欲を更新
// 感情値 (0-100) をそのまま報酬として扱う
// Float64に変更して互換性を確保
func (bg *BasalGanglia) RewardFromEmotion(emotionValue float64) {
	// 内部でLockを取るのでここでは取らない
	bg.UpdateMotivation(emotionValue)
}

// clampValues は値を範囲内に制限（内部用）
func (bg *BasalGanglia) clampValues() {
	if bg.Motivation < bg.minMotivation {
		bg.Motivation = bg.minMotivation
	}
	if bg.Motivation > bg.maxMotivation {
		bg.Motivation = bg.maxMotivation
	}

	// 期待値を0-100に制限
	if bg.PredictedReward < 0 {
		bg.PredictedReward = 0
	}
	if bg.PredictedReward > 100 {
		bg.PredictedReward = 100
	}
}

// Reset は意欲を初期状態にリセット
func (bg *BasalGanglia) Reset() {
	bg.mu.Lock()
	defer bg.mu.Unlock()
	bg.Motivation = 50.0
	bg.PredictedReward = 50.0
}
