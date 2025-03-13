package user

import (
	"log"
	"net/http"
	"time"

	"github.com/donghquinn/blog_back_go/auth"
	crypt "github.com/donghquinn/blog_back_go/libraries/crypto"
	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/users"
	types "github.com/donghquinn/blog_back_go/types/user"
	"github.com/google/uuid"
)

func CreateLoginToken(res http.ResponseWriter, req *http.Request, loginRequst types.UserLoginRequest) types.LoginResponse {
	// 복호화
	decodeEmail, decodePassword, decodeErr := decodeLoginRequest(loginRequst)

	if decodeErr != nil {
		log.Printf("[LOGIN] Decode Requested User Info Error: %v", decodeErr)

		return types.LoginResponse{
			Status:  402,
			Code:    "02",
			Message: "Decode Login Request Error",
		}
	}

	// DB에서 유저 데이터 체크
	queryResult, queryErr := getUserInfo(loginRequst.Email)

	if queryErr != nil {

		return types.LoginResponse{
			Status:  403,
			Code:    "03",
			Message: "Query User Info Error",
		}
	}

	// 패스워드 비교 (암호화 해싱된 패스워드)
	isMatch, matchErr := crypt.PasswordCompare(queryResult.UserPassword, decodePassword)

	if matchErr != nil {
		log.Printf("[LOGIN] Match Hashed Password Error: %v", matchErr)

		return types.LoginResponse{
			Status:  404,
			Code:    "04",
			Message: "Matching User Password Error",
		}
	}

	// 패스워드 일치하지 않을 때
	if !isMatch {
		log.Printf("[LOGIN] Password Does Not Match: %v", isMatch)

		return types.LoginResponse{
			Status:  405,
			Code:    "05",
			Message: "Password Does not Match",
		}
	}

	uuid1, uuidErr1 := uuid.NewUUID()

	if uuidErr1 != nil {
		log.Printf("[REDIS] Create UUID Error: %v", uuidErr1)
		return types.LoginResponse{
			Status:  406,
			Code:    "06",
			Message: "Create UUID Error",
		}
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

		return types.LoginResponse{
			Status:  407,
			Code:    "07",
			Message: "Create JWT Token Error",
		}
	}

	uuid2, uuidErr2 := uuid.NewUUID()

	if uuidErr2 != nil {
		log.Printf("[REDIS] Create UUID Error: %v", uuidErr2)
		return types.LoginResponse{
			Status:  406,
			Code:    "06",
			Message: "Create UUID Error",
		}
	}

	// JWT 토큰 생성
	refreshToken, refreshTokenErr := auth.CreateJwtToken(queryResult.UserId, uuid2.String(), decodeEmail, queryResult.UserStatus, queryResult.BlogId, 7*24*time.Hour)

	if refreshTokenErr != nil {
		return types.LoginResponse{
			Status:  407,
			Code:    "07",
			Message: "Create JWT Token Error",
		}
	}

	userData, getRedisErr := GetRefreshUserInfoFromRedis(refreshToken)

	if getRedisErr != nil {
		log.Printf("[JWT] Get Token Error: %v", getRedisErr)

		return types.LoginResponse{
			Status:  407,
			Code:    "07",
			Message: "Create JWT Token Error",
		}
	}

	// 이미 등록된 토큰이 있다면 삭제하고 새로 등록
	if (userData != types.LoginRedisStruct{}) {
		log.Printf("[JWT] Found Already Set Token")
		deletErr := DeleteAlreadySetToken(refreshToken)

		if deletErr != nil {
			log.Printf("[JWT] Delete Token Error: %v", deletErr)

			return types.LoginResponse{
				Status:  407,
				Code:    "07",
				Message: "Create JWT Token Error",
			}
		}
	}

	setRedisErr := SetRefreshToken(refreshToken, decodeEmail, queryResult.UserStatus, queryResult.UserId, queryResult.BlogId)

	if setRedisErr != nil {

		return types.LoginResponse{
			Status:  407,
			Code:    "07",
			Message: "Create JWT Token Error",
		}
	}

	accessTokenCookie := http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Secure:   true, // 로컬 환경에서는 FALSE, 실제에서는 TRUE
		HttpOnly: true, // 로컬 환경에서는 FALSE, 실제에서는 TRUE
		SameSite: http.SameSiteNoneMode,
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		Secure:   true, // 로컬 환경에서는 FALSE, 실제에서는 TRUE
		HttpOnly: true, // 로컬 환경에서는 FALSE, 실제에서는 TRUE
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(res, &accessTokenCookie)
	http.SetCookie(res, &refreshTokenCookie)

	return types.LoginResponse{
		Status:       200,
		Code:         "0000",
		Message:      "SUCCESS",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func RefreshToken(res http.ResponseWriter, req *http.Request) types.LoginResponse {
	refreshToken, getTokenErr := RefreshTokenCookie(req)

	if getTokenErr != nil {
		return types.LoginResponse{
			Status:  401,
			Code:    "01",
			Message: "Get Refresh Token Error",
		}
	}

	userData, getErr := GetRefreshUserInfoFromRedis(refreshToken)

	if getErr != nil {
		return types.LoginResponse{
			Status:  402,
			Code:    "02",
			Message: "Get Refresh Token Error",
		}
	}

	// 이미 등록된 토큰이 있다면 삭제하고 새로 등록
	if (userData != types.LoginRedisStruct{}) {
		return types.LoginResponse{
			Status:  403,
			Code:    "03",
			Message: "No User Data Token Error",
		}
	}

	uuid1, uuidErr1 := uuid.NewUUID()

	if uuidErr1 != nil {
		log.Printf("[REDIS] Create UUID Error: %v", uuidErr1)
		return types.LoginResponse{
			Status:  406,
			Code:    "06",
			Message: "Create UUID Error",
		}
	}

	// JWT 토큰 생성
	accessToken, tokenErr := auth.CreateJwtToken(userData.UserId, uuid1.String(), userData.Email, userData.UserStatus, userData.BlogId, 3*time.Hour)

	if tokenErr != nil {
		return types.LoginResponse{
			Status:  407,
			Code:    "07",
			Message: "Create JWT Token Error",
		}
	}

	uuid2, uuidErr2 := uuid.NewUUID()

	if uuidErr2 != nil {
		log.Printf("[REDIS] Create UUID Error: %v", uuidErr2)
		return types.LoginResponse{
			Status:  406,
			Code:    "06",
			Message: "Create UUID Error",
		}
	}

	// JWT 토큰 생성
	newRefreshToken, refreshTokenErr := auth.CreateJwtToken(userData.UserId, uuid2.String(), userData.Email, userData.UserStatus, userData.BlogId, 7*24*time.Hour)

	if refreshTokenErr != nil {
		return types.LoginResponse{
			Status:  407,
			Code:    "07",
			Message: "Create JWT Token Error",
		}
	}

	setErr := SetRefreshToken(refreshToken, userData.Email, userData.UserStatus, userData.UserId, userData.BlogId)

	if setErr != nil {
		return types.LoginResponse{
			Status:  404,
			Code:    "04",
			Message: "Set Refresh Token Error",
		}
	}

	return types.LoginResponse{
		Status:       200,
		Code:         "01",
		Message:      "SUCCESS",
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}
}

func RefreshTokenCookie(request *http.Request) (string, error) {
	refreshToken, refreshErr := GetRefreshTokenFromCookie(request)

	if refreshErr != nil {
		return "", refreshErr
	}

	return refreshToken, nil
}

func GetRefreshTokenFromCookie(request *http.Request) (string, error) {
	refreshToken, refreshErr := request.Cookie("refreshToken")

	if refreshErr != nil {
		log.Printf("Get Refresh Token Error: %v", refreshErr)
		return "", refreshErr
	}

	log.Printf("Giot Refresh Token : %s", refreshToken)

	return refreshToken.Value, nil
}

func GetRefreshUserInfoFromRedis(refreshToken string) (types.LoginRedisStruct, error) {
	redisCon, redisErr := database.RedisInstance()

	if redisErr != nil {
		return types.LoginRedisStruct{}, redisErr
	}

	userData, getErr := redisCon.RedisLoginGet(refreshToken)

	if getErr != nil {
		return types.LoginRedisStruct{}, getErr
	}

	userInfo := types.LoginRedisStruct{
		Email:      userData["email"],
		UserStatus: userData["userStatus"],
		UserId:     userData["userId"],
		BlogId:     userData["blogId"],
	}
	return userInfo, nil
}

func SetRefreshToken(refreshToken string, email string, userStatus string, userId string, blogId string) error {
	redisCon, redisErr := database.RedisInstance()

	if redisErr != nil {
		return redisErr
	}

	redisSetErr := redisCon.RedisLoginSet(refreshToken, email, userStatus, userId, blogId)

	if redisSetErr != nil {
		return redisSetErr
	}

	return nil
}

func DeleteAlreadySetToken(refreshToken string) error {
	redisCon, redisErr := database.RedisInstance()

	if redisErr != nil {
		return redisErr
	}

	delErr := redisCon.Delete(refreshToken)

	if delErr != nil {
		return delErr
	}

	return nil
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
