package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/http/limit"

	"weavatar/app/http"
	"weavatar/routes"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Register(app foundation.Application) {
}

func (receiver *RouteServiceProvider) Boot(app foundation.Application) {
	// Add HTTP middlewares
	facades.Route().GlobalMiddleware(http.Kernel{}.Middleware()...)

	receiver.configureRateLimiting()

	routes.Api()
}

func (receiver *RouteServiceProvider) configureRateLimiting() {
	facades.RateLimiter().For("global", func(ctx contractshttp.Context) contractshttp.Limit {
		return limit.PerMinute(1000).Response(func(ctx contractshttp.Context) {
			ctx.Request().AbortWithStatusJson(contractshttp.StatusTooManyRequests, contractshttp.Json{
				"code":    contractshttp.StatusTooManyRequests,
				"message": "达到请求上限，请稍后再试",
			})
		})
	})
	facades.RateLimiter().ForWithLimits("verify_code", func(ctx contractshttp.Context) []contractshttp.Limit {
		return []contractshttp.Limit{
			limit.PerMinute(2).Response(func(ctx contractshttp.Context) {
				ctx.Request().AbortWithStatusJson(contractshttp.StatusTooManyRequests, contractshttp.Json{
					"code":    contractshttp.StatusTooManyRequests,
					"message": "达到请求上限，请稍后再试",
				})
			}),
			limit.PerDay(50).By(ctx.Request().Ip()).Response(func(ctx contractshttp.Context) {
				ctx.Request().AbortWithStatusJson(contractshttp.StatusTooManyRequests, contractshttp.Json{
					"code":    contractshttp.StatusTooManyRequests,
					"message": "达到请求上限，请明天再试",
				})
			}),
		}
	})
	facades.RateLimiter().For("captcha", func(ctx contractshttp.Context) contractshttp.Limit {
		return limit.PerMinute(30).Response(func(ctx contractshttp.Context) {
			ctx.Request().AbortWithStatusJson(contractshttp.StatusTooManyRequests, contractshttp.Json{
				"code":    contractshttp.StatusTooManyRequests,
				"message": "达到请求上限，请稍后再试",
			})
		})
	})
}
