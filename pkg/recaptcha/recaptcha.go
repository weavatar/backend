package recaptcha

import (
	"sync"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"
)

type Recaptcha struct {
	secret string
}

type recaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

var once sync.Once

var internalRecaptcha *Recaptcha

// NewRecaptcha 创建 Recaptcha 实例
func NewRecaptcha() *Recaptcha {
	once.Do(func() {
		internalRecaptcha = &Recaptcha{
			secret: facades.Config().GetString("recaptcha.secret"),
		}
	})
	return internalRecaptcha
}

// Confirm 验证 Recaptcha
func (re *Recaptcha) Confirm(remoteIp, response, action string) bool {
	resp, err := re.check(remoteIp, response)
	if err != nil {
		return false
	}

	return resp.Success && resp.Score >= 0.7 && resp.Action == action
}

func (re *Recaptcha) check(remoteIp, response string) (r recaptchaResponse, err error) {
	client := req.C()
	var resp recaptchaResponse

	_, err = client.R().SetFormData(map[string]string{
		"secret":   re.secret,
		"remoteip": remoteIp,
		"response": response,
	}).SetSuccessResult(&resp).SetErrorResult(&resp).Post("https://recaptcha.net/recaptcha/api/siteverify")
	if err != nil {
		facades.Log().Error("[Recaptcha] ", " HTTP请求失败 "+err.Error())
		return resp, err
	}

	return resp, nil
}
