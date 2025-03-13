package post

import (
	"database/sql"
	"log"

	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/posts"
	types "github.com/donghquinn/blog_back_go/types/post"
)

// 특정 게시글 가져오기
func GetPostData(postSeq string, blogId string) (types.SelectSpecificPostDataResult, error){
	updateErr := UpdateViewCount(postSeq)

	if updateErr != nil {
		return types.SelectSpecificPostDataResult{}, updateErr
	}

	postList, getPostErr := GetPostContents(postSeq, blogId)

	if getPostErr != nil {
		return types.SelectSpecificPostDataResult{}, getPostErr
	}

	return postList, nil
}

// 조회수 올리기
func UpdateViewCount(postSeq string) error {
	connect, connectErr := database.InitDatabaseConnection()

	if connectErr != nil {
		log.Printf("[CONTENTS] Init Database Connection Error for Post Data: %v", connectErr)
		return connectErr
	}

	// 조회수 업데이트
	_, updateErr := connect.InsertQuery(queries.UpdateViewCount, postSeq)

	defer connect.Close()

	if updateErr != nil {
		log.Printf("[CONTENTS] Update View Count Error: %v", updateErr)
		return connectErr
	}

	return nil
}

// 개사굴 콘탠추 상세 조회
func GetPostContents(postSeq string, blogId string) (types.SelectSpecificPostDataResult, error) {
	var queryResult types.SelectSpecificPostDataResult


	connect, connectErr := database.InitDatabaseConnection()

	if connectErr != nil {
		log.Printf("[CONTENTS] Init Database Connection Error for Post Data: %v", connectErr)
		return types.SelectSpecificPostDataResult{}, connectErr
	}

	// 특정 게시글 조회
	result, queryErr := connect.QueryOne(queries.SelectSpecificPostContents, postSeq, blogId)

	defer connect.Close()

	if queryErr != nil {
		log.Printf("[CONTENTS] Query A Post Contents Error: %v", queryErr)
		return types.SelectSpecificPostDataResult{}, queryErr
	}

	postScanErr := result.Scan(
		&queryResult.PostSeq,
		&queryResult.PostTitle,
		&queryResult.PostContents,
		&queryResult.PostStatus,
		&queryResult.Tags,
		&queryResult.CategoryName,
		&queryResult.UserName,
		&queryResult.Viewed,
		&queryResult.IsPinned,
		&queryResult.RegDate,
		&queryResult.ModDate)


	if postScanErr != nil {
		if postScanErr == sql.ErrNoRows {
			queryResult.CategoryName = nil
			queryResult.Tags = nil

			return queryResult, nil
		}

		log.Printf("[CONTENTS] Can Post Data Error: %v", postScanErr)
		return types.SelectSpecificPostDataResult{}, postScanErr
	}
		// log.Println(*queryResult.CategoryName)
		// log.Println(*queryResult.Tags)
	return queryResult, nil
}

// func GetTagList(postSeq string) (types.SelectSpeicificPostTagDataResult, error) {
// 			// 태그 쿼리
// 	var tagQueryResult types.SelectSpeicificPostTagDataResult

// 	connect, connectErr := database.InitDatabaseConnection()

// 	if connectErr != nil {
// 		log.Printf("[CONTENTS] Init Database Connection Error for Post Data: %v", connectErr)
// 		return types.SelectSpeicificPostTagDataResult{},connectErr
// 	}

// 		// 태그들 조회
// 	tagResult, tagErr := database.QueryOne(connect, queries.SelectPostTags, postSeq)

// 	defer connect.Close()

// 	if tagErr != nil {
// 		log.Printf("[CONTENTS] Query Tags Error: %v", tagErr)
// 		return types.SelectSpeicificPostTagDataResult{}, tagErr
// 	}

// 	if tagResult!= nil {
// 		scanErr := tagResult.Scan(
// 			&tagQueryResult.TagName)

// 		if scanErr != nil {
// 			log.Printf("[CONTENTS] Scan Tag Query Data Error: %v", scanErr)
// 			return types.SelectSpeicificPostTagDataResult{}, scanErr
// 		}
// 	}

// 	return tagQueryResult, nil
// }