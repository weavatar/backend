package services

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/HaoZi-Team/letteravatar"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
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
	Sanitize(ctx http.Context) (string, string, string, int, bool, string)
	GetQqAvatar(hash string) (image.Image, carbon.Carbon, error)
	GetGravatarAvatar(hash string) (image.Image, carbon.Carbon, error)
	GetDefaultAvatar(defaultAvatar string, option []string) (image.Image, carbon.Carbon, error)
	GetDefaultAvatarByType(avatarType string, option []string) (image.Image, carbon.Carbon, error)
	GetAvatar(appid string, hash string, defaultAvatar string, option []string) (image.Image, carbon.Carbon, string, error)
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
	appid := ctx.Request().Input("appid", "0")

	hash := strings.ToLower(hashExt[0]) // Hash 转小写
	imageExt := "webp"                  // 默认为 WEBP 格式

	// 如果提供了图片格式，则使用提供的图片格式
	if len(hashExt) > 1 {
		imageExt = hashExt[1]
	}
	// 检查图片格式是否支持
	imageSlices := []string{"png", "jpg", "jpeg", "gif", "webp"}
	if !slices.Contains(imageSlices, imageExt) {
		imageExt = "webp"
	}

	// 从 URL 中获取图片大小
	sizeStr := ctx.Request().Input("s", "")
	if len(sizeStr) == 0 {
		sizeStr = ctx.Request().Input("size", "80")
	}
	size, err := strconv.Atoi(sizeStr)
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
	forceDefault := ctx.Request().Input("f", "")
	if len(forceDefault) == 0 {
		forceDefault = ctx.Request().Input("forcedefault", "n")
	}
	forceDefaultBool := false
	forceDefaultSlices := []string{"y", "yes"}
	if slices.Contains(forceDefaultSlices, forceDefault) {
		forceDefaultBool = true
	}

	// 从 URL 中获取默认头像
	defaultAvatar := ctx.Request().Input("d", "")
	if len(defaultAvatar) == 0 {
		defaultAvatar = ctx.Request().Input("default", "")
	}
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

	return appid, hash, imageExt, size, forceDefaultBool, defaultAvatar
}

// GetQqAvatar 通过 QQ 号获取头像
func (r *AvatarImpl) GetQqAvatar(hash string) (image.Image, carbon.Carbon, error) {
	type qqHash struct {
		Hash string `gorm:"primaryKey"`
		Qq   uint   `gorm:"type:bigint;not null"`
	}
	var qq qqHash

	if facades.Storage().Exists("cache/qq/" + hash[:2] + "/" + hash) {
		gmt, err := time.LoadLocation("GMT")
		if err != nil {
			return nil, carbon.Now(), err
		}
		img, imgErr := imaging.Open(facades.Storage().Path("cache/qq/" + hash[:2] + "/" + hash))
		lastModified, err := facades.Storage().LastModified("cache/qq/" + hash[:2] + "/" + hash)
		lastModified = lastModified.In(gmt)
		if imgErr == nil && err == nil {
			return img, carbon.FromStdTime(lastModified), nil
		}
	}

	err := facades.Orm().Connection("hash").Query().Table("qq_mails").Where("hash", hash).First(&qq)
	if err != nil {
		return nil, carbon.Now(), err
	}

	if qq.Qq == 0 {
		return nil, carbon.Now(), errors.New("未找到对应的 QQ 号")
	}

	client := req.C()
	resp, reqErr := client.R().SetQueryParams(map[string]string{
		"b":  "qq",
		"nk": strconv.Itoa(int(qq.Qq)),
		"s":  "640",
	}).Get("http://q1.qlogo.cn/g")
	if !resp.IsSuccessState() {
		if reqErr != nil {
			return nil, carbon.Now(), reqErr
		} else {
			return nil, carbon.Now(), errors.New("获取 QQ头像 失败")
		}
	}

	length, lengthErr := strconv.Atoi(resp.GetHeader("Content-Length"))
	if length < 6400 || lengthErr != nil {
		// 如果图片小于 6400 字节，则尝试获取 100 尺寸的图片
		resp, reqErr = client.R().SetQueryParams(map[string]string{
			"b":  "qq",
			"nk": strconv.Itoa(int(qq.Qq)),
			"s":  "100",
		}).Get("http://q1.qlogo.cn/g")
		if !resp.IsSuccessState() {
			if reqErr != nil {
				return nil, carbon.Now(), reqErr
			} else {
				return nil, carbon.Now(), errors.New("获取 QQ头像 失败")
			}
		}
	}

	// 检查图片是否正常
	reader := bytes.NewReader(resp.Bytes())
	img, imgErr := imaging.Decode(reader)
	if err != nil {
		facades.Log().Warning("QQ头像[图片不正常] ", err.Error())
		return nil, carbon.Now(), imgErr
	}

	// 保存图片
	err = facades.Storage().Put("cache/qq/"+hash[:2]+"/"+hash, resp.String())
	if err != nil {
		return nil, carbon.Now(), err
	}
	gmt, err := time.LoadLocation("GMT")
	if err != nil {
		return nil, carbon.Now(), err
	}
	lastModified, err := facades.Storage().LastModified("cache/qq/" + hash[:2] + "/" + hash)
	if err != nil {
		return nil, carbon.Now(), err
	}
	lastModified = lastModified.In(gmt)

	return img, carbon.FromStdTime(lastModified), nil
}

