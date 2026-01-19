package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSConfig はCORS設定を保持する構造体
type CORSConfig struct {
	AllowOrigins []string
}

// CORSMiddleware はCORSヘッダーを設定するミドルウェア
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 許可するオリジン（開発環境用）
		// 本番環境では設定ファイル等から読み込むことを推奨
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// OPTIONSメソッド（プリフライトリクエスト）の処理
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
