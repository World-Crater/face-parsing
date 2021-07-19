package repo

import (
	"face-parsing/domain"
	"fmt"

	"github.com/antchfx/htmlquery"
	"github.com/pkg/errors"
)

type ActressResource struct {
	BaseUrl string
	Page    uint
}

func NewActressResourceUrl(baseUrl string) domain.ActressResourceService {
	actressResourceUrl := ActressResource{baseUrl, 1}
	return &actressResourceUrl
}

func (a ActressResource) GetUrl() string {
	return fmt.Sprintf("%s?page=%d", a.BaseUrl, a.Page)
}

func (a *ActressResource) SetNextPage() {
	if a.Page == 0 {
		a.Page = 1
	}
	a.Page += 1
}

func (a *ActressResource) GetActressesFromResourceUrl() ([]domain.ActressResourceInfo, error) {
	doc, err := htmlquery.LoadURL(a.GetUrl())
	nodes, err := htmlquery.QueryAll(doc, "//*[@id=\"main-area\"]/section/table/tbody/tr[*]/td[1]/a/img")
	if err != nil {
		return nil, errors.Wrap(err, "query html fail")
	}

	actressResourceInfo := []domain.ActressResourceInfo{}
	for _, value := range nodes {
		actressResourceInfo = append(actressResourceInfo, domain.ActressResourceInfo{SubUrlPath: value.Attr[0].Val, Name: value.Attr[1].Val})
	}
	return actressResourceInfo, nil
}
