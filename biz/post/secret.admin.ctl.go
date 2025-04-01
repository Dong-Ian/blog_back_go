package post

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/middlewares"
	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/post"
	"github.com/donghquinn/blog_back_go/utils"
)

// 비공개 게시글로 변경
func ChangeToSecretPostController(res http.ResponseWriter, request *http.Request) {
	_, ok := middlewares.GetUserFromContext(request.Context())

	if !ok {
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "JWT Verifying Error",
			Result:  false,
		})

		return
	}

	var changeRequest types.ChangeToRequest

	parseErr := utils.DecodeBody(request, &changeRequest)

	if parseErr != nil {
		log.Printf("[CHANGE_SECRET] Parse Request Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Invalid Request Body",
			Result:  false,
		})

		return
	}

	changeErr := ChangeToSecretPost(changeRequest.PostSeq)

	if changeErr != nil {
		log.Printf("[CHANGE_SECRET] Change Secret Failed: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Change Secret Failed",
			Result:  false,
		})

		return
	}

	response.Response(res, response.CommonResponseWithMessage{
		Status:  http.StatusOK,
		Code:    "01",
		Message: "SUCCESS",
		Result:  true,
	})
}

// 비공개 게시글에서 공개 게시글로 변경
func ChangeToNotSecretPostController(res http.ResponseWriter, request *http.Request) {
	_, ok := middlewares.GetUserFromContext(request.Context())

	if !ok {
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "JWT Verifying Error",
			Result:  false,
		})
		return
	}

	var changeRequest types.ChangeToRequest

	parseErr := utils.DecodeBody(request, &changeRequest)

	if parseErr != nil {
		log.Printf("[CHANGE_NOT_SECRET] Parse Request Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Invalid Request Body",
			Result:  false,
		})

		return
	}

	changeErr := ChangeToNotSecretPost(changeRequest.PostSeq)

	if changeErr != nil {
		log.Printf("[CHANGE_SECRET] Change Rror: %v", changeErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Change Not Secret Failed",
			Result:  false,
		})

		return
	}

	response.Response(res, response.CommonResponseWithMessage{
		Status:  http.StatusOK,
		Code:    "01",
		Message: "SUCCESS",
		Result:  true,
	})
}
