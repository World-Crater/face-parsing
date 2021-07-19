package usecase

import (
	"face-parsing/domain"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ActressStore struct {
	name            string
	imageUrlSubPath string
	baseUrl         string
	savePath        string
	domain.ActressResourceService
}

func NewActressStore(actressResourceService domain.ActressResourceService, savePath string, baseUrl string) domain.ActressStoreService {
	return &ActressStore{ActressResourceService: actressResourceService, savePath: savePath, baseUrl: baseUrl}
}

func (f *ActressStore) SetActress(name string, imageUrlSubPath string) {
	f.name = name
	f.imageUrlSubPath = imageUrlSubPath
}

func (f ActressStore) DownloadImage() {
	log.Info("download image: ", f.name, f.imageUrlSubPath)
	startTime := time.Now()

	req, e := http.NewRequest("GET", fmt.Sprintf("%s/%s", f.baseUrl, f.imageUrlSubPath), nil)
	req.Header.Set("Referer", f.GetUrl())
	client := &http.Client{}

	// don't worry about errors
	response, e := client.Do(req)

	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(fmt.Sprintf("%s/%s.jpg", f.savePath, f.name))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("success. cost time: ", time.Since(startTime))
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
