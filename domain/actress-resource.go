package domain

import "strings"

type ActressResourceService interface {
	GetUrl() string
	SetNextPage()
	GetActressesFromResourceUrl() ([]ActressResourceInfo, error)
}

type ActressResourceInfo struct {
	SubUrlPath string
	Name       string
}

func (a ActressResourceInfo) GetFormatName() string {
	return strings.Replace(a.Name, "/", "-", -1)
}
