package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateNicknameRequest struct {
	Nickname string `form:"nickname" json:"nickname"`
}

func (r *UpdateNicknameRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateNicknameRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"nickname": "required|string",
	}
}

func (r *UpdateNicknameRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"nickname.required": "昵称不能为空",
		"nickname.string":   "昵称必须是字符串",
	}
}

func (r *UpdateNicknameRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateNicknameRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
