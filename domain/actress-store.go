package domain

type ActressStoreService interface {
	DownloadImage() error
	SetActress(name string, imageUrlSubPath string)
	SetActressWithImageURL(name, url string)
	DetectImageThenCropImage(actressName string) error
	CropImage(faceRectangle FaceRectangle) error
	GetImage() []byte
	GetCropImage() []byte
}
