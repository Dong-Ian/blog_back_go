package post

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/dto"
	types "github.com/donghquinn/blog_back_go/types/post"
	"github.com/donghquinn/blog_back_go/utils"
)

func GetCategoryController(res http.ResponseWriter, req *http.Request) {
	var getCategoryListRequest types.GetPostListRequest

	parseErr := utils.DecodeBody(req, &getCategoryListRequest)

	if parseErr != nil {
		dto.SetErrorResponse(res, 401, "01", "Parse Request Body Error", parseErr)
		return
	}

	categoryList, categoryErr := GetAllCategoryList(getCategoryListRequest.BlogId)

	if categoryErr != nil {
		dto.SetErrorResponse(res, 400, "01", "Get All Category Error", categoryErr)
		return
	}

	dto.SetCategoryResponse(res, 200, "01", categoryList)
}
