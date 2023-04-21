package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type OauthCallbackRequest struct {
	Code  string `form:"code" json:"code"`
	State string `form:"state" json:"state"`
}

func (r *OauthCallbackRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *OauthCallbackRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"code":  "required|string",
		"state": "required|string",
	}
}

func (r *OauthCallbackRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"code.required":  "授权码不能为空",
		"code.string":    "授权码必须是字符串",
		"state.required": "状态码不能为空",
		"state.string":   "状态码必须是字符串",
	}
}

func (r *OauthCallbackRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *OauthCallbackRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
