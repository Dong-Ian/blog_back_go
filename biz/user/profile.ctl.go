package user

import (
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/dto"
	crypt "github.com/donghquinn/blog_back_go/libraries/crypto"
	types "github.com/donghquinn/blog_back_go/types/admin/users"
	"github.com/donghquinn/blog_back_go/utils"
)

func GetUserProfileController(res http.ResponseWriter, req *http.Request) {
	var getUserProfileRequest types.UserGetProfileRequest

	parseErr := utils.DecodeBody(req, &getUserProfileRequest)

	if parseErr != nil {
		log.Printf("[LOGIN] Parse Body Error: %v", parseErr)

		dto.SetErrorResponse(res, 401, "01", "SignUp Parsing Error", parseErr)
		return
	}

	profile, querErr := GetUserProfile(getUserProfileRequest.BlogId, getUserProfileRequest.UserId)

	if querErr != nil {
		dto.SetErrorResponse(res, 402, "02", "Profile query Error", querErr)
		return
	}

	decodedName, nameErr := crypt.DecryptString(profile.UserName)

	if nameErr != nil {
		log.Printf("[PROFILE] Decode User Name: %v", nameErr)
		dto.SetErrorResponse(res, 403, "03", "Decode User Name Error", nameErr)
		return
	}

	decodedEmail, emailErr := crypt.DecryptString(profile.UserEmail)

	if emailErr != nil {
		log.Printf("[PROFILE] Decode User Email: %v", emailErr)
		dto.SetErrorResponse(res, 404, "04", "Decode User Email Error", emailErr)
		return
	}

	profile.UserName = decodedName
	profile.UserEmail = decodedEmail

	dto.SetProfileResponse(res, 200, "01", profile)
}
