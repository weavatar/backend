package cdn

import (
	"github.com/bytedance/sonic"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

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
		facades.Log().With(map[string]any{
			"urls":     urls,
			"response": response,
		}).Error("CDN[网宿][URL缓存刷新失败]")
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
		facades.Log().With(map[string]any{
			"paths":    paths,
			"response": response,
		}).Error("CDN[网宿][路径缓存刷新失败]")
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

	var config auth.AkskConfig
	config.AccessKey = receiver.AccessKey
	config.SecretKey = receiver.SecretKey
	config.EndPoint = "open.chinanetcenter.com"
	config.Uri = "/cdn/report/reqSummary?startdate=" + startTime.SetTimezone(carbon.UTC).ToRfc3339String() + "&enddate=" + endTime.SetTimezone(carbon.UTC).ToRfc3339String() + "&scheme=all"
	config.Method = "POST"
	response := auth.Invoke(config, getASummaryOfRequestsRequest.String())

	var data summary.GetASummaryOfRequestsResponse
	err := sonic.UnmarshalString(response, &data)
	if err != nil {
		facades.Log().With(map[string]any{
			"response": response,
		}).Error("CDN[网宿][获取用量失败]")
		return 0
	}

	return 0
}
