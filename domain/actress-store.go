package domain

type ActressStoreService interface {
	DownloadImage() error
	DeleteImage() error
	GetImagePath() string
	SetActress(name string, imageUrlSubPath string)
	CropImage(imagePath string, faceRectangle FaceRectangle) error
}
