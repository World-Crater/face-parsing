package usecase

import (
	"face-parsing/domain"

	"github.com/pkg/errors"
)

type ActressValidator struct {
	ActressList map[string]bool
	faceService domain.FaceService
}

func NewActressValidator(faceService domain.FaceService) ActressValidator {
	actressValidator := ActressValidator{
		faceService: faceService,
	}
	actressValidator.ActressList = make(map[string]bool)

	actressValidator.UpdateActressInfos()

	return actressValidator
}

func (a ActressValidator) IsInActressList(actress string) bool {
	return a.ActressList[actress]
}

func (a *ActressValidator) UpdateActressInfos() error {
	getInfosAllActressesResponse, err := a.faceService.GetInfosAllActresses()
	if err != nil {
		return errors.Wrap(err, "get infos all actresses fail")
	}

	for _, actress := range getInfosAllActressesResponse {
		a.ActressList[actress.Name] = true
	}
	return nil
}
