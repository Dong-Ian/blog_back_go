package post

import (
	"log"
	"net/http"
	"strconv"

	crypt "github.com/donghquinn/blog_back_go/libraries/crypto"
	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/post"
	"github.com/donghquinn/blog_back_go/utils"
)

// 전체 포스트 가져오기 - 페이징
func GetPostController(res http.ResponseWriter, req *http.Request) {
	var getPostListRequest types.GetPostListRequest

	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	size, _ := strconv.Atoi(req.URL.Query().Get("size"))

	parseErr := utils.DecodeBody(req, &getPostListRequest)

	if parseErr != nil {
		log.Printf("[POST_EDIT] Parse View Specific Post Contents Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "Parse View Specific Post Contents Error",
			Result:  false,
		})

		return
	}

	unpinnedQueryResult, queryErr := QueryUnpinnedPostData(getPostListRequest.BlogId, page, size)

	if queryErr != nil {
		log.Printf("[POST_EDIT] Query Post Data Error: %v", queryErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Query Post Data Error",
			Result:  false,
		})

		return
	}

	pinnedQueryResult, pinnedErr := QueryisPinnedPostData(getPostListRequest.BlogId)

	if pinnedErr != nil {
		log.Printf("[POST_EDIT] Query Post Data Error: %v", pinnedErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Pin Error",
			Result:  false,
		})

		return
	}

	unpinnedTotalCount, unpinnedTotalCountErr := GetTotalUnPinnedPostCount(getPostListRequest.BlogId)

	if unpinnedTotalCountErr != nil {
		log.Printf("[POST_EDIT] Query Post Data Error: %v", unpinnedTotalCountErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  404,
			Code:    "04",
			Message: "UnPin Error",
			Result:  false,
		})

		return
	}

	var pinnedData []types.SelectAllPostDataResponse
	var unpinnedData []types.SelectAllPostDataResponse

	for _, data := range unpinnedQueryResult {
		decodedName, decodeErr := crypt.DecryptString(data.UserName)

		if decodeErr != nil {
			log.Printf("[LIST] Decoding User Name Error: %v", decodeErr)
			response.Response(res, response.CommonResponseWithMessage{
				Status:  405,
				Code:    "05",
				Message: "Decode Name Error",
				Result:  false,
			})

			return
		}

		unpinnedData = append(unpinnedData, types.SelectAllPostDataResponse{
			PostSeq:      data.PostSeq,
			PostTitle:    data.PostTitle,
			PostContents: data.PostContents,
			CategoryName: data.CategoryName,
			UserName:     decodedName,
			IsPinned:     data.IsPinned,
			Viewed:       data.Viewed,
			RegDate:      data.RegDate,
			ModDate:      data.ModDate,
		})
	}

	// 이름 디코딩 위해
	for _, data := range pinnedQueryResult {
		decodedName, decodeErr := crypt.DecryptString(data.UserName)

		if decodeErr != nil {
			log.Printf("[LIST] Decoding User Name Error: %v", decodeErr)
			response.Response(res, response.CommonResponseWithMessage{
				Status:  406,
				Code:    "06",
				Message: "Decode Name Error",
				Result:  false,
			})

			return
		}

		pinnedData = append(pinnedData, types.SelectAllPostDataResponse{
			PostSeq:      data.PostSeq,
			PostTitle:    data.PostTitle,
			PostContents: data.PostContents,
			CategoryName: data.CategoryName,
			UserName:     decodedName,
			IsPinned:     data.IsPinned,
			Viewed:       data.Viewed,
			RegDate:      data.RegDate,
			ModDate:      data.ModDate,
		})
	}

	response.Response(res, types.ResponsePostListType{
		Status:           http.StatusOK,
		Code:             "01",
		Message:          "SUCCESS",
		Result:           true,
		PinnedPostList:   pinnedData,
		UnpinnedPostList: unpinnedData,
		PostCount:        unpinnedTotalCount.Count,
		Page:             page,
		Size:             size,
	})

}

