package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	frameworkmiddleware "github.com/goravel/framework/http/middleware"

	"weavatar/app/http/controllers"
	"weavatar/app/http/middleware"
)

func Api() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().String(http.StatusOK, "WeAvatar API")
	})

	avatarController := controllers.NewAvatarController()
	facades.Route().Get("avatar/{hash}", avatarController.Avatar) // 用于获取头像
	facades.Route().Get("avatar", avatarController.Avatar)        // 用于获取头像
	facades.Route().Get("avatar/", avatarController.Avatar)       // 用于获取头像

	facades.Route().Prefix("captcha").Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Router) {
		captchaController := controllers.NewCaptchaController()
		route.Middleware(frameworkmiddleware.Throttle("captcha")).Get("image", captchaController.Image)
		route.Middleware(frameworkmiddleware.Throttle("verify_code")).Post("sms", captchaController.Sms)
		route.Middleware(frameworkmiddleware.Throttle("verify_code")).Post("email", captchaController.Email)
	})
	facades.Route().Prefix("user").Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Router) {
		userController := controllers.NewUserController()
		route.Post("oauthLogin", userController.OauthLogin)
		route.Post("oauthCallback", userController.OauthCallback)
		route.Middleware(middleware.Jwt()).Get("info", userController.GetInfo)
		route.Middleware(middleware.Jwt()).Put("info", userController.UpdateInfo)
		route.Middleware(middleware.Jwt()).Post("logout", userController.Logout)
		route.Middleware(middleware.Jwt()).Get("qqAvatar", userController.GetQQAvatar)
	})
	facades.Route().Middleware(frameworkmiddleware.Throttle("global"), middleware.Jwt()).Group(func(route route.Router) {
		route.Resource("avatars", avatarController)
	})
	/*facades.Route().Middleware(frameworkmiddleware.Throttle("global"), middleware.Jwt()).Group(func(route route.Router) {
		appController := controllers.NewAppController()
		route.Resource("apps", appController)
	})*/

	facades.Route().Prefix("system").Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Router) {
		systemController := controllers.NewSystemController()
		route.Middleware(middleware.Jwt()).Get("checkBind", systemController.CheckBind)
		route.Get("cdnUsage", systemController.CdnUsage)
	})

	facades.Route().Fallback(func(ctx http.Context) http.Response {
		return ctx.Response().Status(http.StatusNotFound).String(http.StatusText(http.StatusNotFound))
	})
}