// GetGravatarAvatar 通过 Gravatar 获取头像
func (r *AvatarImpl) GetGravatarAvatar(hash string) (image.Image, carbon.Carbon, error) {
	var img image.Image
	var imgErr error

	if facades.Storage().Exists("cache/gravatar/" + hash[:2] + "/" + hash) {
		gmt, err := time.LoadLocation("GMT")
		if err != nil {
			return nil, carbon.Now(), err
		}
		img, imgErr = imaging.Open(facades.Storage().Path("cache/gravatar/" + hash[:2] + "/" + hash))
		lastModified, err := facades.Storage().LastModified("cache/gravatar/" + hash[:2] + "/" + hash)
		lastModified = lastModified.In(gmt)
		if imgErr == nil && err == nil {
			return img, carbon.FromStdTime(lastModified), nil
		}
	}

	client := req.C()
	// 有一些头像请求1000尺寸的大图(http://0.gravatar.com/avatar/1b6a1437577086c55c785980430123ce.png?s=1000&r=g&d=404), Gravatar会返回404, 不知道为什么，所以使用600尺寸的图片代替
	resp, reqErr := client.R().Get("http://proxy.server/http://0.gravatar.com/avatar/" + hash + ".png?s=600&r=g&d=404")
	if reqErr != nil || !resp.IsSuccessState() {
		return nil, carbon.Now(), errors.New("获取 Gravatar头像 失败")
	}

	// 检查图片是否正常
	reader := bytes.NewReader(resp.Bytes())
	img, imgErr = imaging.Decode(reader)
	if imgErr != nil {
		facades.Log().Warning("Gravatar[图片不正常] ", imgErr.Error())
		return nil, carbon.Now(), imgErr
	}

	// 保存图片
	err := facades.Storage().Put("cache/gravatar/"+hash[:2]+"/"+hash, resp.String())
	if err != nil {
		return nil, carbon.Now(), err
	}
	gmt, err := time.LoadLocation("GMT")
	if err != nil {
		return nil, carbon.Now(), err
	}
	lastModified, err := facades.Storage().LastModified("cache/gravatar/" + hash[:2] + "/" + hash)
	lastModified = lastModified.In(gmt)
	if err != nil {
		return nil, carbon.Now(), err
	}

	return img, carbon.FromStdTime(lastModified), nil
}

