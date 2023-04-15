// Package verifycode 用以发送手机验证码和邮箱验证码
package verifycode

import (
	"fmt"
	"strings"
	"sync"

	"github.com/goravel/framework/contracts/mail"
	"github.com/goravel/framework/facades"

	"weavatar/packages/helpers"
	"weavatar/packages/sms"
)

type VerifyCode struct {
	Store Store
}

var once sync.Once
var internalVerifyCode *VerifyCode

// NewVerifyCode 单例模式获取
func NewVerifyCode() *VerifyCode {
	once.Do(func() {
		internalVerifyCode = &VerifyCode{
			Store: &CacheStore{
				// 增加前缀保持数据库整洁，出问题调试时也方便
				KeyPrefix: "verify_code:",
			},
		}
	})

	return internalVerifyCode
}

// SendSMS 发送短信验证码
func (vc *VerifyCode) SendSMS(phone string, forName string) bool {

	// 生成验证码
	code := vc.generateVerifyCode(phone, forName)

	// 方便本地和 API 自动测试
	if facades.Config.GetBool("app.debug") && strings.HasPrefix(phone, facades.Config.GetString("verifycode.debug_phone_prefix")) {
		return true
	}

	// 发送短信
	return sms.NewSMS().Send(phone, sms.Message{
		Data: map[string]string{"code": code},
	})
}

// SendEmail 发送邮件验证码
func (vc *VerifyCode) SendEmail(email string, forName string) bool {

	// 生成验证码
	code := vc.generateVerifyCode(email, forName)

	// 方便本地和 API 自动测试
	if facades.Config.GetBool("app.debug") && strings.HasSuffix(email, facades.Config.GetString("verifycode.debug_email_suffix")) {
		return true
	}

	content := fmt.Sprintf("<h1>您的 Email 验证码是 %v </h1>", code)
	// 发送邮件
	err := facades.Mail.To([]string{email}).
		Content(mail.Content{Subject: facades.Config.GetString("app.name") + " - 验证码", Html: content}).
		Send()

	return err == nil
}

// Check 检查用户提交的验证码是否正确，key 可以是手机号或者 Email
func (vc *VerifyCode) Check(key string, answer string, useFor string, clear bool) bool {

	return vc.Store.Verify(useFor+":"+key, answer, clear)
}

// generateVerifyCode 生成验证码，并放置于 Redis 中
func (vc *VerifyCode) generateVerifyCode(key string, forName string) string {

	// 生成随机码
	code := helpers.RandomNumber(facades.Config.GetInt("verifycode.code_length"))

	// 为方便开发，本地环境使用固定验证码
	if facades.Config.GetBool("app.debug") {
		code = facades.Config.GetString("verifycode.debug_code")
	}

	// 存储验证码
	vc.Store.Set(forName+":"+key, code)
	return code
}
