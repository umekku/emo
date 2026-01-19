package store

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// DB はデータベース接続を管理
type DB struct {
	*sql.DB
}

// NewDB は新しいデータベース接続を作成・初期化
func NewDB(dbPath string) (*DB, error) {
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &DB{db}
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

// initSchema はテーブルスキーマを初期化
func (d *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS memories (
		uuid TEXT PRIMARY KEY,
		text TEXT NOT NULL,
		emotions TEXT NOT NULL, -- JSON string
		weight REAL NOT NULL,
		type TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		last_access DATETIME NOT NULL,
		tags TEXT NOT NULL -- JSON string
	);
	
	CREATE INDEX IF NOT EXISTS idx_memories_type ON memories(type);
	CREATE INDEX IF NOT EXISTS idx_memories_created_at ON memories(created_at);
	CREATE INDEX IF NOT EXISTS idx_memories_weight ON memories(weight);
	`

	_, err := d.Exec(schema)
	return err
}

// Close はデータベース接続を閉じる
func (d *DB) Close() error {
	return d.DB.Close()
}
