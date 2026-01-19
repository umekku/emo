package hippocampus

import (
	"os"
	"testing"
	"time"

	"github.com/umekku/mind-os/internal/models"
	"github.com/umekku/mind-os/internal/store"
)

// setupTest はテスト用のHippocampusインスタンスとクリーンアップ関数を返す
func setupTest(t *testing.T) (*Hippocampus, func()) {
	// 一時ファイルを使うのがベストだが、簡易化のため固定ファイル
	dbPath := "test_hippocampus.db"
	// 既存があれば削除
	os.Remove(dbPath)

	db, err := store.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}

	h := New(db)

	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return h, cleanup
}

// TestNew は Hippocampus インスタンスの生成をテスト
func TestNew(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()

	if h == nil {
		t.Fatal("New() returned nil")
	}
	if h.STM == nil {
		t.Fatal("STM not initialized")
	}
	// LTM removed from struct, check store
	if h.store == nil {
		t.Fatal("Store not initialized")
	}
	if len(h.STM) != 0 {
		t.Errorf("Initial STM size = %d, want 0", len(h.STM))
	}
	if count := h.GetLTMCount(); count != 0 {
		t.Errorf("Initial LTM size = %d, want 0", count)
	}
}

// TestAddEpisode はエピソード追加をテスト
func TestAddEpisode(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()

	emotions := []models.EmotionValue{
		{Code: models.EmotionJoy, Value: 80},
		{Code: models.EmotionHope, Value: 60},
	}

	h.AddEpisode("テスト記憶", emotions)

	if len(h.STM) != 1 {
		t.Fatalf("STM size = %d, want 1", len(h.STM))
	}

	memory := h.STM[0]
	if memory.Text != "テスト記憶" {
		t.Errorf("Memory text = %v, want テスト記憶", memory.Text)
	}
	if memory.Type != models.MemorySTM {
		t.Errorf("Memory type = %v, want STM", memory.Type)
	}
	if len(memory.Emotions) != 2 {
		t.Errorf("Emotions count = %d, want 2", len(memory.Emotions))
	}
	if memory.UUID == "" {
		t.Error("UUID not generated")
	}
}

// TestAddEpisode_MultipleMemories は複数記憶の追加をテスト
func TestAddEpisode_MultipleMemories(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()

	for i := 0; i < 5; i++ {
		emotions := []models.EmotionValue{
			{Code: models.EmotionJoy, Value: 50 + i*10},
		}
		h.AddEpisode("記憶"+string(rune('A'+i)), emotions)
	}

	if len(h.STM) != 5 {
		t.Errorf("STM size = %d, want 5", len(h.STM))
	}
}

// TestAddEpisode_MaxSTMSize はSTMサイズ制限をテスト
func TestAddEpisode_MaxSTMSize(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()
	h.maxSTMSize = 10 // テスト用に小さく設定

	// 最大サイズを超えて追加
	for i := 0; i < 15; i++ {
		emotions := []models.EmotionValue{
			{Code: models.EmotionNeutral, Value: 50},
		}
		h.AddEpisode("記憶", emotions)
	}

	if len(h.STM) != h.maxSTMSize {
		t.Errorf("STM size = %d, want %d", len(h.STM), h.maxSTMSize)
	}
}

// TestGetRecentContext は直近記憶の取得をテスト
func TestGetRecentContext(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()

	// STMに記憶を追加
	for i := 0; i < 5; i++ {
		emotions := []models.EmotionValue{
			{Code: models.EmotionJoy, Value: 70},
		}
		h.AddEpisode("STM記憶", emotions)
		time.Sleep(10 * time.Millisecond) // 時間差をつける
	}

	// LTMに記憶を追加 (DB経由)
	for i := 0; i < 3; i++ {
		memory := models.RuneMemory{
			UUID:       "ltm-" + string(rune('A'+i)),
			Text:       "LTM記憶",
			Type:       models.MemoryLTM,
			Weight:     0.8,
			CreatedAt:  time.Now(),
			LastAccess: time.Now(),
			Tags:       []string{"ltm"},
		}
		// DBに直接保存
		h.store.SaveMemory(memory)
		time.Sleep(10 * time.Millisecond)
	}

	recent := h.GetRecentContext()

	if len(recent) != 8 {
		t.Errorf("Recent context size = %d, want 8", len(recent))
	}

	// 最新の記憶が最初に来ることを確認
	if len(recent) > 1 && recent[0].LastAccess.Before(recent[len(recent)-1].LastAccess) {
		t.Error("Recent context not sorted by LastAccess")
	}
}

// TestSleepAndConsolidate は睡眠処理をテスト
func TestSleepAndConsolidate(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()
	h.consolidationThreshold = 0.6

	// 高重みの記憶（LTMへ移行されるべき）
	highWeightEmotions := []models.EmotionValue{
		{Code: models.EmotionJoy, Value: 90},
		{Code: models.EmotionLove, Value: 80},
	}
	h.AddEpisode("重要な記憶", highWeightEmotions)

	// 低重みの記憶（忘却されるべき）
	lowWeightEmotions := []models.EmotionValue{
		{Code: models.EmotionNeutral, Value: 30},
	}
	h.AddEpisode("どうでもいい記憶", lowWeightEmotions)

	initialSTMCount := len(h.STM)
	if initialSTMCount != 2 {
		t.Fatalf("Initial STM count = %d, want 2", initialSTMCount)
	}

	// 睡眠処理
	h.SleepAndConsolidate()

	// STMはクリアされる
	if len(h.STM) != 0 {
		t.Errorf("STM count after sleep = %d, want 0", len(h.STM))
	}

	// 高重みの記憶のみLTMへ移行 (DBチェック)
	ltmCount, _ := h.store.GetLTMCount()
	if ltmCount != 1 {
		t.Errorf("LTM count after sleep = %d, want 1", ltmCount)
	}

	// DBから取得して確認
	memos, _ := h.store.GetRecentMemories(1)
	if len(memos) > 0 {
		ltmMemory := memos[0]
		if ltmMemory.Type != models.MemoryLTM {
			t.Errorf("LTM memory type = %v, want LTM", ltmMemory.Type)
		}
		if ltmMemory.Text != "重要な記憶" {
			t.Errorf("LTM memory text = %v, want 重要な記憶", ltmMemory.Text)
		}
	} else {
		t.Error("Failed to fetch LTM memory")
	}
}

