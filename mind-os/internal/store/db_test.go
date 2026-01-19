package store

import (
	"os"
	"testing"
	"time"

	"github.com/umekku/mind-os/internal/models"
)

func TestDB(t *testing.T) {
	dbPath := "test_mind.db"
	defer os.Remove(dbPath)

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	defer db.Close()

	// 1. SaveMemory
	memo := models.RuneMemory{
		UUID: "test-uuid-1",
		Text: "Test Memory",
		Emotions: []models.EmotionValue{
			{Code: models.EmotionJoy, Value: 80},
		},
		Weight:     0.8,
		Type:       models.MemoryLTM,
		CreatedAt:  time.Now(),
		LastAccess: time.Now(),
		Tags:       []string{"test"},
	}

	if err := db.SaveMemory(memo); err != nil {
		t.Errorf("SaveMemory failed: %v", err)
	}

	// 2. GetRecentMemories
	memos, err := db.GetRecentMemories(10)
	if err != nil {
		t.Errorf("GetRecentMemories failed: %v", err)
	}
	if len(memos) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(memos))
	}
	if memos[0].Text != "Test Memory" {
		t.Errorf("Content mismatch: %s", memos[0].Text)
	}

	// 3. GetLTMCount
	count, err := db.GetLTMCount()
	if err != nil {
		t.Errorf("GetLTMCount failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// 4. GetMemoryByUUID
	fetched, err := db.GetMemoryByUUID("test-uuid-1")
	if err != nil {
		t.Errorf("GetMemoryByUUID failed: %v", err)
	}
	if fetched == nil {
		t.Fatal("Memory not found")
	}
	if fetched.Weight != 0.8 {
		t.Errorf("Weight mismatch: %f", fetched.Weight)
	}

	// 5. DeleteOldMemories
	// 追加でいくつか保存
	for i := 0; i < 5; i++ {
		m := memo
		m.UUID = string(rune('a' + i)) // simple unique id mockup
		m.Weight = 0.1 * float64(i)
		db.SaveMemory(m)
	}

	// Total should be 1 + 5 = 6
	// 'a'...'e' + 'test-uuid-1'

	if err := db.DeleteOldMemories(3); err != nil { // Keep 3
		t.Errorf("DeleteOldMemories failed: %v", err)
	}

	newCount, _ := db.GetLTMCount()
	if newCount != 3 {
		t.Errorf("Expected 3 memories after delete, got %d", newCount)
	}
}
