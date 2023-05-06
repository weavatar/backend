package controllers

import (
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	packagecdn "weavatar/packages/cdn"
)

type SystemController struct {
	//Dependent services
}

func NewSystemController() *SystemController {
	return &SystemController{
		//Inject services
	}
}

// CdnUsage 获取CDN使用情况
func (r *SystemController) CdnUsage(ctx http.Context) {
	// 取昨日0点时间
	yesterday := carbon.Now().SubDay().StartOfDay().ToDateString()
	// 取今日0点时间
	today := carbon.Now().StartOfDay().ToDateString()
	// 域名
	domain := facades.Config.GetString("http.host")

	// 先判断下有没有缓存
	usage := facades.Cache.GetInt64("cdn_usage", -1)
	if usage != -1 {
		ctx.Response().Json(http.StatusOK, http.Json{
			"code":    0,
			"message": "获取成功",
			"data": http.Json{
				"usage": usage,
			},
		})
		return
	}

	cdn := packagecdn.NewCDN()
	usage = int64(cdn.GetUsage(domain, yesterday, today))

	cacheTime := time.Duration(carbon.Now().EndOfDay().Timestamp() - carbon.Now().Timestamp())
	err := facades.Cache.Put("cdn_usage", usage, cacheTime*time.Second)

	if err != nil {
		facades.Log.Error("[SystemController][CdnUsage] 缓存CDN使用情况失败 " + err.Error())
		ctx.Response().Json(http.StatusOK, http.Json{
			"code":    0,
			"message": "获取成功",
			"data":    usage,
		})
		return
	}

	ctx.Response().Json(http.StatusOK, http.Json{
		"code":    0,
		"message": "获取成功",
		"data":    usage,
	})
}
