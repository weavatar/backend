package bootstrap

import (
	"github.com/goravel/framework/foundation"
	"github.com/goravel/framework/support/carbon"

	"weavatar/config"
)

func Boot() {
	app := foundation.NewApplication()

	// Bootstrap the application
	app.Boot()

	// Bootstrap the config.
	config.Boot()

	// 设置Carbon的时区
	carbon.SetTimezone(carbon.PRC)
}
