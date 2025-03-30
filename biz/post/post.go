package post

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	crypt "github.com/donghquinn/blog_back_go/libraries/crypto"
	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/posts"
	types "github.com/donghquinn/blog_back_go/types/post"
	"github.com/donghquinn/gqbd"
)

// TODO List 조회에 비밀글 값 추가

// 포스트들 가져오기 - 모듈함수
func QueryUnpinnedPostData(blogId string, isPinned string, category string, tag string, limit int, offset int) ([]types.SelectAllPostDataResponse, error) {
	qb := gqbd.BuildSelect(gqbd.MariaDB, "post_table p", "p.post_seq", "p.post_title", "p.post_contents",
		"c.category_name", "IFNULL(u.user_name, 'unknown') AS user_name", "p.is_pinned", "p.viewd",
		"p.reg_date", "p.mod_date").
		LeftJoin("user_table u", "u.user_id = p.user_id AND u.user_status = 1").
		LeftJoin("category_table c", "c.post_seq = p.post_seq AND c.category_status = 1").
		Where("p.blog_owner = ?", blogId).
		Where("p.post_status = ?", "1")

	if isPinned != "" {
		qb = qb.Where("p.is_pinned = ?", isPinned)
	}

	if category != "" {
		qb = qb.Where("p.category_name LIKE ?", "%"+category+"%")
	}

	if tag != "" {
		// tagList := strings.Split(tag, ",")
		qb = qb.Where("p.tags LIKE ?", "%"+tag+"%")
	}

	qb = qb.OrderBy("p.reg_date", "DESC", nil).
		Limit(limit).
		Offset(offset)

	query, args, queryBuildErr := qb.Build()

	if queryBuildErr != nil {
		log.Printf("[LIST] Create Query Builder Error: %v", queryBuildErr)
		return nil, queryBuildErr
	}

	// parseBodyErr :=utils.DecodeBody(&req.Body)
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return nil, dbErr
	}

	// 페이징 파라미터 파싱
	result, queryErr := connect.QueryBuilderRows(query, args)

	if queryErr != nil {
		log.Printf("[LIST] Get Unpinned Post Data Error: %v", queryErr)

		return nil, queryErr
	}

	var queryResult = []types.SelectAllPostDataResponse{}

	for result.Next() {
		var row types.SelectAllPostDataResponse
		var encodedName string // 인코딩된 사용자 이름을 저장할 임시 변수

		scanErr := result.Scan(
			&row.PostSeq,
			&row.PostTitle,
			&row.PostContents,
			&row.CategoryName,
			&encodedName,
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

		decodedName, decodeErr := crypt.DecryptString(encodedName)

		if decodeErr != nil {
			log.Printf("[LIST] Decode user name Error: %v", decodeErr)
			return nil, decodeErr
		}

		row.UserName = decodedName

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

	var queryResult = []types.SelectAllPostDataResponse{}

	for result.Next() {
		var row types.SelectAllPostDataResponse
		var encoded string

		scanErr := result.Scan(
			&row.PostSeq,
			&row.PostTitle,
			&row.PostContents,
			&row.CategoryName,
			&encoded,
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

		decrypted, decryptErr := crypt.DecryptString(encoded)

		if decryptErr != nil {
			log.Printf("[LIST] Decrypt User Name Err: %v", decryptErr)
			row.UserName = encoded
		} else {
			row.UserName = decrypted
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
func GetTotalPostCount(blogId string, isPinned string, category string, tag string) (int64, error) {
	qb := gqbd.BuildSelect(gqbd.MariaDB, "post_table p", "COUNT(p.post_seq)").
		Where("p.blog_owner = ?", blogId).
		Where("p.post_status = ?", "1")

	if isPinned != "" {
		qb = qb.Where("p.is_pinned = ?", isPinned)
	}

	if category != "" {
		qb = qb.Where("p.category_name LIKE ?", "%"+category+"%")
	}

	if tag != "" {
		// tagList := strings.Split(tag, ",")
		qb = qb.Where("p.tags LIKE ?", "%"+tag+"%")
	}

	query, args, queryBuildErr := qb.Build()

	if queryBuildErr != nil {
		log.Printf("[LIST] Create Query Builder Error: %v", queryBuildErr)
		return -999, queryBuildErr
	}

	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return -999, dbErr
	}

	queryResult, queryErr := connect.QueryBuilderOneRow(query, args)

	if queryErr != nil {
		log.Printf("[LIST] Get UnPinned Post Count Error: %v", queryErr)

		return -999, queryErr
	}

	var totalCount int64

	queryResult.Scan(&totalCount)

	return totalCount, nil
}

// 고정 전체 게시글 개수 조회
func GetTotalPinnedPostCount(blogId string) (types.PostTotalCountType, error) {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return types.PostTotalCountType{}, dbErr
	}

	queryResult, queryErr := connect.QueryOne(queries.SelectPinnedPostCount, blogId)

	if queryErr != nil {
		log.Printf("[LIST] Get Pinned Post Count Error: %v", queryErr)

		return types.PostTotalCountType{}, queryErr
	}

	var unPinnedTotalCount types.PostTotalCountType

	queryResult.Scan(&unPinnedTotalCount.Count)

	return unPinnedTotalCount, nil
}

// 게시글 태그로 조회
func GetPostByTag(data types.GetPostsByTagRequest, page int, size int) ([]types.PostsByTagsResponseType, types.PostTotalCountType, error) {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return []types.PostsByTagsResponseType{}, types.PostTotalCountType{}, dbErr
	}

	posts, selectErr := connect.GetMultiple(queries.SelectPostByTags, "%"+data.TagName+"%", data.BlogId, fmt.Sprintf("%d", size), fmt.Sprintf("%d", (page-1)*size))

	if selectErr != nil {
		log.Printf("[POST_TAG] GET Post by TagName Error: %v", selectErr)
		return []types.PostsByTagsResponseType{}, types.PostTotalCountType{}, selectErr
	}

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
			return []types.PostsByTagsResponseType{}, types.PostTotalCountType{}, scanErr
		}

		postsData = append(postsData, row)
	}

	connect2, dbErr2 := database.InitDatabaseConnection()

	if dbErr2 != nil {
		return []types.PostsByTagsResponseType{}, types.PostTotalCountType{}, dbErr2
	}

	count, countErr := connect2.QueryOne(queries.SelectTotalPostCountByTags, "%"+data.TagName+"%", data.BlogId)

	if countErr != nil {
		log.Printf("[POST_TAG] GET Post Total Count by TagName Error: %v", countErr)
		return []types.PostsByTagsResponseType{}, types.PostTotalCountType{}, countErr
	}

	var totalPostCount types.PostTotalCountType

	count.Scan(&totalPostCount.Count)

	// stringify된 array를 array로
	var postByTagsList []types.PostsByTagsResponseType

	for _, d := range postsData {
		var tempTag []string

		jsonErr := json.Unmarshal([]byte(d.TagName), &tempTag)

		if jsonErr != nil {
			log.Printf("[POST_TAG] Unmarshing Array Error: %v", jsonErr)
			return []types.PostsByTagsResponseType{}, types.PostTotalCountType{}, jsonErr
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
func GetPostByCategory(data types.GetPostsByCategoryRequest, page int, size int) ([]types.PostByCategoryResponseType, types.PostTotalCountType, error) {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return []types.PostByCategoryResponseType{}, types.PostTotalCountType{}, dbErr
	}

	log.Println(data.BlogId)
	posts, selectErr := connect.GetMultiple(queries.SelectPostByCategory, "%"+data.CategoryName+"%", data.BlogId, fmt.Sprintf("%d", size), fmt.Sprintf("%d", (page-1)*size))

	if selectErr != nil {
		log.Printf("[POST_CATEGORY] GET Post by CategoryName Error: %v", selectErr)
		return []types.PostByCategoryResponseType{}, types.PostTotalCountType{}, selectErr
	}

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
				return []types.PostByCategoryResponseType{}, types.PostTotalCountType{}, scanErr
			}

		}

		postsData = append(postsData, row)
	}

	connect2, dbErr2 := database.InitDatabaseConnection()

	if dbErr2 != nil {
		return []types.PostByCategoryResponseType{}, types.PostTotalCountType{}, dbErr2
	}

	count, countErr := connect2.QueryOne(queries.SelectTotalPostCountByCategory, "%"+data.CategoryName+"%", data.BlogId)

	if countErr != nil {
		log.Printf("[POST_TAG] GET Post Total Count by Category Error: %v", countErr)
		return []types.PostByCategoryResponseType{}, types.PostTotalCountType{}, countErr
	}

	var totalPostCount types.PostTotalCountType

	count.Scan(&totalPostCount.Count)

	// stringify된 array를 array로
	var postByCategoryList []types.PostByCategoryResponseType

	for _, d := range postsData {
		var tempTag []string

		if d.TagName != "NULL" {
			jsonErr := json.Unmarshal([]byte(d.TagName), &tempTag)
			if jsonErr != nil {
				log.Printf("[POST_CATEGORY] Unmarshing Array Error: %v", jsonErr)
				return []types.PostByCategoryResponseType{}, types.PostTotalCountType{}, jsonErr
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
