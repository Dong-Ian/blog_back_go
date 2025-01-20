package routers

import (
	"net/http"

	admincontrollers "github.com/donghquinn/blog_back_go/controllers/admin/users"
	"github.com/gorilla/mux"
)

func AdminUserRouter(server *mux.Router) {
	server.HandleFunc("/admin/token/check", admincontrollers.CheckTokenController).Methods(http.MethodGet)
	server.HandleFunc("/admin/user/profile/update", admincontrollers.UpdateProfileController).Methods(http.MethodPost)
	server.HandleFunc("/admin/user/profile/title", admincontrollers.UpdateTitleController).Methods(http.MethodPost)
	server.HandleFunc("/admin/user/profile/color", admincontrollers.UpdateColorController).Methods(http.MethodPost)
}
