package upload

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/donghquinn/blog_back_go/libraries/database"
	"github.com/donghquinn/blog_back_go/middlewares"
	queries "github.com/donghquinn/blog_back_go/queries/upload"
	"github.com/donghquinn/blog_back_go/response"
	types "github.com/donghquinn/blog_back_go/types/post"
)

// 게시글 이미지 업로드
func UploadPostImageController(res http.ResponseWriter, req *http.Request) {
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

	// 요청으로부터 이미지 파일 가져오기
	file, handler, fileErr := GetImagefileFromRequest(res, req)

	if fileErr != nil {
		log.Printf("[UPLOAD_POST] File Getting Error: %v", fileErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  402,
			Code:    "02",
			Message: "File Getting Error",
			Result:  false,
		})

		return
	}

	// 파일 생성
	tempFile, tempErr := CreateFileImage(res, req, file, handler)

	if tempErr != nil {
		log.Printf("[UPLOAD_POST] Create Temp Image File Error: %v", tempErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  403,
			Code:    "03",
			Message: "Create Temp Image File",
			Result:  false,
		})

		return
	}

	contentType := handler.Header["Content-Type"][0]

	// 이미지 업로드 - minio

	_, uploadErr := database.UploadImage(handler.Filename, tempFile.Name(), contentType)

	if uploadErr != nil {
		log.Printf("[UPLOAD_POST] Upload Image Error: %v", uploadErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  404,
			Code:    "04",
			Message: "Upload Image Error",
			Result:  false,
		})

		return
	}

	connect, _ := database.InitDatabaseConnection()

	var insertId int64

	// 데이터 입력 - DB
	seq, insertErr := connect.InsertQuery(
		queries.InsertPostImageData,
		// USER ID from JWT
		"1",
		user.UserId,
		"post_table",
		"POST_IMAGE",
		strconv.Itoa(int(handler.Size)),
		handler.Filename,
		contentType)

	if insertErr != nil {
		log.Printf("[UPLOAD_POST] Insert Image Info Error: %v", insertErr)
		response.Response(res, response.CommonResponseWithMessage{
			Status:  405,
			Code:    "05",
			Message: "Insert Image Info Error",
			Result:  false,
		})

		return
	}

	insertId = seq

	defer connect.Close()

	removeErr := os.Remove(tempFile.Name())

	if removeErr != nil {
		log.Printf("[UPLOAD] Remove Saved Image Error: %v", removeErr)

		response.Response(res, response.CommonResponseWithMessage{
			Status:  406,
			Code:    "06",
			Message: "Removing Saved Image Error",
			Result:  false,
		})
		return
	}

	response.Response(res, types.ResponseInsertIdType{
		Status:   http.StatusOK,
		Code:     "01",
		Message:  "SUCCESS",
		Result:   true,
		InsertId: fmt.Sprintf("%d", insertId),
	})
}
