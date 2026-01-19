package thalamus

import (
	"testing"

	"github.com/umekku/mind-os/internal/models"
)

// TestNew は初期化をテスト
func TestNew(t *testing.T) {
	th := New()
	if th == nil {
		t.Fatal("New() returned nil")
	}
	if th.SatiationLevel != 0.5 {
		t.Errorf("Initial SatiationLevel = %f, want 0.5", th.SatiationLevel)
	}
}

// TestFilter_NovelInput は新しい入力の処理をテスト
func TestFilter_NovelInput(t *testing.T) {
	th := New()
	input := models.SensoryInput{InputText: "こんにちは"}

	gain, err := th.Filter(input)
	if err != nil {
		t.Fatalf("Filter error: %v", err)
	}

	if gain != 1.0 {
		t.Errorf("Gain for novel input = %f, want 1.0", gain)
	}
	if th.RepetitionCount != 0 {
		t.Errorf("RepetitionCount = %d, want 0", th.RepetitionCount)
	}
}

// TestFilter_RepeatedInput は繰り返し入力の減衰をテスト
func TestFilter_RepeatedInput(t *testing.T) {
	th := New()
	input := models.SensoryInput{InputText: "繰り返し"}

	// 1回目
	th.Filter(input)

	// 2回目 (繰り返し1回)
	gain, _ := th.Filter(input)
	if gain >= 1.0 {
		t.Errorf("Gain should decrease for repeated input, got %f", gain)
	}
	if th.RepetitionCount != 1 {
		t.Errorf("RepetitionCount = %d, want 1", th.RepetitionCount)
	}

	// 3回目 (繰り返し2回)
	gain2, _ := th.Filter(input)
	if gain2 >= gain {
		t.Errorf("Gain should decrease further, %f -> %f", gain, gain2)
	}
}

// TestCheckSimilarity は類似性判定をテスト
func TestCheckSimilarity(t *testing.T) {
	th := New()

	tests := []struct {
		text1    string
		text2    string
		expected bool
	}{
		{"Hello", "Hello", true},
		{"Hello", "hello", true},          // 大文字小文字無視
		{"Hello ", "hello", true},         // 空白無視
		{"abcdefg", "abcd", false},        // 部分一致率が低い(4/7 < 0.8)
		{"abcdefghij", "abcdefghi", true}, // 部分一致率が高い(9/10 >= 0.8)
		{"", "hello", false},
	}

	for _, tt := range tests {
		result := th.checkSimilarity(tt.text1, tt.text2)
		if result != tt.expected {
			t.Errorf("checkSimilarity(%q, %q) = %v, want %v", tt.text1, tt.text2, result, tt.expected)
		}
	}
}
