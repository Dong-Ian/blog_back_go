package dto

import (
	"encoding/json"
	"net/http"

	"github.com/donghquinn/blog_back_go/types"
)

// 기본 응답
func SetResponse(res http.ResponseWriter, statusCode int, code string) {
	responseObject, _ := json.Marshal(types.ResponseType{Result: true, Code: code})

	res.WriteHeader(200)
	res.Write(responseObject)
}

func SetSignupResponse(res http.ResponseWriter, statusCode int, code string, blogId string) {
	responseObject, _ := json.Marshal(types.ResponseSignupType{Result: true, Code: code, BlogId: blogId})

	res.WriteHeader(200)
	res.Write(responseObject)
}

func SetImageUrlResponse(res http.ResponseWriter, statusCode int, code string, urls []string) {
	responseObject, _ := json.Marshal(types.ResponseImageUrl{Code: code, Result: true, ImageResult: urls})

	res.WriteHeader(200)
	res.Write(responseObject)
}

// 기본 응답
func SetProfileResponse(res http.ResponseWriter, statusCode int, code string, profile types.UserProfileDataResponseType) {
	responseObject, _ := json.Marshal(types.ResponseProfileType{Code: code, Result: true, ProfileResult: profile})

	res.WriteHeader(200)
	res.Write(responseObject)
}

// 메세지 담은 응답
func SetResponseWithMessage(res http.ResponseWriter, statusCode int, code string, message string) {
	responseObject, _ := json.Marshal(types.ResponseMessageType{Code: code, Result: true, Message: message})

	res.WriteHeader(200)
	res.Write(responseObject)
}

// 토큰 담은 메세지
func SetTokenResponse(res http.ResponseWriter, statusCode int, code string, token types.LoginResponse) {
	responseObject, _ := json.Marshal(types.ResponseTokenType{Code: code, Result: true, AccessToken: token.AccessToken, RefreshToken: token.RefreshToken})

	res.WriteHeader(200)
	res.Write(responseObject)
}

// 이메일 담음 응답
func SetEmailResponse(res http.ResponseWriter, statusCode int, code string, email string) {
	responseObject, _ := json.Marshal(types.ResponseFoundEmailType{Code: code, Result: true, Email: email})

	res.WriteHeader(200)
	res.Write(responseObject)
}

func SetPasswordResponse(res http.ResponseWriter, statusCode int, code string, password string) {
	responseObject, _ := json.Marshal(types.ResponseFoundPasswordType{Code: code, Result: true, Password: password})

	res.WriteHeader(200)
	res.Write(responseObject)
}

// 에러 응답
func SetErrorResponse(res http.ResponseWriter, statusCode int, code string, message string, err error) {
	responseObject, _ := json.Marshal(types.ErrorResponseType{Code: code, Result: false, Message: message})

	res.WriteHeader(500)
	res.Write(responseObject)
}
