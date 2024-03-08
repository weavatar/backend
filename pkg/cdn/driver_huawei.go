package cdn

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	cdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"
	"github.com/spf13/cast"
)

type HuaWei struct {
	AccessKey, SecretKey string // 密钥
}

// RefreshUrl 刷新URL
func (r *HuaWei) RefreshUrl(urls []string) bool {
	for i, url := range urls {
		urls[i] = "https://" + url
	}

	auth, err := global.NewCredentialsBuilder().
		WithAk(r.AccessKey).
		WithSk(r.SecretKey).
		SafeBuild()
	if err != nil {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"err": err.Error(),
		}).Warning("初始化华为云CDN客户端失败")
	}

	build, err := cdn.CdnClientBuilder().
		WithRegion(region.CN_NORTH_1).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"err": err.Error(),
		}).Warning("初始化华为云CDN客户端失败")

	}

	client := cdn.NewCdnClient(build)
	request := &model.CreateRefreshTasksRequest{}
	typeRefreshTask := model.GetRefreshTaskRequestBodyTypeEnum().PREFIX
	modeRefreshTask := model.GetRefreshTaskRequestBodyModeEnum().ALL
	refreshTaskbody := &model.RefreshTaskRequestBody{
		Type: &typeRefreshTask,
		Mode: &modeRefreshTask,
		Urls: urls,
	}
	request.Body = &model.RefreshTaskRequest{
		RefreshTask: refreshTaskbody,
	}

	response, err := client.CreateRefreshTasks(request)
	if err != nil {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"err": err.Error(),
		}).Warning("刷新URL失败")
		return false
	}

	if response.HttpStatusCode != 200 {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"resp": response.RefreshTask,
		}).Warning("刷新URL失败")
		return false
	}

	return true
}

// RefreshPath 刷新路径
func (r *HuaWei) RefreshPath(paths []string) bool {
	for i, url := range paths {
		paths[i] = "https://" + url
	}

	auth, err := global.NewCredentialsBuilder().
		WithAk(r.AccessKey).
		WithSk(r.SecretKey).
		SafeBuild()
	if err != nil {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"err": err.Error(),
		}).Warning("初始化华为云CDN客户端失败")
	}

	build, err := cdn.CdnClientBuilder().
		WithRegion(region.CN_NORTH_1).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"err": err.Error(),
		}).Warning("初始化华为云CDN客户端失败")

	}

	client := cdn.NewCdnClient(build)
	request := &model.CreateRefreshTasksRequest{}
	typeRefreshTask := model.GetRefreshTaskRequestBodyTypeEnum().DIRECTORY
	modeRefreshTask := model.GetRefreshTaskRequestBodyModeEnum().ALL
	refreshTaskbody := &model.RefreshTaskRequestBody{
		Type: &typeRefreshTask,
		Mode: &modeRefreshTask,
		Urls: paths,
	}
	request.Body = &model.RefreshTaskRequest{
		RefreshTask: refreshTaskbody,
	}

	response, err := client.CreateRefreshTasks(request)
	if err != nil {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"err": err.Error(),
		}).Warning("刷新路径失败")
		return false
	}

	if response.HttpStatusCode != 200 {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"resp": response.RefreshTask,
		}).Warning("刷新路径失败")
		return false
	}

	return true
}

// GetUsage 获取用量
func (r *HuaWei) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	auth, err := global.NewCredentialsBuilder().
		WithAk(r.AccessKey).
		WithSk(r.SecretKey).
		SafeBuild()
	if err != nil {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"err": err.Error(),
		}).Warning("初始化华为云CDN客户端失败")
	}

	build, err := cdn.CdnClientBuilder().
		WithRegion(region.CN_NORTH_1).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"err": err.Error(),
		}).Warning("初始化华为云CDN客户端失败")

	}

	client := cdn.NewCdnClient(build)
	request := &model.ShowDomainStatsRequest{}
	request.Action = "summary"
	request.StartTime = startTime.TimestampMilli()
	request.EndTime = endTime.TimestampMilli()
	request.DomainName = domain
	request.StatType = "req_num"
	response, err := client.ShowDomainStats(request)
	if err != nil {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"err":  err.Error(),
			"resp": response.Result,
		}).Warning("获取用量失败")
		return 0
	}

	if response.HttpStatusCode != 200 {
		facades.Log().Tags("CDN", "华为云").With(map[string]any{
			"resp": response.Result,
		}).Warning("获取用量失败")
		return 0
	}

	if _, ok := response.Result["req_num"]; ok {
		return cast.ToUint(response.Result["req_num"])
	}

	return uint(0)
}
