package domain

type ActressStoreService interface {
	DownloadImage()
	DeleteImage() error
	GetImagePath() string
	SetActress(name string, imageUrlSubPath string)
}
