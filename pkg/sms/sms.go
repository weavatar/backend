// Package sms 发送短信
package sms

import (
	"sync"

	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
)

// Message 短信结构体
type Message struct {
	Data map[string]string

	Content string
}

type SMS struct {
	Driver Driver
	Config map[string]string
}

var once sync.Once

var internalSMS *SMS

// NewSMS 创建短信实例
func NewSMS() *SMS {
	driver := facades.Config().Get("sms.driver")
	config := make(map[string]string)
	switch driver {
	case "aliyun":
		config = cast.ToStringMapString(facades.Config().Get("sms.aliyun"))
		config["expire_time"] = facades.Config().GetString("verifycode.expire_time")
		once.Do(func() {
			internalSMS = &SMS{
				Driver: &Tencent{},
				Config: config,
			}
		})
	case "tencent":
		config = cast.ToStringMapString(facades.Config().Get("sms.tencent"))
		config["expire_time"] = facades.Config().GetString("verifycode.expire_time")
		once.Do(func() {
			internalSMS = &SMS{
				Driver: &Tencent{},
				Config: config,
			}
		})
	}

	return internalSMS
}

// Send 发送短信
func (s *SMS) Send(phone string, message Message) error {
	return s.Driver.Send(phone, message, s.Config)
}
