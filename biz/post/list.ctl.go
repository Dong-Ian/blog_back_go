package post

import (
	"log"
	"net/http"
	"strconv"

	"github.com/donghquinn/blog_back_go/dto"
	crypt "github.com/donghquinn/blog_back_go/libraries/crypto"
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
		dto.SetErrorResponse(res, 401, "01", "Parse View Specific Post Contents Error", parseErr)
		return
	}

	unpinnedQueryResult, queryErr := QueryUnpinnedPostData(getPostListRequest.BlogId, page, size)

	if queryErr != nil {
		dto.SetErrorResponse(res, 402, "02", "Query Post Data Error", queryErr)

		return
	}

	pinnedQueryResult, pinnedErr := QueryisPinnedPostData(getPostListRequest.BlogId)

	if pinnedErr != nil {
		dto.SetErrorResponse(res, 402, "02", "Query Post Data Error", pinnedErr)

		return
	}

	unpinnedTotalCount, unpinnedTotalCountErr := GetTotalUnPinnedPostCount(getPostListRequest.BlogId)

	if unpinnedTotalCountErr != nil {
		dto.SetErrorResponse(res, 402, "02", "Query Post Data Error", unpinnedTotalCountErr)

		return
	}

	var pinnedData []types.SelectAllPostDataResponse
	var unpinnedData []types.SelectAllPostDataResponse

	for _, data := range unpinnedQueryResult {
		decodedName, decodeErr := crypt.DecryptString(data.UserName)

		if decodeErr != nil {
			log.Printf("[LIST] Decoding User Name Error: %v", decodeErr)
			dto.SetErrorResponse(res, 403, "03", "Decode Name Error", decodeErr)
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
			dto.SetErrorResponse(res, 403, "03", "Decode Name Error", decodeErr)
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

	dto.SetPostListResponse(res, 200, "01", unpinnedData, pinnedData, unpinnedTotalCount.Count, page, size)
}

// 전체 포스트 가져오기 - 페이징
func GetPinnedPostController(res http.ResponseWriter, req *http.Request) {
	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	size, _ := strconv.Atoi(req.URL.Query().Get("size"))
	var getPinnedPostRequest types.GetPostListRequest

	parseErr := utils.DecodeBody(req, &getPinnedPostRequest)

	if parseErr != nil {
		dto.SetErrorResponse(res, 401, "01", "Parse View Specific Post Contents Error", parseErr)
		return
	}

	pinnedQueryResult, pinnedErr := QueryisPinnedPostList(getPinnedPostRequest.BlogId, page, size)

	if pinnedErr != nil {
		dto.SetErrorResponse(res, 401, "01", "Query Post Data Error", pinnedErr)

		return
	}

	pinnedTotalCount, pinnedTotalCountErr := GetTotalPinnedPostCount(getPinnedPostRequest.BlogId)

	if pinnedTotalCountErr != nil {
		dto.SetErrorResponse(res, 401, "01", "Query Post Data Error", pinnedTotalCountErr)

		return
	}

	var pinnedData []types.SelectAllPostDataResponse

	// 이름 디코딩 위해
	for _, data := range pinnedQueryResult {
		decodedName, decodeErr := crypt.DecryptString(data.UserName)

		if decodeErr != nil {
			log.Printf("[LIST] Decoding User Name Error: %v", decodeErr)
			dto.SetErrorResponse(res, 402, "02", "Decode Name Error", decodeErr)
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

	dto.SetPinnedPostListResponse(res, 200, "01", pinnedData, pinnedTotalCount.Count, page, size)
}

// 태그로 포스트 찾기
func GetPostsByTagController(res http.ResponseWriter, req *http.Request) {
	var getPostByTagRequest types.GetPostsByTagRequest

	parseErr := utils.DecodeBody(req, &getPostByTagRequest)

	if parseErr != nil {
		dto.SetErrorResponse(res, 401, "01", "Parse Request Body Error", parseErr)
		return
	}

	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	size, _ := strconv.Atoi(req.URL.Query().Get("size"))

	postList, totalPostCount, postErr := GetPostByTag(getPostByTagRequest, page, size)

	if postErr != nil {
		dto.SetErrorResponse(res, 402, "02", "Get Post List By Tag Error", postErr)
		return
	}

	dto.SetPostByTagResponse(res, 200, "01", postList, totalPostCount.Count)
}

// 태그로 포스트 찾기
func GetPostsByCategoryController(res http.ResponseWriter, req *http.Request) {
	var getPostByCategoryRequest types.GetPostsByCategoryRequest

	parseErr := utils.DecodeBody(req, &getPostByCategoryRequest)

	if parseErr != nil {
		dto.SetErrorResponse(res, 401, "01", "Parse Request Body Error", parseErr)
		return
	}

	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	size, _ := strconv.Atoi(req.URL.Query().Get("size"))

	postList, totalCount, postErr := GetPostByCategory(getPostByCategoryRequest, page, size)

	if postErr != nil {
		dto.SetErrorResponse(res, 402, "02", "Get Post List By Tag Error", postErr)
		return
	}

	dto.SetPostByCategoryResponse(res, 200, "01", postList, totalCount.Count)
}
