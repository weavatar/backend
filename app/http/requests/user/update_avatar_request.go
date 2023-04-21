package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateAvatarRequest struct {
	Name string `form:"name" json:"name"`
}

func (r *UpdateAvatarRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateAvatarRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateAvatarRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateAvatarRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateAvatarRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
