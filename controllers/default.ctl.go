package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/donghquinn/blog_back_go/auth"
	"github.com/donghquinn/blog_back_go/dto"
	"github.com/donghquinn/blog_back_go/middlewares"
	"github.com/donghquinn/blog_back_go/types"
	"github.com/google/uuid"
)

func DefaultController(res http.ResponseWriter, req *http.Request) {
	_, _, _, _, err := auth.ValidateJwtToken(req)

	if err != nil {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", err)

		return
	}

	dto.SetResponseWithMessage(res, 200, "01", "Hello World")
}

func RefreshController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", nil)

		return
	}

	uuid1, uuidErr1 := uuid.NewUUID()

	if uuidErr1 != nil {
		log.Printf("[REDIS] Create UUID Error: %v", uuidErr1)
		dto.SetErrorResponse(res, 406, "06", "Create Uuid Error", uuidErr1)
	}

	uuid2, uuidErr2 := uuid.NewUUID()

	if uuidErr2 != nil {
		log.Printf("[REDIS] Create UUID Error: %v", uuidErr2)
		dto.SetErrorResponse(res, 406, "06", "Create Uuid Error", uuidErr2)
	}

	// JWT 토큰 생성
	accessToken, tokenErr := auth.CreateJwtToken(user.UserId, uuid1.String(), user.UserEmail, user.UserStatus, user.BlogId, 3*time.Hour)

	if tokenErr != nil {
		dto.SetErrorResponse(res, 407, "07", "Create JWT Token Error", tokenErr)
		return
	}

	// JWT 토큰 생성
	refreshToken, tokenErr := auth.CreateJwtToken(user.UserId, uuid2.String(), user.UserEmail, user.UserStatus, user.BlogId, 7*24*time.Hour)

	if tokenErr != nil {
		dto.SetErrorResponse(res, 407, "07", "Create JWT Token Error", tokenErr)
		return
	}

	dto.SetTokenResponse(res, 200, "01", types.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

func CorsTestController(res http.ResponseWriter, req *http.Request) {
	dto.SetResponseWithMessage(res, 200, "01", "Hi")
}
