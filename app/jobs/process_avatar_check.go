package jobs

import (
	"strconv"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"

	"weavatar/app/models"
	packagecdn "weavatar/pkg/cdn"
	"weavatar/pkg/helper"
	"weavatar/pkg/imagecheck"
)

type ProcessAvatarCheck struct {
}

// Signature The name and signature of the job.
func (receiver *ProcessAvatarCheck) Signature() string {
	return "process_avatar_check"
}

// Handle Execute the job.
func (receiver *ProcessAvatarCheck) Handle(args ...any) error {
	if status := facades.Config().GetString("app.status", "main"); status != "main" {
		return nil
	}

	if len(args) < 2 {
		facades.Log().With(map[string]any{
			"args": args,
		}).Error("图片审核[队列参数不足]")
		return nil
	}

	hash, ok := args[0].(string)
	if !ok {
		facades.Log().With(map[string]any{
			"hash": hash,
		}).Error("图片审核[队列参数断言失败]")
		return nil
	}

	appID, ok2 := args[1].(string)
	if !ok2 {
		facades.Log().With(map[string]any{
			"appID": appID,
		}).Error("图片审核[队列参数断言失败]")
		return nil
	}

	if appID == "0" {
		var avatar models.Avatar
		err := facades.Orm().Query().Where("hash", hash).First(&avatar)
		if err != nil {
			facades.Log().With(map[string]any{
				"hash": hash,
				"err":  err.Error(),
			}).Error("图片审核[数据查询失败]")
			return nil
		}
		if avatar.Checked || avatar.Hash == nil {
			return nil
		}

		avatar.Checked = true
		err = facades.Orm().Query().Save(&avatar)
		if err != nil {
			facades.Log().With(map[string]any{
				"avatar": avatar,
				"err":    err.Error(),
			}).Error("图片审核[数据库更新失败]")
			return nil
		}

		// 检查WeAvatar头像是否存在
		var imageHash string
		exist := facades.Storage().Exists("upload/default/" + hash[:2] + "/" + hash)
		if exist {
			fileString, fileErr := facades.Storage().Get("upload/default/" + hash[:2] + "/" + hash)
			if fileErr != nil {
				facades.Log().With(map[string]any{
					"hash": hash,
					"err":  fileErr.Error(),
				}).Error("图片审核[文件读取失败]")
				return nil
			}

			imageHash = helper.MD5(fileString)
			err = facades.Storage().Put("checker/"+imageHash[:2]+"/"+imageHash, fileString)
			if err != nil {
				facades.Log().With(map[string]any{
					"avatarHash": hash,
					"imageHash":  imageHash,
					"err":        err.Error(),
				}).Error("图片审核[文件缓存失败]")
				return nil
			}
		} else {
			client := req.C()
			client.SetTimeout(5 * time.Second)
			client.SetCommonRetryCount(2)
			client.ImpersonateSafari()

			resp, reqErr := client.R().Get("http://proxy.server/https://secure.gravatar.com/avatar/" + hash + ".png?s=600&r=g&d=404")
			if reqErr != nil || !resp.IsSuccessState() {
				return nil
			}

			imageHash = helper.MD5(resp.String())
			err = facades.Storage().Put("checker/"+imageHash[:2]+"/"+imageHash, resp.String())
			if err != nil {
				facades.Log().With(map[string]any{
					"avatarHash": hash,
					"imageHash":  imageHash,
					"err":        err.Error(),
				}).Error("图片审核[文件缓存失败]")
				return nil
			}
		}

		var image models.Image
		err = facades.Orm().Query().Where("hash", imageHash).FirstOrFail(&image)
		if err != nil {
			checker := imagecheck.NewChecker()
			ban, checkErr := checker.Check("https://weavatar.com/avatar/" + hash + ".png?s=600&d=404")
			if checkErr != nil {
				facades.Log().With(map[string]any{
					"hash": hash,
					"err":  checkErr.Error(),
				}).Error("图片审核[审核失败]")
				avatar.Checked = false
				err = facades.Orm().Query().Save(&avatar)
				if err != nil {
					facades.Log().With(map[string]any{
						"avatar": avatar,
						"err":    err.Error(),
					}).Error("图片审核[数据更新失败]")
				}

				return nil
			}

			err = facades.Orm().Query().UpdateOrCreate(&image, &models.Image{
				Hash: imageHash,
			}, &models.Image{
				Ban: ban,
			})
			if err != nil {
				facades.Log().With(map[string]any{
					"hash": hash,
					"err":  err.Error(),
				}).Error("图片审核[缓存数据创建失败]")
			}

			avatar.Ban = ban
		} else {
			avatar.Ban = image.Ban
		}

		err = facades.Orm().Query().Save(&avatar)
		if err != nil {
			facades.Log().With(map[string]any{
				"avatar": avatar,
				"err":    err.Error(),
			}).Error("图片审核[数据更新失败]")
			return nil
		}

		if avatar.Ban {
			cdn := packagecdn.NewCDN()
			cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
		}

		return nil
	}

	var avatar models.AppAvatar
	err := facades.Orm().Query().Where("avatar_hash", hash).First(&avatar)
	if err != nil {
		facades.Log().With(map[string]any{
			"hash": hash,
			"err":  err.Error(),
		}).Error("图片审核[数据查询失败]")
		return nil
	}
	if avatar.Checked || len(avatar.AvatarHash) == 0 {
		return nil
	}

	avatar.Checked = true
	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().With(map[string]any{
			"avatar": avatar,
			"err":    err.Error(),
		}).Error("图片审核[数据库更新失败]")
		return nil
	}

	// 检查WeAvatar APP头像是否存在
	var imageHash string
	exist := facades.Storage().Exists("upload/app/" + strconv.Itoa(int(avatar.AppID)) + "/" + hash[:2] + "/" + hash)
	if exist {
		fileString, fileErr := facades.Storage().Get("upload/app/" + strconv.Itoa(int(avatar.AppID)) + "/" + hash[:2] + "/" + hash)
		if fileErr != nil {
			facades.Log().With(map[string]any{
				"hash": hash,
				"err":  fileErr.Error(),
			}).Error("图片审核[文件读取失败]")
			return nil
		}

		imageHash = helper.MD5(fileString)
		err = facades.Storage().Put("checker/"+imageHash[:2]+"/"+imageHash, fileString)
		if err != nil {
			facades.Log().With(map[string]any{
				"avatarHash": hash,
				"imageHash":  imageHash,
				"err":        err.Error(),
			}).Error("图片审核[文件缓存失败]")
			return nil
		}
	} else {
		return nil
	}

	var image models.Image
	err = facades.Orm().Query().Where("hash", imageHash).FirstOrFail(&image)
	if err != nil {
		checker := imagecheck.NewChecker()
		ban, checkErr := checker.Check("https://weavatar.com/avatar/" + hash + ".png?appid=" + strconv.Itoa(int(avatar.AppID)) + "&s=600&d=404")
		if checkErr != nil {
			facades.Log().With(map[string]any{
				"hash": hash,
				"err":  checkErr.Error(),
			}).Error("图片审核[审核失败]")
			avatar.Checked = false
			err = facades.Orm().Query().Save(&avatar)
			if err != nil {
				facades.Log().With(map[string]any{
					"avatar": avatar,
					"err":    err.Error(),
				}).Error("图片审核[数据更新失败]")
			}

			return nil
		}
		err = facades.Orm().Query().UpdateOrCreate(&image, &models.Image{
			Hash: imageHash,
		}, &models.Image{
			Ban: ban,
		})
		if err != nil {
			facades.Log().With(map[string]any{
				"hash": hash,
				"err":  err.Error(),
			}).Error("图片审核[缓存数据创建失败]")
		}
		avatar.Ban = ban
	} else {
		avatar.Ban = image.Ban
	}

	err = facades.Orm().Query().Save(&avatar)
	if err != nil {
		facades.Log().With(map[string]any{
			"avatar": avatar,
			"err":    err.Error(),
		}).Error("图片审核[数据更新失败]")
		return nil
	}

	if avatar.Ban {
		cdn := packagecdn.NewCDN()
		cdn.RefreshUrl([]string{"weavatar.com/avatar/" + hash})
	}

	return nil
}
