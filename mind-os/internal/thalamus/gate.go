package thalamus

import (
	"strings"
	"sync"

	"github.com/umekku/mind-os/internal/models"
)

// Thalamus は視床モジュール - 感覚入力のフィルタリングと順応を管理
// 繰り返しの刺激に対して応答を減衰させ（順応）、新しい刺激を優先する
type Thalamus struct {
	mu sync.RWMutex

	// 順応管理
	LastInputText   string  // 直前の入力テキスト
	RepetitionCount int     // 同じ入力が連続した回数
	SatiationLevel  float64 // 飽和度 (0.0-1.0) 高いほど新しい刺激を求めている

	// 内部パラメータ
	similarityThreshold float64 // 類似判定の閾値 (0.0-1.0)
}

// New は新しい Thalamus インスタンスを作成
func New() *Thalamus {
	return &Thalamus{
		LastInputText:       "",
		RepetitionCount:     0,
		SatiationLevel:      0.5, // 初期値: 中立
		similarityThreshold: 0.8, // 80%以上の類似度で「同じ」と判定
	}
}

// Filter は入力信号の強度係数 (Gain) を計算
// 繰り返し入力に対して順応（慣れ）を適用し、反応を減衰させる
func (t *Thalamus) Filter(input models.SensoryInput) (float64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	currentText := input.InputText

	// 1. 重複チェック（類似度判定）
	isSimilar := t.checkSimilarity(currentText, t.LastInputText)

	if isSimilar {
		// 同じ入力が繰り返されている
		t.RepetitionCount++
	} else {
		// 新しい入力
		t.RepetitionCount = 0
		t.LastInputText = currentText
	}

	// 2. ゲイン計算（順応による減衰）
	// 基本ゲイン: 1.0
	// 繰り返しが多いほど減衰: Gain = 1.0 / (1.0 + 0.5 * RepetitionCount)
	// 例: 0回目 -> 1.0, 1回目 -> 0.67, 2回目 -> 0.5, 3回目 -> 0.4
	gain := 1.0 / (1.0 + 0.5*float64(t.RepetitionCount))

	// 3. 飽和度の更新
	// 繰り返しが多いと飽和度が上がる（新しい刺激を求める）
	t.SatiationLevel = float64(t.RepetitionCount) / 10.0
	if t.SatiationLevel > 1.0 {
		t.SatiationLevel = 1.0
	}

	return gain, nil
}

// checkSimilarity は2つのテキストの類似度を判定
// 簡易実装: 完全一致または高い部分一致で true
func (t *Thalamus) checkSimilarity(text1, text2 string) bool {
	// 空文字列の場合は類似していないとみなす
	if text1 == "" || text2 == "" {
		return false
	}

	// 完全一致チェック
	if text1 == text2 {
		return true
	}

	// 正規化（小文字化、空白除去）
	norm1 := strings.ToLower(strings.TrimSpace(text1))
	norm2 := strings.ToLower(strings.TrimSpace(text2))

	if norm1 == norm2 {
		return true
	}

	// 部分一致チェック（長い方の文字列に対する短い方の包含率）
	shorter, longer := norm1, norm2
	if len(norm1) > len(norm2) {
		shorter, longer = norm2, norm1
	}

	// 短い方が長い方に含まれているか
	if strings.Contains(longer, shorter) {
		// 長さの比率で類似度を判定
		ratio := float64(len(shorter)) / float64(len(longer))
		return ratio >= t.similarityThreshold
	}

	return false
}

// GetSatiationLevel は現在の飽和度を返す
func (t *Thalamus) GetSatiationLevel() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.SatiationLevel
}

// GetRepetitionCount は現在の繰り返し回数を返す
func (t *Thalamus) GetRepetitionCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.RepetitionCount
}

// Reset は状態をリセット
func (t *Thalamus) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.LastInputText = ""
	t.RepetitionCount = 0
	t.SatiationLevel = 0.5
}
