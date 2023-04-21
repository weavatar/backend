package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateNicknameRequest struct {
	Name string `form:"name" json:"name"`
}

func (r *UpdateNicknameRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateNicknameRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateNicknameRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateNicknameRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateNicknameRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
