package post

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/libraries/database"
	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/post"
	"github.com/donghquinn/blog_back_go/utils"
)

func GetImageUrl(res http.ResponseWriter, req *http.Request) {
	var getPostRequest types.GetPostByPostSeq

	err := utils.DecodeBody(req, &getPostRequest)

	if err != nil {
		log.Printf("[GET_IMAGE] Request Is Not Valid: %v", err)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "Request Is Not Valid",
			Result:  false,
		})

	}

	imageData, imageErr := GetImageData(getPostRequest.PostSeq)

	if imageErr != nil {
		log.Printf("[GET_IMAGE] Image Data Error: %v", err)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Image Data Error",
			Result:  false,
		})

		return
	}

	var urlArray []string

	// 게시글 URL 배열 만들기
	for _, data := range imageData {
		url, getErr := database.GetImageUrl(data.ObjectName, data.FileFormat)

		if getErr != nil {
			log.Printf("[GET_IMAGE] Get Presigned URL Error: %v", err)
			response.Response(res, response.CommonResponseWithMessage{
				Status:  403,
				Code:    "03",
				Message: "Get Presigned URL Error",
				Result:  false,
			})

			return
		}

		urlArray = append(urlArray, url.String())
	}

	// responseData := types.ViewImageUrl {
	// 	Urls: urlArray}

	response.Response(res, types.ResponseGetImageUrlType{
		Status:      http.StatusOK,
		Code:        "01",
		Message:     "Get Presigned URL Error",
		ImageResult: urlArray,
		Result:      true,
	})

}
