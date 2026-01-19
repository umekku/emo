package pfc

import (
	"sync"
	"testing"

	"github.com/umekku/mind-os/internal/models"
)

// TestNew は PrefrontalCortex インスタンスの生成をテスト
func TestNew(t *testing.T) {
	pfc := New()
	if pfc == nil {
		t.Fatal("New() returned nil")
	}
	if pfc.Sanity != 80 {
		t.Errorf("Initial Sanity = %d, want 80", pfc.Sanity)
	}
}

// TestGetSanity は理性値取得をテスト
func TestGetSanity(t *testing.T) {
	pfc := New()
	sanity := pfc.GetSanity()
	if sanity != 80 {
		t.Errorf("GetSanity() = %d, want 80", sanity)
	}
}

// TestSetSanity は理性値設定をテスト
func TestSetSanity(t *testing.T) {
	pfc := New()

	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"Normal value", 60, 60},
		{"Max value", 100, 100},
		{"Min value", 0, 0},
		{"Over max", 150, 100},
		{"Under min", -50, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfc.SetSanity(tt.value)
			if pfc.GetSanity() != tt.expected {
				t.Errorf("SetSanity(%d) = %d, want %d",
					tt.value, pfc.GetSanity(), tt.expected)
			}
		})
	}
}

// TestArbitrate_HighSanity は高理性時の感情調整をテスト
func TestArbitrate_HighSanity(t *testing.T) {
	pfc := New()
	pfc.SetSanity(100)

	raw := []models.EmotionValue{
		{Code: models.EmotionAnger, Value: 80},
		{Code: models.EmotionFear, Value: 70},
		{Code: models.EmotionJoy, Value: 60},
	}

	// Arbitrate(raw, 0, 0)
	arbitrated := pfc.Arbitrate(raw, 0.0, 0.0)

	// ネガティブ感情が減衰されているか
	for i, emotion := range arbitrated {
		switch emotion.Code {
		case models.EmotionAnger, models.EmotionFear:
			if emotion.Value >= raw[i].Value {
				t.Errorf("Negative emotion %v not suppressed: %d >= %d",
					emotion.Code, emotion.Value, raw[i].Value)
			}
		case models.EmotionJoy:
			if emotion.Value <= raw[i].Value {
				t.Errorf("Positive emotion %v not boosted: %d <= %d",
					emotion.Code, emotion.Value, raw[i].Value)
			}
		}
	}
}

// TestArbitrate_LowSanity は低理性時の感情調整をテスト
func TestArbitrate_LowSanity(t *testing.T) {
	pfc := New()
	pfc.SetSanity(20)

	raw := []models.EmotionValue{
		{Code: models.EmotionAnger, Value: 80},
		{Code: models.EmotionFear, Value: 70},
	}

	arbitrated := pfc.Arbitrate(raw, 0.0, 0.0)

	// 低理性時はそのまま通る
	for i, emotion := range arbitrated {
		if emotion.Value != raw[i].Value {
			t.Errorf("Low sanity should not modify emotions: %d != %d",
				emotion.Value, raw[i].Value)
		}
	}
}

// TestArbitrate_NeutralEmotions は中立感情の処理をテスト
func TestArbitrate_NeutralEmotions(t *testing.T) {
	pfc := New()
	pfc.SetSanity(80)

	raw := []models.EmotionValue{
		{Code: models.EmotionNeutral, Value: 50},
		{Code: models.EmotionSurprise, Value: 60},
	}

	arbitrated := pfc.Arbitrate(raw, 0.0, 0.0)

	// 中立感情は変更されない
	for i, emotion := range arbitrated {
		if emotion.Value != raw[i].Value {
			t.Errorf("Neutral emotion should not be modified: %d != %d",
				emotion.Value, raw[i].Value)
		}
	}
}

// TestArbitrate_HighCortisol はストレス過多時の感情増幅をテスト
func TestArbitrate_HighCortisol(t *testing.T) {
	pfc := New()
	pfc.SetSanity(80) // 理性は高いが...

	raw := []models.EmotionValue{
		{Code: models.EmotionAnger, Value: 50},
	}

	// Cortisol 100 (Max) -> 増幅1.5倍
	// ただしSanity80による抑制もかかる
	// Suppression: 0.8 * (1.0 - (100-50)/100) = 0.8 * 0.5 = 0.4
	// Boost: 1.5
	// Value: 50 * 1.5 = 75
	// Reduction: 75 * 0.4 * 0.5 = 15
	// Final: 75 - 15 = 60
	// 元の50より増えるはず

	arbitrated := pfc.Arbitrate(raw, 100.0, 0.0)

	if len(arbitrated) == 0 {
		t.Fatal("Arbitrated is empty")
	}

	if arbitrated[0].Value <= raw[0].Value {
		t.Errorf("High cortisol should boost anger: %d <= %d",
			arbitrated[0].Value, raw[0].Value)
	}
}

