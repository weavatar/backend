package sms

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/bytedance/sonic"

	"github.com/goravel/framework/facades"
)

type Aliyun struct{}

func (a *Aliyun) Send(phone string, message Message, config map[string]string) bool {
	client, err := CreateClient(tea.String(config["access_key_id"]), tea.String(config["access_key_secret"]))
	if err != nil {
		facades.Log().Error("短信[阿里云] ", " 解析绑定错误 ", err.Error())
		return false
	}

	param, err := sonic.Marshal(message.Data)
	if err != nil {
		facades.Log().Error("短信[阿里云] ", " 短信模板参数解析错误 ", err.Error())
		return false
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(config["sign_name"]),
		TemplateCode:  tea.String(config["template_code"]),
		PhoneNumbers:  tea.String(phone),
		TemplateParam: tea.String(string(param)),
	}

	runtime := &util.RuntimeOptions{}

	_result, err := client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		var errs = &tea.SDKError{}
		if _t, ok := err.(*tea.SDKError); ok {
			errs = _t
		} else {
			errs.Message = tea.String(err.Error())
		}

		var r dysmsapi20170525.SendSmsResponseBody
		err = sonic.Unmarshal([]byte(*errs.Data), &r)
		if err != nil {
			facades.Log().Error("短信[阿里云] ", " 解析JSON失败 ", errs)
		}

		return false
	}

	if tea.StringValue(_result.Body.Message) != "OK" {
		facades.Log().Error("短信[阿里云] ", " 发送失败 ", _result.Body)
		return false
	}

	return true
}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}

	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}
