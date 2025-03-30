package routers

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/biz/post"
	"github.com/donghquinn/blog_back_go/middlewares"
	"github.com/gorilla/mux"
)

func PostRouter(server *mux.Router) {
	server.HandleFunc("/post/{postSeq}", post.PostContentsController).Methods(http.MethodGet)
	server.HandleFunc("/post/list", post.GetPostController).Methods(http.MethodGet)

	server.HandleFunc("/post/url", post.GetImageUrl).Methods(http.MethodPost)

	server.HandleFunc("/post/category/list", post.GetCategoryController).Methods(http.MethodPost)
}

func PostAdminRouter(server *mux.Router) {
	sub := server.PathPrefix("/admin").Subrouter()
	sub.Use(middlewares.AuthMiddleware)

	sub.HandleFunc("/post/register", post.RegisterPostController).Methods(http.MethodPost)
	sub.HandleFunc("/post/edit", post.EditPostController).Methods(http.MethodPost)

	sub.HandleFunc("/post/delete", post.DeletePostController).Methods(http.MethodPost)

	sub.HandleFunc("/post/update/pin", post.UpdatePinPostController).Methods(http.MethodPost)
	sub.HandleFunc("/post/update/unpin", post.UpdateUnPinPostController).Methods(http.MethodPost)
	sub.HandleFunc("/post/update/secret", post.ChangeToSecretPostController).Methods(http.MethodPost)
	sub.HandleFunc("/post/update/unsecret", post.ChangeToNotSecretPostController).Methods(http.MethodPost)

}
