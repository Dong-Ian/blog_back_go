package routers

import (
	"net/http"

	admincontrollers "github.com/donghquinn/blog_back_go/controllers/admin/posts"
	"github.com/gorilla/mux"
)

func AdminPostRouter(server *mux.Router) {
	server.HandleFunc("/admin/post/register", admincontrollers.RegisterPostController).Methods(http.MethodPost)
	server.HandleFunc("/admin/post/edit", admincontrollers.EditPostController).Methods(http.MethodPost)

	server.HandleFunc("/admin/post/delete", admincontrollers.DeletePostController).Methods(http.MethodPost)

	server.HandleFunc("/admin/post/update/pin", admincontrollers.UpdatePinPostController).Methods(http.MethodPost)
	server.HandleFunc("/admin/post/update/unpin", admincontrollers.UpdateUnPinPostController).Methods(http.MethodPost)
	server.HandleFunc("/admin/post/update/secret", admincontrollers.ChangeToSecretPostController).Methods(http.MethodPost)
	server.HandleFunc("/admin/post/update/unsecret", admincontrollers.ChangeToNotSecretPostController).Methods(http.MethodPost)
}
