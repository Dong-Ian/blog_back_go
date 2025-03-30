package biz

import (
	"log"
	"net/http"
	"time"

	"github.com/donghquinn/blog_back_go/auth"
	"github.com/donghquinn/blog_back_go/middlewares"
	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/user"
)

func DefaultController(res http.ResponseWriter, req *http.Request) {
	_, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "JWT Verifying Error",
			Result:  false,
		})

		return
	}

	response.Response(res, response.CommonResponseWithMessage{
		Status:  http.StatusOK,
		Code:    "01",
		Message: "Hello, World!",
		Result:  true,
	})
}

func RefreshController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {

		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "JWT Verifying Error",
			Result:  false,
		})

		return
	}

	// JWT 토큰 생성
	accessToken, tokenErr := auth.CreateJwtToken(user.UserId, user.UserEmail, user.UserStatus, user.BlogId, 3*time.Hour)

	if tokenErr != nil {
		log.Printf("[REFRESH] Create Token Error: %v", tokenErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  407,
			Code:    "07",
			Message: "Create Token Error",
			Result:  false,
		})

		return
	}

	// JWT 토큰 생성
	refreshToken, tokenErr := auth.CreateJwtToken(user.UserId, user.UserEmail, user.UserStatus, user.BlogId, 7*24*time.Hour)

	if tokenErr != nil {
		log.Printf("[REFRESH] Create Refresh Token Error: %v", tokenErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  408,
			Code:    "08",
			Message: "Create Refresh Token Error",
			Result:  false,
		})

		return
	}

	response.Response(res, types.LoginResponse{
		Status:       http.StatusOK,
		Code:         "01",
		Message:      "SUCCESS",
		Result:       true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func CorsTestController(res http.ResponseWriter, req *http.Request) {

	response.Response(res, response.CommonResponseWithMessage{
		Status:  http.StatusOK,
		Code:    "01",
		Message: "Hi",
		Result:  true,
	})
}
