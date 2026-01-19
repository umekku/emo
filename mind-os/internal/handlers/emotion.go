package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umekku/mind-os/internal/amygdala"
	"github.com/umekku/mind-os/internal/models"
)

// EmotionRequest は感情分析リクエストの構造体
type EmotionRequest struct {
	Text string `json:"text" binding:"required" validate:"required,max=1000"`
}

// EmotionResponse は感情分析レスポンスの構造体
type EmotionResponse struct {
	Text     string                `json:"text"`
	Emotions []models.EmotionValue `json:"emotions"`
}

// EmotionHandler は感情分析ハンドラー
type EmotionHandler struct {
	amygdala *amygdala.Amygdala
}

// NewEmotionHandler は新しい EmotionHandler を作成
func NewEmotionHandler() *EmotionHandler {
	am, err := amygdala.New()
	if err != nil {
		// 初期化失敗時はpanicせず、エラーを考慮すべきだが
		// ここでは簡易的にpanic、実際はログ出力してnilで返す等の処理が必要
		panic(err)
	}
	return &EmotionHandler{
		amygdala: am,
	}
}

// Assess はテキストから感情を分析
// POST /api/v1/emotions/assess
// [神経科学] 扁桃体(Amygdala)の機能に基づき、入力されたテキストパターンから
// 6つの基本感情（喜び・怒り・恐れ・悲しみ・嫌悪・驚き）および社会的感情（信頼・希望）を抽出します。
// [制御機構] 高次の前頭前野(PFC)によるトップダウン制御（抑制）はここでは適用されず、
// 純粋なボトムアップの情動反応（一次反応）を返します。
// @Summary      Assess Emotions (Amygdala)
// @Description  テキスト入力に対して扁桃体モジュールが生成する即時的な感情反応を分析します。PFCによる抑制前の「生の感情」です。
// @Tags         brain
// @Accept       json
// @Produce      json
// @Param        input  body      handlers.EmotionRequest  true  "Text Input"
// @Success      200    {object}  handlers.EmotionResponse
// @Failure      400    {object}  models.ProblemDetails
// @Router       /api/v1/emotions/assess [post]
func (h *EmotionHandler) Assess(c *gin.Context) {
	var req EmotionRequest
	if err := BindStrict(c, &req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}

	emotions := h.amygdala.Assess(req.Text)

	c.JSON(http.StatusOK, EmotionResponse{
		Text:     req.Text,
		Emotions: emotions,
	})
}
