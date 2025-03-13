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

func InsertPostData(registerPostRequest types.RegisterPostRequest, userId string, blogId string) (int64, error) {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return -9999, dbErr
	}

	// 데이터 입력
	insertId, queryErr := connect.InsertQuery(
		queries.InsertPost, 
		userId, 
		registerPostRequest.PostTitle, 
		registerPostRequest.PostContents,
		registerPostRequest.IsPinned,
		// strconv.Itoa(registerPostRequest.IsSecret),
		blogId)

	if queryErr != nil {
		log.Printf("[REGISTER] Insert Post Data Error: %v", queryErr)
		return -9999, queryErr
	}

	defer connect.Close()

	postSeq := strconv.Itoa(int(insertId))

	categories := registerPostRequest.Category

	if categories != ""  {
		insertCategoriesErr := InsertCategories(categories, postSeq, blogId)

		if insertCategoriesErr != nil {
			log.Printf("[REGISTER] Insert Categories Error: %v", insertCategoriesErr)

			return -999999, insertCategoriesErr
		}
	}

	tags := registerPostRequest.Tags

	if len(tags) > 0 {
		insertTagErr := InsertTags(tags, postSeq, blogId)
		if insertTagErr != nil {
			log.Printf("[REGISTER] Insert TagList Error: %v", insertTagErr)
			return -99999, insertTagErr
		}
	}

	for _, seq := range(registerPostRequest.ImageSeqs) {
		// 파일 데이터 업데이트
		_, insertUpdateErr := connect.InsertQuery(queries.InsertUpdatePostImage, postSeq, seq)

		if insertUpdateErr != nil {
			log.Printf("[REGISTER] Insert Update File Data Error: %v", insertUpdateErr)
			return -999999, insertUpdateErr
		}
	}

	return insertId, nil
}

func InsertCategories(categories string, postSeq string, blogId string) error {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	isValidCategory := utils.ValidateRequestValue(categories)

	if isValidCategory {
		_, categoryErr := connect.InsertQuery(queries.InsertCategory, postSeq, categories, blogId)

		if categoryErr != nil {
			log.Printf("[REGISTER] Insert category data Error: %v", categoryErr)
			return categoryErr
		}
	}

	defer connect.Close()

	return nil
}

func InsertTags(tags []string, postSeq string, blogId string) error {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	// Array https://www.infracody.com/2023/08/how-to-deal-with-array-data-in-mysql.html
	tagArray, _ := json.Marshal(tags)

	_, tagQueryErr := connect.InsertQuery(queries.InsertTag, postSeq, string(tagArray), blogId)

	defer connect.Close()

	if tagQueryErr != nil {
		log.Printf("[REGISTER] Insert Tag Data Error: %v", tagQueryErr)

		return tagQueryErr
	}

	return nil
}
