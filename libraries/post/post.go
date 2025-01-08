package post

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/posts"
	types "github.com/donghquinn/blog_back_go/types/post"
)

// TODO List 조회에 비밀글 값 추가

// 포스트들 가져오기 - 모듈함수
func QueryUnpinnedPostData(blogId string, page int, size int) ([]types.SelectAllPostDataResponse, error) {
	// parseBodyErr :=utils.DecodeBody(&req.Body)
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return nil, dbErr
	}

	// 페이징 파라미터 파싱
	result, queryErr := connect.GetMultiple(queries.SelectUnPinnedPosts, blogId, fmt.Sprintf("%d", size), fmt.Sprintf("%d", (page-1)*size))

	if queryErr != nil {
		log.Printf("[LIST] Get Unpinned Post Data Error: %v", queryErr)

		return nil, queryErr
	}

	defer connect.Close()

	var queryResult = []types.SelectAllPostDataResponse{}

	for result.Next() {
		var row types.SelectAllPostDataResponse

		scanErr := result.Scan(
			&row.PostSeq,
			&row.PostTitle,
			&row.PostContents,
			&row.CategoryName,
			&row.UserName,
			&row.IsPinned,
			&row.Viewed,
			&row.RegDate,
			&row.ModDate)

		if scanErr != nil {
			if scanErr == sql.ErrNoRows {
				return []types.SelectAllPostDataResponse{}, nil
			} else {
				log.Printf("[LIST] Scan and Assign Unpinned Query Result Error: %v", scanErr)

				return nil, scanErr
			}

		}

		queryResult = append(queryResult, row)
	}

	return queryResult, nil
}

// 포스트들 가져오기 - 모듈함수
func QueryisPinnedPostList(blogId string, page int, size int) ([]types.SelectAllPostDataResponse, error) {
	// parseBodyErr :=utils.DecodeBody(&req.Body)
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return nil, dbErr
	}

	// 페이징 파라미터 파싱
	result, queryErr := connect.GetMultiple(queries.SelectAllPinnedPosts, blogId, fmt.Sprintf("%d", size), fmt.Sprintf("%d", (page-1)*size))

	if queryErr != nil {
		log.Printf("[LIST] Get Pinned Post Data Error: %v", queryErr)

		return nil, queryErr
	}

	defer connect.Close()

	var queryResult = []types.SelectAllPostDataResponse{}

	for result.Next() {
		var row types.SelectAllPostDataResponse

		scanErr := result.Scan(
			&row.PostSeq,
			&row.PostTitle,
			&row.PostContents,
			&row.CategoryName,
			&row.UserName,
			&row.IsPinned,
			&row.Viewed,
			&row.RegDate,
			&row.ModDate)

		if scanErr != nil {
			if scanErr == sql.ErrNoRows {
				return []types.SelectAllPostDataResponse{}, nil
			} else {
				log.Printf("[LIST] Scan and Assign Pinned Query Result Error: %v", scanErr)

				return nil, scanErr
			}

		}

		queryResult = append(queryResult, row)
	}

	return queryResult, nil
}

// 포스트들 가져오기 - 모듈함수
func QueryisPinnedPostData(blogId string) ([]types.SelectAllPostDataResponse, error) {
	// parseBodyErr :=utils.DecodeBody(&req.Body)
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return nil, dbErr
	}

	// 페이징 파라미터 파싱
	result, queryErr := connect.GetMultiple(queries.SelectPinnedPosts, blogId)

	if queryErr != nil {
		log.Printf("[LIST] Get Pinned Post Data Error: %v", queryErr)

		return nil, queryErr
	}

	defer connect.Close()

	var queryResult = []types.SelectAllPostDataResponse{}

	for result.Next() {
		var row types.SelectAllPostDataResponse

		scanErr := result.Scan(
			&row.PostSeq,
			&row.PostTitle,
			&row.PostContents,
			&row.CategoryName,
			&row.UserName,
			&row.IsPinned,
			&row.Viewed,
			&row.RegDate,
			&row.ModDate)

		if scanErr != nil {
			if scanErr == sql.ErrNoRows {
				return []types.SelectAllPostDataResponse{}, nil
			} else {
				log.Printf("[LIST] Scan and Assign Pinned Query Result Error: %v", scanErr)

				return nil, scanErr
			}

		}

		queryResult = append(queryResult, row)
	}

	return queryResult, nil
}

