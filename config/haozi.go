package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("haozi", map[string]any{
		// 通行证配置
		"account": map[string]interface{}{
			"base_url":      config.Env("HAOZI_ACCOUNT_BASE_URL"),
			"client_id":     config.Env("HAOZI_ACCOUNT_CLIENT_ID"),
			"client_secret": config.Env("HAOZI_ACCOUNT_CLIENT_SECRET"),
		},
	})
}
