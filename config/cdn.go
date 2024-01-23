package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("cdn", map[string]any{
		"driver": config.Env("CDN_DRIVER", "ctyun"),
		// 天翼云
		"ctyun": map[string]any{
			"app_id":     config.Env("CDN_CTYUN_APP_ID", ""),
			"app_secret": config.Env("CDN_CTYUN_APP_SECRET", ""),
		},
		// 网宿
		"wangsu": map[string]any{
			"access_key": config.Env("CDN_WANGSU_ACCESS_KEY", ""),
			"secret_key": config.Env("CDN_WANGSU_SECRET_KEY", ""),
		},
		// 星盾
		"starshield": map[string]any{
			"access_key":  config.Env("CDN_STARSHIELD_ACCESS_KEY", ""),
			"secret_key":  config.Env("CDN_STARSHIELD_SECRET_KEY", ""),
			"instance_id": config.Env("CDN_STARSHIELD_INSTANCE_ID", ""),
			"zone_id":     config.Env("CDN_STARSHIELD_ZONE_ID", ""),
		},
		// 白山云
		"baishan": map[string]any{
			"token": config.Env("CDN_BAISHAN_TOKEN", ""),
		},
		// 云盾
		"yundun": map[string]any{
			"username": config.Env("CDN_YUNDUN_USERNAME", ""),
			"password": config.Env("CDN_YUNDUN_PASSWORD", ""),
		},
		// CloudFlare
		"cloudflare": map[string]any{
			"key":     config.Env("CDN_CLOUDFLARE_KEY", ""),
			"email":   config.Env("CDN_CLOUDFLARE_EMAIL", ""),
			"zone_id": config.Env("CDN_CLOUDFLARE_ZONE_ID", ""),
		},
	})
}
