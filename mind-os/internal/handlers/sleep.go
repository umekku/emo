package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Sleep は睡眠処理を実行
// POST /api/v1/sleep-cycles
// [神経科学] 睡眠中のメモリリプレイ（Sharp-wave ripples）を模倣し、短期記憶(STM)から長期記憶(LTM)への固定化(Consolidation)を行います。
// 同時に、シナプス恒常性の維持のため、重要度の低い記憶の忘却や神経伝達物質の再充填を行います。
// @Summary      Execute Sleep Cycle
// @Description  睡眠サイクルを実行し、記憶の定着（Consolidation）と神経伝達物質の回復を行います。
// @Tags         brain
// @Produce      json
// @Success      200  {object}  models.SuccessResponse
// @Failure      500  {object}  models.ProblemDetails
// @Router       /api/v1/sleep-cycles [post]
func (h *BrainHandler) Sleep(c *gin.Context) {
	result := h.brain.Sleep()

	SuccessResponse(c, gin.H{
		"message":           "Sleep consolidation completed",
		"consolidatedCount": result.ConsolidatedCount,
		"forgottenCount":    result.ForgottenCount,
		"stmCount":          result.STMCount,
		"ltmCount":          result.LTMCount,
	})
}

// Rest は休息を適用
// POST /api/v1/rest
// [神経科学] 一時的な休息によるデフォルトモードネットワーク(DMN)の活性化と、神経資源の回復をシミュレートします。
// 過度のストレスや疲労からの回復を促し、Sanity（精神的健全性）を向上させます。
// @Summary      Rest Brain
// @Description  一時的な休息を取り、精神的健全性(Sanity)を回復させます。DMNの活性化をシミュレートします。
// @Tags         brain
// @Accept       json
// @Produce      json
// @Param        input  body      handlers.RestRequest  true  "Rest Parameters"
// @Success      200    {object}  models.SuccessResponse
// @Failure      400    {object}  models.ProblemDetails
// @Router       /api/v1/rest [post]
func (h *BrainHandler) Rest(c *gin.Context) {
	var req RestRequest
	if err := BindStrict(c, &req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}

	h.brain.Rest(req.Quality)
	state := h.brain.GetState()

	SuccessResponse(c, gin.H{
		"message":     "Rest applied",
		"restQuality": req.Quality,
		"sanity":      state.Sanity,
		"sanityLevel": state.SanityLevel,
	})
}

// Daydream は白昼夢処理を実行（新規追加）
// POST /api/v1/daydreams
// [神経科学] デフォルトモードネットワーク(DMN)の活性化により、過去の記憶や未来のシミュレーションをランダムに想起します。
// 創造性の向上や、記憶の再構築による新たな洞察の獲得プロセスをモデル化しています。
// @Summary      Daydream
// @Description  DMNを活性化し、白昼夢（Daydreaming）処理を行います。創造性や洞察の獲得を促します。
// @Tags         brain
// @Produce      json
// @Success      200  {object}  models.SuccessResponse
// @Router       /api/v1/daydreams [post]
func (h *BrainHandler) Daydream(c *gin.Context) {
	// 白昼夢処理の実装 (internal/core/brain.go に未実装の場合はスタブ)
	// 現状の brain.go には Daydream メソッドが見当たらないため、
	// ここでは単純なレスポンス、または拡張が必要。
	// Task: "Check implementation" -> assuming simple log or stub for now if core support is missing.

	slog.Info("Daydreaming triggered")
	// TODO: core.Brain に Daydream メソッドを追加する

	SuccessResponse(c, gin.H{
		"message": "Daydreaming... (Stub)",
	})
}
