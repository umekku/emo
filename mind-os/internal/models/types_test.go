package models

import (
	"testing"
	"time"
)

// TestEmotionValue_Validate は EmotionValue のバリデーションをテスト
func TestEmotionValue_Validate(t *testing.T) {
	tests := []struct {
		name     string
		emotion  EmotionValue
		expected bool
	}{
		{
			name:     "有効な値: 0",
			emotion:  EmotionValue{Code: EmotionJoy, Value: 0},
			expected: true,
		},
		{
			name:     "有効な値: 50",
			emotion:  EmotionValue{Code: EmotionJoy, Value: 50},
			expected: true,
		},
		{
			name:     "有効な値: 100",
			emotion:  EmotionValue{Code: EmotionJoy, Value: 100},
			expected: true,
		},
		{
			name:     "無効な値: -1",
			emotion:  EmotionValue{Code: EmotionJoy, Value: -1},
			expected: false,
		},
		{
			name:     "無効な値: 101",
			emotion:  EmotionValue{Code: EmotionJoy, Value: 101},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.emotion.Validate()
			if result != tt.expected {
				t.Errorf("Validate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsValidEmotionCode は感情コードの有効性チェックをテスト
func TestIsValidEmotionCode(t *testing.T) {
	tests := []struct {
		name     string
		code     EmotionCode
		expected bool
	}{
		{"Joy", EmotionJoy, true},
		{"Surprise", EmotionSurprise, true},
		{"Anger", EmotionAnger, true},
		{"Fear", EmotionFear, true},
		{"Love", EmotionLove, true},
		{"Disgust", EmotionDisgust, true},
		{"Hope", EmotionHope, true},
		{"Neutral", EmotionNeutral, true},
		{"Invalid", EmotionCode("X"), false},
		{"Empty", EmotionCode(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmotionCode(tt.code)
			if result != tt.expected {
				t.Errorf("IsValidEmotionCode(%v) = %v, want %v", tt.code, result, tt.expected)
			}
		})
	}
}

// TestEmotionMap_Conversion は EmotionMap の変換をテスト
func TestEmotionMap_Conversion(t *testing.T) {
	// EmotionMap -> EmotionValue[]
	em := EmotionMap{
		EmotionJoy:  80,
		EmotionLove: 60,
		EmotionHope: 40,
	}

	values := em.ToEmotionValues()
	if len(values) != 3 {
		t.Errorf("ToEmotionValues() length = %d, want 3", len(values))
	}

	// EmotionValue[] -> EmotionMap
	converted := FromEmotionValues(values)
	if len(converted) != 3 {
		t.Errorf("FromEmotionValues() length = %d, want 3", len(converted))
	}

	// 値の一致確認
	for code, value := range em {
		if converted[code] != value {
			t.Errorf("Value mismatch for %v: got %d, want %d", code, converted[code], value)
		}
	}
}

// TestRuneMemory_Structure は RuneMemory の構造をテスト
func TestRuneMemory_Structure(t *testing.T) {
	now := time.Now()
	memory := RuneMemory{
		UUID: "test-uuid-123",
		Text: "テスト記憶",
		Emotions: []EmotionValue{
			{Code: EmotionJoy, Value: 70},
			{Code: EmotionHope, Value: 50},
		},
		Weight:     0.8,
		Type:       MemoryLTM,
		CreatedAt:  now,
		LastAccess: now,
		Tags:       []string{"test", "memory"},
	}

	if memory.UUID != "test-uuid-123" {
		t.Errorf("UUID = %v, want test-uuid-123", memory.UUID)
	}
	if memory.Text != "テスト記憶" {
		t.Errorf("Text = %v, want テスト記憶", memory.Text)
	}
	if memory.Type != MemoryLTM {
		t.Errorf("Type = %v, want LTM", memory.Type)
	}
	if memory.Weight != 0.8 {
		t.Errorf("Weight = %v, want 0.8", memory.Weight)
	}
	if len(memory.Emotions) != 2 {
		t.Errorf("Emotions length = %d, want 2", len(memory.Emotions))
	}
	if len(memory.Tags) != 2 {
		t.Errorf("Tags length = %d, want 2", len(memory.Tags))
	}
	if memory.CreatedAt != now {
		t.Errorf("CreatedAt = %v, want %v", memory.CreatedAt, now)
	}
	if memory.LastAccess != now {
		t.Errorf("LastAccess = %v, want %v", memory.LastAccess, now)
	}
}

// TestMindStateResponse_Structure は MindStateResponse の構造をテスト
func TestMindStateResponse_Structure(t *testing.T) {
	response := MindStateResponse{
		CurrentReaction: []EmotionValue{
			{Code: EmotionJoy, Value: 85},
		},
		MoodStability: 0.75,
		PersonalityBias: []EmotionValue{
			{Code: EmotionHope, Value: 60},
			{Code: EmotionLove, Value: 55},
		},
		Motivation: 0.80,
		Sanity:     0.90,
	}

	if len(response.CurrentReaction) != 1 {
		t.Errorf("CurrentReaction length = %d, want 1", len(response.CurrentReaction))
	}
	if response.MoodStability != 0.75 {
		t.Errorf("MoodStability = %v, want 0.75", response.MoodStability)
	}
	if len(response.PersonalityBias) != 2 {
		t.Errorf("PersonalityBias length = %d, want 2", len(response.PersonalityBias))
	}
	if response.Motivation != 0.80 {
		t.Errorf("Motivation = %v, want 0.80", response.Motivation)
	}
	if response.Sanity != 0.90 {
		t.Errorf("Sanity = %v, want 0.90", response.Sanity)
	}
}
