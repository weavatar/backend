package jobs

import (
	"errors"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"

	"weavatar/app/models"
	packagecdn "weavatar/packages/cdn"
	"weavatar/packages/helpers"
)

type ProcessAvatarUpdate struct {
}

// Signature The name and signature of the job.
func (receiver *ProcessAvatarUpdate) Signature() string {
	return "process_avatar_update"
}

// Handle Execute the job.
func (receiver *ProcessAvatarUpdate) Handle(args ...any) error {
	if len(args) < 2 {
		facades.Log.Error("头像更新[队列参数不足]")
		return nil
	}

	// 断言参数
	hash, ok := args[0].(string)
	if !ok {
		facades.Log.Error("头像更新[队列参数断言失败] HASH:" + hash)
		return nil
	}
	path, ok2 := args[1].(string)
	if !ok2 {
		facades.Log.Error("头像更新[队列参数断言失败] PATH:" + path)
		return nil
	}

	// 检查图片是否存在
	if !facades.Storage.Exists(path) {
		facades.Log.Error("头像更新[文件不存在] " + path)
		return nil
	}

	// 计算文件的 MD5 值
	fileString, fileErr := facades.Storage.Get(path)
	if fileErr != nil {
		facades.Log.Error("头像更新[读取文件失败] " + path)
		return nil
	}
	fileHash := helpers.MD5(fileString)

	// 下载新头像
	var imgErr error
	from := "gravatar"
	client := req.C().SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.34")

	// 首先从 Gravatar 下载
	resp, reqErr := client.R().Get("http://proxy.server/http://0.gravatar.com/avatar/" + hash + ".png?s=1000&r=g&d=404")
	if reqErr != nil || !resp.IsSuccessState() {
		// 如果 Gravatar 下载失败，再从 QQ 下载
		from = "qq"
		type qqHash struct {
			Hash string `gorm:"primaryKey"`
			Qq   uint   `gorm:"type:bigint;not null"`
		}
		var qq qqHash

		err := facades.Orm.Connection("hash").Query().Table("qq_mails").Where("hash", hash).First(&qq)
		if err != nil {
			return nil
		}

		if qq.Qq == 0 {
			return nil
		}

		resp, reqErr = client.R().SetQueryParams(map[string]string{
			"b":  "qq",
			"nk": strconv.Itoa(int(qq.Qq)),
			"s":  "640",
		}).Get("http://q1.qlogo.cn/g")
		if reqErr != nil || !resp.IsSuccessState() {
			return nil
		}

		length, lengthErr := strconv.Atoi(resp.GetHeader("Content-Length"))
		if length < 6400 || lengthErr != nil {
			// 如果图片小于 6400 字节，则尝试获取 100 尺寸的图片
			resp, reqErr = client.R().SetQueryParams(map[string]string{
				"b":  "qq",
				"nk": strconv.Itoa(int(qq.Qq)),
				"s":  "100",
			}).Get("http://q1.qlogo.cn/g")
			if reqErr != nil || !resp.IsSuccessState() {
				return nil
			}
		}
	}

	reader := strings.NewReader(resp.String())
	_, imgErr = imaging.Decode(reader)
	if imgErr != nil {
		facades.Log.Warning("头像更新[图片不正常] ", imgErr.Error())
		return nil
	}

	delErr := facades.Storage.Delete(path)
	if delErr != nil {
		facades.Log.Error("头像更新[删除原有头像失败] " + path)
		return nil
	}

	saveErr := facades.Storage.Put("cache/"+from+"/"+hash[:2]+"/"+hash, resp.String())
	if saveErr != nil {
		return nil
	}

	if fileHash != helpers.MD5(resp.String()) {
		var avatar models.Avatar
		err := facades.Orm.Query().Where("hash", hash).First(&avatar)
		if err != nil || avatar.Hash == nil || avatar.UserID != nil {
			if err != nil {
				facades.Log.Error("头像更新[数据库查询失败] " + err.Error())
				return err
			} else {
				// 这里有2种可能：1.数据库中没有这个头像，但是缓存中有，所以需要删除缓存 2.用户已经上传了头像，不再需要缓存，所以也需要删除缓存
				facades.Log.Error("头像更新[数据库查询为空] HASH:" + hash)
				_ = facades.Storage.Delete("cache/" + from + "/" + hash[:2] + "/" + hash)
			}

			return errors.New("头像更新[数据库查询为空] HASH:" + hash)
		}

		avatar.Checked = false
		avatar.Ban = false
		err = facades.Orm.Query().Save(&avatar)
		if err != nil {
			facades.Log.Error("头像更新[数据库更新失败] " + err.Error())
			return err
		}

		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
	}

	return nil
}
