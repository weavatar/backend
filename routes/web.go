package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	frameworkmiddleware "github.com/goravel/framework/http/middleware"
	"weavatar/app/http/middleware"

	"weavatar/app/http/controllers"
)

func Web() {
	facades.Route.Get("/", func(ctx http.Context) {
		ctx.Response().String(http.StatusOK, "WeAvatar API")
	})

	avatarController := controllers.NewAvatarController()
	facades.Route.Get("avatar/{hash}", avatarController.Avatar) // 用于获取头像
	facades.Route.Get("avatar", avatarController.Avatar)        // 用于获取头像

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
	facades.Route.Prefix("avatar").Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Route) {

		route.Middleware(middleware.Jwt()).Get("get", avatarController.Get)
		route.Middleware(middleware.Jwt()).Post("create", avatarController.Create)
		route.Middleware(middleware.Jwt()).Post("update", avatarController.Update)
		route.Middleware(middleware.Jwt()).Post("delete", avatarController.Delete)
		route.Middleware(middleware.Jwt()).Get("get/{id}", avatarController.GetSingle)
		route.Middleware(middleware.Jwt()).Post("checkBind", avatarController.CheckBind)
	})
	facades.Route.Prefix("app").Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Route) {
		appController := controllers.NewAppController()
		route.Middleware(middleware.Jwt()).Get("get", appController.Get)
		route.Middleware(middleware.Jwt()).Post("create", appController.Create)
		route.Middleware(middleware.Jwt()).Post("update", appController.Update)
		route.Middleware(middleware.Jwt()).Post("delete", appController.Delete)
		route.Middleware(middleware.Jwt()).Get("get/{id}", appController.GetSingle)
	})
}
