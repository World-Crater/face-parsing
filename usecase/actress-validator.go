package usecase

import (
	"face-parsing/domain"
	"fmt"

	"github.com/pkg/errors"
)

type ActressValidator struct {
	actressList    map[string]bool
	faceService    domain.FaceService
	cantDetectList map[string]bool
}

func NewActressValidator(faceService domain.FaceService) (*ActressValidator, error) {
	actressValidator := ActressValidator{
		faceService: faceService,
	}
	actressValidator.actressList = make(map[string]bool)
	actressValidator.cantDetectList = make(map[string]bool)

	if err := actressValidator.UpdateActressInfos(); err != nil {
		return nil, errors.Wrap(err, "update actress infos fail")
	}

	return &actressValidator, nil
}

func (a ActressValidator) IsInActressList(actress string) bool {
	return a.actressList[actress]
}

func (a ActressValidator) IsInCantDetectList(name string, subUrlPath string) bool {
	return a.cantDetectList[a.getCantDetectListKey(name, subUrlPath)]
}

func (a ActressValidator) getCantDetectListKey(name string, subUrlPath string) string {
	return fmt.Sprintf("%s|%s", name, subUrlPath)
}

func (a *ActressValidator) AddToCantDetectList(name string, subUrlPath string) {
	a.cantDetectList[a.getCantDetectListKey(name, subUrlPath)] = true
}

func (a *ActressValidator) UpdateActressInfos() error {
	getInfosAllActressesResponse, err := a.faceService.GetInfosAllActresses()
	if err != nil {
		return errors.Wrap(err, "get infos all actresses fail")
	}

	for _, actress := range getInfosAllActressesResponse {
		a.actressList[actress.Name] = true
	}
	return nil
}
