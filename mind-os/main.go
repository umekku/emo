package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/umekku/mind-os/docs" // Swagger docs
	"github.com/umekku/mind-os/internal/config"
	"github.com/umekku/mind-os/internal/handlers"
	"github.com/umekku/mind-os/internal/middleware"
)

// @title           Mind-OS API
// @version         1.0
// @description     Emotional Intelligence & Brain OS Microservice.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /
// @schemes         http

func main() {
	// 構造化ロガーの初期化 (JSON形式)
	// 本番用ではログレベルを環境変数から制御する
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// .envファイルを読み込み（存在しない場合はスキップ）
	if err := godotenv.Load(); err != nil {
		slog.Info(".env file not found, using environment variables or defaults")
	}

	// 設定読み込み (エラーがある場合は内部でos.Exit(1))
	cfg := config.LoadConfig()

	// Ginモード設定
	gin.SetMode(cfg.Mode)

	// Ginルーターの初期化 (DefaultではなくNewを使用してカスタムミドルウェアを適用)
	r := gin.New()

	// ミドルウェア適用
	isDev := cfg.Mode != "release"
	r.Use(gin.Recovery())                       // パニック回復
	r.Use(middleware.LoggerMiddleware())        // カスタムロガー
	r.Use(middleware.SecurityMiddleware(isDev)) // セキュリティヘッダー (HSTS, etc.)
	r.Use(middleware.CORSMiddleware())          // CORS設定

	// Swagger エンドポイント (Dev/Debugモードのみが望ましいが、要件に従い常時有効化)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ヘルスチェックエンドポイント
	// @Summary      Health Check
	// @Description  Get service health status
	// @Tags         system
	// @Produce      json
	// @Success      200  {object}  map[string]string
	// @Router       /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": cfg.ServiceName,
			"version": "0.1.0",
		})
	})

	// Hello Worldエンドポイント
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Welcome to %s - Emotional AI Microservice", cfg.ServiceName),
			"status":  "running",
		})
	})

	// 感情分析API
	emotionHandler := handlers.NewEmotionHandler()
	memoryHandler := handlers.NewMemoryHandler()
	motivationHandler := handlers.NewMotivationHandler()
	brainHandler := handlers.NewBrainHandler(cfg)

	api := r.Group("/api")
	{
		emotion := api.Group("/emotion")
		{
			emotion.POST("/assess", emotionHandler.Assess)
		}

		memory := api.Group("/memory")
		{
			memory.POST("/add", memoryHandler.AddMemory)
			memory.GET("/recent", memoryHandler.GetRecentMemories)
			memory.GET("/stats", memoryHandler.GetMemoryStats)
			memory.POST("/sleep", memoryHandler.Sleep)
		}

		motivation := api.Group("/motivation")
		{
			motivation.GET("", motivationHandler.GetMotivation)
			motivation.POST("/feedback", motivationHandler.UpdateMotivation)
			motivation.POST("/emotion-reward", motivationHandler.RewardFromEmotion)
			motivation.POST("/decay", motivationHandler.ApplyDecay)
			motivation.POST("/reset", motivationHandler.Reset)
		}

		// 統合Brain API (v1)
		v1 := api.Group("/v1")
		{
			// リソースベースのエンドポイント定義
			v1.POST("/sensory-inputs", brainHandler.ProcessSensory)
			v1.POST("/sleep-cycles", brainHandler.Sleep)
			v1.GET("/brain-states/current", brainHandler.GetState)
			v1.POST("/daydreams", brainHandler.Daydream)

			// 既存パスのエイリアス/維持(または移行期間)
			// v1.POST("/sensory", brainHandler.ProcessSensory) // Deprecated
			// v1.POST("/sleep", brainHandler.Sleep) // Deprecated

			// その他のエンドポイント (Grouping pending instructions, keeping as is or grouping logically)
			v1.POST("/feedback", brainHandler.Feedback)
			v1.POST("/stress", brainHandler.ApplyStress)
			v1.POST("/rest", brainHandler.Rest) // sleep-cycles/rest?
			v1.GET("/memories", brainHandler.GetRecentMemories)
		}
	}

	// サーバー起動
	addr := ":" + cfg.Port
	slog.Info("Starting server", "service", cfg.ServiceName, "port", cfg.Port, "mode", cfg.Mode)

	// セキュアなサーバー設定
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
