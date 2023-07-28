package rules

import (
	"weavatar/pkg/recaptcha"

	"github.com/goravel/framework/contracts/validation"
)

type ReCaptcha struct {
}

// Signature The name of the rule.
func (receiver *ReCaptcha) Signature() string {
	return "recaptcha"
}

// Passes Determine if the validation rule passes.
func (receiver *ReCaptcha) Passes(data validation.Data, val any, options ...any) bool {
	ip, exist := data.Get("ip")
	if !exist {
		return false
	}

	action, ok := options[0].(string)
	if !ok {
		return false
	}
	remoteIp, ok := ip.(string)
	if !ok {
		return false
	}
	response, ok := val.(string)
	if !ok {
		return false
	}

	if !recaptcha.NewRecaptcha().Confirm(remoteIp, response, action) {
		return false
	}

	return true
}

// Message Get the validation error message.
func (receiver *ReCaptcha) Message() string {
	return ""
}
