package controllers

import (
	"os"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	requests "weavatar/app/http/requests/avatar"
	"weavatar/app/models"
	"weavatar/app/services"
	packagecdn "weavatar/pkg/cdn"
	"weavatar/pkg/helper"
)

type AvatarController struct {
	avatar services.Avatar
}

func NewAvatarController() *AvatarController {
	return &AvatarController{
		avatar: services.NewAvatarImpl(),
	}
}

// Avatar 获取头像
func (r *AvatarController) Avatar(ctx http.Context) http.Response {
	appid, hash, imageExt, size, forceDefault, defaultAvatar := r.avatar.Sanitize(ctx)

	var avatar []byte
	var lastModified carbon.Carbon
	var option []string
	var err error
	from := "weavatar"

	if defaultAvatar == "letter" {
		option = []string{ctx.Request().Input("letter"), hash}
	} else {
		option = []string{hash}
	}

	if forceDefault {
		avatar, lastModified, err = r.avatar.GetDefaultByType(defaultAvatar, option)
	} else {
		avatar, lastModified, from, err = r.avatar.GetAvatar(appid, hash, defaultAvatar, option)
	}

	// 判断一下 404 请求
	if avatar == nil && defaultAvatar == "404" {
		return ctx.Response().String(http.StatusNotFound, "404 Not Found\nWeAvatar")
	}

	// 判断一下默认头像 302 请求
	if avatar == nil && helper.IsURL(defaultAvatar) {
		return ctx.Response().Redirect(http.StatusFound, defaultAvatar)
	}

	if err != nil || avatar == nil {
		facades.Log().Error("WeAvatar[获取头像错误] ", err, avatar)
		return ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
	}

	img, err := vips.NewImageFromBuffer(avatar)
	if err != nil {
		facades.Log().Error("WeAvatar[生成头像错误] ", err.Error())
		return ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
	}
	err = img.Thumbnail(size, size, vips.InterestingCentre)
	if err != nil {
		facades.Log().Error("WeAvatar[缩放头像错误] ", err.Error())
		return ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
	}

	var imageData []byte
	switch imageExt {
	case "webp":
		imageData, _, err = img.ExportWebp(vips.NewWebpExportParams())
	case "png":
		imageData, _, err = img.ExportPng(vips.NewPngExportParams())
	case "jpg", "jpeg":
		imageData, _, err = img.ExportJpeg(vips.NewJpegExportParams())
	case "gif":
		imageData, _, err = img.ExportGIF(vips.NewGifExportParams())
	case "tiff":
		imageData, _, err = img.ExportTiff(vips.NewTiffExportParams())
	case "heif":
		imageData, _, err = img.ExportHeif(vips.NewHeifExportParams())
	case "avif":
		imageData, _, err = img.ExportAvif(vips.NewAvifExportParams())
	default:
		imageData, _, err = img.ExportWebp(vips.NewWebpExportParams())
	}
	if err != nil {
		facades.Log().Error("WeAvatar[编码头像错误] ", err.Error())
		return ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
	}

	ctx.Response().Header("Cache-Control", "public, max-age=300")
	ctx.Response().Header("Avatar-By", "weavatar.com")
	ctx.Response().Header("Avatar-From", from)
	ctx.Response().Header("Last-Modified", lastModified.SubHours(8).SetTimezone(carbon.GMT).ToRfc7231String())
	ctx.Response().Header("Expires", carbon.Now().SetTimezone(carbon.GMT).AddMinutes(5).ToRfc7231String())

	return ctx.Response().Data(http.StatusOK, imageExt, imageData)
}

