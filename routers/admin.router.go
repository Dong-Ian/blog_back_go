package routers

import (
	routers "github.com/donghquinn/blog_back_go/routers/admin"
	"github.com/gorilla/mux"
)

func AdminRouter(server *mux.Router) {
	routers.AdminPostRouter(server)
	routers.UploadImageController(server)
	routers.AdminUserRouter(server)
}
