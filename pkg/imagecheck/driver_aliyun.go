package imagecheck

import (
	"fmt"
	"net/http"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/goravel/framework/support/json"
)

type Aliyun struct {
	AccessKeyId     string
	AccessKeySecret string
}

// Check 检查图片是否违规 true: 违规 false: 未违规
func (a *Aliyun) Check(url string) (bool, error) {
	client, err := a.createClient(tea.String(a.AccessKeyId), tea.String(a.AccessKeySecret), "shanghai")
	if err != nil {
		return false, err
	}

	parameters, err := json.MarshalString(map[string]string{
		"imageUrl": url,
	})
	if err != nil {
		return false, err
	}

	imageModerationRequest := &green20220302.ImageModerationRequest{
		Service:           tea.String("baselineCheck"),
		ServiceParameters: tea.String(parameters),
	}
	runtime := &util.RuntimeOptions{
		Autoretry:   tea.Bool(true),
		MaxAttempts: tea.Int(3),
		IgnoreSSL:   tea.Bool(true),
	}
	response, _err := client.ImageModerationWithOptions(imageModerationRequest, runtime)

	// 自动切换地域
	flag := false
	if _err != nil {
		var err = &tea.SDKError{}
		if _t, ok := _err.(*tea.SDKError); ok {
			err = _t
			// 系统异常，切换到下个地域调用。
			if *err.StatusCode == 500 {
				flag = true
			}
		}
	}
	if response == nil || *response.StatusCode == 500 || *response.Body.Code == 500 {
		flag = true
	}
	if flag {
		client, err := a.createClient(tea.String(a.AccessKeyId), tea.String(a.AccessKeySecret), "beijing")
		if err != nil {
			return false, err
		}
		response, _err = client.ImageModerationWithOptions(imageModerationRequest, runtime)
		if _err != nil {
			return false, _err
		}
	}

	if response != nil {
		statusCode := tea.IntValue(tea.ToInt(response.StatusCode))
		body := response.Body
		imageModerationResponseData := body.Data
		if statusCode == http.StatusOK {
			if tea.IntValue(tea.ToInt(body.Code)) == 200 {
				result := imageModerationResponseData.Result
				for i := 0; i < len(result); i++ {
					if tea.Float32Value(result[i].Confidence) > 80 {
						return true, nil
					}
				}
			} else {
				return false, fmt.Errorf("审核调用失败 url:%s, httpCode:%d, requestId:%s, msg:%s", url, statusCode, tea.StringValue(body.RequestId), tea.StringValue(body.Msg))
			}
		} else {
			return false, fmt.Errorf("审核调用失败 url:%s, httpCode:%d, requestId:%s, msg:%s", url, statusCode, tea.StringValue(body.RequestId), tea.StringValue(body.Msg))
		}
	}

	return false, nil
}

func (a *Aliyun) createClient(accessKeyId *string, accessKeySecret *string, endpoint string) (_result *green20220302.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	if endpoint == "shanghai" {
		config.RegionId = tea.String("cn-shanghai")
		config.Endpoint = tea.String("green-cip.cn-shanghai.aliyuncs.com")
	}
	if endpoint == "beijing" {
		config.RegionId = tea.String("cn-beijing")
		config.Endpoint = tea.String("green-cip.cn-beijing.aliyuncs.com")
	}

	_result, _err = green20220302.NewClient(config)
	return _result, _err
}
