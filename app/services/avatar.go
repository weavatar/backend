package services

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/disintegration/imaging"
	"github.com/disintegration/letteravatar"
	"github.com/golang/freetype/truetype"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"
	"github.com/ipsn/go-adorable"
	"github.com/issue9/identicon/v2"
	"github.com/o1egl/govatar"
	"github.com/spf13/cast"
	"golang.org/x/exp/slices"
	_ "golang.org/x/image/webp"

	"weavatar/app/models"
)

type Avatar interface {
}

type AvatarImpl struct {
}

func NewAvatarImpl() *AvatarImpl {
	return &AvatarImpl{}
}

// Sanitize 消毒头像请求
func (r *AvatarImpl) Sanitize(ctx http.Context) (string, string, string, int, bool, string) {
	// 从 URL 中获取 Hash 和图片格式
	hashExt := strings.Split(ctx.Request().Input("hash", ""), ".")
	// 从查询字符串中获取 APPID
	appid := ctx.Request().Input("appid", "")

	hash := strings.ToLower(hashExt[0]) // Hash 转小写
	imageExt := "png"                   // 默认为 PNG 格式

	// 如果浏览器支持 WEBP 格式，则默认使用 WEBP 格式
	accept := ctx.Request().Header("Accept", "")
	if strings.Contains(accept, "image/webp") || strings.Contains(accept, "image/*") {
		imageExt = "webp"
	}
	// 如果提供了图片格式，则使用提供的图片格式
	if len(hashExt) > 1 {
		imageExt = hashExt[1]
	}
	// 检查图片格式是否支持
	imageSlices := []string{"png", "jpg", "jpeg", "gif", "webp"}
	if !slices.Contains(imageSlices, imageExt) {
		imageExt = "png"
	}

	// 从 URL 中获取图片大小
	size := cast.ToInt(ctx.Request().Input("s", "80"))
	// 检查图片大小是否合法
	if size > 2000 {
		size = 2000
	}
	if size < 10 {
		size = 10
	}

	// 从 URL 中获取是否强制使用默认头像
	forceDefault := ctx.Request().Input("f", "n")
	forceDefaultBool := false
	forceDefaultSlices := []string{"y", "yes"}
	if slices.Contains(forceDefaultSlices, forceDefault) {
		forceDefaultBool = true
	}

	// 从 URL 中获取默认头像
	defaultAvatar := ctx.Request().Input("d", "")
	defaultAvatarSlices := []string{"404", "mp", "mm", "mystery", "identicon", "monsterid", "wavatar", "retro", "robohash", "blank", "letter"}
	if !slices.Contains(defaultAvatarSlices, defaultAvatar) {
		// 如果不是预设的默认头像，则检查是否是合法的 URL
		parsedURL, err := url.Parse(forceDefault)
		if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
			defaultAvatar = ""
		}
	}

	// 判断 Hash 是否32位
	if len(hash) != 32 {
		forceDefaultBool = true
	}

	return appid, hash, imageExt, size, forceDefaultBool, defaultAvatar
}

// getQqAvatar 通过 QQ 号获取头像
func (r *AvatarImpl) getQqAvatar(hash string) (image.Image, error) {
	tableIndex := (cast.ToInt64("0x"+hash[:10]) % int64(4000)) + 1
	type qqHash struct {
		Hash string `gorm:"column:hash"`
		Qq   string `gorm:"column:qq"`
	}
	var qq qqHash

	if facades.Storage.Exists("cache/qq/" + hash[:2] + "/" + hash) {
		img, imgErr := imaging.Open(facades.Storage.Path("cache/qq/" + hash[:2] + "/" + hash))
		if imgErr != nil {
			facades.Log.Warning("QQ头像[图片解析出错]", imgErr.Error())
			return nil, imgErr
		}

		return img, nil
	} else {
		err := facades.Orm.Connection("hash").Query().Table("qq_"+cast.ToString(tableIndex)).Where("hash", hash).First(&qq)
		if err != nil {
			return nil, err
		}

		if qq.Qq == "" {
			return nil, errors.New("未找到对应的 QQ 号")
		}

		client := req.C().SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.34")
		resp, reqErr := client.R().SetQueryParams(map[string]string{
			"b":  "qq",
			"nk": qq.Qq,
			"s":  "640",
		}).Get("http://q1.qlogo.cn/g")
		if reqErr != nil || !resp.IsSuccessState() {
			facades.Log.Warning("QQ头像[获取图片出错]", reqErr.Error())
			return nil, errors.New("获取 QQ头像 失败")
		}

		if cast.ToInt(resp.GetHeader("Content-Length")) < 6400 {
			// 如果图片小于 6400 字节，则尝试获取 100 尺寸的图片
			resp, reqErr = client.R().SetQueryParams(map[string]string{
				"b":  "qq",
				"nk": qq.Qq,
				"s":  "100",
			}).Get("http://q1.qlogo.cn/g")
			if reqErr != nil || !resp.IsSuccessState() {
				facades.Log.Warning("QQ头像[图片解析出错]", reqErr.Error())
				return nil, reqErr
			}
		}

		// 检查图片是否正常
		reader := strings.NewReader(resp.String())
		img, imgErr := imaging.Decode(reader)
		if err != nil {
			facades.Log.Warning("QQ头像[图片不正常]", err.Error())
			return nil, imgErr
		}

		// 保存图片
		err = facades.Storage.Put("cache/qq/"+hash[:2]+"/"+hash, resp.String())
		if err != nil {
			return nil, err
		}

		return img, nil
	}
}

