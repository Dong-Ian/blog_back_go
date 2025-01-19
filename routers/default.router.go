package routers

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/controllers"
	"github.com/gorilla/mux"
)

func DefaultRouter(server *mux.Router) {
	server.HandleFunc("/api", controllers.DefaultController).Methods(http.MethodGet)

	server.HandleFunc("/", controllers.CorsTestController).Methods(http.MethodGet)

	server.HandleFunc("/refresh", controllers.RefreshController).Methods(http.MethodPost)

}
