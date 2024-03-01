package avatar

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type StoreAvatarRequest struct {
	Raw        string `json:"raw" form:"raw"`
	VerifyCode string `json:"verify_code" form:"verify_code"`

	Captcha string `json:"captcha" form:"captcha"`
}

func (r *StoreAvatarRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *StoreAvatarRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"raw":         "required|string",
		"verify_code": "required|len:6|number|verify_code:raw,avatar",
		"captcha":     "required|recaptcha:avatar",
	}
}

func (r *StoreAvatarRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"raw.required":            "地址不能为空",
		"raw.string":              "地址必须为字符串",
		"verify_code.required":    "验证码不能为空",
		"verify_code.len":         "验证码长度必须为 6 位",
		"verify_code.number":      "验证码必须为数字",
		"verify_code.verify_code": "验证码错误",
		"captcha.required":        "reCAPTCHA不能为空",
		"captcha.recaptcha":       "reCAPTCHA校验失败（更换网络环境或稍后再试）",
	}
}

func (r *StoreAvatarRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *StoreAvatarRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	_ = data.Set("ip", ctx.Request().Ip())
	_ = data.Set("raw", strings.ToLower(ctx.Request().Input("raw")))
	return nil
}
