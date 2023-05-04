package requests

import (
	"errors"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type SmsRequest struct {
	Phone  string `json:"phone" form:"phone"`
	UseFor string `json:"use_for" form:"use_for"`

	CaptchaID string `json:"captcha_id" form:"captcha_id"`
	Captcha   string `json:"captcha" form:"captcha"`
}

func (r *SmsRequest) Authorize(ctx http.Context) error {
	if facades.Cache.Has("verify_code:" + r.Phone) {
		return errors.New("发送过于频繁，请稍后再试")
	}
	return nil
}

func (r *SmsRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"phone":      "required|len:11|number|phone",
		"use_for":    "required|in:avatar",
		"captcha_id": "required|string",
		"captcha":    "required|len:6|number|captcha",
	}
}

func (r *SmsRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"phone.required":      "手机号不能为空",
		"phone.len":           "手机号长度必须为 11 位",
		"phone.number":        "手机号必须为数字",
		"phone.phone":         "手机号格式不正确",
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

func (r *SmsRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *SmsRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
