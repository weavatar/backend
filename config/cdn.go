package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("cdn", map[string]any{
		"driver": "upyun",
		// 星盾
		"starshield": map[string]any{
			"access_key": config.Env("JDCLOUD_ACCESS_KEY", ""),
			"secret_key": config.Env("JDCLOUD_SECRET_KEY", ""),
			"zone_id":    config.Env("JDCLOUD_STARSHIELD_ZONE_ID", ""),
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
