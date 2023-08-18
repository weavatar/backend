package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("verifycode", map[string]any{
		// 验证码的长度
		"code_length": 6,
		// 过期时间，单位是分钟
		"expire_time": 5,
		// debug 模式下的过期时间
		"debug_expire_time": 10080,
		// 开发环境固定验证码
		"debug_code": 123456,
	})
}
