package deliver

import (
	"face-parsing/domain"
	"face-parsing/usecase"
	"fmt"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ImageDeliver struct {
	domain.ActressStoreService
	domain.FaceService
}

func NewImageDeliver(actressStoreService domain.ActressStoreService, faceService domain.FaceService) *ImageDeliver {
	return &ImageDeliver{
		actressStoreService,
		faceService,
	}
}

func (i ImageDeliver) CropToFixR18Image(limit, offsetMax int) error {
	getInfosAllActressesResponse, err := i.FaceService.GetInfosAllActresses(limit, offsetMax)
	if err != nil {
		return errors.Wrap(err, "get infos all actresses fail")
	}

	for _, actress := range getInfosAllActressesResponse {
		i.ActressStoreService.SetActressWithImageURL(actress.Name, actress.Preview)

		log.Infof("fix r18. name: %s", actress.Name)

		if err := i.ActressStoreService.DownloadImage(); err != nil {
			log.Errorf("pass. download image fail. error: %+v\n", err)
			continue
		}

		if err := i.ActressStoreService.DetectImageThenCropImage(actress.Name); err != nil {
			log.Errorf("pass. detect image then crop image fail. error: %+v\n", err)
			continue
		}

		if err := i.FaceService.PutInfo(actress.ID, i.ActressStoreService.GetImagePath()); err != nil {
			log.Errorf("pass. update info fail. error: %+v\n", err)
			continue
		}
	}

	return nil
}

func (i ImageDeliver) AddActress(
	actressResourceUrl domain.ActressResourceService,
	actressValidator *usecase.ActressValidator,
	actressStore domain.ActressStoreService,
	faceService domain.FaceService,
	actressUploadCount, actressUploadCountMax int,
) error {
	for {
		for {
			log.Info("current page: ", actressResourceUrl.GetUrl())

			getActressesFromResourceUrl, err := actressResourceUrl.GetActressesFromResourceUrl()
			if err != nil {
				log.Error("get actresses from resource url fail", err)
				return errors.Wrap(err, "get actresses from resource URL failed")
			}

			for _, resourceInfoFromUrl := range getActressesFromResourceUrl {
				if actressValidator.IsInActressList(resourceInfoFromUrl.GetFormatName()) {
					log.Warn(fmt.Sprintf("%s in actress list", resourceInfoFromUrl.GetFormatName()))
					continue
				}

				if actressValidator.IsInCantDetectList(resourceInfoFromUrl.GetFormatName(), resourceInfoFromUrl.SubUrlPath) {
					log.Warn(fmt.Sprintf("%s in can't detect list", resourceInfoFromUrl.GetFormatName()))
					continue
				}

				actressStore.SetActress(resourceInfoFromUrl.GetFormatName(), resourceInfoFromUrl.SubUrlPath)

				if err := actressStore.DownloadImage(); err != nil {
					log.Error("pass. download image fail. error: ", err)
					continue
				}

				if err := actressStore.DetectImageThenCropImage(resourceInfoFromUrl.GetFormatName()); err != nil {
					log.Warn(fmt.Sprintf("pass. detect then crop image fail. error: %+v", err))
					actressValidator.AddToCantDetectList(resourceInfoFromUrl.GetFormatName(), resourceInfoFromUrl.SubUrlPath)
					continue
				}

				postInfosResponse, err := faceService.PostInfo(actressStore.GetImagePath(), domain.Actress{Name: resourceInfoFromUrl.GetFormatName()})
				if err != nil {
					log.Error("post info fail. error: ", err)
					return errors.Wrap(err, "post info fail")
				}

				if _, err = faceService.PostFace(actressStore.GetImagePath(), postInfosResponse.ID); err != nil {
					log.Error("post face fail. error: ", err)
					log.Info("delete info: ", resourceInfoFromUrl.GetFormatName())
					if err := faceService.DeleteInfo(postInfosResponse.ID); err != nil {
						log.Error("delete info fail. error: ", err)
						return errors.Wrap(err, "delete info fail")
					}
					continue
				}

				if err := actressStore.DeleteImage(); err != nil {
					log.Fatal("delete image fail. error: ", err)
				}

				log.Info(fmt.Sprintf("upload %s to face service", resourceInfoFromUrl.GetFormatName()))
				actressUploadCount++
			}
			actressResourceUrl.SetNextPage()
			actressValidator.UpdateActressInfos()

			if actressUploadCount >= actressUploadCountMax {
				log.Info("uploaded 1000 actresses")
				break
			}
		}
		time.Sleep(time.Hour * 24)
		actressUploadCount = 0
	}
}
