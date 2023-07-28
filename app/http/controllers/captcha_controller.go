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
func (r *CaptchaController) Image(ctx http.Context) {
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	if err != nil {
		facades.Log().Error("[CaptchaController][Image] 生成图片验证码失败 ", err)
		Error(ctx, http.StatusInternalServerError, "系统内部错误")
		return
	}

	Success(ctx, http.Json{
		"captcha_id": id,
		"captcha":    b64s,
	})
}

// Sms 获取短信验证码
func (r *CaptchaController) Sms(ctx http.Context) {
	var smsRequest requests.SmsRequest
	errors, err := ctx.Request().ValidateRequest(&smsRequest)
	if err != nil {
		Error(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if errors != nil {
		Error(ctx, http.StatusUnprocessableEntity, errors.All())
		return
	}

	if ok := verifycode.NewVerifyCode().SendSMS(smsRequest.Phone, smsRequest.UseFor); !ok {
		facades.Log().Error("[CaptchaController][Sms] 发送短信验证码失败 ", err)
		Error(ctx, http.StatusInternalServerError, "发送失败")
		return
	}

	Success(ctx, nil)
}

// Email 获取邮箱验证码
func (r *CaptchaController) Email(ctx http.Context) {
	var emailRequest requests.EmailRequest
	errors, err := ctx.Request().ValidateRequest(&emailRequest)
	if err != nil {
		Error(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if errors != nil {
		Error(ctx, http.StatusUnprocessableEntity, errors.All())
		return
	}

	if ok := verifycode.NewVerifyCode().SendEmail(emailRequest.Email, emailRequest.UseFor); !ok {
		facades.Log().Error("[CaptchaController][Email] 发送邮箱验证码失败 ", err)
		Error(ctx, http.StatusInternalServerError, "发送失败")
		return
	}

	Success(ctx, nil)
}
