package requests

import (
	"errors"
	"strings"

	"github.com/goravel-kit/geetest"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
)

type EmailRequest struct {
	Email  string `json:"email" form:"email"`
	UseFor string `json:"use_for" form:"use_for"`

	Captcha geetest.Ticket `json:"captcha" form:"captcha"`
}

func (r *EmailRequest) Authorize(ctx http.Context) error {
	if facades.Cache().Has("verify_code:" + r.Email) {
		return errors.New("发送过于频繁，请稍后再试")
	}
	return nil
}

func (r *EmailRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":   "required|email",
		"use_for": "required|in:avatar",
		"captcha": "required|geetest",
	}
}

func (r *EmailRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"email.required":   "邮箱不能为空",
		"email.email":      "邮箱格式不正确",
		"use_for.required": "用途不能为空",
		"use_for.in":       "用途不正确",
		"captcha.required": "验证码不能为空",
		"captcha.geetest":  "验证码校验失败（更换设备环境或刷新重试）",
	}
}

func (r *EmailRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *EmailRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	if email, exist := data.Get("email"); exist {
		email = strings.ToLower(cast.ToString(email))
		_ = data.Set("email", email)
	}

	return nil
}
