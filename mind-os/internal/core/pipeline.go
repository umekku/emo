package core

import (
	"log/slog"
	"time"

	"github.com/umekku/mind-os/internal/cortex"
	"github.com/umekku/mind-os/internal/models"
)

// ProcessInput は入力を脳全体で処理
// 【神経科学的意味】感覚入力から感情・認知・記憶・言語までの統合処理パイプライン
// 【処理フロー】
// 1. ホルモン減衰・概日リズム更新
// 2. 視床フィルタリング（順応・ゲイン計算）
// 3. 感情生成（扁桃体） / 共感プロセス（ミラーニューロン）
// 4. ホルモン更新（視床下部）
// 5. 意欲更新（大脳基底核）
// 6. 感情調整（前頭前皮質）
// 7. 記憶保存（海馬）
// 8. 言語理解（ウェルニッケ野）
// 9. 言語生成（ブローカ野）
func (b *Brain) ProcessInput(input models.SensoryInput) models.MindStateResponse {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 1. 時間経過処理 (ホルモン減衰)
	b.Hypothalamus.Decay()

	// 1.5. 概日リズム更新（体内時計）
	b.Hypothalamus.UpdateCircadianRhythm(time.Now())

	// 2. 視床フィルタリング (順応・ゲイン計算)
	gain, err := b.Thalamus.Filter(input)
	if err != nil {
		slog.Warn("Thalamus filter error", "error", err)
		gain = 1.0 // エラー時のゲインは1.0（影響なし）にする
	}

	var rawEmotions []models.EmotionValue
	text := input.InputText

	if input.Type == models.SignalPhysical {
		// 物理的刺激処理: 扁桃体分析をスキップし、信号値を直接感情に変換
		rawEmotions = b.processPhysicalSignal(input.SignalValue, gain)
	} else {
		// 会話（デフォルト）: 扁桃体によるテキスト解析
		rawEmotions = b.processChatInput(text, gain)
	}

	// 5. 前頭前皮質: 理性による感情の調整
	// 視床下部のホルモン状態を取得
	cortisol, oxytocin := b.Hypothalamus.GetStatus()
	controlledEmotions := b.PFC.Arbitrate(rawEmotions, cortisol, oxytocin)

	// 6. 海馬: 記憶として保存
	b.Hippocampus.AddEpisode(text, controlledEmotions)

	// 7. 言語理解（ウェルニッケ野）
	concepts, intent := b.Wernicke.Comprehend(text)

	// 8. レスポンスを生成
	response := b.generateMindState(controlledEmotions)

	// 9. 言語生成（ブローカ野）
	// Chat入力の場合のみテキスト応答を生成
	if input.Type == models.SignalChat {
		replyText := b.Broca.GenerateResponse(
			controlledEmotions,
			response.Motivation,
			response.Sanity,
			concepts,
			intent,
		)
		response.ReplyText = replyText
	}

	return response
}

// processPhysicalSignal は物理的刺激を処理
// 【処理内容】信号値を直接感情に変換し、ホルモンと意欲を更新
func (b *Brain) processPhysicalSignal(signalValue int, gain float64) []models.EmotionValue {
	var rawEmotions []models.EmotionValue
	val := float64(signalValue)
	stressor := 0.0
	affection := 0.0

	if val > 0 {
		// 正の値 -> Joy (快感) -> 愛着(Oxytocin)
		// Gainを適用
		addEmotion(&rawEmotions, models.EmotionJoy, int(val*gain))
		affection = val * gain
	} else if val < 0 {
		// 負の値 -> Disgust (不快感) -> ストレス(Cortisol)
		addEmotion(&rawEmotions, models.EmotionDisgust, int(-val*gain))
		stressor = -val * gain
	} else {
		addEmotion(&rawEmotions, models.EmotionNeutral, 10)
	}

	// 視床下部更新
	b.Hypothalamus.Update(stressor, affection)

	// 意欲への直接作用（報酬系への直接入力）
	// RPE用に0-100にマッピング (-100->0, 0->50, 100->100)
	// Gainを適用
	actualReward := ((val * gain) + 100.0) / 2.0
	b.BasalGanglia.UpdateMotivation(actualReward)

	return rawEmotions
}

// processChatInput はチャット入力を処理
// 【処理内容】テキストから感情を生成し、共感プロセスを適用
func (b *Brain) processChatInput(text string, gain float64) []models.EmotionValue {
	// 3. 感情生成 (Amygdala)
	rawEmotions := b.Amygdala.Assess(text)

	// 3.5. 共感プロセス（ミラーニューロンシステム）
	// ユーザーの感情を推定
	userEmotion, err := b.Mirror.SimulateUserEmotion(text)
	if err != nil {
		slog.Warn("Mirror neuron simulation error", "error", err)
	}

	// Oxytocinレベルに基づいて共感レベルを更新
	_, oxytocin := b.Hypothalamus.GetStatus()
	b.Mirror.UpdateEmpathyLevel(oxytocin)

	// 情動伝染: ユーザーの感情をAI自身の感情にブレンド
	empathyLevel := b.Mirror.GetEmpathyLevel()
	rawEmotions = cortex.BlendEmotions(userEmotion, rawEmotions, empathyLevel)

	// 概日リズムの効果を取得
	_, emotionalSensitivity, _ := b.Hypothalamus.GetCircadianEffects()

	// Gainと概日リズムの感情感度を感情値に適用
	for i := range rawEmotions {
		// 夜間は感情的になる（Grief, Sadness, Love への感度上昇など）
		sensitivity := gain
		if rawEmotions[i].Code == models.EmotionGrief ||
			rawEmotions[i].Code == models.EmotionLove ||
			rawEmotions[i].Code == models.EmotionFear {
			sensitivity *= emotionalSensitivity
		}

		rawEmotions[i].Value = int(float64(rawEmotions[i].Value) * sensitivity)
		if rawEmotions[i].Value > 100 {
			rawEmotions[i].Value = 100
		}
	}

	// 感情からストレス/愛着を算出
	stressor := 0.0
	affection := 0.0
	for _, e := range rawEmotions {
		val := float64(e.Value)
		switch e.Code {
		case models.EmotionJoy, models.EmotionLove:
			affection += val
		case models.EmotionAnger, models.EmotionFear, models.EmotionDisgust, models.EmotionGrief:
			stressor += val
		}
	}

	// 視床下部更新 (感情由来)
	// 感情値そのままでは強すぎる可能性があるため係数を掛ける
	b.Hypothalamus.Update(stressor*0.5, affection*0.5)

	// 4. 意欲更新 (感情由来の報酬)
	// 感情の平均値を計算して報酬とする
	total := 0
	for _, emotion := range rawEmotions {
		total += emotion.Value
	}
	if len(rawEmotions) > 0 {
		avgEmotion := float64(total) / float64(len(rawEmotions))
		// Gainはすでに感情値に適用済みなので、ここではそのまま使用
		b.BasalGanglia.UpdateMotivation(avgEmotion)
	}

	return rawEmotions
}
