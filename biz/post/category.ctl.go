package post

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/post"
	"github.com/donghquinn/blog_back_go/utils"
)

func GetCategoryController(res http.ResponseWriter, req *http.Request) {
	var getCategoryListRequest types.GetPostListRequest

	parseErr := utils.DecodeBody(req, &getCategoryListRequest)

	if parseErr != nil {
		log.Printf("[CATEGORY] Parse Reqeust Error: %v", parseErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "Parse Request Body Error",
			Result:  false,
		})

		return
	}

	categoryList, categoryErr := GetAllCategoryList(getCategoryListRequest.BlogId)

	if categoryErr != nil {
		log.Printf("[CATEGORY] Parse Reqeust Error: %v", parseErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "Get All Category Error",
			Result:  false,
		})

		return
	}

	response.Response(res, types.ResponseCategoryResponseType{
		Status:       http.StatusOK,
		Code:         "01",
		Message:      "SUCCESS",
		CategoryList: categoryList,
		Result:       true,
	})
}
