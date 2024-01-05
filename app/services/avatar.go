package services

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/goki/freetype/truetype"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/haozi-team/letteravatar"
	"github.com/imroc/req/v3"
	"github.com/ipsn/go-adorable"
	"github.com/issue9/identicon/v2"
	"github.com/o1egl/govatar"
	"golang.org/x/exp/slices"
	"weavatar/app/jobs"
	"weavatar/app/models"
	"weavatar/pkg/helper"
)

// Avatar 头像服务
type Avatar interface {
	Sanitize(ctx http.Context) (appid uint, hash string, imageExt string, size int, forceDefault bool, defaultAvatar string)
	GetQQ(hash string) (qq int, img []byte, lastModified carbon.Carbon, err error)
	GetGravatar(hash string) (img []byte, lastModified carbon.Carbon, err error)
	GetDefault(defaultAvatar string, option []string) ([]byte, carbon.Carbon, error)
	GetDefaultByType(avatarType string, option []string) ([]byte, carbon.Carbon, error)
	GetAvatar(appid uint, hash string, defaultAvatar string, option []string) ([]byte, carbon.Carbon, string, error)
}

type AvatarImpl struct {
	BanImage []byte
	Font     *truetype.Font
	Client   *req.Client
}

func NewAvatarImpl() *AvatarImpl {
	ban, err := os.ReadFile(facades.Storage().Path("default/ban.png"))
	if err != nil {
		panic(err)
	}
	fontStr, err := os.ReadFile(facades.Storage().Path("fonts/HarmonyOS_Sans_SC_Medium.ttf"))
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(fontStr)
	if err != nil {
		panic(err)
	}

	client := req.C()
	client.SetTimeout(5 * time.Second)
	client.SetCommonRetryCount(2)
	client.ImpersonateSafari()

	return &AvatarImpl{
		BanImage: ban,
		Font:     font,
		Client:   client,
	}
}

// Sanitize 消毒头像请求
func (r *AvatarImpl) Sanitize(ctx http.Context) (appid uint, hash string, imageExt string, size int, forceDefault bool, defaultAvatar string) {
	hashExt := strings.Split(ctx.Request().Input("hash", ""), ".")
	appid = uint(ctx.Request().InputInt("appid", 0))

	hash = strings.ToLower(hashExt[0]) // Hash 转小写
	imageExt = "webp"                  // 默认为 WEBP 格式

	if len(hashExt) > 1 {
		imageExt = hashExt[1]
	}
	imageSlices := []string{"png", "jpg", "jpeg", "gif", "webp", "tiff", "avif"} // heif 由于 vips 的 bug 暂时移除
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

	forceDefaultRaw := ctx.Request().Input("f", "")
	if len(forceDefaultRaw) == 0 {
		forceDefaultRaw = ctx.Request().Input("forcedefault", "n")
	}
	forceDefault = false
	forceDefaultSlices := []string{"y", "yes"}
	if slices.Contains(forceDefaultSlices, forceDefaultRaw) {
		forceDefault = true
	}

	defaultAvatar = ctx.Request().Input("d", "")
	if len(defaultAvatar) == 0 {
		defaultAvatar = ctx.Request().Input("default", "")
	}
	defaultAvatarSlices := []string{"404", "mp", "mm", "mystery", "identicon", "monsterid", "wavatar", "retro", "robohash", "blank", "letter"}
	if !slices.Contains(defaultAvatarSlices, defaultAvatar) {
		// 如果不是预设的默认头像，则检查是否是合法的 URL
		if !helper.IsURL(defaultAvatar) {
			defaultAvatar = ""
		}
	}

	match, _ := regexp.MatchString("^[a-f0-9]{32}$", hash)
	if !match {
		forceDefault = true
	}

	return appid, hash, imageExt, size, forceDefault, defaultAvatar
}

