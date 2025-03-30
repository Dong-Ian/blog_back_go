package post

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	crypt "github.com/donghquinn/blog_back_go/libraries/crypto"
	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/posts"
	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/post"
	"github.com/donghquinn/blog_back_go/utils"
	"github.com/gorilla/mux"
)

// 전체 포스트 가져오기 - 페이징
func GetPostListController(res http.ResponseWriter, req *http.Request) {
	blogId, getErr := utils.GetBlogIdFromContext(req.Context())

	if !getErr {
		response.Response(res, response.CommonResponseWithMessage{
			Status:  http.StatusBadRequest,
			Code:    "01",
			Message: "Get BlogId Error from context",
			Result:  false,
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

	unpinnedQueryResult, queryErr := QueryPostList(blogId, isPinned, category, tag, limit, offset)

	if queryErr != nil {
		log.Printf("[POST_LIST] Query Post Data Error: %v", queryErr)
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
		log.Printf("[POST_LIST] Query Post Data Error: %v", queryErr)
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

// 게시글 컨텐츠 컨트롤러
func PostContentsController(res http.ResponseWriter, req *http.Request) {
	pathVar := mux.Vars(req)
	postSeq := pathVar["postSeq"]

	blogId, getErr := utils.GetBlogIdFromContext(req.Context())

	if !getErr {
		response.Response(res, response.CommonResponseWithMessage{
			Status:  http.StatusBadRequest,
			Code:    "01",
			Message: "No Blog Id Found",
			Result:  false,
		})

		return
	}

	// 게시글 쿼리
	queryResult, queryErr := GetPostData(postSeq, blogId)

	if queryErr != nil {
		log.Printf("[POST_CONTENT] Query Specific Contents Error: %v", queryErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Query Specific Contents Error",
			Result:  false,
		})

		return
	}

	imageData, imageErr := GetImageData(postSeq)

	if imageErr != nil {
		log.Printf("[POST_CONTENT] Image Data Error: %v", queryErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Image Data Error",
			Result:  false,
		})

		return
	}

	var urlArray []string

	if len(imageData) > 0 {
		// 게시글 URL 배열 만들기
		for _, data := range imageData {
			url, getErr := database.GetImageUrl(data.ObjectName, data.FileFormat)

			if getErr != nil {
				log.Printf("[POST_CONTENT] Get Presigned URL Error: %v", getErr)

				response.Response(res, response.CommonResponseWithMessage{
					Status:  404,
					Code:    "04",
					Message: "Get Presigned URL Error",
					Result:  false,
				})

				return
			}

			if url == nil {
				urlArray = make([]string, 0)
			}

			urlArray = append(urlArray, url.String())
		}
	}

	userName, _ := crypt.DecryptString(queryResult.UserName)

	// 특정 게시글 태그 배열 가공해서 담아 응답
	var tagsArray []string

	if queryResult.Tags != nil {
		jsonErr := json.Unmarshal([]byte(*queryResult.Tags), &tagsArray)
		if jsonErr != nil {
			log.Printf("[CONTENTS] JSON Unmarsh tag array Error: %v", jsonErr)

			response.Response(res, response.CommonResponseWithMessage{
				Status:  405,
				Code:    "05",
				Message: "Unmarshing Tags Error",
				Result:  false,
			})

			return
		}
	} else {
		tagsArray = make([]string, 0)
	}

	var categoryName string

	if queryResult.CategoryName != nil {
		categoryName = *queryResult.CategoryName
	} else {
		categoryName = ""
	}

	// 게시글 컨텐츠 데이터
	postContentsData := types.ViewSpecificPostContentsResponse{
		PostSeq:      queryResult.PostSeq,
		PostTitle:    queryResult.PostTitle,
		Tags:         tagsArray,
		PostContents: queryResult.PostContents,
		CategoryName: categoryName,
		UserName:     userName,
		Urls:         urlArray,
		Viewed:       queryResult.Viewed,
		IsPinned:     queryResult.IsPinned,
		RegDate:      queryResult.RegDate,
		ModDate:      queryResult.ModDate,
	}

	response.Response(res, types.ResponsePostContentsType{
		Status:   http.StatusOK,
		Code:     "01",
		Message:  "SUCCESS",
		PostList: postContentsData,
		Result:   true,
	})

}

// 게시글 번호에 맞는 file 데이터 전부 가져오기
func GetImageData(postSeq string) ([]types.SelectPostImageData, error) {
	var returnImageDate []types.SelectPostImageData

	connect, connectErr := database.InitDatabaseConnection()

	if connectErr != nil {
		log.Printf("[CONTENTS] Init Database Connection Error for Image Data: %v", connectErr)
		return []types.SelectPostImageData{}, connectErr
	}

	result, queryErr := connect.GetMultiple(queries.SelectImageData, postSeq, "POST_IMAGE")

	if queryErr != nil {
		log.Printf("[CONTENTS] Query Image Data Error: %v", queryErr)
		return []types.SelectPostImageData{}, queryErr
	}

	defer connect.Close()

	for result.Next() {
		var row types.SelectPostImageData

		scanErr := result.Scan(
			&row.ObjectName,
			&row.FileFormat,
			&row.TargetPurpose,
			&row.TargetSeq)

		if scanErr != nil {
			if scanErr == sql.ErrNoRows {
				returnImageDate = append(returnImageDate, types.SelectPostImageData{})
			} else {
				log.Printf("[CONTENTS] Scan Files Error: %v", scanErr)
				return []types.SelectPostImageData{}, nil
			}
		}

		returnImageDate = append(returnImageDate, row)
	}

	return returnImageDate, nil
}

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
