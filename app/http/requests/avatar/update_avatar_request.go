package avatar

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateAvatarRequest struct {
	Avatar string `form:"avatar" json:"avatar"`

	CaptchaID string `form:"captcha_id" json:"captcha_id"`
	Captcha   string `form:"captcha" json:"captcha"`
}

func (r *UpdateAvatarRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateAvatarRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"avatar":     "required|image",
		"captcha_id": "required|string",
		"captcha":    "required|len:6|number|captcha",
	}
}

func (r *UpdateAvatarRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"avatar.required":     "头像不能为空",
		"avatar.image":        "头像必须为图片",
		"captcha_id.required": "图形验证码 ID 不能为空",
		"captcha_id.string":   "图形验证码 ID 必须为字符串",
		"captcha.required":    "图形验证码不能为空",
		"captcha.len":         "图形验证码长度必须为 6 位",
		"captcha.number":      "图形验证码必须为数字",
		"captcha.captcha":     "图形验证码错误",
	}
}

func (r *UpdateAvatarRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateAvatarRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
