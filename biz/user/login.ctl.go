package user

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/dto"
	"github.com/donghquinn/blog_back_go/libraries"
	"github.com/donghquinn/blog_back_go/response"
	"github.com/donghquinn/blog_back_go/types"
	"github.com/donghquinn/blog_back_go/utils"
)

func LoginController(res http.ResponseWriter, req *http.Request) {
	var loginRequst types.UserLoginRequest

	parseErr := utils.DecodeBody(req, &loginRequst)

	if parseErr != nil {
		log.Printf("[LOGIN] Parse Body Error: %v", parseErr)

		dto.SetErrorResponse(res, 401, "01", "SignUp Parsing Error", parseErr)
		return
	}

	// res.Header().Set("Set-Cookie", accessTokenCookie.String())
	// res.Header().Add("Set-Cookie", refreshTokenCookie.String())

	result := libraries.CreateLoginToken(res, req, loginRequst)
	response.Response(res, result)
}
