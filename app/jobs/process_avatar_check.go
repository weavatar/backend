package jobs

import (
	"github.com/goravel/framework/facades"

	"weavatar/app/models"
	packagecdn "weavatar/packages/cdn"
	"weavatar/packages/qcloud"
)

type ProcessAvatarCheck struct {
}

// Signature The name and signature of the job.
func (receiver *ProcessAvatarCheck) Signature() string {
	return "process_avatar_check"
}

// Handle Execute the job.
func (receiver *ProcessAvatarCheck) Handle(args ...any) error {
	if len(args) < 1 {
		facades.Log.Error("COS审核[队列参数不足]")
		return nil
	}

	// 断言参数
	hash, ok := args[0].(string)
	if !ok {
		facades.Log.Error("COS审核[队列参数断言失败] HASH:" + hash)
		return nil
	}

	accessKey := facades.Config.GetString("qcloud.cos_check.access_key")
	secretKey := facades.Config.GetString("qcloud.cos_check.secret_key")
	bucket := facades.Config.GetString("qcloud.cos_check.bucket")
	checker := qcloud.NewCreator(accessKey, secretKey, bucket)

	// 检查图片是否违规
	isSafe, err := checker.Check("https://weavatar.com/avatar/" + hash + ".png?s=400&d=404")
	if err != nil {
		return err
	}

	// 更新数据库
	var avatar models.Avatar
	err = facades.Orm.Query().Where("hash", hash).First(&avatar)
	if err != nil {
		facades.Log.Error("COS审核[数据库查询失败] " + err.Error())
		return err
	}

	avatar.Checked = true
	avatar.Ban = !isSafe
	err = facades.Orm.Query().Save(&avatar)
	if err != nil {
		facades.Log.Error("COS审核[数据库更新失败] " + err.Error())
		return err
	}

	// 刷新缓存
	cdn := packagecdn.NewCDN()
	cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})

	return nil
}
