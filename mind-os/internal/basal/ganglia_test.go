package basal

import (
	"sync"
	"testing"
)

// TestNew は BasalGanglia インスタンスの生成をテスト
func TestNew(t *testing.T) {
	bg := New()
	if bg == nil {
		t.Fatal("New() returned nil")
	}
	if bg.Motivation != 50 {
		t.Errorf("Initial Motivation = %f, want 50", bg.Motivation)
	}
}

// TestGetMotivation は意欲取得をテスト
func TestGetMotivation(t *testing.T) {
	bg := New()
	motivation := bg.GetMotivation()
	if motivation != 50 {
		t.Errorf("GetMotivation() = %d, want 50", motivation)
	}
}

// TestUpdateMotivation_Positive はポジティブフィードバックをテスト
func TestUpdateMotivation_Positive(t *testing.T) {
	bg := New()
	initialMotivation := bg.GetMotivation()

	bg.UpdateMotivation(100.0) // 期待値50に対して報酬100 -> RPE=50 -> Motivation+=25

	newMotivation := bg.GetMotivation()
	if newMotivation <= initialMotivation {
		t.Errorf("Motivation after positive feedback = %d, should be greater than %d",
			newMotivation, initialMotivation)
	}
}

// TestUpdateMotivation_Negative はネガティブフィードバックをテスト
func TestUpdateMotivation_Negative(t *testing.T) {
	bg := New()
	initialMotivation := bg.GetMotivation()

	bg.UpdateMotivation(0.0) // 期待値50に対して報酬0 -> RPE=-50 -> Motivation-=25

	newMotivation := bg.GetMotivation()
	if newMotivation >= initialMotivation {
		t.Errorf("Motivation after negative feedback = %d, should be less than %d",
			newMotivation, initialMotivation)
	}
}

// TestUpdateMotivation_MaxLimit は上限をテスト
func TestUpdateMotivation_MaxLimit(t *testing.T) {
	bg := New()
	bg.SetMotivation(95)

	// 複数回のポジティブフィードバック
	for i := 0; i < 5; i++ {
		bg.UpdateMotivation(100.0)
	}

	motivation := bg.GetMotivation()
	if motivation > 100 {
		t.Errorf("Motivation = %d, should not exceed 100", motivation)
	}
	if motivation != 100 {
		t.Errorf("Motivation = %d, want 100 (max)", motivation)
	}
}

// TestUpdateMotivation_MinLimit は下限をテスト
func TestUpdateMotivation_MinLimit(t *testing.T) {
	bg := New()
	bg.SetMotivation(10)

	// 複数回ネガティブフィードバック
	for i := 0; i < 5; i++ {
		bg.UpdateMotivation(0.0)
	}

	motivation := bg.GetMotivation()
	if motivation < 0 {
		t.Errorf("Motivation = %d, should not be less than 0", motivation)
	}
	if motivation != 0 {
		t.Errorf("Motivation = %d, want 0 (min)", motivation)
	}
}

// TestSetMotivation は意欲設定をテスト
func TestSetMotivation(t *testing.T) {
	bg := New()

	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"Normal value", 75, 75},
		{"Max value", 100, 100},
		{"Min value", 0, 0},
		{"Over max", 150, 100},
		{"Under min", -50, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bg.SetMotivation(tt.value)
			if bg.GetMotivation() != tt.expected {
				t.Errorf("SetMotivation(%d) = %d, want %d",
					tt.value, bg.GetMotivation(), tt.expected)
			}
		})
	}
}

// TestApplyDecay は自然減衰をテスト
func TestApplyDecay(t *testing.T) {
	tests := []struct {
		name           string
		initialValue   int
		expectDecrease bool
	}{
		{"High motivation decays", 90, true},
		{"Low motivation recovers", 20, false},
		{"Neutral stays", 50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bg := New()
			bg.SetMotivation(tt.initialValue)
			initialMotivation := bg.GetMotivation()

			bg.ApplyDecay()

			newMotivation := bg.GetMotivation()

			if tt.expectDecrease {
				if newMotivation >= initialMotivation {
					t.Errorf("Expected decrease: %d -> %d", initialMotivation, newMotivation)
				}
			} else if tt.initialValue < 50 {
				if newMotivation <= initialMotivation {
					t.Errorf("Expected increase: %d -> %d", initialMotivation, newMotivation)
				}
			}
		})
	}
}

