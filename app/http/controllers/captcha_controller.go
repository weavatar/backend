package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	requests "weavatar/app/http/requests/captcha"
	"weavatar/pkg/captcha"
	"weavatar/pkg/verifycode"
)

type CaptchaController struct {
	// Dependent services
}

func NewCaptchaController() *CaptchaController {
	return &CaptchaController{
		// Inject services
	}
}

// Image 获取图片验证码
func (r *CaptchaController) Image(ctx http.Context) http.Response {
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	if err != nil {
		facades.Log().Error("[CaptchaController][Image] 生成图片验证码失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "系统内部错误")
	}

	return Success(ctx, http.Json{
		"captcha_id": id,
		"captcha":    b64s,
	})
}

// Sms 获取短信验证码
func (r *CaptchaController) Sms(ctx http.Context) http.Response {
	var smsRequest requests.SmsRequest
	sanitize := Sanitize(ctx, &smsRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := verifycode.NewVerifyCode().SendSMS(smsRequest.Phone, smsRequest.UseFor); err != nil {
		facades.Log().Error("[CaptchaController][Sms] 发送短信验证码失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "发送失败")
	}

	return Success(ctx, nil)
}

// Email 获取邮箱验证码
func (r *CaptchaController) Email(ctx http.Context) http.Response {
	var emailRequest requests.EmailRequest
	sanitize := Sanitize(ctx, &emailRequest)
	if sanitize != nil {
		return sanitize
	}

	if err := verifycode.NewVerifyCode().SendEmail(emailRequest.Email, emailRequest.UseFor); err != nil {
		facades.Log().Error("[CaptchaController][Email] 发送邮箱验证码失败 ", err.Error())
		return Error(ctx, http.StatusInternalServerError, "发送失败")
	}

	return Success(ctx, nil)
}