// TestArbitrate_HighOxytocin は愛着過多時の感情変換をテスト
func TestArbitrate_HighOxytocin(t *testing.T) {
	pfc := New()
	pfc.SetSanity(80)

	raw := []models.EmotionValue{
		{Code: models.EmotionAnger, Value: 80},
	}

	// Oxytocin 100 (Max) -> 変換確率50%
	// 確実なテストは難しいが、何度か試行してGriefが出るか確認
	griefCount := 0
	loopCount := 100

	for i := 0; i < loopCount; i++ {
		res := pfc.Arbitrate(raw, 0.0, 100.0)
		if res[0].Code == models.EmotionGrief {
			griefCount++
		}
	}

	if griefCount == 0 {
		t.Errorf("Anger should be converted to Grief with high oxytocin (0/%d)", loopCount)
	}
	// 確率的テストなので、あまり厳密にしすぎない（10%期待）
	if griefCount < 10 {
		t.Errorf("Conversion rate too low: %d/%d", griefCount, loopCount)
	}
}

// TestUpdateSanity は理性値更新をテスト
func TestUpdateSanity(t *testing.T) {
	pfc := New()
	initialSanity := pfc.GetSanity()

	pfc.UpdateSanity(10)
	if pfc.GetSanity() != initialSanity+10 {
		t.Errorf("UpdateSanity(10) = %d, want %d",
			pfc.GetSanity(), initialSanity+10)
	}

	pfc.UpdateSanity(-20)
	if pfc.GetSanity() != initialSanity-10 {
		t.Errorf("UpdateSanity(-20) = %d, want %d",
			pfc.GetSanity(), initialSanity-10)
	}
}

// TestApplyStress はストレス適用をテスト
func TestApplyStress(t *testing.T) {
	pfc := New()
	initialSanity := pfc.GetSanity()

	pfc.ApplyStress(50)

	newSanity := pfc.GetSanity()
	if newSanity >= initialSanity {
		t.Errorf("Stress should decrease sanity: %d >= %d",
			newSanity, initialSanity)
	}
}

// TestRest は休息による回復をテスト
func TestRest(t *testing.T) {
	pfc := New()
	pfc.SetSanity(50)
	initialSanity := pfc.GetSanity()

	pfc.Rest(80)

	newSanity := pfc.GetSanity()
	if newSanity <= initialSanity {
		t.Errorf("Rest should increase sanity: %d <= %d",
			newSanity, initialSanity)
	}
}

// TestGetSanityLevel は理性レベル文字列をテスト
func TestGetSanityLevel(t *testing.T) {
	tests := []struct {
		sanity   int
		expected string
	}{
		{90, "very_high"},
		{70, "high"},
		{50, "normal"},
		{30, "low"},
		{10, "very_low"},
	}

	pfc := New()
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			pfc.SetSanity(tt.sanity)
			level := pfc.GetSanityLevel()
			if level != tt.expected {
				t.Errorf("GetSanityLevel() = %v, want %v", level, tt.expected)
			}
		})
	}
}

// TestCanControlEmotions は感情制御可能性をテスト
func TestCanControlEmotions(t *testing.T) {
	pfc := New()

	tests := []struct {
		sanity   int
		expected bool
	}{
		{50, true},
		{30, true},
		{29, false},
		{0, false},
	}

	for _, tt := range tests {
		pfc.SetSanity(tt.sanity)
		result := pfc.CanControlEmotions()
		if result != tt.expected {
			t.Errorf("CanControlEmotions() with sanity %d = %v, want %v",
				tt.sanity, result, tt.expected)
		}
	}
}

