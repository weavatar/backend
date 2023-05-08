package cdn

import (
	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"
)

type DDun struct {
	apiKey, apiSecret string
}

type DDunClean struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

type DDunRefreshResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"msg"`
}

type DDunUsageResponse struct {
	Code    uint      `json:"code"`
	Data    [][2]uint `json:"data"`
	Message string    `json:"msg"`
}

// RefreshUrl 刷新URL
func (d *DDun) RefreshUrl(urls []string) bool {
	client := req.C()

	data := make([]DDunClean, len(urls))
	for i, url := range urls {
		data[i] = DDunClean{
			Type: "clean_url",
			Data: map[string]string{"url": url},
		}
	}

	var resp DDunRefreshResponse

	_, err := client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).SetHeaders(map[string]string{
		"api-key":    d.apiKey,
		"api-secret": d.apiSecret,
	}).Post("http://cdn.ddunyun.com/v1/jobs")
	if err != nil {
		return false
	}

	if resp.Code != 0 {
		facades.Log.Error("CDN[盾云] ", " [URL缓存刷新失败] "+resp.Message)
		return false
	}

	return true
}

// RefreshPath 刷新路径
func (d *DDun) RefreshPath(paths []string) bool {
	client := req.C()

	data := make([]DDunClean, len(paths))
	for i, url := range paths {
		data[i] = DDunClean{
			Type: "clean_dir",
			Data: map[string]string{"url": url},
		}
	}

	var resp DDunRefreshResponse

	post, err := client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).SetHeaders(map[string]string{
		"api-key":    d.apiKey,
		"api-secret": d.apiSecret,
	}).Post("http://cdn.ddunyun.com/v1/jobs")
	if err != nil {
		return false
	}
	if !post.IsSuccessState() {
		return false
	}

	if resp.Code != 0 {
		facades.Log.Error("CDN[盾云] ", " [目录缓存刷新失败] "+resp.Message)
		return false
	}

	return true
}

// GetUsage 获取用量
func (d *DDun) GetUsage(domain, startTime, endTime string) uint {
	client := req.C()

	var resp DDunUsageResponse

	// 由于cdnfly这个非标准querystring，所以只能手动把参数拼接到url上了
	post, err := client.R().SetSuccessResult(&resp).SetErrorResult(&resp).SetHeaders(map[string]string{
		"api-key":    d.apiKey,
		"api-secret": d.apiSecret,
	}).Get("http://cdn.ddunyun.com/v1/monitor/site/realtime?type=req&start=" + startTime + "%2000:00:00" + "&end=" + endTime + "%2000:00:00" + "&domain=" + domain + "&server_post=")

	if err != nil {
		return 0
	}
	if !post.IsSuccessState() {
		return 0
	}

	if resp.Code != 0 {
		facades.Log.Error("CDN[盾云] ", " [获取用量失败] "+resp.Message)
		return 0
	}

	sum := uint(0)
	for _, data := range resp.Data {
		sum += data[1]
	}

	return sum
}
