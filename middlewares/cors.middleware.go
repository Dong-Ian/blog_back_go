package middlewares

import (
	"log"
	"net/http"
)

var originList = []string{
	"http://localhost:3000",
	"https://blog.minjae-dev.com/",
	"https://blog.donghyuns.com",
}

func CorsMiddlewares(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		log.Printf("Origin: %s", origin)

		// 요청의 Origin이 허용된 Origin 목록에 있는지 확인
		isAllowed := false
		for _, allowedOrigin := range originList {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			res.Header().Set("Access-Control-Allow-Origin", origin)
		}

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
