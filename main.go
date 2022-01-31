package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"face-parsing/domain"
	"face-parsing/repo"
	"face-parsing/usecase"
)

const BASE_URL, SAVE_PATH = "http://www.minnano-av.com", "./images"

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("fatal error config. %+v\n", err))
	}
	viper.AutomaticEnv()

	actressUploadCountMax := viper.GetInt("ACTRESS_UPLOAD_COUNT_MAX")
	actressUploadCount := viper.GetInt("UPLOAD_COUNT")
	faceServiceURL := viper.GetString("FACE_SERVICE")
	resourceURLPage := viper.GetUint("RESOURCE_URL_PAGE")

	// repo
	faceService := repo.FaceService{
		Url: faceServiceURL,
	}
	actressResourceUrl := repo.NewActressResourceUrl("http://www.minnano-av.com/actress_list.php", resourceURLPage)

	// usecase
	actressValidator, err := usecase.NewActressValidator(&faceService)
	if err != nil {
		panic(fmt.Sprintf("new actress validator fail. %+v\n", err))
	}
	actressStore := usecase.NewActressStore(actressResourceUrl, SAVE_PATH, BASE_URL)

	// deliver
	for {
		for {
			log.Info("current page: ", actressResourceUrl.GetUrl())

			getActressesFromResourceUrl, err := actressResourceUrl.GetActressesFromResourceUrl()
			if err != nil {
				log.Error("get actresses from resource url fail", err)
				return
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

				postDetectResponse, err := faceService.PostDetect(actressStore.GetImagePath())
				if err != nil {
					log.Warn(fmt.Sprintf("%s detect failed", resourceInfoFromUrl.GetFormatName()))
					actressValidator.AddToCantDetectList(resourceInfoFromUrl.GetFormatName(), resourceInfoFromUrl.SubUrlPath)
					continue
				}
				if postDetectResponse.FaceNum == 0 {
					log.Warn(fmt.Sprintf("%s can't detect a face", resourceInfoFromUrl.GetFormatName()))
					actressValidator.AddToCantDetectList(resourceInfoFromUrl.GetFormatName(), resourceInfoFromUrl.SubUrlPath)
					continue
				}

				if err := actressStore.CropImage(actressStore.GetImagePath(), postDetectResponse.Faces[0].FaceRectangle); err != nil {
					log.Warn(fmt.Sprintf("%s can't crop image. err: %+v", resourceInfoFromUrl.GetFormatName(), err))
					actressValidator.AddToCantDetectList(resourceInfoFromUrl.GetFormatName(), resourceInfoFromUrl.SubUrlPath)
					continue
				}

				postInfosResponse, err := faceService.PostInfo(actressStore.GetImagePath(), domain.Actress{Name: resourceInfoFromUrl.GetFormatName()})
				if err != nil {
					log.Error("post info fail. error: ", err)
					return
				}

				if _, err = faceService.PostFace(actressStore.GetImagePath(), postInfosResponse.ID); err != nil {
					log.Error("post face fail. error: ", err)
					log.Info("delete info: ", resourceInfoFromUrl.GetFormatName())
					if err := faceService.DeleteInfo(postInfosResponse.ID); err != nil {
						log.Error("delete info fail. error: ", err)
						return
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
