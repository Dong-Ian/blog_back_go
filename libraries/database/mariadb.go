package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/donghquinn/blog_back_go/configs"
	"github.com/donghquinn/blog_back_go/queries"
	_ "github.com/go-sql-driver/mysql"
)

type DataBaseConnector struct {
	*sql.DB
}

// DB 연결 인스턴스
func InitDatabaseConnection() (*DataBaseConnector, error) {
	dbConfig := configs.DatabaseConfig

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database)

	// driver := sql.Open("mysql", )

	db, err := sql.Open("mysql", dbUrl)

	if err != nil {
		log.Printf("[DATABASE] Start Database Connection Error: %v", err)

		return nil, err
	}

	db.SetConnMaxLifetime(time.Second * 60)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(100)

	connect := &DataBaseConnector{db}
	return connect, nil
}

// 테이블 생성
func CheckConnection() error {
	log.Printf("Waiting for Database Connection,,,")
	time.Sleep(time.Second * 10)

	connect, dbErr := InitDatabaseConnection()

	if dbErr != nil {
		return dbErr
	}

	pingErr := connect.Ping()

	if pingErr != nil {
		log.Printf("[DATABASE] Database Ping Error: %v", pingErr)
		return pingErr
	}

	createErr := connect.CreateTable(queries.CreateTableQueryList)

	if createErr != nil {
		log.Printf("[DATABASE] Create Table Error: %v", createErr)
		return createErr
	}

	defer connect.Close()

	return nil
}

func (connect *DataBaseConnector) CreateTable(queryList []string) error {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		log.Printf("[DATABASE] Begin Transaction Error: %v", txErr)
		return txErr
	}

	defer tx.Rollback()

	for _, queryString := range queryList {
		_, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			tx.Rollback()
			log.Printf("[DATABASE] Create Table Querystring Transaction Exec Error: %v", execErr)
			return execErr
		}
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		log.Printf("[DATABASE] Commit Transaction Error: %v", commitErr)
		return commitErr
	}

	return nil
}

// 쿼리
func (connect *DataBaseConnector) GetMultiple(queryString string, args ...string) (*sql.Rows, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	result, err := connect.Query(queryString, arguments...)

	if err != nil {
		log.Printf("[QUERY] Query Error: %v\n", err)

		return nil, err
	}

	defer connect.Close()

	return result, nil
}

// 쿼리
func (connect *DataBaseConnector) QueryOne(queryString string, args ...string) (*sql.Row, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	result := connect.QueryRow(queryString, arguments...)

	if result.Err() != nil {
		log.Printf("[QUERY] Query Error: %v\n", result.Err())

		return nil, result.Err()
	}

	defer connect.Close()

	return result, nil
}

func (connect *DataBaseConnector) QueryBuilderRows(queryString string, args []interface{}) (*sql.Rows, error) {
	result, err := connect.Query(queryString, args...)

	if err != nil {
		log.Printf("[QUERY] Query Error: %v\n", err)

		return nil, err
	}

	defer connect.Close()

	return result, nil
}
func (connect *DataBaseConnector) QueryBuilderOneRow(queryString string, args []interface{}) (*sql.Row, error) {
	result := connect.QueryRow(queryString, args...)

	if result.Err() != nil {
		log.Printf("[QUERY] Query Error: %v\n", result.Err())

		return nil, result.Err()
	}

	defer connect.Close()

	return result, nil
}

// 인서트 쿼리
func (connect *DataBaseConnector) InsertQuery(queryString string, args ...string) (int64, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	insertResult, insertErr := connect.Exec(queryString, arguments...)

	if insertErr != nil {
		log.Printf("[INSERT] Insert Query Err: %v", insertErr)

		return -99999, insertErr
	}

	defer connect.Close()

	// Insert ID
	insertId, insertIdErr := insertResult.LastInsertId()

	if insertIdErr != nil {
		log.Printf("[INSERT] Get Insert ID Error: %v", insertIdErr)

		return -999999, insertIdErr
	}

	return insertId, nil
}
