package imagecheck

import (
	"fmt"
	"net/http"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green20220302 "github.com/alibabacloud-go/green-20220302/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/bytedance/sonic"
)

type Checker struct {
	accessKeyId, accessKeySecret string
}

// once 单例模式
var once sync.Once

// internal 内部使用的对象
var internal *Checker

func NewCreator(accessKeyId, accessKeySecret string) *Checker {
	once.Do(func() {
		internal = &Checker{accessKeyId: accessKeyId, accessKeySecret: accessKeySecret}
	})
	return internal
}

// Check 检查图片是否违规
func (c *Checker) Check(url string) (bool, error) {
	client, err := createClient(tea.String(c.accessKeyId), tea.String(c.accessKeySecret), "shanghai")
	if err != nil {
		return false, err
	}

	parameters, err := sonic.MarshalString(map[string]string{
		"imageUrl": url,
		"dataId":   url,
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

	// 系统异常，切换到下个地域调用。
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
		client, err := createClient(tea.String(c.accessKeyId), tea.String(c.accessKeySecret), "beijing")
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
				return false, fmt.Errorf("审核调用失败 httpCode:%d, requestId:%s, msg:%s", statusCode, tea.StringValue(body.RequestId), tea.StringValue(body.Msg))
			}
		} else {
			return false, fmt.Errorf("审核调用失败 httpCode:%d, requestId:%s, msg:%s", statusCode, tea.StringValue(body.RequestId), tea.StringValue(body.Msg))
		}
	}

	return false, nil
}

func createClient(accessKeyId *string, accessKeySecret *string, endpoint string) (_result *green20220302.Client, _err error) {
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
	_result = &green20220302.Client{}
	_result, _err = green20220302.NewClient(config)
	return _result, _err
}
