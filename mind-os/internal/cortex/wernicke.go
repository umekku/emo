package cortex

import (
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// WernickeArea はウェルニッケ野 - 言語理解を担当
// 『脳科学的意味』側頭葉に位置し、受容性言語理解に関与する領域
// 「分節化された形態素解析により入力テキストを分解し、意味のある単語（概念）と発話意図を抽出」
type WernickeArea struct {
	tokenizer *tokenizer.Tokenizer
}

// NewWernickeArea は新しいウェルニッケ野インスタンスを作成
// 「分節化されたKagomeトークナイザー（IPA辞書）を初期化」
func NewWernickeArea() (*WernickeArea, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}

	return &WernickeArea{
		tokenizer: t,
	}, nil
}

// Comprehend はテキストを理解し、概念と意図を抽出
// 【アルゴリズム】
// 1. 形態素解析を実行
// 2. 名詞を抽出し「概念(concepts)」とする
// 3. 疑問符や特定キーワードから「意図(intent)」を分類（挨拶、質問、陳述）
func (w *WernickeArea) Comprehend(text string) (concepts []string, intent string) {
	// 形態素解析
	tokens := w.tokenizer.Tokenize(text)

	concepts = make([]string, 0)
	hasQuestion := false
	hasGreeting := false

	for _, token := range tokens {
		features := token.Features()
		if len(features) < 2 {
			continue
		}

		pos := features[0]       // 品詞
		surface := token.Surface // 表層形

		// 名詞を概念として抽出
		if pos == "名詞" && len(surface) > 1 {
			concepts = append(concepts, surface)
		}

		// 疑問文の検知
		if surface == "？" || surface == "?" || strings.Contains(text, "何") ||
			strings.Contains(text, "どう") || strings.Contains(text, "いつ") {
			hasQuestion = true
		}

		// 挨拶の検知
		if strings.Contains(surface, "こんにちは") || strings.Contains(surface, "おはよう") ||
			strings.Contains(surface, "こんばんは") || strings.Contains(surface, "ありがとう") {
			hasGreeting = true
		}
	}

	// 意図の判定
	if hasQuestion {
		intent = "question"
	} else if hasGreeting {
		intent = "greeting"
	} else if len(concepts) > 0 {
		intent = "statement"
	} else {
		intent = "unknown"
	}

	return concepts, intent
}
