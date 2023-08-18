package jobs

import (
	"github.com/goravel/framework/facades"

	"weavatar/app/models"
	packagecdn "weavatar/pkg/cdn"
)

type ProcessAvatarUpdate struct {
}

// urls 待刷新的URL列表
var urls []string

// Signature The name and signature of the job.
func (receiver *ProcessAvatarUpdate) Signature() string {
	return "process_avatar_update"
}

// Handle Execute the job.
func (receiver *ProcessAvatarUpdate) Handle(args ...any) error {
	if len(args) < 2 {
		facades.Log().Error("头像更新[队列参数不足]")
		return nil
	}

	// 断言参数
	hash, ok := args[0].(string)
	if !ok {
		facades.Log().Error("头像更新[队列参数断言失败] HASH:" + hash)
		return nil
	}
	path, ok2 := args[1].(string)
	if !ok2 {
		facades.Log().Error("头像更新[队列参数断言失败] PATH:" + path)
		return nil
	}

	// 检查图片是否存在
	if !facades.Storage().Exists(path) {
		facades.Log().Error("头像更新[文件不存在] " + path)
		return nil
	}

	var avatar models.Avatar
	err := facades.Orm().Query().Where("hash", hash).First(&avatar)
	if err != nil {
		facades.Log().Error("头像更新[数据库查询失败] " + err.Error())
		return nil
	}
	if avatar.Hash == nil || avatar.UserID != nil {
		// 这里有2种可能：1.数据库中没有这个头像，但是缓存中有，所以需要删除缓存 2.用户已经上传了头像，不再需要缓存，所以也需要删除缓存
		facades.Log().Error("头像更新[数据库查询为空] HASH:" + hash)
		_ = facades.Storage().Delete("cache/gravatar/" + hash[:2] + "/" + hash)
		_ = facades.Storage().Delete("cache/qq/" + hash[:2] + "/" + hash)
	}

	delErr := facades.Storage().Delete(path)
	if delErr != nil {
		facades.Log().Error("头像更新[删除原有头像失败] " + path)
	}

	avatar.Checked = false
	avatar.Ban = false
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("头像更新[数据库更新失败] " + err.Error())
		return nil
	}

	urls = append(urls, "weavatar.com/avatar/"+hash)
	if len(urls) >= 30 {
		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl(urls)
		urls = []string{}
	}

	return nil
}
