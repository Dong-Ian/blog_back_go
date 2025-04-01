package user

import (
	"log"

	"github.com/donghquinn/blog_back_go/configs"
	"github.com/donghquinn/blog_back_go/libraries/database"
	queries "github.com/donghquinn/blog_back_go/queries/users"
	"github.com/donghquinn/blog_back_go/types"
)

func GetUserProfile(blogId string, userId string) (types.UserProfileDataResponseType, error) {
	var userProfileResult types.UserProfileDataResponseType

	userProfileData, profileErr := GetDefaultUserProfile(blogId)

	if profileErr != nil {
		return types.UserProfileDataResponseType{}, profileErr
	}

	imageUrlList, imageUrlErr := GetUserProfileImageList(userProfileData.UserId)

	if imageUrlErr != nil {
		return types.UserProfileDataResponseType{}, imageUrlErr
	}

	userProfileResult = types.UserProfileDataResponseType{
		UserId:      userProfileData.UserId,
		UserName:    userProfileData.UserName,
		UserEmail:   userProfileData.UserEmail,
		Color:       userProfileData.Color,
		Title:       userProfileData.Title,
		Instagram:   userProfileData.Instagram,
		GithubUrl:   userProfileData.GithubUrl,
		PersonalUrl: userProfileData.PersonalUrl,
		Memo:        userProfileData.Memo,
		Images:      imageUrlList,
	}

	return userProfileResult, nil
}

// 기본 유저 구하기
func GetDefaultUserProfile(blogId string) (types.SelectUserProfileQueryResult, error) {
	var userProfileData types.SelectUserProfileQueryResult

	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return types.SelectUserProfileQueryResult{}, dbErr
	}

	profile, profileErr := connect.QueryOne(queries.SelectUserProfile, blogId)

	if profileErr != nil {
		log.Printf("[PROFILE] Get Profile Error: %v", profileErr)
		return types.SelectUserProfileQueryResult{}, profileErr
	}

	defer connect.Close()

	profile.Scan(
		&userProfileData.UserId,
		&userProfileData.UserEmail,
		&userProfileData.UserName,
		&userProfileData.Color,
		&userProfileData.Title,
		&userProfileData.Instagram,
		&userProfileData.GithubUrl,
		&userProfileData.PersonalUrl,
		&userProfileData.Memo)

	return userProfileData, nil
}

// 유저 ID 로 조회하기
func GetUserProfileByUserId(userId string) (types.SelectUserProfileQueryResult, error) {
	var userProfileData types.SelectUserProfileQueryResult

	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return types.SelectUserProfileQueryResult{}, dbErr
	}

	profile, profileErr := connect.QueryOne(queries.SelectUserProfileByUserId, userId)

	if profileErr != nil {
		log.Printf("[PROFILE] Get Profile Error: %v", profileErr)
		return types.SelectUserProfileQueryResult{}, profileErr
	}

	defer connect.Close()

	profile.Scan(
		&userProfileData.UserId,
		&userProfileData.UserEmail,
		&userProfileData.UserName,
		&userProfileData.Color,
		&userProfileData.Title,
		&userProfileData.Instagram,
		&userProfileData.GithubUrl,
		&userProfileData.PersonalUrl,
		&userProfileData.Memo)

	return userProfileData, nil
}

func GetUserProfileImageList(userId string) (types.UserImageFileData, error) {
	minioConfig := configs.MinioConfig
	// 이미지 데이터 url 가져오기 시작
	// var userImageData []types.SelectFileQueryResult
	var imageUrlList types.UserImageFileData

	connect, dbErr := database.InitDatabaseConnection()

	if dbErr != nil {
		return types.UserImageFileData{}, dbErr
	}

	images, imagesErr := connect.GetMultiple(queries.SelectUserProfileProfileAndBackground, userId, "USER_PROFILE", "USER_BACKGROUND")

	if imagesErr != nil {
		log.Printf("[PROFILE] Get Profile And Background Images Error: %v", imagesErr)
		return types.UserImageFileData{}, imagesErr
	}

	defer connect.Close()

	for images.Next() {
		var row types.SelectFileQueryResult

		if scanErr := images.Scan(
			&row.FileFormat,
			&row.FileType,
			&row.TargetPurpose,
			&row.TargetId,
			&row.ObjectName,
		); scanErr != nil {
			log.Printf("[PROFILE] Scan Error: %v", scanErr)
			return imageUrlList, scanErr
		}
		if row.TargetPurpose == "USER_BACKGROUND" {
			imageUrlList.BackgroundImage = "https://" + minioConfig.HostUrl + "/" + minioConfig.BlogBucket + "/" + row.ObjectName
			log.Printf("[DEBUGGING] Get Profile Image: %s", imageUrlList.BackgroundImage)
		}

		if row.TargetPurpose == "USER_PROFILE" {
			imageUrlList.ProfileImage = "https://" + minioConfig.HostUrl + "/" + minioConfig.BlogBucket + "/" + row.ObjectName
			log.Printf("[DEBUGGING] Get Background Image: %s", imageUrlList.ProfileImage)

		}

	}

	// imageUrls, urlErr := getImages(userImageData)

	// imageUrlList = imageUrls

	// if urlErr != nil {
	// 	return types.UserImageFileData{}, urlErr
	// }

	return imageUrlList, nil
}

func getImages(imageData []types.SelectFileQueryResult) (types.UserImageFileData, error) {
	var imageUrlData types.UserImageFileData

	for _, row := range imageData {
		imageUrl, err := database.GetImageUrl(row.ObjectName, row.FileType)

		if err != nil {
			log.Printf("[PROFILE] Get Profile IMAGE URL: %v", err)
			return types.UserImageFileData{}, err
		}

		if row.TargetPurpose == "USER_PROFILE" {
			imageUrlData.ProfileImage = imageUrl.String()
		}

		if row.TargetPurpose == "USER_BACKGROUND" {
			imageUrlData.BackgroundImage = imageUrl.String()
		}
	}

	return imageUrlData, nil
}
