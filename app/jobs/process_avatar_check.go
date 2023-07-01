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
		facades.Log().Error("COS审核[队列参数不足]")
		return nil
	}

	hash, ok := args[0].(string)
	if !ok {
		facades.Log().Error("COS审核[队列参数断言失败] HASH:" + hash)
		return nil
	}

	var avatar models.Avatar
	err := facades.Orm().Query().Where("hash", hash).First(&avatar)
	if err != nil {
		facades.Log().Error("COS审核[数据库查询失败] " + err.Error())
		return err
	}

	type qqHash struct {
		Hash string `gorm:"primaryKey"`
		Qq   uint   `gorm:"type:bigint;not null"`
	}
	var qq qqHash

	err = facades.Orm().Connection("hash").Query().Table("qq_mails").Where("hash", hash).First(&qq)
	if err != nil {
		facades.Log().Error("COS审核[数据库查询失败] " + err.Error())
		return err
	}

	if qq.Qq != 0 {
		avatar.Checked = true
		avatar.Ban = false
		err = facades.Orm().Query().Save(&avatar)
		if err != nil {
			facades.Log().Error("COS审核[数据库更新失败] " + err.Error())
			return err
		}

		return nil
	}

	accessKey := facades.Config().GetString("qcloud.cos_check.access_key")
	secretKey := facades.Config().GetString("qcloud.cos_check.secret_key")
	bucket := facades.Config().GetString("qcloud.cos_check.bucket")
	checker := qcloud.NewCreator(accessKey, secretKey, bucket)

	isSafe, err := checker.Check("https://weavatar.com/avatar/" + hash + ".png?s=400&d=404")
	if err != nil {
		return err
	}

	avatar.Checked = true
	avatar.Ban = !isSafe
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("COS审核[数据库更新失败] " + err.Error())
		return err
	}

	cdn := packagecdn.NewCDN()
	cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})

	return nil
}
