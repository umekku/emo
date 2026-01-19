package core

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/umekku/mind-os/internal/amygdala"
	"github.com/umekku/mind-os/internal/basal"
	"github.com/umekku/mind-os/internal/config"
	"github.com/umekku/mind-os/internal/cortex"
	"github.com/umekku/mind-os/internal/hippocampus"
	"github.com/umekku/mind-os/internal/hypothalamus"
	"github.com/umekku/mind-os/internal/pfc"
	"github.com/umekku/mind-os/internal/store"
	"github.com/umekku/mind-os/internal/thalamus"
)

// Brain は脳全体を統合する構造体
// 【役割】各脳機能モジュール（扁桃体、海馬、前頭前皮質など）の保持とライフサイクル管理
// 【構造】スレッドセーフな操作のためのMutexと、各サブシステムへのポインタを持つ
type Brain struct {
	mu sync.RWMutex

	// 脳の構成体
	Amygdala     *amygdala.Amygdala        // 扁桃体 - 反射的感情生成
	Hippocampus  *hippocampus.Hippocampus  // 海馬 - 記憶管理（STM/LTM）
	BasalGanglia *basal.BasalGanglia       // 大脳基底核 - 意欲と報酬系
	PFC          *pfc.PrefrontalCortex     // 前頭前皮質 - 理性による感情制御
	Hypothalamus *hypothalamus.Homeostasis // 視床下部 - ホルモンと恒常性維持
	Thalamus     *thalamus.Thalamus        // 視床 - 感覚入力のフィルタリング
	Mirror       *cortex.SocialCognition   // ミラーニューロン - 共感と社会的認知
	Wernicke     *cortex.WernickeArea      // ウェルニッケ野 - 言語理解
	Broca        *cortex.BrocaArea         // ブローカ野 - 言語生成

	// インフラ
	DB *store.DB // データベース接続
}

// New は新しい Brain インスタンスを作成
// 【処理内容】
// 1. データベース接続の確立
// 2. configに基づく各脳機能モジュールの初期化
// 3. 依存関係の注入
func New(cfg *config.Config) *Brain {
	dbPath := cfg.DBPath
	if dbPath == "" {
		dbPath = "mind.db"
	}

	// 絶対パスに変換（ディレクトリ作成のため）
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		slog.Warn("Failed to resolve absolute path", "error", err)
		absPath = dbPath
	}

	db, err := store.NewDB(absPath)
	if err != nil {
		slog.Warn("Failed to initialize DB. Running in memory-only mode.", "error", err)
		db = nil
	} else {
		slog.Info("Connected to mind database", "path", dbPath)
	}

	am, err := amygdala.New()
	if err != nil {
		slog.Error("Failed to initialize Amygdala", "error", err)
		os.Exit(1)
	}

	wernicke, err := cortex.NewWernickeArea()
	if err != nil {
		slog.Error("Failed to initialize Wernicke's area", "error", err)
		os.Exit(1)
	}

	// 起動時の初期化ログ
	slog.Info("Brain initializing modules",
		"STM_MAX", cfg.STMMaxSize,
		"LTM_MAX", cfg.LTMMaxSize,
		"THRESHOLD", cfg.ConsolidationThreshold,
	)

	// 海馬の初期化
	hc := hippocampus.New(db)

	return &Brain{
		Amygdala:     am,
		Hippocampus:  hc,
		BasalGanglia: basal.New(),
		PFC:          pfc.New(),
		Hypothalamus: hypothalamus.NewHomeostasis(), // ホルモン減衰率・内部定数使用
		Thalamus:     thalamus.New(),
		Mirror:       cortex.New(am),
		Wernicke:     wernicke,
		Broca:        cortex.NewBrocaArea(),
		DB:           db,
	}
}

// Close はリソースを解放
// 【処理内容】データベース接続などのクリーンアップを行う
func (b *Brain) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 終了時に長期記憶への保存（睡眠処理）を行うのが望ましい
	// b.Sleep()

	if b.DB != nil {
		if err := b.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

// 内部定数: 1日のサイクル
const (
	DayDuration = 24 * time.Hour
)
