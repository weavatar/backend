package controllers

import (
	"bytes"
	"github.com/goravel/framework/facades"
	"image"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/goravel/framework/contracts/http"
	"github.com/nickalie/go-webpbin"

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
	var option string
	var err error
	from := "weavatar"

	if defaultAvatar == "letter" {
		option = strings.Trim(ctx.Request().Input("letter"), " ")
	} else {
		option = hash
	}

	if forceDefault {
		avatar, err = avatarService.GetDefaultAvatarByType(defaultAvatar, option)
	} else {
		avatar, from, err = avatarService.GetAvatar(appid, hash, defaultAvatar, option)
	}

	if err != nil {
		facades.Log.Error("WeAvatar[获取头像错误]", err.Error())
		ctx.Response().String(http.StatusInternalServerError, "WeAvatar 服务出现错误")
		return
	}

	// 判断一下 404 请求
	if avatar == nil && defaultAvatar == "404" {
		ctx.Response().String(http.StatusNotFound, "404 Not Found\nWeAvatar")
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

// encodeImage 编码图片
func (r *AvatarController) encodeImage(img image.Image, imageExt string) ([]byte, error) {
	var err error
	writer := bytes.NewBuffer([]byte{})

	switch imageExt {
	case "webp":
		err = webpbin.Encode(writer, img)
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

func (r *AvatarController) Create(ctx http.Context) {

}

func (r *AvatarController) Get(ctx http.Context) {

}

func (r *AvatarController) GetSingle(ctx http.Context) {

}

func (r *AvatarController) Update(ctx http.Context) {

}

func (r *AvatarController) Delete(ctx http.Context) {

}

func (r *AvatarController) CheckBind(ctx http.Context) {

}
