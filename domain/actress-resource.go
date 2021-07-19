package domain

type ActressResourceService interface {
	GetUrl() string
	SetNextPage()
	GetActressesFromResourceUrl() ([]ActressResourceInfo, error)
}

type ActressResourceInfo struct {
	SubUrlPath string
	Name       string
}
