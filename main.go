package main

import (
	"fmt"

	"github.com/spf13/viper"

	"face-parsing/deliver"
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
	cropToFixR18limit := viper.GetInt("CROP_TO_FIX_R18_LIMIT")
	cropToFixR18ImageOffsetMax := viper.GetInt("CROP_TO_FIX_R18_IMAGE_OFFSET_MAX")
	action := viper.GetString("ACTION")

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
	actressStore := usecase.NewActressStore(actressResourceUrl, &faceService, SAVE_PATH, BASE_URL)

	// deliver
	imageDeliver := deliver.NewImageDeliver(actressStore, &faceService)

	switch action {
	case "add_actresses":
		if err := imageDeliver.AddActress(
			actressResourceUrl,
			actressValidator,
			actressStore,
			&faceService,
			actressUploadCount,
			actressUploadCountMax,
		); err != nil {
			panic(err)
		}
	case "fix_r18":
		if err := imageDeliver.CropToFixR18Image(cropToFixR18limit, cropToFixR18ImageOffsetMax); err != nil {
			panic(err)
		}
	}
}
