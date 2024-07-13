package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("cdn", map[string]any{
		"driver": config.Env("CDN_DRIVER", "baishan,huawei"),
		// 白山云
		"baishan": map[string]any{
			"token": config.Env("CDN_BAISHAN_TOKEN", ""),
		},
		// 华为云
		"huawei": map[string]any{
			"access_key": config.Env("CDN_HUAWEI_ACCESS_KEY", ""),
			"secret_key": config.Env("CDN_HUAWEI_SECRET_KEY", ""),
		},
		// 括彩云
		"kuocai": map[string]any{
			"username": config.Env("CDN_KUOCAI_USERNAME", ""),
			"password": config.Env("CDN_KUOCAI_PASSWORD", ""),
		},
		// CloudFlare
		"cloudflare": map[string]any{
			"key":     config.Env("CDN_CLOUDFLARE_KEY", ""),
			"email":   config.Env("CDN_CLOUDFLARE_EMAIL", ""),
			"zone_id": config.Env("CDN_CLOUDFLARE_ZONE_ID", ""),
		},
	})
}
