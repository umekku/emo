package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DeprecationMiddleware は API の非推奨情報と廃止予定日をヘッダーに付与します。
// RFC 8594 (Sunset) および Link Header (Deprecation) に準拠します。
func DeprecationMiddleware(sunsetDate string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Deprecation ヘッダー (boolean ではなく、情報へのリンクなどが推奨されるが、ここでは推奨としてtrueを設定する場合もある)
		// しかし標準的には "Deprecation: <date>" or "Deprecation: true"
		// ここでは Sunset ヘッダーをメインに使用し、Deprecation は true とする
		c.Header("Deprecation", "true")

		// Sunset ヘッダー: 廃止日時(HTTP Date format required)
		// 入力が "2025-12-31" のような形式であればパースして RFC1123 形式に変換
		if t, err := time.Parse("2006-01-02", sunsetDate); err == nil {
			c.Header("Sunset", t.Format(http.TimeFormat))
		} else {
			// パース失敗時はそのまま入れるか、ログ出すなど
			c.Header("Sunset", sunsetDate)
		}

		// 警告ヘッダー (Warning) も付けるとより親切 (code 299)
		c.Header("Warning", fmt.Sprintf(`299 - "This endpoint is deprecated and will be removed on %s"`, sunsetDate))

		c.Next()
	}
}
