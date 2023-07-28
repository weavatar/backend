package controllers

import (
	"bytes"
	"image"
	"os"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	_ "golang.org/x/image/webp"

	requests "weavatar/app/http/requests/avatar"
	"weavatar/app/models"
	"weavatar/app/services"
	packagecdn "weavatar/pkg/cdn"
	"weavatar/pkg/helpers"
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
func (c *AvatarController) Avatar(ctx http.Context) {
	appid, hash, imageExt, size, forceDefault, defaultAvatar := c.avatar.Sanitize(ctx)

	carbon.SetTimezone(carbon.GMT)

	var avatar image.Image
	var lastModified carbon.Carbon
	var option []string
	var err error
	from := "weavatar"

	if defaultAvatar == "letter" {
		option = []string{strings.Trim(ctx.Request().Input("letter"), " "), hash}
	} else {
		option = []string{hash}
	}

	if forceDefault {
		avatar, lastModified, err = c.avatar.GetDefaultAvatarByType(defaultAvatar, option)
	} else {
		avatar, lastModified, from, err = c.avatar.GetAvatar(appid, hash, defaultAvatar, option)
	}

	// 判断一下 404 请求
	if avatar == nil && defaultAvatar == "404" {
		ctx.Response().String(http.StatusNotFound, "404 Not Found\nWeAvatar")
		return
	}

	if err != nil || avatar == nil {
		facades.Log().Error("WeAvatar[获取头像错误] ", err, avatar)
		ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
		return
	}

	img := imaging.Resize(avatar, size, size, imaging.Lanczos)
	imageData, imgErr := c.encodeImage(img, imageExt)
	if imgErr != nil {
		facades.Log().Error("WeAvatar[生成头像错误] ", imgErr.Error())
		ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
		return
	}

	ctx.Response().Header("Cache-Control", "public, max-age=300")
	ctx.Response().Header("Avatar-By", "weavatar.com")
	ctx.Response().Header("Avatar-From", from)
	ctx.Response().Header("Last-Modified", lastModified.ToRfc7231String())
	ctx.Response().Header("Expires", carbon.Now().AddMinutes(5).ToRfc7231String())

	ctx.Response().Data(http.StatusOK, "image/"+imageExt, imageData)

	carbon.SetTimezone(carbon.PRC)
}

// encodeImage 编码图片为指定格式
func (c *AvatarController) encodeImage(img image.Image, imageExt string) ([]byte, error) {
	var err error
	writer := bytes.NewBuffer([]byte{})

	switch imageExt {
	case "webp":
		err = webp.Encode(writer, img, &webp.Options{Lossless: true})
	case "png":
		err = imaging.Encode(writer, img, imaging.PNG)
	case "jpg", "jpeg":
		err = imaging.Encode(writer, img, imaging.JPEG)
	case "gif":
		err = imaging.Encode(writer, img, imaging.GIF)
	default:
		err = imaging.Encode(writer, img, imaging.PNG)
	}

	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

// Index 获取头像列表
func (c *AvatarController) Index(ctx http.Context) {
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		Error(ctx, http.StatusUnauthorized, "登录已过期")
		return
	}

	var avatars []models.Avatar
	err := facades.Orm().Query().Where("user_id", user.ID).Find(&avatars)
	if err != nil {
		facades.Log().Error("[AvatarController][Index] 查询用户头像失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, avatars)
}

// Show 获取头像详情
func (c *AvatarController) Show(ctx http.Context) {
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		Error(ctx, http.StatusUnauthorized, "登录已过期")
		return
	}

	var avatar models.Avatar
	err := facades.Orm().Query().Where("user_id", user.ID).Where("hash", ctx.Request().Input("id")).First(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Show] 查询用户头像失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	if avatar.Hash == nil {
		Error(ctx, http.StatusNotFound, "头像不存在")
		return
	}

	Success(ctx, avatar)
}

// Store 添加头像
func (c *AvatarController) Store(ctx http.Context) {
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		Error(ctx, http.StatusUnauthorized, "登录已过期")
		return
	}

	var storeAvatarRequest requests.StoreAvatarRequest
	errors, err := ctx.Request().ValidateRequest(&storeAvatarRequest)
	if err != nil {
		Error(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if errors != nil {
		Error(ctx, http.StatusUnprocessableEntity, errors.All())
		return
	}

	// 尝试解析图片
	upload, uploadErr := ctx.Request().File("avatar")
	if uploadErr != nil {
		facades.Log().Error("[AvatarController][Store] 解析上传失败 ", uploadErr.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	file, fileErr := os.ReadFile(upload.File())
	if fileErr != nil {
		facades.Log().Error("[AvatarController][Store] 读取上传失败 ", fileErr.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	decode, decodeErr := imaging.Decode(bytes.NewReader(file))
	if decodeErr != nil {
		facades.Log().Error("[AvatarController][Store] 解析图片失败 ", decodeErr.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	// 判断图片长宽是否符合要求
	if decode.Bounds().Dx() != decode.Bounds().Dy() {
		Error(ctx, http.StatusUnprocessableEntity, "图片长宽必须相等")
		return
	}

	var avatar models.Avatar
	hash := helpers.MD5(storeAvatarRequest.Raw)
	_, err = facades.Orm().Query().Exec(`INSERT INTO avatars (hash, created_at, updated_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE updated_at=VALUES(updated_at)`, hash, carbon.DateTime{Carbon: carbon.Now()}, carbon.DateTime{Carbon: carbon.Now()})
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 初始化查询用户头像失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}
	err = facades.Orm().Query().Where("hash", hash).First(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 查询用户头像失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	saveErr := facades.Storage().Put("upload/default/"+hash[:2]+"/"+hash, string(file))
	if saveErr != nil {
		facades.Log().Error("[AvatarController][Store] 保存用户头像失败 ", saveErr.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	avatar.UserID = &user.ID
	avatar.Raw = &storeAvatarRequest.Raw
	avatar.Ban = false
	avatar.Checked = false
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Store] 添加用户头像失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		delErr := facades.Storage().Delete("upload/default/" + hash[:2] + "/" + hash)
		if delErr != nil {
			facades.Log().Error("[AvatarController][Store] 删除用户头像失败 ", delErr.Error())
		}
		return
	}

	// 刷新缓存
	go func() {
		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
	}()

	Success(ctx, nil)
}

// Update 更新头像
func (c *AvatarController) Update(ctx http.Context) {
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		Error(ctx, http.StatusUnauthorized, "登录已过期")
		return
	}

	var updateAvatarRequest requests.UpdateAvatarRequest
	errors, err := ctx.Request().ValidateRequest(&updateAvatarRequest)
	if err != nil {
		Error(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if errors != nil {
		Error(ctx, http.StatusUnprocessableEntity, errors.All())
		return
	}

	hash := ctx.Request().Input("id")
	if len(hash) != 32 {
		Error(ctx, http.StatusUnprocessableEntity, "头像哈希格式错误")
		return
	}

	var avatar models.Avatar
	err = facades.Orm().Query().Where("hash", hash).Where("user_id", user.ID).First(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Update] 查询用户头像失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	if avatar.Hash == nil {
		Error(ctx, http.StatusNotFound, "头像不存在")
		return
	}

	// 尝试解析图片
	upload, uploadErr := ctx.Request().File("avatar")
	if uploadErr != nil {
		facades.Log().Error("[AvatarController][Update] 解析上传失败 ", uploadErr.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	file, fileErr := os.ReadFile(upload.File())
	if fileErr != nil {
		facades.Log().Error("[AvatarController][Update] 读取上传失败 ", fileErr.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	decode, decodeErr := imaging.Decode(bytes.NewReader(file))
	if decodeErr != nil {
		facades.Log().Error("[AvatarController][Update] 解析图片失败 ", decodeErr.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	// 判断图片长宽是否符合要求
	if decode.Bounds().Dx() != decode.Bounds().Dy() {
		Error(ctx, http.StatusUnprocessableEntity, "图片长宽必须相等")
		return
	}

	avatar.Checked = false
	avatar.Ban = false
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Update] 更新用户头像失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	saveErr := facades.Storage().Put("upload/default/"+hash[:2]+"/"+hash, string(file))
	if saveErr != nil {
		facades.Log().Error("[AvatarController][Update] 保存用户头像失败 ", saveErr.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	// 刷新缓存
	go func() {
		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
	}()

	Success(ctx, nil)
}

// Destroy 删除头像
func (c *AvatarController) Destroy(ctx http.Context) {
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		Error(ctx, http.StatusUnauthorized, "登录已过期")
		return
	}

	hash := ctx.Request().Input("id")
	if len(hash) != 32 {
		Error(ctx, http.StatusUnprocessableEntity, "头像哈希格式错误")
		return
	}

	var avatar models.Avatar
	err := facades.Orm().Query().Where("hash", hash).Where("user_id", user.ID).First(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Destroy] 查询用户头像失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	if avatar.Hash == nil {
		Error(ctx, http.StatusNotFound, "头像不存在")
		return
	}

	avatar.Checked = false
	avatar.Ban = false
	avatar.UserID = nil
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("[AvatarController][Destroy] 删除用户头像失败 ", err.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	delErr := facades.Storage().Delete("upload/default/" + hash[:2] + "/" + hash)
	if delErr != nil {
		facades.Log().Error("[AvatarController][Destroy] 删除用户头像失败 ", delErr.Error())
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	// 刷新缓存
	go func() {
		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
	}()

	Success(ctx, nil)
}
