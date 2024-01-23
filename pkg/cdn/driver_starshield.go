package cdn

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	. "github.com/jdcloud-api/jdcloud-sdk-go/core"
	. "github.com/jdcloud-api/jdcloud-sdk-go/services/starshield/apis"
	. "github.com/jdcloud-api/jdcloud-sdk-go/services/starshield/client"
)

type StarShield struct {
	AccessKey, SecretKey string // 密钥
	InstanceID           string // 实例ID
	ZoneID               string // 域名标识
}

type StarShieldClean struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

type StarShieldRefreshResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"msg"`
}

type StarShieldUsageResponse struct {
	Code    uint      `json:"code"`
	Data    [][2]uint `json:"data"`
	Message string    `json:"msg"`
}

// RefreshUrl 刷新URL
func (s *StarShield) RefreshUrl(urls []string) bool {
	credentials := NewCredentials(s.AccessKey, s.SecretKey)
	client := NewStarshieldClient(credentials)
	client.DisableLogger()
	request := NewPurgeFilesByCache_TagsAndHostOrPrefixRequest(s.ZoneID)
	request.AddHeader("x-jdcloud-account-id", s.InstanceID)
	request.SetPrefixes(urls)

	resp, err := client.PurgeFilesByCache_TagsAndHostOrPrefix(request)
	if err != nil {
		facades.Log().Tags("CDN", "星盾").With(map[string]any{
			"err":  err.Error(),
			"resp": resp,
		}).Error("URL缓存刷新失败")
		return false
	}
	if resp.Error.Code != 0 {
		facades.Log().Tags("CDN", "星盾").With(map[string]any{
			"resp": resp,
		}).Error("URL缓存刷新失败")
		return false
	}

	return true
}

// RefreshPath 刷新路径
func (s *StarShield) RefreshPath(paths []string) bool {
	credentials := NewCredentials(s.AccessKey, s.SecretKey)
	client := NewStarshieldClient(credentials)
	client.DisableLogger()
	request := NewPurgeFilesByCache_TagsAndHostOrPrefixRequest(s.ZoneID)
	request.AddHeader("x-jdcloud-account-id", s.InstanceID)
	request.SetPrefixes(paths)

	resp, err := client.PurgeFilesByCache_TagsAndHostOrPrefix(request)
	if err != nil {
		facades.Log().Tags("CDN", "星盾").With(map[string]any{
			"err":  err.Error(),
			"resp": resp,
		}).Error("目录缓存刷新失败")
		return false
	}
	if resp.Error.Code != 0 {
		facades.Log().Tags("CDN", "星盾").With(map[string]any{
			"resp": resp,
		}).Error("目录缓存刷新失败")
		return false
	}

	return true
}

// GetUsage 获取用量
func (s *StarShield) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	credentials := NewCredentials(s.AccessKey, s.SecretKey)
	client := NewStarshieldClient(credentials)
	client.DisableLogger()
	request := NewZoneRequestSumRequest(s.ZoneID, "all", domain, startTime.ToDateString()+"T00:00:00.000Z", endTime.ToDateString()+"T00:00:00.000Z")
	request.AddHeader("x-jdcloud-account-id", s.InstanceID)

	resp, err := client.ZoneRequestSum(request)
	if err != nil {
		facades.Log().Tags("CDN", "星盾").With(map[string]any{
			"err":  err.Error(),
			"resp": resp,
		}).Error("获取用量失败")
		return 0
	}
	if resp.Error.Code != 0 {
		facades.Log().Tags("CDN", "星盾").With(map[string]any{
			"resp": resp,
		}).Error("获取用量失败")
		return 0
	}

	return uint(resp.Result.Value)
}
