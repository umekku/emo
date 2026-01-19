package core

import (
	"time"

	"github.com/umekku/mind-os/internal/models"
)

// Daydream はデイドリーム（白昼夢）処理を実行
// 指定された時間分、マインドワンダリングを行い、結果をレスポンスに含める
func (b *Brain) Daydream(duration time.Duration) models.MindStateResponse {
	// マインドワンダリングを実行
	daydreamLog := b.WanderMind(duration)

	// 現在の状態を取得してレスポンスを生成
	// 空の感情で generateMindState を呼び出す
	response := b.generateMindState([]models.EmotionValue{})

	// DaydreamLogを追加
	response.DaydreamLog = daydreamLog

	return response
}
