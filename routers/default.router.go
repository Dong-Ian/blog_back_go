package routers

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/controllers"
	postRouters "github.com/donghquinn/blog_back_go/routers/posts"
	userRouters "github.com/donghquinn/blog_back_go/routers/users"
	"github.com/gorilla/mux"
)

func DefaultRouter(server *mux.Router) {
	server.HandleFunc("/api", controllers.DefaultController).Methods(http.MethodGet)

	server.HandleFunc("/", controllers.CorsTestController).Methods(http.MethodGet)

	userRouters.UserRouter(server)
	postRouters.PostRouter(server)
}
