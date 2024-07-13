package cdn

import "github.com/goravel/framework/support/carbon"

type Driver interface {
	RefreshUrl(urls []string) error
	RefreshPath(paths []string) error
	GetUsage(domain string, startTime, endTime carbon.Carbon) (uint, error)
}
