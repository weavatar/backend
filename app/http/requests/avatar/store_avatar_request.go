package avatar

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type StoreAvatarRequest struct {
	Raw        string `json:"raw" form:"raw"`
	VerifyCode string `json:"verify_code" form:"verify_code"`
	Avatar     string `json:"avatar" form:"avatar"`

	CaptchaID string `json:"captcha_id" form:"captcha_id"`
	Captcha   string `json:"captcha" form:"captcha"`
}

func (r *StoreAvatarRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *StoreAvatarRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"raw":         "required|string",
		"verify_code": "required|len:6|number|verify_code:raw,avatar",
		"avatar":      "required|image",
		"captcha_id":  "required|string",
		"captcha":     "required|len:6|number|captcha",
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
		"avatar.required":         "头像不能为空",
		"avatar.image":            "头像必须为图片",
		"captcha_id.required":     "图形验证码 ID 不能为空",
		"captcha_id.string":       "图形验证码 ID 必须为字符串",
		"captcha.required":        "图形验证码不能为空",
		"captcha.len":             "图形验证码长度必须为 6 位",
		"captcha.number":          "图形验证码必须为数字",
		"captcha.captcha":         "图形验证码错误",
	}
}

func (r *StoreAvatarRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *StoreAvatarRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
