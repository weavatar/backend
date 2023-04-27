package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	frameworkmiddleware "github.com/goravel/framework/http/middleware"

	"weavatar/app/http/controllers"
	"weavatar/app/http/middleware"
)

func Web() {
	facades.Route.Get("/", func(ctx http.Context) {
		ctx.Response().String(http.StatusOK, "WeAvatar API")
	})

	avatarController := controllers.NewAvatarController()
	facades.Route.Get("avatar/{hash}", avatarController.Avatar) // 用于获取头像
	facades.Route.Get("avatar", avatarController.Avatar)        // 用于获取头像
	facades.Route.Get("avatar/", avatarController.Avatar)       // 用于获取头像

	facades.Route.Prefix("captcha").Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Route) {
		captchaController := controllers.NewCaptchaController()
		route.Get("image", captchaController.Image)
		route.Post("sms", captchaController.Sms)
		route.Post("email", captchaController.Email)
	})
	facades.Route.Prefix("user").Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Route) {
		userController := controllers.NewUserController()
		route.Post("oauthLogin", userController.OauthLogin)
		route.Post("oauthCallback", userController.OauthCallback)
		route.Middleware(middleware.Jwt()).Post("updateNickname", userController.UpdateNickname)
		route.Middleware(middleware.Jwt()).Post("logout", userController.Logout)
		route.Middleware(middleware.Jwt()).Post("refresh", userController.Refresh)
	})
	facades.Route.Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Route) {
		route.Middleware(middleware.Jwt()).Resource("avatars", avatarController)
		route.Middleware(middleware.Jwt()).Post("avatars/checkBind", avatarController.CheckBind)
	})
	facades.Route.Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Route) {
		appController := controllers.NewAppController()
		route.Middleware(middleware.Jwt()).Resource("apps", appController)
	})
}
