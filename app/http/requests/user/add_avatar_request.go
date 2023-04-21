package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type AddAvatarRequest struct {
	Name string `form:"name" json:"name"`
}

func (r *AddAvatarRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *AddAvatarRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *AddAvatarRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *AddAvatarRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *AddAvatarRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
