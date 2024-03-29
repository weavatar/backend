package cdn

import (
	"strings"
	"sync"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/spf13/cast"
)

// CDN CDN操作类
type CDN struct {
	Driver []Driver
}

var internalCDN = &CDN{}
var once sync.Once

// NewCDN 创建CDN实例
func NewCDN() CDN {
	once.Do(func() {
		driver := facades.Config().GetString("cdn.driver", "ctyun")
		drivers := strings.Split(driver, ",")

		for _, d := range drivers {
			config := cast.ToStringMapString(facades.Config().Get("cdn." + d))
			switch d {
			case "ctyun":
				internalCDN.Driver = append(internalCDN.Driver, &CTYun{
					AppID:       config["app_id"],
					AppSecret:   config["app_secret"],
					ApiEndpoint: "https://open.ctcdn.cn",
				})
			case "wangsu":
				internalCDN.Driver = append(internalCDN.Driver, &WangSu{
					AccessKey: config["access_key"],
					SecretKey: config["secret_key"],
				})
			case "starshield":
				internalCDN.Driver = append(internalCDN.Driver, &StarShield{
					AccessKey:  config["access_key"],
					SecretKey:  config["secret_key"],
					InstanceID: config["instance_id"],
					ZoneID:     config["zone_id"],
				})
			case "baishan":
				internalCDN.Driver = append(internalCDN.Driver, &BaiShan{
					Token: config["token"],
				})
			case "huawei":
				internalCDN.Driver = append(internalCDN.Driver, &HuaWei{
					AccessKey: config["access_key"],
					SecretKey: config["secret_key"],
				})
			case "yundun":
				internalCDN.Driver = append(internalCDN.Driver, &YunDun{
					UserName: config["username"],
					PassWord: config["password"],
				})
			case "kuocai":
				internalCDN.Driver = append(internalCDN.Driver, &KuoCai{
					UserName: config["username"],
					PassWord: config["password"],
				})
			case "cloudflare":
				internalCDN.Driver = append(internalCDN.Driver, &CloudFlare{
					Key:    config["key"],
					Email:  config["email"],
					ZoneID: config["zone_id"],
				})
			}
		}
	})

	return *internalCDN
}

// RefreshUrl 刷新URL
func (cdn *CDN) RefreshUrl(urls []string) bool {
	for _, driver := range cdn.Driver {
		driver.RefreshUrl(urls)
	}

	return true
}

// RefreshPath 刷新路径
func (cdn *CDN) RefreshPath(paths []string) bool {
	for _, driver := range cdn.Driver {
		driver.RefreshPath(paths)
	}

	return true
}

// GetUsage 获取CDN使用情况
func (cdn *CDN) GetUsage(domain string, startTime, endTime carbon.Carbon) uint {
	var usage uint
	for _, driver := range cdn.Driver {
		usage += driver.GetUsage(domain, startTime, endTime)
	}

	return usage
}
