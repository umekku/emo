package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// 環境変数を一時的にクリア（または保存）
	os.Unsetenv("SERVICE_NAME")
	os.Unsetenv("PORT")
	os.Unsetenv("GIN_MODE")

	cfg := LoadConfig()

	if cfg.ServiceName != "Mind OS" {
		t.Errorf("Default ServiceName should be 'Mind OS', got %s", cfg.ServiceName)
	}
	if cfg.Port != "8081" {
		t.Errorf("Default Port should be '8081', got %s", cfg.Port)
	}
	if cfg.Mode != "release" {
		t.Errorf("Default Mode should be 'release', got %s", cfg.Mode)
	}
}

func TestLoadConfig_EnvVars(t *testing.T) {
	os.Setenv("SERVICE_NAME", "Test Service")
	os.Setenv("PORT", "9090")
	os.Setenv("GIN_MODE", "debug")
	defer func() {
		os.Unsetenv("SERVICE_NAME")
		os.Unsetenv("PORT")
		os.Unsetenv("GIN_MODE")
	}()

	cfg := LoadConfig()

	if cfg.ServiceName != "Test Service" {
		t.Errorf("ServiceName should be 'Test Service', got %s", cfg.ServiceName)
	}
	if cfg.Port != "9090" {
		t.Errorf("Port should be '9090', got %s", cfg.Port)
	}
	if cfg.Mode != "debug" {
		t.Errorf("Mode should be 'debug', got %s", cfg.Mode)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")

	val := getEnvAsInt("TEST_INT", 0)
	if val != 123 {
		t.Errorf("getEnvAsInt should return 123, got %d", val)
	}

	valDefault := getEnvAsInt("NON_EXISTENT", 999)
	if valDefault != 999 {
		t.Errorf("getEnvAsInt should return default 999, got %d", valDefault)
	}
}
