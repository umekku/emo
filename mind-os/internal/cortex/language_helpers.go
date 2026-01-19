package cortex

import (
	"math/rand"

	"github.com/umekku/mind-os/internal/models"
)

// getDominantEmotion は最も強い感情を取得
// 【アルゴリズム】感情リストの中で最大値を持つ感情コードを返す
// 【計算量】O(N) - Nは感情の種類数
func getDominantEmotion(emotions []models.EmotionValue) models.EmotionCode {
	if len(emotions) == 0 {
		return models.EmotionNeutral
	}

	maxEmotion := emotions[0]
	for _, e := range emotions {
		if e.Value > maxEmotion.Value {
			maxEmotion = e
		}
	}

	return maxEmotion.Code
}

// emotionToTemplateKey は感情コードをテンプレートキーに変換
// 【用途】内部的な感情コードモデルをテンプレートのマップキー文字列にマッピング
func emotionToTemplateKey(emotion models.EmotionCode) string {
	switch emotion {
	case models.EmotionJoy:
		return "joy"
	case models.EmotionAnger:
		return "anger"
	case models.EmotionFear:
		return "fear"
	case models.EmotionLove:
		return "love"
	case models.EmotionDisgust:
		return "disgust"
	case models.EmotionGrief:
		return "grief"
	default:
		return "neutral"
	}
}

// addConfusion は理性が低い時の混乱を追加
// 【演出効果】テキスト末尾に曖昧な表現を付加し、混乱状態を表現
func addConfusion(text string) string {
	confusions := []string{"...", "あれ？", "どうだっけ...", "頭が回らない...", "何か変だな..."}
	return text + " " + confusions[rand.Intn(len(confusions))]
}