// getGravatarAvatar 通过 Gravatar 获取头像
func (r *AvatarImpl) getGravatarAvatar(hash string) (image.Image, error) {
	var img image.Image
	var imgErr error

	if facades.Storage.Exists("cache/gravatar/" + hash[:2] + "/" + hash) {
		img, imgErr = imaging.Open(facades.Storage.Path("cache/gravatar/" + hash[:2] + "/" + hash))
		if imgErr != nil {
			facades.Log.Warning("Gravatar[图片解析出错]", imgErr.Error())
			return nil, imgErr
		}
	} else {
		// 不存在则从 Gravatar 获取头像
		client := req.C().SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.34")
		resp, reqErr := client.R().Get("http://proxy.server/http://0.gravatar.com/avatar/" + hash + ".png?s=1000&r=g&d=404")
		if reqErr != nil || !resp.IsSuccessState() {
			// Gravatar 不需要记录日志
			return nil, errors.New("获取 Gravatar头像 失败")
		}

		// 检查图片是否正常
		reader := strings.NewReader(resp.String())
		img, imgErr = imaging.Decode(reader)
		if imgErr != nil {
			facades.Log.Warning("Gravatar[图片不正常]", imgErr.Error())
			return nil, imgErr
		}

		// 保存图片
		err := facades.Storage.Put("cache/gravatar/"+hash[:2]+"/"+hash, resp.String())
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

// getDefaultAvatar 通过默认参数获取头像
func (r *AvatarImpl) getDefaultAvatar(defaultAvatar string, option string) (image.Image, error) {

	if defaultAvatar == "404" {
		return nil, nil
	}

	if defaultAvatar == "" {
		img, imgErr := imaging.Open(facades.Storage.Path("default/default.png"))
		if imgErr != nil {
			return nil, imgErr
		}

		return img, nil
	}

	if defaultAvatar == "mp" || defaultAvatar == "mm" || defaultAvatar == "mystery" {
		img, imgErr := imaging.Open(facades.Storage.Path("default/mp.png"))
		if imgErr != nil {
			return nil, imgErr
		}

		return img, nil
	}

	if defaultAvatar == "identicon" {
		img := identicon.Make(identicon.Style1, 1200, color.RGBA{R: 255, A: 100}, color.RGBA{R: 102, G: 204, B: 255, A: 255}, []byte(option))
		return img, nil
	}

	if defaultAvatar == "monsterid" {
		img, err := govatar.GenerateForUsername(govatar.FEMALE, option)
		if err != nil {
			return nil, err
		}

		return img, nil
	}

	if defaultAvatar == "wavatar" {
		avatar := adorable.PseudoRandom([]byte(option))
		img, err := imaging.Decode(bytes.NewReader(avatar))

		if err != nil {
			return nil, err
		}

		return img, nil
	}

	if defaultAvatar == "retro" {
		ii := identicon.New(identicon.Style2, 1200, color.RGBA{R: 255, A: 100}, color.RGBA{R: 102, G: 204, B: 255, A: 255})
		img := ii.Make([]byte(option))
		return img, nil
	}

	if defaultAvatar == "robohash" {
		img, err := govatar.GenerateForUsername(govatar.MALE, option)
		if err != nil {
			return nil, err
		}

		return img, nil
	}

	if defaultAvatar == "blank" {
		img, imgErr := imaging.Open(facades.Storage.Path("default/blank.png"))
		if imgErr != nil {
			return nil, imgErr
		}

		return img, nil
	}

	if defaultAvatar == "letter" {
		firstLetter, _ := utf8.DecodeRuneInString(option)
		fontStr, err := facades.Storage.Get("fonts/HarmonyOS_Sans_SC_Medium.ttf")
		if err != nil {
			return nil, err
		}
		font, fontErr := truetype.Parse([]byte(fontStr))
		if fontErr != nil {
			return nil, fontErr
		}

		img, imgErr := letteravatar.Draw(100, firstLetter, &letteravatar.Options{
			Font:       font,
			PaletteKey: option, // 对相同的字符串使用相同的颜色
		})
		if imgErr != nil {
			return nil, err
		}

		return img, nil
	}

	return nil, nil
}

func (r *AvatarImpl) GetDefaultAvatarByType(avatarType, option string) (image.Image, error) {
	var avatar image.Image
	var err error

	switch avatarType {
	case "404":
		avatar, err = r.getDefaultAvatar("404", option)
	case "mp":
		avatar, err = r.getDefaultAvatar("mp", option)
	case "mm":
		avatar, err = r.getDefaultAvatar("mm", option)
	case "mystery":
		avatar, err = r.getDefaultAvatar("mystery", option)
	case "identicon":
		avatar, err = r.getDefaultAvatar("identicon", option)
	case "monsterid":
		avatar, err = r.getDefaultAvatar("monsterid", option)
	case "wavatar":
		avatar, err = r.getDefaultAvatar("wavatar", option)
	case "retro":
		avatar, err = r.getDefaultAvatar("retro", option)
	case "robohash":
		avatar, err = r.getDefaultAvatar("robohash", option)
	case "blank":
		avatar, err = r.getDefaultAvatar("blank", option)
	case "letter":
		avatar, err = r.getDefaultAvatar("letter", option)
	default:
		avatar, err = r.getDefaultAvatar("", option)
	}

	return avatar, err
}

// GetAvatar 获取头像
func (r *AvatarImpl) GetAvatar(appid string, hash string, defaultAvatar string, option string) (image.Image, string, error) {
	var avatar models.Avatar
	var appAvatar models.AppAvatar
	var err error

	banImg, banImgErr := imaging.Open(facades.Storage.Path("default/ban.png"))
	var img image.Image
	var imgErr error
	if banImgErr != nil {
		facades.Log.Warning("WeAvatar[封禁图片解析出错]", banImgErr.Error())
		return nil, "weavatar", banImgErr
	}

	// 检查是否有默认头像
	err = facades.Orm.Query().Where("hash", hash).First(&avatar)
	if err != nil {
		facades.Log.Error("WeAvatar[数据库错误]", err.Error())
		return nil, "weavatar", err
	}
	if avatar.UserID != 0 && avatar.Hash != "" {
		// 检查 Hash 是否有对应的 App
		err = facades.Orm.Query().Where("app_id", appid).Where("avatar_hash", hash).First(&appAvatar)
		if err != nil {
			facades.Log.Error("WeAvatar[数据库错误]", err.Error())
			return nil, "weavatar", err
		}
		// 如果有对应的 App，则检查其 APP 头像是否封禁状态
		if appAvatar.AppID != 0 {
			if appAvatar.Ban {
				return banImg, "weavatar", nil
			} else {
				// 如果不是封禁状态，则检查是否有对应的头像
				img, imgErr = imaging.Open(facades.Storage.Path("upload/app/" + cast.ToString(appAvatar.AppID) + "/" + hash[:2] + "/" + hash))
			}
		} else {
			// 检查默认头像是否封禁状态
			if avatar.Ban {
				return banImg, "weavatar", nil
			} else {
				// 如果不是封禁状态，则检查是否有对应的头像
				img, imgErr = imaging.Open(facades.Storage.Path("upload/default/" + cast.ToString(appAvatar.AppID) + "/" + hash[:2] + "/" + hash))
			}
		}

		// 如果头像获取失败，则使用默认头像
		if imgErr != nil {
			facades.Log.Warning("WeAvatar[头像匹配失败]", imgErr.Error())
			img, _ = r.GetDefaultAvatarByType(defaultAvatar, option)
			return img, "weavatar", nil
		}
		return img, "weavatar", nil
	} else {
		// 优先使用 Gravatar 头像
		img, imgErr = r.getGravatarAvatar(hash)
		from := "gravatar"
		if imgErr != nil {
			// 如果 Gravatar 头像获取失败，则使用 QQ 头像
			img, imgErr = r.getQqAvatar(hash)
			from = "qq"
			if imgErr != nil {
				// 如果 QQ 头像获取失败，则使用默认头像
				from = "weavatar"
				img, _ = r.GetDefaultAvatarByType(defaultAvatar, option)
			}
		}

		return img, from, nil
	}
}
