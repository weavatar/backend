package cdn

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/json"

	"weavatar/pkg/wangsu/api/client/purge"
	"weavatar/pkg/wangsu/api/client/summary"
	"weavatar/pkg/wangsu/common/auth"
)

type WangSu struct {
	AccessKey, SecretKey string // 密钥
}

// RefreshUrl 刷新URL
func (receiver *WangSu) RefreshUrl(urls []string) bool {
	var pointers []*string
	for _, url := range urls {
		url = "https://" + url + "**"
		pointers = append(pointers, &url)
	}

	createAPurgeRequestRequest := purge.CreateAPurgeRequestRequest{}
	createAPurgeRequestRequest.SetDirUrls(pointers)
	createAPurgeRequestRequest.SetTarget("production")

	var config auth.AkskConfig
	config.AccessKey = receiver.AccessKey
	config.SecretKey = receiver.SecretKey
	config.EndPoint = "open.chinanetcenter.com"
	config.Uri = "/cdn/purges"
	config.Method = "POST"
	response := auth.Invoke(config, createAPurgeRequestRequest.String())

	if len(response) != 0 {
		facades.Log().Tags("CDN", "网宿").With(map[string]any{
			"urls":     urls,
			"response": response,
		}).Warning("URL缓存刷新失败")
		return false
	}

	return true
}

// RefreshPath 刷新路径
func (receiver *WangSu) RefreshPath(paths []string) bool {
	var pointers []*string
	for _, path := range paths {
		path = "https://" + path + "**"
		pointers = append(pointers, &path)
	}

	createAPurgeRequestRequest := purge.CreateAPurgeRequestRequest{}
	createAPurgeRequestRequest.SetDirUrls(pointers)
	createAPurgeRequestRequest.SetTarget("production")

	var config auth.AkskConfig
	config.AccessKey = receiver.AccessKey
	config.SecretKey = receiver.SecretKey
	config.EndPoint = "open.chinanetcenter.com"
	config.Uri = "/cdn/purges"
	config.Method = "POST"
	response := auth.Invoke(config, createAPurgeRequestRequest.String())

	if len(response) != 0 {
		facades.Log().Tags("CDN", "网宿").With(map[string]any{
			"paths":    paths,
			"response": response,
		}).Warning("路径缓存刷新失败")
		return false
	}

	return true
}

// GetUsage 获取用量
func (receiver *WangSu) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	getASummaryOfRequestsRequest := summary.GetASummaryOfRequestsRequest{}
	filters := summary.GetASummaryOfRequestsRequestFilters{}

	filters.SetHostnames([]*string{&domain})
	getASummaryOfRequestsRequest.SetFilters(&filters)

	getASummaryOfRequestsParams := summary.Parameters{}
	getASummaryOfRequestsParams.SetStartdate(startTime.SetTimezone(carbon.UTC).ToRfc3339String())
	getASummaryOfRequestsParams.SetEnddate(endTime.SetTimezone(carbon.UTC).ToRfc3339String())
	getASummaryOfRequestsParams.SetScheme("all")

	var config auth.AkskConfig
	config.AccessKey = receiver.AccessKey
	config.SecretKey = receiver.SecretKey
	config.EndPoint = "open.chinanetcenter.com"
	config.Uri = "/cdn/report/reqSummary"
	config.Method = "POST"
	response := auth.Invoke(config, getASummaryOfRequestsRequest.String(), getASummaryOfRequestsParams.String())

	var data summary.GetASummaryOfRequestsResponse
	err := json.UnmarshalString(response, &data)
	if err != nil {
		facades.Log().Tags("CDN", "网宿").With(map[string]any{
			"response": response,
			"error":    err.Error(),
		}).Warning("获取用量失败")
		return 0
	}

	sum := float64(0)
	for _, item := range data.Groups {
		for _, num := range item.Data {
			sum += *num
		}
	}

	return uint(sum)
}
