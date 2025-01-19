package admincontrollers

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/dto"
	post "github.com/donghquinn/blog_back_go/libraries/post/admin"
	"github.com/donghquinn/blog_back_go/middlewares"
	types "github.com/donghquinn/blog_back_go/types/admin/posts"
	"github.com/donghquinn/blog_back_go/utils"
)

// 게시글 등록
func RegisterPostController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", nil)

		return
	}

	var registerPostRequest types.RegisterPostRequest

	parseErr := utils.DecodeBody(req, &registerPostRequest)

	if parseErr != nil {
		dto.SetErrorResponse(res, 402, "02", "Parsing Request Body", parseErr)
		return
	}

	postSeq, insertErr := post.InsertPostData(registerPostRequest, user.UserId, user.BlogId)

	if insertErr != nil {
		dto.SetErrorResponse(res, 403, "03", "Insert Post Data Error", insertErr)
		return
	}

	dto.SetPostRegisterResponse(res, 200, "01", postSeq)
}

func DeletePostController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", nil)

		return
	}

	var deleteRequest types.DeletePostRequest

	parseErr := utils.DecodeBody(req, &deleteRequest)

	if parseErr != nil {
		log.Printf("[DELETE] Parse Delete Request Error: %v", parseErr)
		dto.SetErrorResponse(res, 402, "02", "Delete Post Error", parseErr)
		return
	}

	deleteErr := post.DeletePost(deleteRequest.PostSeq, user.BlogId)

	if deleteErr != nil {
		dto.SetErrorResponse(res, 403, "03", "Delete Post Error", deleteErr)
		return
	}

	dto.SetResponse(res, 200, "01")
}

// 고정 게시글 데이터 업데이트
func UpdatePinPostController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", nil)

		return
	}

	var updatePinRequest types.UpdatePinRequest

	parseErr := utils.DecodeBody(req, &updatePinRequest)

	if parseErr != nil {
		log.Printf("[PIN] Parse Pin Request Error: %v", parseErr)
		dto.SetErrorResponse(res, 402, "02", "Update Pin Error", parseErr)
		return
	}

	updateErr := post.UpdatePinPost(updatePinRequest, user.BlogId)

	if updateErr != nil {
		dto.SetErrorResponse(res, 403, "03", "Update Pin Error", updateErr)
		return
	}

	dto.SetResponse(res, 200, "01")
}

// 고정 게시글 해제 데이터 업데이트
func UpdateUnPinPostController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", nil)

		return
	}

	var updateUnPinRequest types.UpdatePinRequest

	parseErr := utils.DecodeBody(req, &updateUnPinRequest)

	if parseErr != nil {
		log.Printf("[PIN] Parse Un-Pin Request Error: %v", parseErr)
		dto.SetErrorResponse(res, 402, "02", "Update Un-Pin Error", parseErr)
		return
	}

	updateErr := post.UpdateUnPinPost(updateUnPinRequest, user.BlogId)

	if updateErr != nil {
		dto.SetErrorResponse(res, 403, "03", "Update Un-Pin Error", updateErr)
		return
	}

	dto.SetResponse(res, 200, "01")
}
