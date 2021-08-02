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

const BASE_URL, SAVE_PATH, ACTRESS_UPLOAD_COUNT_MAX = "http://www.minnano-av.com", "./images", 1000

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("fatal error config. %+v\n", err))
	}
	viper.AutomaticEnv()

	// repo
	faceService := repo.FaceService{
		Url: viper.GetString("FACE_SERVICE"),
	}
	actressResourceUrl := repo.NewActressResourceUrl("http://www.minnano-av.com/actress_list.php", viper.GetUint("RESOURCE_URL_PAGE"))

	// usecase
	actressValidator, err := usecase.NewActressValidator(&faceService)
	if err != nil {
		panic(fmt.Sprintf("new actress validator fail. %+v\n", err))
	}
	actressStore := usecase.NewActressStore(actressResourceUrl, SAVE_PATH, BASE_URL)

	// deliver
	actressUploadCount := viper.GetInt("UPLOAD_COUNT")
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
				_, err := faceService.PostSearch(actressStore.GetImagePath())
				if err != nil {
					log.Warn(fmt.Sprintf("%s can't detect a face", resourceInfoFromUrl.GetFormatName()))
					actressValidator.AddToCantDetectList(resourceInfoFromUrl.GetFormatName(), resourceInfoFromUrl.SubUrlPath)
					continue
				}
				postInfosResponse, err := faceService.PostInfo(actressStore.GetImagePath(), domain.Actress{Name: resourceInfoFromUrl.GetFormatName()})
				if err != nil {
					log.Error("post info fail. error: ", err)
					return
				}
				_, err = faceService.PostFace(actressStore.GetImagePath(), postInfosResponse.ID)
				if err != nil {
					log.Error("post face fail. error: ", err)
					return
				}
				if err := actressStore.DeleteImage(); err != nil {
					log.Fatal("delete image fail. error: ", err)
				}

				log.Info(fmt.Sprintf("upload %s to face service", resourceInfoFromUrl.GetFormatName()))
				actressUploadCount++
			}
			actressResourceUrl.SetNextPage()
			actressValidator.UpdateActressInfos()

			if actressUploadCount >= ACTRESS_UPLOAD_COUNT_MAX {
				log.Info("uploaded 1000 actresses")
				break
			}
		}
		time.Sleep(time.Hour * 24)
		actressUploadCount = 0
	}
}
