package cdn

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/imroc/req/v3"

	"weavatar/pkg/helper"
)

type YunDun struct {
	UserName, PassWord string
}

type YunDunRefreshResponse struct {
	Status struct {
		Code                        int    `json:"code"`
		Message                     string `json:"message"`
		CreateAt                    string `json:"create_at"`
		ApiTimeConsuming            string `json:"api_time_consuming"`
		FunctionTimeConsuming       string `json:"function_time_consuming"`
		DispatchBeforeTimeConsuming string `json:"dispatch_before_time_consuming"`
	} `json:"status"`
	Data struct {
		Wholesite  []interface{} `json:"wholesite"`
		Specialurl []string      `json:"specialurl"`
		Specialdir []interface{} `json:"specialdir"`
		RequestId  string        `json:"request_id"`
	} `json:"data"`
}

type YunDunUsageRequest struct {
	Router               string   `json:"router"`
	StartTime            string   `json:"start_time"`
	EndTime              string   `json:"end_time"`
	Nodes                []string `json:"nodes"`
	GroupId              []string `json:"group_id"`
	SubDomain            []string `json:"sub_domain"`
	SubDomainsAndNodeIps struct {
	} `json:"sub_domains_and_node_ips"`
	Interval string `json:"interval"`
}

type YunDunUsageResponse struct {
	Status struct {
		Code                        int    `json:"code"`
		Message                     string `json:"message"`
		CreateAt                    string `json:"create_at"`
		ApiTimeConsuming            string `json:"api_time_consuming"`
		FunctionTimeConsuming       string `json:"function_time_consuming"`
		DispatchBeforeTimeConsuming string `json:"dispatch_before_time_consuming"`
	} `json:"status"`
	Data struct {
		HttpsTimes struct {
			Description string `json:"description"`
			Trend       struct {
				XData []string `json:"x_data"`
				YData []int    `json:"y_data"`
			} `json:"trend"`
			Total struct {
				Unit  string `json:"unit"`
				Total int    `json:"total"`
			} `json:"total"`
		} `json:"https_times"`
		TotalTimes struct {
			Description string `json:"description"`
			Trend       struct {
				XData []string `json:"x_data"`
				YData []int    `json:"y_data"`
			} `json:"trend"`
			Total struct {
				Unit  string `json:"unit"`
				Total int    `json:"total"`
			} `json:"total"`
		} `json:"total_times"`
		HitCacheTimes struct {
			Description string `json:"description"`
			Trend       struct {
				XData []string `json:"x_data"`
				YData []int    `json:"y_data"`
			} `json:"trend"`
			Total struct {
				Unit  string `json:"unit"`
				Total int    `json:"total"`
			} `json:"total"`
		} `json:"hit_cache_times"`
	} `json:"data"`
}

type YunDunErrorResponse struct {
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
}

// RefreshUrl 刷新URL
func (y *YunDun) RefreshUrl(urls []string) bool {
	client, err := y.login()
	if err != nil {
		facades.Log().With(map[string]any{
			"err": err.Error(),
		}).Error("CDN[云盾][登录失败]")
		return false
	}

	// 提交刷新请求
	refreshURL := "https://www.yundun.com/api/V4/Web.Domain.DashBoard.saveCache"
	data := map[string][]string{
		"specialurl": urls,
	}

	var refreshResponse YunDunRefreshResponse
	var errorResponse YunDunErrorResponse
	resp, err := client.R().SetBody(data).SetSuccessResult(&refreshResponse).SetErrorResult(&errorResponse).Put(refreshURL)
	if err != nil {
		facades.Log().With(map[string]any{
			"urls":  urls,
			"resp":  resp.String(),
			"error": err.Error(),
		}).Error("CDN[云盾][URL刷新失败]")
		return false
	}

	if refreshResponse.Status.Code == 1 {
		return true
	}

	facades.Log().With(map[string]any{
		"urls":  urls,
		"error": errorResponse.Status.Message,
	}).Error("CDN[云盾][URL刷新失败]")
	return false
}

