package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetState は脳の状態を取得
// GET /api/v1/brain-states/current
// [神経科学] このエンドポイントは、前頭前野(PFC)が監視する現在の脳の全体状態をスナップショットとして提供します。
// 意欲(線条体)、理性(PFC)、記憶負荷(海馬)の統合的なステータスを示し、ホメオスタシスの維持状況を確認できます。
// @Summary      Get Current Brain State
// @Description  現在の脳の状態（意欲、理性、記憶負荷など）を取得します。ETagによるキャッシュ制御をサポートしています。
// @Tags         brain
// @Produce      json
// @Success      200  {object}  models.SuccessResponse
// @Success      304  {string}  string "Not Modified"
// @Header       200,304  {string}  ETag  "State Hash"
// @Router       /api/v1/brain-states/current [get]
func (h *BrainHandler) GetState(c *gin.Context) {
	state := h.brain.GetState()

	// レスポンスデータの構築
	respData := gin.H{
		"motivation":      state.Motivation,
		"motivationLevel": state.MotivationLevel,
		"sanity":          state.Sanity,
		"sanityLevel":     state.SanityLevel,
		"stmCount":        state.STMCount,
		"ltmCount":        state.LTMCount,
	}

	// ETag生成
	jsonData, _ := json.Marshal(respData)
	hash := sha256.Sum256(jsonData)
	etag := hex.EncodeToString(hash[:])

	// If-None-Match チェック
	if c.GetHeader("If-None-Match") == etag {
		c.Status(http.StatusNotModified)
		return
	}

	c.Header("ETag", etag)
	SuccessResponse(c, respData)
}

// ApplyStress はストレスを適用
// POST /api/v1/stress
// [神経科学] 外部からのストレッサーにより、視床下部-下垂体-副腎系(HPA軸)を介してコルチゾールレベルが上昇する反応をシミュレートします。
// 慢性的なストレスは理性を低下させ、海馬の萎縮（記憶形成の阻害）を引き起こす可能性があります。
// @Summary      Apply Stress
// @Description  外部からのストレスを与え、理性（Sanity）を低下させます。HPA軸の反応をシミュレートします。
// @Tags         brain
// @Accept       json
// @Produce      json
// @Param        input  body      handlers.StressRequest  true  "Stress Level"
// @Success      200    {object}  models.SuccessResponse
// @Failure      400    {object}  models.ProblemDetails
// @Router       /api/v1/stress [post]
func (h *BrainHandler) ApplyStress(c *gin.Context) {
	var req StressRequest
	if err := BindStrict(c, &req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}

	h.brain.ApplyStress(req.Level)
	state := h.brain.GetState()

	SuccessResponse(c, gin.H{
		"message":     "Stress applied",
		"stressLevel": req.Level,
		"sanity":      state.Sanity,
		"sanityLevel": state.SanityLevel,
	})
}
