package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("cdn", map[string]any{
		"driver": config.Env("CDN_DRIVER", "baishan"),
		// 白山云
		"baishan": map[string]any{
			"token": config.Env("CDN_BAISHAN_TOKEN", ""),
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
		// AnyCast
		"anycast": map[string]any{
			"api_key":    config.Env("CDN_ANYCAST_API_KEY", ""),
			"api_secret": config.Env("CDN_ANYCAST_API_SECRET", ""),
		},
	})
}
