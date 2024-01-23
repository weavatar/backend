package cdn

import (
	"context"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/imroc/req/v3"
)

type CloudFlare struct {
	Key, Email string // 密钥
	ZoneID     string // 域名标识
}

// CloudFlareGraphQLQuery 结构体用于构造 GraphQL 查询
type CloudFlareGraphQLQuery struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

// CloudFlareHttpRequests 结构体用于解析 GraphQL 查询结果
type CloudFlareHttpRequests struct {
	Data struct {
		Viewer struct {
			Zones []struct {
				HttpRequests1DGroups []struct {
					Sum struct {
						Requests int `json:"requests"`
					} `json:"sum"`
				} `json:"httpRequests1dGroups"`
			} `json:"zones"`
		} `json:"viewer"`
	} `json:"data"`
	Errors any `json:"errors"`
}

// RefreshUrl 刷新URL
func (s *CloudFlare) RefreshUrl(urls []string) bool {
	api, err := cloudflare.New(s.Key, s.Email)
	if err != nil {
		facades.Log().Tags("CDN", "CloudFlare").With(map[string]any{
			"err": err.Error(),
		}).Warning("初始化失败")
		return false
	}

	resp, err := api.PurgeCache(context.Background(), s.ZoneID, cloudflare.PurgeCacheRequest{
		Prefixes: urls,
	})
	if err != nil {
		facades.Log().Tags("CDN", "CloudFlare").With(map[string]any{
			"err":  err.Error(),
			"resp": resp,
		}).Warning("URL缓存刷新失败")
		return false
	}
	if !resp.Response.Success {
		facades.Log().Tags("CDN", "CloudFlare").With(map[string]any{
			"err":  err.Error(),
			"resp": resp.Response.Errors,
		}).Warning("URL缓存刷新失败")
		return false
	}

	return true
}

// RefreshPath 刷新路径
func (s *CloudFlare) RefreshPath(paths []string) bool {
	api, err := cloudflare.New(s.Key, s.Email)
	if err != nil {
		facades.Log().Tags("CDN", "CloudFlare").With(map[string]any{
			"err": err.Error(),
		}).Warning("初始化失败")
		return false
	}

	resp, err := api.PurgeCache(context.Background(), s.ZoneID, cloudflare.PurgeCacheRequest{
		Prefixes: paths,
	})
	if err != nil {
		facades.Log().Tags("CDN", "CloudFlare").With(map[string]any{
			"err":  err.Error(),
			"resp": resp,
		}).Warning("路径缓存刷新失败")
		return false
	}
	if !resp.Response.Success {
		facades.Log().Tags("CDN", "CloudFlare").With(map[string]any{
			"err":  err.Error(),
			"resp": resp.Response.Errors,
		}).Warning("路径缓存刷新失败")
		return false
	}

	return true
}

// GetUsage 获取用量
func (s *CloudFlare) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	client := req.C()
	client.SetBaseURL("https://api.cloudflare.com/client/v4")
	client.SetTimeout(10 * time.Second)
	client.SetCommonRetryCount(2)
	client.SetCommonHeaders(map[string]string{
		"X-Auth-Email": s.Email,
		"X-Auth-Key":   s.Key,
	})

	query := CloudFlareGraphQLQuery{
		Query: `
		{
		  viewer {
			zones(filter: {zoneTag: $zoneTag}) {
			  httpRequests1dGroups(limit: 1, filter: {date_gt: $start, date_lt: $end}) {
				sum {
				  requests
				}
			  }
			}
		  }
		}
        `,
		Variables: map[string]any{
			"zoneTag": s.ZoneID,
			// CloudFlare 不这样写的话取不到数据
			"start": startTime.Yesterday().ToDateString(),
			"end":   endTime.ToDateString(),
		},
	}

	var resp CloudFlareHttpRequests
	_, err := client.R().SetBodyJsonMarshal(query).SetSuccessResult(&resp).SetErrorResult(&resp).Post("/graphql")
	if err != nil {
		facades.Log().Tags("CDN", "CloudFlare").With(map[string]any{
			"err":  err.Error(),
			"resp": resp,
		}).Warning("获取用量失败")
		return 0
	}

	// 数据可能为空，需要判断
	if len(resp.Data.Viewer.Zones) == 0 || len(resp.Data.Viewer.Zones[0].HttpRequests1DGroups) == 0 {
		return 0
	}

	return uint(resp.Data.Viewer.Zones[0].HttpRequests1DGroups[0].Sum.Requests)
}
