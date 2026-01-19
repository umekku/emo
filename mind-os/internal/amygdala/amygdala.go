// amygdala.go: 形態素解析を用いてテキストから感情価（Valence/Arousal）を抽出し、生存本能的な反応生成を行う扁桃体モジュール
package amygdala

import (
	"log/slog"
	"sort"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/umekku/mind-os/internal/models"
)

// Amygdala は扁桃体モジュール - 入力テキストから反射的な感情を生成
type Amygdala struct {
	tokenizer *tokenizer.Tokenizer
	dict      map[string]models.EmotionValue
}

// New は新しい Amygdala インスタンスを作成
func New() (*Amygdala, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}

	return &Amygdala{
		tokenizer: t,
		dict:      initializeDictionary(),
	}, nil
}

// Assess は入力テキストから反射的な感情を評価
// トークン単位で辞書マッチングを行い、感情値を累積させる
func (a *Amygdala) Assess(text string) []models.EmotionValue {
	// Kagomeでトークン化
	tokens := a.tokenizer.Tokenize(text)
	emotionMap := make(models.EmotionMap)
	hit := false

	slog.Debug("Amygdala Assessment", "text", text, "tokens_count", len(tokens))

	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			continue
		}

		// 表層形 (Surface) と 基本形 (BaseForm) を取得
		surface := token.Surface
		baseForm := extractBaseForm(token)

		// 1. 表層形で検索
		if val, ok := a.dict[surface]; ok {
			updateEmotionMap(emotionMap, val)
			hit = true
			slog.Debug("Emotion Hit (Surface)", "token", surface, "value", val)
			continue
		}

		// 2. 基本形で検索
		if val, ok := a.dict[baseForm]; ok {
			updateEmotionMap(emotionMap, val)
			hit = true
			slog.Debug("Emotion Hit (Base)", "token", baseForm, "value", val)
		}
	}

	// 何もヒットしない場合は Neutral
	if !hit {
		return []models.EmotionValue{
			{Code: models.EmotionNeutral, Value: 10},
		}
	}

	// ソートして返す
	result := emotionMap.ToEmotionValues()
	sort.Slice(result, func(i, j int) bool {
		return result[i].Value > result[j].Value
	})

	return result
}

// extractBaseForm はトークンから基本形を抽出するヘルパー
func extractBaseForm(token tokenizer.Token) string {
	features := token.Features()
	if len(features) > 6 && features[6] != "*" {
		return features[6]
	}
	return token.Surface
}

// updateEmotionMap は感情マップを更新（加算）する
func updateEmotionMap(em models.EmotionMap, ev models.EmotionValue) {
	current := em[ev.Code]
	newValue := current + ev.Value
	if newValue > 100 {
		newValue = 100
	}
	em[ev.Code] = newValue
}

// initializeDictionary は感情辞書を初期化
func initializeDictionary() map[string]models.EmotionValue {
	return map[string]models.EmotionValue{
		// Joy (喜び)
		"最高":   {Code: models.EmotionJoy, Value: 90},
		"楽しい":  {Code: models.EmotionJoy, Value: 80},
		"嬉しい":  {Code: models.EmotionJoy, Value: 85},
		"笑":    {Code: models.EmotionJoy, Value: 60},
		"良":    {Code: models.EmotionJoy, Value: 50},
		"好き":   {Code: models.EmotionJoy, Value: 70}, // Joy要素としての好き
		"良い":   {Code: models.EmotionJoy, Value: 50},
		"やった":  {Code: models.EmotionJoy, Value: 80},
		"美味しい": {Code: models.EmotionJoy, Value: 85}, // 食事関連
		"旨い":   {Code: models.EmotionJoy, Value: 80},

		// Love (信頼/愛)
		"愛":     {Code: models.EmotionLove, Value: 90},
		"信頼":    {Code: models.EmotionLove, Value: 80},
		"相棒":    {Code: models.EmotionLove, Value: 85},
		"一緒":    {Code: models.EmotionLove, Value: 60},
		"味方":    {Code: models.EmotionLove, Value: 70},
		"なでなで":  {Code: models.EmotionLove, Value: 65},
		"ありがとう": {Code: models.EmotionLove, Value: 60},

		// Anger (怒り)
		"バカ":    {Code: models.EmotionAnger, Value: 80},
		"うざい":   {Code: models.EmotionAnger, Value: 70},
		"嫌い":    {Code: models.EmotionAnger, Value: 80},
		"クソ":    {Code: models.EmotionAnger, Value: 85},
		"怒":     {Code: models.EmotionAnger, Value: 90},
		"ふざけるな": {Code: models.EmotionAnger, Value: 75},

		// Sadness (悲しみ)
		"悲しい": {Code: models.EmotionSadness, Value: 80},
		"辛い":  {Code: models.EmotionSadness, Value: 85},
		"泣":   {Code: models.EmotionSadness, Value: 70},
		"だめ":  {Code: models.EmotionSadness, Value: 60},
		"無理":  {Code: models.EmotionSadness, Value: 65},
		"最悪":  {Code: models.EmotionSadness, Value: 90},
		"ごめん": {Code: models.EmotionSadness, Value: 50}, // 罪悪感としての悲しみ

		// Surprise (驚き)
		"えっ":   {Code: models.EmotionSurprise, Value: 60},
		"すごい":  {Code: models.EmotionSurprise, Value: 70},
		"まさか":  {Code: models.EmotionSurprise, Value: 80},
		"！？":   {Code: models.EmotionSurprise, Value: 75},
		"びっくり": {Code: models.EmotionSurprise, Value: 80},

		// Fear (恐れ)
		"怖い":  {Code: models.EmotionFear, Value: 85},
		"やばい": {Code: models.EmotionFear, Value: 70},
		"逃げ":  {Code: models.EmotionFear, Value: 80},
		"不安":  {Code: models.EmotionFear, Value: 60},
		"警告":  {Code: models.EmotionFear, Value: 75},
		"エラー": {Code: models.EmotionFear, Value: 65},
		"バグ":  {Code: models.EmotionFear, Value: 80},

		// Disgust (嫌悪)
		"苦い":  {Code: models.EmotionDisgust, Value: 70},
		"不味い": {Code: models.EmotionDisgust, Value: 80},
		"臭い":  {Code: models.EmotionDisgust, Value: 85},
		"キモい": {Code: models.EmotionDisgust, Value: 90},
		"汚い":  {Code: models.EmotionDisgust, Value: 85},
	}
}
