package cdn

import (
	"sync"

	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
)

// CDN CDN操作类
type CDN struct {
	Driver Driver
}

// once 单例模式
var once sync.Once

// internalCDN 内部使用的 SMS 对象
var internalCDN *CDN

// NewCDN 单例模式获取
func NewCDN() *CDN {
	driver := facades.Config.Get("cdn.driver")
	config := make(map[string]string)

	switch driver {
	case "ddun":
		config = cast.ToStringMapString(facades.Config.Get("cdn.ddun"))
		once.Do(func() {
			internalCDN = &CDN{
				Driver: &DDun{
					apiKey:    config["api_key"],
					apiSecret: config["api_secret"],
				},
			}
		})
	case "yundun":
		config = cast.ToStringMapString(facades.Config.Get("cdn.yundun"))
		once.Do(func() {
			internalCDN = &CDN{
				Driver: &YunDun{
					UserName: config["username"],
					PassWord: config["password"],
				},
			}
		})
	}

	return internalCDN
}

// RefreshUrl 刷新URL
func (cdn *CDN) RefreshUrl(urls []string) bool {
	return cdn.Driver.RefreshUrl(urls)
}

// RefreshPath 刷新路径
func (cdn *CDN) RefreshPath(paths []string) bool {
	return cdn.Driver.RefreshPath(paths)
}
