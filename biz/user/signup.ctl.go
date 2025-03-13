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
	"github.com/google/uuid"
)

func SignupController(res http.ResponseWriter, req *http.Request) {
	var signupRequestBody types.UserSignupRequest

	// BODY 파싱
	parseErr := utils.DecodeBody(req, &signupRequestBody)

	if parseErr != nil {
		log.Printf("[SIGN_UP] Parse Body Error: %v", parseErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "Parse Body Error",
			Result:  false,
		})

		return
	}

	// 요청 회원가입 데이터 복호화 (패스워드는 암호화의 암호화된 상태로 전달됨)
	decodedEmail, decodedName, decodedPassword, decodeErr := decodeSignupUserRequest(signupRequestBody)

	// log.Printf("[SIGNUP] decodedEmail: %s, decodedName: %s, decodedPassword: %s", decodedEmail, decodedName, decodedPassword)

	if decodeErr != nil {
		log.Printf("[LOGIN] Not Found User Error: %v", decodeErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Decode Received User Info Error",
			Result:  false,
		})

		return
	}

	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		log.Printf("[LOGIN] Not Found User Error: %v", dbErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Database Connect Error",
			Result:  false,
		})

		return
	}

	// 암호화해서 업로드
	userId, encodedEmail, encodedName, encodedPassword, encodeErr := encodeSignupUserInfo(decodedEmail, decodedPassword, decodedName)

	if encodeErr != nil {
		log.Printf("[LOGIN] Not Found User Error: %v", encodeErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  404,
			Code:    "04",
			Message: "Encoding Process Error",
			Result:  false,
		})

		return
	}

	// log.Printf("[SIGNUP] userId: %s, encodedEmail: %s, encodedName: %s, encodedPassword: %s",userId, encodedEmail, encodedName, encodedPassword)
	// 새로운 유저 데이터 입력
	_, insertErr := connect.InsertQuery(queries.InsertSignupUser, userId, encodedEmail, encodedPassword, encodedName, signupRequestBody.BlogId)

	if insertErr != nil {
		log.Printf("[LOGIN] Insert New User Info Error: %v", insertErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  405,
			Code:    "05",
			Message: "Insert New User Info Error",
			Result:  false,
		})

		return
	}

	response.Response(res, response.CommonResponseWithMessage{
		Status:  http.StatusOK,
		Code:    "01",
		Message: "SUCCESS",
		Result:  true,
	})
}

func decodeSignupUserRequest(signupRequest types.UserSignupRequest) (string, string, string, error) {
	decodedEmail, decodeEmailErr := crypt.DecryptString(signupRequest.Email)

	if decodeEmailErr != nil {
		log.Printf("[SIGNUP] Decode Email Error: %v", decodeEmailErr)
		return "", "", "", decodeEmailErr
	}

	decodedName, decodeNameErr := crypt.DecryptString(signupRequest.Name)

	if decodeNameErr != nil {
		log.Printf("[SIGNUP] Decode Name Error: %v", decodeNameErr)
		return "", "", "", decodeNameErr
	}

	decodedPassword, decodePasswordErr := crypt.DecryptString(signupRequest.Password)

	if decodePasswordErr != nil {
		log.Printf("[SIGNUP] Decode Password Error: %v", decodePasswordErr)
		return "", "", "", decodePasswordErr
	}

	return decodedEmail, decodedName, decodedPassword, nil
}

// 인코딩
func encodeSignupUserInfo(decodeEmail string, decodePassword string, decodeName string) (string, string, string, string, error) {
	userId, uuidErr := uuid.NewV7()

	if uuidErr != nil {
		log.Printf("[SIGN_UP] Creating User UUID Error: %v", uuidErr)

		return "", "", "", "", uuidErr
	}

	encodedEmail, encodeEmailErr := crypt.EncryptString(decodeEmail)

	if encodeEmailErr != nil {
		return "", "", "", "", encodeEmailErr
	}

	encodedName, encodeNameErr := crypt.EncryptString(decodeName)

	if encodeNameErr != nil {
		return "", "", "", "", encodeNameErr
	}

	encodedPassword, encodePasswordErr := crypt.EncryptHashPassword(decodePassword)

	if encodePasswordErr != nil {
		return "", "", "", "", encodePasswordErr
	}

	return userId.String(), encodedEmail, encodedName, encodedPassword, nil
}
