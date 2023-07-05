package jobs

import (
	"github.com/disintegration/imaging"
	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"

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
		return nil
	}
	if avatar.Checked {
		return nil
	}

	// 首先标记为已审核，因为请求审核的时候会再次访问头像触发审核流程导致套娃
	avatar.Checked = true
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("COS审核[数据库更新失败] " + err.Error())
		return nil
	}

	// 检查WeAvatar头像是否存在
	_, imgErr := imaging.Open(facades.Storage().Path("upload/default/" + hash[:2] + "/" + hash))
	if imgErr != nil {
		// 不存在则请求Gravatar头像
		client := req.C()
		resp, reqErr := client.R().Get("http://proxy.server/http://0.gravatar.com/avatar/" + hash + ".png?s=10&r=g&d=404")
		if reqErr != nil || !resp.IsSuccessState() {
			// 也不存在说明是QQ头像或者默认头像，不需要审核
			return nil
		}
	}

	accessKey := facades.Config().GetString("qcloud.cos_check.access_key")
	secretKey := facades.Config().GetString("qcloud.cos_check.secret_key")
	bucket := facades.Config().GetString("qcloud.cos_check.bucket")
	checker := qcloud.NewCreator(accessKey, secretKey, bucket)

	isSafe, err := checker.Check("https://weavatar.com/avatar/" + hash + ".png?s=400&d=404")
	if err != nil {
		avatar.Checked = false
		err = facades.Orm().Query().Save(&avatar)
		if err != nil {
			facades.Log().Error("COS审核[数据库更新失败] " + err.Error())
		}
		return nil
	}

	avatar.Ban = !isSafe
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().Error("COS审核[数据库更新失败] " + err.Error())
		return err
	}

	if avatar.Ban {
		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
	}

	return nil
}
