package cdn

import "github.com/golang-module/carbon/v2"

type Driver interface {
	RefreshUrl(urls []string) bool
	RefreshPath(paths []string) bool
	GetUsage(domain string, startTime, endTime carbon.Carbon) uint
}
