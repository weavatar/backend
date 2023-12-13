// Package captcha 处理图片验证码逻辑
package captcha

import (
	"sync"

	"github.com/goravel/framework/facades"
	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

var once sync.Once

var internalCaptcha *Captcha

// NewCaptcha 创建图片验证码实例
func NewCaptcha() *Captcha {
	once.Do(func() {
		internalCaptcha = &Captcha{}
		store := CacheStore{
			KeyPrefix: facades.Config().GetString("app.name") + ":captcha:",
		}

		driver := base64Captcha.NewDriverDigit(
			facades.Config().GetInt("captcha.height"),         // 宽
			facades.Config().GetInt("captcha.width"),          // 高
			facades.Config().GetInt("captcha.length"),         // 长度
			facades.Config().Get("captcha.maxskew").(float64), // 数字的最大倾斜角度
			facades.Config().GetInt("captcha.dotcount"),       // 图片背景里的混淆点数量
		)

		internalCaptcha.Base64Captcha = base64Captcha.NewCaptcha(driver, &store)
	})

	return internalCaptcha
}

// GenerateCaptcha 生成图片验证码
func (c *Captcha) GenerateCaptcha() (id string, b64s string, err error) {
	id, b64s, _, err = c.Base64Captcha.Generate()
	return id, b64s, err
}

// VerifyCaptcha 验证验证码是否正确
func (c *Captcha) VerifyCaptcha(id string, answer string, clear bool) (match bool) {
	if facades.Config().GetBool("app.debug") && id == facades.Config().GetString("captcha.testing_key") {
		return true
	}

	return c.Base64Captcha.Verify(id, answer, clear)
}
