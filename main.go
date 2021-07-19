package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"face-parsing/domain"
	"face-parsing/repo"
	"face-parsing/usecase"
)

const BASE_URL, SAVE_PATH = "http://www.minnano-av.com", "./images"

func main() {
	// repo
	faceService := repo.FaceService{
		Url: "http://face-service:3000",
	}
	actressResourceUrl := repo.NewActressResourceUrl("http://www.minnano-av.com/actress_list.php")

	// usecase
	actressValidator := usecase.NewActressValidator(&faceService)
	actressStore := usecase.NewActressStore(actressResourceUrl, SAVE_PATH, BASE_URL)

	// deliver
	for {
		for i := 0; i < 10; i++ {
			log.Info("current page: ", actressResourceUrl.GetUrl())

			getActressesFromResourceUrl, err := actressResourceUrl.GetActressesFromResourceUrl()
			if err != nil {
				log.Error("get actresses from resource url fail", err)
				return
			}

			for _, resourceInfoFromUrl := range getActressesFromResourceUrl {
				if actressValidator.IsInActressList(resourceInfoFromUrl.Name) {
					log.Warn(fmt.Sprintf("%s in actress list", resourceInfoFromUrl.Name))
					continue
				}
				actressStore.SetActress(resourceInfoFromUrl.Name, resourceInfoFromUrl.SubUrlPath)
				actressStore.DownloadImage()
				_, err := faceService.PostSearch(actressStore.GetImagePath())
				if err != nil {
					log.Warn(fmt.Sprintf("%s can't detect a face", resourceInfoFromUrl.Name))
					continue
				}
				postInfosResponse, err := faceService.PostInfo(actressStore.GetImagePath(), domain.Actress{Name: resourceInfoFromUrl.Name})
				if err != nil {
					log.Error("post info fail. error: ", err)
					return
				}
				_, err = faceService.PostFace(actressStore.GetImagePath(), postInfosResponse.ID)
				if err != nil {
					log.Error("post face fail. error: ", err)
					return
				}
				log.Info(fmt.Sprintf("upload %s to face service", resourceInfoFromUrl.Name))
				if err := actressStore.DeleteImage(); err != nil {
					log.Fatal("delete image fail. error: ", err)
				}
			}
			actressResourceUrl.SetNextPage()
			actressValidator.UpdateActressInfos()
		}
		time.Sleep(time.Hour * 24)
	}
}
