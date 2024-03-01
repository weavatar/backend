package jobs

import (
	"errors"
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
		return errors.New("图片审核[当前环境不允许执行此操作]")
	}

	if len(args) < 2 {
		facades.Log().With(map[string]any{
			"args": args,
		}).Warning("图片审核[队列参数不足]")
		return errors.New("图片审核[队列参数不足]")
	}

	hash, ok := args[0].(string)
	if !ok {
		facades.Log().With(map[string]any{
			"hash": hash,
		}).Warning("图片审核[队列参数断言失败]")
		return errors.New("图片审核[队列参数断言失败]")
	}

	appID, ok2 := args[1].(uint)
	if !ok2 {
		facades.Log().With(map[string]any{
			"appID": appID,
		}).Warning("图片审核[队列参数断言失败]")
		return errors.New("图片审核[队列参数断言失败]")
	}

	var imageHash string

	if appID == 0 {
		// 默认头像
		if exist := facades.Storage().Exists("upload/default/" + hash[:2] + "/" + hash); exist {
			fileString, err := facades.Storage().Get("upload/default/" + hash[:2] + "/" + hash)
			if err != nil {
				facades.Log().With(map[string]any{
					"hash": hash,
					"err":  err.Error(),
				}).Warning("图片审核[文件读取失败]")
				return err
			}

			imageHash = helper.MD5(fileString)
			err = facades.Storage().Put("checker/"+imageHash[:2]+"/"+imageHash, fileString)
			if err != nil {
				facades.Log().With(map[string]any{
					"avatarHash": hash,
					"imageHash":  imageHash,
					"err":        err.Error(),
				}).Warning("图片审核[文件缓存失败]")
				return err
			}
		} else {
			client := req.C()
			client.SetTimeout(5 * time.Second)
			client.SetCommonRetryCount(2)
			client.ImpersonateSafari()

			resp, err := client.R().Get("http://proxy.server/https://0.gravatar.com/avatar/" + hash + ".png?s=600&r=g&d=404")
			if err != nil || !resp.IsSuccessState() {
				facades.Log().With(map[string]any{
					"hash":     hash,
					"response": resp.String(),
				}).Warning("图片审核[Gravatar头像下载失败]")
				return err
			}

			imageHash = helper.MD5(resp.String())
			err = facades.Storage().Put("checker/"+imageHash[:2]+"/"+imageHash, resp.String())
			if err != nil {
				facades.Log().With(map[string]any{
					"avatarHash": hash,
					"imageHash":  imageHash,
					"err":        err.Error(),
				}).Warning("图片审核[文件缓存失败]")
				return err
			}
		}
	} else {
		// APP头像
		var avatar models.AppAvatar
		err := facades.Orm().Query().Where("avatar_sha256", hash).First(&avatar)
		if err != nil {
			facades.Log().With(map[string]any{
				"avatarSHA256": hash,
				"appID":        appID,
				"err":          err.Error(),
			}).Warning("图片审核[数据查询失败]")
			return err
		}
		if exist := facades.Storage().Exists("upload/app/" + strconv.Itoa(int(avatar.AppID)) + "/" + hash[:2] + "/" + hash); exist {
			fileString, fileErr := facades.Storage().Get("upload/app/" + strconv.Itoa(int(avatar.AppID)) + "/" + hash[:2] + "/" + hash)
			if fileErr != nil {
				facades.Log().With(map[string]any{
					"avatarSHA256": hash,
					"appID":        appID,
					"err":          fileErr.Error(),
				}).Warning("图片审核[文件读取失败]")
				return fileErr
			}

			imageHash = helper.MD5(fileString)
			err = facades.Storage().Put("checker/"+imageHash[:2]+"/"+imageHash, fileString)
			if err != nil {
				facades.Log().With(map[string]any{
					"avatarSHA256": hash,
					"imageHash":    imageHash,
					"appID":        appID,
					"err":          err.Error(),
				}).Warning("图片审核[文件缓存失败]")
				return err
			}
		} else {
			return errors.New("图片审核[APP头像不存在]")
		}
	}

	var image models.Image
	if err := facades.Orm().Query().Where("hash", imageHash).FirstOrFail(&image); err != nil {
		checker := imagecheck.NewChecker()
		ban, checkErr := checker.Check("https://weavatar.com/avatar/" + hash + ".png?s=600&d=404")
		if checkErr != nil {
			facades.Log().With(map[string]any{
				"hash":      hash,
				"imageHash": imageHash,
				"err":       checkErr.Error(),
			}).Warning("图片审核[审核失败]")
			return err
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
			}).Warning("图片审核[数据创建失败]")
		}
	}

	if image.Ban {
		cdn := packagecdn.NewCDN()
		cdn.RefreshPath([]string{"weavatar.com/avatar/"})
	}

	return nil
}
