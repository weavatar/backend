package avatar

import (
	"github.com/goravel-kit/geetest"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/support/json"
	"github.com/spf13/cast"
)

type UpdateAvatarRequest struct {
	Captcha geetest.Ticket `form:"captcha" json:"captcha"`
}

func (r *UpdateAvatarRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateAvatarRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"captcha": "required|geetest",
	}
}

func (r *UpdateAvatarRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"captcha.geetest": "验证码校验失败（更换设备环境或刷新重试）",
	}
}

func (r *UpdateAvatarRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateAvatarRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	ticket := make(map[string]string)
	if captcha, exist := data.Get("captcha"); exist {
		if err := json.UnmarshalString(cast.ToString(captcha), &ticket); err != nil {
			return err
		}
		if err := data.Set("captcha", ticket); err != nil {
			return err
		}
	}

	return nil
}
