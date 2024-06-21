package post

import (
	"database/sql"
	"log"

	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/posts"
	types "github.com/donghquinn/blog_back_go/types/post"
)

func GetAllCategoryList() ([]string, error) {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return []string{}, dbErr
	}

	queryResult, queryErr := database.Query(connect, queries.SelectAllCategories)

	if queryErr != nil {
		log.Printf("[CATEGORY] Get All Categories Error: %v", queryErr)

		return []string{}, queryErr
	}

	var categoryList []string

	for queryResult.Next() {
		var row types.CategoryQueryResult

		scanErr := queryResult.Scan(&row.CategoryName)

		if scanErr != nil {
			if scanErr == sql.ErrNoRows {
				categoryList = make([]string, 0)
			} else {
				return []string{}, scanErr
			}
		}

		categoryList = append(categoryList, row.CategoryName)
	}

	return categoryList, nil
}