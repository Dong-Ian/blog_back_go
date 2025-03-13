package upload

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// 이미지 파일 을 요청으로 부터 가져오기
func GetImagefileFromRequest(res http.ResponseWriter, req *http.Request) (multipart.File, *multipart.FileHeader, error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 500 MB files, 최대 10개.
    req.ParseMultipartForm(500 << 10)

	// FormFile returns the first file for the given key `file`
    // it also returns the FileHeader so we can get the Filename,
    // the Header and the size of the file
	file, handler, err := req.FormFile("file")

	if err != nil {
		log.Printf("[UPLOAD] Get File from Request Error: %v", err)
        return nil, nil, err
    }

	return file, handler, nil
}

// 이미지 파일 생성
func CreateFileImage(res http.ResponseWriter, req *http.Request, file multipart.File, handler *multipart.FileHeader) (*os.File,error) {
	// Create a temporary file within our temp-images directory that follows
    // a particular naming pattern
	fileFormat := strings.Split(handler.Header["Content-Type"][0], "/")[1]
    tempFile, err := os.CreateTemp("assets/image", fmt.Sprintf("upload-*.%s", fileFormat))

    if err != nil {
		log.Printf("[UPLOAD] Craete File Image Error: %v", err)
		return nil, err
    }

	// read all of the contents of our uploaded file into a
    // byte array
    fileBytes, err2 := io.ReadAll(file)
    if err2 != nil {
 
		return nil, err2
    }
	
    // write this byte array to our temporary file
    tempFile.Write(fileBytes)

    defer tempFile.Close()

	return tempFile, nil
}