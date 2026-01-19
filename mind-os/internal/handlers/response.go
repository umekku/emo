package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umekku/mind-os/internal/models"
)

// ErrorResponse は RFC 9457 に準拠したエラーレスポンスを返すヘルパー
// status: HTTPステータスコード
// title: エラーのタイトル
// detail: エラーの詳細 (error.Error()など)
func ErrorResponse(c *gin.Context, status int, title string, detail string) {
	// ログ出力 (Error or Warn based on status)
	if status >= 500 {
		slog.Error("API Error", "status", status, "title", title, "detail", detail, "path", c.Request.URL.Path)
	} else {
		slog.Warn("API Client Error", "status", status, "title", title, "detail", detail, "path", c.Request.URL.Path)
	}

	problem := models.ProblemDetails{
		Type:     "about:blank", // 特に定義されていない場合のデフォルト
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}

	c.Header("Content-Type", "application/problem+json")
	c.JSON(status, problem)
}

// SuccessResponse は標準形式の成功レスポンスを返すヘルパー
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
