package hypothalamus

import (
	"math"
	"sync"
	"time"
)

// Homeostasis は生体の恒常性を管理する構造体
// 【神経科学的意味】視床下部におけるホルモンバランスの維持
// 【役割】ストレス(Cortisol)と愛着(Oxytocin)の拮抗作用、および時間経過による減衰をシミュレート
type Homeostasis struct {
	mu          sync.RWMutex
	Cortisol    float64   // ストレスホルモン (0-100): 高いと不機嫌になりやすい、闘争・逃走反応に関連
	Oxytocin    float64   // 愛着ホルモン (0-100): 高いと精神が安定する、社会的結束に関連
	Melatonin   float64   // 睡眠ホルモン (0-100): 夜間に上昇、概日リズムに関与
	Serotonin   float64   // 覚醒・安心ホルモン (0-100): 日中に上昇、気分調整に関与
	LastUpdated time.Time // 最終更新時間

	// テスト用の時間プロバイダー
	// 実環境では time.Now() を使用するが、テスト時に時間を固定できるようにする
	TimeProvider func() time.Time
}

// NewHomeostasis は新しい Homeostasis インスタンスを作成
func NewHomeostasis() *Homeostasis {
	return &Homeostasis{
		Cortisol:     0,
		Oxytocin:     0,
		Melatonin:    0,
		Serotonin:    50,
		LastUpdated:  time.Now(),
		TimeProvider: time.Now, // デフォルト: システム時間
	}
}

// Update は外部刺激により値を変動させる
// stressor: 負の刺激 (Cortisol増加)
// affection: 正の刺激 (Oxytocin増加, Cortisol減少)
// 【メカニズム】正の刺激はオキシトシンを分泌させ、同時にコルチゾールを抑制する（拮抗作用）
func (h *Homeostasis) Update(stressor float64, affection float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := h.TimeProvider()

	if stressor > 0 {
		h.Cortisol += stressor
	}

	if affection > 0 {
		h.Oxytocin += affection
		// 愛着によるストレス緩和効果 (Buffering effect of social support)
		// ストレスホルモンを減少させる
		h.Cortisol -= affection
	}

	h.clamp()
	h.LastUpdated = now
}

// Decay は時間経過による自然減衰を計算する
// 【神経科学的意味】ホルモンの血中濃度の半減期をシミュレート
func (h *Homeostasis) Decay() {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := h.TimeProvider()
	elapsed := now.Sub(h.LastUpdated).Hours()

	// 経過時間が負または0なら何もしない
	if elapsed <= 0 {
		return
	}

	// 1時間あたり-10 程度 (半減期ではなく線形減衰を採用)
	decayRate := 10.0
	decayAmount := decayRate * elapsed

	if h.Cortisol > 0 {
		h.Cortisol -= decayAmount
		if h.Cortisol < 0 {
			h.Cortisol = 0
		}
	}

	if h.Oxytocin > 0 {
		h.Oxytocin -= decayAmount
		if h.Oxytocin < 0 {
			h.Oxytocin = 0
		}
	}

	h.LastUpdated = now
}

// GetStatus は現在の値を返す
func (h *Homeostasis) GetStatus() (float64, float64) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.Cortisol, h.Oxytocin
}

// clamp は値を 0-100 の範囲に収める
func (h *Homeostasis) clamp() {
	h.Cortisol = math.Max(0, math.Min(100, h.Cortisol))
	h.Oxytocin = math.Max(0, math.Min(100, h.Oxytocin))
	h.Melatonin = math.Max(0, math.Min(100, h.Melatonin))
	h.Serotonin = math.Max(0, math.Min(100, h.Serotonin))
}
