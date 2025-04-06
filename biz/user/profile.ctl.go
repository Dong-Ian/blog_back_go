package user

import (
	"log"
	"net/http"

	crypt "github.com/donghquinn/blog_back_go/libraries/crypto"
	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/user"

	"github.com/donghquinn/blog_back_go/utils"
)

func GetUserProfileController(res http.ResponseWriter, req *http.Request) {
	blogId, getErr := utils.GetBlogIdFromContext(req.Context())

	if !getErr {
		response.Response(res, response.CommonResponseWithMessage{
			Status:  http.StatusBadRequest,
			Code:    "01",
			Message: "Get BlogId Error from context",
			Result:  false,
		})

		return
	}

	profile, querErr := GetUserProfile(blogId)

	if querErr != nil {
		log.Printf("[UPLOAD_PROFILE] Profile Query Error: %v", querErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "Profile Query Error",
			Result:  false,
		})

		return
	}

	decodedName, nameErr := crypt.DecryptString(profile.UserName)

	if nameErr != nil {
		log.Printf("[PROFILE] Decode User Name: %v", nameErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Decode User Name Error",
			Result:  false,
		})

		return
	}

	decodedEmail, emailErr := crypt.DecryptString(profile.UserEmail)

	if emailErr != nil {
		log.Printf("[PROFILE] Decode User Email: %v", emailErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  404,
			Code:    "04",
			Message: "Decode User Email Error",
			Result:  false,
		})

		return
	}

	profile.UserName = decodedName
	profile.UserEmail = decodedEmail

	response.Response(res, types.ResponseProfileType{
		Status:        http.StatusOK,
		Code:          "01",
		Message:       "SUCCESS",
		Result:        true,
		ProfileResult: profile,
	})
}
