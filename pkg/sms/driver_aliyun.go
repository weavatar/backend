package sms

import (
	"errors"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/goravel/framework/support/json"
)

type Aliyun struct{}

func (a *Aliyun) Send(phone string, message Message, config map[string]string) error {
	client, err := CreateClient(tea.String(config["access_key_id"]), tea.String(config["access_key_secret"]))
	if err != nil {
		return err
	}

	param, err := json.Marshal(message.Data)
	if err != nil {
		return err
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
		var _t *tea.SDKError
		if errors.As(err, &_t) {
			errs = _t
		}

		return errs
	}

	if tea.StringValue(_result.Body.Message) != "OK" {
		return errors.New("短信发送失败: " + tea.StringValue(_result.Body.Message))
	}

	return nil
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
