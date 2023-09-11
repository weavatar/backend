package config

import (
	fiberfacades "github.com/goravel/fiber/facades"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("http", map[string]any{
		// HTTP Driver
		"default": "fiber",
		// HTTP Drivers
		"drivers": map[string]any{
			"fiber": map[string]any{
				"route": func() (route.Route, error) {
					return fiberfacades.Route("fiber"), nil
				},
			},
		},
		// HTTP URL
		"url": config.Env("APP_URL", "http://localhost"),
		// HTTP Host
		"host": config.Env("APP_HOST", "127.0.0.1"),
		// HTTP Port
		"port": config.Env("APP_PORT", "3000"),
		// HTTPS Configuration
		"tls": map[string]any{
			// HTTPS Host
			"host": config.Env("APP_HOST", "127.0.0.1"),
			// HTTPS Port
			"port": config.Env("APP_PORT", "3000"),
			// SSL Certificate, you can put the certificate in /public folder
			"ssl": map[string]any{
				// ca.pem
				"cert": "",
				// ca.key
				"key": "",
			},
		},
	})
}
