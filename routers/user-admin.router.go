package routers

import (
	"net/http"

	"github.com/donghquinn/blog_back_go/biz/user"
	"github.com/donghquinn/blog_back_go/middlewares"
	"github.com/gorilla/mux"
)

func UserRouter(server *mux.Router) {
	server.HandleFunc("/user/signup", user.SignupController).Methods(http.MethodPost)
	server.HandleFunc("/user/login", user.LoginController).Methods(http.MethodPost)

	server.HandleFunc("/user/search/email", user.SearchEmailController).Methods(http.MethodPost)
	server.HandleFunc("/user/search/password", user.SearchPasswordController).Methods(http.MethodPost)

	server.HandleFunc("/user/profile", user.GetUserProfileController).Methods(http.MethodPost)
}

func UserAdminRouter(server *mux.Router) {
	sub := server.PathPrefix("/admin").Subrouter()
	sub.Use(middlewares.AuthMiddleware)

	sub.HandleFunc("/token/check", user.CheckTokenController).Methods(http.MethodGet)
	sub.HandleFunc("/user/profile/update", user.UpdateProfileController).Methods(http.MethodPost)
	sub.HandleFunc("/user/profile/title", user.UpdateTitleController).Methods(http.MethodPost)
	sub.HandleFunc("/user/profile/color", user.UpdateColorController).Methods(http.MethodPost)
}
