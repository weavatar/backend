package controllers

import "github.com/goravel/framework/contracts/http"

func Success(ctx http.Context, data any) {
	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

func Error(ctx http.Context, code int, message any) {
	ctx.Response().Json(http.StatusOK, http.Json{
		"code":    code,
		"message": message,
	})
}

func Sanitize(ctx http.Context, request http.FormRequest) bool {
	errors, err := ctx.Request().ValidateRequest(request)
	if err != nil {
		Error(ctx, http.StatusUnprocessableEntity, err.Error())
		return false
	}
	if errors != nil {
		Error(ctx, http.StatusUnprocessableEntity, errors.One())
		return false
	}

	return true
}
