package sms

import (
	"errors"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Tencent struct{}

func (s *Tencent) Send(phone string, message Message, config map[string]string) error {
	credential := common.NewCredential(
		config["access_key"],
		config["secret_key"],
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := tencentsms.NewClient(credential, "ap-guangzhou", cpf)

	request := tencentsms.NewSendSmsRequest()
	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	request.SignName = common.StringPtr(config["sign_name"])
	request.TemplateId = common.StringPtr(config["template_code"])
	request.TemplateParamSet = common.StringPtrs([]string{message.Data["code"], config["expire_time"]})
	request.SmsSdkAppId = common.StringPtr(config["sdk_app_id"])

	response, err := client.SendSms(request)

	if err != nil {
		return err
	}

	statusSet := response.Response.SendStatusSet
	code := *statusSet[0].Code
	if code == "Ok" {
		return nil
	} else {
		return errors.New("短信发送失败: " + response.ToJsonString())
	}

}
