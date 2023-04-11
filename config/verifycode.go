package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("verifycode", map[string]any{
		// 验证码的长度
		"code_length": 6,

		// 过期时间，单位是分钟
		"expire_time": 5,

		// debug 模式下的过期时间，方便本地开发调试
		"debug_expire_time": 10080,

		// 本地开发环境验证码使用 debug_code
		"debug_code": 123456,
	})
}