// RefreshPath 刷新路径
func (y *YunDun) RefreshPath(paths []string) bool {
	client, err := y.login()
	if err != nil {
		facades.Log().With(map[string]any{
			"err": err.Error(),
		}).Error("CDN[云盾][登录失败]")
		return false
	}

	// 提交刷新请求
	refreshURL := "https://www.yundun.com/api/V4/Web.Domain.DashBoard.saveCache"
	data := map[string][]string{
		"specialdir": paths,
	}

	var refreshResponse YunDunRefreshResponse
	var errorResponse YunDunErrorResponse
	resp, err := client.R().SetBody(data).SetSuccessResult(&refreshResponse).SetErrorResult(&errorResponse).Put(refreshURL)
	if err != nil {
		facades.Log().With(map[string]any{
			"paths": paths,
			"resp":  resp.String(),
			"error": err.Error(),
		}).Error("CDN[云盾][路径刷新失败]")
		return false
	}

	if refreshResponse.Status.Code == 1 {
		return true
	}

	facades.Log().With(map[string]any{
		"paths": paths,
		"error": errorResponse.Status.Message,
	}).Error("CDN[云盾][路径刷新失败]")
	return false
}

// GetUsage 获取使用量
func (y *YunDun) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {

	client, err := y.login()
	if err != nil {
		facades.Log().With(map[string]any{
			"err": err.Error(),
		}).Error("CDN[云盾][登录失败]")
		return 0
	}

	var request = YunDunUsageRequest{
		Router:    "cdn.domain.times",
		StartTime: startTime.ToDateTimeString(),
		EndTime:   endTime.ToDateTimeString(),
		Nodes:     []string{},
		GroupId:   []string{},
		SubDomain: []string{domain},
		Interval:  "1d",
	}
	var usageResponse YunDunUsageResponse
	var errorResponse YunDunErrorResponse

	resp, err := client.R().SetBodyJsonMarshal(request).SetSuccessResult(&usageResponse).SetErrorResult(&errorResponse).Post("https://www.yundun.com/api/V4/stati.data.get")
	if err != nil {
		facades.Log().With(map[string]any{
			"domain": domain,
			"start":  startTime.ToDateTimeString(),
			"end":    endTime.ToDateTimeString(),
			"resp":   resp.String(),
			"error":  err.Error(),
		}).Error("CDN[云盾][获取用量失败]")
		return 0
	}

	if usageResponse.Status.Code == 1 {
		return uint(usageResponse.Data.TotalTimes.Total.Total)
	}

	facades.Log().With(map[string]any{
		"domain": domain,
		"start":  startTime.ToDateTimeString(),
		"end":    endTime.ToDateTimeString(),
		"error":  errorResponse.Status.Message,
	}).Error("CDN[云盾][获取用量失败]")
	return 0
}

// login 登录平台
func (y *YunDun) login() (*req.Client, error) {
	timeStamp := strconv.Itoa(int(carbon.Now().TimestampMilli()))
	rand.NewSource(time.Now().UnixNano())
	random := helper.RandomNumber(16)
	callback := "jsonp_" + timeStamp + "_" + random
	attachURL := fmt.Sprintf("https://www.yundun.com/api/sso/V4/attach?callback=%s&_time=%s", callback, timeStamp)

	client := req.C()
	client.ImpersonateSafari()

	// 先获取登录 Token
	_, err := client.R().Get(attachURL)
	if err != nil {
		return nil, err
	}

	// 提交登录请求
	loginURL := "https://www.yundun.com/api/sso/V4/login?sso_version=2"
	loginParams := map[string]string{
		"username": y.UserName,
		"password": y.PassWord,
	}
	_, err = client.R().SetFormData(loginParams).Post(loginURL)
	if err != nil {
		return nil, err
	}

	return client, nil
}
