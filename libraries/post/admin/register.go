package post // 게시글 데이터 입력
import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/admin/posts"
	types "github.com/donghquinn/blog_back_go/types/admin/posts"
	"github.com/donghquinn/blog_back_go/utils"
)

func InsertPostData(registerPostRequest types.RegisterPostRequest, userId string) error {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	// 데이터 입력
	insertId, queryErr := database.InsertQuery(
		connect, 
		queries.InsertPost, 
		userId, 
		registerPostRequest.PostTitle, 
		registerPostRequest.PostContents,
		registerPostRequest.IsPinned)

	if queryErr != nil {
		log.Printf("[REGISTER] Insert Post Data Error: %v", queryErr)
		return queryErr
	}

	defer connect.Close()

	postSeq := strconv.Itoa(int(insertId))

	isValidCategory := utils.ValidateRequestValue(registerPostRequest.Category)

	if isValidCategory {
		_, categoryErr := database.InsertQuery(connect, queries.InsertCategory, postSeq, registerPostRequest.Category)

		if categoryErr != nil {
			log.Printf("[REGISTER] Insert category data Error: %v", categoryErr)
			return categoryErr
		}
	}

	tags := registerPostRequest.Tags

	if len(tags) > 0 {
		// Array https://www.infracody.com/2023/08/how-to-deal-with-array-data-in-mysql.html
		tagArray, _ := json.Marshal(tags)

		_, tagQueryErr := database.InsertQuery(connect, queries.InsertTag, postSeq, string(tagArray))

		if tagQueryErr != nil {
			log.Printf("[REGISTER] Insert Tag Data Error: %v", tagQueryErr)

			return tagQueryErr
		}
	}

	for _, seq := range(registerPostRequest.ImageSeqs) {
		// 파일 데이터 업데이트
		_, insertUpdateErr := database.InsertQuery(connect, queries.InsertUpdatePostImage, postSeq, seq)

		if insertUpdateErr != nil {
			log.Printf("[REGISTER] Insert Update File Data Error: %v", insertUpdateErr)
			return insertUpdateErr
		}
	}

	return nil
}