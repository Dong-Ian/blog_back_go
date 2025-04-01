package user

import (
	"log"
	"net/http"

	crypt "github.com/donghquinn/blog_back_go/libraries/crypto"
	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/users"
	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/user"
	"github.com/donghquinn/blog_back_go/utils"
)

// 이메일 찾기
func SearchEmailController(res http.ResponseWriter, req *http.Request) {
	var findEmailRequest types.UserSearchEmailRequest

	parseErr := utils.DecodeBody(req, &findEmailRequest)

	if parseErr != nil {
		log.Printf("[LOGIN] Parse Body Error: %v", parseErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "Parse Body Error",
			Result:  false,
		})

		return
	}

	// 이메일 쿼리
	foundUserEmail, findErr := getUserEmail(findEmailRequest.Name)

	if findErr != nil {
		log.Printf("[LOGIN] No User Found Error: %v", findErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "No User Found Error",
			Result:  false,
		})

		return
	}

	// 이메일 복호화
	decodedEmail, decodedErr := crypt.DecryptString(foundUserEmail.UserEmail)

	if decodedErr != nil {
		log.Printf("[LOGIN] Parse Body Error: %v", decodedErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Decode Email Error",
			Result:  false,
		})

		return
	}

	response.Response(res, types.ResponseFoundEmailType{
		Status:  http.StatusOK,
		Code:    "01",
		Message: "SUCCESS",
		Result:  true,
		Email:   decodedEmail,
	})
}

func getUserEmail(userName string) (types.SelectUserSearchEmailResult, error) {
	var emailQueryResult types.SelectUserSearchEmailResult

	connect, connectErr := database.InitDatabaseConnection()

	if connectErr != nil {
		return types.SelectUserSearchEmailResult{}, connectErr
	}

	queryResult, queryErr := connect.QueryOne(queries.SelectUserEmail, userName)

	if queryErr != nil {
		return types.SelectUserSearchEmailResult{}, queryErr
	}

	queryResult.Scan(
		&emailQueryResult.UserEmail)

	return emailQueryResult, nil
}

// 페스워드 찾기
func SearchPasswordController(res http.ResponseWriter, req *http.Request) {
	var findEmailRequest types.UserSearchPasswordRequest

	parsErr := utils.DecodeBody(req, &findEmailRequest)

	if parsErr != nil {
		log.Printf("[LOGIN] Parse Body Error: %v", parsErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "Parse Body Error",
			Result:  false,
		})

		return
	}

	// 패스워드 쿼리
	foundUserPassword, findErr := getUserPassword(findEmailRequest.Email, findEmailRequest.Name)

	if findErr != nil {
		log.Printf("[LOGIN] Not Found User Error: %v", findErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Not User Found Error",
			Result:  false,
		})

		return
	}

	response.Response(res, types.ResponseFoundPasswdType{
		Status:   http.StatusOK,
		Code:     "01",
		Message:  "SUCCESS",
		Result:   true,
		Password: foundUserPassword.UserPassword,
	})
}

func getUserPassword(userEmail string, userName string) (types.SelectUserSearchPasswordResult, error) {
	var emailQueryResult types.SelectUserSearchPasswordResult

	connect, connectErr := database.InitDatabaseConnection()

	if connectErr != nil {
		return types.SelectUserSearchPasswordResult{}, connectErr
	}

	queryResult, queryErr := connect.QueryOne(queries.SelectUserPassword, userName, userEmail)

	if queryErr != nil {
		return types.SelectUserSearchPasswordResult{}, queryErr
	}

	queryResult.Scan(
		&emailQueryResult.UserPassword)

	return emailQueryResult, nil
}
