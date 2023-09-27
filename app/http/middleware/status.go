package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// Status 检查程序状态
func Status() http.Middleware {
	return func(ctx http.Context) {
		status := facades.Config().GetString("app.status", "main")
		if status != "main" && ctx.Request().Method() != "GET" {
			ctx.Request().AbortWithStatusJson(http.StatusOK, http.Json{
				"code":    503,
				"message": "当前系统运行在热备模式，暂无法操作",
			})

			return
		}

		ctx.Request().Next()
	}
}
