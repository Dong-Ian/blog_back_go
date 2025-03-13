package post

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/dto"
	"github.com/donghquinn/blog_back_go/middlewares"
	types "github.com/donghquinn/blog_back_go/types/admin/posts"
	"github.com/donghquinn/blog_back_go/utils"
)

// 비공개 게시글로 변경
func ChangeToSecretPostController(response http.ResponseWriter, request *http.Request) {
	_, ok := middlewares.GetUserFromContext(request.Context())

	if !ok {
		dto.SetErrorResponse(response, 401, "01", "JWT Validate Error", nil)
		return
	}

	var changeRequest types.ChangeToRequest

	parseErr := utils.DecodeBody(request, &changeRequest)

	if parseErr != nil {
		dto.SetErrorResponse(response, 402, "02", "Invalid Request Body", parseErr)
		return
	}

	changeErr := ChangeToSecretPost(changeRequest.PostSeq)

	if changeErr != nil {
		dto.SetErrorResponse(response, 403, "03", "Change Secret Failed", changeErr)
		return
	}

	dto.SetResponse(response, 200, "01")
}

// 비공개 게시글에서 공개 게시글로 변경
func ChangeToNotSecretPostController(response http.ResponseWriter, request *http.Request) {
	_, ok := middlewares.GetUserFromContext(request.Context())

	if !ok {
		dto.SetErrorResponse(response, 401, "01", "JWT Validate Error", nil)
		return
	}

	var changeRequest types.ChangeToRequest

	parseErr := utils.DecodeBody(request, &changeRequest)

	if parseErr != nil {
		dto.SetErrorResponse(response, 402, "02", "Invalid Request Body", parseErr)
		return
	}

	changeErr := ChangeToNotSecretPost(changeRequest.PostSeq)

	if changeErr != nil {
		dto.SetErrorResponse(response, 403, "03", "Change Not Secret Failed", changeErr)
		return
	}

	dto.SetResponse(response, 200, "01")
}
