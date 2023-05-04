package requests

import (
	"errors"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type EmailRequest struct {
	Email  string `json:"email" form:"email"`
	UseFor string `json:"use_for" form:"use_for"`

	CaptchaID string `json:"captcha_id" form:"captcha_id"`
	Captcha   string `json:"captcha" form:"captcha"`
}

func (r *EmailRequest) Authorize(ctx http.Context) error {
	if facades.Cache.Has("verify_code:" + r.Email) {
		return errors.New("发送过于频繁，请稍后再试")
	}
	return nil
}

func (r *EmailRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":      "required|email",
		"use_for":    "required|in:avatar",
		"captcha_id": "required|string",
		"captcha":    "required|len:6|number|captcha",
	}
}

func (r *EmailRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"email.required":      "邮箱不能为空",
		"email.email":         "邮箱格式不正确",
		"use_for.required":    "用途不能为空",
		"use_for.in":          "用途不正确",
		"captcha_id.required": "图形验证码 ID 不能为空",
		"captcha_id.string":   "图形验证码 ID 必须为字符串",
		"captcha.required":    "图形验证码不能为空",
		"captcha.len":         "图形验证码长度必须为 6 位",
		"captcha.number":      "图形验证码必须为数字",
		"captcha.captcha":     "图形验证码错误",
	}
}

func (r *EmailRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *EmailRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
