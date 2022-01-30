package repo

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestGetInfos(t *testing.T) {
// 	faceService := FaceService{
// 		"http://localhost:3000",
// 	}
// 	assert.Equal(t, faceService.GetInfos().Rows[0].Name, "AIKA", "equal name")
// }

// func TestPostSearch(t *testing.T) {
// 	faceService := FaceService{
// 		"http://localhost:3000",
// 	}
// 	result, err := faceService.PostSearch("./test-1.jpg")
// 	assert.Equal(t, err.Error(), "internal server error", "equal error")
// 	result, err = faceService.PostSearch("./test-2.jpg")
// 	assert.Equal(t, result[0].Name, "安齋らら", "equal name")
// }

// func TestPostInfo(t *testing.T) {
// 	faceService := FaceService{
// 		"http://localhost:3000",
// 	}
// 	result, _ := faceService.PostInfo("./test-1.jpg", Actress{
// 		Name: "test3",
// 	})
// 	fmt.Println(result)
// }
