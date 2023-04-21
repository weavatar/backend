package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type DeleteAvatarRequest struct {
	Name string `form:"name" json:"name"`
}

func (r *DeleteAvatarRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *DeleteAvatarRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DeleteAvatarRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DeleteAvatarRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DeleteAvatarRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
