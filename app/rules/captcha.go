package rules

import (
	"github.com/goravel/framework/contracts/validation"
	"github.com/spf13/cast"

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

	if !captcha.NewCaptcha().VerifyCaptcha(cast.ToString(captchaID), cast.ToString(val), true) {
		return false
	}

	return true
}

// Message Get the validation error message.
func (receiver *Captcha) Message() string {
	return ""
}
