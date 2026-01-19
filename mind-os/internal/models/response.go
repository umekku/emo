package models

// Standard Response Structures (Enterprise Grade)

// SuccessResponse は標準的な成功レスポンス
type SuccessResponse struct {
	MindState *MindStateResponse `json:"mindState,omitempty"` // マインドステート
	Reply     string             `json:"reply,omitempty"`     // 応答テキスト (Chat mode only)
	Debug     *DebugInfo         `json:"debug,omitempty"`     // デバッグ情報
	Data      interface{}        `json:"data,omitempty"`      // その他の汎用データ
}

// DebugInfo はデバッグ情報
type DebugInfo struct {
	Cortisol        float64 `json:"cortisol"`
	Oxytocin        float64 `json:"oxytocin"`
	PredictedReward float64 `json:"predictedReward"`
	DaydreamLog     string  `json:"daydreamLog,omitempty"`
}

// ProblemDetails は RFC 9457 準拠のエラーレスポンス
type ProblemDetails struct {
	Type     string `json:"type"`               // エラータイプURI (例: "about:blank")
	Title    string `json:"title"`              // 人間が読める短いサマリー
	Status   int    `json:"status"`             // HTTPステータスコード
	Detail   string `json:"detail,omitempty"`   // エラーの詳細説明
	Instance string `json:"instance,omitempty"` // エラー発生リソースのURI
}

// HandlerResponse はハンドラーが返す汎用レスポンス
// これをラップして統一レスポンスを生成する
type HandlerResponse struct {
	Result interface{} `json:"result"`
}
