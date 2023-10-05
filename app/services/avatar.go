package services

import (
	"bytes"
	"errors"
	"image/color"
	"image/png"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/HaoZi-Team/letteravatar"
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

// Avatar 头像服务
type Avatar interface {
	Sanitize(ctx http.Context) (string, string, string, int, bool, string)
	GetQQ(hash string) (uint, []byte, carbon.Carbon, error)
	GetGravatar(hash string) ([]byte, carbon.Carbon, error)
	GetDefault(defaultAvatar string, option []string) ([]byte, carbon.Carbon, error)
	GetDefaultByType(avatarType string, option []string) ([]byte, carbon.Carbon, error)
	GetAvatar(appid string, hash string, defaultAvatar string, option []string) ([]byte, carbon.Carbon, string, error)
}

type AvatarImpl struct {
	BanImage []byte
	Client   *req.Client
}

func NewAvatarImpl() *AvatarImpl {
	ban, _ := os.ReadFile(facades.Storage().Path("default/ban.png"))
	client := req.C()
	client.SetTimeout(5 * time.Second)
	client.SetCommonRetryCount(2)
	client.ImpersonateSafari()

	return &AvatarImpl{
		BanImage: ban,
		Client:   client,
	}
}

// Sanitize 消毒头像请求
func (r *AvatarImpl) Sanitize(ctx http.Context) (string, string, string, int, bool, string) {
	hashExt := strings.Split(ctx.Request().Input("hash", ""), ".")
	appid := ctx.Request().Input("appid", "0")

	hash := strings.ToLower(hashExt[0]) // Hash 转小写
	imageExt := "webp"                  // 默认为 WEBP 格式

	if len(hashExt) > 1 {
		imageExt = hashExt[1]
	}
	imageSlices := []string{"png", "jpg", "jpeg", "gif", "webp", "tiff", "heif", "avif"}
	if !slices.Contains(imageSlices, imageExt) {
		imageExt = "webp"
	}

	sizeStr := ctx.Request().Input("s", "")
	if len(sizeStr) == 0 {
		sizeStr = ctx.Request().Input("size", "80")
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		size = 80
	}
	if size > 2000 {
		size = 2000
	}
	if size < 10 {
		size = 10
	}

	forceDefault := ctx.Request().Input("f", "")
	if len(forceDefault) == 0 {
		forceDefault = ctx.Request().Input("forcedefault", "n")
	}
	forceDefaultBool := false
	forceDefaultSlices := []string{"y", "yes"}
	if slices.Contains(forceDefaultSlices, forceDefault) {
		forceDefaultBool = true
	}

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

	match, _ := regexp.MatchString("^[a-f0-9]{32}$", hash)
	if !match {
		forceDefaultBool = true
	}

	return appid, hash, imageExt, size, forceDefaultBool, defaultAvatar
}

// GetQQ 通过 QQ 号获取头像
func (r *AvatarImpl) GetQQ(hash string) (uint, []byte, carbon.Carbon, error) {
	hashIndex, err := strconv.ParseInt(hash[:10], 16, 64)
	if err != nil {
		return 0, nil, carbon.Now(), err
	}
	tableIndex := (hashIndex % int64(500)) + 1
	type qqHash struct {
		Hash string `gorm:"primaryKey"`
		QQ   uint   `gorm:"type:bigint;not null"`
	}

	var qq qqHash
	if err = facades.Orm().Connection("hash").Query().Table("qq_"+strconv.Itoa(int(tableIndex))).Where("hash", hash).First(&qq); err != nil {
		return 0, nil, carbon.Now(), err
	}
	if qq.QQ == 0 {
		return 0, nil, carbon.Now(), errors.New("未找到对应的 QQ 号")
	}

	if facades.Storage().Exists("cache/qq/" + hash[:2] + "/" + hash) {
		img, imgErr := os.ReadFile("cache/qq/" + hash[:2] + "/" + hash)
		lastModified, err := facades.Storage().LastModified("cache/qq/" + hash[:2] + "/" + hash)
		if imgErr == nil && err == nil {
			return qq.QQ, img, carbon.FromStdTime(lastModified), nil
		}
	}

	resp, reqErr := r.Client.R().SetQueryParams(map[string]string{
		"b":  "qq",
		"nk": strconv.Itoa(int(qq.QQ)),
		"s":  "640",
	}).Get("http://q1.qlogo.cn/g")
	if !resp.IsSuccessState() {
		if reqErr != nil {
			return 0, nil, carbon.Now(), reqErr
		} else {
			return 0, nil, carbon.Now(), errors.New("获取 QQ头像 失败")
		}
	}

	length, lengthErr := strconv.Atoi(resp.GetHeader("Content-Length"))
	if length < 6400 || lengthErr != nil {
		// 如果图片小于 6400 字节，则尝试获取 100 尺寸的图片
		resp, reqErr = r.Client.R().SetQueryParams(map[string]string{
			"b":  "qq",
			"nk": strconv.Itoa(int(qq.QQ)),
			"s":  "100",
		}).Get("http://q1.qlogo.cn/g")
		if !resp.IsSuccessState() {
			if reqErr != nil {
				return 0, nil, carbon.Now(), reqErr
			} else {
				return 0, nil, carbon.Now(), errors.New("获取 QQ头像 失败")
			}
		}
	}

	if err := facades.Storage().Put("cache/qq/"+hash[:2]+"/"+hash, resp.String()); err != nil {
		return 0, nil, carbon.Now(), err
	}
	lastModified, err := facades.Storage().LastModified("cache/qq/" + hash[:2] + "/" + hash)
	if err != nil {
		return 0, nil, carbon.Now(), err
	}

	return qq.QQ, resp.Bytes(), carbon.FromStdTime(lastModified), nil
}

// GetGravatar 通过 Gravatar 获取头像
func (r *AvatarImpl) GetGravatar(hash string) ([]byte, carbon.Carbon, error) {
	var img []byte
	var imgErr error

	if facades.Storage().Exists("cache/gravatar/" + hash[:2] + "/" + hash) {
		img, imgErr = os.ReadFile("cache/gravatar/" + hash[:2] + "/" + hash)
		lastModified, err := facades.Storage().LastModified("cache/gravatar/" + hash[:2] + "/" + hash)
		if imgErr == nil && err == nil {
			return img, carbon.FromStdTime(lastModified), nil
		}
	}

	resp, reqErr := r.Client.R().Get("http://proxy.server/https://secure.gravatar.com/avatar/" + hash + ".png?s=600&r=g&d=404")
	if reqErr != nil || !resp.IsSuccessState() {
		return nil, carbon.Now(), errors.New("获取 Gravatar 头像失败")
	}

	// 保存图片
	if err := facades.Storage().Put("cache/gravatar/"+hash[:2]+"/"+hash, resp.String()); err != nil {
		return nil, carbon.Now(), err
	}

	lastModified, err := facades.Storage().LastModified("cache/gravatar/" + hash[:2] + "/" + hash)
	if err != nil {
		return nil, carbon.Now(), err
	}

	return resp.Bytes(), carbon.FromStdTime(lastModified), nil
}

// GetDefault 通过默认参数获取头像
func (r *AvatarImpl) GetDefault(defaultAvatar string, option []string) ([]byte, carbon.Carbon, error) {
	if defaultAvatar == "404" {
		return nil, carbon.Now(), nil
	}

	if defaultAvatar == "" {
		img, imgErr := os.ReadFile(facades.Storage().Path("default/default.png"))
		if imgErr != nil {
			return nil, carbon.Now(), imgErr
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "mp" {
		img, imgErr := os.ReadFile(facades.Storage().Path("default/mp.png"))
		if imgErr != nil {
			return nil, carbon.Now(), imgErr
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "identicon" {
		img := identicon.Make(identicon.Style1, 1200, color.RGBA{R: 255, A: 100}, color.RGBA{R: 102, G: 204, B: 255, A: 255}, []byte(option[0]))
		var buf bytes.Buffer
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, carbon.Now(), err
		}

		return buf.Bytes(), carbon.Now(), nil
	}

	if defaultAvatar == "monsterid" {
		img, err := govatar.GenerateForUsername(govatar.FEMALE, option[0])
		if err != nil {
			return nil, carbon.Now(), err
		}
		var buf bytes.Buffer
		err = png.Encode(&buf, img)
		if err != nil {
			return nil, carbon.Now(), err
		}

		return buf.Bytes(), carbon.Now(), nil
	}

	if defaultAvatar == "wavatar" {
		return adorable.PseudoRandom([]byte(option[0])), carbon.Now(), nil
	}

	if defaultAvatar == "retro" {
		ii := identicon.New(identicon.Style2, 1200, color.RGBA{R: 255, A: 100}, color.RGBA{R: 102, G: 204, B: 255, A: 255})
		img := ii.Make([]byte(option[0]))
		var buf bytes.Buffer
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, carbon.Now(), err
		}

		return buf.Bytes(), carbon.Now(), nil
	}

	if defaultAvatar == "robohash" {
		img, err := govatar.GenerateForUsername(govatar.MALE, option[0])
		if err != nil {
			return nil, carbon.Now(), err
		}
		var buf bytes.Buffer
		err = png.Encode(&buf, img)
		if err != nil {
			return nil, carbon.Now(), err
		}

		return buf.Bytes(), carbon.Now(), nil
	}

	if defaultAvatar == "blank" {
		img, imgErr := os.ReadFile(facades.Storage().Path("default/blank.png"))
		if imgErr != nil {
			return nil, carbon.Now(), imgErr
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "letter" {
		fontStr, err := os.ReadFile(facades.Storage().Path("fonts/HarmonyOS_Sans_SC_Medium.ttf"))
		if err != nil {
			return nil, carbon.Now(), err
		}
		font, fontErr := truetype.Parse(fontStr)
		if fontErr != nil {
			return nil, carbon.Now(), fontErr
		}

		fontSize := 0
		hasChinese := false
		letters := []rune(option[0])
		length := len(letters)
		if length > 4 {
			letters = letters[:4]
			length = 4
		}

		for _, w := range letters {
			if unicode.Is(unicode.Han, w) {
				hasChinese = true
				break
			}
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

		var buf bytes.Buffer
		err = png.Encode(&buf, img)
		if err != nil {
			return nil, carbon.Now(), err
		}

		return buf.Bytes(), carbon.Now(), nil
	}

	resp, reqErr := r.Client.R().Get(defaultAvatar)
	if reqErr != nil || !resp.IsSuccessState() {
		img, imgErr := os.ReadFile(facades.Storage().Path("default/default.png"))
		if imgErr != nil {
			return nil, carbon.Now(), imgErr
		}

		return img, carbon.Now(), nil
	}

	return resp.Bytes(), carbon.Now(), nil
}

// GetDefaultByType 通过默认头像类型获取头像
func (r *AvatarImpl) GetDefaultByType(avatarType string, option []string) ([]byte, carbon.Carbon, error) {
	var avatar []byte
	var lastModified carbon.Carbon
	var err error

	switch avatarType {
	case "404":
		avatar, lastModified, err = r.GetDefault("404", option)
	case "mp", "mm", "mystery":
		avatar, lastModified, err = r.GetDefault("mp", option)
	case "identicon":
		avatar, lastModified, err = r.GetDefault("identicon", option)
	case "monsterid":
		avatar, lastModified, err = r.GetDefault("monsterid", option)
	case "wavatar":
		avatar, lastModified, err = r.GetDefault("wavatar", option)
	case "retro":
		avatar, lastModified, err = r.GetDefault("retro", option)
	case "robohash":
		avatar, lastModified, err = r.GetDefault("robohash", option)
	case "blank":
		avatar, lastModified, err = r.GetDefault("blank", option)
	case "letter":
		avatar, lastModified, err = r.GetDefault("letter", option)
	default:
		avatar, lastModified, err = r.GetDefault(avatarType, option)
	}

	return avatar, lastModified, err
}

// GetAvatar 获取头像
func (r *AvatarImpl) GetAvatar(appid string, hash string, defaultAvatar string, option []string) ([]byte, carbon.Carbon, string, error) {
	var avatar models.Avatar
	var appAvatar models.AppAvatar
	var err error

	var img []byte
	var lastModified carbon.Carbon
	var imgErr error
	from := "weavatar"

	// 取头像数据
	_, err = facades.Orm().Query().Exec(`INSERT INTO avatars (hash, created_at, updated_at) VALUES (?, ?, ?) ON CONFLICT DO NOTHING`, hash, carbon.DateTime{Carbon: carbon.Now()}, carbon.DateTime{Carbon: carbon.Now()})
	if err != nil {
		facades.Log().With(map[string]any{
			"hash":  hash,
			"error": err.Error(),
		}).Hint("插入哈希到头像表失败").Error("WeAvatar[数据库错误]")
		return nil, carbon.Now(), from, err
	}
	err = facades.Orm().Query().Where("hash", hash).First(&avatar)
	if err != nil {
		facades.Log().With(map[string]any{
			"hash":  hash,
			"error": err.Error(),
		}).Hint("从头像表查询哈希失败").Error("WeAvatar[数据库错误]")
		return nil, carbon.Now(), from, err
	}
	if avatar.UserID != nil && avatar.Raw != nil {
		err = facades.Orm().Query().Where("app_id", appid).Where("avatar_hash", hash).First(&appAvatar)
		if err != nil {
			facades.Log().With(map[string]any{
				"hash":   hash,
				"avatar": avatar,
				"error":  err.Error(),
			}).Hint("从APP头像表查询哈希失败").Error("WeAvatar[数据库错误]")
			return nil, carbon.Now(), from, err
		}

		if appAvatar.AppID != 0 {
			if appAvatar.Ban {
				return r.BanImage, carbon.Now(), from, nil
			} else {
				img, imgErr = os.ReadFile(facades.Storage().Path("upload/app/" + strconv.Itoa(int(appAvatar.AppID)) + "/" + hash[:2] + "/" + hash))
				lastModified = appAvatar.UpdatedAt.Carbon
			}
		} else {
			if avatar.Ban {
				return r.BanImage, carbon.Now(), from, nil
			} else {
				img, imgErr = os.ReadFile(facades.Storage().Path("upload/default/" + hash[:2] + "/" + hash))
				lastModified = avatar.UpdatedAt.Carbon
			}
		}

		if imgErr != nil {
			facades.Log().With(map[string]any{
				"avatar": avatar,
				"error":  imgErr.Error(),
			}).Hint("记录显示用户存在头像，但是无法找到头像文件").Warning("WeAvatar[头像匹配失败]")
			img, lastModified, _ = r.GetDefaultByType(defaultAvatar, option)
			return img, lastModified, from, nil
		}
	} else {
		if avatar.Ban {
			return r.BanImage, carbon.Now(), from, nil
		}

		img, lastModified, imgErr = r.GetGravatar(hash)
		from = "gravatar"
		if imgErr != nil {
			var qq uint
			qq, img, lastModified, imgErr = r.GetQQ(hash)
			from = "qq"
			if imgErr != nil {
				img, lastModified, _ = r.GetDefaultByType(defaultAvatar, option)
				from = "weavatar"
			}
			// 自动免审默认头像和 QQ 头像
			if !avatar.Checked {
				if qq != 0 {
					raw := strconv.Itoa(int(qq)) + "@qq.com"
					avatar.Raw = &raw
				}
				avatar.Checked = true
				avatar.Ban = false
				err = facades.Orm().Query().Save(&avatar)
				if err != nil {
					facades.Log().With(map[string]any{
						"qq":     qq,
						"avatar": avatar,
						"error":  err.Error(),
					}).Hint("更新头像的审核状态失败").Error("WeAvatar[数据库错误]")
				}
			}
		}
	}

	if !avatar.Checked && from != "qq" {
		go func(h string, a string) {
			_ = facades.Queue().Job(&jobs.ProcessAvatarCheck{}, []queue.Arg{
				{Type: "string", Value: h},
				{Type: "string", Value: a},
			}).Dispatch()
		}(hash, appid)
	}

	return img, lastModified, from, nil
}
