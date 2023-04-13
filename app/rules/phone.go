package rules

import (
	"regexp"

	"github.com/goravel/framework/contracts/validation"
)

type Phone struct {
}

// Signature The name of the rule.
func (receiver *Phone) Signature() string {
	return "phone"
}

// Passes Determine if the validation rule passes.
func (receiver *Phone) Passes(data validation.Data, val any, options ...any) bool {
	// 正则匹配手机号
	return regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(val.(string))
}

// Message Get the validation error message.
func (receiver *Phone) Message() string {
	return ""
}