// GetQQ 通过 QQ 号获取头像
func (r *AvatarImpl) GetQQ(hash string) (qq int, img []byte, lastModified carbon.Carbon, err error) {
	hashIndex, err := strconv.ParseInt(hash[:10], 16, 64)
	if err != nil {
		return 0, nil, carbon.Now(), err
	}
	tableIndex := (hashIndex % int64(500)) + 1
	type qqHash struct {
		Hash string `gorm:"primaryKey"`
		QQ   int    `gorm:"type:bigint;not null"`
	}

	var qqModel qqHash
	if err = facades.Orm().Connection("hash").Query().Table("qq_"+strconv.Itoa(int(tableIndex))).Where("hash", hash).First(&qqModel); err != nil {
		return 0, nil, carbon.Now(), err
	}
	if qqModel.QQ == 0 {
		return 0, nil, carbon.Now(), errors.New("未找到对应的 QQ 号")
	}

	if facades.Storage().Exists("cache/qq/" + hash[:2] + "/" + hash) {
		img, err = os.ReadFile(facades.Storage().Path("cache/qq/" + hash[:2] + "/" + hash))
		lastModifiedStd, lastModifiedErr := facades.Storage().LastModified("cache/qq/" + hash[:2] + "/" + hash)
		if err == nil && lastModifiedErr == nil {
			return qqModel.QQ, img, carbon.FromStdTime(lastModifiedStd), nil
		}
	}

	resp, reqErr := r.Client.R().SetQueryParams(map[string]string{
		"b":  "qq",
		"nk": strconv.Itoa(qqModel.QQ),
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
			"nk": strconv.Itoa(qqModel.QQ),
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
	lastModifiedStd, err := facades.Storage().LastModified("cache/qq/" + hash[:2] + "/" + hash)
	if err != nil {
		return 0, nil, carbon.Now(), err
	}

	return qqModel.QQ, resp.Bytes(), carbon.FromStdTime(lastModifiedStd), nil
}

// GetGravatar 通过 Gravatar 获取头像
func (r *AvatarImpl) GetGravatar(hash string) (img []byte, lastModified carbon.Carbon, err error) {
	if facades.Storage().Exists("cache/gravatar/" + hash[:2] + "/" + hash) {
		img, err = os.ReadFile(facades.Storage().Path("cache/gravatar/" + hash[:2] + "/" + hash))
		lastModifiedStd, lastModifiedErr := facades.Storage().LastModified("cache/gravatar/" + hash[:2] + "/" + hash)
		if err == nil && lastModifiedErr == nil {
			return img, carbon.FromStdTime(lastModifiedStd), nil
		}
	}

	resp, reqErr := r.Client.R().Get("http://proxy.server/https://s.gravatar.com/avatar/" + hash + ".png?s=600&r=g&d=404")
	if reqErr != nil || !resp.IsSuccessState() {
		return nil, carbon.Now(), errors.New("获取 Gravatar 头像失败")
	}

	// 保存图片
	if err = facades.Storage().Put("cache/gravatar/"+hash[:2]+"/"+hash, resp.String()); err != nil {
		return nil, carbon.Now(), err
	}

	lastModifiedStd, lastModifiedErr := facades.Storage().LastModified("cache/gravatar/" + hash[:2] + "/" + hash)
	if lastModifiedErr != nil {
		return nil, carbon.Now(), err
	}

	return resp.Bytes(), carbon.FromStdTime(lastModifiedStd), nil
}

// GetDefault 通过默认参数获取头像
func (r *AvatarImpl) GetDefault(defaultAvatar string, option []string) (img []byte, lastModified carbon.Carbon, err error) {
	if defaultAvatar == "404" {
		return nil, carbon.Now(), nil
	}

	if defaultAvatar == "" {
		img, err = os.ReadFile(facades.Storage().Path("default/default.png"))
		if err != nil {
			return nil, carbon.Now(), err
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "mp" {
		img, err = os.ReadFile(facades.Storage().Path("default/mp.png"))
		if err != nil {
			return nil, carbon.Now(), err
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "identicon" {
		imgStd := identicon.Make(identicon.Style1, 1200, color.RGBA{R: 255, A: 100}, color.RGBA{R: 102, G: 204, B: 255, A: 255}, []byte(option[0]))
		var buf bytes.Buffer
		err = png.Encode(&buf, imgStd)
		if err != nil {
			return nil, carbon.Now(), err
		}

		return buf.Bytes(), carbon.Now(), nil
	}

	if defaultAvatar == "monsterid" {
		var imgStd image.Image
		imgStd, err = govatar.GenerateForUsername(govatar.FEMALE, option[0])
		if err != nil {
			return nil, carbon.Now(), err
		}
		var buf bytes.Buffer
		err = png.Encode(&buf, imgStd)
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
		imgStd := ii.Make([]byte(option[0]))
		var buf bytes.Buffer
		err = png.Encode(&buf, imgStd)
		if err != nil {
			return nil, carbon.Now(), err
		}

		return buf.Bytes(), carbon.Now(), nil
	}

	if defaultAvatar == "robohash" {
		var imgStd image.Image
		imgStd, err = govatar.GenerateForUsername(govatar.MALE, option[0])
		if err != nil {
			return nil, carbon.Now(), err
		}
		var buf bytes.Buffer
		err = png.Encode(&buf, imgStd)
		if err != nil {
			return nil, carbon.Now(), err
		}

		return buf.Bytes(), carbon.Now(), nil
	}

	if defaultAvatar == "blank" {
		img, err = os.ReadFile(facades.Storage().Path("default/blank.png"))
		if err != nil {
			return nil, carbon.Now(), err
		}

		return img, carbon.Now(), nil
	}

	if defaultAvatar == "letter" {
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

		var imgStd image.Image
		imgStd, err = letteravatar.Draw(1000, letters, &letteravatar.Options{
			Font:       r.Font,
			FontSize:   fontSize,
			PaletteKey: option[1], // 对相同的字符串使用相同的颜色
		})
		if err != nil {
			return nil, carbon.Now(), err
		}

		var buf bytes.Buffer
		err = png.Encode(&buf, imgStd)
		if err != nil {
			return nil, carbon.Now(), err
		}

		return buf.Bytes(), carbon.Now(), nil
	}

	// 都不是的话，返回 nil（自定义 URL）
	return nil, carbon.Now(), nil
}

// GetDefaultByType 通过默认头像类型获取头像
func (r *AvatarImpl) GetDefaultByType(avatarType string, option []string) (img []byte, lastModified carbon.Carbon, err error) {
	switch avatarType {
	case "404":
		img, lastModified, err = r.GetDefault("404", option)
	case "mp", "mm", "mystery":
		img, lastModified, err = r.GetDefault("mp", option)
	case "identicon":
		img, lastModified, err = r.GetDefault("identicon", option)
	case "monsterid":
		img, lastModified, err = r.GetDefault("monsterid", option)
	case "wavatar":
		img, lastModified, err = r.GetDefault("wavatar", option)
	case "retro":
		img, lastModified, err = r.GetDefault("retro", option)
	case "robohash":
		img, lastModified, err = r.GetDefault("robohash", option)
	case "blank":
		img, lastModified, err = r.GetDefault("blank", option)
	case "letter":
		img, lastModified, err = r.GetDefault("letter", option)
	default:
		img, lastModified, err = r.GetDefault(avatarType, option)
	}

	return img, lastModified, err
}

// GetAvatar 获取头像
func (r *AvatarImpl) GetAvatar(appid uint, hash string, defaultAvatar string, option []string) (img []byte, lastModified carbon.Carbon, from string, err error) {
	var avatar models.Avatar
	var appAvatar models.AppAvatar

	from = "weavatar"

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
				img, err = os.ReadFile(facades.Storage().Path("upload/app/" + strconv.Itoa(int(appAvatar.AppID)) + "/" + hash[:2] + "/" + hash))
				lastModified = appAvatar.UpdatedAt.Carbon
			}
		} else {
			if avatar.Ban {
				return r.BanImage, carbon.Now(), from, nil
			} else {
				img, err = os.ReadFile(facades.Storage().Path("upload/default/" + hash[:2] + "/" + hash))
				lastModified = avatar.UpdatedAt.Carbon
			}
		}

		if err != nil {
			facades.Log().With(map[string]any{
				"avatar": avatar,
				"error":  err.Error(),
			}).Hint("记录显示用户存在头像，但是无法找到头像文件").Warning("WeAvatar[头像匹配失败]")
			img, lastModified, _ = r.GetDefaultByType(defaultAvatar, option)
			return img, lastModified, from, nil
		}
	} else {
		if avatar.Ban {
			return r.BanImage, carbon.Now(), from, nil
		}

		img, lastModified, err = r.GetGravatar(hash)
		from = "gravatar"
		if err != nil {
			var qq int
			qq, img, lastModified, err = r.GetQQ(hash)
			from = "qq"
			if err != nil {
				img, lastModified, _ = r.GetDefaultByType(defaultAvatar, option)
				from = "weavatar"
			}
			// 自动免审默认头像和 QQ 头像
			if !avatar.Checked {
				if qq != 0 {
					raw := strconv.Itoa(qq) + "@qq.com"
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
		go func(h string, a uint) {
			err := facades.Queue().Job(&jobs.ProcessAvatarCheck{}, []queue.Arg{
				{Type: "string", Value: hash},
				{Type: "uint", Value: appid},
			}).Dispatch()
			if err != nil {
				facades.Log().With(map[string]any{
					"hash":  hash,
					"appid": appid,
					"error": err.Error(),
				}).Error("WeAvatar[任务分发失败]")
			}
		}(hash, appid)
	}

	return img, lastModified, from, nil
}
