package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("recaptcha", map[string]any{
		"secret": config.Env("RECAPTCHA_SECRET"),
	})
}
