package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config はアプリケーション設定を保持する構造体
// 【役割】環境変数から設定を読み込み、アプリケーション全体で使用可能にする
type Config struct {
	// サーバー設定
	ServiceName string
	Port        string
	Mode        string // debug, release, test

	// データベース設定
	DBPath string

	// セキュリティ設定
	CORSAllowedOrigins []string
	APIKey             string

	// デバッグ設定
	DebugMode bool
	LogLevel  string

	// 脳パラメータ設定
	STMMaxSize             int
	LTMMaxSize             int
	ConsolidationThreshold float64

	// ホルモン設定
	HormoneDecayRate float64

	// 概日リズム設定
	DayTimeStart   int
	NightTimeStart int
}

// LoadConfig は環境変数から設定を読み込む
// デフォルト値もここで管理
func LoadConfig() *Config {
	cfg := &Config{
		// サーバー設定
		ServiceName: getEnv("SERVICE_NAME", "Mind OS"),
		Port:        getEnv("PORT", "8081"),
		Mode:        getEnv("GIN_MODE", "release"),

		// データベース設定
		DBPath: getEnv("DB_PATH", "mind.db"),

		// セキュリティ設定
		CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		APIKey:             getEnv("API_KEY", ""),

		// デバッグ設定
		DebugMode: getEnvAsBool("DEBUG_MODE", false),
		LogLevel:  getEnv("LOG_LEVEL", "info"),

		// 脳パラメータ設定
		STMMaxSize:             getEnvAsInt("STM_MAX_SIZE", 100),
		LTMMaxSize:             getEnvAsInt("LTM_MAX_SIZE", 1000),
		ConsolidationThreshold: getEnvAsFloat("CONSOLIDATION_THRESHOLD", 0.6),

		// ホルモン設定
		HormoneDecayRate: getEnvAsFloat("HORMONE_DECAY_RATE", 10.0),

		// 概日リズム設定
		DayTimeStart:   getEnvAsInt("DAY_TIME_START", 6),
		NightTimeStart: getEnvAsInt("NIGHT_TIME_START", 22),
	}

	// 必須項目の検証
	if err := cfg.validate(); err != nil {
		// 起動阻止 (Fail Fast)
		log.SetFlags(0)
		log.Fatalf("[CONFIG ERROR] %v", err)
	}

	return cfg
}

// validate は設定の必須項目や整合性をチェックし、問題があればエラーを返す
// 【セキュリティ対策】起動時に設定不備を検出し、不正な状態で稼働することを防ぐ
func (c *Config) validate() error {
	var errs []string

	// 1. APIキーの強制（プロダクションモードの場合）
	// 空文字チェック
	if c.Mode == "release" && strings.TrimSpace(c.APIKey) == "" {
		errs = append(errs, "API_KEY is required in release mode")
	}

	// 2. ポート番号の検証
	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 1 || port > 65535 {
		errs = append(errs, fmt.Sprintf("Invalid PORT: %s", c.Port))
	}

	// 3. DBパスの検証
	if strings.TrimSpace(c.DBPath) == "" {
		errs = append(errs, "DB_PATH is required")
	}
	// 拡張子チェック (SQLiteを想定)
	ext := filepath.Ext(c.DBPath)
	if ext != ".db" && ext != ".sqlite" && ext != ".sqlite3" {
		errs = append(errs, fmt.Sprintf("Invalid DB_PATH extension: %s (expected .db, .sqlite, .sqlite3)", c.DBPath))
	}

	if len(errs) > 0 {
		return fmt.Errorf("configuration validation failed:\n - %s", strings.Join(errs, "\n - "))
	}
	return nil
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返すヘルパー
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvAsInt は環境変数をintとして取得
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

// getEnvAsFloat は環境変数をfloat64として取得
func getEnvAsFloat(key string, fallback float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return fallback
}

// getEnvAsBool は環境変数をboolとして取得
func getEnvAsBool(key string, fallback bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return fallback
}

// getEnvAsSlice は環境変数をカンマ区切りでスライスとして取得
func getEnvAsSlice(key string, fallback []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback
	}
	values := strings.Split(valueStr, ",")
	// トリム処理
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}
	return values
}
