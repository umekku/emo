package store

import (
	"database/sql"
	"encoding/json"

	"github.com/umekku/mind-os/internal/models"
)

// SaveMemory は記憶を保存または更新
func (d *DB) SaveMemory(m models.RuneMemory) error {
	emotionsJSON, err := json.Marshal(m.Emotions)
	if err != nil {
		return err
	}

	tagsJSON, err := json.Marshal(m.Tags)
	if err != nil {
		return err
	}

	query := `
	INSERT OR REPLACE INTO memories (uuid, text, emotions, weight, type, created_at, last_access, tags)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = d.Exec(query,
		m.UUID,
		m.Text,
		string(emotionsJSON),
		m.Weight,
		string(m.Type),
		m.CreatedAt,
		m.LastAccess,
		string(tagsJSON),
	)

	return err
}

// GetRecentMemories は直近の記憶を取得
func (d *DB) GetRecentMemories(limit int) ([]models.RuneMemory, error) {
	query := `
	SELECT uuid, text, emotions, weight, type, created_at, last_access, tags
	FROM memories
	ORDER BY last_access DESC
	LIMIT ?
	`

	rows, err := d.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []models.RuneMemory
	for rows.Next() {
		var m models.RuneMemory
		var emotionsJSON, tagsJSON string
		var typeStr string

		err := rows.Scan(
			&m.UUID,
			&m.Text,
			&emotionsJSON,
			&m.Weight,
			&typeStr,
			&m.CreatedAt,
			&m.LastAccess,
			&tagsJSON,
		)
		if err != nil {
			return nil, err
		}

		m.Type = models.MemoryType(typeStr)

		if err := json.Unmarshal([]byte(emotionsJSON), &m.Emotions); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(tagsJSON), &m.Tags); err != nil {
			return nil, err
		}

		memories = append(memories, m)
	}

	return memories, nil
}

// GetLTMCount は長期記憶の数を取得
func (d *DB) GetLTMCount() (int, error) {
	var count int
	err := d.QueryRow("SELECT COUNT(*) FROM memories WHERE type = ?", models.MemoryLTM).Scan(&count)
	return count, err
}

// DeleteOldMemories は古い記憶を削除してLTMのサイズを制限
// weightが低く、アクセスが古いものを削除
func (d *DB) DeleteOldMemories(keepCount int) error {
	// LTMのサイズ制限として実装

	// 現在の総数を取得
	total, err := d.GetLTMCount()
	if err != nil {
		return err
	}

	if total <= keepCount {
		return nil
	}

	deleteCount := total - keepCount

	// 削除対象: weight ASC (低い順), last_access ASC (古い順)
	// 重要でなく、最近使われていないものから削除
	query := `
	DELETE FROM memories 
	WHERE uuid IN (
		SELECT uuid FROM memories 
		WHERE type = ? 
		ORDER BY weight ASC, last_access ASC 
		LIMIT ?
	)
	`

	_, err = d.Exec(query, models.MemoryLTM, deleteCount)
	return err
}

// GetMemoryByUUID はUUIDで記憶を検索
func (d *DB) GetMemoryByUUID(uuid string) (*models.RuneMemory, error) {
	query := `
	SELECT uuid, text, emotions, weight, type, created_at, last_access, tags
	FROM memories
	WHERE uuid = ?
	`

	var m models.RuneMemory
	var emotionsJSON, tagsJSON string
	var typeStr string

	err := d.QueryRow(query, uuid).Scan(
		&m.UUID,
		&m.Text,
		&emotionsJSON,
		&m.Weight,
		&typeStr,
		&m.CreatedAt,
		&m.LastAccess,
		&tagsJSON,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	m.Type = models.MemoryType(typeStr)
	json.Unmarshal([]byte(emotionsJSON), &m.Emotions)
	json.Unmarshal([]byte(tagsJSON), &m.Tags)

	return &m, nil
}
