package post

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/dto"
	"github.com/donghquinn/blog_back_go/libraries/database"
	types "github.com/donghquinn/blog_back_go/types/post"
	"github.com/donghquinn/blog_back_go/utils"
)

func GetImageUrl(res http.ResponseWriter, req *http.Request) {
	var getPostRequest types.GetPostByPostSeq

	err := utils.DecodeBody(req, &getPostRequest)

	if err != nil {
		dto.SetErrorResponse(res, 401, "01", "Request Is Not Valid", err)
	}

	imageData, imageErr := GetImageData(getPostRequest.PostSeq)

	if imageErr != nil {
		dto.SetErrorResponse(res, 402, "02", "Image Data Error", imageErr)
		return
	}

	var urlArray []string

	// 게시글 URL 배열 만들기
	for _, data := range imageData {
		url, getErr := database.GetImageUrl(data.ObjectName, data.FileFormat)

		if getErr != nil {
			dto.SetErrorResponse(res, 403, "03", "Get Presigned URL Error", getErr)
			return
		}

		urlArray = append(urlArray, url.String())
	}

	// responseData := types.ViewImageUrl {
	// 	Urls: urlArray}

	dto.SetImageUrlResponse(res, 200, "01", urlArray)
}
