package domain

type ActressStoreService interface {
	DownloadImage() error
	DeleteImage() error
	GetImagePath() string
	SetActress(name string, imageUrlSubPath string)
	SetActressWithImageURL(name, url string)
	DetectImageThenCropImage(actressName string) error
	CropImage(imagePath string, faceRectangle FaceRectangle) error
}
