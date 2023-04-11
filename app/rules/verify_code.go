package rules

import (
	"github.com/goravel/framework/contracts/validation"
	"github.com/spf13/cast"

	"weavatar/packages/verifycode"
)

type VerifyCode struct {
}

// Signature The name of the rule.
func (receiver *VerifyCode) Signature() string {
	return "verify_code"
}

// Passes Determine if the validation rule passes.
func (receiver *VerifyCode) Passes(data validation.Data, val any, options ...any) bool {
	// 第一个参数，字段名称，如 phone
	fieldName := cast.ToString(options[0])

	// 第二个参数，验证码类型，如 register
	forName := cast.ToString(options[1])

	// 取字段值
	field, exist := data.Get(fieldName)
	if !exist {
		return false
	}
	if !verifycode.NewVerifyCode().Check(cast.ToString(field), cast.ToString(val), forName, true) {
		return false
	}

	return true
}

// Message Get the validation error message.
func (receiver *VerifyCode) Message() string {
	return ""
}
