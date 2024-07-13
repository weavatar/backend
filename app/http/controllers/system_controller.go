package controllers

import (
	"strings"
	"time"

	cdnfacades "github.com/goravel-kit/cdn/facades"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	"weavatar/app/models"
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

	data, err := cdnfacades.Cdn().GetUsage(domain, yesterday, today)
	if err != nil {
		return Success(ctx, http.Json{
			"usage": 0,
		})
	}

	usage = int64(data)
	cacheTime := time.Duration(carbon.Now().EndOfDay().Timestamp() - carbon.Now().Timestamp() + 7200)
	if err = facades.Cache().Put("cdn_usage", usage, cacheTime*time.Second); err != nil {
		facades.Log().Error("[SystemController][CdnUsage] 缓存CDN使用情况失败 " + err.Error())
	}

	return Success(ctx, http.Json{
		"usage": usage,
	})
}

// CheckBind 检查绑定
func (r *SystemController) CheckBind(ctx http.Context) http.Response {
	raw := strings.ToLower(ctx.Request().Input("raw"))
	sha256 := helper.SHA256(raw)

	var avatar models.Avatar
	if err := facades.Orm().Query().Where("sha256", sha256).First(&avatar); err != nil {
		facades.Log().Error("[AvatarController][CheckBind] 查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if avatar.UserID == "" {
		return Success(ctx, http.Json{
			"bind": false,
		})
	}

	return Success(ctx, http.Json{
		"bind": true,
	})
}
