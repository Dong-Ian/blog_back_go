package post

import (
	"database/sql"
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
		log.Printf("[POST_LIST] Create Query Builder Error: %v", queryBuildErr)
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
		log.Printf("[POST_LIST] Get Unpinned Post Data Error: %v", queryErr)

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
				log.Printf("[POST_LIST] Scan and Assign Unpinned Query Result Error: %v", scanErr)
				return nil, scanErr
			}
		}

		decodedName, decodeErr := crypt.DecryptString(encodedName)

		if decodeErr != nil {
			log.Printf("[POST_LIST] Decode user name Error: %v", decodeErr)
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
		log.Printf("[POST_LIST] Create Query Builder Error: %v", queryBuildErr)
		return -999, queryBuildErr
	}

	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return -999, dbErr
	}

	queryResult, queryErr := connect.QueryBuilderOneRow(query, args)

	if queryErr != nil {
		log.Printf("[POST_LIST] Get UnPinned Post Count Error: %v", queryErr)

		return -999, queryErr
	}

	var totalCount int64

	if scanErr := queryResult.Scan(&totalCount); scanErr != nil {
		log.Printf("[POST_LIST] Total Count Scan Erro: %v", scanErr)
		return -9999, scanErr
	}

	return totalCount, nil
}
