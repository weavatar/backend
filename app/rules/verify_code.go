package rules

import (
	"strconv"

	"github.com/goravel/framework/contracts/validation"

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
	fieldName := options[0].(string)

	// 第二个参数，验证码类型，如 register
	useFor := options[1].(string)

	// 第三个参数（如果有），是否清除验证码，如 false
	clear := true
	var err error
	if len(options) > 2 {
		clear, err = strconv.ParseBool(options[2].(string))
		if err != nil {
			clear = true
		}
	}

	// 取字段值
	field, exist := data.Get(fieldName)
	if !exist {
		return false
	}
	if !verifycode.NewVerifyCode().Check(field.(string), val.(string), useFor, clear) {
		return false
	}

	return true
}

// Message Get the validation error message.
func (receiver *VerifyCode) Message() string {
	return ""
}
