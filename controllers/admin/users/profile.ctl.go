package admincontrollers

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/dto"
	"github.com/donghquinn/blog_back_go/libraries/profile"
	"github.com/donghquinn/blog_back_go/middlewares"
	types "github.com/donghquinn/blog_back_go/types/admin/users"
	"github.com/donghquinn/blog_back_go/utils"
)

// 프로필 변경 컨트롤러
func UpdateProfileController(res http.ResponseWriter, req *http.Request) {
	var updateProfile types.UserChangeProfileRequest

	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", nil)
		return
	}

	parseErr := utils.DecodeBody(req, &updateProfile)

	if parseErr != nil {
		log.Printf("[PROFILE] Change Profile Request Error: %v", parseErr)
		dto.SetErrorResponse(res, 402, "02", "Change Profile Request Error", parseErr)
		return
	}

	updateErr := profile.ChangeProfile(updateProfile, user.UserId, user.BlogId)

	if updateErr != nil {
		dto.SetErrorResponse(res, 403, "03", "Insert Profile Update Error", updateErr)
		return
	}

	dto.SetResponse(res, 200, "01")
}

// 색상 변경 컨트롤러
func UpdateColorController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", nil)

		return
	}

	var changeColorRequest types.UserUpdateProfileColorRequest

	parseErr := utils.DecodeBody(req, &changeColorRequest)

	if parseErr != nil {
		log.Printf("[COLOR] Change Color Request Error: %v", parseErr)
		dto.SetErrorResponse(res, 402, "02", "Change Color Request Error", parseErr)
		return
	}

	changeColorErr := profile.ChangeColor(changeColorRequest, user.UserId, user.BlogId)

	if changeColorErr != nil {
		dto.SetErrorResponse(res, 403, "03", "Change Color Error", changeColorErr)
		return
	}

	dto.SetResponse(res, 200, "01")
}

// 블로그 타이틀 변경 컨트롤러
func UpdateTitleController(res http.ResponseWriter, req *http.Request) {
	user, ok := middlewares.GetUserFromContext(req.Context())

	if !ok {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", nil)

		return
	}

	var changeTitleRequest types.UserUpdateBlogTitleRequest

	parseErr := utils.DecodeBody(req, &changeTitleRequest)

	if parseErr != nil {
		log.Printf("[TITLE] Change Title Request Error: %v", parseErr)
		dto.SetErrorResponse(res, 402, "02", "Change Title Request Error", parseErr)
		return
	}

	changeTitleErr := profile.ChangeBlogTitle(changeTitleRequest, user.UserId, user.BlogId)

	if changeTitleErr != nil {
		dto.SetErrorResponse(res, 403, "03", "Change Title Error", changeTitleErr)
		return
	}

	dto.SetResponse(res, 200, "01")
}
