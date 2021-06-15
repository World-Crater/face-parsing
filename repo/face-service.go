package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type FaceService struct {
	URL string
}

type Actress struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Romanization interface{} `json:"romanization"`
	Detail       interface{} `json:"detail"`
	Preview      string      `json:"preview"`
	Createdat    time.Time   `json:"createdat"`
	Updatedat    time.Time   `json:"updatedat"`
}

type ActressWithRecognition struct {
	Actress
	Token                 string  `json:"token"`
	RecognitionPercentage float64 `json:"recognitionPercentage"`
}

type GetInfosResponse struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
	Rows   []Actress
}

type PostSearchResponse []struct {
	ActressWithRecognition
}

type PostInfosResponse struct {
	ID string `json:"id"`
}

func (service *FaceService) GetInfos() GetInfosResponse {
	res, err := http.Get(fmt.Sprintf("%s/faces/infos", service.URL))
	if err != nil {
		//yorktodo
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//yorktodo
	}
	var infosResponse GetInfosResponse
	json.Unmarshal(body, &infosResponse)
	return infosResponse
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

func (service *FaceService) PostSearch(filePath string) (PostSearchResponse, error) {
	payload, writer, err := service.createImagePayload(filePath, "image")
	if err != nil {
		//yorktodo
	}
	if err := writer.Close(); err != nil {
		//yorktodo
	}
	res, err := http.Post(fmt.Sprintf("%s/faces/search", service.URL), writer.FormDataContentType(), payload)
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
	var postSearchResponse PostSearchResponse
	json.Unmarshal(body, &postSearchResponse)
	return postSearchResponse, nil
}

func (service *FaceService) PostInfo(filePath string, actress Actress) (PostInfosResponse, error) {
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

	res, err := http.Post(fmt.Sprintf("%s/faces/info", service.URL), writer.FormDataContentType(), payload)
	if err != nil {
		fmt.Println("yorkyork", err)
		//yorktodo
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusInternalServerError {
		//yorktodo
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//yorktodo
	}
	var postInfosResponse PostInfosResponse
	json.Unmarshal(body, &postInfosResponse)
	return postInfosResponse, nil
}
