package cdn

import (
	"strings"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/spf13/cast"
)

// CDN CDN操作类
type CDN struct {
	Driver []Driver
}

var internalCDN = &CDN{}

// NewCDN 创建CDN实例
func NewCDN() CDN {
	if len(internalCDN.Driver) > 0 {
		return *internalCDN
	}

	driver := facades.Config().GetString("cdn.driver", "starshield")
	drivers := strings.Split(driver, ",")

	for _, d := range drivers {
		config := cast.ToStringMapString(facades.Config().Get("cdn." + d))
		switch d {
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
		case "upyun":
			internalCDN.Driver = append(internalCDN.Driver, &UpYun{
				Token: config["token"],
			})
		case "ddun":
			internalCDN.Driver = append(internalCDN.Driver, &DDun{
				apiKey:    config["api_key"],
				apiSecret: config["api_secret"],
			})
		case "anycast":
			internalCDN.Driver = append(internalCDN.Driver, &AnyCast{
				apiKey:    config["api_key"],
				apiSecret: config["api_secret"],
			})
		}
	}

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
