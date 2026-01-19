package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

// SecurityMiddleware は包括的なセキュリティヘッダーを設定するミドルウェア
func SecurityMiddleware(isDev bool) gin.HandlerFunc {
	secureMiddleware := secure.New(secure.Options{
		// HSTS設定 (Strict-Transport-Security)
		STSSeconds:           31536000, // 1年
		STSIncludeSubdomains: true,
		STSPreload:           true,

		// クリックジャッキング対策 (X-Frame-Options)
		FrameDeny: true,

		// MIMEスニッフィング対策 (X-Content-Type-Options)
		ContentTypeNosniff: true,

		// XSS対策 (X-XSS-Protection)
		BrowserXssFilter: true,

		// コンテンツセキュリティポリシー (Content-Security-Policy)
		// 必要に応じて緩和するが、基本は厳格に
		ContentSecurityPolicy: "default-src 'self'",

		// 開発モードではHTTPSリダイレクトなどを無効化
		IsDevelopment: isDev,
	})

	return func(c *gin.Context) {
		err := secureMiddleware.Process(c.Writer, c.Request)

		// リダイレクトが発生した場合などは処理を中断
		if err != nil {
			c.Abort()
			return
		}

		// Don't write header if we already wrote it (redirects)
		if c.Writer.Written() {
			c.Abort()
			return
		}
	}
}