// 고정 개시글 전체 개수
func GetTotalUnPinnedPostCount(blogId string) (types.PostTotalUnPinnedCountType, error) {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return types.PostTotalUnPinnedCountType{}, dbErr
	}

	queryResult, queryErr := connect.QueryOne(queries.SelectUnPinnedPostCount, blogId)

	if queryErr != nil {
		log.Printf("[LIST] Get UnPinned Post Count Error: %v", queryErr)

		return types.PostTotalUnPinnedCountType{}, queryErr
	}

	var unPinnedTotalCount types.PostTotalUnPinnedCountType

	queryResult.Scan(&unPinnedTotalCount.Count)

	return unPinnedTotalCount, nil
}

// 고정 전체 게시글 개수 조회
func GetTotalPinnedPostCount(blogId string) (types.PostTotalUnPinnedCountType, error) {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return types.PostTotalUnPinnedCountType{}, dbErr
	}

	queryResult, queryErr := connect.QueryOne(queries.SelectPinnedPostCount, blogId)

	if queryErr != nil {
		log.Printf("[LIST] Get Pinned Post Count Error: %v", queryErr)

		return types.PostTotalUnPinnedCountType{}, queryErr
	}

	var unPinnedTotalCount types.PostTotalUnPinnedCountType

	queryResult.Scan(&unPinnedTotalCount.Count)

	return unPinnedTotalCount, nil
}

// 게시글 태그로 조회
func GetPostByTag(data types.GetPostsByTagRequest, page int, size int) ([]types.PostsByTagsResponseType, types.PostTotalUnPinnedCountType, error) {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return []types.PostsByTagsResponseType{}, types.PostTotalUnPinnedCountType{}, dbErr
	}

	posts, selectErr := connect.GetMultiple(queries.SelectPostByTags, "%"+data.TagName+"%", data.BlogId, fmt.Sprintf("%d", size), fmt.Sprintf("%d", (page-1)*size))

	if selectErr != nil {
		log.Printf("[POST_TAG] GET Post by TagName Error: %v", selectErr)
		return []types.PostsByTagsResponseType{}, types.PostTotalUnPinnedCountType{}, selectErr
	}

	defer connect.Close()

	var postsData []types.SelectPostsByTags

	// Array https://stackoverflow.com/questions/14477941/read-select-columns-into-string-in-go
	for posts.Next() {
		var row types.SelectPostsByTags

		scanErr := posts.Scan(
			&row.TagName,
			&row.CategoryName,
			&row.UserName,
			&row.PostSeq,
			&row.PostTitle,
			&row.PostContents,
			&row.Viewed,
			&row.RegDate,
			&row.ModDate)

		if scanErr != nil {
			log.Printf("[POST_TAG] Scan Query Result Error: %v", scanErr)
			return []types.PostsByTagsResponseType{}, types.PostTotalUnPinnedCountType{}, scanErr
		}

		postsData = append(postsData, row)
	}

	connect2, dbErr2 := database.InitDatabaseConnection()

	if dbErr2 != nil {
		return []types.PostsByTagsResponseType{}, types.PostTotalUnPinnedCountType{}, dbErr2
	}

	count, countErr := connect2.QueryOne(queries.SelectTotalPostCountByTags, "%"+data.TagName+"%", data.BlogId)

	if countErr != nil {
		log.Printf("[POST_TAG] GET Post Total Count by TagName Error: %v", countErr)
		return []types.PostsByTagsResponseType{}, types.PostTotalUnPinnedCountType{}, countErr
	}

	var totalPostCount types.PostTotalUnPinnedCountType

	count.Scan(&totalPostCount.Count)

	// stringify된 array를 array로
	var postByTagsList []types.PostsByTagsResponseType

	for _, d := range postsData {
		var tempTag []string

		jsonErr := json.Unmarshal([]byte(d.TagName), &tempTag)

		if jsonErr != nil {
			log.Printf("[POST_TAG] Unmarshing Array Error: %v", jsonErr)
			return []types.PostsByTagsResponseType{}, types.PostTotalUnPinnedCountType{}, jsonErr
		}

		data := types.PostsByTagsResponseType{
			TagName:      tempTag,
			CategoryName: d.CategoryName,
			PostTitle:    d.PostTitle,
			PostContents: d.PostContents,
			PostSeq:      d.PostSeq,
			Viewed:       d.Viewed,
			RegDate:      d.RegDate,
			ModDate:      d.ModDate}

		postByTagsList = append(postByTagsList, data)
	}

	return postByTagsList, totalPostCount, nil
}

