package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateProfileRequest struct {
	Nickname string `form:"nickname" json:"nickname"`
	Avatar   string `form:"avatar" json:"avatar"`
}

func (r *UpdateProfileRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateProfileRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"nickname": "required|string",
		"avatar":   "full_url",
	}
}

func (r *UpdateProfileRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"nickname.required": "昵称不能为空",
		"nickname.string":   "昵称必须是字符串",
		"avatar.full_url":   "头像必须是一个完整的URL",
	}
}

func (r *UpdateProfileRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateProfileRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	avatar := ctx.Request().Input("avatar")
	if len(avatar) == 0 {
		avatar = "https://cdn.goravel.net/avatars/default.png"
	}
	_ = data.Set("avatar", avatar)
	return nil
}
