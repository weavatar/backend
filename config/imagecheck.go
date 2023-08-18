package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("imagecheck", map[string]any{
		"driver": config.Env("CHECK_DRIVER", "aliyun"), // "aliyun", "cos
		"aliyun": map[string]interface{}{
			"access_key_id":     config.Env("CHECK_ALIYUN_ACCESS_ID"),
			"access_key_secret": config.Env("CHECK_ALIYUN_ACCESS_SECRET"),
		},
		"cos": map[string]interface{}{
			"access_key": config.Env("CHECK_COS_ACCESS_KEY"),
			"secret_key": config.Env("CHECK_COS_SECRET_KEY"),
			"bucket":     config.Env("CHECK_COS_BUCKET"),
		},
	})
}
