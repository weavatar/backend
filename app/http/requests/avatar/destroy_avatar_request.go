package avatar

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type DestroyAvatarRequest struct {
	Hash string `form:"hash" json:"hash"`
}

func (r *DestroyAvatarRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *DestroyAvatarRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"hash": "required|string|len:32",
	}
}

func (r *DestroyAvatarRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"hash.required": "头像哈希不能为空",
		"hash.string":   "头像哈希必须为字符串",
		"hash.len":      "头像哈希长度必须为 32 位",
	}
}

func (r *DestroyAvatarRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DestroyAvatarRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
