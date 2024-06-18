package admincontrollers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/donghquinn/blog_back_go/auth"
	"github.com/donghquinn/blog_back_go/dto"
	"github.com/donghquinn/blog_back_go/libraries/database"
	upload "github.com/donghquinn/blog_back_go/libraries/upload/image"
	queries "github.com/donghquinn/blog_back_go/queries/upload"
)

// 게시글 이미지 업로드
func UploadPostImageController(res http.ResponseWriter, req *http.Request) {
	userId, _, _, err := auth.ValidateJwtToken(req)

	if err != nil {
		dto.SetErrorResponse(res, 401, "01", "JWT Verifying Error", err)

		return
	}

	// 요청으로부터 이미지 파일 가져오기
	file, handler, fileErr := upload.GetImagefileFromRequest(res, req)

	if fileErr != nil {
		dto.SetErrorResponse(res, 402, "02", "File Getting Error", fileErr)

		return
	}

	// 파일 생성
	tempFile, tempErr := upload.CreateFileImage(res, req, file, handler)

	if tempErr != nil {
		dto.SetErrorResponse(res, 403, "03", "Create Temp Image File", tempErr)

		return
	}

	contentType := handler.Header["Content-Type"][0]

	// 이미지 업로드 - minio

	_, uploadErr := database.UploadImage(handler.Filename, tempFile.Name(), contentType)

	if uploadErr != nil {
		dto.SetErrorResponse(res, 404, "04", "Upload Image Error", uploadErr)
		return
	}

	connect, _ := database.InitDatabaseConnection()

	var insertId int64

// 데이터 입력 - DB
	seq, insertErr := database.InsertQuery(
		connect, 
		queries.InsertPostImageData,
		// USER ID from JWT
		"1",
		userId,
		"post_table",
		"POST_IMAGE",
		strconv.Itoa(int(handler.Size)),
		handler.Filename, 
		contentType)
    
	if insertErr != nil {
 		dto.SetErrorResponse(res, 405, "05", "Insert Image Info Error", insertErr)
		return
    }

	insertId = seq

	defer connect.Close()

	removeErr := os.Remove(tempFile.Name())

	if removeErr != nil {
		log.Printf("[UPLOAD] Remove Saved Image Error: %v", removeErr)

		dto.SetErrorResponse(res, 406, "06", "Remove Image Error", removeErr)
		return
	}

	dto.SetFileInsertIdResponse(res, 200, "01", fmt.Sprintf("%d", insertId))
}