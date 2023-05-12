package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateInfoRequest struct {
	Nickname string `form:"nickname" json:"nickname"`
	Avatar   string `form:"avatar" json:"avatar"`
}

func (r *UpdateInfoRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateInfoRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"nickname": "required|string",
		"avatar":   "full_url",
	}
}

func (r *UpdateInfoRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"nickname.required": "昵称不能为空",
		"nickname.string":   "昵称必须是字符串",
		"avatar.full_url":   "头像必须是一个完整的URL",
	}
}

func (r *UpdateInfoRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateInfoRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	avatar, exist := data.Get("avatar")
	if !exist {
		err := data.Set("avatar", "https://weavatar.com/avatar/?d=mp")
		if err != nil {
			return err
		}
	} else {
		check, ok := avatar.(string)
		if !ok || len(check) == 0 {
			err := data.Set("avatar", "https://weavatar.com/avatar/?d=mp")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
