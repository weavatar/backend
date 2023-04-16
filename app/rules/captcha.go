package rules

import (
	"github.com/goravel/framework/contracts/validation"

	"weavatar/packages/captcha"
)

type Captcha struct {
}

// Signature The name of the rule.
func (receiver *Captcha) Signature() string {
	return "captcha"
}

// Passes Determine if the validation rule passes.
func (receiver *Captcha) Passes(data validation.Data, val any, options ...any) bool {
	captchaID, exist := data.Get("captcha_id")
	if !exist {
		return false
	}

	// 第一个参数（如果有），是否清除验证码，如 false
	clear := true
	if len(options) > 0 {
		clear = options[0].(bool)
	}

	if !captcha.NewCaptcha().VerifyCaptcha(captchaID.(string), val.(string), clear) {
		return false
	}

	return true
}

// Message Get the validation error message.
func (receiver *Captcha) Message() string {
	return ""
}
