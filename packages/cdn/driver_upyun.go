package cdn

import (
	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"
)

type UpYun struct {
	Token string
}

type UpYunPurgeBatch struct {
	NoIf      uint   `json:"noif"`
	SourceUrl string `json:"source_url"`
}

type UpYunPurgeBatchSuccessResponse struct {
	Code   uint   `json:"code"`
	Status string `json:"status"`
}

type UpYunUsageSuccessResponse struct {
	Rps         float64          `json:"rps"`
	Reqs        float64          `json:"reqs"`
	Bytes       float64          `json:"bytes"`
	Bandwidth   float64          `json:"bandwidth"`
	Treqs       float64          `json:"treqs"`
	Tbytes      float64          `json:"tbytes"`
	Tbandwidth  float64          `json:"tbandwidth"`
	Hreqs       float64          `json:"hreqs"`
	Hbytes      float64          `json:"hbytes"`
	Hbandwidth  float64          `json:"hbandwidth"`
	Dreqs       float64          `json:"dreqs"`
	Dbytes      float64          `json:"dbytes"`
	Dbandwidth  float64          `json:"dbandwidth"`
	Wsbytes     float64          `json:"wsbytes"`
	Wsbandwidth float64          `json:"wsbandwidth"`
	Time        carbon.Timestamp `json:"time"`
}

type UpYunErrorResponse struct {
	ErrorCode uint   `json:"error_code"`
	Request   string `json:"request"`
	Message   string `json:"message"`
}

// RefreshUrl 刷新URL
func (u *UpYun) RefreshUrl(urls []string) bool {
	client := req.C()

	var sourceUrl string
	for _, url := range urls {
		sourceUrl += url + "\n"
	}

	data := UpYunPurgeBatch{
		NoIf:      1,
		SourceUrl: sourceUrl,
	}

	var successResp []UpYunPurgeBatchSuccessResponse
	var errorResp UpYunErrorResponse

	_, err := client.R().SetBody(data).SetSuccessResult(&successResp).SetErrorResult(&errorResp).SetBearerAuthToken(u.Token).Post("https://api.upyun.com/buckets/purge/batch")
	if err != nil {
		facades.Log.Error("CDN[又拍云] ", " [URL缓存刷新错误] ", err.Error())
		return false
	}

	for _, resp := range successResp {
		if resp.Code != 1 {
			facades.Log.Error("CDN[又拍云] ", " [URL缓存刷新失败] ", resp.Code, " ", resp.Status, " ", errorResp.ErrorCode, " ", errorResp.Message)
			return false
		}
	}

	if errorResp.ErrorCode != 0 {
		facades.Log.Error("CDN[又拍云] ", " [URL缓存刷新错误] ", errorResp.ErrorCode, " ", errorResp.Message, " ", errorResp.Request)
		return false
	}

	return true
}

// RefreshPath 刷新路径
func (u *UpYun) RefreshPath(paths []string) bool {
	client := req.C()

	var sourceUrl string
	for _, path := range paths {
		sourceUrl += path + "\n"
	}

	data := UpYunPurgeBatch{
		NoIf:      1,
		SourceUrl: sourceUrl,
	}

	var successResp []UpYunPurgeBatchSuccessResponse
	var errorResp UpYunErrorResponse

	_, err := client.R().SetBody(data).SetSuccessResult(&successResp).SetErrorResult(&errorResp).SetBearerAuthToken(u.Token).Post("https://api.upyun.com/buckets/purge/batch")
	if err != nil {
		facades.Log.Error("CDN[又拍云] ", " [路径缓存刷新错误] ", err)
		return false
	}

	for _, resp := range successResp {
		if resp.Code != 1 {
			facades.Log.Error("CDN[又拍云] ", " [路径缓存刷新失败] ", resp.Code, " ", resp.Status, " ", errorResp.ErrorCode, " ", errorResp.Message)
			return false
		}
	}

	if errorResp.ErrorCode != 0 {
		facades.Log.Error("CDN[又拍云] ", " [URL缓存刷新错误] ", errorResp.ErrorCode, " ", errorResp.Message, " ", errorResp.Request)
		return false
	}

	return true
}

// GetUsage 获取用量
func (u *UpYun) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	client := req.C()

	var successResp []UpYunUsageSuccessResponse
	var errorResp UpYunErrorResponse

	_, err := client.R().SetSuccessResult(&successResp).SetErrorResult(&errorResp).SetBearerAuthToken(u.Token).SetQueryParams(map[string]string{
		"start_time":  startTime.ToIso8601MilliString(),
		"end_time":    endTime.ToIso8601MilliString(),
		"query_type":  "domain",
		"query_value": domain,
		"flow_type":   "cdn",
		"flow_source": "cdn",
	}).Get("https://api.upyun.com/flow/common_data")

	if err != nil {
		facades.Log.Error("CDN[又拍云] ", " [获取用量错误] ", err)
		return 0
	}

	if errorResp.ErrorCode != 0 {
		facades.Log.Error("CDN[又拍云] ", " [获取用量失败] ", errorResp.ErrorCode, " ", errorResp.Message, " ", errorResp.Request)
		return 0
	}

	sum := uint(0)
	for _, data := range successResp {
		sum += uint(data.Reqs)
	}

	return sum
}