// GetDefaultAvatar 通过默认参数获取头像
func (r *AvatarImpl) GetDefaultAvatar(defaultAvatar string, option []string) (image.Image, carbon.Carbon, error) {
	if defaultAvatar == "404" {
		return nil, carbon.Now(), nil
	}

	if defaultAvatar == "" {
		img, imgErr := imaging.Open(facades.Storage().Path("default/default.png"))
		if imgErr != nil {
			return nil, carbon.Now(), imgErr
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "mp" {
		img, imgErr := imaging.Open(facades.Storage().Path("default/mp.png"))
		if imgErr != nil {
			return nil, carbon.Now(), imgErr
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "identicon" {
		img := identicon.Make(identicon.Style1, 1200, color.RGBA{R: 255, A: 100}, color.RGBA{R: 102, G: 204, B: 255, A: 255}, []byte(option[0]))
		return img, carbon.Now(), nil
	}

	if defaultAvatar == "monsterid" {
		img, err := govatar.GenerateForUsername(govatar.FEMALE, option[0])
		if err != nil {
			return nil, carbon.Now(), err
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "wavatar" {
		avatar := adorable.PseudoRandom([]byte(option[0]))
		img, err := imaging.Decode(bytes.NewReader(avatar))

		if err != nil {
			return nil, carbon.Now(), err
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "retro" {
		ii := identicon.New(identicon.Style2, 1200, color.RGBA{R: 255, A: 100}, color.RGBA{R: 102, G: 204, B: 255, A: 255})
		img := ii.Make([]byte(option[0]))
		return img, carbon.Now(), nil
	}

	if defaultAvatar == "robohash" {
		img, err := govatar.GenerateForUsername(govatar.MALE, option[0])
		if err != nil {
			return nil, carbon.Now(), err
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "blank" {
		img, imgErr := imaging.Open(facades.Storage().Path("default/blank.png"))
		if imgErr != nil {
			return nil, carbon.Now(), imgErr
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "letter" {
		fontStr, err := facades.Storage().Get("fonts/HarmonyOS_Sans_SC_Medium.ttf")
		if err != nil {
			return nil, carbon.Now(), err
		}
		font, fontErr := truetype.Parse([]byte(fontStr))
		if fontErr != nil {
			return nil, carbon.Now(), fontErr
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
			return nil, carbon.Now(), err
		}

		return img, carbon.Now(), nil
	}

	client := req.C()
	resp, reqErr := client.R().Get(defaultAvatar)
	if reqErr != nil || !resp.IsSuccessState() {
		img, imgErr := imaging.Open(facades.Storage().Path("default/default.png"))
		if imgErr != nil {
			return nil, carbon.Now(), imgErr
		}

		return img, carbon.Now(), nil
	}

	// 检查图片是否正常
	reader := bytes.NewReader(resp.Bytes())
	img, imgErr := imaging.Decode(reader)
	if imgErr != nil {
		img, imgErr = imaging.Open(facades.Storage().Path("default/default.png"))
		if imgErr != nil {
			return nil, carbon.Now(), imgErr
		}
	}

	return img, carbon.Now(), nil
}

// GetDefaultAvatarByType 通过默认头像类型获取头像
func (r *AvatarImpl) GetDefaultAvatarByType(avatarType string, option []string) (image.Image, carbon.Carbon, error) {
	var avatar image.Image
	var lastModified carbon.Carbon
	var err error

	switch avatarType {
	case "404":
		avatar, lastModified, err = r.GetDefaultAvatar("404", option)
	case "mp", "mm", "mystery":
		avatar, lastModified, err = r.GetDefaultAvatar("mp", option)
	case "identicon":
		avatar, lastModified, err = r.GetDefaultAvatar("identicon", option)
	case "monsterid":
		avatar, lastModified, err = r.GetDefaultAvatar("monsterid", option)
	case "wavatar":
		avatar, lastModified, err = r.GetDefaultAvatar("wavatar", option)
	case "retro":
		avatar, lastModified, err = r.GetDefaultAvatar("retro", option)
	case "robohash":
		avatar, lastModified, err = r.GetDefaultAvatar("robohash", option)
	case "blank":
		avatar, lastModified, err = r.GetDefaultAvatar("blank", option)
	case "letter":
		avatar, lastModified, err = r.GetDefaultAvatar("letter", option)
	default:
		avatar, lastModified, err = r.GetDefaultAvatar(avatarType, option)
	}

	return avatar, lastModified, err
}

// GetAvatar 获取头像
func (r *AvatarImpl) GetAvatar(appid string, hash string, defaultAvatar string, option []string) (image.Image, carbon.Carbon, string, error) {
	var avatar models.Avatar
	var appAvatar models.AppAvatar
	var err error

	var img image.Image
	var lastModified carbon.Carbon

	banImg, banImgErr := imaging.Open(facades.Storage().Path("default/ban.png"))
	from := "weavatar"
	var imgErr error
	if banImgErr != nil {
		facades.Log().Warning("WeAvatar[封禁图片解析出错] ", banImgErr.Error())
		return nil, carbon.Now(), from, banImgErr
	}

	// 取头像数据
	_, err = facades.Orm().Query().Exec(`INSERT IGNORE INTO avatars (hash, created_at, updated_at) VALUES (?, ?, ?)`, hash, carbon.DateTime{Carbon: carbon.Now()}, carbon.DateTime{Carbon: carbon.Now()})
	if err != nil {
		facades.Log().Error("WeAvatar[数据库错误] ", err.Error())
		return nil, carbon.Now(), from, err
	}
	err = facades.Orm().Query().Where("hash", hash).First(&avatar)
	if err != nil {
		facades.Log().Error("WeAvatar[数据库错误] ", err.Error())
		return nil, carbon.Now(), from, err
	}
	if avatar.UserID != nil && avatar.Raw != nil {
		// 检查 Hash 是否有对应的 App
		err = facades.Orm().Query().Where("app_id", appid).Where("avatar_hash", hash).First(&appAvatar)
		if err != nil {
			facades.Log().Error("WeAvatar[数据库错误] ", err.Error())
			return nil, carbon.Now(), from, err
		}
		// 如果有对应的 App，则检查其 APP 头像是否封禁状态
		if appAvatar.AppID != 0 {
			if appAvatar.Ban {
				return banImg, carbon.Now(), from, nil
			} else {
				// 如果不是封禁状态，则检查是否有对应的头像
				img, imgErr = imaging.Open(facades.Storage().Path("upload/app/" + strconv.Itoa(int(appAvatar.AppID)) + "/" + hash[:2] + "/" + hash))
			}
		} else {
			if avatar.Ban {
				return banImg, carbon.Now(), from, nil
			} else {
				// 如果不是封禁状态，则检查是否有对应的头像
				img, imgErr = imaging.Open(facades.Storage().Path("upload/default/" + hash[:2] + "/" + hash))
			}
		}

		// 如果头像获取失败，则使用默认头像
		if imgErr != nil {
			facades.Log().Warning("WeAvatar[头像匹配失败] ", imgErr.Error())
			img, lastModified, _ = r.GetDefaultAvatarByType(defaultAvatar, option)
			return img, lastModified, from, nil
		}
	} else {
		if avatar.Ban {
			return banImg, carbon.Now(), from, nil
		}
		// 优先使用 Gravatar 头像
		img, lastModified, imgErr = r.GetGravatarAvatar(hash)
		from = "gravatar"
		if imgErr != nil {
			// 如果 Gravatar 头像获取失败，则使用 QQ 头像
			img, lastModified, imgErr = r.GetQqAvatar(hash)
			from = "qq"
			if imgErr != nil {
				// 如果 QQ 头像获取失败，则使用默认头像
				img, lastModified, _ = r.GetDefaultAvatarByType(defaultAvatar, option)
				from = "weavatar"
			}
		}
	}

	// 审核头像
	if !avatar.Checked {
		go func() {
			_ = facades.Queue().Job(&jobs.ProcessAvatarCheck{}, []queue.Arg{
				{Type: "string", Value: hash},
				{Type: "string", Value: appid},
			}).Dispatch()
		}()
	}

	return img, lastModified, from, nil
}
