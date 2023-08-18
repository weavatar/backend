package imagecheck

type Driver interface {
	Check(url string) (bool, error)
}
