package bootstrap

import (
	"github.com/goravel/framework/foundation"

	"weavatar/config"
)

func Boot() {
	app := foundation.NewApplication()

	// 框架，启动！
	app.Boot()

	// 配置，启动！
	config.Boot()
}
