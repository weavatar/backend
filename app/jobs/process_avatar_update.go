package jobs

import (
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
		hashIndex, hashErr := strconv.ParseInt(hash[:10], 16, 64)
		if hashErr != nil {
			return nil
		}
		tableIndex := (hashIndex % int64(4000)) + 1
		type qqHash struct {
			Hash string `gorm:"column:hash"`
			Qq   string `gorm:"column:qq"`
		}
		var qq qqHash

		err := facades.Orm.Connection("hash").Query().Table("qq_"+strconv.Itoa(int(tableIndex))).Where("hash", hash).First(&qq)
		if err != nil {
			return nil
		}

		if qq.Qq == "" {
			return nil
		}

		resp, reqErr = client.R().SetQueryParams(map[string]string{
			"b":  "qq",
			"nk": qq.Qq,
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
				"nk": qq.Qq,
				"s":  "100",
			}).Get("http://q1.qlogo.cn/g")
			if reqErr != nil || !resp.IsSuccessState() {
				return nil
			}
		}
	}

	// 检查图片是否正常
	reader := strings.NewReader(resp.String())
	_, imgErr = imaging.Decode(reader)
	if imgErr != nil {
		facades.Log.Warning("头像更新[图片不正常] ", imgErr.Error())
		return nil
	}

	// 删除原有头像
	delErr := facades.Storage.Delete(path)
	if delErr != nil {
		facades.Log.Error("头像更新[删除原有头像失败] " + path)
		return nil
	}

	// 保存图片
	saveErr := facades.Storage.Put("cache/"+from+"/"+hash[:2]+"/"+hash, resp.String())
	if saveErr != nil {
		return nil
	}

	// 检查图片是否变化
	if fileHash == helpers.MD5(resp.String()) {
		// 更新数据库
		var avatar models.Avatar
		err := facades.Orm.Query().Where("hash", hash).First(&avatar)
		if err != nil || avatar.Hash == nil {
			if err != nil {
				facades.Log.Error("头像更新[数据库查询失败] " + err.Error())
			} else {
				facades.Log.Error("头像更新[数据库查询失败] HASH:" + hash)
			}

			return err
		}

		avatar.Checked = false
		avatar.Ban = false
		err = facades.Orm.Query().Save(&avatar)
		if err != nil {
			facades.Log.Error("头像更新[数据库更新失败] " + err.Error())
			return err
		}

		// 刷新缓存
		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl([]string{"https://weavatar.com/avatar/" + hash + "*", "http://weavatar.com/avatar/" + hash + "*"})
	}

	return nil
}
