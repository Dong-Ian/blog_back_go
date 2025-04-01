package routers

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/biz"
	"github.com/gorilla/mux"
)

func DefaultRouter(server *mux.Router) {
	server.HandleFunc("/api", biz.DefaultController).Methods(http.MethodGet)

	server.HandleFunc("/", biz.CorsTestController).Methods(http.MethodGet)

	server.HandleFunc("/refresh", biz.RefreshController).Methods(http.MethodPost)

}
