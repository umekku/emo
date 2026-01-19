package cortex

import (
	"fmt"
	"math/rand"

	"github.com/umekku/mind-os/internal/models"
)

// BrocaArea はブローカ野 - 言語生成を司る
// 『脳科学的意味』前頭葉に位置し、可動性言語生成に関与する領域
// 「分節化された現在の感情・意欲・理性状態に基づいて適切な応答テキストを選択・生成」
type BrocaArea struct {
	templates map[string][]string // 感情ごとのテンプレート
}

// NewBrocaArea は新しいブローカ野インスタンスを作成
// 「分節化された応答テンプレートをロードして初期化」
func NewBrocaArea() *BrocaArea {
	return &BrocaArea{
		templates: initializeTemplates(),
	}
}

// GenerateResponse は現在の心理状態に基づいて応答を生成
// 【アルゴリズム】
// 1. 意欲チェック: 極端に低い場合は応答拒否
// 2. 感情判定: 支配的な感情を特定
// 3. 意図分類: 挨拶/質問/陳述に応じて生成ロジックを選択
// 4. 理性チェック: 理性が低い場合は混乱表現を追加
func (b *BrocaArea) GenerateResponse(
	emotions []models.EmotionValue,
	motivation float64,
	sanity float64,
	concepts []string,
	intent string,
) string {
	// 意欲が極端に低い場合は短文または無言
	if motivation < 0.2 {
		return "..."
	}

	// 主要な感情を判定
	dominantEmotion := getDominantEmotion(emotions)

	// 意図に応じた基本応答
	var baseResponse string
	switch intent {
	case "greeting":
		baseResponse = b.generateGreeting(dominantEmotion, motivation)
	case "question":
		baseResponse = b.generateQuestionResponse(dominantEmotion, sanity, concepts)
	case "statement":
		baseResponse = b.generateStatementResponse(dominantEmotion, concepts)
	default:
		baseResponse = b.generateDefaultResponse(dominantEmotion, motivation)
	}

	// 理性が低い場合、文脈が乱れる
	if sanity < 0.3 {
		baseResponse = addConfusion(baseResponse)
	}

	return baseResponse
}

// generateGreeting は挨拶応答を生成
func (b *BrocaArea) generateGreeting(emotion models.EmotionCode, motivation float64) string {
	if motivation < 0.3 {
		return "...こんにちは"
	}

	switch emotion {
	case models.EmotionJoy:
		return "こんにちは、元気だね！"
	case models.EmotionAnger:
		return "...何？"
	case models.EmotionGrief:
		return "...こんにちは..."
	default:
		return "こんにちは"
	}
}

// generateQuestionResponse は質問への応答を生成
func (b *BrocaArea) generateQuestionResponse(emotion models.EmotionCode, sanity float64, concepts []string) string {
	if sanity < 0.3 {
		return "よくわからない..."
	}

	conceptText := ""
	if len(concepts) > 0 {
		conceptText = concepts[0]
	}

	switch emotion {
	case models.EmotionJoy:
		if conceptText != "" {
			return fmt.Sprintf("%sのこと？知ってるよ！", conceptText)
		}
		return "何だろう？教えて！"
	case models.EmotionAnger:
		return "今はそんな気分じゃない"
	case models.EmotionFear:
		return "わからない...怖い..."
	default:
		if conceptText != "" {
			return fmt.Sprintf("%sについて？うーん...", conceptText)
		}
		return "何だろう..."
	}
}

// generateStatementResponse は陳述への応答を生成
func (b *BrocaArea) generateStatementResponse(emotion models.EmotionCode, concepts []string) string {
	conceptText := ""
	if len(concepts) > 0 {
		conceptText = concepts[0]
	}

	templates := b.templates[emotionToTemplateKey(emotion)]
	if len(templates) == 0 {
		templates = b.templates["neutral"]
	}

	baseReply := templates[rand.Intn(len(templates))]

	if conceptText != "" && emotion == models.EmotionJoy {
		return fmt.Sprintf("%sって%s", conceptText, baseReply)
	}

	return baseReply
}

// generateDefaultResponse はデフォルト応答を生成
func (b *BrocaArea) generateDefaultResponse(emotion models.EmotionCode, motivation float64) string {
	if motivation < 0.3 {
		return "..."
	}

	templates := b.templates[emotionToTemplateKey(emotion)]
	if len(templates) == 0 {
		return "..."
	}

	return templates[rand.Intn(len(templates))]
}
