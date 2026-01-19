package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CustomValidator はバリデーションロジックをカプセル化した構造体
type CustomValidator struct {
	v *validator.Validate
}

// GlobalValidator はアプリケーション全体で使用するバリデータインスタンス
var GlobalValidator *CustomValidator

func init() {
	GlobalValidator = NewCustomValidator()
}

// NewCustomValidator は新しいバリデーターを作成
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	// JSONタグ名をエラーメッセージに使用するための設定
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &CustomValidator{v: v}
}

// ValidateStruct は構造体を検証
func (cv *CustomValidator) ValidateStruct(obj interface{}) error {
	return cv.v.Struct(obj)
}

// BindStrict はリクエストボディをJSONとしてパースし、
// 未知のフィールドの拒否とバリデーションを行う
func BindStrict(c *gin.Context, obj interface{}) error {
	// 1. リクエストボディの読み込み
	// Bodyは一度読むと消えるため、必要に応じてバッファリングするが、
	// ここでは Decode して終わるので直接読み込む

	// JSONデコーダーの設定
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields() // 未知のフィールドを許可しない

	// デコード実行
	if err := decoder.Decode(obj); err != nil {
		if err == io.EOF {
			return fmt.Errorf("request body is empty")
		}
		return fmt.Errorf("json parse error: %w", err)
	}

	// 2. バリデーション実行
	if err := GlobalValidator.ValidateStruct(obj); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
