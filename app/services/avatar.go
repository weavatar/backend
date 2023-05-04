package services

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/HaoZi-Team/letteravatar"
	"github.com/disintegration/imaging"
	"github.com/golang-module/carbon/v2"
	"github.com/golang/freetype/truetype"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"
	"github.com/ipsn/go-adorable"
	"github.com/issue9/identicon/v2"
	"github.com/o1egl/govatar"
	"golang.org/x/exp/slices"
	_ "golang.org/x/image/webp"

	"weavatar/app/jobs"
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
func (r *AvatarImpl) Sanitize(ctx http.Context) (uint, string, string, int, bool, string) {
	// 从 URL 中获取 Hash 和图片格式
	hashExt := strings.Split(ctx.Request().Input("hash", ""), ".")
	// 从查询字符串中获取 APPID
	appid := ctx.Request().InputInt("appid", 0)

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
	size, err := strconv.Atoi(ctx.Request().Input("s", "80"))
	if err != nil {
		size = 80
	}
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
		parsedURL, parseErr := url.Parse(defaultAvatar)
		if parseErr != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || len(parsedURL.Host) == 0 || strings.Contains(parsedURL.Host, "weavatar.com") || len(parsedURL.RawQuery) != 0 {
			defaultAvatar = ""
		}
	}

	// 判断 Hash 是否32位
	if len(hash) != 32 {
		forceDefaultBool = true
	}

	return uint(appid), hash, imageExt, size, forceDefaultBool, defaultAvatar
}

