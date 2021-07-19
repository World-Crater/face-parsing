package repo

import (
	"bytes"
	"encoding/json"
	"face-parsing/domain"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type FaceService struct {
	Url string
}

func (service *FaceService) GetInfos(limit uint, offset uint) (*domain.GetInfosResponse, error) {
	if limit == 0 {
		return nil, errors.New("require limit")
	}

	res, err := http.Get(fmt.Sprintf("%s/faces/infos?limit=%d&offset=%d", service.Url, limit, offset))
	if err != nil {
		//yorktodo
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//yorktodo
	}
	var infosResponse domain.GetInfosResponse
	json.Unmarshal(body, &infosResponse)
	return &infosResponse, nil
}

func (service *FaceService) GetInfosAllActresses() ([]domain.Actress, error) {
	actresses := []domain.Actress{}
	offset := 0
	offsetIncrease := 1000

	for {
		GetInfosResponse, err := service.GetInfos(1000, uint(offset))
		if err != nil {
			return nil, errors.Wrap(err, "get infos fail")
		}
		actresses = append(actresses, GetInfosResponse.Rows...)
		offset = offset + offsetIncrease
		if offset >= GetInfosResponse.Count {
			log.Info("current actress infos count", len(actresses))
			return actresses, nil
		}
	}
}

func (service *FaceService) createImagePayload(filePath string, keyName string) (*bytes.Buffer, *multipart.Writer, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "file open fail")
	}
	defer file.Close()
	part1, err := writer.CreateFormFile(keyName, filepath.Base(filePath))
	_, err = io.Copy(part1, file)
	if err != nil {
		return nil, nil, errors.Wrap(err, "io copy fail")
	}
	return payload, writer, nil
}

func (service *FaceService) PostSearch(filePath string) (domain.PostSearchResponse, error) {
	payload, writer, err := service.createImagePayload(filePath, "image")
	if err != nil {
		//yorktodo
	}
	if err := writer.Close(); err != nil {
		//yorktodo
	}
	res, err := http.Post(fmt.Sprintf("%s/faces/search", service.Url), writer.FormDataContentType(), payload)
	if err != nil {
		//yorktodo
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("internal server error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//yorktodo
	}
	var postSearchResponse domain.PostSearchResponse
	json.Unmarshal(body, &postSearchResponse)
	return postSearchResponse, nil
}

func (service *FaceService) PostInfo(filePath string, actress domain.Actress) (*domain.PostInfosResponse, error) {
	if actress.Name == "" {
		//yorktodo
	}

	payload, writer, err := service.createImagePayload(filePath, "preview")
	if err != nil {
		//yorktodo
	}
	_ = writer.WriteField("name", actress.Name)
	if actress.Romanization != "" && actress.Romanization != nil {
		//yorktodo
		_ = writer.WriteField("romanization", actress.Romanization.(string))
	}
	if actress.Detail != "" && actress.Detail != nil {
		//yorktodo
		_ = writer.WriteField("detail", actress.Detail.(string))
	}
	if err := writer.Close(); err != nil {
		//yorktodo
	}

	res, err := http.Post(fmt.Sprintf("%s/faces/info", service.Url), writer.FormDataContentType(), payload)
	if err != nil {
		fmt.Println("yorkyork", err)
		//yorktodo
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("internal error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//yorktodo
	}
	var postInfosResponse domain.PostInfosResponse
	json.Unmarshal(body, &postInfosResponse)
	return &postInfosResponse, nil
}

func (service *FaceService) PostFace(filePath string, infoId string) (*domain.PostFaceResponse, error) {
	if infoId == "" {
		return nil, errors.New("require infoId")
	}

	payload, writer, err := service.createImagePayload(filePath, "image")
	if err != nil {
		//yorktodo
	}

	_ = writer.WriteField("infoId", infoId)

	if err := writer.Close(); err != nil {
		//yorktodo
	}

	res, err := http.Post(fmt.Sprintf("%s/faces/face", service.Url), writer.FormDataContentType(), payload)
	if err != nil {
		fmt.Println("yorkyork", err)
		//yorktodo
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("internal error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//yorktodo
	}
	var postFaceResponse domain.PostFaceResponse
	json.Unmarshal(body, &postFaceResponse)
	return &postFaceResponse, nil
}
