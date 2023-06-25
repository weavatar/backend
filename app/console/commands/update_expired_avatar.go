package commands

import (
	"os"
	"path/filepath"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/spf13/cast"

	"weavatar/app/jobs"
)

type UpdateExpiredAvatar struct {
}

// Signature The name and signature of the console command.
func (receiver *UpdateExpiredAvatar) Signature() string {
	return "avatar:update-expired"
}

// Description The console command description.
func (receiver *UpdateExpiredAvatar) Description() string {
	return "更新过期的头像"
}

// Extend The console command extend.
func (receiver *UpdateExpiredAvatar) Extend() command.Extend {
	return command.Extend{
		Category: "avatar",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "force",
				Value:   "false",
				Aliases: []string{"f"},
				Usage:   "强制更新",
			},
		}}
}

// Handle Execute the console command.
func (receiver *UpdateExpiredAvatar) Handle(ctx console.Context) error {
	dir := facades.Storage().Path("cache")
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			facades.Log().Error("更新过期头像[无法获取文件信息] ", err.Error())
			return err
		}

		if !info.IsDir() {
			filename := filepath.Base(path)
			modTime := carbon.FromStdTime(info.ModTime())

			relPath, relErr := filepath.Rel(dir, path)
			if relErr != nil {
				facades.Log().Error("更新过期头像[无法获取相对路径] ", relErr.Error())
				return err
			}

			// 修改时间超过7天或者强制更新
			if (modTime.DiffAbsInSeconds(carbon.Now()) > 604800 || cast.ToBool(ctx.Option("force"))) && len(filename) == 32 {
				facades.Log().Info("更新过期头像[文件] " + filename)
				_ = facades.Queue().Job(&jobs.ProcessAvatarUpdate{}, []queue.Arg{
					{Type: "string", Value: filename},
					{Type: "string", Value: filepath.Join("cache", relPath)},
				}).Dispatch()
			}
		}

		return nil
	})

	if err != nil {
		facades.Log().Error("更新过期头像[遍历目录时出错] ", err.Error())
	}
	return nil
}