// TestGetSuppressionRate は抑制率取得をテスト
func TestGetSuppressionRate(t *testing.T) {
	pfc := New()

	tests := []struct {
		sanity      int
		minExpected float64
		maxExpected float64
	}{
		{100, 1.0, 1.0},
		{50, 0.5, 0.5},
		{20, 0.0, 0.0},
		{0, 0.0, 0.0},
	}

	for _, tt := range tests {
		pfc.SetSanity(tt.sanity)
		rate := pfc.GetSuppressionRate()
		if rate < tt.minExpected || rate > tt.maxExpected {
			t.Errorf("GetSuppressionRate() with sanity %d = %v, want between %v and %v",
				tt.sanity, rate, tt.minExpected, tt.maxExpected)
		}
	}
}

// TestCalculateEmotionalImpact は感情影響度計算をテスト
func TestCalculateEmotionalImpact(t *testing.T) {
	pfc := New()

	emotions := []models.EmotionValue{
		{Code: models.EmotionAnger, Value: 80},
		{Code: models.EmotionFear, Value: 60},
		{Code: models.EmotionJoy, Value: 70},
	}

	// 高理性時
	pfc.SetSanity(100)
	highSanityImpact := pfc.CalculateEmotionalImpact(emotions)

	// 低理性時
	pfc.SetSanity(20)
	lowSanityImpact := pfc.CalculateEmotionalImpact(emotions)

	// 高理性時の方が影響が小さいはず
	if highSanityImpact >= lowSanityImpact {
		t.Errorf("High sanity impact (%d) should be less than low sanity impact (%d)",
			highSanityImpact, lowSanityImpact)
	}
}

// TestReset はリセットをテスト
func TestReset(t *testing.T) {
	pfc := New()
	pfc.SetSanity(30)

	pfc.Reset()

	if pfc.GetSanity() != 80 {
		t.Errorf("After Reset, Sanity = %d, want 80", pfc.GetSanity())
	}
}

// TestConcurrency はスレッドセーフ性をテスト
func TestConcurrency(t *testing.T) {
	pfc := New()
	var wg sync.WaitGroup

	emotions := []models.EmotionValue{
		{Code: models.EmotionAnger, Value: 70},
	}

	// 複数のゴルーチンで同時にアクセス
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(delta int) {
			defer wg.Done()
			pfc.UpdateSanity(delta)
			pfc.GetSanity()
			pfc.Arbitrate(emotions, 0.0, 0.0)
		}(i%2*2 - 1) // -1 or +1
	}

	wg.Wait()
}

// TestArbitrate_AllEmotionTypes はすべての感情タイプをテスト
func TestArbitrate_AllEmotionTypes(t *testing.T) {
	pfc := New()
	pfc.SetSanity(80)

	raw := []models.EmotionValue{
		{Code: models.EmotionJoy, Value: 60},
		{Code: models.EmotionSurprise, Value: 50},
		{Code: models.EmotionAnger, Value: 70},
		{Code: models.EmotionFear, Value: 65},
		{Code: models.EmotionLove, Value: 55},
		{Code: models.EmotionDisgust, Value: 60},
		{Code: models.EmotionHope, Value: 58},
		{Code: models.EmotionNeutral, Value: 50},
	}

	arbitrated := pfc.Arbitrate(raw, 0.0, 0.0)

	if len(arbitrated) != len(raw) {
		t.Errorf("Arbitrated length = %d, want %d", len(arbitrated), len(raw))
	}

	// すべての感情が処理されていることを確認
	for i, emotion := range arbitrated {
		if emotion.Code != raw[i].Code {
			t.Errorf("Emotion code changed: %v != %v", emotion.Code, raw[i].Code)
		}
	}
}

// TestArbitrate_EmptyInput は空入力をテスト
func TestArbitrate_EmptyInput(t *testing.T) {
	pfc := New()
	raw := []models.EmotionValue{}

	arbitrated := pfc.Arbitrate(raw, 0.0, 0.0)

	if len(arbitrated) != 0 {
		t.Errorf("Arbitrated length = %d, want 0", len(arbitrated))
	}
}

// TestStressAndRest はストレスと休息のサイクルをテスト
func TestStressAndRest(t *testing.T) {
	pfc := New()
	initialSanity := pfc.GetSanity()

	// ストレス適用
	pfc.ApplyStress(80)
	afterStress := pfc.GetSanity()

	if afterStress >= initialSanity {
		t.Error("Stress should decrease sanity")
	}

	// 休息
	pfc.Rest(100)
	afterRest := pfc.GetSanity()

	if afterRest <= afterStress {
		t.Error("Rest should increase sanity")
	}
}
