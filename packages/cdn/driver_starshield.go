package cdn

import (
	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/facades"
	. "github.com/jdcloud-api/jdcloud-sdk-go/core"
	. "github.com/jdcloud-api/jdcloud-sdk-go/services/starshield/apis"
	. "github.com/jdcloud-api/jdcloud-sdk-go/services/starshield/client"
)

type StarShield struct {
	AccessKey, SecretKey string // 密钥
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
	request.SetPrefixes(urls)

	resp, err := client.PurgeFilesByCache_TagsAndHostOrPrefix(request)
	if err != nil {
		facades.Log.Error("CDN[星盾] ", " [URL缓存刷新失败] RequestID: "+resp.RequestID, " ", resp.Error.Message)
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
	request.SetPrefixes(paths)

	resp, err := client.PurgeFilesByCache_TagsAndHostOrPrefix(request)
	if err != nil {
		facades.Log.Error("CDN[星盾] ", " [目录缓存刷新失败] RequestID: "+resp.RequestID, " ", resp.Error.Message)
		return false
	}

	return true
}

// GetUsage 获取用量
func (s *StarShield) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	credentials := NewCredentials(s.AccessKey, s.SecretKey)
	client := NewStarshieldClient(credentials)
	client.DisableLogger()
	request := NewZoneRequestSumRequest(s.ZoneID, "all", domain, startTime.ToIso8601MilliString(), endTime.ToIso8601MilliString())
	request.AddHeader("x-jdcloud-account-id", s.ZoneID)

	resp, err := client.ZoneRequestSum(request)
	if err != nil {
		facades.Log.Error("CDN[星盾] ", " [获取用量失败] RequestID: "+resp.RequestID, " ", resp.Error.Message)
		return 0
	}

	return uint(resp.Result.Value)
}
