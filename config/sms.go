package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("sms", map[string]any{
		// 默认的短信驱动
		"driver": "tencent",
		// 默认是阿里云的测试 sign_name 和 template_code
		"aliyun": map[string]interface{}{
			"access_key_id":     config.Env("SMS_ALIYUN_ACCESS_ID"),
			"access_key_secret": config.Env("SMS_ALIYUN_ACCESS_SECRET"),
			"sign_name":         config.Env("SMS_ALIYUN_SIGN_NAME"),
			"template_code":     config.Env("SMS_ALIYUN_TEMPLATE_CODE"),
		},

		// 腾讯云
		"tencent": map[string]interface{}{
			"access_key_id":     config.Env("SMS_TENCENT_ACCESS_ID"),
			"access_key_secret": config.Env("SMS_TENCENT_ACCESS_SECRET"),
			"sign_name":         config.Env("SMS_TENCENT_SIGN_NAME"),
			"template_code":     config.Env("SMS_TENCENT_TEMPLATE_CODE"),
			"sdk_app_id":        config.Env("SMS_TENCENT_SDK_APP_ID"),
		},
	})
}
