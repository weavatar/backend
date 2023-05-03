package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("cdn", map[string]any{
		"driver": "ddun",

		// 盾云CDN
		"ddun": map[string]interface{}{
			"api_key":    config.Env("CDN_DDUN_API_KEY", ""),
			"api_secret": config.Env("CDN_DDUN_API_SECRET", ""),
		},
		// 上海云盾
		"yundun": map[string]interface{}{
			"username": config.Env("CDN_YUNDUN_USERNAME", ""),
			"password": config.Env("CDN_YUNDUN_PASSWORD", ""),
		},
	})
}
