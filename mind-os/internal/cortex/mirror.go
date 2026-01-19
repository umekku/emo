package cortex

import (
	"sync"

	"github.com/umekku/mind-os/internal/amygdala"
	"github.com/umekku/mind-os/internal/models"
)

// SocialCognition はミラーニューロンシステム - 社会的認知と共感を管理
type SocialCognition struct {
	mu sync.RWMutex

	EmpathyLevel float64 // 共感の強さ (0.0-1.0)

	// 内部参照
	amygdala *amygdala.Amygdala // ユーザー感情推定のため扁桃体を参照
}

// New は新しい SocialCognition インスタンスを作成
func New(amyg *amygdala.Amygdala) *SocialCognition {
	return &SocialCognition{
		EmpathyLevel: 0.5, // 初期値: 中程度の共感性
		amygdala:     amyg,
	}
}

// SimulateUserEmotion はユーザーが抱いている感情を推測
// テキストから「ユーザーの感情状態」を推定し、最も強い感情を返す
func (sc *SocialCognition) SimulateUserEmotion(text string) (models.EmotionValue, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// Amygdalaを使ってテキストを解析
	// （本来は文脈を考慮した推定が必要だが、簡易実装として扁桃体の結果を流用）
	emotions := sc.amygdala.Assess(text)

	if len(emotions) == 0 {
		// 感情が検出されない場合はNeutralを返す
		return models.EmotionValue{
			Code:  models.EmotionNeutral,
			Value: 50,
		}, nil
	}

	// 最も強い感情を「ユーザーの感情」として推定
	maxEmotion := emotions[0]
	for _, e := range emotions {
		if e.Value > maxEmotion.Value {
			maxEmotion = e
		}
	}

	// 文脈の反転処理（簡易）
	// 例: テキストに「痛い」「辛い」などのネガティブワードがある場合、
	// ユーザーはFear/Griefを感じていると推定
	userEmotion := sc.contextualizeEmotion(maxEmotion)

	return userEmotion, nil
}

// contextualizeEmotion は感情を文脈に応じて調整
// AI自身の反応ではなく、「ユーザーが感じている感情」に変換
func (sc *SocialCognition) contextualizeEmotion(emotion models.EmotionValue) models.EmotionValue {
	// 簡易実装: 基本的にはそのまま返すが、
	// 将来的には文脈解析（一人称表現の検知など）を追加可能

	// 例: Angerが検出された場合、ユーザーが怒っているか、
	// ユーザーが怒りを向けられているかを区別する必要がある
	// 現状は単純にそのまま返す

	return emotion
}

// UpdateEmpathyLevel は共感レベルを更新
// Oxytocinレベルに基づいて共感の強さを調整
func (sc *SocialCognition) UpdateEmpathyLevel(oxytocin float64) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Oxytocin 0-100 を 0.0-1.0 の共感レベルに変換
	// Oxytocin が高いほど共感性が高まる
	// 最低でも0.2（基本的な共感）、最大1.0（完全な共感）
	sc.EmpathyLevel = 0.2 + (oxytocin / 100.0 * 0.8)

	if sc.EmpathyLevel > 1.0 {
		sc.EmpathyLevel = 1.0
	}
	if sc.EmpathyLevel < 0.0 {
		sc.EmpathyLevel = 0.0
	}
}

// GetEmpathyLevel は現在の共感レベルを返す
func (sc *SocialCognition) GetEmpathyLevel() float64 {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.EmpathyLevel
}

// BlendEmotions はユーザー感情とAI自身の感情をブレンド（情動伝染）
// userEmotion: 推定されたユーザーの感情
// myEmotions: AI自身が生成した感情
// empathyStrength: 共感の強さ（通常はEmpathyLevel）
func BlendEmotions(userEmotion models.EmotionValue, myEmotions []models.EmotionValue, empathyStrength float64) []models.EmotionValue {
	// 情動伝染: ユーザーの感情をAI自身の感情に加算
	// empathyStrength が高いほど、ユーザーの感情の影響が大きい

	// ユーザー感情の影響度を計算
	contagionValue := int(float64(userEmotion.Value) * empathyStrength * 0.5)

	if contagionValue <= 0 {
		return myEmotions
	}

	// 既存の感情リストに同じ感情コードがあれば加算、なければ追加
	blended := make([]models.EmotionValue, len(myEmotions))
	copy(blended, myEmotions)

	found := false
	for i, e := range blended {
		if e.Code == userEmotion.Code {
			// 同じ感情があれば加算
			blended[i].Value += contagionValue
			if blended[i].Value > 100 {
				blended[i].Value = 100
			}
			found = true
			break
		}
	}

	if !found {
		// 新しい感情として追加
		blended = append(blended, models.EmotionValue{
			Code:  userEmotion.Code,
			Value: contagionValue,
		})
	}

	return blended
}
