package post

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/post"
	"github.com/donghquinn/blog_back_go/utils"
)

// 전체 포스트 가져오기 - 페이징
func GetPostController(res http.ResponseWriter, req *http.Request) {
	blogId := req.Header.Get("BlogId")

	if blogId == "" {
		log.Printf("[GET_POST] BlogId header is missing")
		response.Response(res, response.CommonResponseWithMessage{
			Status:  http.StatusBadRequest,
			Code:    "01",
			Message: "No blog Id Provided",
		})

		return
	}

	log.Printf("[DEBUGGING] blogId: %s", blogId)

	page := req.URL.Query().Get("page")
	size := req.URL.Query().Get("size")
	tag := req.URL.Query().Get("tag")
	category := req.URL.Query().Get("category")
	isPinned := req.URL.Query().Get("pin")

	log.Printf("[DEBUGGING] url queries - page:%s, size: %s, tag: %s", page, size, tag)

	offset, limit := utils.ParsePaginationParams(page, size)

	unpinnedQueryResult, queryErr := QueryUnpinnedPostData(blogId, isPinned, category, tag, limit, offset)

	if queryErr != nil {
		log.Printf("[POST_EDIT] Query Post Data Error: %v", queryErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  http.StatusInternalServerError,
			Code:    "02",
			Message: "Query Post Data Error",
			Result:  false,
		})

		return
	}

	totalCount, countErr := GetTotalPostCount(blogId, isPinned, category, tag)

	if countErr != nil {
		log.Printf("[POST_EDIT] Query Post Data Error: %v", queryErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  http.StatusInternalServerError,
			Code:    "03",
			Message: "Query Post Data Error",
			Result:  false,
		})

		return
	}
	// 이름 디코딩 위해

	response.Response(res, types.ResponsePostListType{
		Status:    http.StatusOK,
		Code:      "01",
		Message:   "SUCCESS",
		Result:    true,
		PostList:  unpinnedQueryResult,
		PostCount: totalCount,
		Page:      page,
		Size:      size,
	})

}
