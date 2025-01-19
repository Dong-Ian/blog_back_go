package routers

import (
	"net/http"

	admincontrollers "github.com/donghquinn/blog_back_go/controllers/admin/upload"
	"github.com/gorilla/mux"
)

func UploadImageController(server *mux.Router) {
	server.HandleFunc("/admin/upload/image/profile", admincontrollers.UploadProfileImageController).Methods(http.MethodPost)
	server.HandleFunc("/admin/upload/image/background", admincontrollers.UploadBackgroundImageController).Methods(http.MethodPost)

	server.HandleFunc("/admin/upload/image/post", admincontrollers.UploadPostImageController).Methods(http.MethodPost)
}
