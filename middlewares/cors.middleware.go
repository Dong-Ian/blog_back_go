package middlewares

import (
	"log"
	"net/http"
)

var originList = []string{
	"http://localhost:3000",
	"https://blog.minjae-dev.com/",
	"https://blog.donghyuns.com",
	"unknown",
}

func CorsMiddlewares(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		log.Printf("Origin: %s", origin)

		// Origin 헤더가 없으면 기본 설정
		if origin == "" {
			origin = "unknown"
		}

		// 요청의 Origin이 허용된 Origin 목록에 있는지 확인
		for _, allowedOrigin := range originList {
			if allowedOrigin == origin {
				log.Printf("Allowed Origin: %s", origin)

				res.Header().Set("Access-Control-Allow-Origin", origin)
				res.Header().Set("Access-Control-Max-Age", "86400")
				res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
				res.Header().Set("Access-Control-Allow-Credentials", "true")
				res.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

				break
			}
		}

		// Handle preflight request
		if req.Method == http.MethodOptions {
			res.Header().Set("Access-Control-Allow-Origin", origin)
			res.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(res, req)
	})
}
