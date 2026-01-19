package hypothalamus

import (
	"math"
	"time"
)

// 概日リズムの定数
const (
	DayTimeStart   = 6  // 日中開始 (6:00)
	NightTimeStart = 22 // 夜間開始 (22:00)
)

// UpdateCircadianRhythm は現在時刻に基づいて概日リズムホルモンを更新
// 時間帯によってMelatonin（睡眠）とSerotonin（覚醒）のバランスを調整
func (h *Homeostasis) UpdateCircadianRhythm(currentTime time.Time) {
	h.mu.Lock()
	defer h.mu.Unlock()

	hour := currentTime.Hour()

	// 夜間 (22:00 - 06:00)
	if hour >= NightTimeStart || hour < DayTimeStart {
		// メラトニン上昇（睡眠促進）
		// 深夜2時前後で最大値
		h.Melatonin = calculateNightMelatonin(hour)

		// セロトニン低下（覚醒度低下）
		h.Serotonin = 20.0

	} else {
		// 日中 (06:00 - 22:00)
		// メラトニン低下（覚醒）
		h.Melatonin = 10.0

		// セロトニン上昇（覚醒・安定）
		// 正午前後で最大値
		h.Serotonin = calculateDaySerotonin(hour)
	}

	h.clamp()
}

// calculateNightMelatonin は夜間のメラトニン値を計算
// 深夜2時前後（2-2時）で最大値 (80-100)
func calculateNightMelatonin(hour int) float64 {
	// 22時以降、翌日0時に向けて上昇
	if hour >= NightTimeStart {
		// 22時: 0.8, 23時: 0.9, 24時:0時: 1.0のイメージ
		progress := float64(hour-NightTimeStart) / 2.0
		return 80.0 + (progress * 20.0)
	}

	// 0-6時の朝に向けて下降
	// 0時: 1.0, 3時: 0.9, 6時: 0.5
	progress := float64(DayTimeStart-hour) / float64(DayTimeStart)
	return 50.0 + (progress * 50.0)
}

// calculateDaySerotonin は日中のセロトニン値を計算
// 正午前後（10-14時）で最大値 (60-100)
func calculateDaySerotonin(hour int) float64 {
	// 朝6時から正午に向けて上昇
	if hour >= DayTimeStart && hour < 12 {
		// 6時: 60, 12時: 100
		progress := float64(hour-DayTimeStart) / 6.0
		return 60.0 + (progress * 40.0)
	}

	// 正午から夜22時に向けて緩やかに下降
	if hour >= 12 && hour < NightTimeStart {
		// 12時: 100, 22時: 60
		progress := float64(hour-12) / float64(NightTimeStart-12)
		return 100.0 - (progress * 40.0)
	}

	return 60.0 // デフォルト
}

// GetCircadianEffects は概日リズムによる効果を返す
// motivationCap: 意欲の最大値制限 (0.5-1.0)
// emotionalSensitivity: 感情感度の係数 (1.0-1.2)
// cortisolDecayBoost: ストレス減衰の加速率 (1.0-2.0)
func (h *Homeostasis) GetCircadianEffects() (motivationCap float64, emotionalSensitivity float64, cortisolDecayBoost float64) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// メラトニンが高い（夜間）ほど意欲にキャップ
	// Melatonin 0 -> cap 1.0, Melatonin 100 -> cap 0.5
	motivationCap = 1.0 - (h.Melatonin / 200.0)
	motivationCap = math.Max(0.5, motivationCap)

	// メラトニンが高い（夜間）ほど感情的になる（感傷的）
	// Melatonin 0 -> 1.0x, Melatonin 100 -> 1.2x
	emotionalSensitivity = 1.0 + (h.Melatonin / 500.0)

	// セロトニンが高い（日中）ほどストレス回復が早い
	// Serotonin 0 -> 1.0x, Serotonin 100 -> 2.0x
	cortisolDecayBoost = 1.0 + (h.Serotonin / 100.0)

	return
}