// Index 获取头像列表
func (r *AvatarController) Index(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth().User(ctx, &user); err != nil {
		return Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	page := ctx.Request().QueryInt("page", 1)
	limit := ctx.Request().QueryInt("limit", 10)

	var avatars []models.Avatar
	var total int64
	err := facades.Orm().Query().Where("user_id", user.ID).Paginate(page, limit, &avatars, &total)
	if err != nil {
		facades.Log().Error("[AvatarController][Index] 查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	return Success(ctx, http.Json{
		"total": total,
		"items": avatars,
	})
}

// Show 获取头像详情
func (r *AvatarController) Show(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth().User(ctx, &user); err != nil {
		return Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	var avatar models.Avatar
	err := facades.Orm().Query().Where("user_id", user.ID).Where("hash", ctx.Request().Input("id")).First(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Show] 查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if avatar.Hash == nil {
		return Error(ctx, http.StatusNotFound, "头像不存在")
	}

	return Success(ctx, avatar)
}

// Store 添加头像
func (r *AvatarController) Store(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth().User(ctx, &user); err != nil {
		Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	var storeAvatarRequest requests.StoreAvatarRequest
	sanitize := Sanitize(ctx, &storeAvatarRequest)
	if sanitize != nil {
		return sanitize
	}

	upload, uploadErr := ctx.Request().File("avatar")
	if uploadErr != nil {
		facades.Log().Error("[AvatarController][Store] 解析上传失败 ", uploadErr.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	file, fileErr := os.ReadFile(upload.File())
	if fileErr != nil {
		facades.Log().Error("[AvatarController][Store] 读取上传失败 ", fileErr.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	img, err := vips.NewImageFromBuffer(file)
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 解析图片失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "无法解析图片")
	}
	if img.Width() != img.Height() {
		return Error(ctx, http.StatusUnprocessableEntity, "图片长宽必须相等")
	}

	var avatar models.Avatar
	hash := helper.MD5(storeAvatarRequest.Raw)
	err = facades.Orm().Query().FirstOrCreate(&avatar, models.Avatar{Hash: &hash})
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 初始化查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	saveErr := facades.Storage().Put("upload/default/"+hash[:2]+"/"+hash, string(file))
	if saveErr != nil {
		facades.Log().Error("[AvatarController][Store] 保存用户头像失败 ", saveErr.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	avatar.UserID = &user.ID
	avatar.Raw = &storeAvatarRequest.Raw
	avatar.Ban = false
	avatar.Checked = false
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 添加用户头像失败 ", err.Error())
		delErr := facades.Storage().Delete("upload/default/" + hash[:2] + "/" + hash)
		if delErr != nil {
			facades.Log().Error("[AvatarController][Store] 删除用户头像失败 ", delErr.Error())
		}
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	// 刷新缓存
	go func() {
		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
	}()

	return Success(ctx, nil)
}

// Update 更新头像
func (r *AvatarController) Update(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth().User(ctx, &user); err != nil {
		return Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	var updateAvatarRequest requests.UpdateAvatarRequest
	var sanitize = Sanitize(ctx, &updateAvatarRequest)
	if sanitize != nil {
		return sanitize
	}

	hash := ctx.Request().Input("id")
	if len(hash) != 32 {
		return Error(ctx, http.StatusUnprocessableEntity, "头像哈希格式错误")
	}

	var avatar models.Avatar
	err := facades.Orm().Query().Where("hash", hash).Where("user_id", user.ID).First(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Update] 查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if avatar.Hash == nil {
		return Error(ctx, http.StatusNotFound, "头像不存在")
	}

	// 尝试解析图片
	upload, uploadErr := ctx.Request().File("avatar")
	if uploadErr != nil {
		facades.Log().Error("[AvatarController][Update] 解析上传失败 ", uploadErr.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	file, fileErr := os.ReadFile(upload.File())
	if fileErr != nil {
		facades.Log().Error("[AvatarController][Update] 读取上传失败 ", fileErr.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	img, err := vips.NewImageFromBuffer(file)
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 解析图片失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "无法解析图片")
	}
	if img.Width() != img.Height() {
		return Error(ctx, http.StatusUnprocessableEntity, "图片长宽必须相等")
	}

	avatar.Checked = false
	avatar.Ban = false
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Update] 更新用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	saveErr := facades.Storage().Put("upload/default/"+hash[:2]+"/"+hash, string(file))
	if saveErr != nil {
		facades.Log().Error("[AvatarController][Update] 保存用户头像失败 ", saveErr.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	// 刷新缓存
	go func() {
		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
	}()

	return Success(ctx, nil)
}

// Destroy 删除头像
func (r *AvatarController) Destroy(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth().User(ctx, &user); err != nil {
		return Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	hash := ctx.Request().Input("id")
	if len(hash) != 32 {
		return Error(ctx, http.StatusUnprocessableEntity, "头像哈希格式错误")
	}

	var avatar models.Avatar
	err := facades.Orm().Query().Where("hash", hash).Where("user_id", user.ID).First(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Destroy] 查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if avatar.Hash == nil {
		return Error(ctx, http.StatusNotFound, "头像不存在")
	}

	avatar.Checked = false
	avatar.Ban = false
	avatar.UserID = nil
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Destroy] 删除用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	delErr := facades.Storage().Delete("upload/default/" + hash[:2] + "/" + hash)
	if delErr != nil {
		facades.Log().Error("[AvatarController][Destroy] 删除用户头像失败 ", delErr.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	// 刷新缓存
	go func() {
		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
	}()

	return Success(ctx, nil)
}
