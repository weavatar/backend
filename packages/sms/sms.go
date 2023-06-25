// Package sms 发送短信
package sms

import (
	"sync"

	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
)

// Message 短信的结构体
type Message struct {
	Data map[string]string

	Content string
}

// SMS 发送短信操作类
type SMS struct {
	Driver Driver
	Config map[string]string
}

// once 单例模式
var once sync.Once

// internalSMS 内部使用的 SMS 对象
var internalSMS *SMS

// NewSMS 单例模式获取
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

func (sms *SMS) Send(phone string, message Message) bool {
	return sms.Driver.Send(phone, message, sms.Config)
}
