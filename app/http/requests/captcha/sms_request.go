package requests

import (
	"errors"

	"github.com/goravel-kit/geetest"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type SmsRequest struct {
	Phone  string `json:"phone" form:"phone"`
	UseFor string `json:"use_for" form:"use_for"`

	Captcha geetest.Ticket `json:"captcha" form:"captcha"`
}

func (r *SmsRequest) Authorize(ctx http.Context) error {
	if facades.Cache().Has("verify_code:" + r.Phone) {
		return errors.New("发送过于频繁，请稍后再试")
	}
	return nil
}

func (r *SmsRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"phone":   "required|len:11|number|phone",
		"use_for": "required|in:avatar",
		"captcha": "required|geetest",
	}
}

func (r *SmsRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"phone.required":   "手机号不能为空",
		"phone.len":        "手机号长度必须为 11 位",
		"phone.number":     "手机号必须为数字",
		"phone.phone":      "手机号格式不正确",
		"use_for.required": "用途不能为空",
		"use_for.in":       "用途不正确",
		"captcha.required": "验证码不能为空",
		"captcha.geetest":  "验证码校验失败（更换设备环境或刷新重试）",
	}
}

func (r *SmsRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *SmsRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