// TestGetMotivationLevel は意欲レベル文字列をテスト
func TestGetMotivationLevel(t *testing.T) {
	tests := []struct {
		motivation int
		expected   string
	}{
		{90, "very_high"},
		{70, "high"},
		{50, "normal"},
		{30, "low"},
		{10, "very_low"},
	}

	bg := New()
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			bg.SetMotivation(tt.motivation)
			level := bg.GetMotivationLevel()
			if level != tt.expected {
				t.Errorf("GetMotivationLevel() = %v, want %v", level, tt.expected)
			}
		})
	}
}

// TestShouldTakeAction は行動判定をテスト
func TestShouldTakeAction(t *testing.T) {
	bg := New()

	tests := []struct {
		motivation int
		threshold  int
		expected   bool
	}{
		{70, 50, true},
		{50, 50, true},
		{30, 50, false},
		{80, 90, false},
	}

	for _, tt := range tests {
		bg.SetMotivation(tt.motivation)
		result := bg.ShouldTakeAction(tt.threshold)
		if result != tt.expected {
			t.Errorf("ShouldTakeAction(%d) with motivation %d = %v, want %v",
				tt.threshold, tt.motivation, result, tt.expected)
		}
	}
}

// TestRewardFromEmotion は感情からの報酬計算をテスト
func TestRewardFromEmotion(t *testing.T) {
	tests := []struct {
		name           string
		emotionValue   float64 // Changed to float64
		expectIncrease bool
	}{
		{"High positive emotion", 90.0, true},
		// Neutral(50)の場合、期待値50と一致 -> RPE=0 -> 変動なし
		{"Neutral emotion", 50.0, false},
		{"Negative emotion", 20.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bg := New()
			initialMotivation := bg.GetMotivation()

			bg.RewardFromEmotion(tt.emotionValue)

			newMotivation := bg.GetMotivation()

			if tt.expectIncrease {
				if newMotivation <= initialMotivation {
					t.Errorf("Expected increase: %d -> %d", initialMotivation, newMotivation)
				}
			} else if tt.emotionValue < 50 {
				if newMotivation >= initialMotivation {
					// 20 < 50 なので RPE < 0 -> motivation decrease
					t.Errorf("Expected decrease or stay: %d -> %d", initialMotivation, newMotivation)
				}
			}
		})
	}
}

// TestReset はリセットをテスト
func TestReset(t *testing.T) {
	bg := New()
	bg.SetMotivation(80)

	bg.Reset()

	if bg.GetMotivation() != 50 {
		t.Errorf("After Reset, Motivation = %d, want 50", bg.GetMotivation())
	}
}

// TestConcurrency はスレッドセーフ性をテスト
func TestConcurrency(t *testing.T) {
	bg := New()
	var wg sync.WaitGroup

	// 複数のゴルーチンで同時にアクセス
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(isPositive bool) {
			defer wg.Done()
			val := 0.0
			if isPositive {
				val = 100.0
			}
			bg.UpdateMotivation(val)
			bg.GetMotivation()
		}(i%2 == 0)
	}

	wg.Wait()

	// パニックせずに完了すればOK
	motivation := bg.GetMotivation()
	if motivation < 0 || motivation > 100 {
		t.Errorf("Motivation = %d, out of valid range after concurrent access", motivation)
	}
}

// TestMultipleFeedbacks は複数フィードバックの累積をテスト
func TestMultipleFeedbacks(t *testing.T) {
	bg := New()
	initialMotivation := bg.GetMotivation()

	// 3回ポジティブ
	for i := 0; i < 3; i++ {
		bg.UpdateMotivation(100.0)
	}

	afterPositive := bg.GetMotivation()
	if afterPositive <= initialMotivation {
		t.Errorf("After 3 positive feedbacks, motivation should increase")
	}

	// 2回ネガティブ
	for i := 0; i < 2; i++ {
		bg.UpdateMotivation(0.0)
	}

	afterNegative := bg.GetMotivation()
	if afterNegative >= afterPositive {
		t.Errorf("After negative feedbacks, motivation should decrease")
	}
}
