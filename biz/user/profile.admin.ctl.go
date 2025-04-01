package user

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/middlewares"
	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/user"
	"github.com/donghquinn/blog_back_go/utils"
)

// 프로필 변경 컨트롤러
func UpdateProfileController(res http.ResponseWriter, req *http.Request) {
	var updateProfile types.UserChangeProfileRequest

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

	parseErr := utils.DecodeBody(req, &updateProfile)

	if parseErr != nil {
		log.Printf("[PROFILE] Change Profile Request Error: %v", parseErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Change Profile Request Error",
			Result:  false,
		})

		return
	}

	updateErr := ChangeProfile(updateProfile, user.UserId, user.BlogId)

	if updateErr != nil {
		log.Printf("[UPLOAD_PROFILE] Insert Profile Update Error: %v", updateErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Insert Profile Update Error",
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

// 색상 변경 컨트롤러
func UpdateColorController(res http.ResponseWriter, req *http.Request) {
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

	var changeColorRequest types.UserUpdateProfileColorRequest

	parseErr := utils.DecodeBody(req, &changeColorRequest)

	if parseErr != nil {
		log.Printf("[COLOR] Change Color Request Error: %v", parseErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Chnage Color Request Error",
			Result:  false,
		})

		return
	}

	changeColorErr := ChangeColor(changeColorRequest, user.UserId, user.BlogId)

	if changeColorErr != nil {
		log.Printf("[UPLOAD_PROFILE] Change Color Error: %v", changeColorErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Change Color Error",
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

// 블로그 타이틀 변경 컨트롤러
func UpdateTitleController(res http.ResponseWriter, req *http.Request) {
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

	var changeTitleRequest types.UserUpdateBlogTitleRequest

	parseErr := utils.DecodeBody(req, &changeTitleRequest)

	if parseErr != nil {
		log.Printf("[TITLE] Change Title Request Error: %v", parseErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Chnage Title Requeset Error",
			Result:  false,
		})

		return
	}

	changeTitleErr := ChangeBlogTitle(changeTitleRequest, user.UserId, user.BlogId)

	if changeTitleErr != nil {
		log.Printf("[UPLOAD_PROFILE] File Getting Error: %v", changeTitleErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Change Title Error",
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
