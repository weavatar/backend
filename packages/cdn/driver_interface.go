package cdn

import "github.com/goravel/framework/support/carbon"

type Driver interface {
	RefreshUrl(urls []string) bool
	RefreshPath(paths []string) bool
	GetUsage(domain string, startTime, endTime carbon.Carbon) uint
}
