package jobs

import (
	"strconv"

	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"

	"weavatar/app/models"
	packagecdn "weavatar/packages/cdn"
	"weavatar/packages/helpers"
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
	if len(args) < 2 {
		facades.Log().Error("COS审核[队列参数不足]")
		return nil
	}

	hash, ok := args[0].(string)
	if !ok {
		facades.Log().Error("COS审核[队列参数断言失败] HASH:" + hash)
		return nil
	}

	appID, ok2 := args[1].(string)
	if !ok2 {
		facades.Log().Error("COS审核[队列参数断言失败] APPID:" + appID)
		return nil
	}

	if appID == "0" {
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
		var imageHash string
		exist := facades.Storage().Exists("upload/default/" + hash[:2] + "/" + hash)
		if exist {
			fileString, fileErr := facades.Storage().Get("upload/default/" + hash[:2] + "/" + hash)
			if fileErr != nil {
				facades.Log().Error("COS审核[文件读取失败] " + fileErr.Error())
				return nil
			}
			imageHash = helpers.MD5(fileString)
		} else {
			client := req.C()
			resp, reqErr := client.R().Get("http://proxy.server/http://0.gravatar.com/avatar/" + hash + ".png?s=600&r=g&d=404")
			if reqErr != nil || !resp.IsSuccessState() {
				return nil
			}
			imageHash = helpers.MD5(resp.String())
		}

		accessKey := facades.Config().GetString("qcloud.cos_check.access_key")
		secretKey := facades.Config().GetString("qcloud.cos_check.secret_key")
		bucket := facades.Config().GetString("qcloud.cos_check.bucket")
		checker := qcloud.NewCreator(accessKey, secretKey, bucket)

		var image models.Image
		err = facades.Orm().Query().Where("hash", imageHash).FirstOrFail(&image)
		if err != nil {
			isSafe, checkErr := checker.Check("https://weavatar.com/avatar/" + hash + ".png?s=400&d=404")
			if checkErr != nil {
				avatar.Checked = false
				err = facades.Orm().Query().Save(&avatar)
				if err != nil {
					facades.Log().Error("COS审核[数据更新失败] " + err.Error())
				}
				return nil
			}
			err = facades.Orm().Query().Create(&models.Image{
				Hash: imageHash,
				Ban:  !isSafe,
			})
			if err != nil {
				facades.Log().Error("COS审核[缓存数据创建失败] " + err.Error())
			}
		} else {
			avatar.Ban = image.Ban
			err = facades.Orm().Query().Save(&avatar)
			if err != nil {
				facades.Log().Error("COS审核[数据更新失败] " + err.Error())
				return err
			}

			if avatar.Ban {
				cdn := packagecdn.NewCDN()
				cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
			}
		}

		return nil
	}

	var avatar models.AppAvatar
	err := facades.Orm().Query().Where("avatar_hash", hash).First(&avatar)
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

	// 检查WeAvatar APP头像是否存在
	var imageHash string
	exist := facades.Storage().Exists("upload/app/" + strconv.Itoa(int(avatar.AppID)) + "/" + hash[:2] + "/" + hash)
	if exist {
		fileString, fileErr := facades.Storage().Get("upload/app/" + strconv.Itoa(int(avatar.AppID)) + "/" + hash[:2] + "/" + hash)
		if fileErr != nil {
			facades.Log().Error("COS审核[文件读取失败] " + fileErr.Error())
			return nil
		}
		imageHash = helpers.MD5(fileString)
	} else {
		return nil
	}

	accessKey := facades.Config().GetString("qcloud.cos_check.access_key")
	secretKey := facades.Config().GetString("qcloud.cos_check.secret_key")
	bucket := facades.Config().GetString("qcloud.cos_check.bucket")
	checker := qcloud.NewCreator(accessKey, secretKey, bucket)

	var image models.Image
	err = facades.Orm().Query().Where("hash", imageHash).FirstOrFail(&image)
	if err != nil {
		isSafe, checkErr := checker.Check("https://weavatar.com/avatar/" + hash + ".png?appid=" + strconv.Itoa(int(avatar.AppID)) + "&s=400&d=404")
		if checkErr != nil {
			avatar.Checked = false
			err = facades.Orm().Query().Save(&avatar)
			if err != nil {
				facades.Log().Error("COS审核[数据更新失败] " + err.Error())
			}
			return nil
		}
		err = facades.Orm().Query().Create(&models.Image{
			Hash: imageHash,
			Ban:  !isSafe,
		})
		if err != nil {
			facades.Log().Error("COS审核[缓存数据创建失败] " + err.Error())
		}
	} else {
		avatar.Ban = image.Ban
		err = facades.Orm().Query().Save(&avatar)
		if err != nil {
			facades.Log().Error("COS审核[数据更新失败] " + err.Error())
			return err
		}

		if avatar.Ban {
			cdn := packagecdn.NewCDN()
			cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
		}
	}

	return nil
}
