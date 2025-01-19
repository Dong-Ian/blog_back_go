package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/donghquinn/blog_back_go/auth"
	"github.com/donghquinn/blog_back_go/dto"
	crypt "github.com/donghquinn/blog_back_go/libraries/crypto"
	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/users"
	"github.com/donghquinn/blog_back_go/types"
	"github.com/donghquinn/blog_back_go/utils"
	"github.com/google/uuid"
)

func LoginController(res http.ResponseWriter, req *http.Request) {
	var loginRequst types.UserLoginRequest

	parseErr := utils.DecodeBody(req, &loginRequst)

	if parseErr != nil {
		log.Printf("[LOGIN] Parse Body Error: %v", parseErr)

		dto.SetErrorResponse(res, 401, "01", "SignUp Parsing Error", parseErr)
		return
	}

	// 복호화
	decodeEmail, decodePassword, decodeErr := decodeLoginRequest(loginRequst)

	if decodeErr != nil {
		log.Printf("[LOGIN] Decode Requested User Info Error: %v", decodeErr)
		dto.SetErrorResponse(res, 402, "02", "Decode Login Request", decodeErr)
		return
	}

	// DB에서 유저 데이터 체크
	queryResult, queryErr := getUserInfo(loginRequst.Email)

	if queryErr != nil {
		dto.SetErrorResponse(res, 403, "03", "Query User Info Error", queryErr)
		return
	}

	// 패스워드 비교 (암호화 해싱된 패스워드)
	isMatch, matchErr := crypt.PasswordCompare(queryResult.UserPassword, decodePassword)

	if matchErr != nil {
		log.Printf("[LOGIN] Match Hashed Password Error: %v", matchErr)
		dto.SetErrorResponse(res, 404, "04", "Matching User Password Error", matchErr)
		return
	}

	// 패스워드 일치하지 않을 때
	if !isMatch {
		log.Printf("[LOGIN] Password Does Not Match: %v", isMatch)
		dto.SetErrorResponse(res, 405, "05", "Password Does not Match", fmt.Errorf("password does not match"))
		return
	}

	uuid1, uuidErr1 := uuid.NewUUID()

	if uuidErr1 != nil {
		log.Printf("[REDIS] Create UUID Error: %v", uuidErr1)
		dto.SetErrorResponse(res, 406, "06", "Create Uuid Error", uuidErr1)
	}

	// dbCon, dbErr := database.InitDatabaseConnection()

	// if dbErr != nil {
	// 	log.Printf("[JWT] Start Db CONNECTION Error: %v", dbErr)
	// 	dto.SetErrorResponse(res, 408, "08", "Insert Session Data Error", dbErr)
	// }

	// insertId, insertErr := database.InsertQuery(dbCon, queries.InsertSessionData, userId)

	// if insertErr != nil {
	// 	log.Printf("[JWT] Insert Seesion Data Error")
	// }

	// JWT 토큰 생성
	accessToken, tokenErr := auth.CreateJwtToken(queryResult.UserId, uuid1.String(), decodeEmail, queryResult.UserStatus, queryResult.BlogId, 3*time.Hour)

	if tokenErr != nil {
		dto.SetErrorResponse(res, 407, "07", "Create JWT Token Error", tokenErr)
		return
	}

	uuid2, uuidErr2 := uuid.NewUUID()

	if uuidErr2 != nil {
		log.Printf("[REDIS] Create UUID Error: %v", uuidErr2)
		dto.SetErrorResponse(res, 406, "06", "Create Uuid Error", uuidErr2)
	}

	// JWT 토큰 생성
	refreshToken, tokenErr := auth.CreateJwtToken(queryResult.UserId, uuid2.String(), decodeEmail, queryResult.UserStatus, queryResult.BlogId, 7*24*time.Hour)

	if tokenErr != nil {
		dto.SetErrorResponse(res, 407, "07", "Create JWT Token Error", tokenErr)
		return
	}

	dto.SetTokenResponse(res, 200, "01", types.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

func decodeLoginRequest(loginRequest types.UserLoginRequest) (string, string, error) {
	decodeEmail, decodeEmailErr := crypt.DecryptString(loginRequest.Email)

	if decodeEmailErr != nil {
		log.Printf("[LOGIN] Decode Email Err: %v", decodeEmailErr)
		return "", "", decodeEmailErr
	}

	decodePassword, decodePassErr := crypt.DecryptString(loginRequest.Password)

	if decodePassErr != nil {
		log.Printf("[LOGIN] Decode Password Err: %v", decodePassErr)
		return "", "", decodePassErr
	}

	return decodeEmail, decodePassword, nil
}

// func insertSessionData(userId string) {

// }

func getUserInfo(encodedEmail string) (types.UserLoginQueryResult, error) {
	connect, connectErr := database.InitDatabaseConnection()

	if connectErr != nil {
		return types.UserLoginQueryResult{}, connectErr
	}

	result, queryErr := connect.QueryOne(queries.SelectUserInfo, encodedEmail)

	if queryErr != nil {
		return types.UserLoginQueryResult{}, queryErr
	}

	defer connect.Close()

	var queryUserInfoResult types.UserLoginQueryResult

	result.Scan(
		&queryUserInfoResult.UserId,
		&queryUserInfoResult.UserPassword,
		&queryUserInfoResult.UserStatus,
		&queryUserInfoResult.BlogId)

	return queryUserInfoResult, nil
}
