package usecase

import (
	"bytes"
	"face-parsing/domain"
	"fmt"
	"image"
	"image/jpeg"
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
	name        string
	imageURL    string
	baseUrl     string
	savePath    string
	imageSource []byte
	cropImage   []byte
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
	postDetectResponse, err := a.FaceService.PostDetect(a.imageSource)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s detect failed", actressName))
	}
	if postDetectResponse.FaceNum == 0 {
		return errors.New(fmt.Sprintf("%s can't detect a face", actressName))
	}

	if err := a.CropImage(postDetectResponse.Faces[0].FaceRectangle); err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s can't crop image. err: %+v", actressName, err))
	}

	return nil
}

func (f *ActressStore) CropImage(faceRectangle domain.FaceRectangle) error {
	img, _, err := image.Decode(bytes.NewReader(f.imageSource))
	if err != nil {
		return errors.Wrap(err, "decode byte failed")
	}

	cropImage, err := cutter.Crop(img, cutter.Config{
		Width:  img.Bounds().Dx(),
		Height: faceRectangle.Top + faceRectangle.Height + cropHeightOffset,
	})
	if err != nil {
		return errors.Wrap(err, "crop image failed")
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, cropImage, &jpeg.Options{100}); err != nil {
		return errors.Wrap(err, "encode image failed")
	}

	f.cropImage = buf.Bytes()

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

func (f *ActressStore) DownloadImage() error {
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

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(response.Body); err != nil {
		return errors.Wrap(err, "decode file to image failed")
	}

	f.imageSource = buf.Bytes()

	log.Info("success. cost time: ", time.Since(startTime))
	return nil
}

func (f *ActressStore) GetImage() []byte {
	return f.imageSource
}

func (f *ActressStore) GetCropImage() []byte {
	return f.cropImage
}
