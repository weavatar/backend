package avatar

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateAvatarRequest struct {
	Captcha string `form:"captcha" json:"captcha"`
}

func (r *UpdateAvatarRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateAvatarRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"avatar":  "required|image",
		"captcha": "required|geetest",
	}
}

func (r *UpdateAvatarRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"avatar.required":  "头像不能为空",
		"avatar.image":     "头像必须为图片",
		"captcha.required": "验证码不能为空",
		"captcha.geetest":  "验证码校验失败（更换设备环境或刷新重试）",
	}
}

func (r *UpdateAvatarRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateAvatarRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
