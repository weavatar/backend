package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("cdn", map[string]any{
		"driver": "starshield",
		// 星盾
		"starshield": map[string]interface{}{
			"access_key": config.Env("JDCLOUD_ACCESS_KEY", ""),
			"secret_key": config.Env("JDCLOUD_SECRET_KEY", ""),
			"zone_id":    config.Env("JDCLOUD_STARSHIELD_ZONE_ID", ""),
		},
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