// 게시글 카테고리로 조회
func GetPostByCategory(data types.GetPostsByCategoryRequest, page int, size int) ([]types.PostByCategoryResponseType, types.PostTotalUnPinnedCountType, error) {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return []types.PostByCategoryResponseType{}, types.PostTotalUnPinnedCountType{}, dbErr
	}

	log.Println(data.BlogId)
	posts, selectErr := connect.GetMultiple(queries.SelectPostByCategory, "%"+data.CategoryName+"%", data.BlogId, fmt.Sprintf("%d", size), fmt.Sprintf("%d", (page-1)*size))

	if selectErr != nil {
		log.Printf("[POST_CATEGORY] GET Post by CategoryName Error: %v", selectErr)
		return []types.PostByCategoryResponseType{}, types.PostTotalUnPinnedCountType{}, selectErr
	}

	defer connect.Close()

	var postsData []types.SelectPostsByTags

	// Array https://stackoverflow.com/questions/14477941/read-select-columns-into-string-in-go
	for posts.Next() {
		var row types.SelectPostsByTags

		scanErr := posts.Scan(
			&row.TagName,
			&row.CategoryName,
			&row.UserName,
			&row.PostSeq,
			&row.PostTitle,
			&row.PostContents,
			&row.Viewed,
			&row.RegDate,
			&row.ModDate)

		if scanErr != nil {
			if scanErr == sql.ErrNoRows {
				postsData = make([]types.SelectPostsByTags, 0)
			} else {
				log.Printf("[POST_CATEGORY] Scan Query Result Error: %v", scanErr)
				return []types.PostByCategoryResponseType{}, types.PostTotalUnPinnedCountType{}, scanErr
			}

		}

		postsData = append(postsData, row)
	}

	connect2, dbErr2 := database.InitDatabaseConnection()

	if dbErr2 != nil {
		return []types.PostByCategoryResponseType{}, types.PostTotalUnPinnedCountType{}, dbErr2
	}

	count, countErr := connect2.QueryOne(queries.SelectTotalPostCountByCategory, "%"+data.CategoryName+"%", data.BlogId)

	if countErr != nil {
		log.Printf("[POST_TAG] GET Post Total Count by Category Error: %v", countErr)
		return []types.PostByCategoryResponseType{}, types.PostTotalUnPinnedCountType{}, countErr
	}

	var totalPostCount types.PostTotalUnPinnedCountType

	count.Scan(&totalPostCount.Count)

	// stringify된 array를 array로
	var postByCategoryList []types.PostByCategoryResponseType

	for _, d := range postsData {
		var tempTag []string

		if d.TagName != "NULL" {
			jsonErr := json.Unmarshal([]byte(d.TagName), &tempTag)
			if jsonErr != nil {
				log.Printf("[POST_CATEGORY] Unmarshing Array Error: %v", jsonErr)
				return []types.PostByCategoryResponseType{}, types.PostTotalUnPinnedCountType{}, jsonErr
			}
		}

		data := types.PostByCategoryResponseType{
			TagName:      tempTag,
			CategoryName: d.CategoryName,
			PostTitle:    d.PostTitle,
			PostContents: d.PostContents,
			PostSeq:      d.PostSeq,
			Viewed:       d.Viewed,
			RegDate:      d.RegDate,
			ModDate:      d.ModDate}

		postByCategoryList = append(postByCategoryList, data)
	}

	return postByCategoryList, totalPostCount, nil
}
