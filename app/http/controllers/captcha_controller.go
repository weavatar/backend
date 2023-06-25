package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	requests "weavatar/app/http/requests/captcha"
	"weavatar/packages/captcha"
	"weavatar/packages/verifycode"
)

type CaptchaController struct {
	//Dependent services
}

func NewCaptchaController() *CaptchaController {
	return &CaptchaController{
		//Inject services
	}
}

// Image 获取图片验证码
func (r *CaptchaController) Image(ctx http.Context) {
	// 生成验证码
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	if err != nil {
		facades.Log().Error("[CaptchaController][Image] 生成图片验证码失败 ", err)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "生成失败",
		})
		return
	}
	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "生成成功",
		"data": http.Json{
			"captcha_id": id,
			"captcha":    b64s,
		},
	})
}

// Sms 获取短信验证码
func (r *CaptchaController) Sms(ctx http.Context) {
	var smsRequest requests.SmsRequest
	errors, err := ctx.Request().ValidateRequest(&smsRequest)
	if err != nil {
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": err.Error(),
		})
		return
	}
	if errors != nil {
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": errors.All(),
		})
		return
	}

	if ok := verifycode.NewVerifyCode().SendSMS(smsRequest.Phone, smsRequest.UseFor); !ok {
		facades.Log().Error("[CaptchaController][Sms] 发送短信验证码失败 ", err)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "发送失败",
		})
		return
	}

	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "发送成功",
	})
}

// Email 获取邮箱验证码
func (r *CaptchaController) Email(ctx http.Context) {
	var emailRequest requests.EmailRequest
	errors, err := ctx.Request().ValidateRequest(&emailRequest)
	if err != nil {
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": err.Error(),
		})
		return
	}
	if errors != nil {
		ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"code":    422,
			"message": errors.All(),
		})
		return
	}

	if ok := verifycode.NewVerifyCode().SendEmail(emailRequest.Email, emailRequest.UseFor); !ok {
		facades.Log().Error("[CaptchaController][Email] 发送邮箱验证码失败 ", err)
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"code":    500,
			"message": "发送失败",
		})
		return
	}

	ctx.Response().Success().Json(http.Json{
		"code":    0,
		"message": "发送成功",
	})
}
