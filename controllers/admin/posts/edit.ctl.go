package admincontrollers

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/dto"
	post "github.com/donghquinn/blog_back_go/libraries/post/admin"
	"github.com/donghquinn/blog_back_go/middlewares"
	types "github.com/donghquinn/blog_back_go/types/admin/posts"
	"github.com/donghquinn/blog_back_go/utils"
)

func EditPostController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", nil)

		return
	}

	var editPostRequest types.EditPostRequest

	parseErr := utils.DecodeBody(req, &editPostRequest)

	if parseErr != nil {
		dto.SetErrorResponse(res, 402, "02", "Parse Request Body Error", parseErr)
		return
	}

	editErr := post.EditPost(editPostRequest, user.UserId, user.BlogId)

	if editErr != nil {
		dto.SetErrorResponse(res, 403, "03", "Edit Post Data Error", editErr)
		return
	}

	dto.SetResponse(res, 200, "01")
}
