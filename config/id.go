package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("id", map[string]any{
		// 默认的节点 ID
		"node": config.Env("APP_NODE", "0"),
	})
}
