package domain

import "time"

type GetInfosResponse struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
	Rows   []Actress
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

type PostSearchResponse []struct {
	ActressWithRecognition
}

type PostInfosResponse struct {
	ID string `json:"id"`
}

type PostFaceResponse struct {
	FacesetToken string `json:"facesetToken"`
	FaceToken    string `json:"faceToken"`
}

type FaceService interface {
	GetInfos(limit uint, offset uint) (*GetInfosResponse, error)
	GetInfosAllActresses() ([]Actress, error)
}
