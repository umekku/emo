package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware はリクエスト情報を構造化ログ（slog）として出力するミドルウェア
func LoggerMiddleware() gin.HandlerFunc {
	// JSONハンドラーを使用（実運用ではログ収集基盤に合わせて調整）
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// リクエスト処理
		c.Next()

		// 処理後の情報収集
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		// ログ属性の構築
		attrs := []slog.Attr{
			slog.Int("status", statusCode),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
		}

		if errorMessage != "" {
			attrs = append(attrs, slog.String("error", errorMessage))
		}

		// ステータスコードに応じたログレベル
		if statusCode >= 500 {
			logger.LogAttrs(c, slog.LevelError, "Request Failed", attrs...)
		} else if statusCode >= 400 {
			logger.LogAttrs(c, slog.LevelWarn, "Request Warn", attrs...)
		} else {
			logger.LogAttrs(c, slog.LevelInfo, "Request Success", attrs...)
		}
	}
}
