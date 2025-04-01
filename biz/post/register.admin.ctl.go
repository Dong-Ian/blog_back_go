package post

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/middlewares"
	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/post"
	"github.com/donghquinn/blog_back_go/utils"
)

// 게시글 등록
func RegisterPostController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "JWT Verifying Error",
			Result:  false,
		})

		return
	}

	var registerPostRequest types.RegisterPostRequest

	parseErr := utils.DecodeBody(req, &registerPostRequest)

	if parseErr != nil {
		log.Printf("[POST_REGIST] Parsing Request Body: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Parsing Request Body",
			Result:  false,
		})

		return
	}

	postSeq, insertErr := InsertPostData(registerPostRequest, user.UserId, user.BlogId)

	if insertErr != nil {
		log.Printf("[POST_REGIST] Insert Post Data Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Insert Post Data Error",
			Result:  false,
		})

		return
	}

	response.Response(res, types.ResponsePostRegisterType{
		Status:  http.StatusOK,
		Code:    "01",
		Message: "SUCCESS",
		Result:  true,
		PostSeq: postSeq,
	})
}

func DeletePostController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "JWT Verifying Error",
			Result:  false,
		})

		return
	}

	var deleteRequest types.DeletePostRequest

	parseErr := utils.DecodeBody(req, &deleteRequest)

	if parseErr != nil {
		log.Printf("[DELETE] Parse Delete Request Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Delete Post Error",
			Result:  false,
		})

		return
	}

	deleteErr := DeletePost(deleteRequest.PostSeq, user.BlogId)

	if deleteErr != nil {
		log.Printf("[DELETE] Parse Delete Request Error: %v", deleteErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Delete Post Error",
			Result:  false,
		})

		return
	}

	response.Response(res, response.CommonResponseWithMessage{
		Status:  http.StatusOK,
		Code:    "01",
		Message: "SUCESS",
		Result:  true,
	})
}

// 고정 게시글 데이터 업데이트
func UpdatePinPostController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "JWT Verifying Error",
			Result:  false,
		})

		return
	}

	var updatePinRequest types.UpdatePinRequest

	parseErr := utils.DecodeBody(req, &updatePinRequest)

	if parseErr != nil {
		log.Printf("[PIN] Parse Pin Request Error: %v", parseErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Update Pin Error",
			Result:  false,
		})

		return
	}

	updateErr := UpdatePinPost(updatePinRequest, user.BlogId)

	if updateErr != nil {
		log.Printf("[PIN] Update Pin Error: %v", updateErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Update Pin Error",
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

// 고정 게시글 해제 데이터 업데이트
func UpdateUnPinPostController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "JWT Verifying Error",
			Result:  false,
		})

		return
	}

	var updateUnPinRequest types.UpdatePinRequest

	parseErr := utils.DecodeBody(req, &updateUnPinRequest)

	if parseErr != nil {
		log.Printf("[PIN] Parse Un-Pin Request Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Parse Un-Pin Request Error",
			Result:  false,
		})
		return
	}

	updateErr := UpdateUnPinPost(updateUnPinRequest, user.BlogId)

	if updateErr != nil {
		log.Printf("[PIN] Parse Un-Pin Request Error: %v", updateErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Update Un-Pin Error",
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
