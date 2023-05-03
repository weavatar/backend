package cdn

import (
	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"
)

type DDun struct {
	apiKey, apiSecret string
}

type DDunClean struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

type DDunResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// RefreshUrl 刷新URL
func (d *DDun) RefreshUrl(urls []string) bool {
	client := req.C()

	data := make([]DDunClean, len(urls))
	for i, url := range urls {
		data[i] = DDunClean{
			Type: "clean_url",
			Data: map[string]string{"url": url},
		}
	}

	var resp DDunResponse

	_, err := client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).SetHeaders(map[string]string{
		"api-key":    d.apiKey,
		"api-secret": d.apiSecret,
	}).Post("http://cdn.ddunyun.com/v1/jobs")
	if err != nil {
		return false
	}

	if resp.Code != 0 {
		facades.Log.Error("CDN[盾云][URL缓存刷新失败] " + resp.Message)
		return false
	}

	return true
}

// RefreshPath 刷新路径
func (d *DDun) RefreshPath(paths []string) bool {
	client := req.C()

	data := make([]DDunClean, len(paths))
	for i, url := range paths {
		data[i] = DDunClean{
			Type: "clean_dir",
			Data: map[string]string{"url": url},
		}
	}

	var resp DDunResponse

	post, err := client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).SetHeaders(map[string]string{
		"api-key":    d.apiKey,
		"api-secret": d.apiSecret,
	}).Post("http://cdn.ddunyun.com/v1/jobs")
	if err != nil {
		return false
	}
	if !post.IsSuccessState() {
		return false
	}

	if resp.Code != 0 {
		facades.Log.Error("CDN[盾云][目录缓存刷新失败] " + resp.Message)
		return false
	}

	return true
}
