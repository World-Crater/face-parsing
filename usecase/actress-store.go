package usecase

import (
	"face-parsing/domain"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/oliamb/cutter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	cropHeightOffset int = 20
)

type ActressStore struct {
	name     string
	imageURL string
	baseUrl  string
	savePath string
	domain.ActressResourceService
	domain.FaceService
}

func NewActressStore(actressResourceService domain.ActressResourceService, faceService domain.FaceService, savePath string, baseUrl string) domain.ActressStoreService {
	return &ActressStore{
		ActressResourceService: actressResourceService,
		savePath:               savePath,
		baseUrl:                baseUrl,
		FaceService:            faceService,
	}
}

func (a *ActressStore) DetectImageThenCropImage(actressName string) error {
	postDetectResponse, err := a.FaceService.PostDetect(a.GetImagePath())
	if err != nil {
		return errors.New(fmt.Sprintf("%s detect failed", actressName))
	}
	if postDetectResponse.FaceNum == 0 {
		return errors.New(fmt.Sprintf("%s can't detect a face", actressName))
	}

	if err := a.CropImage(a.GetImagePath(), postDetectResponse.Faces[0].FaceRectangle); err != nil {
		return errors.New(fmt.Sprintf("%s can't crop image. err: %+v", actressName, err))
	}

	return nil
}

func (f *ActressStore) CropImage(imagePath string, faceRectangle domain.FaceRectangle) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return errors.Wrap(err, "open file failed")
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return errors.Wrap(err, "decode file to image failed")
	}

	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:  img.Bounds().Dx(),
		Height: faceRectangle.Top + faceRectangle.Height + cropHeightOffset,
	})
	if err != nil {
		return errors.Wrap(err, "crop image failed")
	}

	ioWriter, err := os.Create(imagePath)
	if err != nil {
		return errors.Wrap(err, "create io writer failed")
	}
	if err := jpeg.Encode(ioWriter, croppedImg, &jpeg.Options{Quality: 100}); err != nil {
		return errors.Wrap(err, "save image failed")
	}

	return nil
}

func (f *ActressStore) SetActress(name string, imageUrlSubPath string) {
	f.name = name
	f.imageURL = fmt.Sprintf("%s/%s", f.baseUrl, imageUrlSubPath)
}

func (f *ActressStore) SetActressWithImageURL(name, url string) {
	f.name = name
	f.imageURL = url
}

func (f ActressStore) DownloadImage() error {
	log.Info(fmt.Sprintf("download %s image to %s", f.name, f.imageURL))
	startTime := time.Now()

	req, err := http.NewRequest("GET", f.imageURL, nil)
	if err != nil {
		return errors.Wrap(err, "new request fail")
	}
	req.Header.Set("Referer", f.GetUrl())
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request fail")
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(fmt.Sprintf("%s/%s.jpg", f.savePath, f.name))
	if err != nil {
		return errors.Wrap(err, "open file fail")
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return errors.Wrap(err, "copy file from body fail")
	}

	log.Info("success. cost time: ", time.Since(startTime))
	return nil
}

func (f ActressStore) DeleteImage() error {
	e := os.Remove(f.GetImagePath())
	if e != nil {
		return errors.Wrap(e, "remove fail")
	}
	return nil
}

func (f ActressStore) GetImagePath() string {
	return fmt.Sprintf("./images/%s.jpg", f.name)
}
