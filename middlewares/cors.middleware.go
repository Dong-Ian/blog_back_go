package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
)

// var originList = []string{
// 	"http://localhost:3000",
// 	"https://blog.minjae-dev.com/",
// 	"https://blog.donghyuns.com",
// 	"unknown",
// }

func CorsMiddlewares(next *mux.Router) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Max-Age", "86400")
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		res.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight request
		if req.Method == http.MethodOptions {
			res.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(res, req)
	})
}
