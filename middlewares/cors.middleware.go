package middlewares

import (
	"log"
	"net/http"
)

func CorsMiddlewares(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Max-Age", "86400")
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		res.Header().Set("Access-Control-Allow-Credentials", "true")


		// Handle preflight request
	if req.Method == http.MethodOptions {
		log.Printf("[MIDDLEWARE] Handling preflight request")
		res.WriteHeader(http.StatusOK)
		return
	}

		next.ServeHTTP(res, req)
	})
}