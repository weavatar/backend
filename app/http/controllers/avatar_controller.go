package controllers

import (
	"bytes"
	"image"
	"os"
	"os/exec"
	"strconv"

	cdnfacades "github.com/goravel-kit/cdn/facades"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	requests "weavatar/app/http/requests/avatar"
	"weavatar/app/models"
	"weavatar/app/services"
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
	appid, hash, ext, size, forceDefault, defaultAvatar := r.avatar.Sanitize(ctx)

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

	// 创建一个临时文件
	file, err := os.CreateTemp("", "weavatar-")
	if err != nil {
		facades.Log().Error("WeAvatar[创建临时文件错误] ", err.Error())
		return ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
	}
	defer os.Remove(file.Name())

	// 写入临时文件
	_, err = file.Write(avatar)
	if err != nil {
		facades.Log().Error("WeAvatar[写入临时文件错误] ", err.Error())
		return ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
	}

	// 调用 vips 处理图片
	output, err := exec.Command("vipsthumbnail", file.Name(), "-s", strconv.Itoa(size), "--smartcrop", "attention", "-o", file.Name()+"."+ext).Output()
	if err != nil {
		facades.Log().Error("WeAvatar[调用 vips 处理图片错误] ", err.Error(), string(output))
		return ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
	}
	defer os.Remove(file.Name() + "." + ext)

	data, err := os.ReadFile(file.Name() + "." + ext)
	if err != nil {
		facades.Log().Error("WeAvatar[读取临时文件错误] ", err.Error())
		return ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
	}

	ctx.Response().Header("Cache-Control", "public, max-age=300")
	ctx.Response().Header("Avatar-By", "weavatar.com")
	ctx.Response().Header("Avatar-From", from)
	ctx.Response().Header("Last-Modified", lastModified.SubHours(8).SetTimezone(carbon.GMT).ToRfc7231String())
	ctx.Response().Header("Expires", carbon.Now().SetTimezone(carbon.GMT).AddHour().ToRfc7231String())

	return ctx.Response().Data(http.StatusOK, "image/"+ext, data)
}

// Index 获取头像列表
func (r *AvatarController) Index(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
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
	if err := facades.Auth(ctx).User(&user); err != nil {
		return Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	var avatar models.Avatar
	err := facades.Orm().Query().Where("user_id", user.ID).Where("sha256", ctx.Request().Input("id")).First(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Show] 查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if avatar.UserID == "" {
		return Error(ctx, http.StatusNotFound, "头像不存在")
	}

	return Success(ctx, avatar)
}

// Store 添加头像
func (r *AvatarController) Store(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	var storeAvatarRequest requests.StoreAvatarRequest
	sanitize := Sanitize(ctx, &storeAvatarRequest)
	if sanitize != nil {
		return sanitize
	}

	upload, err := ctx.Request().File("avatar")
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 解析上传失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	file, err := os.ReadFile(upload.File())
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 读取上传失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 解析图片失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "无法解析图片")
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	if width != height {
		return Error(ctx, http.StatusUnprocessableEntity, "图片长宽必须相等")
	}

	var avatar models.Avatar
	md5 := helper.MD5(storeAvatarRequest.Raw)
	sha256 := helper.SHA256(storeAvatarRequest.Raw)
	err = facades.Orm().Query().FirstOrNew(&avatar, models.Avatar{MD5: md5, SHA256: sha256})
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 初始化查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if err = facades.Storage().Put("upload/default/"+sha256[:2]+"/"+sha256, string(file)); err != nil {
		facades.Log().Error("[AvatarController][Store] 保存用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	avatar.UserID = user.ID
	avatar.Raw = storeAvatarRequest.Raw
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 添加用户头像失败 ", err.Error())
		delErr := facades.Storage().Delete("upload/default/" + sha256[:2] + "/" + sha256)
		if delErr != nil {
			facades.Log().Error("[AvatarController][Store] 删除用户头像失败 ", delErr.Error())
		}
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	go func() {
		if err = cdnfacades.Cdn().RefreshUrl([]string{"weavatar.com/avatar/" + avatar.MD5, "weavatar.com/avatar/" + avatar.SHA256}); err != nil {
			facades.Log().Error("[AvatarController][Store] CDN 刷新失败 ", err.Error())
		}
	}()

	return Success(ctx, nil)
}

// Update 更新头像
func (r *AvatarController) Update(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	var updateAvatarRequest requests.UpdateAvatarRequest
	var sanitize = Sanitize(ctx, &updateAvatarRequest)
	if sanitize != nil {
		return sanitize
	}

	sha256 := ctx.Request().Input("id")
	if len(sha256) != 64 {
		return Error(ctx, http.StatusUnprocessableEntity, "头像哈希格式错误")
	}

	var avatar models.Avatar
	err := facades.Orm().Query().Where("sha256", sha256).Where("user_id", user.ID).First(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Update] 查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if avatar.UserID == "" {
		return Error(ctx, http.StatusNotFound, "头像不存在")
	}

	// 尝试解析图片
	upload, err := ctx.Request().File("avatar")
	if err != nil {
		facades.Log().Error("[AvatarController][Update] 解析上传失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	file, err := os.ReadFile(upload.File())
	if err != nil {
		facades.Log().Error("[AvatarController][Update] 读取上传失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 解析图片失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "无法解析图片")
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	if width != height {
		return Error(ctx, http.StatusUnprocessableEntity, "图片长宽必须相等")
	}

	// 这里保存一下是为了刷新 updated_at
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Update] 更新用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if err = facades.Storage().Put("upload/default/"+sha256[:2]+"/"+sha256, string(file)); err != nil {
		facades.Log().Error("[AvatarController][Update] 保存用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	go func() {
		if err = cdnfacades.Cdn().RefreshUrl([]string{"weavatar.com/avatar/" + avatar.MD5, "weavatar.com/avatar/" + avatar.SHA256}); err != nil {
			facades.Log().Error("[AvatarController][Update] CDN 刷新失败 ", err.Error())
		}
	}()

	return Success(ctx, nil)
}

// Destroy 删除头像
func (r *AvatarController) Destroy(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return Error(ctx, http.StatusUnauthorized, "登录已过期")
	}

	sha256 := ctx.Request().Input("id")
	if len(sha256) != 64 {
		return Error(ctx, http.StatusUnprocessableEntity, "头像哈希格式错误")
	}

	var avatar models.Avatar
	if err := facades.Orm().Query().Where("sha256", sha256).Where("user_id", user.ID).First(&avatar); err != nil {
		facades.Log().Error("[AvatarController][Destroy] 查询用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if avatar.UserID == "" {
		return Error(ctx, http.StatusNotFound, "头像不存在")
	}

	if _, err := facades.Orm().Query().Delete(&avatar); err != nil {
		facades.Log().Error("[AvatarController][Destroy] 删除用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	if err := facades.Storage().Delete("upload/default/" + sha256[:2] + "/" + sha256); err != nil {
		facades.Log().Error("[AvatarController][Destroy] 删除用户头像失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	go func() {
		if err := cdnfacades.Cdn().RefreshUrl([]string{"weavatar.com/avatar/" + avatar.MD5, "weavatar.com/avatar/" + avatar.SHA256}); err != nil {
			facades.Log().Error("[AvatarController][Destroy] CDN 刷新失败 ", err.Error())
		}
	}()

	return Success(ctx, nil)
}
