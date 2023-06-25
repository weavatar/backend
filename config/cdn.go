package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("cdn", map[string]any{
		"driver": config.Env("CDN_DRIVER", "starshield"),
		// 星盾
		"starshield": map[string]any{
			"access_key":  config.Env("CDN_STARSHIELD_ACCESS_KEY", ""),
			"secret_key":  config.Env("CDN_STARSHIELD_SECRET_KEY", ""),
			"instance_id": config.Env("CDN_STARSHIELD_INSTANCE_ID", ""),
			"zone_id":     config.Env("CDN_STARSHIELD_ZONE_ID", ""),
		},
		// 又拍云
		"upyun": map[string]any{
			"token": config.Env("CDN_UPYUN_TOKEN", ""),
		},
		// 盾云CDN
		"ddun": map[string]any{
			"api_key":    config.Env("CDN_DDUN_API_KEY", ""),
			"api_secret": config.Env("CDN_DDUN_API_SECRET", ""),
		},
		// 上海云盾
		"yundun": map[string]any{
			"username": config.Env("CDN_YUNDUN_USERNAME", ""),
			"password": config.Env("CDN_YUNDUN_PASSWORD", ""),
		},
	})
}
