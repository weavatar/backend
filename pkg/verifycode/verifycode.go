// Package verifycode 用以发送手机验证码和邮箱验证码
package verifycode

import (
	"fmt"
	"sync"

	"github.com/goravel/framework/contracts/mail"
	"github.com/goravel/framework/facades"

	"weavatar/pkg/helper"
	"weavatar/pkg/sms"
)

type VerifyCode struct {
	Store Store
}

var once sync.Once
var internalVerifyCode *VerifyCode

// NewVerifyCode 创建验证码实例
func NewVerifyCode() *VerifyCode {
	once.Do(func() {
		internalVerifyCode = &VerifyCode{
			Store: &CacheStore{
				KeyPrefix: "verify_code:",
			},
		}
	})

	return internalVerifyCode
}

// SendSMS 发送短信验证码
func (vc *VerifyCode) SendSMS(phone string, useFor string) bool {
	code := vc.generateVerifyCode(phone, useFor)

	if facades.Config().GetBool("app.debug") {
		return true
	}

	return sms.NewSMS().Send(phone, sms.Message{
		Data: map[string]string{"code": code},
	})
}

// SendEmail 发送邮件验证码
func (vc *VerifyCode) SendEmail(email string, useFor string) bool {
	code := vc.generateVerifyCode(email, useFor)

	if facades.Config().GetBool("app.debug") {
		return true
	}

	content := fmt.Sprintf("<h1>您的 Email 验证码是 %v </h1>", code)

	err := facades.Mail().To([]string{email}).
		Content(mail.Content{Subject: facades.Config().GetString("app.name") + " - 验证码", Html: content}).
		Send()

	return err == nil
}

// Check 检查用户提交的验证码是否正确
func (vc *VerifyCode) Check(key string, answer string, useFor string, clear bool) bool {

	return vc.Store.Verify(useFor+":"+key, answer, clear)
}

// generateVerifyCode 生成验证码
func (vc *VerifyCode) generateVerifyCode(key string, useFor string) string {
	code := helper.RandomNumber(facades.Config().GetInt("verifycode.code_length"))

	if facades.Config().GetBool("app.debug") {
		code = facades.Config().GetString("verifycode.debug_code")
	}

	vc.Store.Set(useFor+":"+key, code)
	return code
}
