package cdn

type Driver interface {
	RefreshUrl(urls []string) bool
	RefreshPath(paths []string) bool
}
