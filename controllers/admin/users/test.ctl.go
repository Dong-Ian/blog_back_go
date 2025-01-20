package admincontrollers

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/dto"
	"github.com/donghquinn/blog_back_go/middlewares"
)

// 프로필 변경 컨트롤러
func CheckTokenController(res http.ResponseWriter, req *http.Request) {
	_, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", nil)
		return
	}

	dto.SetResponse(res, 200, "01")
}
