package cdn

import (
	"time"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/imroc/req/v3"
)

type BaiShan struct {
	Token string
}

type BaiShanRefreshResponse struct {
	Code uint `json:"code"`
	Data any  `json:"data"`
}

type BaiShanUsageResponse struct {
	Code int `json:"code"`
	Data map[string]struct {
		Domain string   `json:"domain"`
		Data   [][]uint `json:"data"`
	} `json:"data"`
}

// RefreshUrl 刷新URL
func (b *BaiShan) RefreshUrl(urls []string) bool {
	client := req.C()
	client.SetTimeout(60 * time.Second)

	for i, url := range urls {
		urls[i] = "https://" + url
	}

	data := map[string]any{
		"urls": urls,
		"type": "url",
	}

	var resp BaiShanRefreshResponse
	_, err := client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).Post("https://cdn.api.baishan.com/v2/cache/refresh?token=" + b.Token)
	if err != nil {
		facades.Log().With(map[string]any{
			"code": resp.Code,
			"data": resp.Data,
			"err":  err.Error(),
		}).Error("CDN[白山][URL缓存刷新失败]")
		return false
	}

	if resp.Code != 0 {
		facades.Log().With(map[string]any{
			"code": resp.Code,
			"data": resp.Data,
		}).Error("CDN[白山][URL缓存刷新失败]")
		return false
	}

	return true
}

// RefreshPath 刷新路径
func (b *BaiShan) RefreshPath(paths []string) bool {
	client := req.C()
	client.SetTimeout(60 * time.Second)

	for i, path := range paths {
		paths[i] = "https://" + path
	}

	refreshURL := "https://cdn.api.baishan.com/v2/cache/refresh?token=" + b.Token
	data := map[string]any{
		"urls": paths,
		"type": "dir",
	}

	var resp BaiShanRefreshResponse
	_, err := client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).Post(refreshURL)
	if err != nil {
		facades.Log().With(map[string]any{
			"code": resp.Code,
			"data": resp.Data,
			"err":  err.Error(),
		}).Error("CDN[白山][路径缓存刷新失败]")
		return false
	}

	if resp.Code != 0 {
		facades.Log().With(map[string]any{
			"code": resp.Code,
			"data": resp.Data,
		}).Error("CDN[白山][路径缓存刷新失败]")
		return false
	}

	return true
}

// GetUsage 获取使用量
func (b *BaiShan) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	client := req.C()
	client.SetTimeout(60 * time.Second)

	var usage BaiShanUsageResponse
	resp, err := client.R().SetQueryParams(map[string]string{
		"token":      b.Token,
		"domains":    domain,
		"start_time": startTime.ToDateTimeString(),
		"end_time":   endTime.ToDateTimeString(),
	}).SetSuccessResult(&usage).Get("https://cdn.api.baishan.com/v2/stat/request/eachDomain")
	if err != nil {
		facades.Log().With(map[string]any{
			"code": usage.Code,
			"data": usage.Data,
			"err":  err.Error(),
		}).Error("CDN[白山][获取用量失败]")
		return 0
	}

	if usage.Code != 0 {
		facades.Log().With(map[string]any{
			"code":     usage.Code,
			"data":     usage.Data,
			"response": resp.String(),
		}).Error("CDN[白山][获取用量失败]")
		return 0
	}

	sum := uint(0)
	for _, domain := range usage.Data {
		for _, data := range domain.Data {
			sum += data[1]
		}
	}

	return sum
}
