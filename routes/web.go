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
		route.Middleware(frameworkmiddleware.Throttle("captcha")).Get("image", captchaController.Image)
		route.Middleware(frameworkmiddleware.Throttle("verify_code")).Post("sms", captchaController.Sms)
		route.Middleware(frameworkmiddleware.Throttle("verify_code")).Post("email", captchaController.Email)
	})
	facades.Route.Prefix("user").Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Route) {
		userController := controllers.NewUserController()
		route.Post("oauthLogin", userController.OauthLogin)
		route.Post("oauthCallback", userController.OauthCallback)
		route.Middleware(middleware.Jwt()).Get("info", userController.Info)
		route.Middleware(middleware.Jwt()).Post("updateProfile", userController.UpdateProfile)
		route.Middleware(middleware.Jwt()).Post("logout", userController.Logout)
		route.Middleware(middleware.Jwt()).Post("refresh", userController.Refresh)
	})
	facades.Route.Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Route) {
		//route.Middleware(middleware.Jwt()).Resource("avatars", avatarController)
		route.Middleware(middleware.Jwt()).Get("avatars", avatarController.Index)
		route.Middleware(middleware.Jwt()).Post("avatars", avatarController.Store)
		route.Middleware(middleware.Jwt()).Get("avatars/{id}", avatarController.Show)
		route.Middleware(middleware.Jwt()).Put("avatars/{id}", avatarController.Update)
		route.Middleware(middleware.Jwt()).Patch("avatars/{id}", avatarController.Update)
		route.Middleware(middleware.Jwt()).Delete("avatars/{id}", avatarController.Destroy)
	})
	facades.Route.Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Route) {
		appController := controllers.NewAppController()
		//route.Middleware(middleware.Jwt()).Resource("apps", appController)
		route.Middleware(middleware.Jwt()).Get("apps", appController.Index)
		route.Middleware(middleware.Jwt()).Post("apps", appController.Store)
		route.Middleware(middleware.Jwt()).Get("apps/{id}", appController.Show)
		route.Middleware(middleware.Jwt()).Put("apps/{id}", appController.Update)
		route.Middleware(middleware.Jwt()).Patch("apps/{id}", appController.Update)
		route.Middleware(middleware.Jwt()).Delete("apps/{id}", appController.Destroy)
	})

	facades.Route.Prefix("system").Middleware(frameworkmiddleware.Throttle("global")).Group(func(route route.Route) {
		systemController := controllers.NewSystemController()
		route.Middleware(middleware.Jwt()).Get("checkBind", systemController.CheckBind)
		route.Get("cdnUsage", systemController.CdnUsage)
	})
}
