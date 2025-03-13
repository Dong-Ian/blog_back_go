package post

import (
	"log"

	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/admin/posts"
	types "github.com/donghquinn/blog_back_go/types/admin/posts"
)

// 게시글 삭제
func DeletePost(postSeq string, blogId string) error {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	_, deletePostErr := connect.InsertQuery(queries.DeletePost, "0", postSeq, blogId)

	if deletePostErr != nil {
		log.Printf("[DELETE] Delete Post Error: %v", deletePostErr)
		return deletePostErr
	}

	defer connect.Close()

	deleteCategoryErr := DeleteCategory(postSeq, blogId)

	if deleteCategoryErr != nil {
		return deleteCategoryErr
	}

	deleteTagErr := DeleteTag(postSeq, blogId)

	if deleteTagErr != nil {
		return deleteTagErr
	}
	
	return nil
}

func DeleteCategory(postSeq string, blogOwner string) error {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	_, deleteCategoryErr := connect.InsertQuery(queries.DeleteCategory, "0", postSeq, blogOwner)

	if deleteCategoryErr != nil {
		log.Printf("[DELETE] Delete Post Category Error: %v", deleteCategoryErr)
		return deleteCategoryErr
	}

	defer connect.Close()

	return nil
}

func DeleteTag(postSeq string, blogOwner string) error {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	_, deleteTagErr := connect.InsertQuery(queries.DeleteTags, "0", postSeq, blogOwner)

	if deleteTagErr != nil {
		log.Printf("[DELETE] Delete Post Tags Error: %v", deleteTagErr)
		return deleteTagErr
	}

	defer connect.Close()

	return nil
}

// 고정 게시글 업데이트
func UpdatePinPost(data types.UpdatePinRequest, blogId string) error {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	_, updateErr := connect.InsertQuery(queries.UpdatePinPost, "1", data.PostSeq, blogId)

	if updateErr != nil {
		log.Printf("[PIN] Update Pin Post Error: %v", updateErr)
		return updateErr
	}

	defer connect.Close()

	return nil
}

// 고정 게시글 해제
func UpdateUnPinPost(data types.UpdatePinRequest, blogId string) error {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	_, updateErr := connect.InsertQuery(queries.UpdatePinPost, "0", data.PostSeq, blogId)

	if updateErr != nil {
		log.Printf("[PIN] Update Un-Pin Post Error: %v", updateErr)
		return updateErr
	}

	defer connect.Close()

	return nil
}

