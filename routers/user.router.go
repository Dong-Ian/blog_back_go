package routers

import (
	"net/http"

	controllers "github.com/donghquinn/blog_back_go/controllers/users"
	"github.com/gorilla/mux"
)

func UserRouter(server *mux.Router) {
	server.HandleFunc("/user/signup", controllers.SignupController).Methods(http.MethodPost)
	server.HandleFunc("/user/login", controllers.LoginController).Methods(http.MethodPost)

	server.HandleFunc("/user/search/email", controllers.SearchEmailController).Methods(http.MethodPost)
	server.HandleFunc("/user/search/password", controllers.SearchPasswordController).Methods(http.MethodPost)

	server.HandleFunc("/user/profile", controllers.GetUserProfileController).Methods(http.MethodPost)
}
