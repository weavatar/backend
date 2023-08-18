package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("sms", map[string]any{
		"driver": "tencent",
		// 阿里云
		"aliyun": map[string]interface{}{
			"access_key_id":     config.Env("SMS_ALIYUN_ACCESS_ID"),
			"access_key_secret": config.Env("SMS_ALIYUN_ACCESS_SECRET"),
			"sign_name":         config.Env("SMS_ALIYUN_SIGN_NAME"),
			"template_code":     config.Env("SMS_ALIYUN_TEMPLATE_CODE"),
		},
		// 腾讯云
		"tencent": map[string]interface{}{
			"access_key":    config.Env("SMS_TENCENT_ACCESS_KEY"),
			"secret_key":    config.Env("SMS_TENCENT_SECRET_KEY"),
			"sign_name":     config.Env("SMS_TENCENT_SIGN_NAME"),
			"template_code": config.Env("SMS_TENCENT_TEMPLATE_CODE"),
			"sdk_app_id":    config.Env("SMS_TENCENT_SDK_APP_ID"),
		},
	})
}
