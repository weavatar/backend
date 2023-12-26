package cdn

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/imroc/req/v3"
)

type CTYun struct {
	AppID       string
	AppSecret   string
	ApiEndpoint string
}

type CTYunRefreshResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	SubmitID string `json:"submit_id"`
	Result   []struct {
		TaskID string `json:"task_id"`
		URL    string `json:"url"`
	} `json:"result"`
}

type CTYunUsageResponse struct {
	StartTime                 int64  `json:"start_time"`
	Code                      int    `json:"code"`
	EndTime                   int64  `json:"end_time"`
	Interval                  string `json:"interval"`
	Message                   string `json:"message"`
	ReqRequestNumDataInterval []struct {
		HitRequestRate           float64 `json:"hit_request_rate"`
		TimeStamp                int64   `json:"time_stamp"`
		MissRequestNum           int     `json:"miss_request_num"`
		RequestNum               int     `json:"request_num"`
		ApplicationLayerProtocol string  `json:"application_layer_protocol"`
	} `json:"req_request_num_data_interval"`
}

// RefreshUrl 刷新URL
func (c *CTYun) RefreshUrl(urls []string) bool {
	api := "/api/v1/refreshmanage/create"

	timestamp, signature, err := c.getSignature(api)
	if err != nil {
		facades.Log().With(map[string]any{
			"err": err.Error(),
		}).Error("CDN[天翼云][计算签名失败]")
		return false
	}

	client := req.C()
	client.SetTimeout(60 * time.Second)

	client.SetCommonHeaders(map[string]string{
		"x-alogic-now":       timestamp,
		"x-alogic-app":       c.AppID,
		"x-alogic-ac":        "app",
		"x-alogic-signature": signature,
	})

	for i, url := range urls {
		urls[i] = "https://" + url
	}

	data := map[string]any{
		"values":    urls,
		"task_type": 1,
	}

	var resp CTYunRefreshResponse
	_, err = client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).Post(c.ApiEndpoint + api)
	if err != nil {
		facades.Log().With(map[string]any{
			"code":    resp.Code,
			"message": resp.Message,
			"err":     err.Error(),
		}).Error("CDN[天翼云][URL缓存刷新失败]")
		return false
	}

	if resp.Code != 100000 {
		facades.Log().With(map[string]any{
			"code":    resp.Code,
			"message": resp.Message,
		}).Error("CDN[天翼云][URL缓存刷新失败]")
		return false
	}

	return true
}

// RefreshPath 刷新路径
func (c *CTYun) RefreshPath(paths []string) bool {
	api := "/api/v1/refreshmanage/create"

	timestamp, signature, err := c.getSignature(api)
	if err != nil {
		facades.Log().With(map[string]any{
			"err": err.Error(),
		}).Error("CDN[天翼云][计算签名失败]")
		return false
	}

	client := req.C()
	client.SetTimeout(60 * time.Second)

	client.SetCommonHeaders(map[string]string{
		"x-alogic-now":       timestamp,
		"x-alogic-app":       c.AppID,
		"x-alogic-ac":        "app",
		"x-alogic-signature": signature,
	})

	// 天翼云文档要求统一使用 http 协议
	for i, path := range paths {
		paths[i] = "http://" + path
	}

	data := map[string]any{
		"values":    paths,
		"task_type": 2,
	}

	var resp CTYunRefreshResponse
	_, err = client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).Post(c.ApiEndpoint + api)
	if err != nil {
		facades.Log().With(map[string]any{
			"code":    resp.Code,
			"message": resp.Message,
			"err":     err.Error(),
		}).Error("CDN[天翼云][路径缓存刷新失败]")
		return false
	}

	if resp.Code != 100000 {
		facades.Log().With(map[string]any{
			"code":    resp.Code,
			"message": resp.Message,
		}).Error("CDN[天翼云][路径缓存刷新失败]")
		return false
	}

	return true
}

// GetUsage 获取使用量
func (c *CTYun) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	api := "/api/v2/statisticsanalysis/query_request_num_data"

	timestamp, signature, err := c.getSignature(api)
	if err != nil {
		facades.Log().With(map[string]any{
			"err": err.Error(),
		}).Error("CDN[天翼云][计算签名失败]")
		return 0
	}

	client := req.C()
	client.SetTimeout(60 * time.Second)

	client.SetCommonHeaders(map[string]string{
		"x-alogic-now":       timestamp,
		"x-alogic-app":       c.AppID,
		"x-alogic-ac":        "app",
		"x-alogic-signature": signature,
	})

	var usage CTYunUsageResponse
	resp, err := client.R().SetBodyJsonMarshal(map[string]any{
		"interval":   "24h",
		"domain":     []string{domain},
		"start_time": startTime.Timestamp(),
		"end_time":   endTime.Timestamp(),
	}).SetSuccessResult(&usage).Post(c.ApiEndpoint + api)
	if err != nil {
		facades.Log().With(map[string]any{
			"code":    usage.Code,
			"message": usage.Message,
			"err":     err.Error(),
		}).Error("CDN[天翼云][获取用量失败]")
		return 0
	}

	if usage.Code != 100000 {
		facades.Log().With(map[string]any{
			"code":     usage.Code,
			"message":  usage.Message,
			"response": resp.String(),
		}).Error("CDN[天翼云][获取用量失败]")
		return 0
	}

	sum := uint(0)
	for _, data := range usage.ReqRequestNumDataInterval {
		sum += uint(data.RequestNum)
	}

	return sum
}

func (c *CTYun) hmacSha256Byte(target, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(target))
	hashBytes := h.Sum(nil)

	return hashBytes
}

func (c *CTYun) encrypt(content, key string) (signature string, err error) {
	// 替换空格为+
	key = strings.ReplaceAll(key, " ", "+")
	// 替换-为+号
	key = strings.ReplaceAll(key, "-", "+")
	// 替换_为/号
	key = strings.ReplaceAll(key, "_", "/")
	// 填充=，字节为4的倍数
	for len(key)%4 != 0 {
		key += "="
	}
	b64Code, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err

	}

	signedByte := c.hmacSha256Byte(content, string(b64Code))

	signedStr := base64.URLEncoding.EncodeToString(signedByte)

	signature = strings.Replace(signedStr, "=", "", -1)

	return signature, nil
}

func (c *CTYun) getSignature(url string) (string, string, error) {

	timestampMs := time.Now().Unix() * 1000

	timestampDay := timestampMs / 86400000

	timestampMsStr := strconv.FormatInt(timestampMs, 10)

	signStr := fmt.Sprintf("%s\n%v\n%s", c.AppID, timestampMs, url)

	identity := fmt.Sprintf("%s:%v", c.AppID, timestampDay)

	tmpSignature, err := c.encrypt(identity, c.AppSecret)

	if err != nil {
		return "", "", err

	}

	signature, err := c.encrypt(signStr, tmpSignature)

	if err != nil {
		return "", "", err

	}

	return timestampMsStr, signature, nil
}
