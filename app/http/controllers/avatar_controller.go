package controllers

import (
	"bytes"
	"image"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"weavatar/app/models"
	"weavatar/app/services"
)

type AvatarController struct {
	//Dependent services
}

func NewAvatarController() *AvatarController {
	return &AvatarController{
		//Inject services
	}
}

// Avatar 获取头像
func (r *AvatarController) Avatar(ctx http.Context) {
	avatarService := services.NewAvatarImpl()
	appid, hash, imageExt, size, forceDefault, defaultAvatar := avatarService.Sanitize(ctx)

	var avatar image.Image
	var option []string
	var err error
	from := "weavatar"

	if defaultAvatar == "letter" {
		option = []string{strings.Trim(ctx.Request().Input("letter"), " "), hash}
	} else {
		option = []string{hash}
	}

	if forceDefault {
		avatar, err = avatarService.GetDefaultAvatarByType(defaultAvatar, option)
	} else {
		avatar, from, err = avatarService.GetAvatar(appid, hash, defaultAvatar, option)
	}

	// 判断一下 404 请求
	if avatar == nil && defaultAvatar == "404" {
		ctx.Response().String(http.StatusNotFound, "404 Not Found\nWeAvatar")
		return
	}

	if err != nil || avatar == nil {
		facades.Log.Error("WeAvatar[获取头像错误]", err, avatar)
		ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
		return
	}

	img := imaging.Resize(avatar, size, size, imaging.Lanczos)
	imageData, imgErr := r.encodeImage(img, imageExt)
	if imgErr != nil {
		facades.Log.Error("WeAvatar[生成头像错误]", imgErr.Error())
		ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
		return
	}

	ctx.Response().Header("Cache-Control", "public, max-age=300")
	ctx.Response().Header("Avatar-By", "weavatar.com")
	ctx.Response().Header("Avatar-From", from)
	ctx.Response().Header("Vary", "Accept")

	ctx.Response().Data(http.StatusOK, "image/"+imageExt, imageData)
}

// encodeImage 编码图片为指定格式
func (r *AvatarController) encodeImage(img image.Image, imageExt string) ([]byte, error) {
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
func (r *AvatarController) Index(ctx http.Context) {
	// 取出用户信息
	user, ok := ctx.Value("user").(models.User)
	if !ok {
		ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
			"code":    401,
			"message": "登录已过期",
		})
		return
	}

	var avatars []models.Avatar
	err := facades.Orm.Query().Where("user_id", user.ID).Find(&avatars)
	if err != nil {
		facades.Log.WithContext(ctx).Error("[AvatarController][Index] 查询用户头像失败: ", err.Error())
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "系统内部错误",
		})
		return
	}
}

// Store 添加头像
func (r *AvatarController) Store(ctx http.Context) {

}

// Update 更新头像
func (r *AvatarController) Update(ctx http.Context) {

}

// Destroy 删除头像
func (r *AvatarController) Destroy(ctx http.Context) {

}

// CheckBind 检查绑定
func (r *AvatarController) CheckBind(ctx http.Context) {

}
