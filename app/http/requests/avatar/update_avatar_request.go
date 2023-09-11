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
		"captcha": "required|recaptcha:avatar",
	}
}

func (r *UpdateAvatarRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"captcha.required":  "reCAPTCHA不能为空",
		"captcha.recaptcha": "reCAPTCHA校验失败（更换网络环境或稍后再试）",
	}
}

func (r *UpdateAvatarRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateAvatarRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	_ = data.Set("ip", ctx.Request().Ip())
	return nil
}
