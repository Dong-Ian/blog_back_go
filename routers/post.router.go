package routers

import (
	"net/http"

	controllers "github.com/donghquinn/blog_back_go/controllers/posts"
	"github.com/gorilla/mux"
)

func PostRouter(server *mux.Router) {
	server.HandleFunc("/post/contents", controllers.PostContentsController).Methods(http.MethodPost)
	server.HandleFunc("/post/list", controllers.GetPostController).Methods(http.MethodPost)
	server.HandleFunc("/post/list/pinned", controllers.GetPinnedPostController).Methods(http.MethodPost)

	server.HandleFunc("/post/list/tag", controllers.GetPostsByTagController).Methods(http.MethodPost)
	server.HandleFunc("/post/list/category", controllers.GetPostsByCategoryController).Methods(http.MethodPost)

	server.HandleFunc("/post/url", controllers.GetImageUrl).Methods(http.MethodPost)

	server.HandleFunc("/post/category/list", controllers.GetCategoryController).Methods(http.MethodPost)
}
