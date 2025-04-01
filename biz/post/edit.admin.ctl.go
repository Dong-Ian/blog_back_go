package post

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/post"

	"github.com/donghquinn/blog_back_go/middlewares"

	"github.com/donghquinn/blog_back_go/utils"
)

func EditPostController(res http.ResponseWriter, req *http.Request) {
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

	var editPostRequest types.EditPostRequest

	parseErr := utils.DecodeBody(req, &editPostRequest)

	if parseErr != nil {
		log.Printf("[POST_EDIT] Parse Request Body Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Parse Request Body Error",
			Result:  false,
		})

		return
	}

	editErr := EditPost(editPostRequest, user.UserId, user.BlogId)

	if editErr != nil {
		log.Printf("[POST_EDIT] Edit Post Data Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Edit Post Data Error",
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
