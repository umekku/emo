package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umekku/mind-os/internal/amygdala"
	"github.com/umekku/mind-os/internal/hippocampus"
	"github.com/umekku/mind-os/internal/models"
)

// MemoryRequest は記憶追加リクエストの構造体
type MemoryRequest struct {
	Text string `json:"text" binding:"required" validate:"required,max=1000"`
}

// MemoryResponse は記憶レスポンスの構造体
type MemoryResponse struct {
	UUID      string                `json:"uuid"`
	Text      string                `json:"text"`
	Emotions  []models.EmotionValue `json:"emotions"`
	Weight    float64               `json:"weight"`
	Type      string                `json:"type"`
	CreatedAt string                `json:"createdAt"` // consistent with other models
	Tags      []string              `json:"tags"`
}

// MemoryStatsResponse は記憶統計レスポンスの構造体
type MemoryStatsResponse struct {
	STMCount int `json:"stmCount"`
	LTMCount int `json:"ltmCount"`
}

// MemoryHandler は記憶管理ハンドラー
type MemoryHandler struct {
	hippocampus *hippocampus.Hippocampus
	amygdala    *amygdala.Amygdala
}

// NewMemoryHandler は新しい MemoryHandler を作成
// Note: DB接続は現在MemoryHandlerでは使用しない(STMのみ)か、将来的に注入する
func NewMemoryHandler() *MemoryHandler {
	am, _ := amygdala.New()
	return &MemoryHandler{
		hippocampus: hippocampus.New(nil), // DBなし（STMのみモード）
		amygdala:    am,                   // 記憶の感情付けに使用
	}
}

// AddMemory は新しいエピソード記憶を追加
// POST /api/v1/memories
// [神経科学] 入力された体験をエピソード記憶として海馬(Hippocampus)にエンコードします。
// 同時に扁桃体(Amygdala)による感情評価を行い、記憶に感情的な重み付け（サリエンス）を付与します。
// @Summary      Create New Memory
// @Description  新しい記憶を作成し、感情評価を行って短期記憶(STM)として保存します。
// @Tags         brain
// @Accept       json
// @Produce      json
// @Param        input  body      handlers.MemoryRequest  true  "Memory Content"
// @Success      201    {object}  handlers.MemoryResponse
// @Failure      400    {object}  models.ProblemDetails
// @Router       /api/v1/memories [post]
func (h *MemoryHandler) AddMemory(c *gin.Context) {
	var req MemoryRequest
	if err := BindStrict(c, &req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}

	// 1. 感情分析
	emotions := h.amygdala.Assess(req.Text)

	// 2. 記憶として追加
	memory := h.hippocampus.AddMemory(req.Text, emotions)

	// レスポンス構築
	// モデル変換
	emotionResponses := make([]models.EmotionValue, len(memory.Emotions))
	for i, e := range memory.Emotions {
		emotionResponses[i] = models.EmotionValue{Code: e.Code, Value: e.Value}
	}

	resp := MemoryResponse{
		UUID:      memory.UUID,
		Text:      memory.Text,
		Emotions:  emotionResponses,
		Weight:    memory.Weight,
		Type:      string(memory.Type),
		CreatedAt: memory.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Tags:      memory.Tags,
	}

	c.JSON(http.StatusCreated, resp)
}

// GetRecentMemories は直近の記憶を取得
// GET /api/v1/memory/recent
func (h *MemoryHandler) GetRecentMemories(c *gin.Context) {
	memories := h.hippocampus.GetRecentContext()

	// レスポンス変換
	responses := make([]MemoryResponse, len(memories))
	for i, m := range memories {
		emotionResponses := make([]models.EmotionValue, len(m.Emotions))
		for j, e := range m.Emotions {
			emotionResponses[j] = models.EmotionValue{Code: e.Code, Value: e.Value}
		}
		responses[i] = MemoryResponse{
			UUID:      m.UUID,
			Text:      m.Text,
			Emotions:  emotionResponses,
			Weight:    m.Weight,
			Type:      string(m.Type),
			CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Tags:      m.Tags,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// GetMemoryStats は記憶の統計情報を取得
// GET /api/v1/memory/stats
func (h *MemoryHandler) GetMemoryStats(c *gin.Context) {
	c.JSON(http.StatusOK, MemoryStatsResponse{
		STMCount: h.hippocampus.GetSTMCount(),
		LTMCount: h.hippocampus.GetLTMCount(),
	})
}

// Sleep は睡眠処理（固定化）を実行
// POST /api/v1/memory/sleep
func (h *MemoryHandler) Sleep(c *gin.Context) {
	h.hippocampus.SleepAndConsolidate()
	c.JSON(http.StatusOK, gin.H{"message": "Sleep cycle completed", "stmCount": h.hippocampus.GetSTMCount(), "ltmCount": h.hippocampus.GetLTMCount()})
}
