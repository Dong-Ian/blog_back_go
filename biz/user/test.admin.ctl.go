package user

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/middlewares"
	"github.com/donghquinn/blog_back_go/response"
)

// 프로필 변경 컨트롤러
func CheckTokenController(res http.ResponseWriter, req *http.Request) {
	_, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {

		response.Response(res, response.CommonResponseWithMessage{
			Status:  401,
			Code:    "01",
			Message: "JWT Verifying Eror",
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
