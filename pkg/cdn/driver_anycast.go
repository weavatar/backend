package cdn

import (
	"github.com/imroc/req/v3"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
)

type AnyCast struct {
	apiKey, apiSecret string
}

type AnyCastClean struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

type AnyCastRefreshResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"msg"`
}

type AnyCastUsageResponse struct {
	Code    uint      `json:"code"`
	Data    [][2]uint `json:"data"`
	Message string    `json:"msg"`
}

// RefreshUrl 刷新URL
func (d *AnyCast) RefreshUrl(urls []string) bool {
	client := req.C()

	data := make([]AnyCastClean, len(urls))
	for i, url := range urls {
		data[i] = AnyCastClean{
			Type: "clean_url",
			Data: map[string]string{"url": "https://" + url + "*"},
		}
	}

	var resp AnyCastRefreshResponse
	_, err := client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).SetHeaders(map[string]string{
		"api-key":    d.apiKey,
		"api-secret": d.apiSecret,
	}).Post("https://console.anycast.ai/v1/jobs")
	if err != nil {
		return false
	}

	if resp.Code != 0 {
		facades.Log().Error("CDN[AnyCast][URL缓存刷新失败] " + resp.Message)
		return false
	}

	return true
}

// RefreshPath 刷新路径
func (d *AnyCast) RefreshPath(paths []string) bool {
	client := req.C()

	data := make([]AnyCastClean, len(paths))
	for i, url := range paths {
		data[i] = AnyCastClean{
			Type: "clean_dir",
			Data: map[string]string{"url": "https://" + url},
		}
	}

	var resp AnyCastRefreshResponse
	post, err := client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).SetHeaders(map[string]string{
		"api-key":    d.apiKey,
		"api-secret": d.apiSecret,
	}).Post("https://console.anycast.ai/v1/jobs")
	if err != nil {
		return false
	}
	if !post.IsSuccessState() {
		return false
	}

	if resp.Code != 0 {
		facades.Log().Error("CDN[AnyCast][目录缓存刷新失败] " + resp.Message)
		return false
	}

	return true
}

// GetUsage 获取用量
func (d *AnyCast) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	client := req.C()

	var resp AnyCastUsageResponse
	post, err := client.R().SetSuccessResult(&resp).SetErrorResult(&resp).SetHeaders(map[string]string{
		"api-key":    d.apiKey,
		"api-secret": d.apiSecret,
	}).Get("https://console.anycast.ai/v1/monitor/site/realtime?type=req&start=" + startTime.ToDateString() + "%2000:00:00" + "&end=" + endTime.ToDateString() + "%2000:00:00" + "&domain=" + domain + "&server_post=")

	if err != nil {
		return 0
	}
	if !post.IsSuccessState() {
		return 0
	}

	if resp.Code != 0 {
		facades.Log().Error("CDN[AnyCast][获取用量失败] " + resp.Message)
		return 0
	}

	sum := uint(0)
	for _, data := range resp.Data {
		sum += data[1]
	}

	return sum
}
