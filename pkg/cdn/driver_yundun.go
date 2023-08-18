package cdn

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/imroc/req/v3"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	"weavatar/pkg/helper"
)

type YunDun struct {
	UserName, PassWord string
}

// RefreshUrl 刷新URL
func (y *YunDun) RefreshUrl(urls []string) bool {
	timeStamp := strconv.Itoa(int(carbon.Now().TimestampMilli()))
	rand.NewSource(time.Now().UnixNano())
	random := helper.RandomNumber(16)
	callback := "jsonp_" + timeStamp + "_" + random
	attachURL := fmt.Sprintf("https://www.yundun.com/api/sso/V4/attach?callback=%s&_time=%s", callback, timeStamp)

	client := req.C()

	// 先获取登录 Token
	_, err := client.R().Get(attachURL)
	if err != nil {
		facades.Log().Error("CDN[云盾] ", " [获取登录Token失败] "+err.Error())
		return false
	}

	// 提交登录请求
	loginURL := "https://www.yundun.com/api/sso/V4/login?sso_version=2"
	loginParams := map[string]string{
		"username": y.UserName,
		"password": y.PassWord,
	}
	_, err = client.R().SetFormData(loginParams).Post(loginURL)
	if err != nil {
		facades.Log().Error("CDN[云盾] ", " [账号登录失败] "+err.Error())
		return false
	}

	// 提交刷新请求
	refreshURL := "https://www.yundun.com/api/V4/Web.Domain.DashBoard.saveCache"
	data := map[string][]string{
		"specialurl": urls,
	}

	resp, err := client.R().SetBody(data).Post(refreshURL)
	if err != nil {
		facades.Log().Error("CDN[云盾] ", " [URL刷新失败] "+err.Error())
		return false
	}

	// 判断是否刷新成功
	var refreshResponse map[string]interface{}
	err = sonic.Unmarshal(resp.Bytes(), &refreshResponse)
	if err != nil {
		facades.Log().Error("CDN[云盾] ", " [JSON解析失败] "+err.Error())
		return false
	}

	if status, ok := refreshResponse["status"].(map[string]interface{}); ok {
		if code, ok := status["code"].(float64); ok && code == 1 {
			return true
		}
	}

	return false
}

// RefreshPath 刷新路径
func (y *YunDun) RefreshPath(paths []string) bool {
	timeStamp := strconv.Itoa(int(carbon.Now().TimestampMilli()))
	rand.NewSource(time.Now().UnixNano())
	random := helper.RandomNumber(16)
	callback := "jsonp_" + timeStamp + "_" + random
	attachURL := fmt.Sprintf("https://www.yundun.com/api/sso/V4/attach?callback=%s&_time=%s", callback, timeStamp)

	client := req.C()

	// 先获取登录 Token
	_, err := client.R().Get(attachURL)
	if err != nil {
		facades.Log().Error("CDN[云盾] ", " [获取登录Token失败] "+err.Error())
		return false
	}

	// 提交登录请求
	loginURL := "https://www.yundun.com/api/sso/V4/login?sso_version=2"
	loginParams := map[string]string{
		"username": y.UserName,
		"password": y.PassWord,
	}
	_, err = client.R().SetFormData(loginParams).Post(loginURL)
	if err != nil {
		facades.Log().Error("CDN[云盾] ", " [账号登录失败] "+err.Error())
		return false
	}

	// 提交刷新请求
	refreshURL := "https://www.yundun.com/api/V4/Web.Domain.DashBoard.saveCache"
	data := map[string][]string{
		"specialdir": paths,
	}

	resp, err := client.R().SetBody(data).Post(refreshURL)
	if err != nil {
		facades.Log().Error("CDN[云盾] ", " [URL刷新失败] "+err.Error())
		return false
	}

	// 判断是否刷新成功
	var refreshResponse map[string]interface{}
	err = sonic.Unmarshal(resp.Bytes(), &refreshResponse)
	if err != nil {
		facades.Log().Error("CDN[云盾] ", " [JSON解析失败] "+err.Error())
		return false
	}

	if status, ok := refreshResponse["status"].(map[string]interface{}); ok {
		if code, ok := status["code"].(float64); ok && code == 1 {
			return true
		}
	}

	return false
}

// GetUsage 获取使用量
func (y *YunDun) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	return 0
}
