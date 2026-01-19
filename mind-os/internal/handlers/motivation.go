package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umekku/mind-os/internal/basal"
)

// MotivationRequest は意欲更新リクエストの構造体
type MotivationRequest struct {
	IsPositive bool `json:"isPositive"`
}

// EmotionRewardRequest は感情報酬リクエストの構造体
type EmotionRewardRequest struct {
	EmotionValue int `json:"emotionValue" binding:"required" validate:"required,min=0,max=100"`
}

// MotivationResponse は意欲レスポンスの構造体
type MotivationResponse struct {
	Motivation int    `json:"motivation"`
	Level      string `json:"level"`
}

// MotivationHandler は意欲管理ハンドラー
type MotivationHandler struct {
	basalGanglia *basal.BasalGanglia
}

// NewMotivationHandler は新しい MotivationHandler を作成
func NewMotivationHandler() *MotivationHandler {
	return &MotivationHandler{
		basalGanglia: basal.New(),
	}
}

// UpdateMotivation は意欲を直接更新
// POST /api/v1/motivation
// [神経科学] フィードバック（報酬/罰）に基づいて線条体でのドーパミン放出を調整します。
// ポジティブな結果は意欲を高め、ネガティブな結果は意欲を減退させます。
// @Summary      Update Motivation (Direct Feedback)
// @Description  外部要因（ユーザーからの直接的なフィードバックなど）に基づいて意欲レベルを更新します。
// @Tags         brain
// @Accept       json
// @Produce      json
// @Param        input  body      handlers.MotivationRequest  true  "Motivation Update"
// @Success      200    {object}  handlers.MotivationResponse
// @Failure      400    {object}  models.ProblemDetails
// @Router       /api/v1/motivation [post]
func (h *MotivationHandler) UpdateMotivation(c *gin.Context) {
	var req MotivationRequest
	if err := BindStrict(c, &req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}

	// bool -> reward value conversion
	reward := 0.0
	if req.IsPositive {
		reward = 100.0
	}
	h.basalGanglia.UpdateMotivation(reward)
	motivation := h.basalGanglia.GetMotivation()

	c.JSON(http.StatusOK, MotivationResponse{
		Motivation: int(motivation),
		Level:      h.basalGanglia.GetMotivationLevel(),
	})
}

// GetMotivation は現在の意欲を取得
// GET /api/v1/motivation
func (h *MotivationHandler) GetMotivation(c *gin.Context) {
	motivation := h.basalGanglia.GetMotivation()
	c.JSON(http.StatusOK, MotivationResponse{
		Motivation: int(motivation),
		Level:      h.basalGanglia.GetMotivationLevel(),
	})
}

// Reset は意欲をリセット
// POST /api/v1/motivation/reset
func (h *MotivationHandler) Reset(c *gin.Context) {
	h.basalGanglia.Reset()
	c.JSON(http.StatusOK, gin.H{"message": "Motivation reset"})
}

// RewardFromEmotion は感情価に基づき報酬系を刺激
// POST /api/v1/motivation/reward
// [神経科学] 扁桃体や眼窩前頭皮質から送られる感情シグナル（快/不快）を報酬予測誤差(RPE)として処理し、意欲を調整します。
// @Summary      Process Emotional Reward
// @Description  感情的なポジティブ/ネガティブな出来事を報酬として処理し、意欲レベルに反映させます。
// @Tags         brain
// @Accept       json
// @Produce      json
// @Param        input  body      handlers.EmotionRewardRequest  true  "Emotional Value"
// @Success      200    {object}  handlers.MotivationResponse
// @Failure      400    {object}  models.ProblemDetails
// @Router       /api/v1/motivation/reward [post]
func (h *MotivationHandler) RewardFromEmotion(c *gin.Context) {
	var req EmotionRewardRequest
	if err := BindStrict(c, &req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}

	h.basalGanglia.RewardFromEmotion(float64(req.EmotionValue))
	motivation := h.basalGanglia.GetMotivation()

	c.JSON(http.StatusOK, MotivationResponse{
		Motivation: int(motivation),
		Level:      h.basalGanglia.GetMotivationLevel(),
	})
}

// ApplyDecay は時間経過による意欲減衰を適用
// POST /api/v1/motivation/decay
// [神経科学] 刺激がない状態が続くと、ドーパミン受容体の感度低下やトニックドーパミンレベルの自然減少により、
// 意欲は徐々にベースラインに戻ろうとします。
// @Summary      Apply Motivation Decay
// @Description  時間経過による意欲の自然減衰をシミュレートします。定期的に呼び出されることを想定しています。
// @Tags         brain
// @Produce      json
// @Success      200    {object}  handlers.MotivationResponse
// @Router       /api/v1/motivation/decay [post]
func (h *MotivationHandler) ApplyDecay(c *gin.Context) {
	h.basalGanglia.ApplyDecay()
	motivation := h.basalGanglia.GetMotivation()

	c.JSON(http.StatusOK, MotivationResponse{
		Motivation: int(motivation),
		Level:      h.basalGanglia.GetMotivationLevel(),
	})
}
