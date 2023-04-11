package controllers

import (
	"bytes"
	"image"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/goravel/framework/contracts/http"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"

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
	// WEBP 参数
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetIcon, 85)
	if err != nil {
		return nil, err
	}
	writer := bytes.NewBuffer([]byte{})

	switch imageExt {
	case "webp":
		err = webp.Encode(writer, img, options)
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