// getQqAvatar 通过 QQ 号获取头像
func (r *AvatarImpl) getQqAvatar(hash string) (image.Image, error) {
	hashIndex, hashErr := strconv.ParseInt(hash[:10], 16, 64)
	if hashErr != nil {
		return nil, hashErr
	}
	tableIndex := (hashIndex % int64(4000)) + 1
	type qqHash struct {
		Hash string `gorm:"column:hash"`
		Qq   string `gorm:"column:qq"`
	}
	var qq qqHash

	if facades.Storage.Exists("cache/qq/" + hash[:2] + "/" + hash) {
		img, imgErr := imaging.Open(facades.Storage.Path("cache/qq/" + hash[:2] + "/" + hash))
		if imgErr != nil {
			facades.Log.Warning("QQ头像[图片解析出错] ", imgErr.Error())
			return nil, imgErr
		}

		return img, nil
	} else {
		err := facades.Orm.Connection("hash").Query().Table("qq_"+strconv.Itoa(int(tableIndex))).Where("hash", hash).First(&qq)
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
		if !resp.IsSuccessState() {
			if reqErr != nil {
				return nil, reqErr
			} else {
				return nil, errors.New("获取 QQ头像 失败")
			}
		}

		length, lengthErr := strconv.Atoi(resp.GetHeader("Content-Length"))
		if length < 6400 || lengthErr != nil {
			// 如果图片小于 6400 字节，则尝试获取 100 尺寸的图片
			resp, reqErr = client.R().SetQueryParams(map[string]string{
				"b":  "qq",
				"nk": qq.Qq,
				"s":  "100",
			}).Get("http://q1.qlogo.cn/g")
			if !resp.IsSuccessState() {
				if reqErr != nil {
					return nil, reqErr
				} else {
					return nil, errors.New("获取 QQ头像 失败")
				}
			}
		}

		// 检查图片是否正常
		reader := strings.NewReader(resp.String())
		img, imgErr := imaging.Decode(reader)
		if err != nil {
			facades.Log.Warning("QQ头像[图片不正常] ", err.Error())
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
			facades.Log.Warning("Gravatar[图片解析出错] ", imgErr.Error())
			return nil, imgErr
		}
	} else {
		// 不存在则从 Gravatar 获取头像
		client := req.C().SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.34")
		resp, reqErr := client.R().Get("http://proxy.server/http://0.gravatar.com/avatar/" + hash + ".png?s=1000&r=g&d=404")
		if reqErr != nil || !resp.IsSuccessState() {
			return nil, errors.New("获取 Gravatar头像 失败")
		}

		// 检查图片是否正常
		reader := strings.NewReader(resp.String())
		img, imgErr = imaging.Decode(reader)
		if imgErr != nil {
			facades.Log.Warning("Gravatar[图片不正常] ", imgErr.Error())
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
func (r *AvatarImpl) getDefaultAvatar(defaultAvatar string, option []string) (image.Image, error) {

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

	if defaultAvatar == "mp" {
		img, imgErr := imaging.Open(facades.Storage.Path("default/mp.png"))
		if imgErr != nil {
			return nil, imgErr
		}

		return img, nil
	}

	if defaultAvatar == "identicon" {
		img := identicon.Make(identicon.Style1, 1200, color.RGBA{R: 255, A: 100}, color.RGBA{R: 102, G: 204, B: 255, A: 255}, []byte(option[0]))
		return img, nil
	}

	if defaultAvatar == "monsterid" {
		img, err := govatar.GenerateForUsername(govatar.FEMALE, option[0])
		if err != nil {
			return nil, err
		}

		return img, nil
	}

	if defaultAvatar == "wavatar" {
		avatar := adorable.PseudoRandom([]byte(option[0]))
		img, err := imaging.Decode(bytes.NewReader(avatar))

		if err != nil {
			return nil, err
		}

		return img, nil
	}

	if defaultAvatar == "retro" {
		ii := identicon.New(identicon.Style2, 1200, color.RGBA{R: 255, A: 100}, color.RGBA{R: 102, G: 204, B: 255, A: 255})
		img := ii.Make([]byte(option[0]))
		return img, nil
	}

	if defaultAvatar == "robohash" {
		img, err := govatar.GenerateForUsername(govatar.MALE, option[0])
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
		fontStr, err := facades.Storage.Get("fonts/HarmonyOS_Sans_SC_Medium.ttf")
		if err != nil {
			return nil, err
		}
		font, fontErr := truetype.Parse([]byte(fontStr))
		if fontErr != nil {
			return nil, fontErr
		}

		fontSize := 0

		// 判断中文
		hasChinese := false
		for _, w := range option[0] {
			if unicode.Is(unicode.Han, w) {
				hasChinese = true
				break
			}
		}

		// 判断长度
		letters := []rune(option[0])
		length := len(letters)
		if length > 4 {
			letters = letters[:4]
			length = 4
		}
		switch {
		case length == 1:
			if hasChinese {
				fontSize = 500
			} else {
				fontSize = 600
			}
		case length == 2:
			if hasChinese {
				fontSize = 400
			} else {
				fontSize = 500
			}
		case length == 3:
			if hasChinese {
				fontSize = 300
			} else {
				fontSize = 400
			}
		case length == 4:
			if hasChinese {
				fontSize = 200
			} else {
				fontSize = 300
			}
		default:
			fontSize = 200
		}

		img, imgErr := letteravatar.Draw(1000, letters, &letteravatar.Options{
			Font:       font,
			FontSize:   fontSize,
			PaletteKey: option[1], // 对相同的字符串使用相同的颜色
		})
		if imgErr != nil {
			return nil, err
		}

		return img, nil
	}

	client := req.C().SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.34")
	resp, reqErr := client.R().Get(defaultAvatar)
	if reqErr != nil || !resp.IsSuccessState() {
		img, imgErr := imaging.Open(facades.Storage.Path("default/default.png"))
		if imgErr != nil {
			return nil, imgErr
		}

		return img, nil
	}

	// 检查图片是否正常
	reader := strings.NewReader(resp.String())
	img, imgErr := imaging.Decode(reader)
	if imgErr != nil {
		img, imgErr = imaging.Open(facades.Storage.Path("default/default.png"))
		if imgErr != nil {
			return nil, imgErr
		}
	}

	return img, nil
}

// GetDefaultAvatarByType 通过默认头像类型获取头像
func (r *AvatarImpl) GetDefaultAvatarByType(avatarType string, option []string) (image.Image, error) {
	var avatar image.Image
	var err error

	switch avatarType {
	case "404":
		avatar, err = r.getDefaultAvatar("404", option)
	case "mp", "mm", "mystery":
		avatar, err = r.getDefaultAvatar("mp", option)
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
		avatar, err = r.getDefaultAvatar(avatarType, option)
	}

	return avatar, err
}

// GetAvatar 获取头像
func (r *AvatarImpl) GetAvatar(appid uint, hash string, defaultAvatar string, option []string) (image.Image, string, error) {
	var avatar models.Avatar
	var appAvatar models.AppAvatar
	var err error

	banImg, banImgErr := imaging.Open(facades.Storage.Path("default/ban.png"))
	var img image.Image
	from := "weavatar"
	var imgErr error
	if banImgErr != nil {
		facades.Log.Warning("WeAvatar[封禁图片解析出错] ", banImgErr.Error())
		return nil, from, banImgErr
	}

	// 取头像数据
	_, err = facades.Orm.Query().Exec(`INSERT INTO "avatars" ("hash", "created_at", "updated_at") VALUES (?, ?, ?) ON CONFLICT DO NOTHING`, hash, carbon.DateTime{Carbon: carbon.Now()}, carbon.DateTime{Carbon: carbon.Now()})
	if err != nil {
		facades.Log.Error("WeAvatar[数据库错误] ", err.Error())
		return nil, from, err
	}
	err = facades.Orm.Query().Where("hash", hash).First(&avatar)
	if err != nil {
		facades.Log.Error("WeAvatar[数据库错误] ", err.Error())
		return nil, from, err
	}
	if avatar.UserID != nil && avatar.Raw != nil {
		// 检查 Hash 是否有对应的 App
		err = facades.Orm.Query().Where("app_id", appid).Where("avatar_hash", hash).First(&appAvatar)
		if err != nil {
			facades.Log.Error("WeAvatar[数据库错误] ", err.Error())
			return nil, from, err
		}
		// 如果有对应的 App，则检查其 APP 头像是否封禁状态
		if appAvatar.AppID != 0 {
			if appAvatar.Ban {
				return banImg, from, nil
			} else {
				// 如果不是封禁状态，则检查是否有对应的头像
				img, imgErr = imaging.Open(facades.Storage.Path("upload/app/" + strconv.Itoa(int(appAvatar.AppID)) + "/" + hash[:2] + "/" + hash))
			}
		} else {
			// 检查默认头像是否封禁状态
			if avatar.Ban {
				return banImg, from, nil
			} else {
				// 如果不是封禁状态，则检查是否有对应的头像
				img, imgErr = imaging.Open(facades.Storage.Path("upload/default/" + hash[:2] + "/" + hash))
			}
		}

		// 如果头像获取失败，则使用默认头像
		if imgErr != nil {
			facades.Log.Warning("WeAvatar[头像匹配失败] ", imgErr.Error())
			img, _ = r.GetDefaultAvatarByType(defaultAvatar, option)
			return img, from, nil
		}
	} else {
		// 检查头像是否封禁状态
		if avatar.Ban {
			return banImg, from, nil
		}
		// 优先使用 Gravatar 头像
		img, imgErr = r.getGravatarAvatar(hash)
		from = "gravatar"
		if imgErr != nil {
			// 如果 Gravatar 头像获取失败，则使用 QQ 头像
			img, imgErr = r.getQqAvatar(hash)
			from = "qq"
			if imgErr != nil {
				// 如果 QQ 头像获取失败，则使用默认头像
				img, _ = r.GetDefaultAvatarByType(defaultAvatar, option)
				from = "weavatar"
			}
		}
	}

	// 审核头像
	if !avatar.Checked && from != "weavatar" {
		go func() {
			_ = facades.Queue.Job(&jobs.ProcessAvatarCheck{}, []queue.Arg{
				{Type: "string", Value: hash},
			}).Dispatch()
		}()
	}

	return img, from, nil
}
