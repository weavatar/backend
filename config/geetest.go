package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("geetest", map[string]any{
		"captcha_id":  config.Env("GEETEST_CAPTCHA_ID"),
		"captcha_key": config.Env("GEETEST_CAPTCHA_KEY"),
		"api_url":     config.Env("GEETEST_API_URL", "https://gcaptcha4.geetest.com"),
	})
}
