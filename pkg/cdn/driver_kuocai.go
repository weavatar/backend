package cdn

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/imroc/req/v3"
)

func init() {
	if !driverInUse("kuocai") {
		return
	}

	register(&KuoCai{
		UserName: facades.Config().GetString("cdn.kuocai.username"),
		PassWord: facades.Config().GetString("cdn.kuocai.password"),
	})
}

type KuoCai struct {
	UserName, PassWord string
}

type KuoCaiCommonResponse struct {
	Code                string `json:"code"`
	Message             string `json:"message"`
	Data                string `json:"data"`
	Success             bool   `json:"success"`
	SuccessWithDateResp bool   `json:"successWithDateResp"`
}

type KuoCaiUsageResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		VisitsSummary struct {
			HitFlux string `json:"hit_flux"`
			ReqNum  int    `json:"req_num"`
			HitNum  int    `json:"hit_num"`
		} `json:"visits_summary"`
		VisitsDetail struct {
			HitFlux struct {
				Unit string    `json:"unit"`
				Data []float64 `json:"data"`
			} `json:"hit_flux"`
			ReqNum []int `json:"req_num"`
			HitNum []int `json:"hit_num"`
		} `json:"visits_detail"`
		Labels []string `json:"labels"`
	} `json:"data"`
	Success             bool `json:"success"`
	SuccessWithDateResp bool `json:"successWithDateResp"`
}

// RefreshUrl 刷新URL
func (r *KuoCai) RefreshUrl(urls []string) error {
	client, err := r.login()
	if err != nil {
		return err
	}

	// 提交刷新请求
	var refreshResponse KuoCaiCommonResponse
	_, err = client.R().SetFormDataFromValues(url.Values{
		"urls[]": urls,
		"type":   {"file"},
	}).
		SetSuccessResult(&refreshResponse).
		Post("https://kuocai.cn/CdnDomainCache/submitCacheRefresh")
	if err != nil {
		return err
	}

	if refreshResponse.Success && refreshResponse.Code == "SUCCESS" {
		return nil
	}

	return fmt.Errorf("URL刷新失败: %s", refreshResponse.Message)
}

// RefreshPath 刷新路径
func (r *KuoCai) RefreshPath(paths []string) error {
	client, err := r.login()
	if err != nil {
		return err
	}

	// 提交刷新请求
	var refreshResponse KuoCaiCommonResponse
	_, err = client.R().SetFormDataFromValues(url.Values{
		"urls[]": paths,
		"type":   {"directory"},
	}).
		SetSuccessResult(&refreshResponse).
		Post("https://kuocai.cn/CdnDomainCache/submitCacheRefresh")
	if err != nil {
		return err
	}

	if refreshResponse.Success && refreshResponse.Code == "SUCCESS" {
		return nil
	}

	return fmt.Errorf("路径刷新失败: %s", refreshResponse.Message)
}

// GetUsage 获取使用量
func (r *KuoCai) GetUsage(domain string, startTime, endTime carbon.Carbon) (uint, error) {
	client, err := r.login()
	if err != nil {
		return 0, err
	}

	var request = map[string]string{
		"startTime": startTime.ToDateTimeString(),
		"endTime":   endTime.ToDateTimeString(),
		"domains":   domain,
		"type":      "Visits",
	}
	var usageResponse KuoCaiUsageResponse
	_, err = client.R().
		SetQueryParams(request).
		SetSuccessResult(&usageResponse).
		Get("https://kuocai.cn/CdnDomainStatistics/queryStatistics")
	if err != nil {
		return 0, err
	}

	if usageResponse.Success && usageResponse.Code == "SUCCESS" {
		return uint(usageResponse.Data.VisitsSummary.ReqNum), nil
	}

	return 0, fmt.Errorf("获取用量失败: %s", usageResponse.Message)
}

// login 登录平台
func (r *KuoCai) login() (*req.Client, error) {
	client := req.C()
	client.SetTimeout(10 * time.Second)
	client.SetCommonRetryCount(2)
	client.ImpersonateSafari()

	// 提交登录请求
	loginURL := "https://kuocai.cn/login/loginUser"
	loginParams := map[string]string{
		"userAccount": r.UserName,
		"userPwd":     r.PassWord,
		"remember":    "true",
	}
	var loginResponse KuoCaiCommonResponse
	_, err := client.R().SetFormData(loginParams).SetSuccessResult(&loginResponse).Post(loginURL)
	if err != nil {
		return nil, err
	}
	if !loginResponse.Success || loginResponse.Code != "SUCCESS" {
		return nil, errors.New(loginResponse.Message)
	}

	client.SetCommonCookies(&http.Cookie{
		Name:  "kuocai_cdn_token",
		Value: loginResponse.Data,
	})

	return client, nil
}
