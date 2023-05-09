package controllers

import (
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"weavatar/app/models"
	packagecdn "weavatar/packages/cdn"
	"weavatar/packages/helpers"
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
	domain := "weavatar.com"

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
	if usage != 0 {
		cacheTime := time.Duration(carbon.Now().EndOfDay().Timestamp() - carbon.Now().Timestamp())
		err := facades.Cache.Put("cdn_usage", usage, cacheTime*time.Second)
		if err != nil {
			facades.Log.Error("[SystemController][CdnUsage] 缓存CDN使用情况失败 " + err.Error())
			ctx.Response().Json(http.StatusOK, http.Json{
				"code":    0,
				"message": "获取成功",
				"data": http.Json{
					"usage": usage,
				},
			})
			return
		}
	}

	ctx.Response().Json(http.StatusOK, http.Json{
		"code":    0,
		"message": "获取成功",
		"data": http.Json{
			"usage": usage,
		},
	})
}

// CheckBind 检查绑定
func (r *SystemController) CheckBind(ctx http.Context) {
	raw := ctx.Request().Input("raw", "12345")
	hash := helpers.MD5(raw)

	var avatar models.Avatar
	err := facades.Orm.Query().Where("hash", hash).First(&avatar)
	if err != nil {
		facades.Log.WithContext(ctx).Error("[AvatarController][CheckBind] 查询用户头像失败 ", err.Error())
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "系统内部错误",
		})
		return
	}

	if avatar.UserID == nil {
		ctx.Response().Json(http.StatusOK, http.Json{
			"code":    0,
			"message": "地址未被其他用户绑定",
			"data": http.Json{
				"bind": false,
			},
		})
		return
	}

	ctx.Response().Json(http.StatusOK, http.Json{
		"code":    0,
		"message": "地址已被其他用户绑定",
		"data": http.Json{
			"bind": true,
		},
	})
}
