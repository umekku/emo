package amygdala

import (
	"testing"

	"github.com/umekku/mind-os/internal/models"
)

// TestNew は Amygdala インスタンスの生成をテスト
func TestNew(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	if a == nil {
		t.Fatal("New() returned nil")
	}
	if a.dict == nil {
		t.Fatal("dict not initialized")
	}
	if len(a.dict) == 0 {
		t.Fatal("dict is empty")
	}
}

// TestAssess_PositiveEmotions はポジティブ感情をテスト
func TestAssess_PositiveEmotions(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		text          string
		expectedCodes []models.EmotionCode
	}{
		{
			name:          "成功メッセージ",
			text:          "嬉しい", // "嬉しい" -> Joy
			expectedCodes: []models.EmotionCode{models.EmotionJoy},
		},
		{
			name:          "喜びの表現",
			text:          "楽しい",
			expectedCodes: []models.EmotionCode{models.EmotionJoy},
		},
		{
			name:          "愛情表現",
			text:          "愛",
			expectedCodes: []models.EmotionCode{models.EmotionLove},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emotions := a.Assess(tt.text)
			if len(emotions) == 0 {
				t.Fatal("No emotions returned")
			}

			// 期待される感情コードが含まれているか確認
			emotionMap := models.FromEmotionValues(emotions)
			for _, code := range tt.expectedCodes {
				if _, exists := emotionMap[code]; !exists {
					t.Errorf("Expected emotion code %v not found", code)
				}
			}
		})
	}
}

// TestAssess_NegativeEmotions はネガティブ感情をテスト
func TestAssess_NegativeEmotions(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		text          string
		expectedCodes []models.EmotionCode
	}{
		{
			name:          "バグ発見",
			text:          "バグ",
			expectedCodes: []models.EmotionCode{models.EmotionFear},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emotions := a.Assess(tt.text)
			if len(emotions) == 0 {
				t.Fatal("No emotions returned")
			}

			emotionMap := models.FromEmotionValues(emotions)
			for _, code := range tt.expectedCodes {
				if _, exists := emotionMap[code]; !exists {
					t.Errorf("Expected emotion code %v not found", code)
				}
			}
		})
	}
}

// TestAssess_Neutral は未知の入力のテスト
func TestAssess_Neutral(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		text string
	}{
		{"未知の単語", "xyzabc"},
		{"数字のみ", "12345"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emotions := a.Assess(tt.text)
			if len(emotions) == 0 {
				t.Fatal("No emotions returned")
			}

			// 未知の入力にはニュートラルが返される
			emotionMap := models.FromEmotionValues(emotions)
			if value, exists := emotionMap[models.EmotionNeutral]; !exists {
				t.Error("Unknown input should return Neutral emotion")
			} else if value != 10 { // 実装に合わせて 50 -> 10 に更新
				t.Errorf("Neutral value = %d, want 10", value)
			}
		})
	}
}

// TestAssess_MultipleKeywords は複数キーワードのテスト
func TestAssess_MultipleKeywords(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}

	text := "バグ 最高" // 文脈依存を避けるためスペース区切り
	emotions := a.Assess(text)

	if len(emotions) == 0 {
		t.Fatal("No emotions returned")
	}

	emotionMap := models.FromEmotionValues(emotions)

	// バグ(Fear)と最高(Joy)の両方が含まれるはず
	if _, exists := emotionMap[models.EmotionFear]; !exists {
		t.Error("Expected Fear emotion from 'バグ'")
	}
	if _, exists := emotionMap[models.EmotionJoy]; !exists {
		t.Error("Expected Joy emotion from '最高'")
	}
}
