// brain.go: 脳の各部位（扁桃体、海馬、大脳基底核など）を統合し、意識のメインループを制御するコアモジュール
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umekku/mind-os/internal/config"
	"github.com/umekku/mind-os/internal/core"
	"github.com/umekku/mind-os/internal/models"
)

// SensoryRequest は感覚入力リクエストの構造体
type SensoryRequest struct {
	Type        string `json:"type" validate:"omitempty,oneof=chat physical"` // "chat" or "physical"
	Text        string `json:"text" binding:"required" validate:"required,max=500"`
	SignalValue int    `json:"signalValue" validate:"min=-100,max=100"` // -100 to 100
}

// FeedbackRequest はフィードバックリクエストの構造体
type FeedbackRequest struct {
	Positive bool `json:"positive"`
}

// StressRequest はストレスリクエストの構造体
type StressRequest struct {
	Level int `json:"level" binding:"required" validate:"min=0,max=100"`
}

// RestRequest は休息リクエストの構造体
type RestRequest struct {
	Quality int `json:"quality" binding:"required" validate:"min=0,max=100"`
}

// MindStateResponse はマインドステートレスポンスの構造体
type MindStateResponse struct {
	CurrentReaction []models.EmotionValue `json:"currentReaction"`
	MoodStability   float64               `json:"moodStability"`
	PersonalityBias []models.EmotionValue `json:"personalityBias"`
	Motivation      float64               `json:"motivation"`
	Sanity          float64               `json:"sanity"`
}

// BrainHandler は脳統合APIハンドラー
type BrainHandler struct {
	brain *core.Brain
}

// NewBrainHandler は新しい BrainHandler を作成
func NewBrainHandler(cfg *config.Config) *BrainHandler {
	return &BrainHandler{
		brain: core.New(cfg),
	}
}

// Feedback はフィードバックを処理
// POST /api/v1/feedback
// [神経科学] 報酬系（VTA-NAc回路）へのドーパミン入力をシミュレートし、行動に対する強化あるいは罰を与えます。
// 予測誤差に基づく学習（強化学習）の基礎となるメカニズムです。
// @Summary      Process Feedback
// @Description  報酬系への入力をシミュレートし、意欲を更新します。
// @Tags         brain
// @Accept       json
// @Produce      json
// @Param        input  body      handlers.FeedbackRequest  true  "Feedback"
// @Success      200    {object}  models.SuccessResponse
// @Failure      400    {object}  models.ProblemDetails
// @Router       /api/v1/feedback [post]
func (h *BrainHandler) Feedback(c *gin.Context) {
	var req FeedbackRequest
	if err := BindStrict(c, &req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}

	motivation := h.brain.UpdateMotivation(req.Positive)

	SuccessResponse(c, gin.H{
		"message":    "Feedback processed",
		"positive":   req.Positive,
		"motivation": motivation,
	})
}

// GetRecentMemories は直近の記憶を取得
// GET /api/v1/memories
// [神経科学] 短期記憶(STM)として海馬に保持されているエピソード記憶のバッファへのアクセスを提供します。
// これらはまだ長期記憶(LTM)への固定化(Consolidation)が完了していない、不安定で揮発性の高い記憶群です。
// @Summary      Get Recent Memories (STM)
// @Description  海馬に存在する短期記憶(STM)を取得します。長期固定化前のエピソード記憶です。
// @Tags         brain
// @Produce      json
// @Success      200    {object}  models.SuccessResponse
// @Router       /api/v1/memories [get]
func (h *BrainHandler) GetRecentMemories(c *gin.Context) {
	memories := h.brain.GetRecentMemories()

	responses := make([]MemoryResponse, len(memories))
	for i, memory := range memories {
		emotionResponses := make([]models.EmotionValue, len(memory.Emotions))
		for j, emotion := range memory.Emotions {
			emotionResponses[j] = models.EmotionValue{
				Code:  models.EmotionCode(emotion.Code),
				Value: emotion.Value,
			}
		}

		responses[i] = MemoryResponse{
			UUID:      memory.UUID,
			Text:      memory.Text,
			Emotions:  emotionResponses,
			Weight:    memory.Weight,
			Type:      string(memory.Type),
			CreatedAt: memory.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Tags:      memory.Tags,
		}
	}

	SuccessResponse(c, gin.H{
		"memories": responses,
		"count":    len(responses),
	})
}