// 전체 포스트 가져오기 - 페이징
func GetPinnedPostController(res http.ResponseWriter, req *http.Request) {
	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	size, _ := strconv.Atoi(req.URL.Query().Get("size"))
	var getPinnedPostRequest types.GetPostListRequest

	parseErr := utils.DecodeBody(req, &getPinnedPostRequest)

	if parseErr != nil {
		log.Printf("[POST_PINNED] Parse View Specific Post Contents Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "Parse View Specific Post Contents Error",
			Result:  false,
		})

		return
	}

	pinnedQueryResult, pinnedErr := QueryisPinnedPostList(getPinnedPostRequest.BlogId, page, size)

	if pinnedErr != nil {
		log.Printf("[POST_PINNED] Query Post Data Error: %v", pinnedErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Query Post Data Error",
			Result:  false,
		})

		return
	}

	pinnedTotalCount, pinnedTotalCountErr := GetTotalPinnedPostCount(getPinnedPostRequest.BlogId)

	if pinnedTotalCountErr != nil {
		log.Printf("[POST_PINNED] Query Post Data Error: %v", pinnedTotalCountErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Query Post Data Error",
			Result:  false,
		})

		return
	}

	var pinnedData []types.SelectAllPostDataResponse

	// 이름 디코딩 위해
	for _, data := range pinnedQueryResult {
		decodedName, decodeErr := crypt.DecryptString(data.UserName)

		if decodeErr != nil {
			log.Printf("[LIST] Decoding User Name Error: %v", decodeErr)
			response.Response(res, response.CommonResponseWithMessage{
				Status:  403,
				Code:    "03",
				Message: "Decode Name Error",
				Result:  false,
			})

			return
		}

		pinnedData = append(pinnedData, types.SelectAllPostDataResponse{
			PostSeq:      data.PostSeq,
			PostTitle:    data.PostTitle,
			PostContents: data.PostContents,
			CategoryName: data.CategoryName,
			UserName:     decodedName,
			IsPinned:     data.IsPinned,
			Viewed:       data.Viewed,
			RegDate:      data.RegDate,
			ModDate:      data.ModDate,
		})
	}

	response.Response(res, types.ResponsePinnedPostListType{
		Status:         http.StatusOK,
		Code:           "01",
		Message:        "SUCCESS",
		Result:         true,
		PinnedPostList: pinnedData,
		PostCount:      pinnedTotalCount.Count,
		Page:           page,
		Size:           size,
	})
}

// 태그로 포스트 찾기
func GetPostsByTagController(res http.ResponseWriter, req *http.Request) {
	var getPostByTagRequest types.GetPostsByTagRequest

	parseErr := utils.DecodeBody(req, &getPostByTagRequest)

	if parseErr != nil {
		log.Printf("[LIST] Parse Request Body Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "Parse Request Body Error",
			Result:  false,
		})

		return
	}

	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	size, _ := strconv.Atoi(req.URL.Query().Get("size"))

	postList, totalPostCount, postErr := GetPostByTag(getPostByTagRequest, page, size)

	if postErr != nil {
		log.Printf("[LIST] Get Post List By Tag Error: %v", postErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Get Post List By Tag Error",
			Result:  false,
		})

		return
	}

	response.Response(res, types.ResponsePostByTagListType{
		Status:    http.StatusOK,
		Code:      "01",
		Message:   "SUCCES",
		Result:    true,
		PostList:  postList,
		PostCount: totalPostCount.Count,
	})
}

// 태그로 포스트 찾기
func GetPostsByCategoryController(res http.ResponseWriter, req *http.Request) {
	var getPostByCategoryRequest types.GetPostsByCategoryRequest

	parseErr := utils.DecodeBody(req, &getPostByCategoryRequest)

	if parseErr != nil {
		log.Printf("[LIST] Parse Request Body Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "Parse Request Body Error",
			Result:  false,
		})

		return
	}

	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	size, _ := strconv.Atoi(req.URL.Query().Get("size"))

	postList, totalCount, postErr := GetPostByCategory(getPostByCategoryRequest, page, size)

	if postErr != nil {
		log.Printf("[LIST] Get Post List By Tag Error: %v", parseErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Get Post List By Tag Error",
			Result:  false,
		})

		return
	}

	response.Response(res, types.ResponsePostByCategoryListType{
		Status:    http.StatusOK,
		Code:      "01",
		Message:   "SUCCESS",
		Result:    true,
		PostList:  postList,
		PostCount: totalCount.Count,
	})

}
