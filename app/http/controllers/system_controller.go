package controllers

import (
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	"weavatar/app/models"
	packagecdn "weavatar/pkg/cdn"
	"weavatar/pkg/helper"
)

type SystemController struct {
	// Dependent services
}

func NewSystemController() *SystemController {
	return &SystemController{
		// Inject services
	}
}

// CdnUsage 获取CDN使用情况
func (r *SystemController) CdnUsage(ctx http.Context) http.Response {
	yesterday := carbon.Now().SubDay().StartOfDay()
	today := carbon.Now().StartOfDay()
	domain := "weavatar.com"

	// 先判断下有没有缓存
	usage := facades.Cache().GetInt64("cdn_usage", -1)
	if usage != -1 {
		return Success(ctx, http.Json{
			"usage": usage,
		})
	}

	cdn := packagecdn.NewCDN()
	usage = int64(cdn.GetUsage(domain, yesterday, today))
	if usage != 0 {
		cacheTime := time.Duration(carbon.Now().EndOfDay().Timestamp() - carbon.Now().Timestamp())
		err := facades.Cache().Put("cdn_usage", usage, cacheTime*time.Second)
		if err != nil {
			facades.Log().Error("[SystemController][CdnUsage] 缓存CDN使用情况失败 " + err.Error())
		}
	}

	return Success(ctx, http.Json{
		"usage": usage,
	})
}

// CheckBind 检查绑定
func (r *SystemController) CheckBind(ctx http.Context) http.Response {
	raw := ctx.Request().Input("raw", "12345")
	hash := helper.MD5(raw)

	var avatar models.Avatar
	if err := facades.Orm().Query().Where("hash", hash).First(&avatar); err != nil {
		facades.Log().Error("[AvatarController][CheckBind] 查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if avatar.UserID == 0 {
		return Success(ctx, http.Json{
			"bind": false,
		})
	}

	return Success(ctx, http.Json{
		"bind": true,
	})
}
