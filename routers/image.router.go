package routers

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/biz/upload"
	"github.com/donghquinn/blog_back_go/middlewares"
	"github.com/gorilla/mux"
)

func UploadImageRouter(server *mux.Router) {
	sub := server.PathPrefix("/admin").Subrouter()
	sub.Use(middlewares.AuthMiddleware)

	sub.HandleFunc("/upload/image/profile", upload.UploadProfileImageController).Methods(http.MethodPost)
	sub.HandleFunc("/upload/image/background", upload.UploadBackgroundImageController).Methods(http.MethodPost)

	sub.HandleFunc("/upload/image/post", upload.UploadPostImageController).Methods(http.MethodPost)
}
