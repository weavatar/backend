package bootstrap

import (
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/goravel/framework/foundation"

	"weavatar/config"
)

func Boot() {
	app := foundation.NewApplication()

	// Bootstrap the application
	app.Boot()

	// Bootstrap the config.
	config.Boot()

	// Bootstrap the vips.
	vips.LoggingSettings(nil, vips.LogLevelError)
	vips.Startup(nil)
}
