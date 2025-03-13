package post

import (
	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/admin/posts"
)

func ChangeToSecretPost(postSeq string) error {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	_, insertErr := connect.InsertQuery( queries.ChangeToSecretPostQuery, postSeq)

	if insertErr != nil {
		return insertErr
	}

	defer connect.Close()

	return nil
}

func ChangeToNotSecretPost(postSeq string) error {
	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	_, insertErr := connect.InsertQuery( queries.ChangeToNotSecretPostQuery, postSeq)

	if insertErr != nil {
		return insertErr
	}

	defer connect.Close()

	return nil
}