// TestSleepAndConsolidate_AllForgotten はすべて忘却されるケースをテスト
func TestSleepAndConsolidate_AllForgotten(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()
	h.consolidationThreshold = 0.8 // 高い閾値

	// すべて低重みの記憶
	for i := 0; i < 5; i++ {
		emotions := []models.EmotionValue{
			{Code: models.EmotionNeutral, Value: 40},
		}
		h.AddEpisode("低重み記憶", emotions)
	}

	h.SleepAndConsolidate()

	if len(h.STM) != 0 {
		t.Errorf("STM count = %d, want 0", len(h.STM))
	}
	ltmCount, _ := h.store.GetLTMCount()
	if ltmCount != 0 {
		t.Errorf("LTM count = %d, want 0 (all forgotten)", ltmCount)
	}
}

// TestCalculateWeight は重み計算をテスト
func TestCalculateWeight(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name      string
		emotions  []models.EmotionValue
		minWeight float64
		maxWeight float64
	}{
		{
			name: "高感情値",
			emotions: []models.EmotionValue{
				{Code: models.EmotionJoy, Value: 90},
				{Code: models.EmotionLove, Value: 80},
			},
			minWeight: 0.8,
			maxWeight: 1.0,
		},
		{
			name: "低感情値",
			emotions: []models.EmotionValue{
				{Code: models.EmotionNeutral, Value: 30},
			},
			minWeight: 0.0,
			maxWeight: 0.5,
		},
		{
			name: "中感情値",
			emotions: []models.EmotionValue{
				{Code: models.EmotionJoy, Value: 50},
			},
			minWeight: 0.4,
			maxWeight: 0.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weight := h.calculateWeight(tt.emotions)
			if weight < tt.minWeight || weight > tt.maxWeight {
				t.Errorf("Weight = %v, want between %v and %v", weight, tt.minWeight, tt.maxWeight)
			}
		})
	}
}

// TestExtractTags はタグ抽出をテスト
func TestExtractTags(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()

	emotions := []models.EmotionValue{
		{Code: models.EmotionJoy, Value: 80},
		{Code: models.EmotionHope, Value: 60},
	}

	// 短いテキスト
	shortTags := h.extractTags("短い", emotions)
	hasShortTag := false
	for _, tag := range shortTags {
		if tag == "short" {
			hasShortTag = true
			break
		}
	}
	if !hasShortTag {
		t.Error("Short text should have 'short' tag")
	}

	// 長いテキスト
	longText := "これは非常に長いテキストです。" + string(make([]byte, 100))
	longTags := h.extractTags(longText, emotions)
	hasLongTag := false
	for _, tag := range longTags {
		if tag == "long" {
			hasLongTag = true
			break
		}
	}
	if !hasLongTag {
		t.Error("Long text should have 'long' tag")
	}

	// 感情コードタグ
	hasJoyTag := false
	for _, tag := range shortTags {
		if tag == "J" {
			hasJoyTag = true
			break
		}
	}
	if !hasJoyTag {
		t.Error("Should have Joy emotion tag")
	}
}

// TestGetSTMCount はSTMカウント取得をテスト
func TestGetSTMCount(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()

	if h.GetSTMCount() != 0 {
		t.Errorf("Initial STM count = %d, want 0", h.GetSTMCount())
	}

	h.AddEpisode("記憶1", []models.EmotionValue{{Code: models.EmotionJoy, Value: 70}})
	h.AddEpisode("記憶2", []models.EmotionValue{{Code: models.EmotionJoy, Value: 70}})

	if h.GetSTMCount() != 2 {
		t.Errorf("STM count = %d, want 2", h.GetSTMCount())
	}
}

// TestGetLTMCount はLTMカウント取得をテスト
func TestGetLTMCount(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()

	if h.GetLTMCount() != 0 {
		t.Errorf("Initial LTM count = %d, want 0", h.GetLTMCount())
	}

	// 高重みの記憶を追加して睡眠処理
	h.AddEpisode("重要記憶", []models.EmotionValue{
		{Code: models.EmotionJoy, Value: 90},
		{Code: models.EmotionLove, Value: 85},
	})
	h.SleepAndConsolidate()

	if h.GetLTMCount() != 1 {
		t.Errorf("LTM count = %d, want 1", h.GetLTMCount())
	}
}

// TestGetMemoryByUUID はUUID検索をテスト
func TestGetMemoryByUUID(t *testing.T) {
	h, cleanup := setupTest(t)
	defer cleanup()

	h.AddEpisode("テスト記憶", []models.EmotionValue{
		{Code: models.EmotionJoy, Value: 70},
	})

	uuid := h.STM[0].UUID

	// STMから検索
	found := h.GetMemoryByUUID(uuid)
	if found == nil {
		t.Fatal("Memory not found by UUID")
	}
	if found.UUID != uuid {
		t.Errorf("Found UUID = %v, want %v", found.UUID, uuid)
	}

	// 存在しないUUID
	notFound := h.GetMemoryByUUID("non-existent-uuid")
	if notFound != nil {
		t.Error("Should return nil for non-existent UUID")
	}
}
