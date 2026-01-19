package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umekku/mind-os/internal/models"
)

// ProcessSensory は感覚入力を処理
// POST /api/v1/sensory-inputs
// [神経科学] 視覚や聴覚などの感覚情報を視床(Thalamus)が受け取り、粗いフィルタリングを行った後、
// 扁桃体(Amygdala)での情動評価と海馬(Hippocampus)での文脈照合を経て、最終的な意識(MindState)を形成します。
// @Summary      Process Sensory Input
// @Description  感覚入力を受信し、脳内の感情・意欲・記憶システムを通じて処理し、マインドステートと応答テキストを返します。
// @Tags         brain
// @Accept       json
// @Produce      json
// @Param        input  body      handlers.SensoryRequest  true  "Sensory Input"
// @Success      200    {object}  models.SuccessResponse{mindState=models.MindStateResponse,debug=models.DebugInfo}
// @Failure      400    {object}  models.ProblemDetails
// @Router       /api/v1/sensory-inputs [post]
func (h *BrainHandler) ProcessSensory(c *gin.Context) {
	var req SensoryRequest
	if err := BindStrict(c, &req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}

	// モデルに変換
	input := models.SensoryInput{
		Type:        models.SignalType(req.Type),
		InputText:   req.Text,
		SignalValue: req.SignalValue,
	}
	// デフォルト値
	if input.Type == "" {
		input.Type = models.SignalChat
	}

	// 脳で処理
	mindState := h.brain.ProcessInput(input)

	// レスポンス用に変換
	currentReaction := make([]models.EmotionValue, len(mindState.CurrentReaction))
	for i, emotion := range mindState.CurrentReaction {
		currentReaction[i] = models.EmotionValue{
			Code:  models.EmotionCode(emotion.Code),
			Value: emotion.Value,
		}
	}

	personalityBias := make([]models.EmotionValue, len(mindState.PersonalityBias))
	for i, emotion := range mindState.PersonalityBias {
		personalityBias[i] = models.EmotionValue{
			Code:  models.EmotionCode(emotion.Code),
			Value: emotion.Value,
		}
	}

	resp := models.SuccessResponse{
		MindState: &models.MindStateResponse{
			CurrentReaction: currentReaction,
			MoodStability:   mindState.MoodStability,
			PersonalityBias: personalityBias,
			Motivation:      mindState.Motivation,
			Sanity:          mindState.Sanity,
			ReplyText:       mindState.ReplyText,
		},
		Reply: mindState.ReplyText,
		Debug: &models.DebugInfo{
			Cortisol:        mindState.Cortisol,
			Oxytocin:        mindState.Oxytocin,
			PredictedReward: mindState.PredictedReward,
			DaydreamLog:     mindState.DaydreamLog,
		},
	}
	SuccessResponse(c, resp)
}
