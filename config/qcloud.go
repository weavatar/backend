package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("qcloud", map[string]any{
		// 头像审核COS配置
		"cos_check": map[string]interface{}{
			"access_key": config.Env("QCLOUD_COS_CHECK_ACCESS_KEY"),
			"secret_key": config.Env("QCLOUD_COS_CHECK_SECRET_KEY"),
			"bucket":     config.Env("QCLOUD_COS_CHECK_BUCKET"),
		},
	})
}
