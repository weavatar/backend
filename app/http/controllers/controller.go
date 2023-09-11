package controllers

import "github.com/goravel/framework/contracts/http"

func Success(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

func Error(ctx http.Context, code int, message any) http.Response {
	return ctx.Response().Json(http.StatusOK, http.Json{
		"code":    code,
		"message": message,
	})
}

func Sanitize(ctx http.Context, request http.FormRequest) http.Response {
	errors, err := ctx.Request().ValidateRequest(request)
	if err != nil {
		return Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	if errors != nil {
		return Error(ctx, http.StatusUnprocessableEntity, errors.One())
	}

	return nil
}
