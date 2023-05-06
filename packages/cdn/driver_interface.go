package cdn

type Driver interface {
	RefreshUrl(urls []string) bool
	RefreshPath(paths []string) bool
	GetUsage(domain, startTime, endTime string) uint
}
